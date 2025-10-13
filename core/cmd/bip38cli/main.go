package main

import (
	"fmt"
	"os"

	"github.com/carlosrabelo/bip38cli/core/internal/app/cli"
	"github.com/carlosrabelo/bip38cli/core/internal/pkg/logger"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	// Set version info so cobra shows the right number
	cli.SetVersionInfo(version, buildTime)

	if err := cli.Execute(); err != nil {
		logger.WithError(err).Error("Application failed")
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
