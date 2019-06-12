package loader

import (
	"github.com/jmoiron/sqlx"
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
	Schema    string `db:"schema"`
	Table     string `db:"table"`
	IndexName string `db:"index_name"`
	IsUnique  bool   `db:"is_unique"`
	IsPrimary bool   `db:"is_primary"`
	IndexKeys string `db:"index_keys"`
}

type DatabaseLoader struct {
	Schema                  func() (string, error)
	ConstraintList          func(*sqlx.DB, string, string) ([]ConstraintMetaData, error)
	EnumList                func(*sqlx.DB, string) ([]EnumMetadata, error)
	EnumValueList           func(*sqlx.DB, string, string) ([]EnumValueMetadata, error)
	ReferenceTableValueList func(*sqlx.DB, string, string) ([]EnumValueMetadata, error)
	TableList               func(*sqlx.DB, string) ([]TableMetadata, error)
	ColumnList              func(*sqlx.DB, string, string) ([]ColumnMetadata, error)
	ParseType               func(string) (string, error)
}
