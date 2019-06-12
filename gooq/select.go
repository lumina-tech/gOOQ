package gooq

import (
	"bytes"
	"fmt"
	"io"

	"github.com/0x6e6562/gosnow"
	"github.com/jmoiron/sqlx"
)

var flake, _ = gosnow.Default()

type selection struct {
	selection  Selectable
	projection []Field
	joins      []join
	joinTarget Selectable
	joinType   JoinType
	predicate  []Condition
	groups     []Field
	havings    []Condition
	ordering   []Field
	unions     []SelectFinalStep
	count      bool
	alias      string
	limit      int
	offset     int
}

func Select(f ...Field) SelectFromStep {
	return &selection{projection: f}
}

func SelectCount() SelectFromStep {
	return &selection{count: true}
}

func (sl *selection) From(s Selectable) SelectJoinStep {
	sl.selection = s
	return sl
}

func (s *selection) Join(t Selectable) SelectOnStep {
	s.joinTarget = t
	s.joinType = Join
	return s
}

func (s *selection) LeftOuterJoin(t Selectable) SelectOnStep {
	// TODO copy and paste from Join(.)
	s.joinTarget = t
	s.joinType = LeftOuterJoin
	return s
}

func (s *selection) Union(t SelectFinalStep) SelectOrderByStep {
	s.unions = append(s.unions, t)
	return s
}

func (s *selection) On(c ...JoinCondition) SelectJoinStep {
	j := join{
		target:   s.joinTarget,
		joinType: s.joinType,
		conds:    c,
	}
	s.joinTarget = nil
	s.joinType = NotJoined
	s.joins = append(s.joins, j)
	return s
}

func (s *selection) As(alias string) SelectFinalStep {
	s.alias = alias
	return s
}

func (s *selection) Where(c ...Condition) SelectGroupByStep {
	s.predicate = c
	return s
}

func (s *selection) GroupBy(f ...Field) SelectHavingStep {
	s.groups = f
	return s
}

func (s *selection) Having(c ...Condition) SelectOrderByStep {
	s.havings = c
	return s
}

func (s *selection) OrderBy(f ...Field) SelectLimitStep {
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

func (s *selection) IsSelectable() {}

func (s *selection) Alias() string {
	return s.alias
}

func (s *selection) MaybeAlias() string {
	if s.alias == "" {
		switch sub := s.selection.(type) {
		case TableLike:
			return sub.Name()
		default:
			return ""
		}
	} else {
		return s.alias
	}
}

///////////////////////////////////////////////////////////////////////////////
// Fetchable
///////////////////////////////////////////////////////////////////////////////

func (s *selection) Fetch(d Dialect, db DBInterface) (*sqlx.Rows, error) {
	var buf bytes.Buffer
	args := s.Render(d, &buf)
	return db.Queryx(buf.String(), args...)
}

func (s *selection) FetchRow(d Dialect, db DBInterface) (*sqlx.Row, error) {
	var buf bytes.Buffer
	args := s.Render(d, &buf)
	return db.QueryRowx(buf.String(), args...), nil
}

///////////////////////////////////////////////////////////////////////////////
// Renderable
///////////////////////////////////////////////////////////////////////////////

func (s *selection) String(d Dialect) string {
	return toString(d, s)
}

func (s *selection) Render(d Dialect, w io.Writer) (placeholders []interface{}) {
	placeholders = []interface{}{}
	return s.RenderWithPlaceholders(d, w, placeholders)
}

func (s *selection) RenderWithPlaceholders(
	d Dialect, w io.Writer, placeholders []interface{},
) []interface{} {

	alias := ""
	if al, ok := s.selection.(Aliasable); ok {
		if al.Alias() != "" {
			alias = al.Alias()
		}
	}

	fmt.Fprint(w, "SELECT ")

	if s.count {
		fmt.Fprint(w, "COUNT(*)")
	} else {
		if len(s.projection) == 0 {
			fmt.Fprint(w, "*")
		} else {
			// It is incorrect to always override projection namespace with selection alias.
			// The original this logic turns the following into
			// e.g. select item.*, foo.bar from (select * from boo) as item ...
			// e.g. select item.*, item.bar from (select * from boo) as item ...
			// colClause := renderProjections(alias, s.projection)
			placeholders = s.renderProjections(w, d, "", s.projection, placeholders)
		}
	}

	// render FROM clause
	fmt.Fprintf(w, " FROM ")
	switch sub := s.selection.(type) {
	case TableLike:
		fmt.Fprint(w, sub.Name())
	case *selection:
		fmt.Fprint(w, "(")
		placeholders = sub.RenderWithPlaceholders(d, w, placeholders)
		fmt.Fprint(w, ")")
	}

	if alias != "" {
		fmt.Fprintf(w, " AS %s", alias)
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

		fmt.Fprintf(w, " %s ", joinString)
		switch sub := join.target.(type) {
		case TableLike:
			fmt.Fprint(w, renderTable(sub))
		case *selection:
			fmt.Fprint(w, "(")
			placeholders = sub.RenderWithPlaceholders(d, w, placeholders)
			fmt.Fprintf(w, ") AS %s", sub.Alias())
		}
		fmt.Fprintf(w, " ON %s", renderJoinConditions(join.conds))
	}

	// render WHERE clause
	if len(s.predicate) > 0 {
		newPlaceholders := renderWhereClause(
			alias, s.predicate, d, len(placeholders), w)
		placeholders = append(placeholders, newPlaceholders...)
	}

	// render GROUP BY clause
	if (len(s.groups)) > 0 {
		fmt.Fprint(w, " GROUP BY ")
		placeholders = s.renderFields(w, d, alias, s.groups, placeholders)
	}

	// render HAVINGS clause
	if len(s.havings) > 0 {
		panic("having clause is not implemented")
	}

	// render UNION clause
	for _, union := range s.unions {
		fmt.Fprintf(w, " UNION (")
		placeholders = union.RenderWithPlaceholders(d, w, placeholders)
		fmt.Fprintf(w, ")")
	}

	// render ORDER BY clause
	if (len(s.ordering)) > 0 {
		fmt.Fprint(w, " ORDER BY ")
		placeholders = s.renderFields(w, d, alias, s.ordering, placeholders)
	}

	// render LIMIT clause
	if s.limit > 0 {
		fmt.Fprintf(w, " LIMIT %d", s.limit)
	}

	// render OFFSET clause
	if s.offset > 0 {
		fmt.Fprintf(w, " OFFSET %d", s.offset)
	}

	return placeholders
}

func renderTable(t TableLike) string {
	if t.Alias() != "" {
		return fmt.Sprintf("%s AS %s", t.Name(), t.Alias())
	}
	return t.Name()
}

// renderProjections - renders project and always render the field expressions
// e.g. time_bucket('5 minutes', table1.creation_date) AS five_min
func (sl *selection) renderProjections(
	w io.Writer, d Dialect, alias string,
	cols []Field, placeholders []interface{},
) []interface{} {
	for i, col := range cols {
		placeholders = sl.renderFieldExpression(d, w, alias, col, placeholders)
		if col.Alias() != "" {
			fmt.Fprint(w, " AS "+col.Alias())
		}
		if i != len(cols)-1 {
			fmt.Fprint(w, ", ")
		}
	}
	return placeholders
}

// renderFields - renders expressions in GROUP BY, ORDER BY etc... it prioritize
// alias over expression. Given an aliased field, renderFields will render the
// alias instead of the expression e.g.
// time_bucket('5 minutes', table1.creation_date) AS five_min -> five_min
func (sl *selection) renderFields(
	w io.Writer, d Dialect, alias string,
	cols []Field, placeholders []interface{},
) []interface{} {
	for i, col := range cols {
		if col.Alias() != "" {
			fmt.Fprint(w, col.Alias())
		} else {
			placeholders = sl.renderFieldExpression(d, w, alias, col, placeholders)
		}
		if i != len(cols)-1 {
			fmt.Fprint(w, ", ")
		}
	}
	return placeholders
}

func (sl *selection) renderFieldExpression(
	d Dialect, w io.Writer, alias string, col Field, placeholders []interface{},
) []interface{} {
	al := resolveParentAlias(alias, col)
	aliased := col.Name()
	if al != "" {
		aliased = fmt.Sprintf("%s.%s", al, aliased)
	}
	var f string
	f = renderFunction(aliased, col.Function())
	fmt.Fprint(w, f)
	filters := col.Filters()
	if filters != nil {
		placeholders = append(placeholders, renderFilter(w, d, filters, len(placeholders))...)
	}
	return placeholders
}
