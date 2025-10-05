package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/carlosrabelo/bip38cli/core/internal/domain/bip38"
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
	// Validate when flags fight each other
	if forceCompressed && forceUncompressed {
		return fmt.Errorf("cannot specify both --compressed and --uncompressed")
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
		return fmt.Errorf("invalid WIF private key: %v", err)
	}

	// Apply compression flags when user ask
	if forceCompressed {
		wif.CompressPubKey = true
	} else if forceUncompressed {
		wif.CompressPubKey = false
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

	// Confirm passphrase second time to avoid typo
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
		return fmt.Errorf("encryption failed: %v", err)
	}

	// Output result for user view
	fmt.Printf("Encrypted key: %s\n", encryptedKey)

	// Show extra info while verbose stay on
	if isVerbose(cmd) {
		compression := "uncompressed"
		if wif.CompressPubKey {
			compression = "compressed"
		}
		fmt.Printf("Key format: %s\n", compression)
	}

	return nil
}

func getPassphrase(prompt string) ([]byte, error) {
	fmt.Print(prompt)
	bytePassword, err := readPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	fmt.Println() // Print newline after hidden input so shell not messy
	return bytePassword, nil
}

// secureZero wipe buffer small way for security
func secureZero(buf []byte) {
	if buf == nil {
		return
	}
	for i := range buf {
		buf[i] = 0
	}
}
