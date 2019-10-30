package modelgen

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/knq/snaker"
	"github.com/lumina-tech/gooq/pkg/generator/metadata"
	"github.com/lumina-tech/gooq/pkg/generator/utils"
)

type ModelGenerator struct {
	templateString string
	outputFile     string
	packageName    string
}

func NewModelGenerator(
	outputFile, packageName string,
) *ModelGenerator {
	return &ModelGenerator{
		templateString: modelTemplate,
		outputFile:     outputFile,
		packageName:    packageName,
	}
}

func NewTableGenerator(
	outputFile, packageName string,
) *ModelGenerator {
	return &ModelGenerator{
		templateString: tableTemplate,
		outputFile:     outputFile,
		packageName:    packageName,
	}
}

func (gen *ModelGenerator) GenerateCode(
	data *metadata.Data,
) error {
	args := TemplateArgs{
		Timestamp: time.Now().Format(time.RFC3339),
		Package:   gen.packageName,
		Schema:    data.Schema,
		Tables:    make([]TableTemplateArgs, 0),
	}
	for _, table := range data.Tables {
		tableName := table.Table.TableName
		fields, err := getFieldArgs(data, table)
		if err != nil {
			return err
		}
		constraints, err := getConstraintArgs(table)
		if err != nil {
			return err
		}
		args.Tables = append(args.Tables, TableTemplateArgs{
			Name:                   strings.ToLower(tableName),
			ModelType:              snaker.SnakeToCamel(tableName),
			ModelTypeWithPackage:   fmt.Sprintf("%s.%s", gen.packageName, snaker.SnakeToCamel(tableName)),
			ModelTableName:         snaker.SnakeToCamel(tableName),
			IsReferenceTable:       isReferenceTable(tableName),
			ReferenceTableEnumType: getEnumTypeFromReferenceTableName(tableName),
			Fields:                 fields,
			Constraints:            constraints,
		})
	}
	enumTemplate := utils.GetTemplate(gen.templateString)
	return utils.RenderToFile(enumTemplate, gen.outputFile, args)
}

func getColumnToTypeMapping(
	table metadata.Table,
) map[string]string {
	result := make(map[string]string)
	for _, fk := range table.ForeignKeyConstraints {
		if isReferenceTable(fk.ForeignTableName) {
			result[fk.ColumnName] = getEnumTypeFromReferenceTableName(fk.ForeignTableName)
		}
	}
	return result
}

func getFieldArgs(
	data *metadata.Data, table metadata.Table,
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
	table metadata.Table,
) ([]ConstraintTemplateArgs, error) {
	var results []ConstraintTemplateArgs
	for _, constraint := range table.Constraints {
		var columns []string
		err := json.Unmarshal([]byte(constraint.IndexKeys), &columns)
		if err != nil {
			return nil, err
		}
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
	return results, nil
}

func isReferenceTable(
	tableName string,
) bool {
	return strings.HasSuffix(tableName, "_reference_table")
}

func getEnumTypeFromReferenceTableName(
	tableName string,
) string {
	enumNameSnakeCase := strings.ReplaceAll(tableName, "_reference_table", "")
	return snaker.SnakeToCamelIdentifier(enumNameSnakeCase)
}
