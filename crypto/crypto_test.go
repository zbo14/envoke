package crypto

import (
	"bytes"
	conds "github.com/zballs/go_resonate/crypto/conditions"
	"github.com/zballs/go_resonate/crypto/ed25519"
	. "github.com/zballs/go_resonate/util"
	"testing"
)

func TestCrypto(t *testing.T) {
	// Sha256 Pre-Image
	preimage := []byte("hello world")
	f1 := conds.NewFulfillmentPreImage(preimage, 1)
	// Validate the fulfillment
	if !f1.Validate(preimage) {
		t.Error("Failed to validate pre-image fulfillment")
	}
	// Print the condition
	c1 := f1.Condition()
	Println(c1)
	// Ed25519
	msg := []byte("deadbeef")
	priv, _ := ed25519.GenerateKeypair("password")
	f2 := conds.NewFulfillmentEd25519(msg, priv, 2)
	if !f2.Validate(msg) {
		t.Error("Failed to validate ed25519 fulfillment")
	}
	// Another Pre-Image
	preimage2 := []byte("foobar")
	f3 := conds.NewFulfillmentPreImage(preimage2, 1)
	if !f3.Validate(preimage) {
		t.Error("Failed to validate pre-image fulfillment")
	}
	// Sha256 Threshold
	subs := conds.Fulfillments{f1, f2, f3}
	threshold := 3
	f4 := conds.NewFulfillmentThreshold(subs, threshold, 1)
	Println(subs)
	buf := new(bytes.Buffer)
	MustWriteVarBytes(msg, buf)
	MustWriteVarBytes(preimage, buf)
	MustWriteVarBytes(preimage2, buf)
	if !f4.Validate(buf.Bytes()) {
		t.Error("Failed to validate threshold fulfillment")
	}
}
