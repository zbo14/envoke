package bigchain

import (
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/ed25519"
	"github.com/zbo14/envoke/spec"
	"github.com/zbo14/envoke/spec/coala"
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
	// Send POST request with tx
	response, err := PostTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response)
}
