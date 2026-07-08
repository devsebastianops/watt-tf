package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pluginCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new plugin",
	Long:  `Create a new plugin with the specified name, description, author, version, and programming language.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return pluginCreateHandler()
	},
}

type PluginCreateOptions struct {
	OutputFile  string
	Name        string
	Description string
	Author      string
	Version     string
	Language    string
}

var pluginCreateOptions = PluginCreateOptions{}

func init() {
	pluginCreateCmd.Flags().StringVarP(&pluginCreateOptions.OutputFile, "output", "o", ".", "Where to create the plugin file (default: current directory)")
	pluginCreateCmd.Flags().StringVarP(&pluginCreateOptions.Name, "name", "n", "", "Name of the plugin (required)")
	pluginCreateCmd.Flags().StringVarP(&pluginCreateOptions.Description, "description", "d", "", "Description of the plugin (optional)")
	pluginCreateCmd.Flags().StringVarP(&pluginCreateOptions.Author, "author", "a", "", "Author of the plugin (optional)")
	pluginCreateCmd.Flags().StringVarP(&pluginCreateOptions.Version, "version", "V", "", "Version of the plugin (optional)")
	pluginCreateCmd.Flags().StringVarP(&pluginCreateOptions.Language, "language", "l", "", "Programming language for the plugin (optional)")
}

func pluginCreateHandler() error {
	fmt.Println("Creating plugin with the following details:")
	fmt.Println("Output File:", pluginCreateOptions.OutputFile)
	fmt.Println("Name:", pluginCreateOptions.Name)
	fmt.Println("Description:", pluginCreateOptions.Description)
	fmt.Println("Author:", pluginCreateOptions.Author)
	fmt.Println("Version:", pluginCreateOptions.Version)
	fmt.Println("Language:", pluginCreateOptions.Language)
	return nil
}
