package bip38

import (
	"testing"

	"github.com/btcsuite/btcd/btcutil/base58"
)

func TestGenerateIntermediateCode(t *testing.T) {
	tests := []struct {
		name       string
		passphrase string
		lot        *uint32
		seq        *uint32
		wantErr    bool
	}{
		{
			name:       "simple passphrase without lot/seq",
			passphrase: "TestingOneTwoThree",
			lot:        nil,
			seq:        nil,
			wantErr:    false,
		},
		{
			name:       "simple passphrase with lot/seq",
			passphrase: "TestingOneTwoThree",
			lot:        uint32Ptr(123),
			seq:        uint32Ptr(456),
			wantErr:    false,
		},
	}

	if !testing.Short() {
		tests = append(tests,
			struct {
				name       string
				passphrase string
				lot        *uint32
				seq        *uint32
				wantErr    bool
			}{
				name:       "empty passphrase",
				passphrase: "",
				lot:        nil,
				seq:        nil,
				wantErr:    false, // BIP38 allow empty passphrases
			},
			struct {
				name       string
				passphrase string
				lot        *uint32
				seq        *uint32
				wantErr    bool
			}{
				name:       "complex passphrase",
				passphrase: "!@#$%^&*()_+-=[]{}|;':\",./<>?`~",
				lot:        nil,
				seq:        nil,
				wantErr:    false,
			},
			struct {
				name       string
				passphrase string
				lot        *uint32
				seq        *uint32
				wantErr    bool
			}{
				name:       "max lot and seq values",
				passphrase: "test",
				lot:        uint32Ptr(1048575), // 2^20 minus 1
				seq:        uint32Ptr(4095),    // 2^12 minus 1
				wantErr:    false,
			},
		)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, err := GenerateIntermediateCode([]byte(tt.passphrase), tt.lot, tt.seq)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateIntermediateCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify generated code stay valid
				if !IsValidIntermediateCode(code) {
					t.Errorf("Generated intermediate code is not valid: %s", code)
				}

				// Parse the code to verify structure
				parsed, err := ParseIntermediateCode(code)
				if err != nil {
					t.Errorf("Failed to parse generated intermediate code: %v", err)
				}

				// Verify lot/seq encoding looks right
				if tt.lot != nil && tt.seq != nil {
					if !parsed.HasLotSeq {
						t.Error("Expected HasLotSeq to be true")
					}
					if parsed.LotNumber == nil || *parsed.LotNumber != *tt.lot {
						t.Errorf("Lot number mismatch: got %v, want %v", parsed.LotNumber, *tt.lot)
					}
					if parsed.SeqNumber == nil || *parsed.SeqNumber != *tt.seq {
						t.Errorf("Sequence number mismatch: got %v, want %v", parsed.SeqNumber, *tt.seq)
					}
				} else {
					if parsed.HasLotSeq {
						t.Error("Expected HasLotSeq to be false")
					}
					if parsed.LotNumber != nil || parsed.SeqNumber != nil {
						t.Error("Expected lot and sequence numbers to be nil")
					}
				}
			}
		})
	}
}

func TestIsValidIntermediateCode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "empty string",
			code:     "",
			expected: false,
		},
		{
			name:     "too short",
			code:     "passphrase",
			expected: false,
		},
		{
			name:     "invalid base58",
			code:     "passphrase0IL",
			expected: false,
		},
		{
			name:     "wrong magic",
			code:     "6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidIntermediateCode(tt.code)
			if result != tt.expected {
				t.Errorf("IsValidIntermediateCode(%q) = %v, want %v", tt.code, result, tt.expected)
			}
		})
	}
}

func TestParseIntermediateCode(t *testing.T) {
	// Generate a known code for testing scenario
	passphrase := "TestingOneTwoThree"
	lot := uint32Ptr(123)
	seq := uint32Ptr(456)

	code, err := GenerateIntermediateCode([]byte(passphrase), lot, seq)
	if err != nil {
		t.Fatalf("Failed to generate test intermediate code: %v", err)
	}

	// Parse it back for checks
	parsed, err := ParseIntermediateCode(code)
	if err != nil {
		t.Fatalf("Failed to parse intermediate code: %v", err)
	}

	// Verify the parsed data carefully
	if parsed.Code != code {
		t.Errorf("Code mismatch: got %s, want %s", parsed.Code, code)
	}

	if !parsed.HasLotSeq {
		t.Error("Expected HasLotSeq to be true")
	}

	if parsed.LotNumber == nil || *parsed.LotNumber != *lot {
		t.Errorf("Lot number mismatch: got %v, want %v", parsed.LotNumber, *lot)
	}

	if parsed.SeqNumber == nil || *parsed.SeqNumber != *seq {
		t.Errorf("Sequence number mismatch: got %v, want %v", parsed.SeqNumber, *seq)
	}

	if len(parsed.OwnerSalt) != 4 {
		t.Errorf("Expected owner salt length 4, got %d", len(parsed.OwnerSalt))
	}

	if len(parsed.PassPoint) != 33 {
		t.Errorf("Expected pass point length 33, got %d", len(parsed.PassPoint))
	}
}

func TestParseIntermediateCodeNoLotSeq(t *testing.T) {
	// Generate a code without lot/seq scenario
	passphrase := "TestingOneTwoThree"

	code, err := GenerateIntermediateCode([]byte(passphrase), nil, nil)
	if err != nil {
		t.Fatalf("Failed to generate test intermediate code: %v", err)
	}

	// Parse it back again
	parsed, err := ParseIntermediateCode(code)
	if err != nil {
		t.Fatalf("Failed to parse intermediate code: %v", err)
	}

	// Verify the parsed data for no lot case
	if parsed.HasLotSeq {
		t.Error("Expected HasLotSeq to be false")
	}

	if parsed.LotNumber != nil {
		t.Errorf("Expected lot number to be nil, got %v", parsed.LotNumber)
	}

	if parsed.SeqNumber != nil {
		t.Errorf("Expected sequence number to be nil, got %v", parsed.SeqNumber)
	}

	if len(parsed.OwnerSalt) != 8 {
		t.Errorf("Expected owner salt length 8, got %d", len(parsed.OwnerSalt))
	}
}

func TestParseInvalidIntermediateCode(t *testing.T) {
	t.Run("too short", func(t *testing.T) {
		_, err := ParseIntermediateCode("passphrase")
		if err == nil {
			t.Fatal("expected error for short code")
		}
		if err.Error() != "invalid intermediate code length" {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("invalid checksum", func(t *testing.T) {
		valid, err := GenerateIntermediateCode([]byte("TestingOneTwoThree"), nil, nil)
		if err != nil {
			t.Fatalf("failed to generate code: %v", err)
		}

		decoded := base58.Decode(valid)
		if len(decoded) == 0 {
			t.Fatal("failed to decode generated code")
		}

		decoded[len(decoded)-1] ^= 0x01
		tampered := base58.Encode(decoded)

		_, parseErr := ParseIntermediateCode(tampered)
		if parseErr == nil {
			t.Fatal("expected checksum error for tampered code")
		}
		if parseErr.Error() != "invalid checksum" {
			t.Fatalf("unexpected error: %v", parseErr)
		}
	})
}

func BenchmarkGenerateIntermediateCode(b *testing.B) {
	passphrase := []byte("TestingOneTwoThree")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GenerateIntermediateCode(passphrase, nil, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGenerateIntermediateCodeWithLotSeq(b *testing.B) {
	passphrase := []byte("TestingOneTwoThree")
	lot := uint32Ptr(123)
	seq := uint32Ptr(456)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GenerateIntermediateCode(passphrase, lot, seq)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseIntermediateCode(b *testing.B) {
	// Generate a test code one time before benchmark loop
	code, err := GenerateIntermediateCode([]byte("TestingOneTwoThree"), nil, nil)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseIntermediateCode(code)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Helper function to create uint32 pointer quickly
func uint32Ptr(v uint32) *uint32 {
	return &v
}
