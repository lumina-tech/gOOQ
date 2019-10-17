package gooq

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type DeleteUsingStep interface {
	DeleteWhereStep
	Using(Selectable) DeleteOnStep
}

type DeleteOnStep interface {
	DeleteWhereStep
	On(...Expression) DeleteWhereStep
}

type DeleteWhereStep interface {
	DeleteResultStep
	Where(...Expression) DeleteReturningStep
}

type DeleteReturningStep interface {
	DeleteFinalStep
	Returning(...Expression) DeleteResultStep
}

type DeleteResultStep interface {
	Fetchable
	Renderable
}

type DeleteFinalStep interface {
	Executable
	Renderable
}

///////////////////////////////////////////////////////////////////////////////
// Implementation
///////////////////////////////////////////////////////////////////////////////

// https://www.postgresql.org/docs/11/sql-delete.html

type deletion struct {
	table          Table
	using          Selectable
	conditions     []Expression
	usingPredicate []Expression // where clause for using
	returning      []Expression
}

func Delete(t Table) DeleteUsingStep {
	return &deletion{table: t}
}

func (d *deletion) Using(s Selectable) DeleteOnStep {
	d.using = s
	return d
}

func (d *deletion) On(c ...Expression) DeleteWhereStep {
	d.usingPredicate = c
	return d
}

func (d *deletion) Where(c ...Expression) DeleteReturningStep {
	d.conditions = c
	return d
}

func (d *deletion) Returning(f ...Expression) DeleteResultStep {
	d.returning = f
	return d
}

///////////////////////////////////////////////////////////////////////////////
// Executable
///////////////////////////////////////////////////////////////////////////////

func (d *deletion) Exec(dl Dialect, db DBInterface) (sql.Result, error) {
	builder := d.Build(dl)
	return db.Exec(builder.String(), builder.arguments...)
}

///////////////////////////////////////////////////////////////////////////////
// Fetchable
///////////////////////////////////////////////////////////////////////////////

func (d *deletion) Fetch(dl Dialect, db DBInterface) (*sqlx.Rows, error) {
	builder := d.Build(dl)
	return db.Queryx(builder.String(), builder.arguments...)
}

func (d *deletion) FetchRow(dl Dialect, db DBInterface) *sqlx.Row {
	builder := d.Build(dl)
	return db.QueryRowx(builder.String(), builder.arguments...)
}

///////////////////////////////////////////////////////////////////////////////
// Renderable
///////////////////////////////////////////////////////////////////////////////

func (d *deletion) Build(dl Dialect) *Builder {
	builder := Builder{}
	d.Render(&builder)
	return &builder
}

// https://www.postgresql.org/docs/10/sql-delete.html
func (d *deletion) Render(
	builder *Builder,
) {

	// DELETE FROM table_name
	builder.Printf("DELETE FROM %s", d.table.GetQualifiedName())

	conditions := d.conditions
	if d.using != nil {
		// render USING clause
		builder.Printf(" USING ")
		d.using.Render(builder)
		// there is no "ON" clause in postgres, this is pattern is from jOOQ.
		// https://www.jooq.org/doc/3.12/manual-single-page/#delete-statement
		conditions = append(d.usingPredicate, conditions...)
	}

	if len(conditions) > 0 {
		// [ WHERE condition ]
		builder.Print(" WHERE ")
		builder.RenderConditions(conditions)
	}

	// [ RETURNING output_expression ]
	if d.returning != nil {
		builder.Print(" RETURNING ")
		builder.RenderExpressions(d.returning)
	}
}
