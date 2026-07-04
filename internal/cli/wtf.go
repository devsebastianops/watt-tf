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
	fmt.Println("Commands:")
	fmt.Println("  build    Build the project")
	fmt.Println("  help     Show this help message")
}
