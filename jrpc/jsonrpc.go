package jrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Request JSONRPC body structure
type Request struct {
	JsonRpc string        `json:"jsonrpc,omitempty"`
	ID      string        `json:"id,omitempty"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// JsonRpc is fields that required to access node via JSONRPC
type JsonRpc struct {
	User     string
	Password string
	Host     string
	Port     string
}

// Do build a JSONRPC request and call POST method with a basic auth to running btc node with given RPC method and RPC params
// return response body in bytes, http code response and error
func (jr *JsonRpc) Do(method string, params []interface{}) (resBytes []byte, code int, err error) {

	var (
		client   = http.DefaultClient
		req      *http.Request
		res      *http.Response
		data     []byte
		endpoint = fmt.Sprintf("http://%s:%s", jr.Host, jr.Port)
		tag      = "JsonRpc.Do"
	)

	client.Timeout = 5 * time.Second

	data, err = json.Marshal(&Request{
		Method: method,
		Params: params,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("%s: json.Marshal - %w", tag, err)
	}

	req, err = http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, 0, fmt.Errorf("%s: http.NewRequestWithContext - %w", tag, err)
	}

	req.Header.Add("Content-Type", "text/plain")
	req.SetBasicAuth(jr.User, jr.Password)

	res, err = client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: client.Do - %w", tag, err)
	}
	defer res.Body.Close()

	resBytes, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: io.ReadAll - %w", tag, err)
	}

	code = res.StatusCode
	return
}
