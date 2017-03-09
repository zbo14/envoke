package schema

import (
	jsonschema "github.com/xeipuuv/gojsonschema"
	"github.com/zbo14/envoke/bigchain"
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/crypto"
)

const (
	COMPOSITION                = ""
	COMPOSITION_RIGHT          = ""
	COMPOSITION_RIGHT_TRANSFER = ""
	MASTER_LICENSE             = ""
	MECHANICAL_LICENSE         = ""
	PUBLICATION                = ""
	RECORDING                  = ""
	RECORDING_RIGHT            = ""
	RECORDING_RIGHT_TRANSFER   = ""
	RELEASE                    = ""
)

func ValidateModel(model Data, source string) (bool, error) {
	schemaLoader := jsonschema.NewReferenceLoader(source)
	goLoader := jsonschema.NewGoLoader(model)
	result, err := jsonschema.Validate(schemaLoader, goLoader)
	if err != nil {
		return false, err
	}
	return result.Valid(), nil
}

func QueryAndValidateModel(id, source string) (Data, error) {
	tx, err := bigchain.GetTx(id)
	if err != nil {
		return nil, err
	}
	model := bigchain.GetTxData(tx)
	success, err := ValidateModel(model, source)
	if err != nil {
		return nil, err
	}
	if !success {
		return nil, Error("Validation failed")
	}
	return tx, nil
}

// Linked-Data

func ValidateComposition(compositionId string) (Data, error) {
	tx, err := QueryAndValidateModel(compositionId, COMPOSITION)
	if err != nil {
		return nil, err
	}
	composition := bigchain.GetTxData(tx)
	composerId := GetComposerId(composition)
	if _, err = QueryAndValidateModel(composerId, ""); err != nil {
		return nil, err
	}
	proId := GetProId(composition)
	if _, err = QueryAndValidateModel(proId, ""); err != nil {
		return nil, err
	}
	return composition, nil
}

func ValidateCompositionRight(compositionRightId string) (Data, crypto.PublicKey, crypto.PublicKey, error) {
	tx, err := QueryAndValidateModel(compositionRightId, COMPOSITION_RIGHT)
	if err != nil {
		return nil, nil, nil, err
	}
	compositionRight := bigchain.GetTxData(tx)
	recipientId := GetRecipientId(compositionRight)
	recipientPub := bigchain.DefaultGetTxRecipient(tx)
	recipientShares := bigchain.GetTxShares(tx)
	senderId := GetSenderId(compositionRight)
	senderPub := bigchain.DefaultGetTxSender(tx)
	tx, err = QueryAndValidateModel(recipientId, "")
	if err != nil {
		return nil, nil, nil, err
	}
	if !recipientPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, nil, nil, ErrorAppend(ErrInvalidKey, recipientPub.String())
	}
	tx, err = QueryAndValidateModel(senderId, "")
	if err != nil {
		return nil, nil, nil, err
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, nil, nil, ErrorAppend(ErrInvalidKey, senderPub.String())
	}
	compositionRight.Set("recipientShares", recipientShares)
	return compositionRight, recipientPub, senderPub, nil
}

func ValidatePublication(publicationId string) (Data, []Data, error) {
	tx, err := QueryAndValidateModel(publicationId, PUBLICATION)
	if err != nil {
		return nil, nil, err
	}
	publication := bigchain.GetTxData(tx)
	senderPub := bigchain.DefaultGetTxSender(tx)
	compositionId := GetCompositionId(publication)
	composition, err := ValidateComposition(compositionId)
	if err != nil {
		return nil, nil, err
	}
	composerId := GetComposerId(composition)
	compositionRightIds := GetCompositionRightIds(publication)
	compositionRights := make([]Data, len(compositionRightIds))
	publisherId := GetPublisherId(publication)
	recipientIds := make(map[string]struct{})
	rightHolder := false
	totalShares := 0
	for i, compositionRightId := range compositionRightIds {
		compositionRight, recipientPub, _, err := ValidateCompositionRight(compositionRightId)
		if err != nil {
			return nil, nil, err
		}
		if composerId != GetSenderId(compositionRight) {
			return nil, nil, ErrorAppend(ErrCriteriaNotMet, "composer must be right sender")
		}
		if compositionId != GetCompositionId(compositionRight) {
			return nil, nil, ErrorAppend(ErrInvalidId, "right links to wrong composition")
		}
		recipientId := GetRecipientId(compositionRight)
		if _, ok := recipientIds[recipientId]; ok {
			return nil, nil, ErrorAppend(ErrCriteriaNotMet, "recipient cannot hold multiple composition rights")
		}
		if !rightHolder && publisherId == recipientId {
			if !senderPub.Equals(recipientPub) {
				return nil, nil, ErrorAppend(ErrCriteriaNotMet, "publisher is not publication sender")
			}
			rightHolder = true
		}
		recipientIds[recipientId] = struct{}{}
		shares := GetRecipientShares(compositionRight)
		if shares <= 0 {
			return nil, nil, ErrorAppend(ErrCriteriaNotMet, "percentage shares must be greater than 0")
		}
		if totalShares += shares; totalShares > 100 {
			return nil, nil, ErrorAppend(ErrCriteriaNotMet, "total percentage shares cannot exceed 100")
		}
		compositionRight.Set("id", compositionRightId)
		compositionRights[i] = compositionRight
	}
	if !rightHolder {
		return nil, nil, ErrorAppend(ErrCriteriaNotMet, "publisher is not composition right holder")
	}
	if totalShares != 100 {
		return nil, nil, ErrorAppend(ErrCriteriaNotMet, "total percentage shares do not equal 100")
	}
	return publication, compositionRights, nil
}

func ValidateCompositionRightTransfer(compositionRightTransferId string) (Data, error) {
	tx, err := QueryAndValidateModel(compositionRightTransferId, COMPOSITION_RIGHT_TRANSFER)
	if err != nil {
		return nil, err
	}
	compositionRightTransfer := bigchain.GetTxData(tx)
	senderPub := bigchain.DefaultGetTxSender(tx)
	recipientId := GetRecipientId(compositionRightTransfer)
	tx, err = QueryAndValidateModel(recipientId, "")
	if err != nil {
		return nil, err
	}
	recipientPub := bigchain.DefaultGetTxSender(tx)
	if senderPub.Equals(recipientPub) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "recipient and sender keys must be different")
	}
	senderId := GetSenderId(compositionRightTransfer)
	tx, err = QueryAndValidateModel(senderId, "")
	if err != nil {
		return nil, err
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, ErrorAppend(ErrInvalidKey, senderPub.String())
	}
	publicationId := GetPublicationId(compositionRightTransfer)
	_, compositionRights, err := ValidatePublication(publicationId)
	if err != nil {
		return nil, err
	}
	txId := GetTxId(compositionRightTransfer)
	txTransfer, err := bigchain.GetTx(txId)
	if err != nil {
		return nil, err
	}
	if bigchain.TRANSFER != bigchain.GetTxOperation(txTransfer) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "expected TRANSFER tx")
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(txTransfer)) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "sender is not signer of TRANSFER tx")
	}
	n := len(bigchain.GetTxOutputs(txTransfer))
	if n != 1 && n != 2 {
		return nil, ErrorAppend(ErrInvalidSize, "tx outputs must have size 1 or 2")
	}
	if !recipientPub.Equals(bigchain.GetTxRecipient(txTransfer, 0)) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "recipient does not hold primary output of TRANSFER tx")
	}
	if n == 2 {
		if !senderPub.Equals(bigchain.GetTxRecipient(txTransfer, 1)) {
			return nil, ErrorAppend(ErrCriteriaNotMet, "sender does not hold secondary output of TRANSFER tx")
		}
	}
	found := false
	compositionRightId := GetCompositionRightId(compositionRightTransfer)
	if compositionRightId != bigchain.GetTxAssetId(txTransfer) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "TRANSFER tx does not link to correct composition right")
	}
	for _, compositionRight := range compositionRights {
		if compositionRightId == bigchain.GetId(compositionRight) {
			found = true
			break
		}
	}
	if !found {
		return nil, ErrorAppend(ErrCriteriaNotMet, "publication does not link to composition right")
	}
	compositionRightTransfer.Set("recipientShares", bigchain.GetTxOutputAmount(txTransfer, 0))
	if n == 2 {
		compositionRightTransfer.Set("senderShares", bigchain.GetTxOutputAmount(txTransfer, 1))
	}
	return compositionRightTransfer, nil
}

func ValidateMechanicalLicense(mechanicalLicenseId string) (Data, error) {
	tx, err := QueryAndValidateModel(mechanicalLicenseId, MECHANICAL_LICENSE)
	mechanicalLicense := bigchain.GetTxData(tx)
	senderPub := bigchain.DefaultGetTxSender(tx)
	senderId := GetSenderId(mechanicalLicense)
	tx, err = bigchain.GetTx(senderId)
	if err != nil {
		return nil, err
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, ErrorAppend(ErrInvalidKey, senderPub.String())
	}
	publicationId := GetPublicationId(mechanicalLicense)
	_, compositionRights, err := ValidatePublication(publicationId)
	if err != nil {
		return nil, err
	}
	compositionRightId := GetCompositionRightId(mechanicalLicense)
	licenseTerritory := GetTerritory(mechanicalLicense)
	compositionRightTransferHolder := false
	if EmptyStr(compositionRightId) {
		compositionRightTransferId := GetCompositionRightTransferId(mechanicalLicense)
		compositionRightTransfer, err := ValidateCompositionRightTransfer(compositionRightTransferId)
		if err != nil {
			return nil, err
		}
		if publicationId != GetPublicationId(compositionRightTransfer) {
			return nil, ErrorAppend(ErrCriteriaNotMet, "compositionRightTransfer links to wrong publication")
		}
		if senderId == GetRecipientId(compositionRightTransfer) {
			//..
		} else if senderId == GetSenderId(compositionRightTransfer) {
			if GetSenderShares(compositionRightTransfer) == 0 {
				return nil, ErrorAppend(ErrCriteriaNotMet, "sender does not have shares in compositionRightTransfer")
			}
		} else {
			return nil, ErrorAppend(ErrCriteriaNotMet, "sender does not have compositionRightTransfer")
		}
		compositionRightId = GetCompositionRightId(compositionRightTransfer)
		compositionRightTransferHolder = true
	}
	var compositionRight Data = nil
	for _, right := range compositionRights {
		if compositionRightId == bigchain.GetId(right) {
			if !compositionRightTransferHolder {
				if senderId != GetRecipientId(right) {
					return nil, ErrorAppend(ErrCriteriaNotMet, "sender does not hold composition right")
				}
			}
			compositionRight = right
			break
		}
	}
	if compositionRight == nil {
		return nil, ErrorAppend(ErrCriteriaNotMet, "could not find composition right")
	}
	rightTerritory := GetTerritory(compositionRight)
OUTER:
	for i := range licenseTerritory {
		for j := range rightTerritory {
			if licenseTerritory[i] == rightTerritory[j] {
				rightTerritory = append(rightTerritory[:j], rightTerritory[j+1:]...)
				continue OUTER
			}
		}
		return nil, ErrorAppend(ErrCriteriaNotMet, "license territory not part of right territory")
	}
	recipientId := GetRecipientId(mechanicalLicense)
	if _, err = QueryAndValidateModel(recipientId, ""); err != nil {
		return nil, err
	}
	return mechanicalLicense, nil
}
