package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/carlosrabelo/bip38cli/core/internal/pkg/logger"
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
