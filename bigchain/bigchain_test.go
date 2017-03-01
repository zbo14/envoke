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
	// Keys
	privAlice, pubAlice := ed25519.GenerateKeypairFromSeed(BytesFromB58(Alice))
	privBob, pubBob := ed25519.GenerateKeypairFromSeed(BytesFromB58(Bob))
	// Data
	data := Data{"dummy": "dummy"}
	// Individual create tx
	tx := DefaultIndividualCreateTx(data, pubAlice)
	FulfillTx(tx, privAlice)
	// Check if it's fulfilled
	if !FulfilledTx(tx) {
		t.Error(ErrInvalidFulfillment)
	}
	// PrintJSON(tx)
	createTxId, err := PostTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(createTxId)
	// Individual transfer tx
	tx = DefaultIndividualTransferTx(createTxId, createTxId, 0, pubBob, pubAlice)
	FulfillTx(tx, privAlice)
	if !FulfilledTx(tx) {
		t.Error(ErrInvalidFulfillment)
	}
	PrintJSON(tx)
	transferTxId, err := PostTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(transferTxId)
	// Transfer the transfer tx
	tx = DefaultIndividualTransferTx(createTxId, transferTxId, 0, pubAlice, pubBob)
	FulfillTx(tx, privBob)
	if !FulfilledTx(tx) {
		t.Error(ErrInvalidFulfillment)
	}
	transfer2TxId, err := PostTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(transfer2TxId)
	// Multiple owners create tx
	tx = MultipleOwnersCreateTx([]int{2, 3}, data, []crypto.PublicKey{pubAlice, pubBob}, pubAlice)
	FulfillTx(tx, privAlice)
	if !FulfilledTx(tx) {
		t.Error(ErrInvalidFulfillment)
	}
	PrintJSON(tx)
	// PrintJSON(tx)
	multipleOwnersTxId, err := PostTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(multipleOwnersTxId)
	tx, err = GetTx(multipleOwnersTxId)
	if err != nil {
		t.Fatal(err)
	}
	outputs := GetTxOutputs(tx)
	pubs := GetOutputsPublicKeys(outputs)
	t.Log(pubs)
	inputs := GetTxInputs(tx)
	pubs = GetInputsPublicKeys(inputs)
	t.Log(pubs)
}
