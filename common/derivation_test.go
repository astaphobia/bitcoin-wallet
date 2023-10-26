package common

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
)

func TestGetKey(t *testing.T) {
	type args struct {
		rootKey        string
		derivationPath string
	}
	tests := []struct {
		name       string
		args       args
		wantExtKey *hdkeychain.ExtendedKey
		wantErr    bool
	}{
		{
			name: "Success",
			args: args{
				rootKey:        "xprv9s21ZrQH143K2CBBPTGvuJgd9xJ5k9232EPYhX7HvEv6wFsGPa4c6eGdbZzYd5PiqdKa8k1HEH8oMiqJ2scyvGwzKwwDjzcrTMfGXXpPnRc",
				derivationPath: "m/84'/0'/0'/112",
			},
			wantErr:    false,
			wantExtKey: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExtKey, err := GetKey(tt.args.rootKey, tt.args.derivationPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			cfg := chaincfg.MainNetParams
			priv, _ := gotExtKey.ECPrivKey()
			wif, _ := btcutil.NewWIF(priv, &cfg, true)
			derivInfo := &DerivationInfo{
				RootKey: tt.args.rootKey,
				Type:    "wpkh",
			}
			add, _ := GetAddress(gotExtKey, &cfg, derivInfo)

			w, _ := btcutil.DecodeWIF("L3j2VaDGaHhst9tHXfagLg1MorQbV5FRiYyt5HrxViGMXv65itvY")
			fmt.Println(w.String())

			fmt.Printf("address: %s, WIF: %s", add.EncodeAddress(), wif.String())
			if !reflect.DeepEqual(gotExtKey, tt.wantExtKey) {
				t.Errorf("GetKey() gotExtKey = %v, want %v", gotExtKey, tt.wantExtKey)
			}
		})
	}
}
