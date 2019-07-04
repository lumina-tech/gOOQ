package generator

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/knq/snaker"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"

	"github.com/lumina-tech/gooq/common"
	"github.com/lumina-tech/gooq/loader"
)

const (
	EnumTemplateFilename        = "enum.go.tmpl"
	EnumGraphqlTemplateFilename = "enum.graphql.tmpl"
	referenceTableSuffix        = "_reference_table"
)

type EnumGenerator struct {
	db                  *sqlx.DB
	modelTemplatePath   string
	modelOutputFile     string
	graphqlTemplatePath string
	graphqlOutputFile   string
	dbname              string
}

func NewEnumGenerator(
	db *sqlx.DB,
	modelTemplatePath, modelOutputPath,
	graphqlTemplatePath, graphqlOutputPath,
	dbname string,
) *EnumGenerator {
	graphqlOutputFile := fmt.Sprintf("%s/%s_enum.generated.graphql",
		graphqlOutputPath, dbname)
	modelOutputFile := fmt.Sprintf("%s/%s_enum.generated.go",
		modelOutputPath, dbname)
	return &EnumGenerator{
		db,
		modelTemplatePath,
		modelOutputFile,
		graphqlTemplatePath,
		graphqlOutputFile,
		dbname,
	}
}

func (generator *EnumGenerator) Run() {
	args := generator.getTemplateArguments()
	enumTemplate := getTemplate(generator.graphqlTemplatePath)
	err := common.RenderToFile(enumTemplate, generator.graphqlOutputFile, args)
	check(err)
	enumTemplate = getTemplate(generator.modelTemplatePath)
	err = common.RenderToFile(enumTemplate, generator.modelOutputFile, args)
	check(err)

	fmt.Println()
	fmt.Println("################################################################################")
	fmt.Println("NOTE: copy the following to gqlgen.yml. This part has yet to be automated")
	fmt.Println("################################################################################")
	fmt.Println()
	for _, enum := range args.Enums {
		enumName := snaker.SnakeToCamelIdentifier(enum.Name)
		fmt.Printf("  %s:\n", enumName)
		fmt.Printf("    model: 'github.com/lumina-tech/lumina/apps/server/model.%s'\n", enumName)
	}
	fmt.Println()
}

func (generator *EnumGenerator) getTemplateArguments() EnumsTemplateArgs {
	results := generator.getEnumTypesFromDatabaseTypes()
	results = append(results, generator.getEnumTypesFromReferenceTables()...)

	// get existing enum description mapping from graphql file
	existingDescription := generator.getExistingEnumOptionDescription()
	// copy over description from existing graphql enum option
	for index := range results {
		result := &results[index]
		enumName := snaker.SnakeToCamelIdentifier(result.Name)
		if optionDescriptions, ok := existingDescription[enumName]; ok {
			for index := range result.Values {
				enumValue := &result.Values[index]
				if description, ok := optionDescriptions[enumValue.EnumValue]; ok {
					enumValue.Description = description
				}
			}
		}
	}
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
	loader := loader.NewPostgresLoader()
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
	loader := loader.NewPostgresLoader()
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
			Name:   strings.ReplaceAll(table.TableName, referenceTableSuffix, ""),
			Values: enumValues,
		})
	}
	return results
}

// getExistingEnumOptionDescription creates a mapping from
// enumType -> (enumOption -> description)
func (generator *EnumGenerator) getExistingEnumOptionDescription() map[string]map[string]string {
	bytes, err := ioutil.ReadFile(generator.graphqlOutputFile)
	check(err)
	schema := gqlparser.MustLoadSchema(&ast.Source{
		Name:  "graph/schema/lumina_enum.generated.graphql",
		Input: string(bytes),
	})
	enumsMap := map[string]map[string]string{}
	for _, definition := range schema.Types {
		if definition.Kind != ast.Enum {
			continue
		}
		optionsMap := map[string]string{}
		enumsMap[definition.Name] = optionsMap
		enumDefinitions := []*ast.EnumValueDefinition(definition.EnumValues)
		for _, enumDefinition := range enumDefinitions {
			optionsMap[enumDefinition.Name] = enumDefinition.Description
		}
	}
	return enumsMap
}
