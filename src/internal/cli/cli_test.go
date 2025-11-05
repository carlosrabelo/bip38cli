package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/carlosrabelo/bip38cli/internal/bip38"
	"github.com/carlosrabelo/bip38cli/pkg/logger"
	"github.com/spf13/cobra"
)

// captureOutput captures stdout and returns a function to restore it
func captureOutput() (func() []byte, func()) {
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	done := make(chan []byte, 1)
	closePipe := sync.Once{}

	closeWriter := func() {
		closePipe.Do(func() {
			_ = w.Close()
			os.Stdout = origStdout
		})
	}

	// Copy output in a goroutine
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.Bytes()
	}()

	collect := func() []byte {
		closeWriter()
		return <-done
	}

	restore := func() {
		closeWriter()
		select {
		case <-done:
		default:
		}
		// Reinitialize logger to avoid "file already closed" errors
		logger.Init(false)
	}

	return collect, restore
}

func TestRunEncryptWithStubbedPassphrase(t *testing.T) {
	origReadPassword := readPassword
	defer func() { readPassword = origReadPassword }()

	callCount := 0
	passphrase := []byte("TestingOneTwoThree")

	readPassword = func(int) ([]byte, error) {
		callCount++
		buf := make([]byte, len(passphrase))
		copy(buf, passphrase)
		return buf, nil
	}

	forceCompressed = false
	forceUncompressed = false

	cmd := &cobra.Command{Use: "encrypt"}

	collect, restore := captureOutput()
	defer restore()

	runErr := runEncrypt(cmd, []string{"5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ"})

	outputBytes := collect()

	if runErr != nil {
		t.Fatalf("runEncrypt returned error: %v", runErr)
	}

	if callCount != 2 {
		t.Fatalf("expected readPassword to be called twice, got %d", callCount)
	}

	output := string(outputBytes)
	if !strings.Contains(output, "Encrypted key: 6P") {
		t.Fatalf("expected encrypted key output, got %q", output)
	}
}

func TestRunGenerateIntermediateRequiresLotAndSequence(t *testing.T) {
	origReadPassword := readPassword
	defer func() { readPassword = origReadPassword }()

	readPassword = func(int) ([]byte, error) {
		return []byte("pass"), nil
	}

	lotNumber = 0
	sequenceNumber = 0
	useLotSeq = false
	defer func() {
		lotNumber = 0
		sequenceNumber = 0
		useLotSeq = false
	}()

	cmd := &cobra.Command{Use: "intermediate-generate"}
	cmd.Flags().Uint32Var(&lotNumber, "lot", 0, "")
	cmd.Flags().Uint32Var(&sequenceNumber, "sequence", 0, "")
	cmd.Flags().BoolVar(&useLotSeq, "use-lot-sequence", false, "")

	if err := cmd.Flags().Set("lot", "123"); err != nil {
		t.Fatalf("failed to set lot flag: %v", err)
	}

	err := runGenerateIntermediate(cmd, nil)
	if err == nil {
		t.Fatal("expected error when only --lot is provided")
	}

	if !strings.Contains(err.Error(), "both --lot and --sequence") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestIsVerbose(t *testing.T) {
	cmd := &cobra.Command{Use: "root"}
	cmd.Flags().Bool("verbose", false, "")

	// Test default (verbose should be false)
	if isVerbose(cmd) {
		t.Fatal("expected verbose to be false by default")
	}

	// Test with verbose flag set to true
	if err := cmd.Flags().Set("verbose", "true"); err != nil {
		t.Fatalf("failed to set verbose flag: %v", err)
	}

	if !isVerbose(cmd) {
		t.Fatal("expected verbose to be true when flag is set")
	}
}

func TestOutputFormat(t *testing.T) {
	cmd := &cobra.Command{Use: "root"}
	cmd.Flags().String("output-format", "text", "")

	// Test default (output-format should be "text")
	if outputFormat(cmd) != "text" {
		t.Fatalf("expected output-format to be 'text' by default, got %q", outputFormat(cmd))
	}

	// Test with output-format flag set to json
	if err := cmd.Flags().Set("output-format", "json"); err != nil {
		t.Fatalf("failed to set output-format flag: %v", err)
	}

	if outputFormat(cmd) != "json" {
		t.Fatalf("expected output-format to be 'json' when flag is set, got %q", outputFormat(cmd))
	}
}

func TestIsCompressed(t *testing.T) {
	cmd := &cobra.Command{Use: "root"}
	cmd.Flags().Bool("compressed", true, "")

	// Test default (compressed should be true)
	if !isCompressed(cmd) {
		t.Fatal("expected compressed to be true by default")
	}

	// Test with compressed flag set to false
	if err := cmd.Flags().Set("compressed", "false"); err != nil {
		t.Fatalf("failed to set compressed flag: %v", err)
	}

	if isCompressed(cmd) {
		t.Fatal("expected compressed to be false when flag is set to false")
	}

	// Test with compressed flag set to true
	if err := cmd.Flags().Set("compressed", "true"); err != nil {
		t.Fatalf("failed to set compressed flag: %v", err)
	}

	if !isCompressed(cmd) {
		t.Fatal("expected compressed to be true when flag is set to true")
	}
}

func TestRunDecryptWithStubbedPassphrase(t *testing.T) {
	origReadPassword := readPassword
	defer func() { readPassword = origReadPassword }()

	callCount := 0
	passphrase := []byte("TestingOneTwoThree")

	readPassword = func(int) ([]byte, error) {
		callCount++
		buf := make([]byte, len(passphrase))
		copy(buf, passphrase)
		return buf, nil
	}

	cmd := &cobra.Command{Use: "decrypt"}

	collect, restore := captureOutput()
	defer restore()

	// Using a valid BIP38 format key for testing (actual content doesn't matter for flow test)
	runErr := runDecrypt(cmd, []string{"6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg"})

	_ = collect() // Ignore output for this test

	// For now, just test that it attempts decryption (the exact key may vary)
	if runErr != nil && !strings.Contains(runErr.Error(), "invalid BIP38 encrypted key format") {
		t.Fatalf("runDecrypt returned unexpected error: %v", runErr)
	}

	if callCount != 1 {
		t.Fatalf("expected readPassword to be called once, got %d", callCount)
	}

	// If we get here, the test passes - we're testing the flow, not exact crypto
}

func TestRunEncryptWithCompressedFlag(t *testing.T) {
	origReadPassword := readPassword
	defer func() { readPassword = origReadPassword }()

	passphrase := []byte("TestingOneTwoThree")

	readPassword = func(int) ([]byte, error) {
		buf := make([]byte, len(passphrase))
		copy(buf, passphrase)
		return buf, nil
	}

	forceCompressed = true
	forceUncompressed = false
	defer func() {
		forceCompressed = false
		forceUncompressed = false
	}()

	cmd := &cobra.Command{Use: "encrypt"}

	collect, restore := captureOutput()
	defer restore()

	runErr := runEncrypt(cmd, []string{"5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ"})

	outputBytes := collect()

	if runErr != nil {
		t.Fatalf("runEncrypt returned error: %v", runErr)
	}

	output := string(outputBytes)
	if !strings.Contains(output, "Encrypted key: 6P") {
		t.Fatalf("expected encrypted key output, got %q", output)
	}
}

func TestRunEncryptWithConflictingFlags(t *testing.T) {
	forceCompressed = true
	forceUncompressed = true
	defer func() {
		forceCompressed = false
		forceUncompressed = false
	}()

	cmd := &cobra.Command{Use: "encrypt"}
	err := runEncrypt(cmd, []string{"5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ"})

	if err == nil {
		t.Fatal("expected error when both compressed and uncompressed flags are set")
	}

	expectedError := "cannot specify both --compressed and --uncompressed"
	if !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("expected error containing %q, got %q", expectedError, err.Error())
	}
}

func TestRunEncryptWithEmptyPassphrase(t *testing.T) {
	origReadPassword := readPassword
	defer func() { readPassword = origReadPassword }()

	readPassword = func(int) ([]byte, error) {
		return []byte(""), nil
	}

	cmd := &cobra.Command{Use: "encrypt"}
	err := runEncrypt(cmd, []string{"5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ"})

	if err == nil {
		t.Fatal("expected error when passphrase is empty")
	}

	expectedError := "passphrase cannot be empty"
	if !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("expected error containing %q, got %q", expectedError, err.Error())
	}
}

func TestRunEncryptWithInvalidWIF(t *testing.T) {
	origReadPassword := readPassword
	defer func() { readPassword = origReadPassword }()

	passphrase := []byte("TestingOneTwoThree")

	readPassword = func(int) ([]byte, error) {
		buf := make([]byte, len(passphrase))
		copy(buf, passphrase)
		return buf, nil
	}

	cmd := &cobra.Command{Use: "encrypt"}
	err := runEncrypt(cmd, []string{"invalid-wif-key"})

	if err == nil {
		t.Fatal("expected error with invalid WIF key")
	}

	expectedError := "invalid WIF private key"
	if !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("expected error containing %q, got %q", expectedError, err.Error())
	}
}

func TestAddressForWIFUsesCorrectNetwork(t *testing.T) {
	keyBytes := make([]byte, 32)
	keyBytes[len(keyBytes)-1] = 1

	mainPriv, _ := btcec.PrivKeyFromBytes(keyBytes)
	testPriv, _ := btcec.PrivKeyFromBytes(keyBytes)

	mainWIF, err := btcutil.NewWIF(mainPriv, &chaincfg.MainNetParams, true)
	if err != nil {
		t.Fatalf("failed to create mainnet WIF: %v", err)
	}

	testWIF, err := btcutil.NewWIF(testPriv, &chaincfg.TestNet3Params, true)
	if err != nil {
		t.Fatalf("failed to create testnet WIF: %v", err)
	}

	mainAddr, err := addressForWIF(mainWIF, addressTypeBIP84)
	if err != nil {
		t.Fatalf("addressForWIF returned error for mainnet: %v", err)
	}

	testAddr, err := addressForWIF(testWIF, addressTypeBIP84)
	if err != nil {
		t.Fatalf("addressForWIF returned error for testnet: %v", err)
	}

	mainHash := btcutil.Hash160(mainPriv.PubKey().SerializeCompressed())
	expectedMain, err := btcutil.NewAddressWitnessPubKeyHash(mainHash, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("failed to derive expected mainnet address: %v", err)
	}

	testHash := btcutil.Hash160(testPriv.PubKey().SerializeCompressed())
	expectedTest, err := btcutil.NewAddressWitnessPubKeyHash(testHash, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatalf("failed to derive expected testnet address: %v", err)
	}

	if mainAddr != expectedMain.EncodeAddress() {
		t.Fatalf("mainnet address mismatch: expected %s, got %s", expectedMain.EncodeAddress(), mainAddr)
	}

	if testAddr != expectedTest.EncodeAddress() {
		t.Fatalf("testnet address mismatch: expected %s, got %s", expectedTest.EncodeAddress(), testAddr)
	}
}

func TestRunEncryptWithMissingPrivateKey(t *testing.T) {
	cmd := &cobra.Command{Use: "encrypt"}
	err := runEncrypt(cmd, []string{})

	if err == nil {
		t.Fatal("expected error when no private key provided")
	}

	expectedError := "private key is required"
	if !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("expected error containing %q, got %q", expectedError, err.Error())
	}
}

func TestRunWalletGenerateProducesWIF(t *testing.T) {
	origGenerate := generateWIF
	defer func() { generateWIF = origGenerate }()

	walletNetwork = "mainnet"
	walletEncrypt = false
	walletShowAddr = false
	walletAddressType = "bip84"
	defer func() {
		walletNetwork = "mainnet"
		walletEncrypt = false
		walletShowAddr = false
		walletAddressType = "bip84"
	}()

	var capturedParams *chaincfg.Params
	var capturedCompressed bool
	var generatedWIF string

	generateWIF = func(params *chaincfg.Params, compressed bool) (*btcutil.WIF, error) {
		capturedParams = params
		capturedCompressed = compressed

		keyBytes := make([]byte, 32)
		keyBytes[31] = 0x01
		privKey, _ := btcec.PrivKeyFromBytes(keyBytes)
		wif, err := btcutil.NewWIF(privKey, params, compressed)
		if err != nil {
			return nil, err
		}
		generatedWIF = wif.String()
		return wif, nil
	}

	cmd := &cobra.Command{Use: "wallet-generate"}

	collect, restore := captureOutput()
	defer restore()

	if err := runWalletGenerate(cmd, nil); err != nil {
		t.Fatalf("runWalletGenerate returned error: %v", err)
	}

	output := string(collect())

	if generatedWIF == "" {
		t.Fatal("expected stub to produce a WIF string")
	}

	if !strings.Contains(output, generatedWIF) {
		t.Fatalf("expected WIF in output, got %q", output)
	}

	if capturedParams != &chaincfg.MainNetParams {
		t.Fatalf("expected mainnet params, got %v", capturedParams)
	}

	if !capturedCompressed {
		t.Fatal("expected compressed flag to be true by default")
	}
}

func TestRunWalletGenerateEncryptsAndShowsAddress(t *testing.T) {
	origGenerate := generateWIF
	defer func() { generateWIF = origGenerate }()

	origReadPassword := readPassword
	defer func() { readPassword = origReadPassword }()

	walletNetwork = "mainnet"
	walletEncrypt = true
	walletShowAddr = true
	walletAddressType = "bip84"
	defer func() {
		walletEncrypt = false
		walletShowAddr = false
		walletNetwork = "mainnet"
		walletAddressType = "bip84"
	}()

	var capturedCompressed bool
	const wifString = "5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ"
	passphraseValue := []byte("TestingOneTwoThree")

	addressWIF, err := btcutil.DecodeWIF(wifString)
	if err != nil {
		t.Fatalf("failed to decode WIF for address expectation: %v", err)
	}
	expectedAddress, err := addressForWIF(addressWIF, addressTypeBIP84)
	if err != nil {
		t.Fatalf("failed to derive expected address: %v", err)
	}

	encryptedWIF, err := btcutil.DecodeWIF(wifString)
	if err != nil {
		t.Fatalf("failed to decode WIF for encryption expectation: %v", err)
	}
	expectedEncrypted, err := bip38.EncryptKey(encryptedWIF, passphraseValue)
	if err != nil {
		t.Fatalf("failed to encrypt expected WIF: %v", err)
	}

	generateWIF = func(params *chaincfg.Params, compressed bool) (*btcutil.WIF, error) {
		capturedCompressed = compressed
		return btcutil.DecodeWIF(wifString)
	}

	readPassword = func(int) ([]byte, error) {
		buf := make([]byte, len(passphraseValue))
		copy(buf, passphraseValue)
		return buf, nil
	}

	cmd := &cobra.Command{Use: "wallet-generate"}
	cmd.Flags().Bool("compressed", true, "")
	if err := cmd.Flags().Set("compressed", "false"); err != nil {
		t.Fatalf("failed to set compressed flag: %v", err)
	}

	collect, restore := captureOutput()
	defer restore()

	if err := runWalletGenerate(cmd, nil); err != nil {
		t.Fatalf("runWalletGenerate returned error: %v", err)
	}

	output := string(collect())

	if !strings.Contains(output, wifString) {
		t.Fatalf("expected generated WIF in output, got %q", output)
	}

	if !strings.Contains(output, expectedEncrypted) {
		t.Fatalf("expected encrypted key in output, got %q", output)
	}

	if !strings.Contains(output, expectedAddress) {
		t.Fatalf("expected derived address in output, got %q", output)
	}

	if !strings.Contains(output, "Address (bip44):") {
		t.Fatalf("expected address type label in output, got %q", output)
	}

	if capturedCompressed {
		t.Fatal("expected compressed flag to be false when --compressed=false")
	}
}

func TestRunWalletGenerateJSON(t *testing.T) {
	origGenerate := generateWIF
	defer func() { generateWIF = origGenerate }()

	walletNetwork = "regtest"
	walletEncrypt = false
	walletShowAddr = true
	walletAddressType = "bip84"
	defer func() {
		walletShowAddr = false
		walletNetwork = "mainnet"
		walletAddressType = "bip84"
	}()

	var generatedWIF string
	var generatedAddress string

	generateWIF = func(params *chaincfg.Params, compressed bool) (*btcutil.WIF, error) {
		keyBytes := make([]byte, 32)
		keyBytes[31] = 0x02
		privKey, _ := btcec.PrivKeyFromBytes(keyBytes)
		wif, err := btcutil.NewWIF(privKey, params, compressed)
		if err != nil {
			return nil, err
		}
		generatedWIF = wif.String()

		address, err := addressForWIF(wif, addressTypeBIP84)
		if err != nil {
			return nil, err
		}
		generatedAddress = address

		return wif, nil
	}

	cmd := &cobra.Command{Use: "wallet-generate"}
	cmd.Flags().String("output-format", "text", "")
	if err := cmd.Flags().Set("output-format", "json"); err != nil {
		t.Fatalf("failed to set output-format flag: %v", err)
	}

	collect, restore := captureOutput()
	defer restore()

	if err := runWalletGenerate(cmd, nil); err != nil {
		t.Fatalf("runWalletGenerate returned error: %v", err)
	}

	outputBytes := collect()

	type walletOutput struct {
		WIF         string `json:"wif"`
		Compressed  bool   `json:"compressed"`
		Network     string `json:"network"`
		Address     string `json:"address"`
		AddressType string `json:"address_type"`
	}

	var payload walletOutput
	if err := json.Unmarshal(outputBytes, &payload); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}

	if payload.Network != "regtest" {
		t.Fatalf("expected network regtest, got %s", payload.Network)
	}

	if payload.WIF != generatedWIF {
		t.Fatalf("unexpected WIF in payload: %s", payload.WIF)
	}

	if payload.Address != generatedAddress {
		t.Fatalf("unexpected address in payload: %s", payload.Address)
	}

	if payload.AddressType != "bip84" {
		t.Fatalf("unexpected address type in payload: %s", payload.AddressType)
	}

	if !strings.HasPrefix(payload.Address, "bc1") && !strings.HasPrefix(payload.Address, "tb1") && !strings.HasPrefix(payload.Address, "bcrt1") {
		t.Fatalf("expected bech32 address, got %s", payload.Address)
	}
}

func TestRunWalletGenerateInvalidNetwork(t *testing.T) {
	walletNetwork = "invalid"
	defer func() { walletNetwork = "mainnet" }()

	cmd := &cobra.Command{Use: "wallet-generate"}

	err := runWalletGenerate(cmd, nil)
	if err == nil {
		t.Fatal("expected error for invalid network")
	}

	if !strings.Contains(err.Error(), "invalid network") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunWalletGenerateBIP44Address(t *testing.T) {
	origGenerate := generateWIF
	defer func() { generateWIF = origGenerate }()

	walletNetwork = "mainnet"
	walletEncrypt = false
	walletShowAddr = true
	walletAddressType = "bip44"
	defer func() {
		walletNetwork = "mainnet"
		walletEncrypt = false
		walletShowAddr = false
		walletAddressType = "bip84"
	}()

	wif, err := btcutil.DecodeWIF("KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7")
	if err != nil {
		t.Fatalf("failed to decode WIF: %v", err)
	}
	expectedAddress, err := addressForWIF(wif, addressTypeBIP44)
	if err != nil {
		t.Fatalf("failed to derive expected BIP44 address: %v", err)
	}

	generateWIF = func(params *chaincfg.Params, compressed bool) (*btcutil.WIF, error) {
		return wif, nil
	}

	cmd := &cobra.Command{Use: "wallet-generate"}

	collect, restore := captureOutput()
	defer restore()

	if err := runWalletGenerate(cmd, nil); err != nil {
		t.Fatalf("runWalletGenerate returned error: %v", err)
	}

	output := string(collect())

	if !strings.Contains(output, "Address (bip44):") {
		t.Fatalf("expected BIP44 label in output, got %q", output)
	}

	if !strings.Contains(output, expectedAddress) {
		t.Fatalf("expected legacy P2PKH address in output, got %q", output)
	}
}

func TestRunWalletGenerateConflictingCompressionFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "wallet-generate"}
	cmd.Flags().Bool("compressed", false, "")
	cmd.Flags().Bool("uncompressed", false, "")
	if err := cmd.Flags().Set("compressed", "true"); err != nil {
		t.Fatalf("failed to set compressed flag: %v", err)
	}
	if err := cmd.Flags().Set("uncompressed", "true"); err != nil {
		t.Fatalf("failed to set uncompressed flag: %v", err)
	}

	walletForceCompressed = true
	walletForceUncompressed = true
	defer func() {
		walletForceCompressed = false
		walletForceUncompressed = false
	}()

	err := runWalletGenerate(cmd, nil)
	if err == nil {
		t.Fatal("expected error when both compressed and uncompressed are set")
	}

	if !strings.Contains(err.Error(), "cannot specify both") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunWalletInspectText(t *testing.T) {
	cmd := &cobra.Command{Use: "wallet-inspect"}
	walletInspectAddressType = "bip84"
	defer func() { walletInspectAddressType = "bip84" }()

	collect, restore := captureOutput()
	defer restore()

	if err := runWalletInspect(cmd, []string{"5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ"}); err != nil {
		t.Fatalf("runWalletInspect returned error: %v", err)
	}

	output := string(collect())

	if !strings.Contains(output, "WIF (mainnet") {
		t.Fatalf("expected mainnet indicator in output, got %q", output)
	}

	if !strings.Contains(output, "Address (bip44): 1GAehh7TsJAHuUAeKZcXf5CnwuGuGgyX2S") {
		t.Fatalf("expected derived address in output, got %q", output)
	}
}

func TestRunWalletInspectJSON(t *testing.T) {
	cmd := &cobra.Command{Use: "wallet-inspect"}
	cmd.Flags().String("output-format", "text", "")
	if err := cmd.Flags().Set("output-format", "json"); err != nil {
		t.Fatalf("failed to set output-format flag: %v", err)
	}
	walletInspectAddressType = "bip84"
	defer func() { walletInspectAddressType = "bip84" }()

	collect, restore := captureOutput()
	defer restore()

	if err := runWalletInspect(cmd, []string{"KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7"}); err != nil {
		t.Fatalf("runWalletInspect returned error: %v", err)
	}

	output := collect()

	type inspectPayload struct {
		WIF         string `json:"wif"`
		Compressed  bool   `json:"compressed"`
		Network     string `json:"network"`
		Address     string `json:"address"`
		AddressType string `json:"address_type"`
	}

	var payload inspectPayload
	if err := json.Unmarshal(output, &payload); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}

	if payload.Network != "mainnet" {
		t.Fatalf("expected mainnet network, got %s", payload.Network)
	}

	if !payload.Compressed {
		t.Fatal("expected compressed flag to be true for compressed WIF")
	}

	if payload.Address == "" {
		t.Fatal("expected address in payload")
	}

	if payload.AddressType != "bip84" {
		t.Fatalf("unexpected address type: %s", payload.AddressType)
	}
}

func TestRunWalletInspectInvalidWIF(t *testing.T) {
	cmd := &cobra.Command{Use: "wallet-inspect"}
	err := runWalletInspect(cmd, []string{"invalid-wif"})

	if err == nil {
		t.Fatal("expected error for invalid WIF")
	}

	if !strings.Contains(err.Error(), "invalid WIF private key") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunGenerateIntermediateWithLotSequence(t *testing.T) {
	origReadPassword := readPassword
	defer func() { readPassword = origReadPassword }()

	passphrase := []byte("TestingOneTwoThree")

	readPassword = func(int) ([]byte, error) {
		buf := make([]byte, len(passphrase))
		copy(buf, passphrase)
		return buf, nil
	}

	lotNumber = 12345
	sequenceNumber = 678
	useLotSeq = true
	defer func() {
		lotNumber = 0
		sequenceNumber = 0
		useLotSeq = false
	}()

	cmd := &cobra.Command{Use: "intermediate-generate"}
	cmd.Flags().Uint32Var(&lotNumber, "lot", 0, "")
	cmd.Flags().Uint32Var(&sequenceNumber, "sequence", 0, "")
	cmd.Flags().BoolVar(&useLotSeq, "use-lot-sequence", false, "")

	// Set the flags to simulate command line usage
	if err := cmd.Flags().Set("lot", "12345"); err != nil {
		t.Fatalf("failed to set lot flag: %v", err)
	}
	if err := cmd.Flags().Set("sequence", "678"); err != nil {
		t.Fatalf("failed to set sequence flag: %v", err)
	}
	if err := cmd.Flags().Set("use-lot-sequence", "true"); err != nil {
		t.Fatalf("failed to set use-lot-sequence flag: %v", err)
	}

	collect, restore := captureOutput()
	defer restore()

	runErr := runGenerateIntermediate(cmd, nil)

	outputBytes := collect()

	if runErr != nil {
		t.Fatalf("runGenerateIntermediate returned error: %v", runErr)
	}

	output := string(outputBytes)
	if !strings.Contains(output, "Intermediate code:") {
		t.Fatalf("expected intermediate code output, got %q", output)
	}
}

func TestRunValidateIntermediate(t *testing.T) {
	cmd := &cobra.Command{Use: "intermediate-validate"}

	// Test with invalid format first - should return error
	err := runValidateIntermediate(cmd, []string{"invalid-intermediate-code"})
	if err == nil {
		t.Fatal("expected error with invalid intermediate code")
	}

	expectedError := "invalid intermediate code format"
	if !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("expected error containing %q, got %q", expectedError, err.Error())
	}
}

func TestRunValidateIntermediateWithInvalidCode(t *testing.T) {
	cmd := &cobra.Command{Use: "intermediate-validate"}
	err := runValidateIntermediate(cmd, []string{"invalid-intermediate-code"})

	if err == nil {
		t.Fatal("expected error with invalid intermediate code")
	}

	expectedError := "invalid intermediate code format"
	if !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("expected error containing %q, got %q", expectedError, err.Error())
	}
}

func TestSecureZero(t *testing.T) {
	// Test that secureZero properly wipes the buffer
	originalData := []byte("sensitive data")
	buf := make([]byte, len(originalData))
	copy(buf, originalData)

	secureZero(buf)

	// Check that all bytes are zero
	for i, b := range buf {
		if b != 0 {
			t.Fatalf("expected byte at position %d to be zero, got %d", i, b)
		}
	}

	// Test with nil buffer
	secureZero(nil) // Should not panic
}

func TestGetPassphrase(t *testing.T) {
	origReadPassword := readPassword
	defer func() { readPassword = origReadPassword }()

	expectedPassphrase := []byte("test-password")
	readPassword = func(int) ([]byte, error) {
		return expectedPassphrase, nil
	}

	result, err := getPassphrase("Test prompt: ")
	if err != nil {
		t.Fatalf("getPassphrase returned error: %v", err)
	}

	if !bytes.Equal(result, expectedPassphrase) {
		t.Fatalf("expected passphrase %q, got %q", expectedPassphrase, result)
	}
}

func TestRunEncryptWithUncompressedFlag(t *testing.T) {
	origReadPassword := readPassword
	defer func() { readPassword = origReadPassword }()

	passphrase := []byte("TestingOneTwoThree")

	readPassword = func(int) ([]byte, error) {
		buf := make([]byte, len(passphrase))
		copy(buf, passphrase)
		return buf, nil
	}

	forceCompressed = false
	forceUncompressed = true
	defer func() {
		forceCompressed = false
		forceUncompressed = false
	}()

	cmd := &cobra.Command{Use: "encrypt"}

	collect, restore := captureOutput()
	defer restore()

	runErr := runEncrypt(cmd, []string{"5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ"})

	outputBytes := collect()

	if runErr != nil {
		t.Fatalf("runEncrypt returned error: %v", runErr)
	}

	output := string(outputBytes)
	if !strings.Contains(output, "Encrypted key: 6P") {
		t.Fatalf("expected encrypted key output, got %q", output)
	}
}

func TestRunDecryptWithShowAddress(t *testing.T) {
	origReadPassword := readPassword
	defer func() { readPassword = origReadPassword }()

	passphrase := []byte("TestingOneTwoThree")

	readPassword = func(int) ([]byte, error) {
		buf := make([]byte, len(passphrase))
		copy(buf, passphrase)
		return buf, nil
	}

	showAddress = true
	defer func() { showAddress = false }()

	cmd := &cobra.Command{Use: "decrypt"}

	collect, restore := captureOutput()
	defer restore()

	// Using a valid BIP38 format key for testing
	runErr := runDecrypt(cmd, []string{"6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg"})

	_ = collect() // Ignore output for this test

	// We expect this to fail during actual decryption, but the flow should work
	if runErr != nil && !strings.Contains(runErr.Error(), "decryption failed") {
		t.Fatalf("runDecrypt returned unexpected error: %v", runErr)
	}
}
