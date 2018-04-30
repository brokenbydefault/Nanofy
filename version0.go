package Nanofy

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"io"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"bytes"
)

type version0 struct {
	address Wallet.Address
	amount  *Numbers.RawAmount
}

func NewNanofierVersion0() Nanofier {
	amm, _ := Numbers.NewRawFromString("1")
	return &version0{
		address: CreateAddress(0),
		amount:  amm,
	}
}

// CreateFileBlocks will return two blocks, the first with the address as the hash of the file and the second one is the
// flag block. It returns a non-nil error if something go wrong.
func (v version0) CreateBlock(file io.Reader, sk Wallet.SecretKey, _ Wallet.Address, balance *Numbers.RawAmount, previous *Block.UniversalBlock) (blks []Block.BlockTransaction, err error) {
	hash, err := CreateHash(file)
	if err != nil {
		return
	}

	bfile, remain, err := v.createBlock(sk, hash.CreateAddress(), balance, previous.SwitchTo(Block.Send).Hash())
	if err != nil {
		return
	}
	blks = append(blks, bfile)

	bflag, remain, err := v.createBlock(sk, v.address, remain, bfile.SwitchToUniversalBlock().SwitchTo(Block.Send).Hash())
	if err != nil {
		return
	}
	blks = append(blks, bflag)

	return
}

func (v version0) createBlock(sk Wallet.SecretKey, dest Wallet.Address, balance *Numbers.RawAmount, frontier []byte) (blk Block.BlockTransaction, remain *Numbers.RawAmount, err error) {
	remain = balance.Subtract(v.amount)
	blk, err = Block.CreateSignedSendBlock(sk, v.amount, balance, frontier, dest)

	return
}

// VerifySignature compares two blocks and the file, returns true if the combination of blocks, pk and file is valid.
func (v version0) VerifySignature(pk *Wallet.PublicKey, flagb *Block.UniversalBlock, sigb *Block.UniversalBlock, _ *Block.UniversalBlock, filekey []byte) (ok bool) {
	if !v.VerifyBlock(pk, flagb, sigb, nil) {
		return false
	}

	destinationaddr, _ := sigb.GetTarget()
	destination, _ := destinationaddr.GetPublicKey()

	if !bytes.Equal(destination, filekey) {
		return false
	}

	return true
}

// VerifyBlock lacks the file verification, which can be used to verify if one block is one correct Nanofy transaction,
// in this case we don't care about the hash of file, because we can don't have the file itself.
func (v version0) VerifyBlock(pk *Wallet.PublicKey, flagb *Block.UniversalBlock, sigb *Block.UniversalBlock, _ *Block.UniversalBlock) (ok bool) {

	// The destination of the flag MUST be the Version0 address
	if flagb.Destination != v.address {
		return false
	}

	// The previous MUST point to each other, "Flag" previous point to "Sig".
	if !bytes.Equal(flagb.Previous, sigb.SwitchTo(Block.Send).Hash()) {
		return false
	}

	// The type MUST be "send" (legacy one)
	if sigb.GetType() != Block.Send || flagb.GetType() != Block.Send {
		return false
	}

	// The signature of the block itself MUST be a correct one
	if !pk.CompareSignature(flagb.SwitchTo(Block.Send).Hash(), flagb.Signature) || !pk.CompareSignature(sigb.SwitchTo(Block.Send).Hash(), sigb.Signature) {
		return false
	}

	// The blocks MUST send only 1 raw.
	if sigb.Balance.Subtract(flagb.Balance).Compare(v.amount) != 0 {
		return false
	}

	return true
}

func (v version0) Address() Wallet.Address {
	return v.address
}

func (v version0) Amount() *Numbers.RawAmount {
	return v.amount
}
