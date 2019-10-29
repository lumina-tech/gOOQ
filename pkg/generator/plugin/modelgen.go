package plugin

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/knq/snaker"
	"github.com/lumina-tech/gooq/pkg/generator/data"
)

const (
	ModelTemplateFilename = "model.go.tmpl"
)

type ModelGenerator struct {
	templateString string
	outputFile     string
	packageName    string
	dbName         string
}

func NewModelGenerator(
	templateString, modelOutputPath, packageName, dbname string,
) *ModelGenerator {
	return &ModelGenerator{
		templateString: templateString,
		outputFile:     modelOutputPath,
		packageName:    packageName,
		dbName:         dbname,
	}
}

func (gen *ModelGenerator) GenerateCode(
	data *data.Data,
) error {
	args := TablesTemplateArgs{
		Timestamp: time.Now().Format(time.RFC3339),
		Package:   gen.packageName,
		Schema:    data.Schema,
		Tables:    make([]TableTemplateArgs, 0),
	}
	for _, table := range data.Tables {
		tableName := table.Table.TableName
		tableTemplateArgs := TableTemplateArgs{
			Name:                   strings.ToLower(tableName),
			DatabaseName:           fmt.Sprintf("%sDatabase", snaker.SnakeToCamel(gen.dbName)),
			ModelType:              snaker.SnakeToCamel(tableName),
			ModelTypeWithPackage:   fmt.Sprintf("%s.%s", "model", snaker.SnakeToCamel(tableName)),
			ModelTableName:         snaker.SnakeToCamel(tableName),
			RepositoryName:         fmt.Sprintf("%sRepository", snaker.SnakeToCamel(tableName)),
			IsReferenceTable:       isReferenceTable(tableName),
			ReferenceTableEnumType: getReferenceTableEnumType(tableName),
			//Fields:                 make([]FieldTemplateArgs, len(columns)),
			//Constraints:            make([]ConstraintTemplateArgs, len(constraints)),
		}
		var err error
		tableTemplateArgs.Fields, err = getFieldArgs(data, table)
		if err != nil {
			return err
		}
		tableTemplateArgs.Constraints = getConstraintArgs(data, table)
	}
	enumTemplate := getTemplate(gen.modelTemplateString)
	return RenderToFile(enumTemplate, gen.modelOutputFile, args)
}

func getColumnToTypeMapping(
	table data.Table,
) map[string]string {
	result := make(map[string]string)
	for _, fk := range table.ForeignKeyConstraints {
		if isReferenceTable(fk.ForeignTableName) {
			result[fk.ColumnName] = getReferenceTableEnumType(fk.ForeignTableName)
		}
	}
	return result
}

func getFieldArgs(
	data *data.Data, table data.Table,
) ([]FieldTemplateArgs, error) {
	columnToRefTableMapping := getColumnToTypeMapping(table)
	var results []FieldTemplateArgs
	for _, column := range table.Columns {
		dataType, err := data.Loader.GetDataType(column.DataType)
		if err != nil {
			return nil, err
		}
		literal := dataType.Literal
		if column.IsNullable {
			literal = dataType.NullableLiteral
		}
		if enumName, ok := columnToRefTableMapping[column.ColumnName]; ok {
			literal = enumName
		} else if column.DataType == "USER-DEFINED" && column.UserDefinedTypeName != "citext" {
			// citext is the only user-defined type that is not an enum
			literal = snaker.SnakeToCamelIdentifier(column.UserDefinedTypeName)
		}
		results = append(results, FieldTemplateArgs{
			Name:     column.ColumnName,
			GooqType: dataType.Name,
			Type:     literal,
		})
	}
	return results, nil
}

func getConstraintArgs(
	data *data.Data, table data.Table,
) []ConstraintTemplateArgs {
	var results []ConstraintTemplateArgs
	for _, constraint := range table.Constraints {
		columns := []string{}
		err := json.Unmarshal([]byte(constraint.IndexKeys), &columns)
		check(err)
		for index := range columns {
			column := columns[index]
			columns[index] = strings.Replace(column, "\"", "\\\"", -1)
		}
		columnsString := fmt.Sprintf("{\"%s\"}", strings.Join(columns, "\",\""))
		results = append(results, ConstraintTemplateArgs{
			Name:      constraint.IndexName,
			Columns:   columnsString,
			Predicate: constraint.IndexPredicate,
		})
	}
	return results
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
