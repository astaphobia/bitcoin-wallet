package common

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
)

// Derivation is structure of final result
type Derivation struct {
	Address   string
	Type      string
	Path      string
	WIF       string
	PublicKey string
	ParentKey string
}

// GetKey to derive a new child from given root key and derivation path
func GetKey(rootKey string, derivationPath string) (extKey *hdkeychain.ExtendedKey, err error) {
	var (
		key   *hdkeychain.ExtendedKey
		rgx   = regexp.MustCompile(`[0-9']+`)
		paths = rgx.FindAllString(derivationPath, -1)
	)

	key, err = hdkeychain.NewKeyFromString(rootKey)
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		var n uint32
		n, err = getNumber(path)
		if err != nil {
			return nil, err
		}

		key, err = key.Derive(n)
		if err != nil {
			return nil, err
		}
	}

	return key, nil

}

func getNumber(path string) (uint32, error) {
	if strings.HasSuffix(path, "'") {
		s := strings.Replace(path, "'", "", 1)
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return hdkeychain.HardenedKeyStart + uint32(n), nil
	}
	n, err := strconv.ParseUint(path, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint32(n), nil
}

// GetAddress to retrieve btcutil.Address from given key, chain, and derivation info.
// For now only support btcutil.AddressPubKeyHash or PKH and btcutil.AddressWithnessPubKeyHash or WPKH
func GetAddress(
	key *hdkeychain.ExtendedKey,
	cfg *chaincfg.Params,
	derivationInfo *DerivationInfo,
) (pubKeyHash btcutil.Address, err error) {

	if 0 == strings.Compare(derivationInfo.Type, "pkh") {
		pubKeyHash, err = key.Address(cfg)
		return
	}

	pubKey, err := key.ECPubKey()
	if err != nil {
		return nil, err
	}

	hash := btcutil.Hash160(pubKey.SerializeCompressed())
	pubKeyHash, err = btcutil.NewAddressWitnessPubKeyHash(hash, cfg)

	return
}

// ToCsv to store data to a CSV file
func ToCsv(derivations []*Derivation) error {

	var (
		fullPath, err = filepath.Abs("./_data/derivation.csv")
		file          *os.File
		data          [][]string
	)

	if err != nil {
		return err
	}

	file, err = os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	data = append(data, []string{
		"Address",
		"Type",
		"Path",
		"PublicKey",
		"WIF",
		"ParentKey",
	})
	for _, derivation := range derivations {
		data = append(data, []string{
			derivation.Address,
			derivation.Type,
			derivation.Path,
			derivation.PublicKey,
			derivation.WIF,
			derivation.ParentKey,
		})
	}

	w := csv.NewWriter(file)
	return w.WriteAll(data)
}
