package cli

import (
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/spf13/cobra"

	"github.com/carlosrabelo/bip38cli/core/internal/domain/bip38"
)

type addressType string

const (
	addressTypeBIP84 addressType = "bip84"
	addressTypeBIP44 addressType = "bip44"
)

func parseAddressType(value string) (addressType, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" || normalized == string(addressTypeBIP84) {
		return addressTypeBIP84, nil
	}
	if normalized == string(addressTypeBIP44) {
		return addressTypeBIP44, nil
	}
	return "", fmt.Errorf("unsupported address type: %s", value)
}

func isVerbose(cmd *cobra.Command) bool {
	if cmd == nil {
		return false
	}
	if cmd.Flag("verbose") != nil && cmd.Flag("verbose").Changed {
		return true
	}
	return false
}

func outputFormat(cmd *cobra.Command) string {
	if cmd == nil {
		return "text"
	}
	if cmd.Flag("output-format") != nil && cmd.Flag("output-format").Changed {
		return cmd.Flag("output-format").Value.String()
	}
	return "text"
}

func isCompressed(cmd *cobra.Command) bool {
	if cmd == nil {
		return true // default to compressed
	}
	if cmd.Flag("compressed") != nil && cmd.Flag("compressed").Changed {
		val, _ := cmd.Flags().GetBool("compressed")
		return val
	}
	return true // default to compressed
}

func addressForWIF(wif *btcutil.WIF, mode addressType) (string, error) {
	netParams, err := bip38.NetworkFromWIF(wif)
	if err != nil {
		return "", err
	}

	pubKey := wif.PrivKey.PubKey()
	if mode == addressTypeBIP84 && !wif.CompressPubKey {
		mode = addressTypeBIP44
	}

	switch mode {
	case addressTypeBIP84:
		compressed := pubKey.SerializeCompressed()
		hash := btcutil.Hash160(compressed)
		addr, err := btcutil.NewAddressWitnessPubKeyHash(hash, netParams)
		if err != nil {
			return "", err
		}
		return addr.EncodeAddress(), nil
	case addressTypeBIP44:
		var serialized []byte
		if wif.CompressPubKey {
			serialized = pubKey.SerializeCompressed()
		} else {
			serialized = pubKey.SerializeUncompressed()
		}
		hash := btcutil.Hash160(serialized)
		addr, err := btcutil.NewAddressPubKeyHash(hash, netParams)
		if err != nil {
			return "", err
		}
		return addr.EncodeAddress(), nil
	default:
		return "", fmt.Errorf("unsupported address type: %s", mode)
	}
}
