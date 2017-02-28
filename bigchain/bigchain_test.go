package bigchain

import (
	"testing"

	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/ed25519"
)

var (
	Alice = "3Ee5HXoheGJf7DDjYSGPmL1bfZ2K2ZkRZ91jxZ7vtFnf"
	Bob   = "4oHxCNBADpwoWHGcBueLb36UQecUKLASGj7pWTT1scf6"
)

func TestBigchain(t *testing.T) {
	// Create tx
	privAlice, pubAlice := ed25519.GenerateKeypairFromSeed(BytesFromB58(Alice))
	tx := IndividualCreateTx(Data{"dummy": "dummy"}, pubAlice)
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
	// Transfer tx
	_, pubBob := ed25519.GenerateKeypairFromSeed(BytesFromB58(Bob))
	tx = IndividualTransferTx(txId, 0, pubBob, pubAlice)
	FulfillTx(tx, privAlice)
	if !FulfilledTx(tx) {
		t.Error(ErrInvalidFulfillment)
	}
	// PrintJSON(tx)
	txId, err = PostTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txId)
}
