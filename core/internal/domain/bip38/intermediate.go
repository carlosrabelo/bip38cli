package bip38

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil/base58"
	"golang.org/x/crypto/scrypt"
)

// Magic bytes for intermediate passphrase code defined here
var (
	intermediateMagic = []byte{0x2C, 0xE9, 0xB3, 0xE1, 0xFF, 0x39, 0xE2, 0x51}
	intermediateLot   = []byte{0x2C, 0xE9, 0xB3, 0xE1, 0xFF, 0x39, 0xE2, 0x53}
)

// IntermediateCode represent BIP38 intermediate passphrase data bag
type IntermediateCode struct {
	Code      string
	OwnerSalt []byte
	PassPoint []byte
	HasLotSeq bool
	LotNumber *uint32
	SeqNumber *uint32
}

// GenerateIntermediateCode build one intermediate code using raw passphrase bytes
func GenerateIntermediateCode(passphrase []byte, lotNumber, sequenceNumber *uint32) (string, error) {
	var ownerSalt []byte
	var ownerEntropy []byte
	hasLotSeq := lotNumber != nil && sequenceNumber != nil

	if hasLotSeq {
		// Generate 4 random bytes for owner salt chunk
		ownerSalt = make([]byte, 4)
		if _, err := rand.Read(ownerSalt); err != nil {
			return "", fmt.Errorf("failed to generate random salt: %v", err)
		}

		// Encode lot and sequence numbers as 4-byte big-endian value
		lotSeq := make([]byte, 4)
		lotSeqValue := (*lotNumber * 4096) + *sequenceNumber
		binary.BigEndian.PutUint32(lotSeq, lotSeqValue)

		// Concatenate owner salt plus lot sequence bits
		ownerEntropy = append(ownerSalt, lotSeq...)
	} else {
		// Generate 8 random bytes for owner entropy (same bytes work as salt)
		ownerEntropy = make([]byte, 8)
		if _, err := rand.Read(ownerEntropy); err != nil {
			return "", fmt.Errorf("failed to generate random entropy: %v", err)
		}
		ownerSalt = ownerEntropy
	}

	// Derive a key from passphrase using scrypt parameters
	prefactor, err := scrypt.Key(passphrase, ownerSalt, 16384, 8, 8, 32)
	if err != nil {
		return "", fmt.Errorf("scrypt derivation failed: %v", err)
	}

	var passfactor []byte
	if hasLotSeq {
		// Take SHA256(SHA256(prefactor + ownerentropy)) like spec say
		combined := append(prefactor, ownerEntropy...)
		hash := sha256.Sum256(combined)
		hash2 := sha256.Sum256(hash[:])
		passfactor = hash2[:]
	} else {
		passfactor = prefactor
	}

	// Compute the elliptic curve point G * passfactor for passpoint
	privKey, _ := btcec.PrivKeyFromBytes(passfactor)
	if privKey.Key.IsZero() {
		return "", errors.New("invalid passfactor scalar")
	}
	passPoint := privKey.PubKey().SerializeCompressed()

	// Build intermediate code bytes step by step
	var magic []byte
	if hasLotSeq {
		magic = make([]byte, len(intermediateLot))
		copy(magic, intermediateLot)
	} else {
		magic = make([]byte, len(intermediateMagic))
		copy(magic, intermediateMagic)
	}

	intermediate := append(magic, ownerEntropy...)
	intermediate = append(intermediate, passPoint...)

	// Add checksum to finish structure
	hash := sha256.Sum256(intermediate)
	hash2 := sha256.Sum256(hash[:])
	intermediate = append(intermediate, hash2[:4]...)

	return base58.Encode(intermediate), nil
}

// ParseIntermediateCode parse a BIP38 intermediate passphrase code
func ParseIntermediateCode(code string) (*IntermediateCode, error) {
	decoded := base58.Decode(code)
	if len(decoded) < 49 {
		return nil, errors.New("invalid intermediate code length")
	}

	// Check magic bytes to detect lot/seq usage
	magic := decoded[:8]
	hasLotSeq := false

	if constantTimeEqual(magic, intermediateLot) {
		hasLotSeq = true
	} else if !constantTimeEqual(magic, intermediateMagic) {
		return nil, errors.New("invalid intermediate code magic")
	}

	// Verify checksum to protect data
	payload := decoded[:len(decoded)-4]
	checksum := decoded[len(decoded)-4:]
	hash := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash[:])

	if !constantTimeEqual(hash2[:4], checksum) {
		return nil, errors.New("invalid checksum")
	}

	// Extract components into typed struct
	ownerEntropy := decoded[8:16]
	passPoint := decoded[16:49]

	var ownerSalt []byte
	var lotNumber, seqNumber *uint32

	if hasLotSeq {
		ownerSalt = ownerEntropy[:4]
		lotSeqBytes := ownerEntropy[4:]
		lotSeqValue := binary.BigEndian.Uint32(lotSeqBytes)
		lot := lotSeqValue / 4096
		seq := lotSeqValue % 4096
		lotNumber = &lot
		seqNumber = &seq
	} else {
		ownerSalt = ownerEntropy
	}

	return &IntermediateCode{
		Code:      code,
		OwnerSalt: ownerSalt,
		PassPoint: passPoint,
		HasLotSeq: hasLotSeq,
		LotNumber: lotNumber,
		SeqNumber: seqNumber,
	}, nil
}

// IsValidIntermediateCode check if string look like valid intermediate code
func IsValidIntermediateCode(code string) bool {
	decoded := base58.Decode(code)
	if len(decoded) != 53 {
		return false
	}

	magic := decoded[:8]
	return constantTimeEqual(magic, intermediateMagic) || constantTimeEqual(magic, intermediateLot)
}
