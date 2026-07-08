package cli

import (
	"os"
	"strings"

	"github.com/devsebastianops/watt-tf/internal/config"
	"github.com/devsebastianops/watt-tf/internal/logger"
	"github.com/devsebastianops/watt-tf/internal/parser"
	"github.com/devsebastianops/watt-tf/internal/plugin"
	"github.com/devsebastianops/watt-tf/internal/transformer"
	"github.com/devsebastianops/watt-tf/internal/writer"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the project",
	Long:  `Build the project using the specified configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return buildHandler()
	},
}

var (
	buildConfigFile string
	buildInputFile  string
	buildOutputFile string
	buildStrict     bool
	buildSchemaFile string
	buildVerbose    bool
)

func init() {
	buildCmd.Flags().StringVarP(&buildConfigFile, "config", "c", ".wtf.yaml", "Path to the configuration file")
	buildCmd.Flags().StringVarP(&buildInputFile, "input", "i", "", "Path to the input file")
	buildCmd.Flags().StringVarP(&buildOutputFile, "output", "o", "watt.tf.json", "Path to the output file")
	buildCmd.Flags().BoolVarP(&buildStrict, "strict", "s", false, "Fail on missing keys (default: false = missing keys replaced with null)")
	buildCmd.Flags().StringVarP(&buildSchemaFile, "schema", "S", "", "Path to JSON Schema file for input validation (optional)")
	buildCmd.Flags().BoolVarP(&buildVerbose, "verbose", "v", false, "Enable verbose output")
}

func buildHandler() error {
	logger.SetUp(buildVerbose)
	logger.Debug("building project",
		"config", buildConfigFile,
		"input", buildInputFile,
		"output", buildOutputFile,
		"strict", buildStrict,
		"schema", buildSchemaFile,
		"verbose", buildVerbose)

	config, configErr := config.LoadConfig(buildConfigFile)
	if configErr != nil {
		return configErr
	}

	logger.Debug("configuration loaded", "transform_count", len(config.Transform))
	logger.Debug("Plugins", "plugins", config.Plugins)

	input, inputErr := parser.ParseInput(buildInputFile)
	if inputErr != nil {
		return inputErr
	}

	envVars := getEnvVars()

	logger.Debug("input parsed successfully", "input_keys", len(input))

	// Validate input against schema if provided
	if buildSchemaFile != "" {
		logger.Info("validating input against schema", "schema", buildSchemaFile)
		validationErr := validateInputSchema(input, buildSchemaFile)
		if validationErr != nil {
			return validationErr
		}
		logger.Debug("input validation passed")
	}

	registry := plugin.NewRegistry()
	registry.RegisterPlugins(config.Plugins)

	dispatchConfig := plugin.DispatchConfig{
		Event:       plugin.EventBeforeTransform,
		Registry:    registry,
		Input:       input,
		Environment: envVars,
		BasePath:    buildConfigFile,
		Result:      nil,
	}
	context, err := plugin.DispatchEvents(dispatchConfig)
	if err != nil {
		return err
	}

	result, transformErr := transformer.Transform(context.Input, envVars, config, buildStrict)
	if transformErr != nil {
		return transformErr
	}

	dispatchConfigAfter := plugin.DispatchConfig{
		Event:       plugin.EventAfterTransform,
		Registry:    registry,
		Input:       input,
		Environment: envVars,
		BasePath:    buildConfigFile,
		Result:      result,
	}
	contextAfter, err := plugin.DispatchEvents(dispatchConfigAfter)
	if err != nil {
		return err
	}

	result = contextAfter.Result

	logger.Debug("transformation completed successfully")

	writeErr := writer.WriteJSON(result, buildOutputFile)
	if writeErr != nil {
		return writeErr
	}

	return nil
}

func getEnvVars() map[string]string {
	envMap := make(map[string]string)
	for _, envVar := range os.Environ() {
		key, value, _ := strings.Cut(envVar, "=")
		envMap[key] = value
	}
	return envMap
}
