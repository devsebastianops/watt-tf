package cli

import (
	"github.com/devsebastianops/watt-tf/internal/config"
	"github.com/devsebastianops/watt-tf/internal/environment"
	"github.com/devsebastianops/watt-tf/internal/logger"
	"github.com/devsebastianops/watt-tf/internal/parser"
	"github.com/devsebastianops/watt-tf/internal/plugin"
	"github.com/devsebastianops/watt-tf/internal/transformer"
	"github.com/devsebastianops/watt-tf/internal/validator"
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

type BuildOptions struct {
	ConfigFile string
	InputFile  string
	OutputFile string
	Strict     bool
	SchemaFile string
	Verbose    bool
}

var buildOptions = BuildOptions{}

func init() {
	buildCmd.Flags().StringVarP(&buildOptions.ConfigFile, "config", "c", ".wtf.yaml", "Path to the configuration file")
	buildCmd.Flags().StringVarP(&buildOptions.InputFile, "input", "i", "", "Path to the input file")
	buildCmd.Flags().StringVarP(&buildOptions.OutputFile, "output", "o", "watt.tf.json", "Path to the output file")
	buildCmd.Flags().BoolVarP(&buildOptions.Strict, "strict", "s", false, "Fail on missing keys (default: false = missing keys replaced with null)")
	buildCmd.Flags().StringVarP(&buildOptions.SchemaFile, "schema", "S", "", "Path to JSON Schema file for input validation (optional)")
	buildCmd.Flags().BoolVarP(&buildOptions.Verbose, "verbose", "v", false, "Enable verbose output")
}

func buildHandler() error {
	logger.SetUp(buildOptions.Verbose)
	logger.Debug("building project",
		"config", buildOptions.ConfigFile,
		"input", buildOptions.InputFile,
		"output", buildOptions.OutputFile,
		"strict", buildOptions.Strict,
		"schema", buildOptions.SchemaFile,
		"verbose", buildOptions.Verbose)

	config, configErr := config.LoadConfig(buildOptions.ConfigFile)
	if configErr != nil {
		return configErr
	}

	logger.Debug("configuration loaded", "transform_count", len(config.Transform))
	logger.Debug("Plugins", "plugins", config.Plugins)

	input, inputErr := parser.ParseInput(buildOptions.InputFile)
	if inputErr != nil {
		return inputErr
	}

	envVars := environment.GetEnvVars()

	logger.Debug("input parsed successfully", "input_keys", len(input))

	// Validate input against schema if provided
	if buildOptions.SchemaFile != "" {
		logger.Info("validating input against schema", "schema", buildOptions.SchemaFile)
		validationErr := validator.ValidateInputSchema(input, buildOptions.SchemaFile)
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
		BasePath:    buildOptions.ConfigFile,
		Result:      nil,
	}
	context, err := plugin.DispatchEvents(dispatchConfig)
	if err != nil {
		return err
	}

	result, transformErr := transformer.Transform(context.Input, envVars, config, buildOptions.Strict)
	if transformErr != nil {
		return transformErr
	}

	dispatchConfigAfter := plugin.DispatchConfig{
		Event:       plugin.EventAfterTransform,
		Registry:    registry,
		Input:       input,
		Environment: envVars,
		BasePath:    buildOptions.ConfigFile,
		Result:      result,
	}
	contextAfter, err := plugin.DispatchEvents(dispatchConfigAfter)
	if err != nil {
		return err
	}

	result = contextAfter.Result

	logger.Debug("transformation completed successfully")

	writeErr := writer.WriteJSON(result, buildOptions.OutputFile)
	if writeErr != nil {
		return writeErr
	}

	return nil
}
