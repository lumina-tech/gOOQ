package gooq

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"strings"

	"github.com/jmoiron/sqlx"
)

type update struct {
	table          TableLike
	selection      Selectable
	bindings       []TableFieldBinding // set clause
	predicate      []Condition         // where clause
	joinPredicate  []JoinCondition     // where clause for join
	conflictAction *ConflictAction
	returning      TableField
}

func Update(t TableLike) UpdateSetStep {
	return &update{table: t}
}

func (u *update) Set(f TableField, v interface{}) UpdateSetStep {
	binding := TableFieldBinding{Field: f, Value: v}
	u.bindings = append(u.bindings, binding)
	return u
}

func (u *update) Join(s Selectable) UpdateOnStep {
	u.selection = s
	return u
}

func (u *update) On(c ...JoinCondition) UpdateOnConflictStep {
	u.joinPredicate = c
	return u
}

func (u *update) Where(c ...Condition) UpdateOnConflictStep {
	u.predicate = c
	return u
}

func (u *update) OnConflict(action ConflictAction) UpdateReturningStep {
	u.conflictAction = &action
	return u
}

func (u *update) Returning(f TableField) UpdateResultStep {
	u.returning = f
	return u
}

///////////////////////////////////////////////////////////////////////////////
// Executable
///////////////////////////////////////////////////////////////////////////////

func (u *update) Exec(d Dialect, db DBInterface) (sql.Result, error) {
	return exec(d, u, db)
}

///////////////////////////////////////////////////////////////////////////////
// Fetchable
///////////////////////////////////////////////////////////////////////////////

func (u *update) Fetch(d Dialect, db DBInterface) (*sqlx.Rows, error) {
	var buf bytes.Buffer
	args := u.Render(d, &buf)
	return db.Queryx(buf.String(), args...)
}

func (u *update) FetchRow(d Dialect, db DBInterface) (*sqlx.Row, error) {
	var buf bytes.Buffer
	args := u.Render(d, &buf)
	return db.QueryRowx(buf.String(), args...), nil
}

///////////////////////////////////////////////////////////////////////////////
// Renderable
///////////////////////////////////////////////////////////////////////////////

func (u *update) String(d Dialect) string {
	return toString(d, u)
}

func (u *update) Render(
	dl Dialect, w io.Writer,
) (placeholders []interface{}) {
	placeholders = []interface{}{}
	return u.RenderWithPlaceholders(dl, w, placeholders)
}

// https://www.postgresql.org/docs/10/sql-update.html
func (u *update) RenderWithPlaceholders(
	d Dialect, w io.Writer, placeholders []interface{},
) []interface{} {

	// UPDATE table_name SET
	fmt.Fprintf(w, "UPDATE %s SET ", u.table.Name())

	// In the case of UPDATE without JOIN
	if u.selection == nil {
		// render SET clause
		setFragments := make([]string, len(u.bindings))
		for i, binding := range u.bindings {
			col := binding.Field.Name()
			setFragments[i] = fmt.Sprintf("%s = %s", col, d.renderPlaceholder(i+1))
			placeholders = append(placeholders, binding.Value)
		}
		setClause := strings.Join(setFragments, ", ")
		fmt.Fprintf(w, setClause)

		// render WHERE clause
		whereValues := renderWhereClause(
			u.table.Name(), u.predicate, d, len(placeholders), w)
		placeholders = append(placeholders, whereValues...)

	} else {
		// render SET clause
		setFragments := make([]string, len(u.bindings))
		for i, binding := range u.bindings {
			col := binding.Field.Name()
			switch value := binding.Value.(type) {
			case TableField:
				rhsAlias, _ := renderFieldAlias(value.Parent().MaybeAlias(), value)
				setFragments[i] = fmt.Sprintf("%s = %s", col, rhsAlias)
			default:
				placeholderPos := len(placeholders) + 1
				setFragments[i] = fmt.Sprintf("%s = %s", col, d.renderPlaceholder(placeholderPos))
				placeholders = append(placeholders, value)
			}
		}
		setClause := strings.Join(setFragments, ", ")
		fmt.Fprintf(w, setClause)

		// render FROM clause
		fmt.Fprint(w, " FROM ")
		switch sub := u.selection.(type) {
		case TableLike:
			fmt.Fprint(w, sub.Name())
		case *selection:
			fmt.Fprint(w, "(")
			placeholders = sub.RenderWithPlaceholders(d, w, placeholders)
			// postgres requires that subquery in FROM must have an alias
			fmt.Fprintf(w, ") AS %s", sub.Alias())
		}

		// render WHERE clause
		if len(u.joinPredicate) > 0 {
			joinWhereClause := renderJoinConditions(u.joinPredicate)
			fmt.Fprintf(w, " WHERE %s", joinWhereClause)
		}
	}

	// render on conflict
	if u.conflictAction != nil {
		fmt.Fprintf(w, " ON CONFLICT %s", string(*u.conflictAction))
	}

	// render returning
	if u.returning != nil {
		fmt.Fprint(w, " RETURNING ")
		fmt.Fprint(w, u.returning.Name())
	}

	return placeholders
}
