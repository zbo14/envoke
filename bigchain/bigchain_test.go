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
	// Dummy data
	dummy := spec.Data{"dummy": "dummy"}
	// Generate tx
	tx := GenerateTx(dummy, dummy, CREATE, pub)
	// Print prepared tx
	Println(string(MustMarshalIndentJSON(tx)))
	// Fulfill the tx
	FulfillTx(tx, priv)
	// Check if it's fulfilled
	if !FulfilledTx(tx) {
		t.Error("tx is not fulfilled")
	}
	// Print fulfilled tx
	Println(string(MustMarshalIndentJSON(tx)))
	// Send POST request with tx
	response, err := PostTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response)
}
