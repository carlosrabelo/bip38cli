//go:build integration
// +build integration

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	binaryPath     = "../bin/bip38cli"
	testPassphrase = "MySecretPassphrase123!"
	testWIF        = "5Kb8kLf9zgWQnogidDA76MzPL6TsZZY36hWXMssSzNydYXYB9KF"
)

func TestIntegrationEncryptDecrypt(t *testing.T) {
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skip("Binary not found, run 'make build' first")
	}

	// Test encryption
	encryptCmd := exec.Command(binaryPath, "encrypt", "--wif", testWIF)
	encryptCmd.Stdin = strings.NewReader(testPassphrase + "\n")

	encryptOutput, err := encryptCmd.CombinedOutput()
	require.NoError(t, err, "Encryption failed: %s", string(encryptOutput))

	encryptedKey := strings.TrimSpace(string(encryptOutput))
	assert.True(t, len(encryptedKey) > 0, "Encrypted key should not be empty")
	assert.True(t, strings.HasPrefix(encryptedKey, "6P"), "Encrypted key should start with 6P")

	// Test decryption
	decryptCmd := exec.Command(binaryPath, "decrypt", "--bip38", encryptedKey)
	decryptCmd.Stdin = strings.NewReader(testPassphrase + "\n")

	decryptOutput, err := decryptCmd.CombinedOutput()
	require.NoError(t, err, "Decryption failed: %s", string(decryptOutput))

	decryptedWIF := strings.TrimSpace(string(decryptOutput))
	assert.Equal(t, testWIF, decryptedWIF, "Decrypted WIF should match original")
}

func TestIntegrationWalletGeneration(t *testing.T) {
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skip("Binary not found, run 'make build' first")
	}

	// Test wallet generation
	cmd := exec.Command(binaryPath, "wallet", "generate", "--json")
	cmd.Stdin = strings.NewReader(testPassphrase + "\n")

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Wallet generation failed: %s", string(output))

	// Should contain JSON output
	assert.Contains(t, string(output), "\"wif\":")
	assert.Contains(t, string(output), "\"bip38\":")
	assert.Contains(t, string(output), "\"address\":")
}

func TestIntegrationIntermediateCode(t *testing.T) {
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skip("Binary not found, run 'make build' first")
	}

	// Test intermediate code generation
	cmd := exec.Command(binaryPath, "intermediate", "generate",
		"--passphrase", testPassphrase,
		"--lot", "123456",
		"--sequence", "1")

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Intermediate code generation failed: %s", string(output))

	intermediateCode := strings.TrimSpace(string(output))
	assert.True(t, len(intermediateCode) > 0, "Intermediate code should not be empty")
	assert.True(t, strings.HasPrefix(intermediateCode, "6P"), "Intermediate code should start with 6P")

	// Test validation
	validateCmd := exec.Command(binaryPath, "intermediate", "validate", "--code", intermediateCode)
	validateOutput, err := validateCmd.CombinedOutput()
	require.NoError(t, err, "Intermediate code validation failed: %s", string(validateOutput))

	assert.Contains(t, string(validateOutput), "valid")
}

func TestIntegrationErrorHandling(t *testing.T) {
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skip("Binary not found, run 'make build' first")
	}

	tests := []struct {
		name     string
		args     []string
		input    string
		contains string
	}{
		{
			name:     "Invalid WIF",
			args:     []string{"encrypt", "--wif", "invalid"},
			input:    testPassphrase + "\n",
			contains: "malformed private key",
		},
		{
			name:     "Empty passphrase",
			args:     []string{"encrypt", "--wif", testWIF},
			input:    "\n",
			contains: "passphrase",
		},
		{
			name:     "Invalid BIP38 key",
			args:     []string{"decrypt", "--bip38", "invalid"},
			input:    testPassphrase + "\n",
			contains: "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			cmd.Stdin = strings.NewReader(tt.input)

			output, err := cmd.CombinedOutput()
			// Command should fail
			assert.Error(t, err)
			assert.Contains(t, string(output), tt.contains)
		})
	}
}

func TestIntegrationConcurrency(t *testing.T) {
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skip("Binary not found, run 'make build' first")
	}

	const numGoroutines = 10
	const operationsPerGoroutine = 5

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*operationsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < operationsPerGoroutine; j++ {
				// Test concurrent encryption
				cmd := exec.Command(binaryPath, "encrypt", "--wif", testWIF)
				cmd.Stdin = strings.NewReader(fmt.Sprintf("%s-%d-%d", testPassphrase, id, j) + "\n")

				output, err := cmd.CombinedOutput()
				if err != nil {
					errors <- fmt.Errorf("encryption failed for goroutine %d, operation %d: %s", id, j, string(output))
					return
				}

				encryptedKey := strings.TrimSpace(string(output))
				if !strings.HasPrefix(encryptedKey, "6P") {
					errors <- fmt.Errorf("invalid encrypted key format for goroutine %d, operation %d", id, j)
					return
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}
}

func TestIntegrationPerformance(t *testing.T) {
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skip("Binary not found, run 'make build' first")
	}

	const numOperations = 100
	start := time.Now()

	for i := 0; i < numOperations; i++ {
		cmd := exec.Command(binaryPath, "encrypt", "--wif", testWIF)
		cmd.Stdin = strings.NewReader(fmt.Sprintf("%s-%d", testPassphrase, i) + "\n")

		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Encryption %d failed: %s", i, string(output))

		encryptedKey := strings.TrimSpace(string(output))
		assert.True(t, strings.HasPrefix(encryptedKey, "6P"), "Invalid encrypted key format")
	}

	duration := time.Since(start)
	avgDuration := duration / numOperations

	t.Logf("Performance: %d operations in %v (avg: %v per operation)",
		numOperations, duration, avgDuration)

	// Should complete reasonably fast (less than 100ms per operation on average)
	assert.Less(t, avgDuration, 100*time.Millisecond, "Performance regression detected")
}

func TestIntegrationHelpAndVersion(t *testing.T) {
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skip("Binary not found, run 'make build' first")
	}

	// Test help
	cmd := exec.Command(binaryPath, "--help")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Help command failed")
	assert.Contains(t, string(output), "Usage:")
	assert.Contains(t, string(output), "Available Commands:")

	// Test version
	cmd = exec.Command(binaryPath, "--version")
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "Version command failed")
	assert.Contains(t, string(output), "bip38cli")
}

func TestIntegrationJSONOutput(t *testing.T) {
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skip("Binary not found, run 'make build' first")
	}

	// Test JSON output for wallet generation
	cmd := exec.Command(binaryPath, "wallet", "generate", "--json")
	cmd.Stdin = strings.NewReader(testPassphrase + "\n")

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "JSON wallet generation failed: %s", string(output))

	// Should be valid JSON
	assert.True(t, strings.HasPrefix(string(output), "{"), "Output should start with {")
	assert.True(t, strings.HasSuffix(string(output), "}\n"), "Output should end with }")
	assert.Contains(t, string(output), "\"wif\":")
	assert.Contains(t, string(output), "\"bip38\":")
}

func TestIntegrationNetworkSupport(t *testing.T) {
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skip("Binary not found, run 'make build' first")
	}

	networks := []string{"mainnet", "testnet", "regtest", "simnet"}

	for _, network := range networks {
		t.Run("network_"+network, func(t *testing.T) {
			cmd := exec.Command(binaryPath, "wallet", "generate",
				"--network", network, "--json")
			cmd.Stdin = strings.NewReader(testPassphrase + "\n")

			output, err := cmd.CombinedOutput()
			require.NoError(t, err, "Wallet generation failed for network %s: %s",
				network, string(output))

			// Should contain valid JSON
			assert.Contains(t, string(output), "\"wif\":")
			assert.Contains(t, string(output), "\"address\":")
		})
	}
}
