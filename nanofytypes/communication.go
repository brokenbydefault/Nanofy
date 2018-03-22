package nanofytypes

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Block"
)

// CLIENT
type DefaultRequest struct {
	Action string `json:"action"`
	App    string `json:"app"`
}

type RequestByFile struct {
	FileKey Wallet.PublicKey `json:"filekey"`
	PubKey  Wallet.PublicKey `json:"pubkey"`
	DefaultRequest
}

type RequestByBlock struct {
	FlagBlock Block.BlockHash
	DefaultRequest
}

// SERVER
type Response struct {
	Exist    bool
	Error    string
	PubKey   Wallet.PublicKey `json:"pubkey"`
	FlagHash Block.BlockHash  `json:"flaghash"`
	SigHash  Block.BlockHash  `json:"sighash"`
}
