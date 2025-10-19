package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/carlosrabelo/bip38cli/core/internal/domain/bip38"
	"github.com/carlosrabelo/bip38cli/core/internal/pkg/errors"
	"github.com/carlosrabelo/bip38cli/core/internal/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	walletNetwork             = "mainnet"
	walletEncrypt             bool
	walletShowAddr            bool
	walletForceCompressed     bool
	walletForceUncompressed   bool
	walletAddressType         = "bip84"
	walletInspectAddressType  = "bip84"
	generateWIF               = bip38.GenerateWIF
)

var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Utilities for generating and inspecting Bitcoin wallets",
	Long: `Manage wallet-related operations.

The generate subcommand creates a fresh private key (WIF) for the selected
network. Optionally encrypt the result with BIP38 and display the derived
Bitcoin address.`,
}

var walletGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new Bitcoin private key (WIF)",
	Long: `Generate a new Bitcoin private key encoded as WIF.

By default the command targets mainnet and produces a compressed key. The
global --compressed flag may be used to opt into uncompressed format. Pass
--encrypt to wrap the generated WIF in BIP38 using an interactive passphrase
prompt.`,
	RunE: runWalletGenerate,
}

var walletInspectCmd = &cobra.Command{
	Use:   "inspect [WIF]",
	Short: "Display metadata for an existing WIF",
	Long: `Inspect an existing Wallet Import Format (WIF) key.

The command validates the provided WIF, infers its Bitcoin network, and
outputs the derived address. The WIF can be provided as an argument or typed
interactively when no argument is passed.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runWalletInspect,
}

func init() {
	rootCmd.AddCommand(walletCmd)
	walletCmd.AddCommand(walletGenerateCmd)
	walletCmd.AddCommand(walletInspectCmd)

	walletGenerateCmd.Flags().StringVar(&walletNetwork, "network", "mainnet", "target network (mainnet|testnet|regtest|simnet|signet)")
	walletGenerateCmd.Flags().BoolVar(&walletEncrypt, "encrypt", false, "encrypt generated key using BIP38")
	walletGenerateCmd.Flags().BoolVar(&walletShowAddr, "show-address", false, "show the Bitcoin address for the generated key")
	walletGenerateCmd.Flags().BoolVar(&walletForceCompressed, "compressed", false, "force compressed public key format")
	walletGenerateCmd.Flags().BoolVar(&walletForceUncompressed, "uncompressed", false, "force uncompressed public key format")
	walletGenerateCmd.Flags().StringVar(&walletAddressType, "address-type", "bip84", "address type (bip84|bip44)")

	walletInspectCmd.Flags().StringVar(&walletInspectAddressType, "address-type", "bip84", "address type (bip84|bip44)")
}

func runWalletGenerate(cmd *cobra.Command, args []string) error {
	if isVerbose(cmd) {
		logger.Init(true)
	}
	defer func() {
		walletForceCompressed = false
		walletForceUncompressed = false
	}()

	if walletForceCompressed && walletForceUncompressed {
		return errors.NewValidationError("cannot specify both --compressed and --uncompressed", nil)
	}

	compressed := isCompressed(cmd)
	if walletForceCompressed {
		compressed = true
	} else if walletForceUncompressed {
		compressed = false
	}

	params, err := bip38.NetworkFromName(walletNetwork)
	if err != nil {
		logger.WithError(err).Error("Unsupported network provided")
		return errors.NewValidationError("invalid network", err).
			WithContext("network", walletNetwork)
	}

	wif, err := generateWIF(params, compressed)
	if err != nil {
		logger.WithError(err).Error("Failed to generate WIF")
		return errors.NewCryptoError("failed to generate private key", err)
	}

	addrType, err := parseAddressType(walletAddressType)
	if err != nil {
		return errors.NewValidationError("invalid address type", err).
			WithContext("address_type", walletAddressType)
	}

	effectiveType := addrType
	if effectiveType == addressTypeBIP84 && !wif.CompressPubKey {
		effectiveType = addressTypeBIP44
	}

	result := map[string]any{
		"wif":          wif.String(),
		"compressed":   wif.CompressPubKey,
		"network":      params.Name,
		"address_type": string(effectiveType),
	}

	if walletShowAddr {
		address, err := addressForWIF(wif, effectiveType)
		if err != nil {
			return fmt.Errorf("failed to derive address: %v", err)
		}
		result["address"] = address
	}

	if walletEncrypt {
		passphrase, err := getPassphrase("Enter passphrase for encryption: ")
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

		encryptedKey, err := bip38.EncryptKey(wif, passphrase)
		if err != nil {
			logger.WithError(err).Error("Failed to encrypt generated WIF")
			return errors.NewCryptoError("failed to encrypt generated key", err)
		}

		result["bip38_encrypted_key"] = encryptedKey
	}

	switch outputFormat(cmd) {
	case "json":
		jsonOutput, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON output: %v", err)
		}
		fmt.Println(string(jsonOutput))
	default:
		compression := "uncompressed"
		if wif.CompressPubKey {
			compression = "compressed"
		}

		fmt.Printf("WIF (%s, %s): %s\n", params.Name, compression, result["wif"])

		if walletShowAddr {
			fmt.Printf("Address (%s): %s\n", result["address_type"], result["address"])
		}

		if walletEncrypt {
			fmt.Printf("BIP38 encrypted key: %s\n", result["bip38_encrypted_key"])
		}
	}

	return nil
}

func runWalletInspect(cmd *cobra.Command, args []string) error {
	if isVerbose(cmd) {
		logger.Init(true)
	}

	var wifStr string
	if len(args) > 0 {
		wifStr = args[0]
	} else {
		fmt.Print("Enter WIF private key: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			wifStr = strings.TrimSpace(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			logger.WithError(err).Error("Failed to read WIF from stdin")
			return errors.NewInputError("failed to read WIF", err)
		}
	}

	if wifStr == "" {
		return errors.NewValidationError("WIF private key is required", nil)
	}

	wif, err := btcutil.DecodeWIF(wifStr)
	if err != nil {
		logger.WithError(err).Error("Failed to decode WIF private key")
		return errors.NewValidationError("invalid WIF private key", err).
			WithContext("wif", wifStr)
	}

	params, err := bip38.NetworkFromWIF(wif)
	if err != nil {
		logger.WithError(err).Error("Unsupported WIF network")
		return errors.NewValidationError("unsupported WIF network", err)
	}

	addrType, err := parseAddressType(walletInspectAddressType)
	if err != nil {
		return errors.NewValidationError("invalid address type", err).
			WithContext("address_type", walletInspectAddressType)
	}

	effectiveType := addrType
	if effectiveType == addressTypeBIP84 && !wif.CompressPubKey {
		effectiveType = addressTypeBIP44
	}

	address, err := addressForWIF(wif, effectiveType)
	if err != nil {
		return fmt.Errorf("failed to derive address: %v", err)
	}

	result := map[string]any{
		"wif":          wif.String(),
		"compressed":   wif.CompressPubKey,
		"network":      params.Name,
		"address":      address,
		"address_type": string(effectiveType),
	}

	switch outputFormat(cmd) {
	case "json":
		jsonOutput, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON output: %v", err)
		}
		fmt.Println(string(jsonOutput))
	default:
		compression := "uncompressed"
		if wif.CompressPubKey {
			compression = "compressed"
		}

		fmt.Printf("WIF (%s, %s): %s\n", params.Name, compression, result["wif"])
		fmt.Printf("Address (%s): %s\n", result["address_type"], address)
	}

	return nil
}
