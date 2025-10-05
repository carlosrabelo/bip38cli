package bip38

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"fmt"
	"regexp"

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
)

// EncryptedKey represent a BIP38 encrypted private key record
type EncryptedKey struct {
	Key        string
	Compressed bool
	ECMultiply bool
}

// IsBIP38Format check if given string look like BIP38 format
func IsBIP38Format(key string) bool {
	return bip38Regex.MatchString(key)
}

// DecryptKey unwrap one BIP38 encrypted private key using given passphrase bytes
func DecryptKey(encryptedKey string, passphrase []byte) (*btcutil.WIF, error) {
	if !IsBIP38Format(encryptedKey) {
		return nil, errors.New("invalid BIP38 format")
	}

	// Decode the base58 key bytes first
	decoded := base58.Decode(encryptedKey)
	if len(decoded) != 43 {
		return nil, errors.New("invalid encrypted key length")
	}

	// Check magic byte stay expected value
	if decoded[0] != bip38Magic {
		return nil, errors.New("invalid magic byte")
	}

	// Extract checksum and verify quickly
	payload := decoded[:39]
	checksum := decoded[39:]
	hash := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash[:])

	if !bytesEqual(hash2[:4], checksum) {
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

	// Create private key object from bytes
	privKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)

	// Create WIF wrapper for private key
	wif, err := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, compressed)
	if err != nil {
		return nil, fmt.Errorf("failed to create WIF: %v", err)
	}

	// Verify by checking address hash matches salt
	pubKey := privKey.PubKey()
	var pubKeyBytes []byte
	if compressed {
		pubKeyBytes = pubKey.SerializeCompressed()
	} else {
		pubKeyBytes = pubKey.SerializeUncompressed()
	}

	// Create address and verify checksum again
	addressPubKey, err := btcutil.NewAddressPubKey(pubKeyBytes, &chaincfg.MainNetParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create address: %v", err)
	}

	address := addressPubKey.EncodeAddress()
	hash = sha256.Sum256([]byte(address))
	hash2 = sha256.Sum256(hash[:])

	if !bytesEqual(hash2[:4], addressHash) {
		return nil, errors.New("incorrect passphrase")
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

	addressPubKey, err := btcutil.NewAddressPubKey(pubKeyBytes, &chaincfg.MainNetParams)
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

	// Encrypt with AES ECB even if library no offer directly
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

	// Add checksum just like BIP38 spec request
	hash = sha256.Sum256(result)
	hash2 = sha256.Sum256(hash[:])
	result = append(result, hash2[:4]...)

	return base58.Encode(result), nil
}

// bytesEqual compare two byte slices for equality without timing tricks
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ECB mode helper because stdlib no expose this cipher mode
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
