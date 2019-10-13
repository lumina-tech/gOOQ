package generator

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/knq/snaker"
	"github.com/lumina-tech/gooq/pkg/generator/postgres"
)

const (
	referenceTableSuffix = "_reference_table"
)

type EnumGenerator struct {
	db                  *sqlx.DB
	modelTemplateString string
	modelOutputFile     string
	dbname              string
}

func NewEnumGenerator(
	db *sqlx.DB,
	modelTemplatePath, modelOutputPath,
	dbname string,
) *EnumGenerator {
	modelOutputFile := fmt.Sprintf("%s/%s_enum.generated.go",
		modelOutputPath, dbname)
	return &EnumGenerator{
		db,
		modelTemplatePath,
		modelOutputFile,
		dbname,
	}
}

func (generator *EnumGenerator) Run() {
	args := generator.getTemplateArguments()
	enumTemplate := getTemplate(generator.modelTemplateString)
	err := RenderToFile(enumTemplate, generator.modelOutputFile, args)
	check(err)
}

func (generator *EnumGenerator) getTemplateArguments() EnumsTemplateArgs {
	results := generator.getEnumTypesFromDatabaseTypes()
	results = append(results, generator.getEnumTypesFromReferenceTables()...)
	sort.SliceStable(results, func(i, j int) bool {
		return strings.Compare(results[i].Name, results[j].Name) < 0
	})
	return EnumsTemplateArgs{
		Package:   "model",
		Timestamp: time.Now().Format(time.RFC3339),
		Enums:     results,
	}
}

func (generator *EnumGenerator) getEnumTypesFromDatabaseTypes() []EnumType {
	loader := postgres.NewPostgresLoader()
	schema, err := loader.Schema()
	check(err)
	dbEnums, err := loader.EnumList(generator.db, schema)
	check(err)
	results := []EnumType{}
	for _, enumType := range dbEnums {
		enumValues, err := loader.EnumValueList(generator.db, schema, enumType.EnumName)
		check(err)
		results = append(results, EnumType{
			Name:   enumType.EnumName,
			Values: enumValues,
		})
	}
	return results
}

func (generator *EnumGenerator) getEnumTypesFromReferenceTables() []EnumType {
	loader := postgres.NewPostgresLoader()
	schema, err := loader.Schema()
	check(err)
	tables, err := loader.TableList(generator.db, schema)
	check(err)
	results := []EnumType{}
	for _, table := range tables {
		if !strings.HasSuffix(table.TableName, referenceTableSuffix) {
			continue
		}
		enumValues, err := loader.ReferenceTableValueList(generator.db, schema, table.TableName)
		check(err)
		if len(enumValues) == 0 {
			fmt.Printf("cannot generate types for %s because there are no entries at the moment\n", table.TableName)
			continue
		}
		results = append(results, EnumType{
			Name:                    strings.ReplaceAll(table.TableName, referenceTableSuffix, ""),
			Values:                  enumValues,
			IsReferenceTable:        true,
			ReferenceTableModelType: snaker.SnakeToCamel(table.TableName),
		})
	}
	return results
}
