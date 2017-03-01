package bigchain

import (
	"testing"

	. "github.com/zbo14/envoke/common"
	conds "github.com/zbo14/envoke/crypto/conditions"
	"github.com/zbo14/envoke/crypto/crypto"
	"github.com/zbo14/envoke/crypto/ed25519"
)

var (
	Alice = "E8CKHfsze3YSKkcmo6Jhw8m57reuQPZJj1mRXwfdihCH"
	Bob   = "2BMuPCVkYbdKAN9qo83gARvaWEfQsv9RrjH6foHxsmTx"
)

func TestBigchain(t *testing.T) {
	// Keys
	privAlice, pubAlice := ed25519.GenerateKeypairFromSeed(BytesFromB58(Alice))
	_, pubBob := ed25519.GenerateKeypairFromSeed(BytesFromB58(Bob))
	// Data
	data := Data{"dummy": "dummy"}
	// Individual create tx
	tx := IndividualCreateTx(data, pubAlice)
	FulfillTx(tx, privAlice)
	// Check if it's fulfilled
	if !FulfilledTx(tx) {
		t.Error(ErrInvalidFulfillment)
	}
	// PrintJSON(tx)
	txId, err := PostTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txId)
	// Individual transfer tx
	tx = IndividualTransferTx(txId, 0, pubBob, pubAlice)
	FulfillTx(tx, privAlice)
	if !FulfilledTx(tx) {
		t.Error(ErrInvalidFulfillment)
	}
	PrintJSON(tx)
	txId, err = PostTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txId)
	// Multiple owners create tx
	fulfillmentThreshold := conds.DefaultFulfillmentThresholdFromPubKeys([]crypto.PublicKey{pubAlice, pubBob})
	tx = MultipleOwnersCreateTx(data, []crypto.PublicKey{pubAlice, pubBob}, pubAlice)
	FulfillTx(tx, privAlice)
	if !FulfilledTx(tx) {
		t.Error(ErrInvalidFulfillment)
	}
	PrintJSON(tx)
	PrintJSON(fulfillmentThreshold)
	// PrintJSON(tx)
	txId, err = PostTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txId)
	tx, err = GetTx(txId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(GetTxPublicKeys(tx))
}
