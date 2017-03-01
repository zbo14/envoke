package bigchain

import (
	"testing"

	. "github.com/zbo14/envoke/common"
	conds "github.com/zbo14/envoke/crypto/conditions"
	"github.com/zbo14/envoke/crypto/crypto"
	"github.com/zbo14/envoke/crypto/ed25519"
)

var (
	Alice = "E1dMSsqjitt1JdZZvSoR2YPuT3Pwmi9HoTVUSJEbuscn"
	Bob   = "HAogmr4HKpvzuYK9fC2ogayRE4KwJab64xQ7yCkF5HMf"
)

func TestBigchain(t *testing.T) {
	// Keys
	_, pubAlice := ed25519.GenerateKeypairFromSeed(BytesFromB58(Alice))
	_, pubBob := ed25519.GenerateKeypairFromSeed(BytesFromB58(Bob))
	// data := Data{"dummy": "dummy"}
	/*
		// Create tx
		tx := IndividualCreateTx(data, pubAlice)
		FulfillTx(tx, privAlice)
		// Check if it's fulfilled
		if !FulfilledTx(tx) {
			t.Error(ErrInvalidFulfillment)
		}
		// PrintJSON(tx)
		txId, err := PostTx(tx)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(txId)
		// Transfer tx
		tx = IndividualTransferTx(txId, 0, pubBob, pubAlice)
		FulfillTx(tx, privAlice)
		if !FulfilledTx(tx) {
			t.Error(ErrInvalidFulfillment)
		}
		// PrintJSON(tx)
		txId, err = PostTx(tx)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(txId)
	*/

	// Threshold Fulfillment from pubkeys
	f1 := conds.UnmarshalURI("cf:4:7Bcrk61eVjv0kyxw4SRQNMNUZ-8u_U1k6_gZaDRn4r-2IpH62UMvjymLnEpIldvik_b_2hpo2t8Mze9fR6DHISpf6jzal6P0wD6p8uisHOyGpR1FISer26CdG28zHAcK", 1)
	f2 := conds.UnmarshalURI("cf:0:", weight)
	fulfillmentThreshold := conds.DefaultFulfillmentThresholdFromPubKeys([]crypto.PublicKey{pubAlice, pubBob})
	// tx := MultipleOwnersCreateTx(data, []crypto.PublicKey{pubAlice, pubBob}, pubAlice)
	PrintJSON(fulfillmentThreshold)
}
