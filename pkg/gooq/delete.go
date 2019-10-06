package gooq

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type DeleteUsingStep interface {
	DeleteWhereStep
	Using(Selectable) DeleteOnStep
}

type DeleteWhereStep interface {
	DeleteResultStep
	Where(...Expression) DeleteReturningStep
}

type DeleteOnStep interface {
	DeleteResultStep
	On(...Expression) DeleteReturningStep
}

type DeleteReturningStep interface {
	DeleteFinalStep
	Returning(...Expression) DeleteResultStep
}

type DeleteResultStep interface {
	Fetchable
}

type DeleteFinalStep interface {
	Executable
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

func (d *deletion) On(c ...Expression) DeleteReturningStep {
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
	//return exec(dl, d, db)
	return nil, nil
}

///////////////////////////////////////////////////////////////////////////////
// Fetchable
///////////////////////////////////////////////////////////////////////////////

func (d *deletion) Fetch(dl Dialect, db DBInterface) (*sqlx.Rows, error) {
	//var buf bytes.Buffer
	//args := d.Render(dl, &buf)
	//return db.Queryx(buf.String(), args...)
	return nil, nil
}

func (d *deletion) FetchRow(dl Dialect, db DBInterface) (*sqlx.Row, error) {
	//var buf bytes.Buffer
	//args := d.Render(dl, &buf)
	//return db.QueryRowx(buf.String(), args...), nil
	return nil, nil
}

///////////////////////////////////////////////////////////////////////////////
// Renderable
///////////////////////////////////////////////////////////////////////////////

func (d *deletion) String(dl Dialect) string {
	builder := Builder{}
	d.Render(&builder)
	return builder.String()
}

// https://www.postgresql.org/docs/10/sql-delete.html
func (d *deletion) Render(
	builder *Builder,
) {

	// DELETE FROM table_name
	builder.Printf("DELETE FROM %s", d.table.GetName())

	if d.using != nil {
		// render USING clause
		builder.Printf(" USING ")
		switch sub := d.using.(type) {
		case Table:
			builder.Printf("%s ", sub.GetName())
		case *selection:
			builder.Printf("(")
			sub.Render(builder)
			builder.Printf(") AS %s ", sub.GetAlias().String)
		}
		// render ON clause
		if len(d.usingPredicate) > 0 {
			builder.RenderConditions(d.usingPredicate)
			builder.Print(" ON ")
			builder.RenderConditions(d.usingPredicate)
		}
	} else if len(d.conditions) > 0 {
		// [ WHERE condition ]
		builder.Print(" WHERE ")
		builder.RenderConditions(d.conditions)
	}

	// [ RETURNING output_expression ]
	if d.returning != nil {
		builder.Print(" RETURNING ")
		builder.RenderProjections(d.returning)
	}
}
