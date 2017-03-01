package bigchain

import (
	"bytes"
	. "github.com/zbo14/envoke/common"
	conds "github.com/zbo14/envoke/crypto/conditions"
	"github.com/zbo14/envoke/crypto/crypto"
	"github.com/zbo14/envoke/crypto/ed25519"
)

const (
	BIGCHAIN_ENDPOINT = ""
	IPDB_ENDPOINT     = "http://cochoa.ipdb.foundation:9984/api/v1/"
	ENDPOINT          = IPDB_ENDPOINT
)

// GET

func GetTx(txId string) (Data, error) {
	url := ENDPOINT + "transactions/" + txId
	response, err := HttpGet(url)
	if err != nil {
		return nil, err
	}
	tx := make(Data)
	if err = ReadJSON(response.Body, &tx); err != nil {
		return nil, err
	}
	return tx, nil
}

// POST

// BigchainDB transaction type
// docs.bigchaindb.com/projects/py-driver/en/latest/handcraft.html

func PostTx(tx Data) (string, error) {
	url := ENDPOINT + "transactions/"
	buf := new(bytes.Buffer)
	buf.Write(MustMarshalJSON(tx))
	response, err := HttpPost(url, "application/json", buf)
	if err != nil {
		return "", err
	}
	data := make(Data)
	if err := ReadJSON(response.Body, &data); err != nil {
		return "", err
	}
	return GetTxId(tx), nil
}

const (
	CREATE   = "CREATE"
	GENESIS  = "GENESIS"
	TRANSFER = "TRANSFER"
	VERSION  = "0.9"
)

func IndividualCreateTx(data Data, owner crypto.PublicKey) Data {
	amounts := []int{1}
	asset := Data{"data": data}
	fulfills := []Data{nil}
	owners := [][]crypto.PublicKey{[]crypto.PublicKey{owner}}
	return CreateTx(amounts, asset, fulfills, owners, owners)
}

func MultipleOwnersCreateTx(data Data, ownersAfter []crypto.PublicKey, ownerBefore crypto.PublicKey) Data {
	amounts := []int{1}
	asset := Data{"data": data}
	fulfills := []Data{nil}
	ownersBefore := []crypto.PublicKey{ownerBefore}
	return CreateTx(amounts, asset, fulfills, [][]crypto.PublicKey{ownersAfter}, [][]crypto.PublicKey{ownersBefore})
}

func IndividualTransferTx(id string, output int, ownerAfter, ownerBefore crypto.PublicKey) Data {
	amounts := []int{1}
	asset := Data{"id": id}
	fulfills := []Data{Data{"txid": id, "output": output}}
	ownersAfter := [][]crypto.PublicKey{[]crypto.PublicKey{ownerAfter}}
	ownersBefore := [][]crypto.PublicKey{[]crypto.PublicKey{ownerBefore}}
	return TransferTx(amounts, asset, fulfills, ownersAfter, ownersBefore)
}

func CreateTx(amounts []int, asset Data, fulfills []Data, ownersAfter, ownersBefore [][]crypto.PublicKey) Data {
	return GenerateTx(amounts, asset, fulfills, nil, CREATE, ownersAfter, ownersBefore)
}

func TransferTx(amounts []int, asset Data, fulfills []Data, ownersAfter, ownersBefore [][]crypto.PublicKey) Data {
	return GenerateTx(amounts, asset, fulfills, nil, TRANSFER, ownersAfter, ownersBefore)
}

func GenerateTx(amounts []int, asset Data, fulfills []Data, metadata Data, operation string, ownersAfter, ownersBefore [][]crypto.PublicKey) Data {
	inputs := NewInputs(fulfills, ownersBefore)
	outputs := NewOutputs(amounts, ownersAfter)
	return NewTx(asset, inputs, metadata, operation, outputs)
}

func NewTx(asset Data, inputs []Data, metadata Data, operation string, outputs []Data) Data {
	tx := Data{
		"asset":     asset,
		"inputs":    inputs,
		"metadata":  metadata,
		"operation": operation,
		"outputs":   outputs,
		"version":   VERSION,
	}
	sum := Checksum256(MustMarshalJSON(tx))
	tx.Set("id", BytesToHex(sum))
	return tx
}

func FulfillTx(tx Data, priv crypto.PrivateKey) {
	json := MustMarshalJSON(tx)
	inputs := tx.Get("inputs").([]Data)
	for _, input := range inputs {
		input.Set("fulfillment", conds.DefaultFulfillmentFromPrivKey(json, priv).String())
	}
}

func FulfilledTx(tx Data) bool {
	var err error
	inputs := tx.GetInterfaceSlice("inputs")
	fulfillments := make([]conds.Fulfillment, len(inputs))
	for i, input := range inputs {
		mapData := AssertMapData(input)
		uri := mapData.GetStr("fulfillment")
		fulfillments[i], err = conds.UnmarshalURI(uri, 1)
		Check(err)
		mapData.Clear("fulfillment")
	}
	fulfilled := true
	json := MustMarshalJSON(tx)
	for _, f := range fulfillments {
		if !f.Validate(json) {
			fulfilled = false
			break
		}
	}
	for i, input := range inputs {
		AssertMapData(input).Set("fulfillment", fulfillments[i])
	}
	return fulfilled
}

// for convenience
func GetTxData(tx Data) Data {
	return tx.GetInnerData("asset", "data")
}

func SetTxData(tx, data Data) {
	tx.SetInnerValue(data, "asset", "data")
}

func GetTxId(tx Data) string {
	return tx.GetStr("id")
}

func GetTxPublicKeys(tx Data) []crypto.PublicKey {
	output := tx.GetInterfaceSlice("outputs")[0]
	condition := AssertMap(output)["condition"]
	details := AssertMap(condition)["details"]
	subs := AssertMapData(details).GetInterfaceSlice("subfulfillments")
	pubs := make([]crypto.PublicKey, len(subs))
	for i, sub := range subs {
		pubstr := AssertMapData(sub).GetStr("public_key")
		pubs[i] = new(ed25519.PublicKey)
		pubs[i].FromString(pubstr)
	}
	return pubs
}

func GetTxPublicKey(tx Data) crypto.PublicKey {
	pub := new(ed25519.PublicKey)
	output := tx.GetInterfaceSlice("outputs")[0]
	condition := AssertMap(output)["condition"]
	details := AssertMap(condition)["details"]
	pubstr := AssertMapData(details).GetStr("public_key")
	pub.FromString(pubstr)
	return pub
}

func NewInputs(fulfills []Data, ownersBefore [][]crypto.PublicKey) []Data {
	n := len(fulfills)
	if n != len(ownersBefore) {
		panic(ErrorAppend(ErrInvalidSize, "slices are different sizes"))
	}
	inputs := make([]Data, n)
	for i := range inputs {
		inputs[i] = NewInput(fulfills[i], ownersBefore[i])
	}
	return inputs
}

func NewInput(fulfills Data, ownersBefore []crypto.PublicKey) Data {
	return Data{
		"fulfillment":   nil,
		"fulfills":      fulfills,
		"owners_before": ownersBefore,
	}
}

func NewOutputs(amounts []int, ownersAfter [][]crypto.PublicKey) []Data {
	n := len(amounts)
	if n != len(ownersAfter) {
		panic(ErrorAppend(ErrInvalidSize, "slices are different sizes"))
	}
	outputs := make([]Data, n)
	for i, owner := range ownersAfter {
		outputs[i] = NewOutput(amounts[i], owner)
	}
	return outputs
}

func NewOutput(amount int, ownersAfter []crypto.PublicKey) Data {
	n := len(ownersAfter)
	if n == 0 {
		return nil
	}
	if n == 1 {
		return Data{
			"amount":      amount,
			"condition":   conds.DefaultFulfillmentFromPubKey(ownersAfter[0]),
			"public_keys": ownersAfter,
		}
	}
	return Data{
		"amount":      amount,
		"condition":   conds.DefaultFulfillmentThresholdFromPubKeys(ownersAfter),
		"public_keys": ownersAfter,
	}
}
