package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wtf",
	Short: "Watt Terraform (wtf) is a tool to generate Terraform configurations from structured input.",
	Long: `Watt Terraform (wtf) is a command-line tool that allows you to generate Terraform configurations
from structured input files. It supports various plugins and transformations to customize the output.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")

	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(pluginCmd)
}
