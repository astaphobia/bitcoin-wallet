package jrpc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AddressInfo is json JSONRPC response format for getaddressinfo rpc method
type AddressInfo struct {
	Timestamp           uint64   `json:"timestamp"`
	WitnessVersion      uint     `json:"witness_version"`
	Address             string   `json:"address"`
	ScriptPubKey        string   `json:"scriptPubKey"`
	Desc                string   `json:"desc"`
	ParentDesc          string   `json:"parent_desc"`
	HDKeyPath           string   `json:"hdkeypath"`
	HDSeedID            string   `json:"hdseedid"`
	HDMasterFingerprint string   `json:"hdmasterfingerprint"`
	WitnessProgram      string   `json:"witness_program"`
	PubKey              string   `json:"pubkey"`
	Labels              []string `json:"labels"`
	IsMine              bool     `json:"ismine"`
	Solvable            bool     `json:"solvable"`
	IsChange            bool     `json:"ischange"`
	IsWatchOnly         bool     `json:"iswatchonly"`
	IsScript            bool     `json:"isscript"`
	IsWitness           bool     `json:"iswitness"`
}

// GetAddressInfo to call JSONRPC method getaddressinfo to running node.
// returns *AddressInfo and error
func (jr *JsonRpc) GetAddressInfo(address string) (resp *AddressInfo, err error) {
	var (
		resBytes []byte
		code     int
		tag      = "JsonRpc.GetAddressInfo"
		r        = struct {
			Result *AddressInfo `json:"result"`
			Error  string       `json:"error"`
		}{}
	)

	resBytes, code, err = jr.Do("getaddressinfo", []interface{}{address})
	if err != nil {
		fmt.Printf("%s - err: %s, code: %d", tag, err.Error(), code)
		return nil, err
	}

	err = json.Unmarshal(resBytes, &r)
	if err != nil {
		fmt.Printf("%s - err: %s, code: %d", tag, err.Error(), code)
		return nil, err
	}

	if code != http.StatusOK {
		fmt.Printf("%s - err: %s, code: %d", tag, r.Error, code)
		return nil, fmt.Errorf("%s: %s", tag, r.Error)
	}

	resp = r.Result
	return
}
