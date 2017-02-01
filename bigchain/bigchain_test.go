package bigchain

import (
	"github.com/zballs/envoke/coala"
	"github.com/zballs/envoke/crypto/ed25519"
	. "github.com/zballs/envoke/util"
	"testing"
)

func TestBigchain(t *testing.T) {
	// Generate keys
	priv, pub := ed25519.GenerateKeypair("password")
	// Create a new data model
	data := coala.NewRecording(coala.JSON, "", "artist", "album", "recording")
	// Generate a BigchainDB tx
	tx := GenerateTx(data, nil, pub)
	// Fulfill the tx
	tx.Fulfill(priv)
	json, _ := MarshalIndentJSON(tx)
	Println(string(json))
	// Check if it's fulfilled
	if !tx.Fulfilled() {
		t.Error("Transaction is not fulfilled")
	}
	response, err := PostTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response)
}
