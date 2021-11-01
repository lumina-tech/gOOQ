package modelgen

import (
	"encoding/json"
	"fmt"
	"strconv"
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
	overrides      map[string]interface{}
}

const (
	OverrideModelsConfig = "models"
	OverrideFieldsConfig = "fields"
	OverrideTypeConfig   = "overridetype"
)

func NewGenerator(
	templateString, outputFile, packageName, modelPackage string, overrides map[string]interface{},
) *ModelGenerator {
	return &ModelGenerator{
		templateString: templateString,
		outputFile:     outputFile,
		packageName:    packageName,
		modelPackage:   modelPackage,
		overrides:      overrides,
	}
}

func NewModelGenerator(
	outputFile, tablePackage, modelPackage string, overrides map[string]interface{},
) *ModelGenerator {
	return NewGenerator(modelTemplate, outputFile, modelPackage, modelPackage, overrides)
}

func NewTableGenerator(
	outputFile, tablePackage, modelPackage string, overrides map[string]interface{},
) *ModelGenerator {
	return NewGenerator(tableTemplate, outputFile, tablePackage, modelPackage, overrides)
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
		fields, err := getFieldArgs(data, table, gen.overrides)
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
	data *metadata.Data, table metadata.Table, overrides map[string]interface{},
) ([]FieldTemplateArgs, error) {
	columnToRefTableMapping := getColumnToTypeMapping(table)
	var results []FieldTemplateArgs
	for _, column := range table.Columns {
		var dataType metadata.DataType
		var err error

		if dataTypeKey, ok := getOverrideDataType(table.Table.TableName, column.ColumnName, overrides); ok {
			dataType, err = data.Loader.GetTypeByName(dataTypeKey)
		} else {
			dataType, err = data.Loader.GetDataType(column.DataType)
		}

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
			// some of the constraint columns have quotes we should unquote them
			if strings.ContainsAny(column, "\"") {
				unquotedColumn, _ := strconv.Unquote(column)
				columns[index] = unquotedColumn
			}
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

func getOverrideDataType(
	tableName string,
	columnName string,
	overrides map[string]interface{},
) (string, bool) {

	models := overrides[OverrideModelsConfig]
	if models != nil {
		if fieldsMap, ok := models.(map[string]interface{})[tableName]; ok {
			fields := fieldsMap.(map[string]interface{})[OverrideFieldsConfig]
			if _, ok := fields.(map[string]interface{})[columnName]; ok {
				overrideTypeMap := fields.(map[string]interface{})[columnName]
				overrideType := overrideTypeMap.(map[string]interface{})[OverrideTypeConfig]
				return fmt.Sprintf("%v", overrideType), true
			}
		}
	}
	return "", false
}
