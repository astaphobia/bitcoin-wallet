package jrpc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ReceivedAddress is json JSONRPC response format for listreceivedbyaddress rpc method
type ReceivedAddress struct {
	Account           string   `json:"account"`
	Address           string   `json:"address"`
	Amount            float64  `json:"amount"`
	Confirmations     uint64   `json:"confirmations"`
	TxIDs             []string `json:"txids,omitempty"`
	InvolvesWatchonly bool     `json:"involvesWatchonly,omitempty"`
}

// ListReceivedByAddress to call JSONRPC method listreceivedbyaddress to running node.
// returns list of *ReceivedAddress and error
func (jr *JsonRpc) ListReceivedByAddress(minConfirms int, includeEmpty bool) (resp []*ReceivedAddress, err error) {

	var (
		resBytes []byte
		code     int
		tag      = "JsonRpc.ListReceivedByAddress"
		r        = struct {
			Result []*ReceivedAddress `json:"result"`
			Error  string             `json:"error"`
		}{}
	)

	resBytes, code, err = jr.Do("listreceivedbyaddress", []interface{}{minConfirms, includeEmpty})
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
