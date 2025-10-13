package cli

import (
	"fmt"

	"github.com/carlosrabelo/bip38cli/core/internal/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	version   string
	buildTime string
)

// rootCmd serves as base command when user doesn't pass subcommands
var rootCmd = &cobra.Command{
	Use:   "bip38cli",
	Short: "A CLI tool for BIP38 Bitcoin private key encryption",
	Long: `bip38cli is a command-line tool that implements BIP38 (Bitcoin Improvement Proposal 38)
for encrypting and decrypting Bitcoin private keys with passphrases.

Features:
- Encrypt/decrypt Bitcoin private keys using BIP38 standard
- Generate intermediate passphrase codes for two-factor encryption
- Support for both compressed and uncompressed keys
- Secure passphrase handling`,
	Version: getVersionString(),
}

// Execute attaches all child commands to the root command and executes it.
// This is the main entry point for the CLI application.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	// Initialize logger with default settings
	logger.Init(false)
}

// SetVersionInfo sets the version information for the CLI.
// This should be called before Execute() to ensure version info is available.
func SetVersionInfo(v, bt string) {
	version = v
	buildTime = bt
	rootCmd.Version = getVersionString()
}

func getVersionString() string {
	if buildTime != "unknown" {
		return fmt.Sprintf("%s (built: %s)", version, buildTime)
	}
	return version
}
