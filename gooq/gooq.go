package gooq

import (
	"bytes"
	"database/sql"
	"io"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Base on https://www.jooq.org/javadoc/latest/

type PredicateType int
type JoinType int
type Dialect int

const (
	EqPredicate PredicateType = iota
	GtPredicate
	GtePredicate
	ILikePredicate
	InPredicate
	IsNotNullPredicate
	IsNullPredicate
	LikePredicate
	LtPredicate
	LtePredicate
	NotEqPredicate
	NotInPredicate
	NotILikePredicate
	NotLikePredicate
	NotSimilarPredicate
	SimilarPredicate
)

const (
	Join JoinType = iota
	LeftOuterJoin
	NotJoined
)

const (
	Sqlite Dialect = iota
	MySQL
	Postgres
)

type DatabaseConstraint struct {
	Name    string
	Columns []string
}

type Aliasable interface {
	Alias() string
	MaybeAlias() string
}

type TableLike interface {
	Selectable
	Name() string
}

type FieldFunction struct {
	Child *FieldFunction
	Name  string
	Expr  string
	Args  []interface{}
}

type Field interface {
	Aliasable
	Functional
	Name() string
	String() string
	FilterWhere(...Condition) Field
	As(string) Field
	Function() FieldFunction
	Filters() []Condition
}

type TableField interface {
	Field
	Parent() Selectable
}

type FieldBinding struct {
	Field Field
	Value interface{}
}

type TableFieldBinding struct {
	Field TableField
	Value interface{}
}

type Condition struct {
	Binding   FieldBinding
	Predicate PredicateType
}

///////////////////////////////////////////////////////////////////////////////
// Select
///////////////////////////////////////////////////////////////////////////////

type SelectFromStep interface {
	SelectWhereStep
	From(Selectable) SelectJoinStep
}

type SelectJoinStep interface {
	SelectWhereStep
	Join(Selectable) SelectOnStep
	LeftOuterJoin(Selectable) SelectOnStep
}

type SelectOnStep interface {
	SelectWhereStep
	On(...JoinCondition) SelectJoinStep
}

type SelectWhereStep interface {
	SelectGroupByStep
	Where(conditions ...Condition) SelectGroupByStep
}

type SelectGroupByStep interface {
	SelectHavingStep
	GroupBy(...Field) SelectHavingStep
}

type SelectHavingStep interface {
	SelectOrderByStep
	Having(conditions ...Condition) SelectOrderByStep
}

type SelectOrderByStep interface {
	SelectLimitStep
	OrderBy(...Field) SelectLimitStep
}

type SelectLimitStep interface {
	SelectOffsetStep
	Limit(limit int) SelectOffsetStep
}

type SelectOffsetStep interface {
	SelectFinalStep
	Offset(offset int) SelectFinalStep
}

type SelectFinalStep interface {
	Selectable
	Fetchable
	As(string) SelectFinalStep
	Union(SelectFinalStep) SelectOrderByStep
}

///////////////////////////////////////////////////////////////////////////////
// Create Table
///////////////////////////////////////////////////////////////////////////////

type CreateTableAsStep interface {
	CreateTableColumnStep
}

type CreateTableColumnStep interface {
	CreateTableOnCommitStep
	Column(f TableField, v interface{}) CreateTableColumnStep
}

type CreateTableOnCommitStep interface {
	CreateTableFinalStep
	OnCommit(CommitAction) CreateTableFinalStep
}

type CreateTableFinalStep interface {
	Executable
}

///////////////////////////////////////////////////////////////////////////////
// Delete
///////////////////////////////////////////////////////////////////////////////

type DeleteUsingStep interface {
	DeleteWhereStep
	Using(Selectable) DeleteOnStep
}

type DeleteWhereStep interface {
	DeleteResultStep
	Where(...Condition) DeleteReturningStep
}

type DeleteOnStep interface {
	DeleteResultStep
	On(...JoinCondition) DeleteReturningStep
}

type DeleteReturningStep interface {
	DeleteFinalStep
	Returning(TableField) DeleteResultStep
}

type DeleteResultStep interface {
	Fetchable
}

type DeleteFinalStep interface {
	Executable
}

///////////////////////////////////////////////////////////////////////////////
// Insert
///////////////////////////////////////////////////////////////////////////////

type InsertSetStep interface {
	InsertSetMoreStep
	Select(s Selectable) InsertOnConflictStep
}

type InsertSetMoreStep interface {
	InsertOnConflictStep
	Set(f TableField, v interface{}) InsertSetMoreStep
}

type InsertOnConflictStep interface {
	InsertReturningStep
	OnConflictDoNothing() InsertReturningStep
	OnConflictDoUpdate(*DatabaseConstraint, ...TableFieldBinding) InsertReturningStep
}

type InsertReturningStep interface {
	InsertFinalStep
	Returning(TableField) InsertResultStep
}

type InsertResultStep interface {
	Fetchable
}

type InsertFinalStep interface {
	Executable
}

///////////////////////////////////////////////////////////////////////////////
// Update
///////////////////////////////////////////////////////////////////////////////

// Update with JOIN
// https://www.postgresql.org/docs/11/sql-update.html
//
// In MySQL
// UPDATE employees
// LEFT JOIN departments ON employees.department_id = departments.id
// SET department_name = departments.name
//
// In Postgres
// UPDATE employees
// SET department_name = departments.name
// FROM departments
// WHERE employees.department_id = departments.id
//
// Update with JOIN in SQL is more explicit so we follow their syntax here
// in lieu of making things up ourselves
type UpdateSetStep interface {
	UpdateJoinStep
	Set(f TableField, v interface{}) UpdateSetStep
}

type UpdateJoinStep interface {
	UpdateWhereStep
	Join(Selectable) UpdateOnStep
}

type UpdateWhereStep interface {
	UpdateOnConflictStep
	Where(conditions ...Condition) UpdateOnConflictStep
}

type UpdateOnStep interface {
	UpdateOnConflictStep
	On(...JoinCondition) UpdateOnConflictStep
}

type UpdateOnConflictStep interface {
	UpdateReturningStep
	OnConflict(ConflictAction) UpdateReturningStep
}

type UpdateReturningStep interface {
	UpdateFinalStep
	Returning(TableField) UpdateResultStep
}

type UpdateResultStep interface {
	Fetchable
}

type UpdateFinalStep interface {
	Executable
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

type TxDBInterface interface {
	DBInterface
	Commit() error
	Rollback() error
}

type Renderable interface {
	Render(Dialect, io.Writer) []interface{}
	RenderWithPlaceholders(Dialect, io.Writer, []interface{}) []interface{}
	String(Dialect) string
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

type Selectable interface {
	Aliasable
	IsSelectable()
}

type JoinCondition struct {
	Lhs, Rhs  TableField
	Predicate PredicateType
}

type join struct {
	target   Selectable
	joinType JoinType
	conds    []JoinCondition
}

func exec(d Dialect, r Renderable, db DBInterface) (sql.Result, error) {
	var buf bytes.Buffer
	args := r.Render(d, &buf)
	return db.Exec(buf.String(), args...)
}

func Qualified(parts ...string) string {
	tmp := []string{}
	for _, part := range parts {
		if part != "" {
			tmp = append(tmp, part)
		}
	}
	return strings.Join(tmp, ".")
}
