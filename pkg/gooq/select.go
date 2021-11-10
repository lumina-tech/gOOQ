package gooq

import (
	"context"
	"fmt"

	"gopkg.in/guregu/null.v3"

	"github.com/jmoiron/sqlx"
)

type JoinClause struct {
	Target     Selectable
	JoinType   JoinType
	Conditions []Expression
}

type SelectWithStep interface {
	SelectFromStep
	Select(projections ...Selectable) SelectFromStep
}

type SelectDistinctStep interface {
	SelectWhereStep
	Distinct() SelectFromStep
	DistinctOn(...Expression) SelectFromStep
	From(Selectable) SelectJoinStep
}

type SelectFromStep interface {
	SelectWhereStep
	From(Selectable) SelectJoinStep
}

type SelectJoinStep interface {
	SelectWhereStep
	Join(Selectable) SelectOnStep
	LeftOuterJoin(Selectable) SelectOnStep
}

type SelectOnStep interface {
	SelectWhereStep
	On(...Expression) SelectJoinStep
}

type SelectWhereStep interface {
	SelectGroupByStep
	Where(conditions ...Expression) SelectGroupByStep
}

type SelectGroupByStep interface {
	SelectHavingStep
	GroupBy(...Expression) SelectHavingStep
}

type SelectHavingStep interface {
	SelectOrderByStep
	Having(conditions ...Expression) SelectOrderByStep
}

type SelectOrderByStep interface {
	SelectOffsetStep
	OrderBy(...Expression) SelectOffsetStep
}

type SelectOffsetStep interface {
	SelectLimitStep
	Offset(offset int) SelectLimitStep
	Seek(v ...interface{}) SelectLimitStep
}

type SelectLimitStep interface {
	SelectFinalStep
	Limit(limit int) SelectFinalStep
}

type SelectFinalStep interface {
	Selectable
	Fetchable
	As(alias string) Selectable
	For(LockingType, LockingOption) SelectFinalStep
	Union(SelectFinalStep) SelectOrderByStep
}

///////////////////////////////////////////////////////////////////////////////
// Implementation
///////////////////////////////////////////////////////////////////////////////

type selection struct {
	with          Selectable
	withAlias     null.String
	selection     Selectable
	distinctOn    []Expression
	projections   []Selectable
	joins         []JoinClause
	joinTarget    Selectable
	joinType      JoinType
	predicate     []Expression
	groups        []Expression
	havings       []Expression
	ordering      []Expression
	unions        []SelectFinalStep
	alias         null.String
	isDistinct    bool
	limit         int
	offset        int
	seek          []interface{}
	lockingType   LockingType
	lockingOption LockingOption
}

func Select(projections ...Selectable) SelectDistinctStep {
	return &selection{projections: projections}
}

func SelectCount() SelectDistinctStep {
	return &selection{
		projections: []Selectable{Count(Asterisk)},
	}
}

func With(withAlias string, t Selectable) SelectWithStep {
	s := &selection{
		withAlias: null.StringFrom(withAlias),
		with:      t,
	}
	return s
}

func (s *selection) Select(projections ...Selectable) SelectFromStep {
	s.projections = projections
	return s
}

func (s *selection) Distinct() SelectFromStep {
	s.isDistinct = true
	return s
}

func (s *selection) DistinctOn(f ...Expression) SelectFromStep {
	s.distinctOn = f
	return s
}

func (s *selection) From(t Selectable) SelectJoinStep {
	s.selection = t
	return s
}

func (s *selection) Join(t Selectable) SelectOnStep {
	s.joinTarget = t
	s.joinType = Join
	return s
}

func (s *selection) LeftOuterJoin(t Selectable) SelectOnStep {
	// TODO copy and paste from From(.)
	s.joinTarget = t
	s.joinType = LeftOuterJoin
	return s
}

func (s *selection) Union(t SelectFinalStep) SelectOrderByStep {
	s.unions = append(s.unions, t)
	return s
}

func (s *selection) On(c ...Expression) SelectJoinStep {
	j := JoinClause{
		Target:     s.joinTarget,
		JoinType:   s.joinType,
		Conditions: c,
	}
	s.joinTarget = nil
	s.joinType = NotJoined
	s.joins = append(s.joins, j)
	return s
}

func (s *selection) As(alias string) Selectable {
	s.alias = null.StringFrom(alias)
	return s
}

func (s *selection) Where(c ...Expression) SelectGroupByStep {
	s.predicate = c
	return s
}

func (s *selection) GroupBy(f ...Expression) SelectHavingStep {
	s.groups = f
	return s
}

func (s *selection) Having(c ...Expression) SelectOrderByStep {
	s.havings = c
	return s
}

func (s *selection) OrderBy(f ...Expression) SelectOffsetStep {
	s.ordering = f
	return s
}

func (s *selection) Offset(offset int) SelectLimitStep {
	s.offset = offset
	return s
}

func (s *selection) Seek(v ...interface{}) SelectLimitStep {
	s.seek = v
	return s
}

func (s *selection) Limit(limit int) SelectFinalStep {
	s.limit = limit
	return s
}

func (s *selection) For(
	lockingType LockingType, lockingOption LockingOption,
) SelectFinalStep {
	s.lockingType = lockingType
	s.lockingOption = lockingOption
	return s
}

func (s *selection) GetAlias() null.String {
	return s.alias
}

///////////////////////////////////////////////////////////////////////////////
// Fetchable
///////////////////////////////////////////////////////////////////////////////

func (s *selection) Fetch(dl Dialect, db DBInterface) (*sqlx.Rows, error) {
	builder := s.Build(dl)
	return db.Queryx(builder.String(), builder.arguments...)
}

func (s *selection) FetchRow(dl Dialect, db DBInterface) *sqlx.Row {
	builder := s.Build(dl)
	return db.QueryRowx(builder.String(), builder.arguments...)
}

func (s *selection) FetchWithContext(
	ctx context.Context, dl Dialect, db DBInterface) (*sqlx.Rows, error) {
	builder := s.Build(dl)
	return db.QueryxContext(ctx, builder.String(), builder.arguments...)
}

func (s *selection) FetchRowWithContext(
	ctx context.Context, dl Dialect, db DBInterface) *sqlx.Row {
	builder := s.Build(dl)
	return db.QueryRowxContext(ctx, builder.String(), builder.arguments...)
}

///////////////////////////////////////////////////////////////////////////////
// Renderable
///////////////////////////////////////////////////////////////////////////////

func (s *selection) Build(d Dialect) *Builder {
	builder := Builder{}
	s.Render(&builder)
	return &builder
}

func (s *selection) Render(
	builder *Builder,
) {

	hasAlias := s.alias.Valid
	if hasAlias {
		builder.Print("(")
	}

	if s.withAlias.Valid {
		builder.Printf("WITH %s AS (", s.withAlias.String)
		s.with.Render(builder)
		builder.Print(") ")
	}

	builder.Print("SELECT ")

	if s.isDistinct {
		builder.Print("DISTINCT ")
	} else if len(s.distinctOn) > 0 {
		builder.Print("DISTINCT ON (")
		builder.RenderExpressions(s.distinctOn)
		builder.Print(") ")
	}

	projections := s.projections
	if len(projections) == 0 {
		projections = []Selectable{Asterisk}
	}
	// It is incorrect to always override projection namespace with selection alias.
	// The original this logic turns the following into
	// e.g. select item.*, foo.bar from (select * from boo) as item ...
	// e.g. select item.*, item.bar from (select * from boo) as item ...
	// colClause := renderProjections(alias, s.projection)
	builder.RenderProjections(projections)

	// render FROM clause
	if s.selection != nil {
		builder.Print(" FROM ")
		s.selection.Render(builder)
	}

	// render JOIN/ON clause
	for _, join := range s.joins {
		var joinString string
		switch join.JoinType {
		case LeftOuterJoin:
			joinString = "LEFT OUTER JOIN"
		case Join:
			joinString = "JOIN"
		}

		builder.Printf(" %s ", joinString)
		join.Target.Render(builder)
		builder.Print(" ON ")
		builder.RenderConditions(join.Conditions)
	}

	predicate := s.predicate
	if len(s.seek) > 0 {
		predicate = append(predicate, s.getSeekCondition())
	}

	// render WHERE clause
	if len(predicate) > 0 {
		builder.Print(" WHERE ")
		builder.RenderConditions(predicate)
	}

	// render GROUP BY clause
	if (len(s.groups)) > 0 {
		builder.Print(" GROUP BY ")
		builder.RenderExpressions(s.groups)
	}

	// render HAVING clause
	if len(s.havings) > 0 {
		builder.Print(" HAVING ")
		builder.RenderConditions(s.havings)
	}

	// render UNION clause
	for _, union := range s.unions {
		builder.Print(" UNION (")
		union.Render(builder)
		builder.Print(")")
	}

	// render ORDER BY clause
	if (len(s.ordering)) > 0 {
		builder.Print(" ORDER BY ")
		builder.RenderExpressions(s.ordering)
	}

	// render LIMIT clause
	if s.limit > 0 {
		builder.Printf(" LIMIT %d", s.limit)
	}

	// render OFFSET clause
	if s.offset > 0 {
		builder.Printf(" OFFSET %d", s.offset)
	}

	// render LOCKING clause
	if s.lockingType != LockingTypeNone {
		builder.Printf(" %s", s.lockingType.String())
		if s.lockingOption != LockingOptionNone {
			builder.Printf(" %s", s.lockingOption.String())
		}
	}

	if hasAlias {
		builder.Printf(") AS \"%s\"", s.alias.String)
	}

}

// faster and stable pagination based on these two articles
// https://blog.jooq.org/2013/10/26/faster-sql-paging-with-jooq-using-the-seek-method/
// https://blog.jooq.org/2013/11/18/faster-sql-pagination-with-keysets-continued/
// WARNING: seekAfter does not support seeking NULL values or the NULLS FIRST and
// NULL LAST clauses.
// e.g. Given the following scenario
// Select().From(Table1).
//   OrderBy(Table1.Column1.Desc(), Table1.Column2.Desc(), Table1.Column3.Desc()).
//   Seek("foo1", "foo2", "foo3"),
// We should generate the following where clause
// WHERE ((column1 < "foo1")
// OR (value1 = "foo1" AND value2 < "foo2")
// OR (value1 = "foo1" AND value2 = "foo2" AND value3 < "foo3"))
func (s *selection) getSeekCondition() Expression {
	if len(s.seek) < len(s.ordering) {
		panic("number of arguments in seek(...) must be gte number of arguments in orderBy")
	}
	// we went with the following approach to deal with mixed ordering
	var orExpressions []BoolExpression
	for i, order := range s.ordering {
		var operator Operator
		switch order.getOperator() {
		case OperatorDesc:
			operator = OperatorLt
		case OperatorAsc:
			operator = OperatorGt
		case OperatorNil:
			operator = OperatorGt
		default:
			panic(fmt.Sprintf("seek does not support operator=%s", order.getOperator()))
		}

		var andExpressions []BoolExpression
		for j := 0; j < i; j++ {
			expr := newBinaryBooleanExpressionImpl(
				OperatorEq, s.getOrderByField(s.ordering[j]), newLiteralExpression(s.seek[j]))
			andExpressions = append(andExpressions, expr)
		}
		expr := newBinaryBooleanExpressionImpl(
			operator, s.getOrderByField(order), newLiteralExpression(s.seek[i]))
		andExpressions = append(andExpressions, expr)
		orExpressions = append(orExpressions, And(andExpressions...))
	}
	return Or(orExpressions...)
}

func (s *selection) getOrderByField(
	order Expression,
) Expression {
	if order.getOperator() == OperatorNil {
		return order.getOriginal()
	}
	return order.getExpressions()[0].getOriginal()
}
