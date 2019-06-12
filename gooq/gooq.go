package gooq

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type DatabaseConstraint struct {
	Name    string
	Columns []Field
}

///////////////////////////////////////////////////////////////////////////////
// Query
///////////////////////////////////////////////////////////////////////////////

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
	FetchRow(Dialect, DBInterface) (*sqlx.Row, error)
}
