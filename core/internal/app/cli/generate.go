package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/carlosrabelo/bip38cli/core/internal/domain/bip38"
	"github.com/carlosrabelo/bip38cli/core/internal/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	generateShowAddress bool
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new Bitcoin private key",
	Long: `Generate a new Bitcoin private key in Wallet Import Format (WIF).

The generated key can be used immediately with other commands:
  bip38cli generate | xargs bip38cli encrypt

By default, compressed format is used which is the modern standard.
Use --uncompressed for legacy compatibility.`,
	RunE: runGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Local flags
	generateCmd.Flags().BoolVar(&generateShowAddress, "show-address", false, "show the Bitcoin address for the generated key")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	outputFormat, _ := cmd.Flags().GetString("output-format")
	compressed, _ := cmd.Flags().GetBool("compressed")

	if verbose {
		logger.Info("Generating new Bitcoin private key...")
	}

	// Generate the key
	generatedKey, err := bip38.GenerateKey(compressed)
	if err != nil {
		return fmt.Errorf("failed to generate key: %v", err)
	}

	// Output based on format
	switch outputFormat {
	case "json":
		return outputJSON(generatedKey)
	default:
		return outputText(generatedKey)
	}
}

func outputText(key *bip38.GeneratedKey) error {
	fmt.Println(key.PrivateKey)

	if generateShowAddress {
		fmt.Printf("Address: %s\n", key.Address)
	}

	return nil
}

func outputJSON(key *bip38.GeneratedKey) error {
	output := map[string]interface{}{
		"private_key": key.PrivateKey,
		"compressed":  key.Compressed,
	}

	if generateShowAddress {
		output["address"] = key.Address
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "")
	return encoder.Encode(output)
}
