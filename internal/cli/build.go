package cli

import (
	"flag"
	"os"
	"strings"

	"github.com/devsebastianops/watt-tf/internal/config"
	"github.com/devsebastianops/watt-tf/internal/logger"
	"github.com/devsebastianops/watt-tf/internal/parser"
	"github.com/devsebastianops/watt-tf/internal/plugin"
	"github.com/devsebastianops/watt-tf/internal/transformer"
	"github.com/devsebastianops/watt-tf/internal/writer"
)

var buildCmd = flag.NewFlagSet("build", flag.ExitOnError)
var buildConfigFile = buildCmd.String("config", ".wtf.yaml", "Path to the configuration file")
var buildInputFile = buildCmd.String("input", "", "Path to the input file")
var buildOutputFile = buildCmd.String("output", "watt.tf.json", "Path to the output file")
var buildStrict = buildCmd.Bool("strict", false, "Fail on missing keys (default: false = missing keys replaced with null)")
var buildSchemaFile = buildCmd.String("schema", "", "Path to JSON Schema file for input validation (optional)")
var buildVerbose = buildCmd.Bool("verbose", false, "Enable verbose output")
var buildStripNulls = buildCmd.Bool("strip-nulls", false, "Strip null values from the output (default: false)")

func build() error {
	buildCmd.Parse(os.Args[2:])

	logger.SetUp(*buildVerbose)

	logger.Debug("building project",
		"config", *buildConfigFile,
		"input", *buildInputFile,
		"output", *buildOutputFile,
		"strict", *buildStrict,
		"schema", *buildSchemaFile,
		"verbose", *buildVerbose)

	config, configErr := config.LoadConfig(*buildConfigFile)
	if configErr != nil {
		return configErr
	}

	logger.Debug("configuration loaded", "transform_count", len(config.Transform))
	logger.Debug("Plugins", "plugins", config.Plugins)

	input, inputErr := parser.ParseInput(*buildInputFile)
	if inputErr != nil {
		return inputErr
	}

	envVars := getEnvVars()

	logger.Debug("input parsed successfully", "input_keys", len(input))

	// Validate input against schema if provided
	if *buildSchemaFile != "" {
		logger.Info("validating input against schema", "schema", *buildSchemaFile)
		validationErr := validateInputSchema(input, *buildSchemaFile)
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
		BasePath:    *buildConfigFile,
		Result:      nil,
	}
	context, err := plugin.DispatchEvents(dispatchConfig)
	if err != nil {
		return err
	}

	result, transformErr := transformer.Transform(context.Input, envVars, config, *buildStrict)
	if transformErr != nil {
		return transformErr
	}

	dispatchConfigAfter := plugin.DispatchConfig{
		Event:       plugin.EventAfterTransform,
		Registry:    registry,
		Input:       input,
		Environment: envVars,
		BasePath:    *buildConfigFile,
		Result:      result,
	}
	contextAfter, err := plugin.DispatchEvents(dispatchConfigAfter)
	if err != nil {
		return err
	}

	result = contextAfter.Result

	if *buildStripNulls {
		result = transformer.StripNullValues(result)
	}

	logger.Debug("transformation completed successfully")

	writeErr := writer.WriteJSON(result, *buildOutputFile)
	if writeErr != nil {
		return writeErr
	}

	return nil
}

// getEnvVars collects all environment variables into a map
func getEnvVars() map[string]string {
	envMap := make(map[string]string)
	for _, envVar := range os.Environ() {
		key, value, _ := strings.Cut(envVar, "=")
		envMap[key] = value
	}
	return envMap
}
