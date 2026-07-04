package cli

import (
	"errors"
	"os"
)

func getSubcommand() (string, error) {
	if len(os.Args) < 2 {
		help()
		return "", errors.New("expected subcommand")
	}

	return os.Args[1], nil
}

func getArgs() []string {
	return os.Args[2:]
}
