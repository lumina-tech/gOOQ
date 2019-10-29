package generator

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/knq/snaker"
	"github.com/lumina-tech/gooq/pkg/generator/metadata"
	"github.com/lumina-tech/gooq/pkg/generator/postgres"
)

const (
	ModelTemplateFilename = "model.go.tmpl"
)

func GenerateModel(
	db *sqlx.DB, templateString, outputPath string, dbName, packageName string,
) {
	dbLoader := postgres.NewPostgresLoader()
	schema, err := dbLoader.Schema()
	check(err)
	tables, err := dbLoader.TableList(db, schema)
	check(err)
	args := TablesTemplateArgs{
		Timestamp: time.Now().Format(time.RFC3339),
		Package:   packageName,
		Schema:    schema,
		Tables:    make([]TableTemplateArgs, 0),
	}
	for _, table := range tables {
		columns, err := dbLoader.ColumnList(db, schema, table.TableName)
		check(err)
		constraints, err := dbLoader.ConstraintList(db, schema, table.TableName)
		check(err)
		tableTemplateArgs := TableTemplateArgs{
			Name:                   strings.ToLower(table.TableName),
			DatabaseName:           fmt.Sprintf("%sDatabase", snaker.SnakeToCamel(dbName)),
			ModelType:              snaker.SnakeToCamel(table.TableName),
			ModelTypeWithPackage:   fmt.Sprintf("%s.%s", "model", snaker.SnakeToCamel(table.TableName)),
			ModelTableName:         snaker.SnakeToCamel(table.TableName),
			RepositoryName:         fmt.Sprintf("%sRepository", snaker.SnakeToCamel(table.TableName)),
			IsReferenceTable:       isReferenceTable(table.TableName),
			ReferenceTableEnumType: getReferenceTableEnumType(table.TableName),
			Fields:                 make([]FieldTemplateArgs, len(columns)),
			Constraints:            make([]ConstraintTemplateArgs, len(constraints)),
		}
		columnToRefTableMapping := getColumnToReferenceTableMapping(db, dbLoader, table.TableName)
		for colIndex, column := range columns {
			fieldTemplateArg := getFieldTemplate(dbLoader, column, columnToRefTableMapping)
			tableTemplateArgs.Fields[colIndex] = fieldTemplateArg
		}
		for conIndex, constraint := range constraints {
			tableTemplateArgs.Constraints[conIndex] = getConstraintTemplate(constraint)
		}
		args.Tables = append(args.Tables, tableTemplateArgs)
	}
	schemaTemplate := getTemplate(templateString)
	err = RenderToFile(schemaTemplate, outputPath, args)
	check(err)
}

func getColumnToReferenceTableMapping(
	db *sqlx.DB, dbLoader *metadata.DatabaseMetadataLoader, tableName string,
) map[string]string {
	foreignKeyConstraints, err := dbLoader.ForeignKeyConstraintList(db, tableName)
	check(err)
	result := make(map[string]string)
	for _, fk := range foreignKeyConstraints {
		if isReferenceTable(fk.ForeignTableName) {
			result[fk.ColumnName] = getReferenceTableEnumType(fk.ForeignTableName)
		}
	}
	return result
}

func getFieldTemplate(
	dbLoader *metadata.DatabaseMetadataLoader,
	column metadata.ColumnMetadata,
	columnToRefTableEnum map[string]string,
) FieldTemplateArgs {
	dataType, err := dbLoader.ParseType(column.DataType)
	check(err)
	gooqType := dataType
	literal := typesMap[gooqType].Literal
	if column.IsNullable {
		literal = typesMap[gooqType].NullableLiteral
	}
	if enumName, ok := columnToRefTableEnum[column.ColumnName]; ok {
		literal = enumName
	} else if column.DataType == "USER-DEFINED" && column.UserDefinedTypeName != "citext" {
		// citext is the only user-defined type that is not an enum
		literal = snaker.SnakeToCamelIdentifier(column.UserDefinedTypeName)
	}
	return FieldTemplateArgs{
		Name:     column.ColumnName,
		GooqType: gooqType,
		Type:     literal,
	}
}

func getConstraintTemplate(
	constraint metadata.ConstraintMetadata,
) ConstraintTemplateArgs {
	columns := []string{}
	err := json.Unmarshal([]byte(constraint.IndexKeys), &columns)
	check(err)
	for index := range columns {
		column := columns[index]
		columns[index] = strings.Replace(column, "\"", "\\\"", -1)
	}
	columnsString := fmt.Sprintf("{\"%s\"}", strings.Join(columns, "\",\""))
	return ConstraintTemplateArgs{
		Name:      constraint.IndexName,
		Columns:   columnsString,
		Predicate: constraint.IndexPredicate,
	}
}

func isReferenceTable(
	tableName string,
) bool {
	return strings.HasSuffix(tableName, "_reference_table")
}

func getReferenceTableEnumType(
	tableName string,
) string {
	enumNameSnakeCase := strings.ReplaceAll(tableName, "_reference_table", "")
	return snaker.SnakeToCamelIdentifier(enumNameSnakeCase)
}
