package Nanofy

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Block"
	"bytes"
	"encoding/binary"
	"github.com/brokenbydefault/Nanollet/Numbers"
)

// VerifySignature compares two blocks and the file, returns true if the combination of blocks, pk and file is valid.
func VerifySignature(publickey *Wallet.PublicKey, filekey []byte, flagb Block.UniversalBlock, sigb Block.UniversalBlock) (ok bool) {
	if !VerifyBlock(publickey, flagb, sigb) {
		return false
	}

	if !bytes.Equal([]byte(sigb.Destination), filekey) {
		return false
	}

	return true
}

// VerifyBlock lacks the file verification, which can be used to verify if one block is one correct Nanofy transaction,
// in this case we don't care about the hash of file, because we can don't have the file itself.
func VerifyBlock(pk *Wallet.PublicKey, flagb Block.UniversalBlock, sigb Block.UniversalBlock) (ok bool) {
	var block string
	switch binary.LittleEndian.Uint64(flagb.Destination[24:]) {
	case uint64(NV0):
		block = "send"
		//case uint32(NV1):
		//block = "utx"
	default:
		return false
	}

	if !pk.CompareSignature(flagb.HashAsSend(), flagb.Signature) || !pk.CompareSignature(sigb.HashAsSend(), sigb.Signature) {
		return false
	}

	if sigb.Type != block || flagb.Type != block {
		return false
	}

	if !bytes.Equal(flagb.Previous, sigb.HashAsSend()) {
		return false
	}

	sendValue, _ := Numbers.NewRawFromString("1")
	if sigb.Balance.Subtract(flagb.Balance).Compare(sendValue) != 0 {
		return false
	}

	return true
}
