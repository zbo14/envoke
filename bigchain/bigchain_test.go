package bigchain

import (
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/ed25519"
	"github.com/zbo14/envoke/spec"
	"testing"
)

const (
	privstr = "iyVaTCBKcuHEn93vxRkSvGehadeMkGL13KQS5Yas2wwRKu4GT3tDgUZ5C2NVz78pYE6NuCUtrASqpaa6rUCoRR1"
	pubstr  = "8w1T7fTbjcbB3si6oD3HMrBZENXNaFAkEStueKzASwxj"
)

func TestBigchain(t *testing.T) {
	// Keys
	priv := new(ed25519.PrivateKey)
	if err := priv.FromString(privstr); err != nil {
		t.Fatal(err.Error())
	}
	pub := new(ed25519.PublicKey)
	if err := pub.FromString(pubstr); err != nil {
		t.Fatal(err.Error())
	}
	// Dummy data
	dummy := spec.Data{"dummy": "dummy"}
	// Generate tx
	tx := GenerateTx(dummy, dummy, CREATE, pub)
	// Fulfill the tx
	tx.Fulfill(priv)
	// Check if it's fulfilled
	if !tx.Fulfilled() {
		t.Fatal("tx is not fulfilled")
	}
	// Print tx
	Println(string(MustMarshalIndentJSON(tx)))
	// Send POST request with tx
	response, err := PostTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response)
}
