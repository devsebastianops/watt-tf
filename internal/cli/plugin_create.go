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

var (
	pluginCreateOutputFile  string
	pluginCreateName        string
	pluginCreateDescription string
	pluginCreateAuthor      string
	pluginCreateVersion     string
	pluginCreateLanguage    string
)

func init() {
	pluginCreateCmd.Flags().StringVarP(&pluginCreateOutputFile, "output", "o", ".", "Where to create the plugin file (default: current directory)")
	pluginCreateCmd.Flags().StringVarP(&pluginCreateName, "name", "n", "", "Name of the plugin (required)")
	pluginCreateCmd.Flags().StringVarP(&pluginCreateDescription, "description", "d", "", "Description of the plugin (optional)")
	pluginCreateCmd.Flags().StringVarP(&pluginCreateAuthor, "author", "a", "", "Author of the plugin (optional)")
	pluginCreateCmd.Flags().StringVarP(&pluginCreateVersion, "version", "V", "", "Version of the plugin (optional)")
	pluginCreateCmd.Flags().StringVarP(&pluginCreateLanguage, "language", "l", "", "Programming language for the plugin (optional)")
}

func pluginCreateHandler() error {
	fmt.Println("Creating plugin with the following details:")
	fmt.Println("Output File:", pluginCreateOutputFile)
	fmt.Println("Name:", pluginCreateName)
	fmt.Println("Description:", pluginCreateDescription)
	fmt.Println("Author:", pluginCreateAuthor)
	fmt.Println("Version:", pluginCreateVersion)
	fmt.Println("Language:", pluginCreateLanguage)
	return nil
}
