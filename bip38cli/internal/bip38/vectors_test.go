package bip38

import (
	"encoding/hex"
	"testing"
	"unicode/utf8"

	"github.com/btcsuite/btcd/btcutil"
)

// Official BIP38 test vectors from https://github.com/bitcoin/bips/blob/master/bip-0038.mediawiki

func TestOfficialVectorsNoCompressionNoEC(t *testing.T) {
	tests := []struct {
		name       string
		passphrase string
		encrypted  string
		wif        string
	}{
		{
			name:       "TestingOneTwoThree",
			passphrase: "TestingOneTwoThree",
			encrypted:  "6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg",
			wif:        "5KN7MzqK5wt2TP1fQCYyHBtDrXdJuXbUzm4A9rKAteGu3Qi5CVR",
		},
		{
			name:       "Satoshi",
			passphrase: "Satoshi",
			encrypted:  "6PRNFFkZc2NZ6dJqFfhRoFNMR9Lnyj7dYGrzdgXXVMXcxoKTePPX1dWByq",
			wif:        "5HtasZ6ofTHP6HCwTqTkLDuLQisYPah7aUnSKfC7h4hMUVw2gi5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/decrypt", func(t *testing.T) {
			wif, err := DecryptKey(tt.encrypted, []byte(tt.passphrase))
			if err != nil {
				t.Fatalf("DecryptKey: %v", err)
			}
			if wif.String() != tt.wif {
				t.Fatalf("WIF = %s, want %s", wif.String(), tt.wif)
			}
		})
		t.Run(tt.name+"/encrypt", func(t *testing.T) {
			wif, err := btcutil.DecodeWIF(tt.wif)
			if err != nil {
				t.Fatalf("DecodeWIF: %v", err)
			}
			encrypted, err := EncryptKey(wif, []byte(tt.passphrase))
			if err != nil {
				t.Fatalf("EncryptKey: %v", err)
			}
			if encrypted != tt.encrypted {
				t.Fatalf("encrypted = %s, want %s", encrypted, tt.encrypted)
			}
		})
	}
}

func TestOfficialVectorNFCPassphrase(t *testing.T) {
	passphrase := string([]rune{0x03D2, 0x0301, 0x0000, 0x10400, 0x1F4A9})
	if !utf8.ValidString(passphrase) {
		t.Fatal("test passphrase is not valid UTF-8")
	}

	const (
		encrypted = "6PRW5o9FLp4gJDDVqJQKJFTpMvdsSGJxMYHtHaQBF3ooa8mwD69bapcDQn"
		wif       = "5Jajm8eQ22H3pGWLEVCXyvND8dQZhiQhoLJNKjYXk9roUFTMSZ4"
		nfcHex    = "cf9300f0909080f09f92a9"
	)

	normalized := normalizePassphrase([]byte(passphrase))
	if hex.EncodeToString(normalized) != nfcHex {
		t.Fatalf("NFC bytes = %x, want %s", normalized, nfcHex)
	}

	got, err := DecryptKey(encrypted, []byte(passphrase))
	if err != nil {
		t.Fatalf("DecryptKey: %v", err)
	}
	if got.String() != wif {
		t.Fatalf("WIF = %s, want %s", got.String(), wif)
	}

	decoded, err := btcutil.DecodeWIF(wif)
	if err != nil {
		t.Fatalf("DecodeWIF: %v", err)
	}
	enc, err := EncryptKey(decoded, []byte(passphrase))
	if err != nil {
		t.Fatalf("EncryptKey: %v", err)
	}
	if enc != encrypted {
		t.Fatalf("encrypted = %s, want %s", enc, encrypted)
	}
}

func TestOfficialVectorsCompressionNoEC(t *testing.T) {
	tests := []struct {
		name       string
		passphrase string
		encrypted  string
		wif        string
	}{
		{
			name:       "TestingOneTwoThree",
			passphrase: "TestingOneTwoThree",
			encrypted:  "6PYNKZ1EAgYgmQfmNVamxyXVWHzK5s6DGhwP4J5o44cvXdoY7sRzhtpUeo",
			wif:        "L44B5gGEpqEDRS9vVPz7QT35jcBG2r3CZwSwQ4fCewXAhAhqGVpP",
		},
		{
			name:       "Satoshi",
			passphrase: "Satoshi",
			encrypted:  "6PYLtMnXvfG3oJde97zRyLYFZCYizPU5T3LwgdYJz1fRhh16bU7u6PPmY7",
			wif:        "KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/decrypt", func(t *testing.T) {
			wif, err := DecryptKey(tt.encrypted, []byte(tt.passphrase))
			if err != nil {
				t.Fatalf("DecryptKey: %v", err)
			}
			if wif.String() != tt.wif {
				t.Fatalf("WIF = %s, want %s", wif.String(), tt.wif)
			}
			if !wif.CompressPubKey {
				t.Fatal("expected compressed WIF")
			}
		})
		t.Run(tt.name+"/encrypt", func(t *testing.T) {
			wif, err := btcutil.DecodeWIF(tt.wif)
			if err != nil {
				t.Fatalf("DecodeWIF: %v", err)
			}
			encrypted, err := EncryptKey(wif, []byte(tt.passphrase))
			if err != nil {
				t.Fatalf("EncryptKey: %v", err)
			}
			if encrypted != tt.encrypted {
				t.Fatalf("encrypted = %s, want %s", encrypted, tt.encrypted)
			}
		})
	}
}

func TestOfficialIntermediateMagicBytes(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		hasLotSeq bool
		lot       uint32
		seq       uint32
	}{
		{
			name:      "no lot TestingOneTwoThree",
			code:      "passphrasepxFy57B9v8HtUsszJYKReoNDV6VHjUSGt8EVJmux9n1J3Ltf1gRxyDGXqnf9qm",
			hasLotSeq: false,
		},
		{
			name:      "no lot Satoshi",
			code:      "passphraseoRDGAXTWzbp72eVbtUDdn1rwpgPUGjNZEc6CGBo8i5EC1FPW8wcnLdq4ThKzAS",
			hasLotSeq: false,
		},
		{
			name:      "lot MOLON LABE",
			code:      "passphraseaB8feaLQDENqCgr4gKZpmf4VoaT6qdjJNJiv7fsKvjqavcJxvuR1hy25aTu5sX",
			hasLotSeq: true,
			lot:       263183,
			seq:       1,
		},
		{
			name:      "lot Greek",
			code:      "passphrased3z9rQJHSyBkNBwTRPkUGNVEVrUAcfAXDyRU1V28ie6hNFbqDwbFBvsTK7yWVK",
			hasLotSeq: true,
			lot:       806938,
			seq:       1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !IsValidIntermediateCode(tt.code) {
				t.Fatal("expected valid intermediate code")
			}
			parsed, err := ParseIntermediateCode(tt.code)
			if err != nil {
				t.Fatalf("ParseIntermediateCode: %v", err)
			}
			if parsed.HasLotSeq != tt.hasLotSeq {
				t.Fatalf("HasLotSeq = %v, want %v", parsed.HasLotSeq, tt.hasLotSeq)
			}
			if tt.hasLotSeq {
				if parsed.LotNumber == nil || *parsed.LotNumber != tt.lot {
					t.Fatalf("lot = %v, want %d", parsed.LotNumber, tt.lot)
				}
				if parsed.SeqNumber == nil || *parsed.SeqNumber != tt.seq {
					t.Fatalf("seq = %v, want %d", parsed.SeqNumber, tt.seq)
				}
			}
		})
	}
}

func TestOfficialECMultiplyDecrypt(t *testing.T) {
	tests := []struct {
		name       string
		passphrase string
		encrypted  string
		wif        string
		confirm    string
		lot        *uint32
		seq        *uint32
		address    string
	}{
		{
			name:       "no lot TestingOneTwoThree",
			passphrase: "TestingOneTwoThree",
			encrypted:  "6PfQu77ygVyJLZjfvMLyhLMQbYnu5uguoJJ4kMCLqWwPEdfpwANVS76gTX",
			wif:        "5K4caxezwjGCGfnoPTZ8tMcJBLB7Jvyjv4xxeacadhq8nLisLR2",
			address:    "1PE6TQi6HTVNz5DLwB1LcpMBALubfuN2z2",
		},
		{
			name:       "no lot Satoshi",
			passphrase: "Satoshi",
			encrypted:  "6PfLGnQs6VZnrNpmVKfjotbnQuaJK4KZoPFrAjx1JMJUa1Ft8gnf5WxfKd",
			wif:        "5KJ51SgxWaAYR13zd9ReMhJpwrcX47xTJh2D3fGPG9CM8vkv5sH",
			address:    "1CqzrtZC6mXSAhoxtFwVjz8LtwLJjDYU3V",
		},
		{
			name:       "lot MOLON LABE",
			passphrase: "MOLON LABE",
			encrypted:  "6PgNBNNzDkKdhkT6uJntUXwwzQV8Rr2tZcbkDcuC9DZRsS6AtHts4Ypo1j",
			wif:        "5JLdxTtcTHcfYcmJsNVy1v2PMDx432JPoYcBTVVRHpPaxUrdtf8",
			confirm:    "cfrm38V8aXBn7JWA1ESmFMUn6erxeBGZGAxJPY4e36S9QWkzZKtaVqLNMgnifETYw7BPwWC9aPD",
			lot:        uint32Ptr(263183),
			seq:        uint32Ptr(1),
			address:    "1Jscj8ALrYu2y9TD8NrpvDBugPedmbj4Yh",
		},
		{
			name:       "lot Greek",
			passphrase: "ΜΟΛΩΝ ΛΑΒΕ",
			encrypted:  "6PgGWtx25kUg8QWvwuJAgorN6k9FbE25rv5dMRwu5SKMnfpfVe5mar2ngH",
			wif:        "5KMKKuUmAkiNbA3DazMQiLfDq47qs8MAEThm4yL8R2PhV1ov33D",
			confirm:    "cfrm38V8G4qq2ywYEFfWLD5Cc6msj9UwsG2Mj4Z6QdGJAFQpdatZLavkgRd1i4iBMdRngDqDs51",
			lot:        uint32Ptr(806938),
			seq:        uint32Ptr(1),
			address:    "1Lurmih3KruL4xDB5FmHof38yawNtP9oGf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/decrypt", func(t *testing.T) {
			wif, err := DecryptKey(tt.encrypted, []byte(tt.passphrase))
			if err != nil {
				t.Fatalf("DecryptKey: %v", err)
			}
			if wif.String() != tt.wif {
				t.Fatalf("WIF = %s, want %s", wif.String(), tt.wif)
			}
		})
		if tt.confirm != "" {
			t.Run(tt.name+"/confirm", func(t *testing.T) {
				result, err := VerifyConfirmationCode(tt.confirm, []byte(tt.passphrase))
				if err != nil {
					t.Fatalf("VerifyConfirmationCode: %v", err)
				}
				if result.Address != tt.address {
					t.Fatalf("address = %s, want %s", result.Address, tt.address)
				}
				if tt.lot != nil {
					if result.LotNumber == nil || *result.LotNumber != *tt.lot {
						t.Fatalf("lot = %v, want %d", result.LotNumber, *tt.lot)
					}
					if result.SeqNumber == nil || *result.SeqNumber != *tt.seq {
						t.Fatalf("seq = %v, want %d", result.SeqNumber, *tt.seq)
					}
				}
			})
		}
	}
}

func TestECMultiplyRoundtrip(t *testing.T) {
	passphrase := []byte("TestingOneTwoThree")
	code, err := GenerateIntermediateCode(passphrase, nil, nil)
	if err != nil {
		t.Fatalf("GenerateIntermediateCode: %v", err)
	}

	result, err := ECMultiplyEncrypt(code, false)
	if err != nil {
		t.Fatalf("ECMultiplyEncrypt: %v", err)
	}

	wif, err := DecryptKey(result.EncryptedKey, passphrase)
	if err != nil {
		t.Fatalf("DecryptKey: %v", err)
	}
	if wif.CompressPubKey {
		t.Fatal("expected uncompressed key")
	}

	confirm, err := VerifyConfirmationCode(result.ConfirmationCode, passphrase)
	if err != nil {
		t.Fatalf("VerifyConfirmationCode: %v", err)
	}
	if confirm.Address == "" {
		t.Fatal("expected confirmation address")
	}
}
