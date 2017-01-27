package bigchain

import (
	"github.com/zballs/go_resonate/coala"
	"github.com/zballs/go_resonate/crypto/ed25519"
	"testing"
)

func TestBigchain(t *testing.T) {
	priv, pub := ed25519.GenerateKeypair("password")
	data := coala.NewRecording("artist", "album", "recording")
	transaction := GenerateTransaction(data, pub)
	transaction.Fulfill(priv, pub)
}
