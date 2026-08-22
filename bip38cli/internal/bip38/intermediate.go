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

// Magic bytes for intermediate passphrase codes (BIP38).
// With lot/sequence:  ... E2 51
// Without lot/sequence: ... E2 53
var (
	intermediateMagicLot   = []byte{0x2C, 0xE9, 0xB3, 0xE1, 0xFF, 0x39, 0xE2, 0x51}
	intermediateMagicNoLot = []byte{0x2C, 0xE9, 0xB3, 0xE1, 0xFF, 0x39, 0xE2, 0x53}
)

// IntermediateCode represent BIP38 intermediate passphrase data bag
type IntermediateCode struct {
	Code         string
	OwnerSalt    []byte // 4 bytes with lot/seq, 8 bytes without
	OwnerEntropy []byte // always 8 bytes
	PassPoint    []byte
	HasLotSeq    bool
	LotNumber    *uint32
	SeqNumber    *uint32
}

// GenerateIntermediateCode build one intermediate code using raw passphrase bytes
func GenerateIntermediateCode(passphrase []byte, lotNumber, sequenceNumber *uint32) (string, error) {
	passphrase = normalizePassphrase(passphrase)

	var ownerSalt []byte
	var ownerEntropy []byte
	hasLotSeq := lotNumber != nil && sequenceNumber != nil

	if hasLotSeq {
		if *lotNumber > 1048575 {
			return "", errors.New("lot number must be between 0 and 1048575")
		}
		if *sequenceNumber > 4095 {
			return "", errors.New("sequence number must be between 0 and 4095")
		}

		ownerSalt = make([]byte, 4)
		if _, err := rand.Read(ownerSalt); err != nil {
			return "", fmt.Errorf("failed to generate random salt: %w", err)
		}

		lotSeq := make([]byte, 4)
		lotSeqValue := (*lotNumber * 4096) + *sequenceNumber
		binary.BigEndian.PutUint32(lotSeq, lotSeqValue)

		ownerEntropy = make([]byte, 0, 8)
		ownerEntropy = append(ownerEntropy, ownerSalt...)
		ownerEntropy = append(ownerEntropy, lotSeq...)
	} else {
		ownerEntropy = make([]byte, 8)
		if _, err := rand.Read(ownerEntropy); err != nil {
			return "", fmt.Errorf("failed to generate random entropy: %w", err)
		}
		ownerSalt = ownerEntropy
	}

	prefactor, err := scrypt.Key(passphrase, ownerSalt, 16384, 8, 8, 32)
	if err != nil {
		return "", fmt.Errorf("scrypt derivation failed: %w", err)
	}
	defer zeroBytes(prefactor)

	var passfactor []byte
	if hasLotSeq {
		combined := append(append([]byte{}, prefactor...), ownerEntropy...)
		hash := sha256.Sum256(combined)
		hash2 := sha256.Sum256(hash[:])
		passfactor = append([]byte{}, hash2[:]...)
		zeroBytes(combined)
	} else {
		passfactor = append([]byte{}, prefactor...)
	}
	defer zeroBytes(passfactor)

	privKey, _ := btcec.PrivKeyFromBytes(passfactor)
	if privKey.Key.IsZero() {
		return "", errors.New("invalid passfactor scalar")
	}
	passPoint := privKey.PubKey().SerializeCompressed()

	var magic []byte
	if hasLotSeq {
		magic = intermediateMagicLot
	} else {
		magic = intermediateMagicNoLot
	}

	intermediate := make([]byte, 0, 53)
	intermediate = append(intermediate, magic...)
	intermediate = append(intermediate, ownerEntropy...)
	intermediate = append(intermediate, passPoint...)

	hash := sha256.Sum256(intermediate)
	hash2 := sha256.Sum256(hash[:])
	intermediate = append(intermediate, hash2[:4]...)

	return base58.Encode(intermediate), nil
}

// ParseIntermediateCode parse a BIP38 intermediate passphrase code
func ParseIntermediateCode(code string) (*IntermediateCode, error) {
	decoded := base58.Decode(code)
	if len(decoded) != 53 {
		return nil, errors.New("invalid intermediate code length")
	}

	magic := decoded[:8]
	hasLotSeq := false

	if constantTimeEqual(magic, intermediateMagicLot) {
		hasLotSeq = true
	} else if !constantTimeEqual(magic, intermediateMagicNoLot) {
		return nil, errors.New("invalid intermediate code magic")
	}

	payload := decoded[:len(decoded)-4]
	checksum := decoded[len(decoded)-4:]
	hash := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash[:])

	if !constantTimeEqual(hash2[:4], checksum) {
		return nil, errors.New("invalid checksum")
	}

	ownerEntropy := make([]byte, 8)
	copy(ownerEntropy, decoded[8:16])
	passPoint := make([]byte, 33)
	copy(passPoint, decoded[16:49])

	var ownerSalt []byte
	var lotNumber, seqNumber *uint32

	if hasLotSeq {
		ownerSalt = make([]byte, 4)
		copy(ownerSalt, ownerEntropy[:4])
		lotSeqValue := binary.BigEndian.Uint32(ownerEntropy[4:])
		lot := lotSeqValue / 4096
		seq := lotSeqValue % 4096
		lotNumber = &lot
		seqNumber = &seq
	} else {
		ownerSalt = make([]byte, 8)
		copy(ownerSalt, ownerEntropy)
	}

	return &IntermediateCode{
		Code:         code,
		OwnerSalt:    ownerSalt,
		OwnerEntropy: ownerEntropy,
		PassPoint:    passPoint,
		HasLotSeq:    hasLotSeq,
		LotNumber:    lotNumber,
		SeqNumber:    seqNumber,
	}, nil
}

// IsValidIntermediateCode check if string look like valid intermediate code
func IsValidIntermediateCode(code string) bool {
	_, err := ParseIntermediateCode(code)
	return err == nil
}
