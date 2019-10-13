package metadata

import (
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v3"
)

type EnumMetadata struct {
	EnumName string `db:"enum_name"`
}

type EnumValueMetadata struct {
	Description string
	EnumValue   string `db:"enum_value"`
	ConstValue  int    `db:"const_value"`
}

type TableMetadata struct {
	Type      string `db:"type"`
	TableName string `db:"table_name"`
	ManualPk  bool   `db:"manual_pk"`
}

type ColumnMetadata struct {
	ColumnName          string `db:"column_name"`
	DataType            string `db:"data_type"`
	IsNullable          bool   `db:"is_nullable"`
	UserDefinedTypeName string `db:"udt_name"`
}

type ConstraintMetaData struct {
	Schema         string      `db:"schema"`
	Table          string      `db:"table"`
	IndexName      string      `db:"index_name"`
	IndexPredicate null.String `db:"index_predicate"`
	IsUnique       bool        `db:"is_unique"`
	IsPrimary      bool        `db:"is_primary"`
	IndexKeys      string      `db:"index_keys"`
}

type ForeignKeyConstraintMetaData struct {
	TableSchema        string `db:"table_schema"`
	ConstraintName     string `db:"constraint_name"`
	TableName          string `db:"table_name"`
	ColumnName         string `db:"column_name"`
	ForeignTableSchema string `db:"foreign_table_schema"`
	ForeignTableName   string `db:"foreign_table_name"`
	ForeignColumnName  string `db:"foreign_column_name"`
}

type DatabaseMetadataLoader struct {
	Schema                   func() (string, error)
	TableList                func(*sqlx.DB, string) ([]TableMetadata, error)
	ColumnList               func(*sqlx.DB, string, string) ([]ColumnMetadata, error)
	ConstraintList           func(*sqlx.DB, string, string) ([]ConstraintMetaData, error)
	ForeignKeyConstraintList func(*sqlx.DB, string) ([]ForeignKeyConstraintMetaData, error)
	EnumList                 func(*sqlx.DB, string) ([]EnumMetadata, error)
	EnumValueList            func(*sqlx.DB, string, string) ([]EnumValueMetadata, error)
	ReferenceTableValueList  func(*sqlx.DB, string, string) ([]EnumValueMetadata, error)
	ParseType                func(string) (string, error)
}
