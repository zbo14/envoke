package linked_data

import (
	jsonschema "github.com/xeipuuv/gojsonschema"
	"github.com/zbo14/envoke/bigchain"
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/crypto"
	"github.com/zbo14/envoke/spec"
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

func QueryAndValidateModel(id string) (Data, error) {
	tx, err := bigchain.GetTx(id)
	if err != nil {
		return nil, err
	}
	model := bigchain.GetTxData(tx)
	source := spec.GetType(model)
	success, err := ValidateModel(model, source)
	if err != nil {
		return nil, err
	}
	if !success {
		return nil, Error("Validation failed")
	}
	return tx, nil
}

func ValidateComposition(compositionId string) (Data, error) {
	tx, err := QueryAndValidateModel(compositionId)
	if err != nil {
		return nil, err
	}
	composition := bigchain.GetTxData(tx)
	composerId := spec.GetComposerId(composition)
	if _, err = QueryAndValidateModel(composerId); err != nil {
		return nil, err
	}
	proId := spec.GetProId(composition)
	if _, err = QueryAndValidateModel(proId); err != nil {
		return nil, err
	}
	return composition, nil
}

func ValidateCompositionRight(compositionRightId string) (Data, crypto.PublicKey, crypto.PublicKey, error) {
	tx, err := QueryAndValidateModel(compositionRightId)
	if err != nil {
		return nil, nil, nil, err
	}
	compositionRight := bigchain.GetTxData(tx)
	recipientId := spec.GetRecipientId(compositionRight)
	recipientPub := bigchain.DefaultGetTxRecipient(tx)
	recipientShares := bigchain.GetTxShares(tx)
	senderId := spec.GetSenderId(compositionRight)
	senderPub := bigchain.DefaultGetTxSender(tx)
	tx, err = QueryAndValidateModel(recipientId)
	if err != nil {
		return nil, nil, nil, err
	}
	if !recipientPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, nil, nil, ErrorAppend(ErrInvalidKey, recipientPub.String())
	}
	tx, err = QueryAndValidateModel(senderId)
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
	tx, err := QueryAndValidateModel(publicationId)
	if err != nil {
		return nil, nil, err
	}
	publication := bigchain.GetTxData(tx)
	senderPub := bigchain.DefaultGetTxSender(tx)
	compositionId := spec.GetCompositionId(publication)
	composition, err := ValidateComposition(compositionId)
	if err != nil {
		return nil, nil, err
	}
	composerId := spec.GetComposerId(composition)
	compositionRightIds := spec.GetCompositionRightIds(publication)
	compositionRights := make([]Data, len(compositionRightIds))
	publisherId := spec.GetPublisherId(publication)
	recipientIds := make(map[string]struct{})
	rightHolder := false
	totalShares := 0
	for i, compositionRightId := range compositionRightIds {
		compositionRight, recipientPub, _, err := ValidateCompositionRight(compositionRightId)
		if err != nil {
			return nil, nil, err
		}
		if composerId != spec.GetSenderId(compositionRight) {
			return nil, nil, ErrorAppend(ErrCriteriaNotMet, "composer must be right sender")
		}
		if compositionId != spec.GetCompositionId(compositionRight) {
			return nil, nil, ErrorAppend(ErrInvalidId, "right links to wrong composition")
		}
		recipientId := spec.GetRecipientId(compositionRight)
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
		shares := spec.GetRecipientShares(compositionRight)
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
	tx, err := QueryAndValidateModel(compositionRightTransferId)
	if err != nil {
		return nil, err
	}
	compositionRightTransfer := bigchain.GetTxData(tx)
	senderPub := bigchain.DefaultGetTxSender(tx)
	recipientId := spec.GetRecipientId(compositionRightTransfer)
	tx, err = QueryAndValidateModel(recipientId)
	if err != nil {
		return nil, err
	}
	recipientPub := bigchain.DefaultGetTxSender(tx)
	if senderPub.Equals(recipientPub) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "recipient and sender keys must be different")
	}
	senderId := spec.GetSenderId(compositionRightTransfer)
	tx, err = QueryAndValidateModel(senderId)
	if err != nil {
		return nil, err
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, ErrorAppend(ErrInvalidKey, senderPub.String())
	}
	publicationId := spec.GetPublicationId(compositionRightTransfer)
	_, compositionRights, err := ValidatePublication(publicationId)
	if err != nil {
		return nil, err
	}
	txId := spec.GetTxId(compositionRightTransfer)
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
	compositionRightId := spec.GetCompositionRightId(compositionRightTransfer)
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
	tx, err := QueryAndValidateModel(mechanicalLicenseId)
	mechanicalLicense := bigchain.GetTxData(tx)
	senderPub := bigchain.DefaultGetTxSender(tx)
	senderId := spec.GetSenderId(mechanicalLicense)
	tx, err = bigchain.GetTx(senderId)
	if err != nil {
		return nil, err
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, ErrorAppend(ErrInvalidKey, senderPub.String())
	}
	publicationId := spec.GetPublicationId(mechanicalLicense)
	_, compositionRights, err := ValidatePublication(publicationId)
	if err != nil {
		return nil, err
	}
	compositionRightId := spec.GetCompositionRightId(mechanicalLicense)
	licenseTerritory := spec.GetTerritory(mechanicalLicense)
	compositionRightTransferHolder := false
	if EmptyStr(compositionRightId) {
		compositionRightTransferId := spec.GetCompositionRightTransferId(mechanicalLicense)
		compositionRightTransfer, err := ValidateCompositionRightTransfer(compositionRightTransferId)
		if err != nil {
			return nil, err
		}
		if publicationId != spec.GetPublicationId(compositionRightTransfer) {
			return nil, ErrorAppend(ErrCriteriaNotMet, "compositionRightTransfer links to wrong publication")
		}
		if senderId == spec.GetRecipientId(compositionRightTransfer) {
			//..
		} else if senderId == spec.GetSenderId(compositionRightTransfer) {
			if spec.GetSenderShares(compositionRightTransfer) == 0 {
				return nil, ErrorAppend(ErrCriteriaNotMet, "sender does not have shares in compositionRightTransfer")
			}
		} else {
			return nil, ErrorAppend(ErrCriteriaNotMet, "sender does not have compositionRightTransfer")
		}
		compositionRightId = spec.GetCompositionRightId(compositionRightTransfer)
		compositionRightTransferHolder = true
	}
	var compositionRight Data = nil
	for _, right := range compositionRights {
		if compositionRightId == bigchain.GetId(right) {
			if !compositionRightTransferHolder {
				if senderId != spec.GetRecipientId(right) {
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
	rightTerritory := spec.GetTerritory(compositionRight)
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
	recipientId := spec.GetRecipientId(mechanicalLicense)
	if _, err = QueryAndValidateModel(recipientId); err != nil {
		return nil, err
	}
	return mechanicalLicense, nil
}

func ValidateRecording(recordingId string) (Data, string, error) {
	tx, err := QueryAndValidateModel(recordingId)
	if err != nil {
		return nil, "", err
	}
	recording := bigchain.GetTxData(tx)
	var senderId string
	senderPub := bigchain.DefaultGetTxSender(tx)
	performerId := spec.GetPerformerId(recording)
	tx, err = QueryAndValidateModel(performerId)
	if err != nil {
		return nil, "", err
	}
	if senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		senderId = performerId
	}
	producerId := spec.GetProducerId(recording)
	tx, err = QueryAndValidateModel(producerId)
	if err != nil {
		return nil, "", err
	}
	if EmptyStr(senderId) {
		if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
			return nil, "", ErrorAppend(ErrCriteriaNotMet, "performer or producer must be sender")
		}
		senderId = producerId
	}
	publicationId := spec.GetPublicationId(recording)
	_, compositionRights, err := ValidatePublication(publicationId)
	if err != nil {
		return nil, "", err
	}
	compositionRightId := spec.GetCompositionRightId(recording)
	if !EmptyStr(compositionRightId) {
		rightHolder := false
		for _, compositionRight := range compositionRights {
			if compositionRightId == spec.GetId(compositionRight) {
				if senderId != spec.GetRecipientId(compositionRight) {
					return nil, "", ErrorAppend(ErrCriteriaNotMet, "sender does not hold composition right")
				}
				rightHolder = true
				break
			}
		}
		if !rightHolder {
			return nil, "", ErrorAppend(ErrCriteriaNotMet, "sender does not hold composition right")
		}
	}
	return recording, senderId, nil
}

func ValidateRecordingRight(recordingRightId string) (Data, crypto.PublicKey, crypto.PublicKey, error) {
	tx, err := QueryAndValidateModel(recordingRightId)
	if err != nil {
		return nil, nil, nil, err
	}
	recordingRight := bigchain.GetTxData(tx)
	recipientId := spec.GetRecipientId(recordingRight)
	recipientPub := bigchain.DefaultGetTxRecipient(tx)
	recipientShares := bigchain.GetTxShares(tx)
	senderPub := bigchain.DefaultGetTxSender(tx)
	senderId := spec.GetSenderId(recordingRight)
	tx, err = QueryAndValidateModel(senderId)
	if err != nil {
		return nil, nil, nil, err
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, nil, nil, ErrorAppend(ErrInvalidKey, senderPub.String())
	}
	tx, err = bigchain.GetTx(recipientId)
	if err != nil {
		return nil, nil, nil, err
	}
	if !recipientPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, nil, nil, ErrorAppend(ErrInvalidKey, recipientPub.String())
	}
	recordingRight.Set("recipientShares", recipientShares)
	return recordingRight, recipientPub, senderPub, nil
}

func ValidateRelease(releaseId string) (Data, []Data, error) {
	tx, err := QueryAndValidateModel(releaseId)
	if err != nil {
		return nil, nil, err
	}
	release := bigchain.GetTxData(tx)
	senderPub := bigchain.DefaultGetTxSender(tx)
	recordLabelId := spec.GetRecordLabelId(release)
	recordingId := spec.GetRecordingId(release)
	recording, senderId, err := ValidateRecording(recordingId)
	if err != nil {
		return nil, nil, err
	}
	publicationId := spec.GetPublicationId(recording)
	compositionRightId := spec.GetCompositionRightId(recording)
	if EmptyStr(compositionRightId) {
		mechanicalLicenseId := spec.GetMechanicalLicenseId(release)
		mechanicalLicense, err := ValidateMechanicalLicense(mechanicalLicenseId)
		if err != nil {
			return nil, nil, err
		}
		if publicationId != spec.GetPublicationId(mechanicalLicense) {
			return nil, nil, ErrorAppend(ErrCriteriaNotMet, "mechanical license does not link to publication")
		}
		if recordLabelId != spec.GetRecipientId(mechanicalLicense) {
			return nil, nil, ErrorAppend(ErrCriteriaNotMet, "record label must be license holder")
		}
	}
	recipientIds := make(map[string]struct{})
	recordingRightIds := spec.GetRecordingRightIds(release)
	recordingRights := make([]Data, len(recordingRightIds))
	rightHolder := false
	totalShares := 0
	for i, rightId := range recordingRightIds {
		recordingRight, recipientPub, _, err := ValidateRecordingRight(rightId)
		if err != nil {
			return nil, nil, err
		}
		if recordingId != spec.GetRecordingId(recordingRight) {
			return nil, nil, ErrorAppend(ErrInvalidId, "right links to wrong recording")
		}
		if senderId != spec.GetSenderId(recordingRight) {
			return nil, nil, ErrorAppend(ErrCriteriaNotMet, "recording sender must be right sender")
		}
		recipientId := spec.GetRecipientId(recordingRight)
		if _, ok := recipientIds[recipientId]; ok {
			return nil, nil, ErrorAppend(ErrCriteriaNotMet, "holder cannot have multiple assignments")
		}
		if !rightHolder && recordLabelId == recipientId {
			if !senderPub.Equals(recipientPub) {
				return nil, nil, ErrorAppend(ErrCriteriaNotMet, "record label is not release sender")
			}
			rightHolder = true
		}
		recipientIds[recipientId] = struct{}{}
		shares := spec.GetRecipientShares(recordingRight)
		if shares <= 0 {
			return nil, nil, ErrorAppend(ErrCriteriaNotMet, "percentage shares must be greater than 0")
		}
		if totalShares += shares; totalShares > 100 {
			return nil, nil, ErrorAppend(ErrCriteriaNotMet, "total percentage shares cannot exceed 100")
		}
		recordingRight.Set("id", rightId)
		recordingRights[i] = recordingRight
	}
	if !rightHolder {
		return nil, nil, ErrorAppend(ErrCriteriaNotMet, "record label is not recording right holder")
	}
	if totalShares != 100 {
		return nil, nil, ErrorAppend(ErrCriteriaNotMet, "total percentage shares do not equal 100")
	}
	return release, recordingRights, nil
}

func ValidateRecordingRightTransfer(recordingRightTransferId string) (Data, error) {
	tx, err := QueryAndValidateModel(recordingRightTransferId)
	if err != nil {
		return nil, err
	}
	recordingRightTransfer := bigchain.GetTxData(tx)
	senderPub := bigchain.DefaultGetTxSender(tx)
	recipientId := spec.GetRecipientId(recordingRightTransfer)
	tx, err = QueryAndValidateModel(recipientId)
	if err != nil {
		return nil, err
	}
	recipientPub := bigchain.DefaultGetTxSender(tx)
	if senderPub.Equals(recipientPub) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "recipient and sender keys must be different")
	}
	senderId := spec.GetSenderId(recordingRightTransfer)
	tx, err = QueryAndValidateModel(senderId)
	if err != nil {
		return nil, err
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, ErrorAppend(ErrInvalidKey, senderPub.String())
	}
	releaseId := spec.GetReleaseId(recordingRightTransfer)
	_, recordingRights, err := ValidateRelease(releaseId)
	if err != nil {
		return nil, err
	}
	txId := spec.GetTxId(recordingRightTransfer)
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
	recordingRightId := spec.GetRecordingRightId(recordingRightTransfer)
	if recordingRightId != bigchain.GetTxAssetId(txTransfer) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "TRANSFER tx does not link to correct recording right")
	}
	for _, recordingRight := range recordingRights {
		if recordingRightId == bigchain.GetId(recordingRight) {
			found = true
			break
		}
	}
	if !found {
		return nil, ErrorAppend(ErrCriteriaNotMet, "release does not link to recording right")
	}
	recordingRightTransfer.Set("recipientShares", bigchain.GetTxOutputAmount(txTransfer, 0))
	if n == 2 {
		recordingRightTransfer.Set("senderShares", bigchain.GetTxOutputAmount(txTransfer, 1))
	}
	return recordingRightTransfer, nil
}

func ValidateMasterLicense(masterLicenseId string) (Data, error) {
	tx, err := QueryAndValidateModel(masterLicenseId)
	if err != nil {
		return nil, err
	}
	masterLicense := bigchain.GetTxData(tx)
	senderPub := bigchain.DefaultGetTxSender(tx)
	senderId := spec.GetSenderId(masterLicense)
	tx, err = QueryAndValidateModel(senderId)
	if err != nil {
		return nil, err
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, ErrorAppend(ErrInvalidKey, senderPub.String())
	}
	releaseId := spec.GetReleaseId(masterLicense)
	_, recordingRights, err := ValidateRelease(releaseId)
	if err != nil {
		return nil, err
	}
	recordingRightId := spec.GetRecordingRightId(masterLicense)
	licenseTerritory := spec.GetTerritory(masterLicense)
	transferHolder := false
	if EmptyStr(recordingRightId) {
		transferId := spec.GetRecordingRightTransferId(masterLicense)
		transfer, err := ValidateRecordingRightTransfer(transferId)
		if err != nil {
			return nil, err
		}
		if releaseId != spec.GetReleaseId(transfer) {
			return nil, ErrorAppend(ErrCriteriaNotMet, "transfer links to wrong release")
		}
		if senderId == spec.GetRecipientId(transfer) {
			//..
		} else if senderId == spec.GetSenderId(transfer) {
			if spec.GetSenderShares(transfer) == 0 {
				return nil, ErrorAppend(ErrCriteriaNotMet, "sender does not have shares in transfer")
			}
		} else {
			return nil, ErrorAppend(ErrCriteriaNotMet, "sender does not have transfer")
		}
		recordingRightId = spec.GetRecordingRightId(transfer)
		transferHolder = true
	}
	var recordingRight Data = nil
	for _, right := range recordingRights {
		if recordingRightId == spec.GetId(right) {
			if !transferHolder {
				if senderId != spec.GetRecipientId(right) {
					return nil, ErrorAppend(ErrCriteriaNotMet, "sender does not hold recording right")
				}
			}
			recordingRight = right
			break
		}
	}
	if recordingRight == nil {
		return nil, ErrorAppend(ErrCriteriaNotMet, "could not find recording right")
	}
	rightTerritory := spec.GetTerritory(recordingRight)
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
	recipientId := spec.GetRecipientId(masterLicense)
	tx, err = QueryAndValidateModel(recipientId)
	if err != nil {
		return nil, err
	}
	return masterLicense, nil
}
