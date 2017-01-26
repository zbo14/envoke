package bigchain

import (
	"github.com/zballs/go_resonate/coala"
	"github.com/zballs/go_resonate/crypto"
	"testing"
)

func TestBigchain(t *testing.T) {
	priv, pub := crypto.GenerateKeypair("password")
	data := coala.NewRecording("artist", "album", "recording")
	transaction := GenerateTransaction(data, pub)
	transaction.Fulfill(priv, pub)
	// TODO:
}
