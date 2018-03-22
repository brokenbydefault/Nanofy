package Nanofy

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"io"
	"github.com/brokenbydefault/Nanollet/Block"
	"errors"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"encoding/binary"
)

var AddressBase = Wallet.PublicKey([]byte{0x51, 0x14, 0xab, 0x7c, 0x6a, 0xd0, 0xd6, 0xc3, 0x14, 0xc5, 0xc2, 0x8e, 0x36, 0xb0, 0x8a, 0x65, 0x0a, 0xd4, 0x2b, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

type Version uint64
const (
	NV0 Version = 0
	//NV1
)

// CreateFileBlocks will return two blocks, the first with the address as the hash of the file and the second one is the
// flag block. It returns a non-nil error if something go wrong.
func (v Version) CreateFileBlocks(file io.Reader, sk *Wallet.SecretKey, balance *Numbers.RawAmount, prev []byte) (blks []*Block.SendBlock, err error) {
	hash, err := CreateFileHash(file)
	if err != nil {
		return
	}

	filepk := Wallet.PublicKey(hash)

	bfile, remain, err := v.createBlock(sk, &filepk, balance, prev)
	if err != nil {
		return
	}
	blks = append(blks, bfile)

	bflag, remain, err := v.createBlock(sk, v.CreateFlagPublicKey(), remain, bfile.Hash())
	if err != nil {
		return
	}
	blks = append(blks, bflag)

	return
}

func (v Version) CreateFileBlocksWithPoW(file io.Reader, sk *Wallet.SecretKey, balance *Numbers.RawAmount, prev []byte) (blks []*Block.SendBlock, err error) {
	blks, err = v.CreateFileBlocks(file, sk, balance, prev)
	for _, blk := range blks {
		blk.CreateProof()
	}

	return
}

func (v Version) createBlock(sk *Wallet.SecretKey, dest *Wallet.PublicKey, balance *Numbers.RawAmount, prev []byte) (blk *Block.SendBlock, remain *Numbers.RawAmount, err error) {
	sendAddr := dest.CreateAddress()
	sendValue := v.Value()

	remain = balance.Subtract(sendValue)

	switch v {
	case NV0:
		blk, err = Block.CreateSignedSendBlock(sk, sendValue, balance, prev, &sendAddr)
		//case NV1:
		//blk, err = Block.CreateUniversalSignedSendBlock(sk, sendValue, balance, previous, &sendAddr)
	default:
		err = errors.New("unsupported version")
	}

	return
}

func (v Version) Value() (*Numbers.RawAmount) {
	n, _ := Numbers.NewRawFromString("1")
	return n
}

// CreateFlagPublicKey creates the public-key of the flag address, which is like
// [xbr_1nanofy8on8preceding8transaction][padding][version]
func (v Version) CreateFlagPublicKey() (address *Wallet.PublicKey) {
	pk := AddressBase

	binary.LittleEndian.PutUint64(pk[24:], uint64(v))
	return &pk
}
