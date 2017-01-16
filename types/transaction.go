package types

import (
	"fmt"
	. "github.com/zballs/go_resonate/util"
)

// BigchainDB-compatible transaction type

type Operation string

const (
	// For ed25519
	BITMASK            = 32
	FULFILLMENT_LENGTH = PUBKEY_LENGTH + SIGNATURE_LENGTH
	TYPE               = "fulfillment"
	TYPE_ID            = 4

	// Operation types
	CREATE   Operation = "CREATE"
	GENESIS  Operation = "GENESIS"
	TRANSFER Operation = "TRANSFER"
)

type Transaction struct {
	Id      string `json:"id"`
	Tx      *Tx    `json:"tx"`
	Version int    `json:"version"`
}

func NewTransaction(tx *Tx, version int) *Transaction {
	transaction := &Transaction{
		Tx:      tx,
		Version: version,
	}
	conditions := transaction.Tx.Conditions
	sigs := make([]string, len(conditions))
	for i, c := range conditions {
		sigs[i] = c.Cond.Details.Signature
		c.Cond.Details.Signature = ""
	}
	json := JSON(transaction)
	sum := Checksum256(json)
	transaction.Id = BytesToHex(sum)
	for i, c := range conditions {
		c.Cond.Details.Signature = sigs[i]
	}
	return transaction
}

type Tx struct {
	Asset        *Asset       `json:"asset"`
	Conditions   Conditions   `json:"conditions"`
	Fulfillments Fulfillments `json:"fulfillments"`
	Metadata     *Metadata    `json:"metadata"`
	Operation    Operation    `json:"operation"`
}

func NewTx(asset *Asset, conditions Conditions, fulfillments Fulfillments, meta *Metadata, op Operation) *Tx {
	return &Tx{
		Asset:        asset,
		Conditions:   conditions,
		Fulfillments: fulfillments,
		Metadata:     meta,
		Operation:    op,
	}
}

func NewCreateTx(asset *Asset, conditions Conditions, fulfillments Fulfillments, meta *Metadata) *Tx {
	return NewTx(asset, conditions, fulfillments, meta, CREATE)
}

type Asset struct {
	Data       map[string]interface{} `json:"data"` //--> coalaip model
	Divisible  bool                   `json:"divisible"`
	Id         string                 `json:"id"`
	Refillable bool                   `json:"refillable"`
	Updatable  bool                   `json:"updatable"`
}

func NewAsset(data map[string]interface{}, divisible, refillable, updatable bool) *Asset {
	id := Uuid4()
	return &Asset{
		Data:       data,
		Divisible:  divisible,
		Id:         id,
		Refillable: refillable,
		Updatable:  updatable,
	}
}

type Condition struct {
	Amount      int          `json:"amount"`
	CID         int          `json:"cid"`
	Cond        *Cond        `json:"condition"`
	OwnersAfter []*PublicKey `json:"owners_after"`
}

type Conditions []*Condition

func NewCondition(amount, cid int, details *Details, ownersAfter []*PublicKey) *Condition {
	sig := details.Signature
	details.Signature = ""
	json := JSON(details)
	sum := Checksum256(json)
	b64 := Base64RawURL(sum)
	uri := fmt.Sprintf("cc:%x:%x:%s:%d", TYPE_ID, BITMASK, b64, FULFILLMENT_LENGTH)
	details.Signature = sig
	return &Condition{
		Amount: amount,
		CID:    cid,
		Cond: &Cond{
			Uri:     uri,
			Details: details,
		},
		OwnersAfter: ownersAfter,
	}
}

type Cond struct {
	Uri     string   `json:"uri"`
	Details *Details `json:"details"`
}

type Details struct {
	Bitmask   int        `json:"bitmask"`
	PublicKey *PublicKey `json:"public_key"`
	Signature string     `json:"signature"`
	Type      string     `json:"type"`
	TypeId    int        `json:"type_id"`
}

func NewDetails(pub *PublicKey) *Details {
	return &Details{
		Bitmask:   BITMASK,
		PublicKey: pub,
		Type:      TYPE,
		TypeId:    TYPE_ID,
	}
}

type Fulfillment struct {
	FID          int                    `json:"fid"`
	Fulfill      map[string]interface{} `json:"fulfillment"`
	Input        map[string]interface{} `json:"input"`
	OwnersBefore []*PublicKey           `json:"owners_before"`
}

type Fulfillments []*Fulfillment

func NewFulfillment(fid int, ownersBefore []*PublicKey) *Fulfillment {
	return &Fulfillment{
		FID:          fid,
		OwnersBefore: ownersBefore,
	}
}

type Metadata struct {
	Data map[string]interface{} `json:"data"`
	Id   string                 `json:"id"`
}

func NewMetadata(data map[string]interface{}) *Metadata {
	id := Uuid4()
	return &Metadata{
		Data: data,
		Id:   id,
	}
}
