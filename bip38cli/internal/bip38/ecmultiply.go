package bip38

import (
	"errors"

	"github.com/btcsuite/btcd/btcutil"
)

func decryptECMultiply(_ []byte, _ []byte) (*btcutil.WIF, error) {
	return nil, errors.New("EC-multiply decryption is not implemented")
}
