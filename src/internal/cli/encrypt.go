package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/carlosrabelo/bip38cli/internal/bip38"
	"github.com/carlosrabelo/bip38cli/pkg/errors"
	"github.com/carlosrabelo/bip38cli/pkg/logger"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt [WIF_PRIVATE_KEY]",
	Short: "Encrypt a Bitcoin private key using BIP38",
	Long: `Encrypt a Bitcoin private key using BIP38 standard encryption.

The private key should be provided in WIF (Wallet Import Format).
If no private key is provided as an argument, you will be prompted to enter it.
The passphrase will always be prompted securely.

Examples:
  bip38cli encrypt 5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ
  bip38cli encrypt --compressed 5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ`,
	Args: cobra.MaximumNArgs(1),
	RunE: runEncrypt,
}

var (
	forceCompressed   bool
	forceUncompressed bool
)

var readPassword = term.ReadPassword

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().BoolVar(&forceCompressed, "compressed", false, "force compressed public key format")
	encryptCmd.Flags().BoolVar(&forceUncompressed, "uncompressed", false, "force uncompressed public key format")
}

func runEncrypt(cmd *cobra.Command, args []string) error {
	// Reinitialize logger with verbose setting if needed
	if isVerbose(cmd) {
		logger.Init(true)
	}

	logger.Debug("Starting encryption process")

	// Validate when flags conflict with each other
	if forceCompressed && forceUncompressed {
		logger.Error("Both compressed and uncompressed flags specified")
		return errors.NewValidationError("cannot specify both --compressed and --uncompressed", nil)
	}

	// Grab private key value
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
			return fmt.Errorf("failed to read private key: %v", err)
		}
	}

	if wifStr == "" {
		return fmt.Errorf("private key is required")
	}

	// Parse WIF so we know key format
	wif, err := btcutil.DecodeWIF(wifStr)
	if err != nil {
		logger.WithError(err).Error("Failed to decode WIF private key")
		return errors.NewValidationError("invalid WIF private key", err).
			WithContext("wif", wifStr)
	}
	logger.WithField("compressed", wif.CompressPubKey).Debug("Successfully decoded WIF private key")

	// Apply compression flags when user requests
	if forceCompressed {
		wif.CompressPubKey = true
	} else if forceUncompressed {
		wif.CompressPubKey = false
	} else {
		// Use global compressed flag when no specific flag is set
		wif.CompressPubKey = isCompressed(cmd)
	}

	// Ask hidden passphrase from user
	passphrase, err := getPassphrase("Enter passphrase for encryption: ")
	if err != nil {
		return fmt.Errorf("failed to read passphrase: %v", err)
	}
	defer secureZero(passphrase)

	if len(passphrase) == 0 {
		return fmt.Errorf("passphrase cannot be empty")
	}

	// Confirm passphrase a second time to avoid typos
	confirmPassphrase, err := getPassphrase("Confirm passphrase: ")
	if err != nil {
		return fmt.Errorf("failed to read passphrase confirmation: %v", err)
	}
	defer secureZero(confirmPassphrase)

	if !bytes.Equal(passphrase, confirmPassphrase) {
		return fmt.Errorf("passphrases do not match")
	}

	// Encrypt the key using domain logic
	encryptedKey, err := bip38.EncryptKey(wif, passphrase)
	if err != nil {
		logger.WithError(err).Error("Failed to encrypt private key")
		return errors.NewCryptoError("encryption failed", err)
	}

	logger.Info("Successfully encrypted private key")

	// Prepare output data
	result := map[string]interface{}{
		"encrypted_key": encryptedKey,
		"compressed":    wif.CompressPubKey,
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
		fmt.Printf("Encrypted key: %s\n", encryptedKey)

		// Show extra info while verbose mode is on
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

func getPassphrase(prompt string) ([]byte, error) {
	fmt.Print(prompt)
	bytePassword, err := readPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	fmt.Println() // Print newline after hidden input so shell isn't messy
	return bytePassword, nil
}

// secureZero wipe buffer small way for security
func secureZero(buf []byte) {
	if buf == nil {
		return
	}
	for i := 0; i < len(buf); i++ {
		buf[i] = 0
	}
}
