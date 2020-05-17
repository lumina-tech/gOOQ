package gooq

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v3"
)

type Named interface {
	GetName() string
	GetQualifiedName() string
}

type Renderable interface {
	Render(builder *Builder)
}

type Selectable interface {
	Renderable
}

type DatabaseConstraint struct {
	Name      string
	Columns   []Field
	Predicate null.String
}

type DBInterface interface {
	sqlx.Execer
	sqlx.ExecerContext
	sqlx.Preparer
	sqlx.PreparerContext
	sqlx.Queryer
	sqlx.QueryerContext
}

type TxInterface interface {
	DBInterface
	Commit() error
	Rollback() error
}

type Executable interface {
	Exec(Dialect, DBInterface) (sql.Result, error)
	ExecWithContext(context.Context, Dialect, DBInterface) (sql.Result, error)
}

type Fetchable interface {
	Fetch(Dialect, DBInterface) (*sqlx.Rows, error)
	FetchRow(Dialect, DBInterface) *sqlx.Row
	FetchWithContext(context.Context, Dialect, DBInterface) (*sqlx.Rows, error)
	FetchRowWithContext(context.Context, Dialect, DBInterface) *sqlx.Row
}
