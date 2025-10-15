package bip38

import (
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

func TestGenerateKey(t *testing.T) {
	tests := []struct {
		name       string
		compressed bool
	}{
		{
			name:       "compressed key",
			compressed: true,
		},
		{
			name:       "uncompressed key",
			compressed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := GenerateKey(tt.compressed)
			if err != nil {
				t.Fatalf("GenerateKey() error = %v", err)
			}

			// Verify private key format
			if key.PrivateKey == "" {
				t.Error("GenerateKey() returned empty private key")
			}

			// Verify address format
			if key.Address == "" {
				t.Error("GenerateKey() returned empty address")
			}

			// Verify compressed flag matches
			if key.Compressed != tt.compressed {
				t.Errorf("GenerateKey() compressed = %v, want %v", key.Compressed, tt.compressed)
			}

			// Verify WIF format by parsing it
			wif, err := btcutil.DecodeWIF(key.PrivateKey)
			if err != nil {
				t.Errorf("GenerateKey() returned invalid WIF: %v", err)
			}

			// Verify compression matches WIF
			if wif.CompressPubKey != tt.compressed {
				t.Errorf("GenerateKey() WIF compression = %v, want %v", wif.CompressPubKey, tt.compressed)
			}

			// Verify network is mainnet
			if !wif.IsForNet(&chaincfg.MainNetParams) {
				t.Error("GenerateKey() should use mainnet parameters")
			}

			// Verify address corresponds to private key
			var pubKeyBytes []byte
			if tt.compressed {
				pubKeyBytes = wif.PrivKey.PubKey().SerializeCompressed()
			} else {
				pubKeyBytes = wif.PrivKey.PubKey().SerializeUncompressed()
			}

			addressPubKey, err := btcutil.NewAddressPubKey(pubKeyBytes, &chaincfg.MainNetParams)
			if err != nil {
				t.Errorf("Failed to create address from generated key: %v", err)
			}

			expectedAddress := addressPubKey.EncodeAddress()
			if key.Address != expectedAddress {
				t.Errorf("GenerateKey() address = %v, want %v", key.Address, expectedAddress)
			}
		})
	}
}

func TestGenerateKeyUniqueness(t *testing.T) {
	// Generate multiple keys and verify they're unique
	keys := make(map[string]bool)

	for i := 0; i < 10; i++ {
		key, err := GenerateKey(true)
		if err != nil {
			t.Fatalf("GenerateKey() error = %v", err)
		}

		if keys[key.PrivateKey] {
			t.Errorf("GenerateKey() produced duplicate key: %s", key.PrivateKey)
		}
		keys[key.PrivateKey] = true
	}
}

func TestGeneratedKeyStruct(t *testing.T) {
	key, err := GenerateKey(true)
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}

	// Test JSON tags work correctly
	if key.PrivateKey == "" {
		t.Error("PrivateKey field is empty")
	}
	if key.Address == "" {
		t.Error("Address field is empty")
	}
}
