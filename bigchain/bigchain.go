package bigchain

import (
	"bytes"

	. "github.com/zbo14/envoke/common"
	conds "github.com/zbo14/envoke/crypto/conditions"
	"github.com/zbo14/envoke/crypto/crypto"
	"github.com/zbo14/envoke/crypto/ed25519"
)

// GET

func GetTx(txId string) (Data, error) {
	url := Getenv("IPDB_ENDPOINT") + "transactions/" + txId
	response, err := HttpGet(url)
	if err != nil {
		return nil, err
	}
	tx := make(Data)
	if err = ReadJSON(response.Body, &tx); err != nil {
		return nil, err
	}
	if !FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return tx, nil
}

// POST

// BigchainDB transaction type
// docs.bigchaindb.com/projects/py-driver/en/latest/handcraft.html

func PostTx(tx Data) (string, error) {
	url := Getenv("IPDB_ENDPOINT") + "transactions/"
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
	return GetId(tx), nil
}

const (
	CREATE   = "CREATE"
	GENESIS  = "GENESIS"
	TRANSFER = "TRANSFER"
	VERSION  = "0.9"
)

func DefaultIndividualCreateTx(data Data, owner crypto.PublicKey) Data {
	return IndividualCreateTx(1, data, owner, owner)
}

func IndividualCreateTx(amount int, data Data, ownerAfter, ownerBefore crypto.PublicKey) Data {
	amounts := []int{amount}
	asset := Data{"data": data}
	fulfills := []Data{nil}
	ownersAfter := [][]crypto.PublicKey{[]crypto.PublicKey{ownerAfter}}
	ownersBefore := [][]crypto.PublicKey{[]crypto.PublicKey{ownerBefore}}
	return CreateTx(amounts, asset, fulfills, ownersAfter, ownersBefore)
}

func MultipleOwnersCreateTx(amounts []int, data Data, ownersAfter []crypto.PublicKey, ownerBefore crypto.PublicKey) Data {
	asset := Data{"data": data}
	fulfills := []Data{nil}
	ownersBefore := []crypto.PublicKey{ownerBefore}
	n := len(amounts)
	if n == 0 {
		panic(ErrorAppend(ErrCriteriaNotMet, "must have at least one amount"))
	}
	owners := make([][]crypto.PublicKey, n)
	if n == 1 {
		owners[0] = ownersAfter
	} else {
		if n != len(ownersAfter) {
			panic(ErrorAppend(ErrCriteriaNotMet, "must have same number of amounts as owners if number > 1"))
		}
		for i, owner := range ownersAfter {
			owners[i] = []crypto.PublicKey{owner}
		}
	}
	return CreateTx(amounts, asset, fulfills, owners, [][]crypto.PublicKey{ownersBefore})
}

func DefaultIndividualTransferTx(assetId, consumeId string, output int, ownerAfter, ownerBefore crypto.PublicKey) Data {
	return IndividualTransferTx(1, assetId, consumeId, output, ownerAfter, ownerBefore)
}

func IndividualTransferTx(amount int, assetId, consumeId string, output int, ownerAfter, ownerBefore crypto.PublicKey) Data {
	amounts := []int{amount}
	asset := Data{"id": assetId}
	fulfills := []Data{Data{"txid": consumeId, "output": output}}
	ownersAfter := [][]crypto.PublicKey{[]crypto.PublicKey{ownerAfter}}
	ownersBefore := [][]crypto.PublicKey{[]crypto.PublicKey{ownerBefore}}
	return TransferTx(amounts, asset, fulfills, ownersAfter, ownersBefore)
}

func DivisibleTransferTx(amounts []int, assetId, consumeId string, output int, ownersAfter []crypto.PublicKey, ownerBefore crypto.PublicKey) Data {
	n := len(amounts)
	if n <= 1 || n != len(ownersAfter) {
		panic(ErrInvalidSize)
	}
	asset := Data{"id": assetId}
	fulfills := []Data{Data{"txid": consumeId, "output": output}}
	owners := make([][]crypto.PublicKey, len(ownersAfter))
	for i, owner := range ownersAfter {
		owners[i] = []crypto.PublicKey{owner}
	}
	ownersBefore := [][]crypto.PublicKey{[]crypto.PublicKey{ownerBefore}}
	return TransferTx(amounts, asset, fulfills, owners, ownersBefore)
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
func GetId(data Data) string {
	return data.GetStr("id")
}

func GetPublicKey(data Data) crypto.PublicKey {
	pub := new(ed25519.PublicKey)
	pub.FromString(data.GetStr("public_key"))
	return pub
}

func GetTxAsset(tx Data) Data {
	return tx.GetMapData("asset")
}

func GetTxAssetId(tx Data) string {
	return GetId(GetTxAsset(tx))
}

func GetTxData(tx Data) Data {
	return tx.GetInnerData("asset", "data")
}

func SetTxData(tx, data Data) {
	tx.SetInnerValue(data, "asset", "data")
}

func GetTxOperation(tx Data) string {
	return tx.GetStr("operation")
}

func GetTxSigners(tx Data) [][]crypto.PublicKey {
	inputs := GetTxInputs(tx)
	return GetInputsPublicKeys(inputs)
}

func DefaultGetTxSigner(tx Data) crypto.PublicKey {
	return GetTxSigner(tx, 0)
}

func GetTxSigner(tx Data, n int) crypto.PublicKey {
	pubs := GetTxSigners(tx)
	return pubs[n][0]
}

func GetTxRecipients(tx Data) [][]crypto.PublicKey {
	outputs := GetTxOutputs(tx)
	return GetOutputsPublicKeys(outputs)
}

func DefaultGetTxRecipient(tx Data) crypto.PublicKey {
	return GetTxRecipient(tx, 0)
}

func GetTxRecipient(tx Data, n int) crypto.PublicKey {
	pubs := GetTxRecipients(tx)
	return pubs[n][0]
}

func GetTxShares(tx Data) int {
	return GetTxOutputAmount(tx, 0)
}

func GetTxInputs(tx Data) []Data {
	inputs := tx.GetInterfaceSlice("inputs")
	datas := make([]Data, len(inputs))
	for i, input := range inputs {
		datas[i] = AssertMapData(input)
	}
	return datas
}

func GetInputPublicKeys(input Data) []crypto.PublicKey {
	owners := input.GetInterfaceSlice("owners_before")
	pubs := make([]crypto.PublicKey, len(owners))
	for i, owner := range owners {
		pubs[i] = new(ed25519.PublicKey)
		pubs[i].FromString(AssertStr(owner))
	}
	return pubs
}

func GetInputsPublicKeys(inputs []Data) [][]crypto.PublicKey {
	pubs := make([][]crypto.PublicKey, len(inputs))
	for i, input := range inputs {
		pubs[i] = GetInputPublicKeys(input)
	}
	return pubs
}

func GetTxOutputAmount(tx Data, n int) int {
	output := GetTxOutput(tx, n)
	return GetOutputAmount(output)
}

func GetTxOutputs(tx Data) []Data {
	outputs := tx.GetInterfaceSlice("outputs")
	datas := make([]Data, len(outputs))
	for i, output := range outputs {
		datas[i] = AssertMapData(output)
	}
	return datas
}

func GetTxOutput(tx Data, n int) Data {
	outputs := GetTxOutputs(tx)
	return outputs[n]
}

func GetOutputAmount(output Data) int {
	return output.GetInt("amount")
}

func GetOutputCondition(output Data) Data {
	return output.GetMapData("condition")
}

func GetConditionDetails(condition Data) Data {
	return condition.GetMapData("details")
}

func GetDetailsSubfulfillments(details Data) []Data {
	subs := details.GetInterfaceSlice("subfulfillments")
	if subs == nil {
		return nil
	}
	datas := make([]Data, len(subs))
	for i, sub := range subs {
		datas[i] = AssertMapData(sub)
	}
	return datas
}

func GetOutputPublicKeys(output Data) []crypto.PublicKey {
	details := output.GetInnerData("condition", "details")
	subs := GetDetailsSubfulfillments(details)
	if subs == nil {
		pub := GetPublicKey(details)
		return []crypto.PublicKey{pub}
	}
	pubs := make([]crypto.PublicKey, len(subs))
	for i, sub := range subs {
		pubs[i] = GetPublicKey(sub)
	}
	return pubs
}

func GetOutputsPublicKeys(outputs []Data) [][]crypto.PublicKey {
	pubs := make([][]crypto.PublicKey, len(outputs))
	for i, output := range outputs {
		pubs[i] = GetOutputPublicKeys(output)
	}
	return pubs
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
