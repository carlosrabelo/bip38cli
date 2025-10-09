package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/carlosrabelo/bip38cli/internal/bip38"
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

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().BoolVar(&forceCompressed, "compressed", false, "force compressed public key format")
	encryptCmd.Flags().BoolVar(&forceUncompressed, "uncompressed", false, "force uncompressed public key format")
}

func runEncrypt(cmd *cobra.Command, args []string) error {
	// Validate conflicting flags
	if forceCompressed && forceUncompressed {
		return fmt.Errorf("cannot specify both --compressed and --uncompressed")
	}

	// Get private key
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

	// Parse WIF
	wif, err := btcutil.DecodeWIF(wifStr)
	if err != nil {
		return fmt.Errorf("invalid WIF private key: %v", err)
	}

	// Apply compression flags if specified
	if forceCompressed {
		wif.CompressPubKey = true
	} else if forceUncompressed {
		wif.CompressPubKey = false
	}

	// Get passphrase
	passphrase, err := getPassphrase("Enter passphrase for encryption: ")
	if err != nil {
		return fmt.Errorf("failed to read passphrase: %v", err)
	}

	if len(passphrase) == 0 {
		return fmt.Errorf("passphrase cannot be empty")
	}

	// Confirm passphrase
	confirmPassphrase, err := getPassphrase("Confirm passphrase: ")
	if err != nil {
		return fmt.Errorf("failed to read passphrase confirmation: %v", err)
	}

	if passphrase != confirmPassphrase {
		return fmt.Errorf("passphrases do not match")
	}

	// Encrypt the key
	encryptedKey, err := bip38.EncryptKey(wif, passphrase)
	if err != nil {
		return fmt.Errorf("encryption failed: %v", err)
	}

	// Output result
	fmt.Printf("Encrypted key: %s\n", encryptedKey)

	// Show additional info in verbose mode
	if cmd.Flag("verbose").Changed {
		compression := "uncompressed"
		if wif.CompressPubKey {
			compression = "compressed"
		}
		fmt.Printf("Key format: %s\n", compression)
	}

	return nil
}

func getPassphrase(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println() // Add newline after password input
	return string(bytePassword), nil
}
