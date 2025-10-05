package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/bip38cli/internal/bip38"
	"github.com/spf13/cobra"
)

var intermediateCmd = &cobra.Command{
	Use:   "intermediate",
	Short: "Generate BIP38 intermediate passphrase codes",
	Long: `Generate BIP38 intermediate passphrase codes for two-factor encryption.

Intermediate codes allow a third party to generate encrypted private keys
without knowing the passphrase. This enables secure key generation in
scenarios where the key generator should not have access to the passphrase.

The generated intermediate code can be used with the 'generate' command
to create encrypted private keys.`,
}

var generateIntermediateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate an intermediate passphrase code",
	Long: `Generate a BIP38 intermediate passphrase code from a passphrase.

You can optionally specify lot and sequence numbers for additional entropy.
If lot and sequence numbers are provided, they will be encoded into the
intermediate code.

Examples:
  bip38cli intermediate generate
  bip38cli intermediate generate --lot 123 --sequence 456`,
	RunE: runGenerateIntermediate,
}

var validateIntermediateCmd = &cobra.Command{
	Use:   "validate [INTERMEDIATE_CODE]",
	Short: "Validate an intermediate passphrase code",
	Long: `Validate the format and integrity of a BIP38 intermediate passphrase code.

If no intermediate code is provided as an argument, you will be prompted to enter it.

Examples:
  bip38cli intermediate validate passphraseabc123...
  bip38cli intermediate validate`,
	Args: cobra.MaximumNArgs(1),
	RunE: runValidateIntermediate,
}

var (
	lotNumber      uint32
	sequenceNumber uint32
	useLotSeq      bool
)

func init() {
	rootCmd.AddCommand(intermediateCmd)
	intermediateCmd.AddCommand(generateIntermediateCmd)
	intermediateCmd.AddCommand(validateIntermediateCmd)

	generateIntermediateCmd.Flags().Uint32Var(&lotNumber, "lot", 0, "lot number (0-1048575)")
	generateIntermediateCmd.Flags().Uint32Var(&sequenceNumber, "sequence", 0, "sequence number (0-4095)")
	generateIntermediateCmd.Flags().BoolVar(&useLotSeq, "use-lot-sequence", false, "use lot and sequence numbers")
}

func runGenerateIntermediate(cmd *cobra.Command, args []string) error {
	// Validate lot and sequence numbers
	var lot, seq *uint32
	if useLotSeq || cmd.Flag("lot").Changed || cmd.Flag("sequence").Changed {
		if lotNumber > 1048575 {
			return fmt.Errorf("lot number must be between 0 and 1048575")
		}
		if sequenceNumber > 4095 {
			return fmt.Errorf("sequence number must be between 0 and 4095")
		}
		lot = &lotNumber
		seq = &sequenceNumber
	}

	// Get passphrase
	passphrase, err := getPassphrase("Enter passphrase: ")
	if err != nil {
		return fmt.Errorf("failed to read passphrase: %v", err)
	}

	if len(passphrase) == 0 {
		return fmt.Errorf("passphrase cannot be empty")
	}

	// Generate intermediate code
	intermediate, err := bip38.GenerateIntermediateCode(passphrase, lot, seq)
	if err != nil {
		return fmt.Errorf("failed to generate intermediate code: %v", err)
	}

	// Output result
	fmt.Printf("Intermediate code: %s\n", intermediate)

	// Show additional info in verbose mode
	if cmd.Flag("verbose").Changed {
		if lot != nil && seq != nil {
			fmt.Printf("Lot number: %d\n", *lot)
			fmt.Printf("Sequence number: %d\n", *seq)
		} else {
			fmt.Println("Type: No lot/sequence")
		}
	}

	return nil
}

func runValidateIntermediate(cmd *cobra.Command, args []string) error {
	// Get intermediate code
	var intermediateCode string
	if len(args) > 0 {
		intermediateCode = args[0]
	} else {
		fmt.Print("Enter intermediate code: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			intermediateCode = strings.TrimSpace(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("failed to read intermediate code: %v", err)
		}
	}

	if intermediateCode == "" {
		return fmt.Errorf("intermediate code is required")
	}

	// Validate format
	if !bip38.IsValidIntermediateCode(intermediateCode) {
		return fmt.Errorf("invalid intermediate code format")
	}

	// Parse the code
	parsed, err := bip38.ParseIntermediateCode(intermediateCode)
	if err != nil {
		return fmt.Errorf("failed to parse intermediate code: %v", err)
	}

	// Output validation result
	fmt.Println("✓ Valid intermediate code")

	// Show details
	fmt.Printf("Type: ")
	if parsed.HasLotSeq {
		fmt.Printf("With lot/sequence\n")
		fmt.Printf("Lot number: %d\n", *parsed.LotNumber)
		fmt.Printf("Sequence number: %d\n", *parsed.SeqNumber)
	} else {
		fmt.Printf("No lot/sequence\n")
	}

	// Show additional info in verbose mode
	if cmd.Flag("verbose").Changed {
		fmt.Printf("Owner salt: %x\n", parsed.OwnerSalt)
		fmt.Printf("Pass point: %x\n", parsed.PassPoint)
	}

	return nil
}