package main

import (
	"os"

	"github.com/devsebastianops/watt-tf/internal/cli"
	"github.com/devsebastianops/watt-tf/internal/logger"
)

func main() {
	err := cli.Run()
	if err != nil {

		logger.Error("wtf command failed", "error", err.Error())
		os.Exit(1)
	}
}
