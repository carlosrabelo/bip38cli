package main

import (
	"fmt"
	"os"

	"github.com/yourusername/bip38cli/cmd/bip38cli/cmd"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	// Set version info for cobra
	cmd.SetVersionInfo(version, buildTime)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}