package bigchain

import (
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/ed25519"
	"github.com/zbo14/envoke/crypto/rsa"
	"github.com/zbo14/envoke/spec"
	"github.com/zbo14/envoke/spec/coala"
	"testing"
)

func TestBigchain(t *testing.T) {
	// Generate ed25519 keys
	privEd25519, pubEd25519 := ed25519.GenerateKeypair("password")
	// Create a data model
	data := coala.NewAlbum(spec.JSON, "", "name1", "artist1")
	// Generate tx
	tx := GenerateTx(data, nil, pubEd25519)
	// Fulfill the tx
	tx.Fulfill(privEd25519)
	// Check if it's fulfilled
	if !tx.Fulfilled() {
		t.Fatal("tx1 is not fulfilled")
	}
	// Print the tx
	json := MustMarshalIndentJSON(tx)
	Println(string(json))
	// Generate RSA key
	privRSA, pubRSA := rsa.GenerateKeypair()
	data = coala.NewAlbum(spec.JSON, "", "name2", "artist2")
	// Generate txs
	tx = GenerateTx(data, nil, pubRSA)
	// Fulfill the tx
	tx.Fulfill(privRSA)
	if !tx.Fulfilled() {
		t.Fatal("tx2 is not fulfilled")
	}
	// Print the tx
	json = MustMarshalIndentJSON(tx)
	Println(string(json))
	// Send POST request with tx
	// response, err := PostTx(tx)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(response)
}
