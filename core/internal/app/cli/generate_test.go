package cli

import (
	"testing"

	"github.com/carlosrabelo/bip38cli/core/internal/domain/bip38"
)

func TestGenerateCmd(t *testing.T) {
	// Test the actual generate command functionality
	t.Run("generate_key_functionality", func(t *testing.T) {
		// Test compressed key generation
		key, err := bip38.GenerateKey(true)
		if err != nil {
			t.Fatalf("GenerateKey() error = %v", err)
		}

		if !key.Compressed {
			t.Error("GenerateKey() should return compressed key when requested")
		}

		if key.PrivateKey == "" {
			t.Error("GenerateKey() returned empty private key")
		}

		if key.Address == "" {
			t.Error("GenerateKey() returned empty address")
		}

		// Test uncompressed key generation
		key2, err := bip38.GenerateKey(false)
		if err != nil {
			t.Fatalf("GenerateKey() error = %v", err)
		}

		if key2.Compressed {
			t.Error("GenerateKey() should return uncompressed key when requested")
		}

		// Keys should be different
		if key.PrivateKey == key2.PrivateKey {
			t.Error("GenerateKey() should produce different keys")
		}
	})
}
