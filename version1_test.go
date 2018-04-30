package Nanofy

import (
	"testing"
	"github.com/brokenbydefault/Nanollet/RPC"
	"github.com/brokenbydefault/Nanollet/RPC/Connectivity"
	"github.com/brokenbydefault/Nanollet/Wallet"
)

func TestVersion1_VerifyBlock(t *testing.T) {

	flagblock, _ := RPCClient.GetBlockByStringHash(Connectivity.Socket, "A7CEB0E504E31D74CD87DB0990B9C44D8211987C464BE758C6007C6DA2D188FB")
	sigblock, _ := RPCClient.GetBlockByStringHash(Connectivity.Socket, "02C03EEFBFAC781971125D04EB0F28D510D7B303E86D1CAB22E0C86271862E9B")
	prevblock, _ := RPCClient.GetBlockByStringHash(Connectivity.Socket, "EF7A37D750DBB692F957C0EA2965B1F48EB2BC2504AE7AE07957904D83D5B267")

	pk, _ := Wallet.Address("xrb_31hbrc4zary87ardrg74pd6xy157z71t4b9edmrqbj5tqgxnriba5e7cf3o6").GetPublicKey()

	nanofier, err := NewNanofierFromFlagBlock(&flagblock)
	if err != nil {
		panic("invalid version")
	}

	if !nanofier.VerifyBlock(&pk, &flagblock, &sigblock, &prevblock) {
		t.Error("one valid block was report as invalid")
	}

}
