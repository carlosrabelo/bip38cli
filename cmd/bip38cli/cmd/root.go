package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	version string
	buildTime string
)

// rootCmd represents the base command when called without any subcommands
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

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bip38cli.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".bip38cli")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil && viper.GetBool("verbose") {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// SetVersionInfo sets the version information for the CLI
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