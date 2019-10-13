package gooq

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v3"
)

type DatabaseConstraint struct {
	Name      string
	Columns   []Field
	Predicate null.String
}

type DBInterface interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
}

type TxInterface interface {
	DBInterface
	Commit() error
	Rollback() error
}

type Executable interface {
	Renderable
	Exec(Dialect, DBInterface) (sql.Result, error)
}

type Fetchable interface {
	Renderable
	Fetch(Dialect, DBInterface) (*sqlx.Rows, error)
	FetchRow(Dialect, DBInterface) *sqlx.Row
}
