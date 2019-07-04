package cmd

import (
	"github.com/spf13/cobra"
)

var generateModelsCmd = &cobra.Command{
	Use:   "generate-models",
	Short: "Generates database models",
	Run: func(cmd *cobra.Command, args []string) {
		// generator.GenerateModel()
	},
}
