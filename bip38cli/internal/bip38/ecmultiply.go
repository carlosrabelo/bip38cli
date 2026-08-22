package bip38

// ECMultiply implements BIP38 Section "Encryption when EC multiply flag used".
//
// The owner generates an intermediate passphrase code (via GenerateIntermediateCode)
// and hands it to a third party (the "printer"). The printer calls ECMultiplyEncrypt
// to produce an encrypted private key without ever knowing the passphrase.
// The owner can later decrypt the key with DecryptKey / decryptECMultiply using their passphrase.

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/chaincfg"
	"golang.org/x/crypto/scrypt"
)

// ECMultiplyResult holds the output of ECMultiplyEncrypt.
type ECMultiplyResult struct {
	EncryptedKey     string // BIP38 encrypted key (6P...)
	ConfirmationCode string // confirmation code (cfrm38...)
	Compressed       bool
}

// ConfirmationResult is returned by VerifyConfirmationCode.
type ConfirmationResult struct {
	Address   string
	HasLotSeq bool
	LotNumber *uint32
	SeqNumber *uint32
}

// ECMultiplyEncrypt generates a BIP38 EC-multiply encrypted key from an
// intermediate passphrase code. The caller does not need to know the passphrase.
func ECMultiplyEncrypt(intermediateCode string, compressed bool) (*ECMultiplyResult, error) {
	ic, err := ParseIntermediateCode(intermediateCode)
	if err != nil {
		return nil, fmt.Errorf("invalid intermediate code: %w", err)
	}

	seedb := make([]byte, 24)
	if _, readErr := rand.Read(seedb); readErr != nil {
		return nil, fmt.Errorf("failed to generate seedb: %w", err)
	}
	defer zeroBytes(seedb)

	h1 := sha256.Sum256(seedb)
	h2 := sha256.Sum256(h1[:])
	factorb := append([]byte{}, h2[:]...)
	defer zeroBytes(factorb)

	passPoint, err := btcec.ParsePubKey(ic.PassPoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse passpoint: %w", err)
	}

	var factorbScalar btcec.ModNScalar
	if overflow := factorbScalar.SetByteSlice(factorb); overflow {
		return nil, errors.New("invalid factorb scalar")
	}
	if factorbScalar.IsZero() {
		return nil, errors.New("invalid factorb scalar")
	}

	var generatedJacobian btcec.JacobianPoint
	passPoint.AsJacobian(&generatedJacobian)
	btcec.ScalarMultNonConst(&factorbScalar, &generatedJacobian, &generatedJacobian)
	generatedJacobian.ToAffine()
	generatedPubKey := btcec.NewPublicKey(&generatedJacobian.X, &generatedJacobian.Y)

	var pubKeyBytes []byte
	if compressed {
		pubKeyBytes = generatedPubKey.SerializeCompressed()
	} else {
		pubKeyBytes = generatedPubKey.SerializeUncompressed()
	}

	netParams := &chaincfg.MainNetParams
	addrPubKey, err := btcutil.NewAddressPubKey(pubKeyBytes, netParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}
	address := addrPubKey.EncodeAddress()

	addrHash1 := sha256.Sum256([]byte(address))
	addrHash2 := sha256.Sum256(addrHash1[:])
	addressHash := addrHash2[:4]

	ownerEntropy := make([]byte, 8)
	copy(ownerEntropy, ic.OwnerEntropy)

	salt := make([]byte, 0, 12)
	salt = append(salt, addressHash...)
	salt = append(salt, ownerEntropy...)
	derived, err := scrypt.Key(ic.PassPoint, salt, 1024, 1, 1, 64)
	if err != nil {
		return nil, fmt.Errorf("scrypt derivation failed: %w", err)
	}
	defer zeroBytes(derived)

	derivedhalf1 := derived[:32]
	derivedhalf2 := derived[32:]

	block, err := aes.NewCipher(derivedhalf2)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	xored1 := make([]byte, 16)
	for i := 0; i < 16; i++ {
		xored1[i] = seedb[i] ^ derivedhalf1[i]
	}
	encryptedhalf1 := make([]byte, 16)
	block.Encrypt(encryptedhalf1, xored1)

	combined := make([]byte, 16)
	copy(combined, encryptedhalf1[8:])
	copy(combined[8:], seedb[16:])
	xored2 := make([]byte, 16)
	for i := 0; i < 16; i++ {
		xored2[i] = combined[i] ^ derivedhalf1[16+i]
	}
	encryptedhalf2 := make([]byte, 16)
	block.Encrypt(encryptedhalf2, xored2)

	flagbyte := byte(0x00)
	if compressed {
		flagbyte |= 0x20
	}
	if ic.HasLotSeq {
		flagbyte |= 0x04
	}

	payload := []byte{bip38Magic, bip38TypeEC, flagbyte}
	payload = append(payload, addressHash...)
	payload = append(payload, ownerEntropy...)
	payload = append(payload, encryptedhalf1[:8]...)
	payload = append(payload, encryptedhalf2...)

	cs1 := sha256.Sum256(payload)
	cs2 := sha256.Sum256(cs1[:])
	payload = append(payload, cs2[:4]...)

	encryptedKey := base58.Encode(payload)

	confirmCode, err := buildConfirmationCode(
		flagbyte, addressHash, ownerEntropy,
		factorb, derivedhalf1, derivedhalf2,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build confirmation code: %w", err)
	}

	return &ECMultiplyResult{
		EncryptedKey:     encryptedKey,
		ConfirmationCode: confirmCode,
		Compressed:       compressed,
	}, nil
}

// decryptECMultiply decrypts a BIP38 EC-multiply encrypted private key.
func decryptECMultiply(decoded []byte, passphrase []byte) (*btcutil.WIF, error) { //nolint:gocyclo
	flagbyte := decoded[2]
	hasLotSeq := flagbyte&0x04 != 0
	compressed := flagbyte&0x20 != 0

	if flagbyte&^byte(0x24) != 0 {
		return nil, errors.New("invalid flag byte")
	}

	addressHash := decoded[3:7]
	ownerEntropy := decoded[7:15]
	encryptedPart1Prefix := decoded[15:23]
	encryptedPart2 := decoded[23:39]

	ownersalt := ownerEntropy
	if hasLotSeq {
		ownersalt = ownerEntropy[:4]
	}

	passfactor, err := derivePassfactor(passphrase, ownersalt, ownerEntropy, hasLotSeq)
	if err != nil {
		return nil, err
	}
	defer zeroBytes(passfactor)

	privFromPass, _ := btcec.PrivKeyFromBytes(passfactor)
	if privFromPass.Key.IsZero() {
		return nil, errors.New("invalid passfactor scalar")
	}
	passpoint := privFromPass.PubKey().SerializeCompressed()

	salt := make([]byte, 0, 12)
	salt = append(salt, addressHash...)
	salt = append(salt, ownerEntropy...)
	derived, err := scrypt.Key(passpoint, salt, 1024, 1, 1, 64)
	if err != nil {
		return nil, fmt.Errorf("scrypt derivation failed: %w", err)
	}
	defer zeroBytes(derived)

	derivedhalf1 := derived[:32]
	derivedhalf2 := derived[32:]

	block, err := aes.NewCipher(derivedhalf2)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	decryptedPart2 := make([]byte, 16)
	block.Decrypt(decryptedPart2, encryptedPart2)
	tmp2 := make([]byte, 16)
	for i := 0; i < 16; i++ {
		tmp2[i] = decryptedPart2[i] ^ derivedhalf1[16+i]
	}

	seedb := make([]byte, 24)
	copy(seedb[16:], tmp2[8:])

	encryptedPart1 := make([]byte, 16)
	copy(encryptedPart1, encryptedPart1Prefix)
	copy(encryptedPart1[8:], tmp2[:8])

	decryptedPart1 := make([]byte, 16)
	block.Decrypt(decryptedPart1, encryptedPart1)
	for i := 0; i < 16; i++ {
		seedb[i] = decryptedPart1[i] ^ derivedhalf1[i]
	}
	defer zeroBytes(seedb)

	fh1 := sha256.Sum256(seedb)
	fh2 := sha256.Sum256(fh1[:])
	factorb := append([]byte{}, fh2[:]...)
	defer zeroBytes(factorb)

	var passfactorScalar, factorbScalar, privScalar btcec.ModNScalar
	if passfactorScalar.SetByteSlice(passfactor) {
		return nil, errors.New("invalid passfactor scalar")
	}
	if factorbScalar.SetByteSlice(factorb) {
		return nil, errors.New("invalid factorb scalar")
	}
	privScalar.Set(&passfactorScalar)
	privScalar.Mul(&factorbScalar)
	if privScalar.IsZero() {
		return nil, errors.New("invalid private key scalar")
	}

	privKeyBytes := privScalar.Bytes()
	privKeyBytesSlice := privKeyBytes[:]
	defer zeroBytes(privKeyBytesSlice)

	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytesSlice)
	pubKey := privKey.PubKey()
	var pubKeyBytes []byte
	if compressed {
		pubKeyBytes = pubKey.SerializeCompressed()
	} else {
		pubKeyBytes = pubKey.SerializeUncompressed()
	}

	var matchedNet *chaincfg.Params
	for _, params := range supportedNetworks {
		addressPubKey, addrErr := btcutil.NewAddressPubKey(pubKeyBytes, params)
		if addrErr != nil {
			continue
		}
		address := addressPubKey.EncodeAddress()
		h1 := sha256.Sum256([]byte(address))
		h2 := sha256.Sum256(h1[:])
		if constantTimeEqual(h2[:4], addressHash) {
			matchedNet = params
			break
		}
	}
	if matchedNet == nil {
		return nil, errors.New("incorrect passphrase")
	}

	wif, err := btcutil.NewWIF(privKey, matchedNet, compressed)
	if err != nil {
		return nil, fmt.Errorf("failed to create WIF: %w", err)
	}
	return wif, nil
}

// VerifyConfirmationCode checks that a confirmation code depends on the passphrase.
func VerifyConfirmationCode(confirmationCode string, passphrase []byte) (*ConfirmationResult, error) { //nolint:gocyclo
	passphrase = normalizePassphrase(passphrase)

	decoded := base58.Decode(confirmationCode)
	if len(decoded) != 55 {
		return nil, errors.New("invalid confirmation code length")
	}

	expectedMagic := []byte{0x64, 0x3B, 0xF6, 0xA8, 0x9A}
	if !constantTimeEqual(decoded[:5], expectedMagic) {
		return nil, errors.New("invalid confirmation code magic")
	}

	payload := decoded[:51]
	checksum := decoded[51:]
	cs1 := sha256.Sum256(payload)
	cs2 := sha256.Sum256(cs1[:])
	if !constantTimeEqual(cs2[:4], checksum) {
		return nil, errors.New("invalid checksum")
	}

	flagbyte := decoded[5]
	hasLotSeq := flagbyte&0x04 != 0
	compressed := flagbyte&0x20 != 0
	addressHash := decoded[6:10]
	ownerEntropy := decoded[10:18]
	pointbprefix := decoded[18]
	encpointb1 := decoded[19:35]
	encpointb2 := decoded[35:51]

	ownersalt := ownerEntropy
	if hasLotSeq {
		ownersalt = ownerEntropy[:4]
	}

	passfactor, err := derivePassfactor(passphrase, ownersalt, ownerEntropy, hasLotSeq)
	if err != nil {
		return nil, err
	}
	defer zeroBytes(passfactor)

	privFromPass, _ := btcec.PrivKeyFromBytes(passfactor)
	if privFromPass.Key.IsZero() {
		return nil, errors.New("invalid passfactor scalar")
	}

	salt := make([]byte, 0, 12)
	salt = append(salt, addressHash...)
	salt = append(salt, ownerEntropy...)
	passpoint := privFromPass.PubKey().SerializeCompressed()
	derived, err := scrypt.Key(passpoint, salt, 1024, 1, 1, 64)
	if err != nil {
		return nil, fmt.Errorf("scrypt derivation failed: %w", err)
	}
	defer zeroBytes(derived)

	derivedhalf1 := derived[:32]
	derivedhalf2 := derived[32:]

	block, err := aes.NewCipher(derivedhalf2)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	dec1 := make([]byte, 16)
	block.Decrypt(dec1, encpointb1)
	dec2 := make([]byte, 16)
	block.Decrypt(dec2, encpointb2)

	pointb := make([]byte, 33)
	pointb[0] = pointbprefix ^ (derivedhalf2[31] & 0x01)
	for i := 0; i < 16; i++ {
		pointb[1+i] = dec1[i] ^ derivedhalf1[i]
		pointb[17+i] = dec2[i] ^ derivedhalf1[16+i]
	}

	pointbPub, err := btcec.ParsePubKey(pointb)
	if err != nil {
		return nil, errors.New("incorrect passphrase")
	}

	var passfactorScalar btcec.ModNScalar
	if passfactorScalar.SetByteSlice(passfactor) {
		return nil, errors.New("invalid passfactor scalar")
	}

	var pointJacobian, resultJacobian btcec.JacobianPoint
	pointbPub.AsJacobian(&pointJacobian)
	btcec.ScalarMultNonConst(&passfactorScalar, &pointJacobian, &resultJacobian)
	resultJacobian.ToAffine()
	addressPub := btcec.NewPublicKey(&resultJacobian.X, &resultJacobian.Y)

	var pubKeyBytes []byte
	if compressed {
		pubKeyBytes = addressPub.SerializeCompressed()
	} else {
		pubKeyBytes = addressPub.SerializeUncompressed()
	}

	addrPubKey, err := btcutil.NewAddressPubKey(pubKeyBytes, &chaincfg.MainNetParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}
	address := addrPubKey.EncodeAddress()

	h1 := sha256.Sum256([]byte(address))
	h2 := sha256.Sum256(h1[:])
	if !constantTimeEqual(h2[:4], addressHash) {
		return nil, errors.New("incorrect passphrase")
	}

	result := &ConfirmationResult{
		Address:   address,
		HasLotSeq: hasLotSeq,
	}
	if hasLotSeq {
		lotSeqValue := uint32(ownerEntropy[4])<<24 | uint32(ownerEntropy[5])<<16 |
			uint32(ownerEntropy[6])<<8 | uint32(ownerEntropy[7])
		lot := lotSeqValue / 4096
		seq := lotSeqValue % 4096
		result.LotNumber = &lot
		result.SeqNumber = &seq
	}
	return result, nil
}

func derivePassfactor(passphrase, ownersalt, ownerEntropy []byte, hasLotSeq bool) ([]byte, error) {
	prefactor, err := scrypt.Key(passphrase, ownersalt, 16384, 8, 8, 32)
	if err != nil {
		return nil, fmt.Errorf("scrypt derivation failed: %w", err)
	}
	defer zeroBytes(prefactor)

	if !hasLotSeq {
		out := make([]byte, len(prefactor))
		copy(out, prefactor)
		return out, nil
	}

	combined := append(append([]byte{}, prefactor...), ownerEntropy...)
	hash := sha256.Sum256(combined)
	hash2 := sha256.Sum256(hash[:])
	zeroBytes(combined)
	out := make([]byte, 32)
	copy(out, hash2[:])
	return out, nil
}

func buildConfirmationCode(
	flagbyte byte,
	addressHash, ownerEntropy, factorb, derivedhalf1, derivedhalf2 []byte,
) (string, error) {
	privKey, _ := btcec.PrivKeyFromBytes(factorb)
	if privKey.Key.IsZero() {
		return "", errors.New("invalid factorb for confirmation code")
	}
	pointb := privKey.PubKey().SerializeCompressed()

	block, err := aes.NewCipher(derivedhalf2)
	if err != nil {
		return "", fmt.Errorf("AES cipher failed: %w", err)
	}

	xored1 := make([]byte, 16)
	for i := 0; i < 16; i++ {
		xored1[i] = pointb[1+i] ^ derivedhalf1[i]
	}
	encpointb1 := make([]byte, 16)
	block.Encrypt(encpointb1, xored1)

	xored2 := make([]byte, 16)
	for i := 0; i < 16; i++ {
		xored2[i] = pointb[17+i] ^ derivedhalf1[16+i]
	}
	encpointb2 := make([]byte, 16)
	block.Encrypt(encpointb2, xored2)

	prefixByte := pointb[0] ^ (derivedhalf2[31] & 0x01)

	cfrmPayload := []byte{0x64, 0x3B, 0xF6, 0xA8, 0x9A, flagbyte}
	cfrmPayload = append(cfrmPayload, addressHash...)
	cfrmPayload = append(cfrmPayload, ownerEntropy...)
	cfrmPayload = append(cfrmPayload, prefixByte)
	cfrmPayload = append(cfrmPayload, encpointb1...)
	cfrmPayload = append(cfrmPayload, encpointb2...)

	cs1 := sha256.Sum256(cfrmPayload)
	cs2 := sha256.Sum256(cs1[:])
	cfrmPayload = append(cfrmPayload, cs2[:4]...)

	return base58.Encode(cfrmPayload), nil
}
