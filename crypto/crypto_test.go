package crypto

import (
	conds "github.com/zballs/go_resonate/crypto/conditions"
	"github.com/zballs/go_resonate/crypto/ed25519"
	// . "github.com/zballs/go_resonate/util"
	"testing"
)

func TestCrypto(t *testing.T) {
	// Sha256 Pre-Image
	preimage := []byte("hello world")
	f1 := conds.NewFulfillmentPreImage(preimage, 1)
	if !f1.Validate(preimage) {
		t.Error("Failed to validate pre-image fulfillment")
	}
	// Ed25519
	msg := []byte("deadbeef")
	priv, _ := ed25519.GenerateKeypair("password")
	f2 := conds.NewFulfillmentEd25519(msg, priv, 1)
	if !f2.Validate(msg) {
		t.Error("Failed to validate ed25519 fulfillment")
	}
	// TODO: sha256 threshold
}
