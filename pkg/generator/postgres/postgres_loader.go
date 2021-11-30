package postgres

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lumina-tech/gooq/pkg/generator/metadata"
)

func NewPostgresLoader() *metadata.Loader {
	return &metadata.Loader{
		ConstraintList:           getConstraintList,
		ForeignKeyConstraintList: getForeignKeyConstraintList,
		Schema:                   getSchema,
		EnumList:                 getEnums,
		EnumValueList:            getEnumValues,
		ReferenceTableValueList:  getReferenceTableValues,
		TableList:                getTable,
		ColumnList:               getColumns,
		GetDataType:              parseType,
		GetTypeByName:            getTypeByName,
	}
}

func getSchema() (string, error) {
	return "public", nil
}

func getConstraintList(
	db *sqlx.DB, schema, tableName string,
) ([]metadata.ConstraintMetadata, error) {
	constraints := []metadata.ConstraintMetadata{}
	err := db.Select(&constraints, constraintValuesQuery, schema, tableName)
	if err != nil {
		return nil, err
	}
	return constraints, nil
}

func getForeignKeyConstraintList(
	db *sqlx.DB, tableName string,
) ([]metadata.ForeignKeyConstraintMetadata, error) {
	constraints := []metadata.ForeignKeyConstraintMetadata{}
	err := db.Select(&constraints, foreignKeyConstraintValuesQuery, tableName)
	if err != nil {
		return nil, err
	}
	return constraints, nil
}

func getEnums(
	db *sqlx.DB, schema string,
) ([]metadata.EnumMetadata, error) {
	enums := []metadata.EnumMetadata{}
	err := db.Select(&enums, enumsQuery, schema)
	if err != nil {
		return nil, err
	}
	return enums, nil
}

func getEnumValues(
	db *sqlx.DB, schema, enumName string,
) ([]metadata.EnumValueMetadata, error) {
	enumValues := []metadata.EnumValueMetadata{}
	err := db.Select(&enumValues, enumValuesQuery, schema, enumName)
	if err != nil {
		return nil, err
	}
	return enumValues, nil
}

func getTable(
	db *sqlx.DB, schema string,
) ([]metadata.TableMetadata, error) {
	tables := []metadata.TableMetadata{}
	err := db.Select(&tables, tablesQuery, schema)
	if err != nil {
		return nil, err
	}
	return tables, nil
}

func getReferenceTableValues(
	db *sqlx.DB, schema, referenceTableName string,
) ([]metadata.EnumValueMetadata, error) {
	enumValues := []metadata.EnumValueMetadata{}
	query := fmt.Sprintf(referenceTableValuesQuery, schema, referenceTableName)
	err := db.Select(&enumValues, query)
	if err != nil {
		return nil, err
	}
	return enumValues, nil
}

func getColumns(
	db *sqlx.DB, schema, tableName string,
) ([]metadata.ColumnMetadata, error) {
	columns := []metadata.ColumnMetadata{}
	err := db.Select(&columns, columnsQuery, schema, tableName)
	if err != nil {
		return nil, err
	}
	return columns, nil
}

func getTypeByName(
	typeName string,
) (metadata.DataType, error) {
	dataType, ok := metadata.NameToType[typeName]
	if !ok {
		return metadata.DataType{}, fmt.Errorf("type with name %s does not exist", typeName)
	}
	return dataType, nil
}

func parseType(
	dataType string,
) (metadata.DataType, error) {
	var typ metadata.DataType
	switch strings.ToLower(dataType) {
	case "array":
		typ = metadata.DataTypeStringArray
	case "boolean":
		typ = metadata.DataTypeBool
	case "character", "character varying", "text", "user-defined":
		typ = metadata.DataTypeString
	case "inet":
		typ = metadata.DataTypeString
	case "smallint", "integer":
		typ = metadata.DataTypeInt
	case "bigint":
		typ = metadata.DataTypeInt64
	case "json":
		typ = metadata.DataTypeString
	case "jsonb":
		typ = metadata.DataTypeJSONB
	case "float":
		typ = metadata.DataTypeFloat32
	case "decimal", "double precision", "numeric":
		typ = metadata.DataTypeFloat64
	case "date", "timestamp with time zone", "time with time zone", "time without time zone", "timestamp without time zone":
		typ = metadata.DataTypeTime
	case "uuid":
		typ = metadata.DataTypeUUID
	default:
		return metadata.DataType{}, fmt.Errorf("invalid type=%s", dataType)
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
	pg_get_expr(idx.indpred, idx.indrelid) AS index_predicate,
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
  ORDER BY
  	indexes.indexname
`

const foreignKeyConstraintValuesQuery = `
SELECT
	tc.table_schema,
	tc.constraint_name,
	tc.table_name,
	kcu.column_name,
	ccu.table_schema AS foreign_table_schema,
	ccu.table_name AS foreign_table_name,
	ccu.column_name AS foreign_column_name
	FROM
		information_schema.table_constraints AS tc
	JOIN information_schema.key_column_usage AS kcu
	ON tc.constraint_name = kcu.constraint_name
	AND tc.table_schema = kcu.table_schema
	JOIN information_schema.constraint_column_usage AS ccu
	ON ccu.constraint_name = tc.constraint_name
	AND ccu.table_schema = tc.table_schema
	WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_name=$1
`
