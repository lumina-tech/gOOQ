package gooq

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type InsertSetStep interface {
	InsertSetMoreStep
	Select(s Selectable) InsertOnConflictStep
}

type InsertSetMoreStep interface {
	InsertValuesStep
	Set(f Field, v interface{}) InsertSetMoreStep
	Columns(f Field, rest ...Field) InsertValuesStep
}

type InsertValuesStep interface {
	InsertOnConflictStep
	Values(v ...interface{}) InsertValuesStep
}

type InsertOnConflictStep interface {
	InsertReturningStep
	OnConflictDoNothing() InsertReturningStep
	OnConflictDoUpdate(*DatabaseConstraint) InsertOnConflictSetStep
}

type InsertOnConflictSetStep interface {
	InsertFinalStep
	SetUpdates(f Field, v interface{}) InsertOnConflictSetStep
	SetUpdateColumns(f ...Field) InsertOnConflictSetStep
}

type InsertReturningStep interface {
	InsertFinalStep
	Returning(...Expression) InsertResultStep
}

type InsertResultStep interface {
	Fetchable
	Renderable
}

type InsertFinalStep interface {
	Executable
	Renderable
}

///////////////////////////////////////////////////////////////////////////////
// Implementation
///////////////////////////////////////////////////////////////////////////////

// https://www.postgresql.org/docs/current/sql-insert.html

type insert struct {
	table                 Table
	selection             Selectable
	columns               []Field
	values                [][]interface{}
	conflictAction        ConflictAction
	conflictConstraint    *DatabaseConstraint
	conflictSetPredicates []setPredicate
	returning             []Expression
}

func InsertInto(t Table) InsertSetStep {
	return &insert{table: t}
}

func (i *insert) Select(s Selectable) InsertOnConflictStep {
	i.selection = s
	return i
}

func (i *insert) Columns(f Field, rest ...Field) InsertValuesStep {
	i.columns = append([]Field{f}, rest...)
	return i
}

func (i *insert) Values(values ...interface{}) InsertValuesStep {
	i.values = append(i.values, values)
	return i
}

func (i *insert) Set(
	field Field, value interface{},
) InsertSetMoreStep {
	i.columns = append(i.columns, field)
	if len(i.values) == 0 {
		i.values = append(i.values, []interface{}{})
	}
	i.values[0] = append(i.values[0], value)
	return i
}

func (i *insert) OnConflictDoNothing() InsertReturningStep {
	i.conflictAction = ConflictActionDoNothing
	return i
}

func (i *insert) OnConflictDoUpdate(
	constraint *DatabaseConstraint,
) InsertOnConflictSetStep {
	i.conflictAction = ConflictActionDoUpdate
	i.conflictConstraint = constraint
	return i
}

func (i *insert) SetUpdates(
	field Field, value interface{},
) InsertOnConflictSetStep {
	i.conflictSetPredicates = append(i.conflictSetPredicates, setPredicate{field, value})
	return i
}

func (i *insert) SetUpdateColumns(
	fields ...Field,
) InsertOnConflictSetStep {
	// NOTE: excluded has to be lowercase
	excludedTable := NewTable("", "excluded")
	for _, field := range fields {
		i.conflictSetPredicates = append(i.conflictSetPredicates, setPredicate{
			field: field,
			value: NewStringField(excludedTable, field.GetName()),
		})
	}
	return i
}

func (i *insert) Returning(f ...Expression) InsertResultStep {
	i.returning = f
	return i
}

///////////////////////////////////////////////////////////////////////////////
// Executable
///////////////////////////////////////////////////////////////////////////////

func (i *insert) Exec(dl Dialect, db DBInterface) (sql.Result, error) {
	builder := i.Build(dl)
	return db.Exec(builder.String(), builder.arguments...)
}

func (i *insert) ExecWithContext(
	ctx context.Context, dl Dialect, db DBInterface) (sql.Result, error) {
	builder := i.Build(dl)
	return db.ExecContext(ctx, builder.String(), builder.arguments...)
}

///////////////////////////////////////////////////////////////////////////////
// Fetchable
///////////////////////////////////////////////////////////////////////////////

func (i *insert) Fetch(dl Dialect, db DBInterface) (*sqlx.Rows, error) {
	builder := i.Build(dl)
	return db.Queryx(builder.String(), builder.arguments...)
}

func (i *insert) FetchRow(dl Dialect, db DBInterface) *sqlx.Row {
	builder := i.Build(dl)
	return db.QueryRowx(builder.String(), builder.arguments...)
}

func (i *insert) FetchWithContext(
	ctx context.Context, dl Dialect, db DBInterface) (*sqlx.Rows, error) {
	builder := i.Build(dl)
	return db.QueryxContext(ctx, builder.String(), builder.arguments...)
}

func (i *insert) FetchRowWithContext(
	ctx context.Context, dl Dialect, db DBInterface) *sqlx.Row {
	builder := i.Build(dl)
	return db.QueryRowxContext(ctx, builder.String(), builder.arguments...)
}

///////////////////////////////////////////////////////////////////////////////
// Renderable
///////////////////////////////////////////////////////////////////////////////

func (i *insert) Build(d Dialect) *Builder {
	builder := Builder{}
	i.Render(&builder)
	return &builder
}

func (i *insert) Render(
	builder *Builder,
) {
	// INSERT INTO table_name
	builder.Printf("INSERT INTO %s ", i.table.GetQualifiedName())

	if i.selection != nil {
		// handle INSERT ...SELECT
		builder.Print("(")
		i.selection.Render(builder)
		builder.Print(")")
	} else {
		// handle INSERT .. SET
		i.renderColumnsAndValues(builder, i.columns, i.values)
	}

	// [ ON CONFLICT conflict_action ]
	if i.conflictAction != ConflictActionNil {
		builder.Printf(" ON CONFLICT")
		if i.conflictConstraint != nil {
			builder.Printf(" ON CONSTRAINT %s", i.conflictConstraint.Name)
		}
		if i.conflictAction == ConflictActionDoNothing {
			builder.Print(" DO NOTHING")
		} else if i.conflictAction == ConflictActionDoUpdate {
			builder.Printf(" %s SET ", i.conflictAction)
			builder.RenderSetPredicates(i.conflictSetPredicates)
		}
	}

	// [ RETURNING output_expression ]
	if i.returning != nil {
		builder.Print(" RETURNING ")
		builder.RenderExpressions(i.returning)
	}
}

// render set columns and values
func (i *insert) renderColumnsAndValues(
	builder *Builder, columns []Field, values [][]interface{},
) *Builder {
	if len(columns) > 0 {
		builder.Print("(")
		for index, column := range columns {
			builder.Printf(column.GetName())
			if index != len(columns)-1 {
				builder.Print(", ")
			}
		}
		builder.Printf(") ")
	}
	builder.Printf("VALUES ")
	for arrayIndex, array := range values {
		builder.Print("(")
		for index, value := range array {
			builder.RenderExpression(newLiteralExpression(value))
			if index != len(array)-1 {
				builder.Print(", ")
			}
		}
		builder.Print(")")
		if arrayIndex != len(values)-1 {
			builder.Print(", ")
		}
	}
	return builder
}
