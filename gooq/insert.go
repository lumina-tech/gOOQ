package gooq

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ConflictAction string

const (
	ConflictActionDoNothing ConflictAction = "DO NOTHING"
	ConflictActionDoUpdate  ConflictAction = "DO UPDATE"
)

type insert struct {
	table              TableLike
	selection          *selection
	bindings           []TableFieldBinding
	conflictAction     *ConflictAction
	conflictConstraint *DatabaseConstraint
	conflictUpdate     []TableFieldBinding
	returning          TableField
}

func InsertInto(t TableLike) InsertSetStep {
	return &insert{table: t}
}

func (i *insert) Select(s Selectable) InsertOnConflictStep {
	selection, ok := s.(*selection)
	if ok {
		i.selection = selection
	} else {
		// TODO(Peter): log warning
	}
	return i
}

func (i *insert) Set(f TableField, v interface{}) InsertSetMoreStep {
	binding := TableFieldBinding{Field: f, Value: v}
	i.bindings = append(i.bindings, binding)
	return i
}

func (i *insert) OnConflictDoNothing() InsertReturningStep {
	action := ConflictActionDoNothing
	i.conflictAction = &action
	return i
}

func (i *insert) OnConflictDoUpdate(constraint *DatabaseConstraint, updateFields ...TableFieldBinding) InsertReturningStep {
	action := ConflictActionDoUpdate
	i.conflictAction = &action
	i.conflictConstraint = constraint
	i.conflictUpdate = updateFields
	return i
}

func (i *insert) Returning(f TableField) InsertResultStep {
	i.returning = f
	return i
}

///////////////////////////////////////////////////////////////////////////////
// Executable
///////////////////////////////////////////////////////////////////////////////

func (i *insert) Exec(d Dialect, db DBInterface) (sql.Result, error) {
	return exec(d, i, db)
}

///////////////////////////////////////////////////////////////////////////////
// Fetchable
///////////////////////////////////////////////////////////////////////////////

func (i *insert) Fetch(dl Dialect, db DBInterface) (*sqlx.Rows, error) {
	var buf bytes.Buffer
	args := i.Render(dl, &buf)
	return db.Queryx(buf.String(), args...)
}

func (i *insert) FetchRow(dl Dialect, db DBInterface) (*sqlx.Row, error) {
	var buf bytes.Buffer
	args := i.Render(dl, &buf)
	return db.QueryRowx(buf.String(), args...), nil
}

///////////////////////////////////////////////////////////////////////////////
// Renderable
///////////////////////////////////////////////////////////////////////////////

func (i *insert) String(d Dialect) string {
	return toString(d, i)
}

func (i *insert) Render(
	dl Dialect, w io.Writer,
) (placeholders []interface{}) {
	placeholders = []interface{}{}
	return i.RenderWithPlaceholders(dl, w, placeholders)
}

// https://www.postgresql.org/docs/current/sql-insert.html
func (i *insert) RenderWithPlaceholders(
	d Dialect, w io.Writer, placeholders []interface{},
) []interface{} {

	// INSERT INTO table_name
	fmt.Fprintf(w, "INSERT INTO %s ", i.table.Name())

	// { VALUES | query }
	if i.selection != nil {
		// render query
		fmt.Fprint(w, "(")
		selectPlaceholders := i.selection.RenderWithPlaceholders(d, w, placeholders)
		placeholders = append(placeholders, selectPlaceholders...)
		fmt.Fprint(w, ")")
	} else {
		// render values clause
		columnFragments := make([]string, len(i.bindings))
		valueFragments := make([]string, len(i.bindings))
		for i, binding := range i.bindings {
			columnFragments[i] = binding.Field.Name()
			valueFragments[i] = d.renderPlaceholder(i + 1)
			placeholders = append(placeholders, binding.Value)
		}
		columnClause := strings.Join(columnFragments, ", ")
		valuesClause := strings.Join(valueFragments, ",")
		fmt.Fprintf(w, "(%s) VALUES (%s)", columnClause, valuesClause)
	}

	// [ ON CONFLICT conflict_action ]
	if i.conflictAction != nil {
		fmt.Fprint(w, " ON CONFLICT ")
		if i.conflictConstraint != nil {
			fmt.Fprintf(w, "( %s )", strings.Join(i.conflictConstraint.Columns, ", "))
		}
		fmt.Fprint(w, string(*i.conflictAction))
		if *i.conflictAction == ConflictActionDoUpdate {
			var setClauses []string
			if len(i.conflictUpdate) > 0 {
				for _, field := range i.conflictUpdate {
					setClauses = append(setClauses, fmt.Sprintf("%s = %s", field.Field.Name(), field.Value))
				}
			} else {
				for _, field := range i.bindings {
					setClauses = append(setClauses, fmt.Sprintf("%s = excluded.%s", field.Field.Name(), field.Field.Name()))
				}
			}
			fmt.Fprintf(w, " SET %s ", strings.Join(setClauses, ","))
		}
	}

	// [ RETURNING output_expression ]
	if i.returning != nil {
		fmt.Fprint(w, " RETURNING ")
		fmt.Fprint(w, i.returning.Name())
	}

	return placeholders
}
