package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/carlosrabelo/bip38cli/bip38cli/internal/bip38"
	"github.com/carlosrabelo/bip38cli/bip38cli/internal/errors"
	"github.com/carlosrabelo/bip38cli/bip38cli/internal/logger"
	"github.com/carlosrabelo/bip38cli/bip38cli/internal/metrics"
	"github.com/spf13/cobra"
)

var intermediateCmd = &cobra.Command{
	Use:   "intermediate",
	Short: "Generate BIP38 intermediate passphrase codes",
	Long: `Generate BIP38 intermediate passphrase codes for two-factor encryption.

Intermediate codes allow a third party to generate encrypted private keys
without knowing the passphrase. This enables secure key generation in
scenarios where the key generator should not have access to the passphrase.`,
}

var generateIntermediateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate an intermediate passphrase code",
	Long: `Generate a BIP38 intermediate passphrase code from a passphrase.

Examples:
  bip38cli intermediate generate
  bip38cli intermediate generate --lot 123 --sequence 456`,
	RunE: runGenerateIntermediate,
}

var validateIntermediateCmd = &cobra.Command{
	Use:   "validate [INTERMEDIATE_CODE]",
	Short: "Validate an intermediate passphrase code",
	Long: `Validate the format and integrity of a BIP38 intermediate passphrase code.

Examples:
  bip38cli intermediate validate passphraseabc123...`,
	Args: cobra.MaximumNArgs(1),
	RunE: runValidateIntermediate,
}

var encryptIntermediateCmd = &cobra.Command{
	Use:   "encrypt [INTERMEDIATE_CODE]",
	Short: "Generate a BIP38 EC-multiply encrypted key from an intermediate code",
	Long: `Generate a BIP38 encrypted private key using the EC-multiply scheme.

The intermediate code is produced by the passphrase owner (via 'intermediate generate')
and handed to a third party. This command generates an encrypted key without
requiring knowledge of the passphrase.

The output includes the encrypted key (6P...) and a confirmation code (cfrm38...)
that the passphrase owner can use to verify the derived address.

Examples:
  bip38cli intermediate encrypt passphraseXXX...
  bip38cli intermediate encrypt --uncompressed passphraseXXX...
  bip38cli intermediate encrypt --output-format json passphraseXXX...`,
	Args: cobra.MaximumNArgs(1),
	RunE: runEncryptIntermediate,
}

var (
	lotNumber                    uint32
	sequenceNumber               uint32
	useLotSeq                    bool
	encryptIntermediateUncompressed bool
)

func init() {
	rootCmd.AddCommand(intermediateCmd)
	intermediateCmd.AddCommand(generateIntermediateCmd)
	intermediateCmd.AddCommand(validateIntermediateCmd)
	intermediateCmd.AddCommand(encryptIntermediateCmd)

	generateIntermediateCmd.Flags().Uint32Var(&lotNumber, "lot", 0, "lot number (0-1048575)")
	generateIntermediateCmd.Flags().Uint32Var(&sequenceNumber, "sequence", 0, "sequence number (0-4095)")
	generateIntermediateCmd.Flags().BoolVar(&useLotSeq, "use-lot-sequence", false, "use lot and sequence numbers")

	encryptIntermediateCmd.Flags().BoolVar(&encryptIntermediateUncompressed, "uncompressed", false, "generate uncompressed key")
}

func runGenerateIntermediate(cmd *cobra.Command, _ []string) error { //nolint:gocyclo
	if isVerbose(cmd) {
		logger.Init(true)
	}

	logger.Debug("Starting intermediate code generation")

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

	passphrase, err := getPassphrase("Enter passphrase: ")
	if err != nil {
		return fmt.Errorf("failed to read passphrase: %v", err)
	}
	defer secureZero(passphrase)

	if len(passphrase) == 0 {
		return fmt.Errorf("passphrase cannot be empty")
	}

	confirmPassphrase, err := getPassphrase("Confirm passphrase: ")
	if err != nil {
		return fmt.Errorf("failed to read passphrase confirmation: %v", err)
	}
	defer secureZero(confirmPassphrase)

	if !bytes.Equal(passphrase, confirmPassphrase) {
		return fmt.Errorf("passphrases do not match")
	}

	timer := metrics.NewTimer("intermediate")
	intermediate, err := bip38.GenerateIntermediateCode(passphrase, lot, seq)
	if err != nil {
		timer.Stop(false)
		logger.WithError(err).Error("Failed to generate intermediate code")
		return errors.NewCryptoError("failed to generate intermediate code", err)
	}
	timer.Stop(true)

	logger.Info("Successfully generated intermediate code")

	result := map[string]interface{}{
		"intermediate_code": intermediate,
		"has_lot_sequence":  lot != nil && seq != nil,
	}

	if lot != nil && seq != nil {
		result["lot_number"] = *lot
		result["sequence_number"] = *seq
	}

	switch outputFormat(cmd) {
	case "json":
		jsonOutput, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON output: %v", err)
		}
		fmt.Println(string(jsonOutput))
	default:
		fmt.Printf("Intermediate code: %s\n", intermediate)
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

	if !bip38.IsValidIntermediateCode(intermediateCode) {
		return fmt.Errorf("invalid intermediate code format")
	}

	parsed, err := bip38.ParseIntermediateCode(intermediateCode)
	if err != nil {
		return fmt.Errorf("failed to parse intermediate code: %v", err)
	}

	fmt.Println("✓ Valid intermediate code")
	fmt.Printf("Type: ")
	if parsed.HasLotSeq {
		fmt.Printf("With lot/sequence\n")
		fmt.Printf("Lot number: %d\n", *parsed.LotNumber)
		fmt.Printf("Sequence number: %d\n", *parsed.SeqNumber)
	} else {
		fmt.Printf("No lot/sequence\n")
	}

	if isVerbose(cmd) {
		fmt.Printf("Owner salt: %x\n", parsed.OwnerSalt)
		fmt.Printf("Pass point: %x\n", parsed.PassPoint)
	}

	return nil
}

func runEncryptIntermediate(cmd *cobra.Command, args []string) error {
	if isVerbose(cmd) {
		logger.Init(true)
	}

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

	if !bip38.IsValidIntermediateCode(intermediateCode) {
		return errors.NewValidationError("invalid intermediate code format", nil)
	}

	compressed := !encryptIntermediateUncompressed

	timer := metrics.NewTimer("intermediate")
	ecResult, err := bip38.ECMultiplyEncrypt(intermediateCode, compressed)
	if err != nil {
		timer.Stop(false)
		logger.WithError(err).Error("EC-multiply encryption failed")
		return errors.NewCryptoError("EC-multiply encryption failed", err)
	}
	timer.Stop(true)

	logger.Info("Successfully generated EC-multiply encrypted key")

	out := map[string]interface{}{
		"encrypted_key":     ecResult.EncryptedKey,
		"confirmation_code": ecResult.ConfirmationCode,
		"compressed":        ecResult.Compressed,
	}

	switch outputFormat(cmd) {
	case "json":
		jsonOutput, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON output: %v", err)
		}
		fmt.Println(string(jsonOutput))
	default:
		fmt.Printf("Encrypted key:     %s\n", ecResult.EncryptedKey)
		fmt.Printf("Confirmation code: %s\n", ecResult.ConfirmationCode)
		if isVerbose(cmd) {
			compression := "uncompressed"
			if ecResult.Compressed {
				compression = "compressed"
			}
			fmt.Printf("Key format: %s\n", compression)
		}
	}

	return nil
}
