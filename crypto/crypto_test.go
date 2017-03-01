package crypto

import (
	"bytes"
	. "github.com/zbo14/envoke/common"
	conds "github.com/zbo14/envoke/crypto/conditions"
	"github.com/zbo14/envoke/crypto/ed25519"
	"github.com/zbo14/envoke/crypto/rsa"
	"sort"
	"testing"
)

func TestCrypto(t *testing.T) {
	// RSA-PEM encoding
	privRSA, pubRSA := rsa.GenerateKeypair()
	privPEM := privRSA.MarshalPEM()
	if err := privRSA.UnmarshalPEM(privPEM); err != nil {
		t.Error(err.Error())
	}
	pubPEM := pubRSA.MarshalPEM()
	if err := pubRSA.UnmarshalPEM(pubPEM); err != nil {
		t.Error(err.Error())
	}
	// Sha256 Pre-Image
	preimage := []byte("helloworld")
	f1 := conds.NewFulfillmentPreImage(preimage, 1)
	// Validate the fulfillment
	if !f1.Validate(preimage) {
		t.Error("Failed to validate pre-image fulfillment")
	}
	// Sha256 Prefix
	prefix := []byte("hello")
	suffix := []byte("world")
	f2 := conds.NewFulfillmentPrefix(prefix, f1, 1)
	// Validate the fulfillment
	if !f2.Validate(suffix) {
		t.Error("Failed to validate prefix fulfillment")
	}
	// Ed25519
	msg := []byte("deadbeef")
	privEd25519, _ := ed25519.GenerateKeypairFromPassword("password")
	f3 := conds.FulfillmentFromPrivKey(msg, privEd25519, 2)
	if !f3.Validate(msg) {
		t.Error("Failed to validate ed25519 fulfillment")
	}
	// RSA
	anotherMsg := []byte("foobar")
	f4 := conds.FulfillmentFromPrivKey(anotherMsg, privRSA, 1)
	if !f4.Validate(anotherMsg) {
		t.Error("Failed to validate pre-image fulfillment")
	}
	// Sha256 Threshold
	subs := conds.Fulfillments{f1, f2, f3, f4}
	sort.Sort(subs)
	threshold := 4
	f5 := conds.NewFulfillmentThreshold(subs, threshold, 1)
	buf := new(bytes.Buffer)
	MustWriteVarBytes(msg, buf)
	MustWriteVarBytes(preimage, buf)
	MustWriteVarBytes(suffix, buf)
	MustWriteVarBytes(anotherMsg, buf)
	if !f5.Validate(buf.Bytes()) {
		t.Error("Failed to validate threshold fulfillment")
	}
	PrintJSON(f5)
	// Get fulfillment uri
	uri := f5.String()
	// Derive new fulfillment from uri, use same weight
	f6, err := conds.UnmarshalURI(uri, 1)
	if err != nil {
		t.Fatal(err.Error())
	}
	// Check whether hashes are the same
	if !bytes.Equal(f5.Hash(), f6.Hash()) {
		t.Error("Expected identical fulfillment hashes")
	}
	// Nested Thresholds
	subs = conds.Fulfillments{f1, f2, f3, f4, f5}
	sort.Sort(subs)
	buf2 := new(bytes.Buffer)
	MustWriteVarBytes(msg, buf2)
	MustWriteVarBytes(preimage, buf2)
	MustWriteVarBytes(suffix, buf2)
	MustWriteVarBytes(buf.Bytes(), buf2)
	MustWriteVarBytes(anotherMsg, buf2)
	threshold = 4
	f7 := conds.NewFulfillmentThreshold(subs, threshold, 1)
	if !f7.Validate(buf2.Bytes()) {
		t.Error("Failed to validate nested thresholds")
	}
}
