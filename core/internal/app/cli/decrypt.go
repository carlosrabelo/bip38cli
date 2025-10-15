package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/carlosrabelo/bip38cli/core/internal/domain/bip38"
	"github.com/carlosrabelo/bip38cli/core/internal/pkg/errors"
	"github.com/carlosrabelo/bip38cli/core/internal/pkg/logger"
	"github.com/spf13/cobra"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt [ENCRYPTED_KEY]",
	Short: "Decrypt a BIP38 encrypted Bitcoin private key",
	Long: `Decrypt a BIP38 encrypted Bitcoin private key using a passphrase.

The encrypted key should be in the standard BIP38 format (starting with 6P).
If no encrypted key is provided as an argument, you will be prompted to enter it.
The passphrase will always be prompted securely.

Examples:
  bip38cli decrypt 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg
  bip38cli decrypt --show-address 6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDecrypt,
}

var showAddress bool

func init() {
	rootCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().BoolVar(&showAddress, "show-address", false, "show the Bitcoin address for the decrypted key")
}

func runDecrypt(cmd *cobra.Command, args []string) error {
	// Reinitialize logger with verbose setting if needed
	if isVerbose(cmd) {
		logger.Init(true)
	}

	logger.Debug("Starting decryption process")

	// Grab encrypted key text
	var encryptedKey string
	if len(args) > 0 {
		encryptedKey = args[0]
	} else {
		fmt.Print("Enter BIP38 encrypted key: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			encryptedKey = strings.TrimSpace(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			logger.WithError(err).Error("Failed to read encrypted key from stdin")
			return errors.NewInputError("failed to read encrypted key", err)
		}
	}

	if encryptedKey == "" {
		logger.Error("No encrypted key provided")
		return errors.NewValidationError("encrypted key is required", nil)
	}

	// Validate format before spending time on decryption
	if !bip38.IsBIP38Format(encryptedKey) {
		logger.WithField("key", encryptedKey).Error("Invalid BIP38 encrypted key format")
		return errors.NewValidationError("invalid BIP38 encrypted key format", nil).
			WithContext("key", encryptedKey)
	}

	logger.Debug("Validated BIP38 encrypted key format")

	// Ask for passphrase the same way as encrypt
	passphrase, err := getPassphrase("Enter passphrase: ")
	if err != nil {
		return fmt.Errorf("failed to read passphrase: %v", err)
	}
	defer secureZero(passphrase)

	if len(passphrase) == 0 {
		return fmt.Errorf("passphrase cannot be empty")
	}

	// Decrypt the key with domain helper
	wif, err := bip38.DecryptKey(encryptedKey, passphrase)
	if err != nil {
		logger.WithError(err).Error("Failed to decrypt private key")
		return errors.NewCryptoError("decryption failed", err)
	}

	logger.Info("Successfully decrypted private key")

	// Prepare output data
	result := map[string]interface{}{
		"private_key": wif.String(),
		"compressed":  wif.CompressPubKey,
	}

	// Show address when user requests
	if showAddress {
		pubKey := wif.PrivKey.PubKey()
		var pubKeyBytes []byte
		if wif.CompressPubKey {
			pubKeyBytes = pubKey.SerializeCompressed()
		} else {
			pubKeyBytes = pubKey.SerializeUncompressed()
		}

		addressPubKey, err := btcutil.NewAddressPubKey(pubKeyBytes, &chaincfg.MainNetParams)
		if err != nil {
			return fmt.Errorf("failed to create address: %v", err)
		}

		result["address"] = addressPubKey.EncodeAddress()
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
		fmt.Printf("Private key (WIF): %s\n", wif.String())

		if showAddress {
			fmt.Printf("Bitcoin address: %s\n", result["address"])
		}

		// Show extra info when verbose flag is active
		if isVerbose(cmd) {
			compression := "uncompressed"
			if wif.CompressPubKey {
				compression = "compressed"
			}
			fmt.Printf("Key format: %s\n", compression)
		}
	}

	return nil
}
