package linked_data

import (
	"github.com/zbo14/envoke/bigchain"
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/crypto"
	"github.com/zbo14/envoke/spec"
)

// Composition

func ValidateCompositionById(compositionId string) (Data, error) {
	tx, err := bigchain.GetTx(compositionId)
	if err != nil {
		return nil, err
	}
	composition := bigchain.GetTxData(tx)
	if err = spec.ValidComposition(composition); err != nil {
		return nil, err
	}
	pub := bigchain.DefaultGetTxSender(tx)
	if err = ValidateComposition(composition, pub); err != nil {
		return nil, err
	}
	return composition, nil
}

func ValidateComposition(composition Data, pub crypto.PublicKey) error {
	composerId := spec.GetComposerId(composition)
	tx, err := bigchain.GetTx(composerId)
	if err != nil {
		return err
	}
	if !pub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return ErrorAppend(ErrCriteriaNotMet, "composition must be signed by composer")
	}
	composer := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(composer); err != nil {
		return err
	}
	return nil
}

func GetComposer(data Data) (Data, error) {
	composerId := spec.GetComposerId(data)
	tx, err := bigchain.GetTx(composerId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

// Publication

func ValidatePublication(publicationId string) (Data, []Data, error) {
	tx, err := bigchain.GetTx(publicationId)
	if err != nil {
		return nil, nil, err
	}
	publication := bigchain.GetTxData(tx)
	if err = spec.ValidPublication(publication); err != nil {
		return nil, nil, err
	}
	senderPub := bigchain.DefaultGetTxSender(tx)
	compositionId := spec.GetCompositionId(publication)
	composition, err := ValidateCompositionById(compositionId)
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
	for i, rightId := range compositionRightIds {
		compositionRight, recipientPub, _, err := ValidateCompositionRight(rightId)
		if err != nil {
			return nil, nil, err
		}
		if compositionId != spec.GetCompositionId(compositionRight) {
			return nil, nil, ErrorAppend(ErrInvalidId, "right links to wrong composition")
		}
		if composerId != spec.GetSenderId(compositionRight) {
			return nil, nil, ErrorAppend(ErrCriteriaNotMet, "composer must be right sender")
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
		compositionRight.Set("id", rightId)
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

func QueryPublicationField(field string, publicationId string) (interface{}, error) {
	publication, compositionRights, err := ValidatePublication(publicationId)
	if err != nil {
		return nil, err
	}
	switch field {
	case "composer":
		composition, err := GetComposition(publication)
		if err != nil {
			return nil, err
		}
		return GetComposer(composition)
	case "composition":
		return GetComposition(publication)
	case "composition_rights":
		return compositionRights, nil
	case "publisher":
		composition, err := GetComposition(publication)
		if err != nil {
			return nil, err
		}
		return GetPublisher(composition)
	default:
		return nil, ErrorAppend(ErrInvalidField, field)
	}
}

func GetComposition(data Data) (Data, error) {
	compositionId := spec.GetCompositionId(data)
	tx, err := bigchain.GetTx(compositionId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetPublisher(data Data) (Data, error) {
	publisherId := spec.GetPublisherId(data)
	tx, err := bigchain.GetTx(publisherId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

// Recording

func ValidateRecording(recordingId string) (Data, string, error) {
	tx, err := bigchain.GetTx(recordingId)
	if err != nil {
		return nil, "", err
	}
	recording := bigchain.GetTxData(tx)
	if err := spec.ValidRecording(recording); err != nil {
		return nil, "", err
	}
	senderPub := bigchain.DefaultGetTxSender(tx)
	performerId := spec.GetPerformerId(recording)
	tx, err = bigchain.GetTx(performerId)
	if err != nil {
		return nil, "", err
	}
	performer := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(performer); err != nil {
		return nil, "", err
	}
	var senderId string
	if senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		senderId = performerId
	}
	producerId := spec.GetProducerId(recording)
	tx, err = bigchain.GetTx(producerId)
	if err != nil {
		return nil, "", err
	}
	producer := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(producer); err != nil {
		return nil, "", err
	}
	if EmptyStr(senderId) && senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		senderId = producerId
	}
	if EmptyStr(senderId) {
		return nil, "", ErrorAppend(ErrCriteriaNotMet, "performer or producer must be sender")
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
			if compositionRightId == bigchain.GetId(compositionRight) {
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

func GetCompositionRight(data Data) (Data, error) {
	compositionRightId := spec.GetCompositionRightId(data)
	tx, err := bigchain.GetTx(compositionRightId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetPerformer(data Data) (Data, error) {
	performerId := spec.GetPerformerId(data)
	tx, err := bigchain.GetTx(performerId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetProducer(data Data) (Data, error) {
	producerId := spec.GetProducerId(data)
	tx, err := bigchain.GetTx(producerId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetPublication(data Data) (Data, error) {
	publicationId := spec.GetPublicationId(data)
	tx, err := bigchain.GetTx(publicationId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func ValidateRelease(releaseId string) (Data, []Data, error) {
	tx, err := bigchain.GetTx(releaseId)
	if err != nil {
		return nil, nil, err
	}
	release := bigchain.GetTxData(tx)
	if err = spec.ValidRelease(release); err != nil {
		return nil, nil, err
	}
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

func QueryReleaseField(field string, releaseId string) (interface{}, error) {
	release, recordingRights, err := ValidateRelease(releaseId)
	if err != nil {
		return nil, err
	}
	switch field {
	case "composition_right":
		recording, err := GetRecording(release)
		if err != nil {
			return nil, err
		}
		return GetCompositionRight(recording)
	case "mechanical_license":
		return GetMechanicalLicense(release)
	case "performer":
		recording, err := GetRecording(release)
		if err != nil {
			return nil, err
		}
		return GetPerformer(recording)
	case "producer":
		recording, err := GetRecording(release)
		if err != nil {
			return nil, err
		}
		return GetProducer(recording)
	case "publication":
		recording, err := GetRecording(release)
		if err != nil {
			return nil, err
		}
		return GetPublication(recording)
	case "recording":
		return GetRecording(release)
	case "recordingRights":
		return recordingRights, nil
	case "record_label":
		recording, err := GetRecording(release)
		if err != nil {
			return nil, err
		}
		return GetRecordLabel(recording)
	default:
		return nil, ErrorAppend(ErrInvalidField, field)
	}
}

func GetRecordLabel(data Data) (Data, error) {
	recordLabelId := spec.GetRecordLabelId(data)
	tx, err := bigchain.GetTx(recordLabelId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetMechanicalLicense(data Data) (Data, error) {
	mechanicalLicenseId := spec.GetMechanicalLicenseId(data)
	tx, err := bigchain.GetTx(mechanicalLicenseId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetRecording(data Data) (Data, error) {
	recordingId := spec.GetRecordingId(data)
	tx, err := bigchain.GetTx(recordingId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

// Mechanical License

func ValidateMechanicalLicense(mechanicalLicenseId string) (Data, error) {
	tx, err := bigchain.GetTx(mechanicalLicenseId)
	if err != nil {
		return nil, err
	}
	mechanicalLicense := bigchain.GetTxData(tx)
	if err = spec.ValidMechanicalLicense(mechanicalLicense); err != nil {
		return nil, err
	}
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
	transferHolder := false
	if EmptyStr(compositionRightId) {
		transferId := spec.GetCompositionRightTransferId(mechanicalLicense)
		transfer, err := ValidateCompositionRightTransfer(transferId)
		if err != nil {
			return nil, err
		}
		if publicationId != spec.GetPublicationId(transfer) {
			return nil, ErrorAppend(ErrCriteriaNotMet, "transfer links to wrong publication")
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
		compositionRightId = spec.GetCompositionRightId(transfer)
		transferHolder = true
	}
	var compositionRight Data = nil
	for _, right := range compositionRights {
		if compositionRightId == bigchain.GetId(right) {
			if !transferHolder {
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
	tx, err = bigchain.GetTx(recipientId)
	if err != nil {
		return nil, err
	}
	recipient := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(recipient); err != nil {
		return nil, err
	}
	return mechanicalLicense, nil
}

func QueryMechanicalLicenseField(field string, mechanicalLicenseId string) (interface{}, error) {
	mechanicalLicense, err := ValidateMechanicalLicense(mechanicalLicenseId)
	if err != nil {
		return nil, err
	}
	switch field {
	case "composition_right":
		return GetCompositionRight(mechanicalLicense)
	case "composition_right_transfer":
		return GetCompositionRightTransfer(mechanicalLicense)
	case "publication":
		return GetPublication(mechanicalLicense)
	case "recipient":
		return GetRecipient(mechanicalLicense)
	case "sender":
		return GetSender(mechanicalLicense)
	default:
		return nil, ErrorAppend(ErrInvalidField, field)
	}
}

func GetRecipient(data Data) (Data, error) {
	recipientId := spec.GetRecipientId(data)
	tx, err := bigchain.GetTx(recipientId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetSender(data Data) (Data, error) {
	senderId := spec.GetSenderId(data)
	tx, err := bigchain.GetTx(senderId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetCompositionRightTransfer(data Data) (Data, error) {
	transferId := spec.GetCompositionRightTransferId(data)
	tx, err := bigchain.GetTx(transferId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

// Master license

func ValidateMasterLicense(masterLicenseId string) (Data, error) {
	tx, err := bigchain.GetTx(masterLicenseId)
	if err != nil {
		return nil, err
	}
	masterLicense := bigchain.GetTxData(tx)
	if err = spec.ValidMasterLicense(masterLicense); err != nil {
		return nil, err
	}
	senderPub := bigchain.DefaultGetTxSender(tx)
	senderId := spec.GetSenderId(masterLicense)
	tx, err = bigchain.GetTx(senderId)
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
		if recordingRightId == bigchain.GetId(right) {
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
	tx, err = bigchain.GetTx(recipientId)
	if err != nil {
		return nil, err
	}
	recipient := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(recipient); err != nil {
		return nil, err
	}
	return masterLicense, nil
}

func QueryMasterLicenseField(field string, masterLicenseId string) (interface{}, error) {
	masterLicense, err := ValidateMasterLicense(masterLicenseId)
	if err != nil {
		return nil, err
	}
	switch field {
	case "recipient":
		return GetRecipient(masterLicense)
	case "recording_right":
		return GetRecordingRight(masterLicense)
	case "recording_right_transfer":
		return GetRecordingRightTransfer(masterLicense)
	case "release":
		return GetRelease(masterLicense)
	case "sender":
		return GetSender(masterLicense)
	default:
		return nil, ErrorAppend(ErrInvalidField, field)
	}
}

func GetRecordingRight(data Data) (Data, error) {
	recordingRightId := spec.GetRecordingRightId(data)
	tx, err := bigchain.GetTx(recordingRightId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetRecordingRightTransfer(data Data) (Data, error) {
	transferId := spec.GetRecordingRightTransferId(data)
	tx, err := bigchain.GetTx(transferId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetRelease(data Data) (Data, error) {
	releaseId := spec.GetReleaseId(data)
	tx, err := bigchain.GetTx(releaseId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

// Right

func ValidateCompositionRight(compositionRightId string) (Data, crypto.PublicKey, crypto.PublicKey, error) {
	tx, err := bigchain.GetTx(compositionRightId)
	if err != nil {
		return nil, nil, nil, err
	}
	compositionRight := bigchain.GetTxData(tx)
	if err = spec.ValidCompositionRight(compositionRight); err != nil {
		return nil, nil, nil, err
	}
	recipientPub := bigchain.DefaultGetTxRecipient(tx)
	recipientShares := bigchain.GetTxShares(tx)
	senderPub := bigchain.DefaultGetTxSender(tx)
	senderId := spec.GetSenderId(compositionRight)
	tx, err = bigchain.GetTx(senderId)
	if err != nil {
		return nil, nil, nil, err
	}
	sender := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(sender); err != nil {
		return nil, nil, nil, err
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, nil, nil, ErrorAppend(ErrInvalidKey, senderPub.String())
	}
	recipientId := spec.GetRecipientId(compositionRight)
	tx, err = bigchain.GetTx(recipientId)
	if err != nil {
		return nil, nil, nil, err
	}
	recipient := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(recipient); err != nil {
		return nil, nil, nil, err
	}
	if !recipientPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, nil, nil, ErrorAppend(ErrInvalidKey, recipientPub.String())
	}
	compositionRight.Set("recipientShares", recipientShares)
	return compositionRight, recipientPub, senderPub, nil
}

func ValidateRecordingRight(recordingRightId string) (Data, crypto.PublicKey, crypto.PublicKey, error) {
	tx, err := bigchain.GetTx(recordingRightId)
	if err != nil {
		return nil, nil, nil, err
	}
	recordingRight := bigchain.GetTxData(tx)
	if err = spec.ValidRecordingRight(recordingRight); err != nil {
		return nil, nil, nil, err
	}
	recipientPub := bigchain.DefaultGetTxRecipient(tx)
	recipientShares := bigchain.GetTxShares(tx)
	senderPub := bigchain.DefaultGetTxSender(tx)
	senderId := spec.GetSenderId(recordingRight)
	tx, err = bigchain.GetTx(senderId)
	if err != nil {
		return nil, nil, nil, err
	}
	sender := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(sender); err != nil {
		return nil, nil, nil, err
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, nil, nil, ErrorAppend(ErrInvalidKey, senderPub.String())
	}
	recipientId := spec.GetRecipientId(recordingRight)
	tx, err = bigchain.GetTx(recipientId)
	if err != nil {
		return nil, nil, nil, err
	}
	recipient := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(recipient); err != nil {
		return nil, nil, nil, err
	}
	if !recipientPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, nil, nil, ErrorAppend(ErrInvalidKey, recipientPub.String())
	}
	recordingRight.Set("recipientShares", recipientShares)
	return recordingRight, recipientPub, senderPub, nil
}

// Right Holders

func ValidateCompositionRightHolder(compositionRightId, publicationId, recipientId string) (Data, error) {
	_, compositionRights, err := ValidatePublication(publicationId)
	if err != nil {
		return nil, err
	}
	for _, compositionRight := range compositionRights {
		if compositionRightId == bigchain.GetId(compositionRight) {
			if recipientId != spec.GetRecipientId(compositionRight) {
				return nil, ErrorAppend(ErrCriteriaNotMet, "agent does not hold composition right")
			}
			return compositionRight, nil
		}
	}
	return nil, ErrorAppend(ErrCriteriaNotMet, "publication does not link to composition right")
}

func ValidateRecordingRightHolder(recipientId, recordingRightId, releaseId string) (Data, error) {
	_, recordingRights, err := ValidateRelease(releaseId)
	if err != nil {
		return nil, err
	}
	for _, recordingRight := range recordingRights {
		if recordingRightId == bigchain.GetId(recordingRight) {
			if recipientId != spec.GetRecipientId(recordingRight) {
				return nil, ErrorAppend(ErrCriteriaNotMet, "agent does not hold recording right")
			}
			return recordingRight, nil
		}
	}
	return nil, ErrorAppend(ErrCriteriaNotMet, "release does not link to recording right")
}

// Transfer

func ValidateCompositionRightTransfer(transferId string) (Data, error) {
	tx, err := bigchain.GetTx(transferId)
	if err != nil {
		return nil, err
	}
	transfer := bigchain.GetTxData(tx)
	if err = spec.ValidCompositionRightTransfer(transfer); err != nil {
		return nil, err
	}
	senderPub := bigchain.DefaultGetTxSender(tx)
	recipientId := spec.GetRecipientId(transfer)
	tx, err = bigchain.GetTx(recipientId)
	if err != nil {
		return nil, err
	}
	recipient := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(recipient); err != nil {
		return nil, err
	}
	recipientPub := bigchain.DefaultGetTxSender(tx)
	if senderPub.Equals(recipientPub) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "recipient and sender keys must be different")
	}
	senderId := spec.GetSenderId(transfer)
	tx, err = bigchain.GetTx(senderId)
	if err != nil {
		return nil, err
	}
	sender := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(sender); err != nil {
		return nil, err
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, ErrorAppend(ErrInvalidKey, senderPub.String())
	}
	publicationId := spec.GetPublicationId(transfer)
	_, compositionRights, err := ValidatePublication(publicationId)
	if err != nil {
		return nil, err
	}
	txId := spec.GetTxId(transfer)
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
	compositionRightId := spec.GetCompositionRightId(transfer)
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
	transfer.Set("recipientShares", bigchain.GetTxOutputAmount(txTransfer, 0))
	if n == 2 {
		transfer.Set("senderShares", bigchain.GetTxOutputAmount(txTransfer, 1))
	}
	return transfer, nil
}

func ValidateRecordingRightTransfer(transferId string) (Data, error) {
	tx, err := bigchain.GetTx(transferId)
	if err != nil {
		return nil, err
	}
	transfer := bigchain.GetTxData(tx)
	if err = spec.ValidRecordingRightTransfer(transfer); err != nil {
		return nil, err
	}
	senderPub := bigchain.DefaultGetTxSender(tx)
	recipientId := spec.GetRecipientId(transfer)
	tx, err = bigchain.GetTx(recipientId)
	if err != nil {
		return nil, err
	}
	recipient := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(recipient); err != nil {
		return nil, err
	}
	recipientPub := bigchain.DefaultGetTxSender(tx)
	if senderPub.Equals(recipientPub) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "recipient and sender keys must be different")
	}
	senderId := spec.GetSenderId(transfer)
	tx, err = bigchain.GetTx(senderId)
	if err != nil {
		return nil, err
	}
	sender := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(sender); err != nil {
		return nil, err
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, ErrorAppend(ErrInvalidKey, senderPub.String())
	}
	releaseId := spec.GetReleaseId(transfer)
	_, recordingRights, err := ValidateRelease(releaseId)
	if err != nil {
		return nil, err
	}
	txId := spec.GetTxId(transfer)
	txTransfer, err := bigchain.GetTx(txId)
	if err != nil {
		return nil, err
	}
	if bigchain.TRANSFER != bigchain.GetTxOperation(txTransfer) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "expected TRANSFER tx")
	}
	n := len(bigchain.GetTxOutputs(txTransfer))
	if n != 1 && n != 2 {
		return nil, ErrorAppend(ErrInvalidSize, "outputs must have size 1 or 2")
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(txTransfer)) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "sender is not signer of TRANSFER tx")
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
	recordingRightId := spec.GetRecordingRightId(transfer)
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
		return nil, ErrorAppend(ErrCriteriaNotMet, "publication does not link to recording right")
	}
	transfer.Set("recipientShares", bigchain.GetTxOutputAmount(txTransfer, 0))
	if n == 2 {
		transfer.Set("senderShares", bigchain.GetTxOutputAmount(txTransfer, 1))
	}
	return transfer, nil
}
