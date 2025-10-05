package main

import (
	"fmt"
	"os"

	"github.com/carlosrabelo/bip38cli/core/internal/app/cli"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	// Set version info so cobra show right number
	cli.SetVersionInfo(version, buildTime)

	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
