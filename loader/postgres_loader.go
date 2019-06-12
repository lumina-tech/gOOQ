package loader

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

func NewPostgresLoader() *DatabaseLoader {
	return &DatabaseLoader{
		ConstraintList:          getConstraintList,
		Schema:                  getSchema,
		EnumList:                getEnums,
		EnumValueList:           getEnumValues,
		ReferenceTableValueList: getReferenceTableValues,
		TableList:               getTable,
		ColumnList:              getColumns,
		ParseType:               parseType,
	}
}

func getSchema() (string, error) {
	return "public", nil
}

func getConstraintList(
	db *sqlx.DB, schema, tableName string,
) ([]ConstraintMetaData, error) {
	constraints := []ConstraintMetaData{}
	err := db.Select(&constraints, constraintValuesQuery, schema, tableName)
	if err != nil {
		return nil, err
	}
	return constraints, nil
}

func getEnums(
	db *sqlx.DB, schema string,
) ([]EnumMetadata, error) {
	enums := []EnumMetadata{}
	err := db.Select(&enums, enumsQuery, schema)
	if err != nil {
		return nil, err
	}
	return enums, nil
}

func getEnumValues(
	db *sqlx.DB, schema, enumName string,
) ([]EnumValueMetadata, error) {
	enumValues := []EnumValueMetadata{}
	err := db.Select(&enumValues, enumValuesQuery, schema, enumName)
	if err != nil {
		return nil, err
	}
	return enumValues, nil
}

func getTable(
	db *sqlx.DB, schema string,
) ([]TableMetadata, error) {
	tables := []TableMetadata{}
	err := db.Select(&tables, tablesQuery, schema)
	if err != nil {
		return nil, err
	}
	return tables, nil
}

func getReferenceTableValues(
	db *sqlx.DB, schema, referenceTableName string,
) ([]EnumValueMetadata, error) {
	enumValues := []EnumValueMetadata{}
	query := fmt.Sprintf(referenceTableValuesQuery, schema, referenceTableName)
	err := db.Select(&enumValues, query)
	if err != nil {
		return nil, err
	}
	return enumValues, nil
}

func getColumns(
	db *sqlx.DB, schema, tableName string,
) ([]ColumnMetadata, error) {
	columns := []ColumnMetadata{}
	err := db.Select(&columns, columnsQuery, schema, tableName)
	if err != nil {
		return nil, err
	}
	return columns, nil
}

func parseType(dt string) (string, error) {
	var typ string
	switch strings.ToLower(dt) {
	case "array":
		typ = "StringArray"
	case "boolean":
		typ = "Bool"
	case "character", "character varying", "text", "user-defined":
		typ = "String"
	case "inet":
		typ = "String"
	case "smallint", "integer":
		typ = "Int"
	case "bigint":
		typ = "Int64"
	case "jsonb":
		typ = "Jsonb"
	case "float":
		typ = "Float32"
	case "decimal", "double precision", "numeric":
		typ = "Float64"
	case "date", "timestamp with time zone", "time with time zone", "time without time zone", "timestamp without time zone":
		typ = "Time"
	case "uuid":
		typ = "UUID"
	default:
		return "", fmt.Errorf("Invalid type=%s", dt)
	}
	return typ, nil
}

const tablesQuery = `
select table_name
from information_schema.tables
where table_schema = $1 AND table_name != 'schema_migrations'
order by table_name
`

const columnsQuery = `
SELECT column_name, data_type, is_nullable::boolean, udt_name
FROM information_schema.columns
WHERE table_schema = $1 and table_name = $2
`

const enumsQuery = `
SELECT DISTINCT t.typname as enum_name
FROM pg_type t
JOIN ONLY pg_namespace n ON n.oid = t.typnamespace
JOIN ONLY pg_enum e ON t.oid = e.enumtypid
WHERE n.nspname = $1
`

const enumValuesQuery = `
SELECT e.enumlabel as enum_value, e.enumsortorder as const_value
FROM pg_type t
JOIN ONLY pg_namespace n ON n.oid = t.typnamespace
LEFT JOIN pg_enum e ON t.oid = e.enumtypid
WHERE n.nspname = $1 AND t.typname = $2
`

const referenceTableValuesQuery = `
SELECT value as enum_value from %s.%s order by value
`

const constraintValuesQuery = `
SELECT
	indexes.schemaname AS schema,
	indexes.tablename AS table,
	indexes.indexname AS index_name,
	idx.indisunique AS is_unique,
	idx.indisprimary AS is_primary,
	array_to_json(ARRAY (
		SELECT
			pg_get_indexdef(idx.indexrelid, k + 1, TRUE)
		FROM
			generate_subscripts(idx.indkey, 1) AS k
		ORDER BY
			k)) AS index_keys
	FROM
		pg_indexes AS indexes
		JOIN pg_class AS i ON i.relname = indexes.indexname
		JOIN pg_index AS idx ON idx.indexrelid = i.oid
	WHERE
		schemaname = $1
		AND tablename = $2
`
