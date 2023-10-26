package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"

	"astaphobia/bitcoin-wallet/common"
	"astaphobia/bitcoin-wallet/jrpc"
)

func main() {

	var (
		env            *common.Env
		descriptors    *common.ListDescriptors
		jr             *jrpc.JsonRpc
		dataToCsv      []*common.Derivation
		descriptorsMap map[string]*common.Descriptor
		err            error
	)

	env, err = common.LoadEnv()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	descriptors, err = common.LoadDescriptors()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	descriptorsMap = descriptors.ToMap()
	jr = &jrpc.JsonRpc{
		Host:     env.Host,
		Port:     env.Port,
		User:     env.RpcUser,
		Password: env.RpcPassword,
	}
	addresses, err := jr.ListReceivedByAddress(0, true)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	for index, address := range addresses {
		var (
			info           *jrpc.AddressInfo
			descriptor     *common.Descriptor
			key            *hdkeychain.ExtendedKey
			cfg            = chaincfg.MainNetParams
			derivation     *common.Derivation
			derivationInfo *common.DerivationInfo
			btcAddress     btcutil.Address
			err            error
			ok             bool
		)

		if env.IsTestnet {
			cfg = chaincfg.TestNet3Params
		}

		info, err = jr.GetAddressInfo(address.Address)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Printf("----: %s :---- %d\n", address.Address, index)
		b, _ := json.Marshal(info)
		fmt.Printf("info: \n%s\n", b)

		descriptor, ok = descriptorsMap[info.ParentDesc]
		if !ok {
			fmt.Printf("%s was not found\n", info.ParentDesc)
			os.Exit(0)
		}
		derivationInfo = descriptor.GetDerivationInfo()

		key, err = common.GetKey(derivationInfo.RootKey, info.HDKeyPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		privateKey, err := key.ECPrivKey()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		publicKey, err := key.ECPubKey()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		wif, err := btcutil.NewWIF(privateKey, &cfg, true)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		derivation = &common.Derivation{
			ParentKey: derivationInfo.RootKey,
			WIF:       wif.String(),
			PublicKey: hex.EncodeToString(publicKey.SerializeCompressed()),
			Type:      derivationInfo.Type,
			Path:      info.HDKeyPath,
		}

		btcAddress, err = common.GetAddress(key, &cfg, derivationInfo)
		if err != nil {
			fmt.Println(err)
			return
		}
		derivation.Address = btcAddress.EncodeAddress()

		b, _ = json.Marshal(derivation)
		fmt.Printf("result: \n%s\n", b)
		fmt.Println("------")

		dataToCsv = append(dataToCsv, derivation)
	}

	err = common.ToCsv(dataToCsv)
	if err != nil {
		fmt.Println(err)
	}
}
