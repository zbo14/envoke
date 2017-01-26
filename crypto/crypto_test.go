package crypto

import (
	"github.com/zballs/go_resonate/crypto/conditions"
	"github.com/zballs/go_resonate/crypto/ed25519"
	. "github.com/zballs/go_resonate/util"
	"testing"
)

func TestCrypto(t *testing.T) {
	priv, pub := ed25519.GenerateKeypair("password")
	fulfillment := conds.NewFulfillmentEd25519([]byte("message"), priv)
	condition := fulfillment.Condition()
	// TODO:
}
