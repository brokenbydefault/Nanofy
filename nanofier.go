package Nanofy

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/Block"
	"io"
	"encoding/binary"
	"golang.org/x/crypto/blake2b"
	"errors"
)

var addressbase = Wallet.PublicKey([]byte{0x51, 0x14, 0xab, 0x7c, 0x6a, 0xd0, 0xd6, 0xc3, 0x14, 0xc5, 0xc2, 0x8e, 0x36, 0xb0, 0x8a, 0x65, 0x0a, 0xd4, 0x2b, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

type Nanofier interface {
	CreateBlock(file io.Reader, sk Wallet.SecretKey, representative Wallet.Address, balance *Numbers.RawAmount, previous *Block.UniversalBlock) (blks []Block.BlockTransaction, err error)

	VerifySignature(pk *Wallet.PublicKey, flagb *Block.UniversalBlock, sigb *Block.UniversalBlock, previous *Block.UniversalBlock, filekey []byte) (ok bool)
	VerifyBlock(pk *Wallet.PublicKey, flagb *Block.UniversalBlock, sigb *Block.UniversalBlock, previous *Block.UniversalBlock) (ok bool)

	Amount() *Numbers.RawAmount
	Address() Wallet.Address
}

func NewNanofier(version uint64) (nanofier Nanofier, err error) {
	switch version {
	case 0:
		nanofier = NewNanofierVersion0()
	case 1:
		nanofier = NewNanofierVersion1()
	default:
		err = errors.New("unknown version")
	}

	return
}

func NewNanofierFromPublicKey(destination Wallet.PublicKey) (nanofier Nanofier, err error) {
	if len(destination) < 32 {
		return nil, errors.New("invalid hash")
	}

	return NewNanofier(binary.LittleEndian.Uint64(destination[24:]))
}

func NewNanofierFromFlagBlock(flagb Block.BlockTransaction) (nanofier Nanofier, err error) {
	destination, _ := flagb.GetTarget()

	pk, err := destination.GetPublicKey()
	if err != nil {
		return nil, err
	}

	return NewNanofierFromPublicKey(pk)
}

func CreatePublicKey(version uint64) Wallet.PublicKey {
	pk := make(Wallet.PublicKey, 32)
	copy(pk, addressbase)

	binary.LittleEndian.PutUint64(pk[24:], version)
	return pk
}

func CreateAddress(version uint64) Wallet.Address {
	return CreatePublicKey(version).CreateAddress()
}

func CreateHash(file io.Reader) (hash Wallet.PublicKey, err error) {
	blake, _ := blake2b.New(32, nil)

	_, err = io.Copy(blake, file)
	if err != nil {
		return nil, err
	}

	return blake.Sum(nil), nil
}
