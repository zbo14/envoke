package bigchain

import (
	"github.com/zballs/go_resonate/coala"
	"github.com/zballs/go_resonate/crypto/ed25519"
	"testing"
)

func TestBigchain(t *testing.T) {
	// Generate public-private keys
	priv, pub := ed25519.GenerateKeypair("password")
	// Create a new data model
	data := coala.NewRecording("artist", "album", "recording")
	// Generate a BigchainDB transaction
	transaction := GenerateTransaction(data, pub)
	// Fulfill the transaction
	transaction.Fulfill(priv, pub)
	// Check if it's fulfilled
	if !transaction.Fulfilled() {
		t.Error("Transaction is not fulfilled")
	}
	id, err := PostTransaction(transaction)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(id)
}
