package bip38

import "golang.org/x/text/unicode/norm"

// normalizePassphrase applies Unicode NFC as required by BIP38 before scrypt.
func normalizePassphrase(passphrase []byte) []byte {
	if len(passphrase) == 0 {
		return passphrase
	}
	return norm.NFC.Bytes(passphrase)
}

// zeroBytes overwrites a buffer so sensitive material does not linger.
func zeroBytes(buf []byte) {
	for i := range buf {
		buf[i] = 0
	}
}
