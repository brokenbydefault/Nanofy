package Nanofy

import (
	"testing"
	"github.com/brokenbydefault/Nanollet/RPC"
	"github.com/brokenbydefault/Nanollet/RPC/Connectivity"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"bytes"
)

func TestVerifyBlock(t *testing.T) {

	flagblock, _ := RPCClient.GetBlockByStringHash(Connectivity.HTTP, "AE95D7936A23D0671DD4E7E0736612F5304A18AD80B2827B2D69A3482A38F1EA")
	sigblock, _ := RPCClient.GetBlockByStringHash(Connectivity.HTTP, "7F0AEA1B6F2E9FD60B120DF01B0FBF4CC1B4B539A1D6B69DC50EAE81FE1A72E7")

	pk, _ := Wallet.Address("xrb_3w73pgb33ht1ws7hwaek5ywyjdteoj4qmcrzayiogpbabbo3i49dkerosn1z").GetPublicKey()

	if !VerifyBlock(&pk, flagblock, sigblock) {
		t.Error("one valid block was report as invalid")
	}

	flagblock, _ = RPCClient.GetBlockByStringHash(Connectivity.HTTP, "AE95D7936A23D0671DD4E7E0736612F5304A18AD80B2827B2D69A3482A38F1EA")
	sigblock, _ = RPCClient.GetBlockByStringHash(Connectivity.HTTP, "72DC2B79C307600FE8187521DB2C0AAA2929D6E10C1E3E3B058ACB6B617EB019")

	pk, _ = Wallet.Address("xrb_3w73pgb33ht1ws7hwaek5ywyjdteoj4qmcrzayiogpbabbo3i49dkerosn1z").GetPublicKey()

	if VerifyBlock(&pk, flagblock, sigblock) {
		t.Error("one invalid block was report as valid")
	}

}

func TestVersion_CreateFlagPublicKey(t *testing.T) {
	if !bytes.Equal(*NV0.CreateFlagPublicKey(), AddressBase) {
		t.Error("wrong address")
	}

}
