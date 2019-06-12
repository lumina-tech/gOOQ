package gooq

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"

	"github.com/jmoiron/sqlx"
)

type deletion struct {
	table          TableLike
	using          Selectable
	predicate      []Condition
	usingPredicate []JoinCondition // where clause for using
	returning      TableField
}

func Delete(t TableLike) DeleteUsingStep {
	return &deletion{table: t}
}

func (d *deletion) Using(s Selectable) DeleteOnStep {
	d.using = s
	return d
}

func (d *deletion) On(c ...JoinCondition) DeleteReturningStep {
	d.usingPredicate = c
	return d
}

func (d *deletion) Where(c ...Condition) DeleteReturningStep {
	d.predicate = c
	return d
}

func (d *deletion) Returning(f TableField) DeleteResultStep {
	d.returning = f
	return d
}

///////////////////////////////////////////////////////////////////////////////
// Executable
///////////////////////////////////////////////////////////////////////////////

func (d *deletion) Exec(dl Dialect, db DBInterface) (sql.Result, error) {
	return exec(dl, d, db)
}

///////////////////////////////////////////////////////////////////////////////
// Fetchable
///////////////////////////////////////////////////////////////////////////////

func (d *deletion) Fetch(dl Dialect, db DBInterface) (*sqlx.Rows, error) {
	var buf bytes.Buffer
	args := d.Render(dl, &buf)
	return db.Queryx(buf.String(), args...)
}

func (d *deletion) FetchRow(dl Dialect, db DBInterface) (*sqlx.Row, error) {
	var buf bytes.Buffer
	args := d.Render(dl, &buf)
	return db.QueryRowx(buf.String(), args...), nil
}

///////////////////////////////////////////////////////////////////////////////
// Renderable
///////////////////////////////////////////////////////////////////////////////

func (d *deletion) String(dl Dialect) string {
	return toString(dl, d)
}

func (d *deletion) Render(
	dl Dialect, w io.Writer,
) (placeholders []interface{}) {
	placeholders = []interface{}{}
	return d.RenderWithPlaceholders(dl, w, placeholders)
}

// https://www.postgresql.org/docs/10/sql-delete.html
func (d *deletion) RenderWithPlaceholders(
	dl Dialect, w io.Writer, placeholders []interface{},
) []interface{} {

	// DELETE FROM table_name
	fmt.Fprintf(w, "DELETE FROM %s", d.table.Name())

	if d.using != nil {
		// render FROM clause
		fmt.Fprint(w, " USING ")
		switch sub := d.using.(type) {
		case TableLike:
			fmt.Fprint(w, sub.Name())
		case *selection:
			fmt.Fprint(w, "(")
			placeholders = sub.RenderWithPlaceholders(dl, w, placeholders)
			fmt.Fprintf(w, ") AS %s ", sub.Alias())
		}
		// render WHERE clause
		if len(d.usingPredicate) > 0 {
			joinWhereClause := renderJoinConditions(d.usingPredicate)
			fmt.Fprintf(w, "WHERE %s", joinWhereClause)
		}
	} else {
		// [ WHERE condition ]
		if len(d.predicate) > 0 {
			placeholdersParams := renderWhereClause(
				d.table.Name(), d.predicate, dl, len(placeholders), w)
			placeholders = append(placeholders, placeholdersParams...)
		}
	}

	// [ RETURNING output_expression ]
	if d.returning != nil {
		fmt.Fprint(w, " RETURNING ")
		fmt.Fprint(w, d.returning.Name())
	}
	return placeholders
}
