package generator

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/lumina-tech/gooq/pkg/database"
	"github.com/lumina-tech/gooq/pkg/generator"
	"github.com/lumina-tech/gooq/pkg/generator/plugin/enumgen"
	"github.com/lumina-tech/gooq/pkg/generator/plugin/modelgen"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	generateDatabaseModelCommandUseDocker bool
	generateDatabaseModelConfigFilePath   string
)

var generateDatabaseModelCommand = &cobra.Command{
	Use:   "generate-database-model",
	Short: "generate Go models by introspecting the database",
	Run: func(cmd *cobra.Command, args []string) {
		if err := loadConfig(); err != nil {
			_, _ = fmt.Fprint(os.Stderr, "cannot read configuration file:", err)
			os.Exit(1)
		}

		var config database.DatabaseConfig
		if err := viper.Unmarshal(&config); err != nil {
			_, _ = fmt.Fprint(os.Stderr, "cannot decode configuration file:", err)
			os.Exit(1)
		}
		if generateDatabaseModelCommandUseDocker {
			db := database.NewDockerizedDB(&config, viper.GetString("dockerTag"))
			defer db.Close()
			database.MigrateDatabase(db.DB.DB, config.MigrationPath)
			generateModelsForDB(db.DB, &config)
		} else {
			db := database.NewDatabase(&config)
			generateModelsForDB(db, &config)
		}
	},
}

func loadConfig() error {
	viper.SetDefault("dockerTag", "11.4-alpine")
	if len(generateDatabaseModelConfigFilePath) != 0 {
		viper.SetConfigFile(generateDatabaseModelConfigFilePath)
		return viper.ReadInConfig()
	}
	viper.SetConfigName("gooq")
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}
	viper.AddConfigPath(wd)
	return viper.ReadInConfig()
}

func generateModelsForDB(
	db *sqlx.DB, config *database.DatabaseConfig,
) {
	enumOutputFile := fmt.Sprintf("%s/%s_enum.generated.go", config.ModelPath, config.DatabaseName)
	modelOutputFile := fmt.Sprintf("%s/%s_model.generated.go", config.ModelPath, config.DatabaseName)
	tableOutputFile := fmt.Sprintf("%s/%s_table.generated.go", config.TablePath, config.DatabaseName)
	err := generator.NewGenerator(
		enumgen.NewEnumGenerator(enumOutputFile),
		modelgen.NewModelGenerator(modelOutputFile, "table", "model", &config.ModelOverrides),
		modelgen.NewTableGenerator(tableOutputFile, "table", "model", nil),
	).Run(db)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "cannot generate code:", err)
		os.Exit(1)
	}
}
