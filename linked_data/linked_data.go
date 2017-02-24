package linked_data

import (
	"github.com/zbo14/envoke/bigchain"
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/crypto"
	"github.com/zbo14/envoke/spec"
)

// Model

func ValidateModelById(modelId string) (Data, error) {
	tx, err := bigchain.GetTx(modelId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	model := bigchain.GetTxData(tx)
	pub := bigchain.GetTxPublicKey(tx)
	return ValidateModel(model, pub)
}

func ValidateModel(model Data, pub crypto.PublicKey) (Data, error) {
	var err error
	_type := spec.GetType(model)
	switch _type {
	case spec.AGENT:
		err = spec.ValidAgent(model)
	case spec.COMPOSITION:
		model, err = ValidateComposition(model, pub)
	case spec.RECORDING:
		model, err = ValidateRecording(model, pub)
	case spec.LICENSE_PUBLISHING:
		model, err = ValidatePublishingLicense(model, pub)
	case spec.LICENSE_RECORDING:
		model, err = ValidateRecordingLicense(model, pub)
	default:
		return nil, ErrorAppend(ErrInvalidType, _type)
	}
	if err != nil {
		return nil, err
	}
	return model, nil
}

func QueryModelIdField(field, modelId string) (interface{}, error) {
	tx, err := bigchain.GetTx(modelId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	model := bigchain.GetTxData(tx)
	pub := bigchain.GetTxPublicKey(tx)
	return QueryModelField(field, model, pub)
}

func QueryModelField(field string, model Data, pub crypto.PublicKey) (interface{}, error) {
	_type := spec.GetType(model)
	switch _type {
	case spec.INFO_COMPOSITION:
		return QueryCompositionInfoField(field, model, pub)
	case spec.INFO_RECORDING:
		return QueryRecordingInfoField(field, model, pub)
	case spec.COMPOSITION:
		return QueryCompositionField(field, model, pub)
	case spec.RECORDING:
		return QueryRecordingField(field, model, pub)
	case spec.LICENSE_PUBLISHING:
		return QueryPublishingLicenseField(field, model, pub)
	case spec.LICENSE_RECORDING:
		return QueryRecordingLicenseField(field, model, pub)
	default:
		return nil, ErrorAppend(ErrInvalidType, _type)
	}
}

// Composition

func ValidateCompositionInfoById(infoId string) (Data, error) {
	tx, err := bigchain.GetTx(infoId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	pub := bigchain.GetTxPublicKey(tx)
	info := bigchain.GetTxData(tx)
	if err := spec.ValidCompositionInfo(info); err != nil {
		return nil, err
	}
	return ValidateCompositionInfo(info, pub)
}

func ValidateCompositionInfo(info Data, pub crypto.PublicKey) (Data, error) {
	composerId := spec.GetInfoComposer(info)
	tx, err := bigchain.GetTx(composerId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	signed := false
	if pub.Equals(bigchain.GetTxPublicKey(tx)) {
		signed = true
	}
	composer := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(composer); err != nil {
		return nil, err
	}
	publisherId := spec.GetInfoPublisher(info)
	tx, err = bigchain.GetTx(publisherId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	if pub.Equals(bigchain.GetTxPublicKey(tx)) {
		signed = true
	}
	if !signed {
		return nil, ErrInvalidKey
	}
	publisher := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(publisher); err != nil {
		return nil, err
	}
	return info, nil
}

func QueryCompositionInfoField(field string, info Data, pub crypto.PublicKey) (interface{}, error) {
	if _, err := ValidateCompositionInfo(info, pub); err != nil {
		return nil, err
	}
	switch field {
	case "composer":
		return GetInfoComposer(info)
	case "publisher":
		return GetInfoPublisher(info)
	default:
		return nil, ErrorAppend(ErrInvalidField, field)
	}
}

func GetInfoComposer(info Data) (Data, error) {
	composerId := spec.GetInfoComposer(info)
	tx, err := bigchain.GetTx(composerId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetInfoPublisher(info Data) (Data, error) {
	publisherId := spec.GetInfoPublisher(info)
	tx, err := bigchain.GetTx(publisherId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func ValidateCompositionById(compositionId string) (Data, error) {
	tx, err := bigchain.GetTx(compositionId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	pub := bigchain.GetTxPublicKey(tx)
	composition := bigchain.GetTxData(tx)
	if err := spec.ValidComposition(composition); err != nil {
		return nil, err
	}
	return ValidateComposition(composition, pub)
}

func ValidateComposition(composition Data, pub crypto.PublicKey) (Data, error) {
	infoId := spec.GetCompositionInfo(composition)
	info, err := ValidateCompositionInfoById(infoId)
	if err != nil {
		return nil, err
	}
	signerId := spec.GetInfoComposer(info)
	tx, err := bigchain.GetTx(signerId)
	if err != nil {
		return nil, err
	}
	if !pub.Equals(bigchain.GetTxPublicKey(tx)) {
		signerId = spec.GetInfoPublisher(info)
		tx, err = bigchain.GetTx(signerId)
		if err != nil {
			return nil, err
		}
		if !pub.Equals(bigchain.GetTxPublicKey(tx)) {
			return nil, ErrorAppend(ErrCriteriaNotMet, "composition must be signed by composer or publisher")
		}
	}
	percentageShares := 0
	rightHolders := make(map[string]struct{})
	rightIds := spec.GetCompositionRights(composition)
	for _, rightId := range rightIds {
		tx, err = bigchain.GetTx(rightId)
		if err != nil {
			return nil, err
		}
		if !bigchain.FulfilledTx(tx) {
			return nil, ErrInvalidFulfillment
		}
		right := bigchain.GetTxData(tx)
		if err = spec.ValidRight(right); err != nil {
			return nil, err
		}
		if _, ok := rightHolders[bigchain.GetTxPublicKey(tx).String()]; ok {
			return nil, ErrorAppend(ErrCriteriaNotMet, "rightHolder cannot have multiple rights to composition")
		}
		if infoId != spec.GetRightInfo(right) {
			return nil, ErrorAppend(ErrInvalidId, "right infoId")
		}
		percentageShares += spec.GetRightPercentageShares(right)
		if percentageShares > 100 {
			return nil, ErrorAppend(ErrCriteriaNotMet, "percentage shares cannot exceed 100")
		}
	}
	if percentageShares != 100 {
		return nil, ErrorAppend(ErrCriteriaNotMet, "total percentage shares do not equal 100")
	}
	return composition, nil
}

func QueryCompositionField(field string, composition Data, pub crypto.PublicKey) (interface{}, error) {
	if _, err := ValidateComposition(composition, pub); err != nil {
		return nil, err
	}
	switch field {
	case "info":
		return GetCompositionInfo(composition)
	case "rights":
		return GetCompositionRights(composition)
	default:
		return nil, ErrorAppend(ErrInvalidField, field)
	}
}

func GetCompositionInfo(composition Data) (Data, error) {
	infoId := spec.GetCompositionInfo(composition)
	tx, err := bigchain.GetTx(infoId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetCompositionRights(composition Data) ([]Data, error) {
	rightIds := spec.GetCompositionRights(composition)
	rights := make([]Data, len(rightIds))
	for i, rightId := range rightIds {
		tx, err := bigchain.GetTx(rightId)
		if err != nil {
			return nil, err
		}
		if !bigchain.FulfilledTx(tx) {
			return nil, ErrInvalidFulfillment
		}
		rights[i] = bigchain.GetTxData(tx)
	}
	return rights, nil
}

// Recording

func ValidateRecordingInfoById(infoId string) (Data, error) {
	tx, err := bigchain.GetTx(infoId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	pub := bigchain.GetTxPublicKey(tx)
	info := bigchain.GetTxData(tx)
	if err := spec.ValidRecordingInfo(info); err != nil {
		return nil, err
	}
	return ValidateRecordingInfo(info, pub)
}

func ValidateRecordingInfo(info Data, pub crypto.PublicKey) (Data, error) {
	labelId := spec.GetInfoLabel(info)
	tx, err := bigchain.GetTx(labelId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	label := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(label); err != nil {
		return nil, err
	}
	signed := false
	var signerId string
	if pub.Equals(bigchain.GetTxPublicKey(tx)) {
		signed = true
		signerId = labelId
	}
	performerId := spec.GetInfoPerformer(info)
	tx, err = bigchain.GetTx(performerId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	performer := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(performer); err != nil {
		return nil, err
	}
	if pub.Equals(bigchain.GetTxPublicKey(tx)) {
		signed = true
		signerId = performerId
	}
	producerId := spec.GetInfoProducer(info)
	tx, err = bigchain.GetTx(producerId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	producer := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(producer); err != nil {
		return nil, err
	}
	if pub.Equals(bigchain.GetTxPublicKey(tx)) {
		signed = true
		signerId = producerId
	}
	if !signed {
		return nil, ErrorAppend(ErrCriteriaNotMet, "recording must be signed by label, performer, or producer")
	}
	compositionId := spec.GetInfoComposition(info)
	composition, err := ValidateCompositionById(compositionId)
	if err != nil {
		return nil, err
	}
	rightHolder := false
	rightIds := spec.GetCompositionRights(composition)
	for _, rightId := range rightIds {
		tx, err := bigchain.GetTx(rightId)
		if err != nil {
			return nil, err
		}
		if pub.Equals(bigchain.GetTxPublicKey(tx)) {
			rightHolder = true
			break
		}
	}
	if !rightHolder {
		licenseId := spec.GetInfoPublishingLicense(info)
		license, err := ValidatePublishingLicenseById(licenseId)
		if err != nil {
			return nil, err
		}
		if compositionId != spec.GetLicenseComposition(license) {
			return nil, ErrorAppend(ErrCriteriaNotMet, "publishing license is not for composition")
		}
		if signerId != spec.GetLicenseLicensee(license) {
			return nil, ErrorAppend(ErrCriteriaNotMet, "signer is not licensee of publishing license")
		}
	}
	return info, nil
}

func QueryRecordingInfoField(field string, info Data, pub crypto.PublicKey) (interface{}, error) {
	if _, err := ValidateRecordingInfo(info, pub); err != nil {
		return nil, err
	}
	switch field {
	case "composition":
		return GetInfoComposition(info)
	case "label":
		return GetInfoLabel(info)
	case "performer":
		return GetInfoPerformer(info)
	case "producer":
		return GetInfoProducer(info)
	case "publishingLicense":
		return GetInfoPublishingLicense(info)
	default:
		return nil, ErrorAppend(ErrInvalidField, field)
	}
}

func GetInfoComposition(info Data) (Data, error) {
	compositionId := spec.GetInfoComposition(info)
	tx, err := bigchain.GetTx(compositionId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetInfoLabel(info Data) (Data, error) {
	labelId := spec.GetInfoLabel(info)
	tx, err := bigchain.GetTx(labelId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetInfoPerformer(info Data) (Data, error) {
	performerId := spec.GetInfoPerformer(info)
	tx, err := bigchain.GetTx(performerId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetInfoProducer(info Data) (Data, error) {
	producerId := spec.GetInfoProducer(info)
	tx, err := bigchain.GetTx(producerId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetInfoPublishingLicense(info Data) (Data, error) {
	licenseId := spec.GetInfoPublishingLicense(info)
	tx, err := bigchain.GetTx(licenseId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func ValidateRecordingById(recordingId string) (Data, error) {
	tx, err := bigchain.GetTx(recordingId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	pub := bigchain.GetTxPublicKey(tx)
	recording := bigchain.GetTxData(tx)
	if err := spec.ValidRecording(recording); err != nil {
		return nil, err
	}
	return ValidateRecording(recording, pub)
}

func ValidateRecording(recording Data, pub crypto.PublicKey) (Data, error) {
	infoId := spec.GetRecordingInfo(recording)
	tx, err := bigchain.GetTx(infoId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	info := bigchain.GetTxData(tx)
	if _, err := ValidateRecordingInfo(info, pub); err != nil {
		return nil, err
	}
	percentageShares := 0
	rightHolders := make(map[string]struct{})
	rightIds := spec.GetRecordingRights(recording)
	for _, rightId := range rightIds {
		tx, err = bigchain.GetTx(rightId)
		if err != nil {
			return nil, err
		}
		if !bigchain.FulfilledTx(tx) {
			return nil, ErrInvalidFulfillment
		}
		pubstr := bigchain.GetTxPublicKey(tx).String()
		if _, ok := rightHolders[pubstr]; ok {
			return nil, ErrorAppend(ErrCriteriaNotMet, "rightHolder cannot have multiple rights to composition")
		}
		right := bigchain.GetTxData(tx)
		if err = spec.ValidRight(right); err != nil {
			return nil, err
		}
		percentageShares += spec.GetRightPercentageShares(right)
		if percentageShares > 100 {
			return nil, ErrorAppend(ErrCriteriaNotMet, "percentage shares cannot exceed 100")
		}
	}
	if percentageShares != 100 {
		return nil, ErrorAppend(ErrCriteriaNotMet, "total percentage shares do not equal 100")
	}
	return recording, nil
}

func QueryRecordingField(field string, recording Data, pub crypto.PublicKey) (interface{}, error) {
	if _, err := ValidateRecording(recording, pub); err != nil {
		return nil, err
	}
	switch field {
	case "info":
		return GetRecordingInfo(recording)
	case "rights":
		return GetRecordingRights(recording)
	default:
		return nil, ErrorAppend(ErrInvalidField, field)
	}
}

func GetRecordingInfo(recording Data) (Data, error) {
	infoId := spec.GetRecordingInfo(recording)
	tx, err := bigchain.GetTx(infoId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetRecordingRights(recording Data) ([]Data, error) {
	rightIds := spec.GetRecordingRights(recording)
	rights := make([]Data, len(rightIds))
	for i, rightId := range rightIds {
		tx, err := bigchain.GetTx(rightId)
		if err != nil {
			return nil, err
		}
		if !bigchain.FulfilledTx(tx) {
			return nil, ErrInvalidFulfillment
		}
		rights[i] = bigchain.GetTxData(tx)
	}
	return rights, nil
}

// Publishing License

func ValidatePublishingLicenseById(licenseId string) (Data, error) {
	tx, err := bigchain.GetTx(licenseId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	pub := bigchain.GetTxPublicKey(tx)
	license := bigchain.GetTxData(tx)
	if err := spec.ValidPublishingLicense(license); err != nil {
		return nil, err
	}
	return ValidatePublishingLicense(license, pub)
}

func ValidatePublishingLicense(license Data, pub crypto.PublicKey) (Data, error) {
	compositionId := spec.GetLicenseComposition(license)
	composition, err := ValidateCompositionById(compositionId)
	if err != nil {
		return nil, err
	}
	rightHolder := false
	rightIds := spec.GetCompositionRights(composition)
	for _, rightId := range rightIds {
		tx, err := bigchain.GetTx(rightId)
		if err != nil {
			return nil, err
		}
		if pub.Equals(bigchain.GetTxPublicKey(tx)) {
			rightHolder = true
			break
		}
	}
	if !rightHolder {
		return nil, ErrorAppend(ErrCriteriaNotMet, "signer is not a composition right holder")
	}
	licenserId := spec.GetLicenseLicenser(license)
	tx, err := bigchain.GetTx(licenserId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, err
	}
	if !pub.Equals(bigchain.GetTxPublicKey(tx)) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "signer must be licenser")
	}
	// necessary?
	licenser := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(licenser); err != nil {
		return nil, err
	}
	licenseeId := spec.GetLicenseLicensee(license)
	tx, err = bigchain.GetTx(licenseeId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	licensee := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(licensee); err != nil {
		return nil, err
	}
	return license, nil
}

func QueryPublishingLicenseField(field string, license Data, pub crypto.PublicKey) (interface{}, error) {
	if _, err := ValidatePublishingLicense(license, pub); err != nil {
		return nil, err
	}
	switch field {
	case "licensee":
		return GetPublishingLicenseLicensee(license)
	case "licenser":
		return GetPublishingLicenseLicenser(license)
	case "composition":
		return GetPublishingLicenseComposition(license)
	default:
		return nil, ErrorAppend(ErrInvalidField, field)
	}
}

func GetPublishingLicenseLicensee(license Data) (Data, error) {
	licenseeId := spec.GetLicenseLicensee(license)
	tx, err := bigchain.GetTx(licenseeId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetPublishingLicenseLicenser(license Data) (Data, error) {
	licenserId := spec.GetLicenseLicenser(license)
	tx, err := bigchain.GetTx(licenserId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetPublishingLicenseComposition(license Data) (Data, error) {
	compositionId := spec.GetLicenseComposition(license)
	tx, err := bigchain.GetTx(compositionId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

// Recording license

func ValidateRecordingLicenseById(licenseId string) (Data, error) {
	tx, err := bigchain.GetTx(licenseId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	pub := bigchain.GetTxPublicKey(tx)
	license := bigchain.GetTxData(tx)
	if err := spec.ValidRecordingLicense(license); err != nil {
		return nil, err
	}
	return ValidateRecordingLicense(license, pub)
}

func ValidateRecordingLicense(license Data, pub crypto.PublicKey) (Data, error) {
	recordingId := spec.GetLicenseRecording(license)
	recording, err := ValidateRecordingById(recordingId)
	if err != nil {
		return nil, err
	}
	rightHolder := false
	rightIds := spec.GetRecordingRights(recording)
	for _, rightId := range rightIds {
		tx, err := bigchain.GetTx(rightId)
		if err != nil {
			return nil, err
		}
		if pub.Equals(bigchain.GetTxPublicKey(tx)) {
			rightHolder = true
			break
		}
	}
	if !rightHolder {
		return nil, ErrorAppend(ErrCriteriaNotMet, "signer is not a recording right holder")
	}
	licenserId := spec.GetLicenseLicenser(license)
	tx, err := bigchain.GetTx(licenserId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, err
	}
	if !pub.Equals(bigchain.GetTxPublicKey(tx)) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "signer must be licenser")
	}
	licenser := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(licenser); err != nil {
		return nil, err
	}
	licenseeId := spec.GetLicenseLicensee(license)
	tx, err = bigchain.GetTx(licenseeId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	licensee := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(licensee); err != nil {
		return nil, err
	}
	return license, nil
}

func QueryRecordingLicenseField(field string, license Data, pub crypto.PublicKey) (interface{}, error) {
	if _, err := ValidateRecordingLicense(license, pub); err != nil {
		return nil, err
	}
	switch field {
	case "licensee":
		return GetRecordingLicenseLicensee(license)
	case "licenser":
		return GetRecordingLicenseLicenser(license)
	case "recording":
		return GetRecordingLicenseRecording(license)
	default:
		return nil, ErrorAppend(ErrInvalidField, field)
	}
}

func GetRecordingLicenseLicensee(license Data) (Data, error) {
	licenseeId := spec.GetLicenseLicensee(license)
	tx, err := bigchain.GetTx(licenseeId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetRecordingLicenseLicenser(license Data) (Data, error) {
	licenserId := spec.GetLicenseLicenser(license)
	tx, err := bigchain.GetTx(licenserId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetRecordingLicenseRecording(license Data) (Data, error) {
	recordingId := spec.GetLicenseRecording(license)
	tx, err := bigchain.GetTx(recordingId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}
