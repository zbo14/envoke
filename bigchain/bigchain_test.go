package bigchain

import (
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/ed25519"
	"github.com/zbo14/envoke/spec"
	"testing"
)

var seedstr = "13jGvCoZsEiqu5kiBLFz8vPVUS5pchjkQFmeP2bNbHae"

func TestBigchain(t *testing.T) {
	// Keys
	seed := BytesFromB58(seedstr)
	priv, pub := ed25519.GenerateKeypairFromSeed(seed)
	// Data model
	info := spec.NewCompositionInfo("composerId", "publisherId", "title")
	// Create tx
	tx := IndividualCreateTx(info, pub)
	FulfillTx(tx, priv)
	// Check if it's fulfilled
	if !FulfilledTx(tx) {
		t.Error("tx is not fulfilled")
	}
	txId, err := PostTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txId)
	// Transfer tx
	_, pubNew := ed25519.GenerateKeypairFromPassword("password")
	tx = IndividualTransferTx(txId, 0, pubNew, pub)
	PrintJSON(tx)
	FulfillTx(tx, priv)
	if !FulfilledTx(tx) {
		t.Error("tx is not fulfilled")
	}
	PrintJSON(tx)
	txId, err = PostTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txId)
}
