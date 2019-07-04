package generator

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/knq/snaker"
	"github.com/lumina-tech/gooq/common"
	"github.com/lumina-tech/gooq/loader"
)

const (
	ModelTemplateFilename = "model.go.tmpl"
)

func GenerateModel(
	db *sqlx.DB, templatePath, outputPath, dbname string,
) {
	dbloader := loader.NewPostgresLoader()
	schema, err := dbloader.Schema()
	check(err)
	tables, err := dbloader.TableList(db, schema)
	check(err)
	args := TablesTemplateArgs{
		Timestamp: time.Now().Format(time.RFC3339),
		Package:   "model",
		Schema:    schema,
		Tables:    make([]TableTemplateArgs, 0),
	}
	for tableIndex, table := range tables {
		columns, err := dbloader.ColumnList(
			db, schema, table.TableName)
		check(err)
		constraints, err := dbloader.ConstraintList(
			db, schema, table.TableName)
		check(err)
		args.Tables = append(args.Tables, TableTemplateArgs{
			Name:        table.TableName,
			Fields:      make([]FieldTemplateArgs, len(columns)),
			Constraints: make([]ConstraintTemplateArgs, len(constraints)),
		})
		for colIndex, column := range columns {
			args.Tables[tableIndex].Fields[colIndex] = getFieldTemplate(dbloader, column)
		}
		for conIndex, constraint := range constraints {
			args.Tables[tableIndex].Constraints[conIndex] = getConstraintTemplate(constraint)
		}
	}

	filename := fmt.Sprintf("%s/%s_model.generated.go", outputPath, dbname)
	schemaTemplate := getTemplate(templatePath)
	err = common.RenderToFile(schemaTemplate, filename, args)
	check(err)
}

func getFieldTemplate(
	dbloader *loader.DatabaseLoader,
	column loader.ColumnMetadata,
) FieldTemplateArgs {
	dataType, err := dbloader.ParseType(column.DataType)
	check(err)
	gooqType := dataType
	if column.IsNullable {
		gooqType = "Null" + dataType
	}
	literal := typesMap[gooqType].Literal
	// citext is the only user-defined type that is not an enum
	if column.DataType == "USER-DEFINED" &&
		column.UserDefinedTypeName != "citext" {
		literal = snaker.SnakeToCamelIdentifier(column.UserDefinedTypeName)
	}
	return FieldTemplateArgs{
		Name:     column.ColumnName,
		GooqType: gooqType,
		Type:     literal,
	}

}

func getConstraintTemplate(
	constraint loader.ConstraintMetaData,
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
		Name:    constraint.IndexName,
		Columns: columnsString,
	}
}
