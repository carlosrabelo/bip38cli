package bip38

import (
	"testing"

	"github.com/btcsuite/btcd/btcutil"
)

func TestIsBIP38Format(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{
			name:     "valid BIP38 key",
			key:      "6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg",
			expected: true,
		},
		{
			name:     "invalid prefix",
			key:      "5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ",
			expected: false,
		},
		{
			name:     "too short",
			key:      "6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2Zo",
			expected: false,
		},
		{
			name:     "too long",
			key:      "6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGgX",
			expected: false,
		},
		{
			name:     "invalid characters",
			key:      "6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2Z0Gg",
			expected: false,
		},
		{
			name:     "empty string",
			key:      "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsBIP38Format(tt.key)
			if result != tt.expected {
				t.Errorf("IsBIP38Format(%q) = %v, want %v", tt.key, result, tt.expected)
			}
		})
	}
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	tests := []struct {
		name       string
		wifString  string
		passphrase string
		compressed bool
	}{
		{
			name:       "uncompressed key",
			wifString:  "5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ",
			passphrase: "TestingOneTwoThree",
			compressed: false,
		},
		{
			name:       "compressed key",
			wifString:  "KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7",
			passphrase: "TestingOneTwoThree",
			compressed: true,
		},
		{
			name:       "testnet compressed key",
			wifString:  "cMai6KJ8sHnNejZctqjiYdg2KSLeiaKcuTbZQrJNEjmMY5JQw6eP",
			passphrase: "TestingOneTwoThree",
			compressed: true,
		},
	}

	if !testing.Short() {
		tests = append(tests,
			struct {
				name       string
				wifString  string
				passphrase string
				compressed bool
			}{
				name:       "simple passphrase",
				wifString:  "5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ",
				passphrase: "test",
				compressed: false,
			},
			struct {
				name       string
				wifString  string
				passphrase string
				compressed bool
			}{
				name:       "complex passphrase",
				wifString:  "5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ",
				passphrase: "!@#$%^&*()_+-=[]{}|;':\",./<>?`~",
				compressed: false,
			},
		)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse original WIF input
			originalWIF, err := btcutil.DecodeWIF(tt.wifString)
			if err != nil {
				t.Fatalf("Failed to decode WIF: %v", err)
			}

			// Set compression flag as needed
			originalWIF.CompressPubKey = tt.compressed

			// Encrypt the key with tested passphrase
			encryptedKey, err := EncryptKey(originalWIF, []byte(tt.passphrase))
			if err != nil {
				t.Fatalf("Failed to encrypt key: %v", err)
			}

			// Verify result stay in BIP38 format
			if !IsBIP38Format(encryptedKey) {
				t.Errorf("Encrypted key is not in valid BIP38 format: %s", encryptedKey)
			}

			// Decrypt the key back
			decryptedWIF, err := DecryptKey(encryptedKey, []byte(tt.passphrase))
			if err != nil {
				t.Fatalf("Failed to decrypt key: %v", err)
			}

			// Verify the keys match after round trip
			if originalWIF.String() != decryptedWIF.String() {
				t.Errorf("Keys do not match after roundtrip:\nOriginal:  %s\nDecrypted: %s",
					originalWIF.String(), decryptedWIF.String())
			}

			// Verify compression flag survive
			if originalWIF.CompressPubKey != decryptedWIF.CompressPubKey {
				t.Errorf("Compression setting mismatch: original=%v, decrypted=%v",
					originalWIF.CompressPubKey, decryptedWIF.CompressPubKey)
			}
		})
	}
}

func TestDecryptWithWrongPassphrase(t *testing.T) {
	// Use a known encrypted key for negative case
	encryptedKey := "6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg"
	wrongPassphrase := "wrongpassphrase"

	_, err := DecryptKey(encryptedKey, []byte(wrongPassphrase))
	if err == nil {
		t.Error("Expected error when decrypting with wrong passphrase, got nil")
	}

	// The error should indicate incorrect passphrase usage
	if err.Error() != "incorrect passphrase" {
		t.Errorf("Expected 'incorrect passphrase' error, got: %v", err)
	}
}

func TestDecryptInvalidKey(t *testing.T) {
	tests := []struct {
		name         string
		encryptedKey string
		passphrase   string
		expectedErr  string
	}{
		{
			name:         "invalid format",
			encryptedKey: "5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ",
			passphrase:   "test",
			expectedErr:  "invalid BIP38 format",
		},
		{
			name:         "invalid base58",
			encryptedKey: "6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2Z0Gg",
			passphrase:   "test",
			expectedErr:  "invalid BIP38 format",
		},
		{
			name:         "empty key",
			encryptedKey: "",
			passphrase:   "test",
			expectedErr:  "invalid BIP38 format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecryptKey(tt.encryptedKey, []byte(tt.passphrase))
			if err == nil {
				t.Error("Expected error for invalid key, got nil")
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("Expected error '%s', got: %v", tt.expectedErr, err)
			}
		})
	}
}

func BenchmarkEncrypt(b *testing.B) {
	wif, _ := btcutil.DecodeWIF("5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ")
	passphrase := []byte("TestingOneTwoThree")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := EncryptKey(wif, passphrase)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecrypt(b *testing.B) {
	encryptedKey := "6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg"
	passphrase := []byte("TestingOneTwoThree")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DecryptKey(encryptedKey, passphrase)
		if err != nil {
			b.Fatal(err)
		}
	}
}
