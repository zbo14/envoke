package crypto

import (
	"bytes"
	conds "github.com/zballs/go_resonate/crypto/conditions"
	"github.com/zballs/go_resonate/crypto/ed25519"
	. "github.com/zballs/go_resonate/util"
	"testing"
)

func TestCryptoConditions(t *testing.T) {
	// Sha256 Pre-Image
	preimage := []byte("hello world")
	f1 := conds.NewFulfillmentPreImage(preimage, 1)
	// Validate the fulfillment
	if !f1.Validate(preimage) {
		t.Error("Failed to validate pre-image fulfillment")
	}
	// Peep the condition
	t.Log(f1.Condition())
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
	buf := new(bytes.Buffer)
	MustWriteVarBytes(msg, buf)
	MustWriteVarBytes(preimage, buf)
	MustWriteVarBytes(preimage2, buf)
	if !f4.Validate(buf.Bytes()) {
		t.Error("Failed to validate threshold fulfillment")
	}
	// Get fulfillment uri
	uri := f4.String()
	// Derive new fulfillment from uri
	f5, err := conds.UnmarshalURI(uri)
	if err != nil {
		t.Fatal(err.Error())
	}
	// Check whether hashes are the same
	if !bytes.Equal(f4.Hash(), f5.Hash()) {
		t.Error("Expected identical fulfillment hashes")
	}
	// Nested Thresholds
	subs = conds.Fulfillments{f1, f2, f3, f4}
	threshold = 4
	f6 := conds.NewFulfillmentThreshold(subs, threshold, 1)
	buf.Write(buf.Bytes())
	if !f6.Validate(buf.Bytes()) {
		t.Error("Failed to validate nested thresholds")
	}
}
