package generator

import (
	"fmt"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/lumina-tech/gooq/pkg/database"
	"github.com/lumina-tech/gooq/pkg/generator"
	"github.com/spf13/cobra"
)

var generateDatabaseModelCommandUseDocker bool
var generateDatabaseModelCommand = &cobra.Command{
	Use:   "generate-database-model",
	Short: "generate Go models by introspecting the database",
	Run: func(cmd *cobra.Command, args []string) {

		config := database.DatabaseConfig{
			Host:          "localhost",
			Port:          5432,
			Username:      "swapi",
			Password:      "swapi",
			Version:       "11.4",
			DatabaseName:  "swapi",
			SSLMode:       "disable",
			MigrationPath: "examples/swapi/migrations",
			ModelPath:     "examples/swapi/model",
			TablePath:     "examples/swapi/table",
		}
		if generateDatabaseModelCommandUseDocker {
			db := database.NewDockerizedDB(&config)
			defer db.Close()
			database.MigrateDatabase(db.DB.DB, config.MigrationPath)
			generateModelsForDB(db.DB, &config)
		} else {
			db := database.NewDatabase(&config)
			generateModelsForDB(db, &config)
		}
	},
}

func generateModelsForDB(
	db *sqlx.DB, config *database.DatabaseConfig,
) {
	templateDir := "./pkg/generator/templates"
	enumTemplatePath := filepath.Join(templateDir, generator.EnumTemplateFilename)
	generator.NewEnumGenerator(db, enumTemplatePath, config.ModelPath, config.DatabaseName).Run()

	modelTemplatePath := filepath.Join(templateDir, generator.ModelTemplateFilename)
	generatedModelFilename := fmt.Sprintf("%s/%s_model.generated.go", config.ModelPath, config.DatabaseName)
	generator.GenerateModel(db, modelTemplatePath, generatedModelFilename, config.DatabaseName)

	tableTemplatePath := filepath.Join(templateDir, generator.TableTemplateFilename)
	generatedTableFilename := fmt.Sprintf("%s/%s_table.generated.go", config.TablePath, config.DatabaseName)
	generator.GenerateModel(db, tableTemplatePath, generatedTableFilename, config.DatabaseName)

}
