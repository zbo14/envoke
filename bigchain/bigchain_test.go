package bigchain

import (
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/ed25519"
	"github.com/zbo14/envoke/spec"
	"github.com/zbo14/envoke/spec/core"
	"testing"
)

var seedstr = "13jGvCoZsEiqu5kiBLFz8vPVUS5pchjkQFmeP2bNbHae"

func TestBigchain(t *testing.T) {
	// Keys
	seed := BytesFromB58(seedstr)
	priv, pub := ed25519.GenerateKeypairFromSeed(seed)
	// Data model
	artist := core.NewArtist("artist@email.com", "artist_name", pub)
	metadata := spec.Data{"dummy": "dummy"}
	// Generate tx
	tx := GenerateTx(artist, metadata, CREATE, pub)
	// Print prepared tx
	// PrintJSON(tx)
	// Fulfill the tx
	FulfillTx(tx, priv)
	// Check if it's fulfilled
	if !FulfilledTx(tx) {
		t.Error("tx is not fulfilled")
	}
	// Print fulfilled tx
	// PrintJSON(tx)
	// Send POST request with tx
	txId, err := PostTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	tx, err = GetTx(txId)
	if err != nil {
		t.Fatal(err)
	}
	data := GetTxData(tx)
	// Println(data)
	artist = new(core.Agent)
	if err := FillStruct(artist, data); err != nil {
		t.Fatal(err)
	}
	t.Log(artist)
}
