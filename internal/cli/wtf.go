package cli

import (
	"errors"
	"fmt"
)

func Wtf() error {
	subcommand, err := getSubcommand()
	if err != nil {
		return err
	}

	switch subcommand {
	case "build":
		return build()
	case "help":
		help()
		return nil
	default:
		return errors.New("unknown subcommand: " + subcommand)
	}
}

func help() {
	fmt.Println("Usage: wtf <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  build    Build the project")
	fmt.Println("  help     Show this help message")
	fmt.Println()
	fmt.Println("Build Options:")
	fmt.Println("  --config <file>   Path to the configuration file (default: .wtf.yaml)")
	fmt.Println("  --input <file>    Path to the input file (required)")
	fmt.Println("  --output <file>   Path to the output file (default: watt.tf.json)")
	fmt.Println("  --strict          Fail on missing keys (default: false = missing keys → null)")
}
