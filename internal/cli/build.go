package cli

import (
	"github.com/devsebastianops/watt-tf/internal/config"
	"github.com/devsebastianops/watt-tf/internal/environment"
	"github.com/devsebastianops/watt-tf/internal/logger"
	"github.com/devsebastianops/watt-tf/internal/parser"
	"github.com/devsebastianops/watt-tf/internal/plugin"
	"github.com/devsebastianops/watt-tf/internal/schema"
	"github.com/devsebastianops/watt-tf/internal/transformer"
	"github.com/devsebastianops/watt-tf/internal/writer"
	"github.com/spf13/cobra"
)

type BuildOptions struct {
	ConfigFile    string
	BlueprintFile string
	InputFile     string
	OutputFile    string
	SchemaFile    string
	StripNulls    bool
	Strict        bool
}

var buildOptions = BuildOptions{}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the project",
	Long:  "Build the project using the specified configuration and input files.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return build()
	},
}

func init() {
	buildCmd.Flags().StringVarP(&buildOptions.ConfigFile, "config", "c", "", "Path to the configuration file  (deprecated, use --blueprint instead)")
	buildCmd.Flags().StringVarP(&buildOptions.BlueprintFile, "blueprint", "b", "blueprint.yaml", "Path to the blueprint YAML file")
	buildCmd.Flags().StringVarP(&buildOptions.InputFile, "input", "i", "", "Path to the input file")
	buildCmd.Flags().StringVarP(&buildOptions.OutputFile, "output", "o", "", "Path to the output file")
	buildCmd.Flags().StringVarP(&buildOptions.SchemaFile, "schema", "s", "", "Path to the schema file")
	buildCmd.Flags().BoolVar(&buildOptions.StripNulls, "strip-nulls", false, "Strip null values from the output")
	buildCmd.Flags().BoolVar(&buildOptions.Strict, "strict", false, "Enable strict mode")
}

func build() error {

	logger.SetUp(persistentFlags.Verbose)

	logger.Debug("building project",
		"config", buildOptions.ConfigFile,
		"blueprint", buildOptions.BlueprintFile,
		"input", buildOptions.InputFile,
		"output", buildOptions.OutputFile,
		"strict", buildOptions.Strict,
		"schema", buildOptions.SchemaFile,
		"verbose", persistentFlags.Verbose)

	var loadPath string
	// Check if user passed config file, which is deprecated
	if buildOptions.ConfigFile != "" {
		logger.Warn("the --config flag is deprecated, please use --blueprint instead")
		loadPath = buildOptions.ConfigFile
	}

	if buildOptions.BlueprintFile != "blueprint.yaml" && buildOptions.ConfigFile != "" {
		logger.Warn("both --config and --blueprint flags are set; --config will be ignored")
	}

	if buildOptions.BlueprintFile != "blueprint.yaml" {
		loadPath = buildOptions.BlueprintFile
	}

	config, configErr := config.LoadConfig(loadPath)
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
		validationErr := schema.ValidateInputSchema(input, buildOptions.SchemaFile)
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
		BasePath:    loadPath,
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
		BasePath:    loadPath,
		Result:      result,
	}
	contextAfter, err := plugin.DispatchEvents(dispatchConfigAfter)
	if err != nil {
		return err
	}

	result = contextAfter.Result

	if buildOptions.StripNulls {
		result = transformer.StripNullValues(result)
	}

	logger.Debug("transformation completed successfully")

	writeErr := writer.WriteJSON(result, buildOptions.OutputFile)
	if writeErr != nil {
		return writeErr
	}

	return nil
}
