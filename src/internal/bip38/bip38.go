// Package bip38 implements the BIP38 standard for encryption and decryption of Bitcoin private keys.
package bip38

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/chaincfg"
	"golang.org/x/crypto/scrypt"
)

const (
	// Magic byte used by classic BIP38 flow
	bip38Magic = 0x01
	// Type byte telling we use non EC multiply form
	bip38Type = 0x42
)

var (
	// Regex used to check BIP38 string look correct
	bip38Regex = regexp.MustCompile(`^6P[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]{56}$`)

	supportedNetworks = []*chaincfg.Params{
		&chaincfg.MainNetParams,
		&chaincfg.TestNet3Params,
		&chaincfg.RegressionNetParams,
		&chaincfg.SimNetParams,
		&chaincfg.SigNetParams,
	}

	networkAliases = map[string]*chaincfg.Params{
		"main":       &chaincfg.MainNetParams,
		"mainnet":    &chaincfg.MainNetParams,
		"bitcoin":    &chaincfg.MainNetParams,
		"livenet":    &chaincfg.MainNetParams,
		"prod":       &chaincfg.MainNetParams,
		"production": &chaincfg.MainNetParams,

		"test":     &chaincfg.TestNet3Params,
		"testnet":  &chaincfg.TestNet3Params,
		"testnet3": &chaincfg.TestNet3Params,
		"tn3":      &chaincfg.TestNet3Params,

		"regtest":       &chaincfg.RegressionNetParams,
		"regression":    &chaincfg.RegressionNetParams,
		"regressionnet": &chaincfg.RegressionNetParams,

		"simnet": &chaincfg.SimNetParams,
		"sim":    &chaincfg.SimNetParams,

		"signet": &chaincfg.SigNetParams,
		"sig":    &chaincfg.SigNetParams,
	}
)

// EncryptedKey represents a BIP38 encrypted private key record with its properties.
type EncryptedKey struct {
	Key        string // The encrypted key in base58 format
	Compressed bool   // Whether the original key was compressed
	ECMultiply bool   // Whether EC multiplication mode was used
}

// NetworkFromName resolves a user-provided network identifier to the matching chain parameters.
func NetworkFromName(name string) (*chaincfg.Params, error) {
	normalized := strings.ToLower(strings.TrimSpace(name))
	if normalized == "" {
		return nil, errors.New("network name is required")
	}

	if params, ok := networkAliases[normalized]; ok {
		return params, nil
	}

	return nil, fmt.Errorf("unsupported network: %s", name)
}

// GenerateWIF creates a fresh private key for the provided network and encodes it as WIF.
func GenerateWIF(params *chaincfg.Params, compressed bool) (*btcutil.WIF, error) {
	if params == nil {
		return nil, errors.New("network parameters are required")
	}

	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	wif, err := btcutil.NewWIF(privKey, params, compressed)
	if err != nil {
		return nil, fmt.Errorf("failed to encode WIF: %v", err)
	}

	return wif, nil
}

// IsBIP38Format checks if the given string matches the BIP38 format pattern.
// Returns true if the string appears to be a valid BIP38 encrypted key.
func IsBIP38Format(key string) bool {
	return bip38Regex.MatchString(key)
}

// DecryptKey decrypts a BIP38 encrypted private key using the given passphrase.
// Returns the decrypted WIF (Wallet Import Format) private key or an error if decryption fails.
func DecryptKey(encryptedKey string, passphrase []byte) (*btcutil.WIF, error) {
	if !IsBIP38Format(encryptedKey) {
		return nil, errors.New("invalid BIP38 format")
	}

	// Decode the base58 key bytes first
	decoded := base58.Decode(encryptedKey)
	if len(decoded) != 43 {
		return nil, errors.New("invalid encrypted key length")
	}

	// Check magic byte matches expected value
	if decoded[0] != bip38Magic {
		return nil, errors.New("invalid magic byte")
	}

	// Extract checksum and verify quickly
	payload := decoded[:39]
	checksum := decoded[39:]
	hash := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash[:])

	if !constantTimeEqual(hash2[:4], checksum) {
		return nil, errors.New("invalid checksum")
	}

	// Check type and flags to avoid unsupported flow
	if decoded[1] != bip38Type {
		return nil, errors.New("unsupported BIP38 type")
	}

	// Extract compression flag bit
	compressed := false
	if decoded[2] == 0xe0 {
		compressed = true
	} else if decoded[2] != 0xc0 {
		return nil, errors.New("invalid flag byte")
	}

	// Extract address hash acting as salt
	addressHash := decoded[3:7]

	// Extract encrypted private key halves
	encryptedHalf1 := decoded[7:23]
	encryptedHalf2 := decoded[23:39]

	// Derive key using scrypt with standard cost
	derivedKey, err := scrypt.Key(passphrase, addressHash, 16384, 8, 8, 64)
	if err != nil {
		return nil, fmt.Errorf("scrypt derivation failed: %v", err)
	}

	// Split derived key in two halves
	derivedhalf1 := derivedKey[:32]
	derivedhalf2 := derivedKey[32:]

	// Decrypt the private key halves via AES
	block, err := aes.NewCipher(derivedhalf2)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	mode := newECBDecrypter(block)

	decrypted1 := make([]byte, 16)
	decrypted2 := make([]byte, 16)

	mode.CryptBlocks(decrypted1, encryptedHalf1)
	mode.CryptBlocks(decrypted2, encryptedHalf2)

	// XOR with first half of derived key to remove mask
	for i := 0; i < 16; i++ {
		decrypted1[i] ^= derivedhalf1[i]
		decrypted2[i] ^= derivedhalf1[i+16]
	}

	// Combine decrypted halves back into private key bytes
	privateKeyBytes := append(decrypted1, decrypted2...)

	// Verify by checking address hash matches salt
	privKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	if privKey.Key.IsZero() {
		return nil, errors.New("invalid private key scalar")
	}
	pubKey := privKey.PubKey()
	var pubKeyBytes []byte
	if compressed {
		pubKeyBytes = pubKey.SerializeCompressed()
	} else {
		pubKeyBytes = pubKey.SerializeUncompressed()
	}

	var matchedNet *chaincfg.Params
	for _, params := range supportedNetworks {
		addressPubKey, err := btcutil.NewAddressPubKey(pubKeyBytes, params)
		if err != nil {
			continue
		}

		address := addressPubKey.EncodeAddress()
		hash = sha256.Sum256([]byte(address))
		hash2 = sha256.Sum256(hash[:])

		if constantTimeEqual(hash2[:4], addressHash) {
			matchedNet = params
			break
		}
	}

	if matchedNet == nil {
		return nil, errors.New("incorrect passphrase")
	}

	wif, err := btcutil.NewWIF(privKey, matchedNet, compressed)
	if err != nil {
		return nil, fmt.Errorf("failed to create WIF: %v", err)
	}

	return wif, nil
}

// EncryptKey wrap the private key with BIP38 using passphrase bytes
func EncryptKey(wif *btcutil.WIF, passphrase []byte) (string, error) {
	privKeyBytes := wif.PrivKey.Serialize()
	compressed := wif.CompressPubKey

	// Compute address based on key format
	pubKey := wif.PrivKey.PubKey()
	var pubKeyBytes []byte
	if compressed {
		pubKeyBytes = pubKey.SerializeCompressed()
	} else {
		pubKeyBytes = pubKey.SerializeUncompressed()
	}

	netParams, err := NetworkFromWIF(wif)
	if err != nil {
		return "", err
	}

	addressPubKey, err := btcutil.NewAddressPubKey(pubKeyBytes, netParams)
	if err != nil {
		return "", fmt.Errorf("failed to create address: %v", err)
	}

	address := addressPubKey.EncodeAddress()

	// Create address hash acting like salt
	hash := sha256.Sum256([]byte(address))
	hash2 := sha256.Sum256(hash[:])
	addressHash := hash2[:4]

	// Derive key using scrypt same as decrypt side
	derivedKey, err := scrypt.Key(passphrase, addressHash, 16384, 8, 8, 64)
	if err != nil {
		return "", fmt.Errorf("scrypt derivation failed: %v", err)
	}

	// Split derived key halves for AES use
	derivedhalf1 := derivedKey[:32]
	derivedhalf2 := derivedKey[32:]

	// XOR private key bytes with first half mask
	xoredKey := make([]byte, 32)
	for i := 0; i < 32; i++ {
		xoredKey[i] = privKeyBytes[i] ^ derivedhalf1[i]
	}

	// Encrypt with AES ECB (library doesn't expose this mode directly)
	block, err := aes.NewCipher(derivedhalf2)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %v", err)
	}

	mode := newECBEncrypter(block)

	encrypted := make([]byte, 32)
	mode.CryptBlocks(encrypted[:16], xoredKey[:16])
	mode.CryptBlocks(encrypted[16:], xoredKey[16:])

	// Build the result byte slice with flags
	flagByte := byte(0xc0)
	if compressed {
		flagByte = 0xe0
	}

	result := []byte{bip38Magic, bip38Type, flagByte}
	result = append(result, addressHash...)
	result = append(result, encrypted...)

	// Add checksum as required by BIP38 spec
	hash = sha256.Sum256(result)
	hash2 = sha256.Sum256(hash[:])
	result = append(result, hash2[:4]...)

	return base58.Encode(result), nil
}

func constantTimeEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare(a, b) == 1
}

// NetworkFromWIF returns the Bitcoin network associated with the provided WIF.
func NetworkFromWIF(wif *btcutil.WIF) (*chaincfg.Params, error) {
	for _, params := range supportedNetworks {
		if wif.IsForNet(params) {
			return params, nil
		}
	}
	return nil, errors.New("unsupported WIF network")
}

// ECB mode helper because stdlib doesn't expose this cipher mode
type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

func newECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int { return x.blockSize }

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

func newECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func (x *ecbDecrypter) BlockSize() int { return x.blockSize }

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
