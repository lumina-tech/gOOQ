package gooq

import (
	"gopkg.in/guregu/null.v3"

	"github.com/jmoiron/sqlx"
)

type join struct {
	target     Selectable
	joinType   JoinType
	conditions []Expression
}

type SelectDistinctStep interface {
	SelectWhereStep
	Distinct() SelectFromStep
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
	SelectLimitStep
	OrderBy(...Expression) SelectLimitStep
}

type SelectLimitStep interface {
	SelectOffsetStep
	Limit(limit int) SelectOffsetStep
}

type SelectOffsetStep interface {
	SelectFinalStep
	Offset(offset int) SelectFinalStep
}

type SelectFinalStep interface {
	Selectable
	Fetchable
	As(alias string) SelectFinalStep
	Union(SelectFinalStep) SelectOrderByStep
}

///////////////////////////////////////////////////////////////////////////////
// Implementation
///////////////////////////////////////////////////////////////////////////////

type selection struct {
	selection   Selectable
	projections []Expression
	joins       []join
	joinTarget  Selectable
	joinType    JoinType
	predicate   []Expression
	groups      []Expression
	havings     []Expression
	ordering    []Expression
	unions      []SelectFinalStep
	alias       null.String
	isDistinct  bool
	limit       int
	offset      int
}

func Select(projections ...Expression) SelectDistinctStep {
	return &selection{projections: projections}
}

func SelectCount() SelectDistinctStep {
	return &selection{
		projections: []Expression{Count(Asterisk)},
	}
}

func (s *selection) Distinct() SelectFromStep {
	s.isDistinct = true
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
	j := join{
		target:     s.joinTarget,
		joinType:   s.joinType,
		conditions: c,
	}
	s.joinTarget = nil
	s.joinType = NotJoined
	s.joins = append(s.joins, j)
	return s
}

func (s *selection) As(alias string) SelectFinalStep {
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

func (s *selection) OrderBy(f ...Expression) SelectLimitStep {
	s.ordering = f
	return s
}

func (s *selection) Limit(limit int) SelectOffsetStep {
	s.limit = limit
	return s
}

func (s *selection) Offset(offset int) SelectFinalStep {
	s.offset = offset
	return s
}

///////////////////////////////////////////////////////////////////////////////
// Aliasable
///////////////////////////////////////////////////////////////////////////////

func (s *selection) GetAlias() null.String {
	return s.alias
}

///////////////////////////////////////////////////////////////////////////////
// Fetchable
///////////////////////////////////////////////////////////////////////////////

func (s *selection) Fetch(d Dialect, db DBInterface) (*sqlx.Rows, error) {
	return db.Queryx(s.String(d))
}

func (s *selection) FetchRow(d Dialect, db DBInterface) (*sqlx.Row, error) {
	return db.QueryRowx(s.String(d)), nil
}

///////////////////////////////////////////////////////////////////////////////
// Renderable
///////////////////////////////////////////////////////////////////////////////

func (s *selection) String(d Dialect) string {
	builder := Builder{}
	s.Render(&builder)
	return builder.String()
}

func (s *selection) Render(
	builder *Builder,
) {

	builder.Print("SELECT ")

	if s.isDistinct {
		builder.Print("DISTINCT ")
	}

	projections := s.projections
	if len(projections) == 0 {
		projections = []Expression{Asterisk}
	}
	// It is incorrect to always override projection namespace with selection alias.
	// The original this logic turns the following into
	// e.g. select item.*, foo.bar from (select * from boo) as item ...
	// e.g. select item.*, item.bar from (select * from boo) as item ...
	// colClause := renderProjections(alias, s.projection)
	builder.RenderProjections(projections)

	// render FROM clause
	builder.Print(" FROM ")
	switch sub := s.selection.(type) {
	case Table:
		builder.Print(sub.GetName())
	case *selection:
		builder.Print("(")
		sub.Render(builder)
		builder.Print(")")
	}

	alias := s.selection.GetAlias()
	if alias.Valid {
		builder.Printf(" AS %s", alias.String)
	}

	// render JOIN/ON clause
	for _, join := range s.joins {
		var joinString string
		switch join.joinType {
		case LeftOuterJoin:
			joinString = "LEFT OUTER JOIN"
		case Join:
			joinString = "JOIN"
		}

		builder.Printf(" %s ", joinString)
		switch sub := join.target.(type) {
		case Table:
			builder.Print(sub.GetAliasOrName())
		case *selection:
			builder.Print("(")
			sub.Render(builder)
			builder.Printf(") AS %s", sub.GetAlias().String)
		}
		builder.Print(" ON ")
		builder.RenderConditions(join.conditions)
	}

	// render WHERE clause
	if len(s.predicate) > 0 {
		builder.Print(" WHERE ")
		builder.RenderConditions(s.predicate)
	}

	// render GROUP BY clause
	if (len(s.groups)) > 0 {
		builder.Print(" GROUP BY ")
		builder.RenderFields(s.groups)
	}

	// render HAVINGS clause
	if len(s.havings) > 0 {
		panic("having clause is not implemented")
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
		builder.RenderFields(s.ordering)
	}

	// render LIMIT clause
	if s.limit > 0 {
		builder.Printf(" LIMIT %d", s.limit)
	}

	// render OFFSET clause
	if s.offset > 0 {
		builder.Printf(" OFFSET %d", s.offset)
	}
}
