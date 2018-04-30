package Nanofy

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"io"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"bytes"
	"errors"
)

type version1 struct {
	address Wallet.Address
	amount  *Numbers.RawAmount
}

func NewNanofierVersion1() Nanofier {
	amm, _ := Numbers.NewRawFromString("1")
	return &version1{
		address: CreateAddress(1),
		amount:  amm,
	}
}

// CreateFileBlocks will return two blocks, the first with the address as the hash of the file and the second one is the
// flag block. It returns a non-nil error if something go wrong.
func (v version1) CreateBlock(file io.Reader, sk Wallet.SecretKey, rep Wallet.Address, balance *Numbers.RawAmount, previous *Block.UniversalBlock) (blks []Block.BlockTransaction, err error) {
	hash, err := CreateHash(file)
	if err != nil {
		return blks, err
	}

	if previous == nil || previous.Type != Block.State {
		return blks, errors.New("not allowed")
	}

	bfile, remain, err := v.createBlock(sk, hash.CreateAddress(), rep, balance, previous.Hash())
	if err != nil {
		return blks, err
	}
	blks = append(blks, bfile)

	bflag, remain, err := v.createBlock(sk, v.address, rep, remain, bfile.Hash())
	if err != nil {
		return blks, err
	}
	blks = append(blks, bflag)

	return blks, err
}

func (v version1) createBlock(sk Wallet.SecretKey, dest Wallet.Address, rep Wallet.Address, balance *Numbers.RawAmount, frontier []byte) (blk Block.BlockTransaction, remain *Numbers.RawAmount, err error) {
	remain = balance.Subtract(v.amount)
	blk, err = Block.CreateSignedUniversalSendBlock(sk, rep, balance, v.amount, frontier, dest)

	return blk, remain, err
}

// VerifySignature compares two blocks and the file, returns true if the combination of blocks, pk and file is valid.
func (v version1) VerifySignature(pk *Wallet.PublicKey, flagb *Block.UniversalBlock, sigb *Block.UniversalBlock, previous *Block.UniversalBlock, filekey []byte) (ok bool) {
	if !v.VerifyBlock(pk, flagb, sigb, previous) {
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
func (v version1) VerifyBlock(pk *Wallet.PublicKey, flagb *Block.UniversalBlock, sigb *Block.UniversalBlock, previous *Block.UniversalBlock) (ok bool) {

	// The destination of the flag MUST be the Version1 address
	destination, _ := flagb.GetTarget()
	if destination != v.address {
		return false
	}

	// The previous MUST point to each other, "Flag" previous point to "Sig", "Sig" previous point to "Previous".
	if !bytes.Equal(flagb.Previous, sigb.Hash()) || !bytes.Equal(sigb.Previous, previous.Hash()) {
		return false
	}

	// The type MUST be "state"
	if sigb.GetType() != Block.State || flagb.GetType() != Block.State || previous.GetType() != Block.State {
		return false
	}

	// The signature of the block itself MUST be a correct one
	if !pk.CompareSignature(flagb.Hash(), flagb.Signature) || !pk.CompareSignature(sigb.Hash(), sigb.Signature) || !pk.CompareSignature(previous.Hash(), previous.Signature) {
		return false
	}

	// The blocks MUST send only 1 raw, it implicitly guaranties that is a send operation over state-block.
	if sigb.Balance.Subtract(flagb.Balance).Compare(v.amount) != 0 || previous.Balance.Subtract(sigb.Balance).Compare(v.amount) != 0 {
		return false
	}

	return true
}

func (v version1) Address() Wallet.Address {
	return v.address
}

func (v version1) Amount() *Numbers.RawAmount {
	return v.amount
}
