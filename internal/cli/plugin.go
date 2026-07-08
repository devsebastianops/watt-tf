package cli

import "github.com/spf13/cobra"

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage plugins",
	Long:  `Commands to create and manage plugins.`,
}

func init() {
	pluginCmd.AddCommand(pluginCreateCmd)
}
