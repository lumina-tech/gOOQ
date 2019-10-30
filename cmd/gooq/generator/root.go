package generator

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd    = &cobra.Command{}
	configFile string
)

func init() {
	generateDatabaseModelCommand.PersistentFlags().BoolVarP(
		&generateDatabaseModelCommandUseDocker, "docker", "d", true, "whether to use dockerize db")
	generateDatabaseModelCommand.PersistentFlags().StringVarP(
		&generateDatabaseModelConfigFilePath, "config-file", "f", "", "path to configuration file")
	rootCmd.AddCommand(generateDatabaseModelCommand)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
