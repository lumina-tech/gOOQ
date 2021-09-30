package metadata

import (
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v3"
)

type DataType struct {
	Name            string
	Literal         string
	NullableLiteral string
}

var (
	DataTypeBool        = DataType{Name: "Bool", Literal: "bool", NullableLiteral: "null.Bool"}
	DataTypeFloat32     = DataType{Name: "Decimal", Literal: "float32", NullableLiteral: "null.Float"}
	DataTypeFloat64     = DataType{Name: "Decimal", Literal: "float64", NullableLiteral: "null.Float"}
	DataTypeInt         = DataType{Name: "Int", Literal: "int", NullableLiteral: "null.Int"}
	DataTypeInt64       = DataType{Name: "Int64", Literal: "int64", NullableLiteral: "null.Int"}
	DataTypeBigInt      = DataType{Name: "BigInt", Literal: "big.Int", NullableLiteral: "big.Int"}
	DataTypeBigFloat    = DataType{Name: "BigFloat", Literal: "big.Float", NullableLiteral: "big.Float"}
	DataTypeJSONB       = DataType{Name: "Jsonb", Literal: "[]byte", NullableLiteral: "nullable.Jsonb"}
	DataTypeString      = DataType{Name: "String", Literal: "string", NullableLiteral: "null.String"}
	DataTypeStringArray = DataType{Name: "StringArray", Literal: "pq.StringArray", NullableLiteral: "pq.StringArray"}
	DataTypeTime        = DataType{Name: "Time", Literal: "time.Time", NullableLiteral: "null.Time"}
	DataTypeUUID        = DataType{Name: "UUID", Literal: "uuid.UUID", NullableLiteral: "nullable.UUID"}
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

type ConstraintMetadata struct {
	Schema         string      `db:"schema"`
	Table          string      `db:"table"`
	IndexName      string      `db:"index_name"`
	IndexPredicate null.String `db:"index_predicate"`
	IsUnique       bool        `db:"is_unique"`
	IsPrimary      bool        `db:"is_primary"`
	IndexKeys      string      `db:"index_keys"`
}

type ForeignKeyConstraintMetadata struct {
	TableSchema        string `db:"table_schema"`
	ConstraintName     string `db:"constraint_name"`
	TableName          string `db:"table_name"`
	ColumnName         string `db:"column_name"`
	ForeignTableSchema string `db:"foreign_table_schema"`
	ForeignTableName   string `db:"foreign_table_name"`
	ForeignColumnName  string `db:"foreign_column_name"`
}

type Loader struct {
	Schema                   func() (string, error)
	TableList                func(*sqlx.DB, string) ([]TableMetadata, error)
	ColumnList               func(*sqlx.DB, string, string) ([]ColumnMetadata, error)
	ConstraintList           func(*sqlx.DB, string, string) ([]ConstraintMetadata, error)
	ForeignKeyConstraintList func(*sqlx.DB, string) ([]ForeignKeyConstraintMetadata, error)
	EnumList                 func(*sqlx.DB, string) ([]EnumMetadata, error)
	EnumValueList            func(*sqlx.DB, string, string) ([]EnumValueMetadata, error)
	ReferenceTableValueList  func(*sqlx.DB, string, string) ([]EnumValueMetadata, error)
	GetDataType              func(string) (DataType, error)
}
