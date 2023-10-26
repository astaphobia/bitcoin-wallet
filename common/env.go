package common

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const envPath = "./_data/env.json"

// Env is env configuration structure
type Env struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	RpcUser     string `json:"rpc_user"`
	RpcPassword string `json:"rpc_password"`
	IsTestnet   bool   `json:"is_testnet"`
}

// LoadEnv to load env.json file from data directory
func LoadEnv() (env *Env, err error) {
	var (
		fullPath string
		file     []byte
	)

	fullPath, err = filepath.Abs(envPath)
	if err != nil {
		return
	}

	file, err = os.ReadFile(fullPath)
	if err != nil {
		return
	}

	env = &Env{}
	err = json.Unmarshal(file, &env)
	return
}
