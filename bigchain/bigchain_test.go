package bigchain

import (
	"github.com/zballs/envoke/crypto/ed25519"
	"github.com/zballs/envoke/spec"
	"github.com/zballs/envoke/spec/coala"
	. "github.com/zballs/envoke/util"
	"testing"
)

func TestBigchain(t *testing.T) {
	// Generate keys
	priv, pub := ed25519.GenerateKeypair("password")
	// Create a data model
	data := coala.NewAlbum(spec.JSON, "", "name", "artist")
	// Generate tx
	tx := GenerateTx(data, nil, pub)
	// Fulfill the tx
	tx.Fulfill(priv)
	json, _ := MarshalIndentJSON(tx)
	Println(string(json))
	// Check if it's fulfilled
	if !tx.Fulfilled() {
		t.Fatal("Transaction is not fulfilled")
	}
	response, err := PostTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response)
}
