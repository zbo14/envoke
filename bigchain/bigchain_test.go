package bigchain

import (
	"testing"

	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/crypto"
	"github.com/zbo14/envoke/crypto/ed25519"
)

var (
	Alice = "E8CKHfsze3YSKkcmo6Jhw8m57reuQPZJj1mRXwfdihCH"
	Bob   = "2BMuPCVkYbdKAN9qo83gARvaWEfQsv9RrjH6foHxsmTx"
)

func TestBigchain(t *testing.T) {
	output := MustOpenWriteFile("output.json")
	// Keys
	privAlice, pubAlice := ed25519.GenerateKeypairFromSeed(BytesFromB58(Alice))
	privBob, pubBob := ed25519.GenerateKeypairFromSeed(BytesFromB58(Bob))
	// Data
	data := Data{"bees": "knees"}
	// Individual create tx
	tx := DefaultIndividualCreateTx(data, pubAlice)
	FulfillTx(tx, privAlice)
	// Check that it's fulfilled
	if !FulfilledTx(tx) {
		t.Error(ErrInvalidFulfillment)
	}
	WriteJSON(output, Data{"createTx": tx})
	createTxId, err := PostTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	// Individual transfer tx
	tx = DefaultIndividualTransferTx(createTxId, createTxId, 0, pubBob, pubAlice)
	FulfillTx(tx, privAlice)
	if !FulfilledTx(tx) {
		t.Error(ErrInvalidFulfillment)
	}
	transferTxId, err := PostTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, Data{"transfer1Tx": tx})
	// Transfer the transfer tx
	tx = DefaultIndividualTransferTx(createTxId, transferTxId, 0, pubAlice, pubBob)
	FulfillTx(tx, privBob)
	if !FulfilledTx(tx) {
		t.Error(ErrInvalidFulfillment)
	}
	if _, err = PostTx(tx); err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, Data{"transfer2Tx": tx})
	// Multiple owners create tx
	tx = MultipleOwnersCreateTx([]int{2, 3}, data, []crypto.PublicKey{pubAlice, pubBob}, pubAlice)
	FulfillTx(tx, privAlice)
	if !FulfilledTx(tx) {
		t.Error(ErrInvalidFulfillment)
	}
	multipleOwnersTxId, err := PostTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, Data{"multipleOwnersTx": tx})
	tx, err = GetTx(multipleOwnersTxId)
	if err != nil {
		t.Fatal(err)
	}
}
