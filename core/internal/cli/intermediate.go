package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/carlosrabelo/bip38cli/internal/bip38"
	"github.com/carlosrabelo/bip38cli/pkg/errors"
	"github.com/carlosrabelo/bip38cli/pkg/logger"
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

func runGenerateIntermediate(cmd *cobra.Command, _ []string) error {
	// Reinitialize logger with verbose setting if needed
	if isVerbose(cmd) {
		logger.Init(true)
	}

	logger.Debug("Starting intermediate code generation")

	// Validate lot and sequence numbers from flags
	var lot, seq *uint32
	lotChanged := cmd.Flag("lot").Changed
	seqChanged := cmd.Flag("sequence").Changed

	if useLotSeq {
		if !lotChanged || !seqChanged {
			logger.Error("Use-lot-sequence flag set but missing lot or sequence")
			return errors.NewValidationError("both --lot and --sequence must be provided when --use-lot-sequence is set", nil)
		}
	}

	if lotChanged != seqChanged {
		logger.Error("Lot and sequence flags must be provided together")
		return errors.NewValidationError("both --lot and --sequence must be provided together", nil)
	}

	if useLotSeq || lotChanged || seqChanged {
		if lotNumber > 1048575 {
			return fmt.Errorf("lot number must be between 0 and 1048575")
		}
		if sequenceNumber > 4095 {
			return fmt.Errorf("sequence number must be between 0 and 4095")
		}
		lot = &lotNumber
		seq = &sequenceNumber
	}

	// Ask hidden passphrase for intermediate flow
	passphrase, err := getPassphrase("Enter passphrase: ")
	if err != nil {
		return fmt.Errorf("failed to read passphrase: %v", err)
	}
	defer secureZero(passphrase)

	if len(passphrase) == 0 {
		return fmt.Errorf("passphrase cannot be empty")
	}

	// Generate intermediate code with domain layer
	intermediate, err := bip38.GenerateIntermediateCode(passphrase, lot, seq)
	if err != nil {
		logger.WithError(err).Error("Failed to generate intermediate code")
		return errors.NewCryptoError("failed to generate intermediate code", err)
	}

	logger.Info("Successfully generated intermediate code")

	// Prepare output data
	result := map[string]interface{}{
		"intermediate_code": intermediate,
		"has_lot_sequence":  lot != nil && seq != nil,
	}

	if lot != nil && seq != nil {
		result["lot_number"] = *lot
		result["sequence_number"] = *seq
	}

	// Output based on format
	switch outputFormat(cmd) {
	case "json":
		jsonOutput, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON output: %v", err)
		}
		fmt.Println(string(jsonOutput))
	default:
		// Text output
		fmt.Printf("Intermediate code: %s\n", intermediate)

		// Show extra detail when verbose flag up
		if isVerbose(cmd) {
			if lot != nil && seq != nil {
				fmt.Printf("Lot number: %d\n", *lot)
				fmt.Printf("Sequence number: %d\n", *seq)
			} else {
				fmt.Println("Type: No lot/sequence")
			}
		}
	}

	return nil
}

func runValidateIntermediate(cmd *cobra.Command, args []string) error {
	// Grab intermediate code input
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

	// Validate format to reject bogus code
	if !bip38.IsValidIntermediateCode(intermediateCode) {
		return fmt.Errorf("invalid intermediate code format")
	}

	// Parse the code to inspect inside
	parsed, err := bip38.ParseIntermediateCode(intermediateCode)
	if err != nil {
		return fmt.Errorf("failed to parse intermediate code: %v", err)
	}

	// Output validation result clearly
	fmt.Println("✓ Valid intermediate code")

	// Show details about code metadata
	fmt.Printf("Type: ")
	if parsed.HasLotSeq {
		fmt.Printf("With lot/sequence\n")
		fmt.Printf("Lot number: %d\n", *parsed.LotNumber)
		fmt.Printf("Sequence number: %d\n", *parsed.SeqNumber)
	} else {
		fmt.Printf("No lot/sequence\n")
	}

	// Show extra info in verbose mode
	if isVerbose(cmd) {
		fmt.Printf("Owner salt: %x\n", parsed.OwnerSalt)
		fmt.Printf("Pass point: %x\n", parsed.PassPoint)
	}

	return nil
}
