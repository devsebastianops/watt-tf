package cli

import (
	"github.com/devsebastianops/watt-tf/internal/logger"
	"github.com/spf13/cobra"
)

var (
	Version   = "development"
	Commit    = "none"
	BuildTime = "2026-01-01T00:00:00Z"
)

var PersistentFlags struct {
	Verbose   bool
	LogFormat string
	Silent    bool
}

var persistentFlags = &PersistentFlags

var RootFlags struct {
	Version bool
}

var rootFlags = &RootFlags

var RootCmd = &cobra.Command{
	Use:   "wtf",
	Short: "Watt TF (wtf) is a tool for building Terraform configurations from structured input.",
	Long: `Watt TF (wtf) is a command-line tool that allows you to build Terraform configurations
from structured input files (like JSON or YAML) using a configuration file that defines transformations.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		isSilentAndVerbose := persistentFlags.Silent && persistentFlags.Verbose
		if isSilentAndVerbose {
			logger.Warn("You used both --silent and --verbose flags. Silent mode will take precedence.")
		}
		logger.SetUp(persistentFlags.Verbose, persistentFlags.LogFormat, persistentFlags.Silent)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if rootFlags.Version {
			logger.Infof("wtf version %s (commit: %s, build time: %s)", Version, Commit, BuildTime)
		}
		return nil
	},
}

func init() {
	RootCmd.Flags().BoolVar(&rootFlags.Version, "version", false, "Show the version of wtf")

	RootCmd.PersistentFlags().BoolVar(&persistentFlags.Verbose, "verbose", false, "Enable verbose output")
	RootCmd.PersistentFlags().StringVar(&persistentFlags.LogFormat, "log-format", "pretty", "Set log format ('pretty', 'text' or 'json')")
	RootCmd.PersistentFlags().BoolVar(&persistentFlags.Silent, "silent", false, "Enable silent mode")
	RootCmd.AddCommand(buildCmd)
}

func Run() error {
	return RootCmd.Execute()
}
