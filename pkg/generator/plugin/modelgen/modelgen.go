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
	modelPackage   string
}

func NewGenerator(
	templateString, outputFile, packageName, modelPackage string,
) *ModelGenerator {
	return &ModelGenerator{
		templateString: templateString,
		outputFile:     outputFile,
		packageName:    packageName,
		modelPackage:   modelPackage,
	}
}

func NewModelGenerator(
	outputFile, tablePackage, modelPackage string,
) *ModelGenerator {
	return NewGenerator(modelTemplate, outputFile, modelPackage, modelPackage)
}

func NewTableGenerator(
	outputFile, tablePackage, modelPackage string,
) *ModelGenerator {
	return NewGenerator(tableTemplate, outputFile, tablePackage, modelPackage)
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
		foreignKeyConstraints, err := getForeignKeyConstraintArgs(table)
		if err != nil {
			return err
		}
		modelType := snaker.SnakeToCamelIdentifier(tableName)
		args.Tables = append(args.Tables, TableTemplateArgs{
			TableName:              table.Table.TableName,
			TableType:              snaker.ForceLowerCamelIdentifier(tableName),
			TableSingletonName:     snaker.SnakeToCamelIdentifier(tableName),
			ModelType:              modelType,
			QualifiedModelType:     fmt.Sprintf("%s.%s", gen.modelPackage, modelType),
			IsReferenceTable:       isReferenceTable(tableName),
			ReferenceTableEnumType: getEnumTypeFromReferenceTableName(tableName),
			Fields:                 fields,
			Constraints:            constraints,
			ForeignKeyConstraints:  foreignKeyConstraints,
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
		results = append(results, ConstraintTemplateArgs{
			Name:      constraint.IndexName,
			Columns:   columns,
			Predicate: constraint.IndexPredicate,
		})
	}
	return results, nil
}

func getForeignKeyConstraintArgs(
	table metadata.Table,
) ([]ForeignKeyConstraintTemplateArgs, error) {
	var results []ForeignKeyConstraintTemplateArgs
	for _, constraint := range table.ForeignKeyConstraints {
		results = append(results, ForeignKeyConstraintTemplateArgs{
			Name:              constraint.ConstraintName,
			ColumnName:        constraint.ColumnName,
			ForeignTableName:  constraint.ForeignTableName,
			ForeignColumnName: constraint.ForeignColumnName,
		})
	}
	return results, nil
}

func isReferenceTable(
	tableName string,
) bool {
	return strings.HasSuffix(tableName, metadata.ReferenceTableSuffix)
}

func getEnumTypeFromReferenceTableName(
	tableName string,
) string {
	enumNameSnakeCase := strings.ReplaceAll(tableName, metadata.ReferenceTableSuffix, "")
	return snaker.SnakeToCamelIdentifier(enumNameSnakeCase)
}
