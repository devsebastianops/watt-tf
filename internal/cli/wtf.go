package cli

import (
	"github.com/devsebastianops/watt-tf/internal/logger"
	"github.com/spf13/cobra"
)

var PersistentFlags struct {
	Verbose   bool
	LogFormat string
}

var persistentFlags = &PersistentFlags

var RootCmd = &cobra.Command{
	Use:   "wtf",
	Short: "Watt TF (wtf) is a tool for building Terraform configurations from structured input.",
	Long: `Watt TF (wtf) is a command-line tool that allows you to build Terraform configurations
from structured input files (like JSON or YAML) using a configuration file that defines transformations.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.SetUp(persistentFlags.Verbose, persistentFlags.LogFormat)
	},
}

func init() {
	RootCmd.PersistentFlags().BoolVar(&persistentFlags.Verbose, "verbose", false, "Enable verbose output")
	RootCmd.PersistentFlags().StringVar(&persistentFlags.LogFormat, "log-format", "pretty", "Set log format ('pretty', 'text' or 'json')")
	RootCmd.AddCommand(buildCmd)
}

func Run() error {
	return RootCmd.Execute()
}
