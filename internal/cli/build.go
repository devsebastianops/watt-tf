package cli

import (
	"flag"
	"os"

	"github.com/devsebastianops/watt-tf/internal/config"
	"github.com/devsebastianops/watt-tf/internal/logger"
	"github.com/devsebastianops/watt-tf/internal/parser"
	"github.com/devsebastianops/watt-tf/internal/transformer"
	"github.com/devsebastianops/watt-tf/internal/writer"
)

var buildCmd = flag.NewFlagSet("build", flag.ExitOnError)
var buildConfigFile = buildCmd.String("config", ".wtf.yaml", "Path to the configuration file")
var buildInputFile = buildCmd.String("input", "", "Path to the input file")
var buildOutputFile = buildCmd.String("output", "watt.tf.json", "Path to the output file")
var buildStrict = buildCmd.Bool("strict", false, "Fail on missing keys (default: false = missing keys replaced with null)")

// var buildVerbose = buildCmd.Bool("verbose", false, "Enable verbose output") - For later

func build() error {
	buildCmd.Parse(os.Args[2:])
	logger.Info("building project",
		"config", *buildConfigFile,
		"input", *buildInputFile,
		"output", *buildOutputFile,
		"strict", *buildStrict)

	config, configErr := config.LoadConfig(*buildConfigFile)
	if configErr != nil {
		return configErr
	}

	logger.Debug("configuration loaded", "transform_count", len(config.Transform))

	input, inputErr := parser.ParseInput(*buildInputFile)
	if inputErr != nil {
		return inputErr
	}

	logger.Debug("input parsed successfully", "input_keys", len(input))

	result, transformErr := transformer.Transform(input, config, *buildStrict)
	if transformErr != nil {
		return transformErr
	}

	logger.Debug("transformation completed successfully")

	writeErr := writer.WriteJSON(result, *buildOutputFile)
	if writeErr != nil {
		return writeErr
	}

	return nil
}
