package cli

import (
	"fmt"
	"os"

	"github.com/carlosrabelo/bip38cli/core/internal/infra/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	version   string
	buildTime string
)

// rootCmd stay as base command when user no pass subcommands
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

// Execute attach every child command and run root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags for config and noise level
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bip38cli.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	// Bind flags to viper for shared usage
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig read config file and ENV variables when present
func initConfig() {
	usedConfig, err := config.Setup(viper.GetViper(), cfgFile)
	cobra.CheckErr(err)

	viper.AutomaticEnv()

	if usedConfig != "" && viper.GetBool("verbose") {
		fmt.Fprintln(os.Stderr, "Using config file:", usedConfig)
	}
}

// SetVersionInfo set the version information for the CLI
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
