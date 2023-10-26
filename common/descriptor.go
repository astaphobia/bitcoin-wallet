package common

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
)

const descriptorsPath = "./_data/descriptors.json"
const descriptorPubsPath = "./_data/descriptor_pubs.json"

// Descriptor is a compact and semi-standard way to easily encode, or "describe", how scripts (and subsequently, addresses) of a wallet should be generated
type Descriptor struct {
	Timestamp      uint64 `json:"timestamp"`
	Next           uint   `json:"next"`
	NextIndex      uint   `json:"next_index"`
	Range          []int  `json:"range"`
	Desc           string `json:"desc"`
	DescPub        string
	Active         bool `json:"active"`
	Internal       bool `json:"internal"`
	derivationInfo *DerivationInfo
}

// DerivationInfo is split of parent descriptor string
type DerivationInfo struct {
	RootKey     string
	Type        string // wpkh,pkh
	Fingerprint string
}

// GetDerivationInfo to fill derivationInfo of Descriptor
func (d *Descriptor) GetDerivationInfo() *DerivationInfo {
	if d.derivationInfo != nil {
		return d.derivationInfo
	}

	var (
		rgx = regexp.MustCompile(`[a-z]+[a-zA-Z0-9]+`)
		res []string
	)

	res = rgx.FindAllString(d.Desc, -1)
	d.derivationInfo = &DerivationInfo{
		Type:        res[0],
		RootKey:     res[1],
		Fingerprint: res[2],
	}

	return d.derivationInfo
}

// ListDescriptors is Format of list descriptors that was received from bitcoin-cli
type ListDescriptors struct {
	WalletName  string        `json:"wallet_name"`
	Descriptors []*Descriptor `json:"descriptors"`
}

// ToMap to build a map version of list descriptors and desc pub as the key of it.
func (ld *ListDescriptors) ToMap() (m map[string]*Descriptor) {

	m = map[string]*Descriptor{}
	if ld == nil && ld.Descriptors == nil {
		return
	}

	for _, d := range ld.Descriptors {
		m[d.DescPub] = d
	}
	return
}

// LoadDescriptors to load descriptors.json file from data directory
func LoadDescriptors() (descriptors *ListDescriptors, err error) {
	var (
		fullPath string
		file     []byte
		pubs     *ListDescriptors
	)

	fullPath, err = filepath.Abs(descriptorsPath)
	if err != nil {
		return
	}

	file, err = os.ReadFile(fullPath)
	if err != nil {
		return
	}

	descriptors = &ListDescriptors{}
	err = json.Unmarshal(file, &descriptors)
	if err != nil {
		return
	}

	fullPath, err = filepath.Abs(descriptorPubsPath)
	if err != nil {
		return
	}

	file, err = os.ReadFile(fullPath)
	if err != nil {
		return
	}

	err = json.Unmarshal(file, &pubs)
	if err != nil {
		return
	}

	for k, desc := range descriptors.Descriptors {
		pub := pubs.Descriptors[k]
		desc.DescPub = pub.Desc
	}
	return
}
