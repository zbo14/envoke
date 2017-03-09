package schema

import (
	jsonschema "github.com/xeipuuv/gojsonschema"
	"github.com/zbo14/envoke/bigchain"
	. "github.com/zbo14/envoke/common"
)

var schemaLoader = jsonschema.NewReferenceLoader(ENVOKE)

func ValidateModel(model Data) (bool, error) {
	goLoader := jsonschema.NewGoLoader(model)
	result, err := jsonschema.Validate(schemaLoader, goLoader)
	if err != nil {
		return false, err
	}
	return result.Valid(), nil
}

func QueryAndValidateModel(txId string) (bool, error) {
	tx, err := bigchain.GetTx(txId)
	if err != nil {
		return false, err
	}
	model := bigchain.GetTxData(tx)
	return ValidateModel(model)
}

// Linked-Data

func ValidateComposition(composition Data) {}
