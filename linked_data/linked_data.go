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
	composerId := spec.GetCompositionComposer(composition)
	tx, err := bigchain.GetTx(composerId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	composerPub := bigchain.GetTxPublicKey(tx)
	composer := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(composer); err != nil {
		return nil, err
	}
	publisherId := spec.GetCompositionPublisher(composition)
	tx, err = bigchain.GetTx(publisherId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	publisherPub := bigchain.GetTxPublicKey(tx)
	publisher := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(publisher); err != nil {
		return nil, err
	}
	if pub.Equals(composerPub) {
		//..
	} else if pub.Equals(publisherPub) {
		//..
	} else {
		return nil, ErrInvalidKey
	}
	rights := spec.GetCompositionRights(composition)
	for _, right := range rights {
		rightHolderId := spec.GetRightHolder(right)
		tx, err = bigchain.GetTx(rightHolderId)
		if err != nil {
			return nil, err
		}
		if !bigchain.FulfilledTx(tx) {
			return nil, ErrInvalidFulfillment
		}
		rightHolder := bigchain.GetTxData(tx)
		if err = spec.ValidAgent(rightHolder); err != nil {
			return nil, err
		}
	}
	return composition, nil
}

func QueryCompositionField(field string, composition Data, pub crypto.PublicKey) (interface{}, error) {
	if _, err := ValidateComposition(composition, pub); err != nil {
		return nil, err
	}
	switch field {
	case "composer":
		return GetCompositionComposer(composition)
	case "publisher":
		return GetCompositionPublisher(composition)
	case "rightHolders":
		return GetCompositionRightHolders(composition)
	default:
		return nil, ErrorAppend(ErrInvalidField, field)
	}
}

func GetCompositionComposer(composition Data) (Data, error) {
	composerId := spec.GetCompositionComposer(composition)
	tx, err := bigchain.GetTx(composerId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetCompositionPublisher(composition Data) (Data, error) {
	publisherId := spec.GetCompositionPublisher(composition)
	tx, err := bigchain.GetTx(publisherId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetCompositionRightHolders(composition Data) ([]Data, error) {
	return GetRightHolders(spec.GetCompositionRights(composition))
}

// Recording

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
	compositionId := spec.GetRecordingComposition(recording)
	composition, err := ValidateCompositionById(compositionId)
	if err != nil {
		return nil, err
	}
	rights := spec.GetCompositionRights(composition)
	rightHolderIds := make(map[string]struct{})
	for _, right := range rights {
		rightHolderIds[spec.GetRightHolder(right)] = struct{}{}
	}
	labelId := spec.GetRecordingLabel(recording)
	tx, err := bigchain.GetTx(labelId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	labelPub := bigchain.GetTxPublicKey(tx)
	label := bigchain.GetTxData(tx)
	if _, ok := rightHolderIds[labelId]; !ok {
		if err = spec.ValidAgent(label); err != nil {
			return nil, err
		}
	}
	performerId := spec.GetRecordingPerformer(recording)
	tx, err = bigchain.GetTx(performerId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	performerPub := bigchain.GetTxPublicKey(tx)
	performer := bigchain.GetTxData(tx)
	if _, ok := rightHolderIds[performerId]; !ok {
		if err = spec.ValidAgent(performer); err != nil {
			return nil, err
		}
	}
	producerId := spec.GetRecordingProducer(recording)
	tx, err = bigchain.GetTx(producerId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	producer := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(producer); err != nil {
		return nil, ErrorAppend(ErrInvalidModel, spec.AGENT)
	}
	if pub.Equals(labelPub) {
		if _, ok := rightHolderIds[labelId]; !ok {
			licenseId := spec.GetRecordingPublishingLicense(recording)
			license, err := ValidatePublishingLicenseById(licenseId)
			if err != nil {
				return nil, err
			}
			if labelId != spec.GetLicenseLicensee(license) || compositionId != spec.GetLicenseComposition(license) {
				return nil, ErrorAppend(ErrCriteriaNotMet, "label does not hold composition right/publishing license to this recording")
			}
		}
	} else if pub.Equals(performerPub) {
		if _, ok := rightHolderIds[performerId]; !ok {
			licenseId := spec.GetRecordingPublishingLicense(recording)
			license, err := ValidatePublishingLicenseById(licenseId)
			if err != nil {
				return nil, err
			}
			if performerId != spec.GetLicenseLicensee(license) || compositionId != spec.GetLicenseComposition(license) {
				return nil, ErrorAppend(ErrCriteriaNotMet, "performer does not hold composition right/publishing license to this recording")
			}
		}
	} else {
		return nil, ErrorAppend(ErrCriteriaNotMet, "recording must be signed by label or performer")
	}
	rights = spec.GetRecordingRights(recording)
	for _, right := range rights {
		rightHolderId := spec.GetRightHolder(right)
		tx, err = bigchain.GetTx(rightHolderId)
		if err != nil {
			return nil, err
		}
		if !bigchain.FulfilledTx(tx) {
			return nil, ErrInvalidFulfillment
		}
		rightHolder := bigchain.GetTxData(tx)
		if err = spec.ValidAgent(rightHolder); err != nil {
			return nil, ErrorAppend(ErrInvalidModel, spec.AGENT)
		}
	}
	return recording, nil
}

func QueryRecordingField(field string, recording Data, pub crypto.PublicKey) (interface{}, error) {
	if _, err := ValidateRecording(recording, pub); err != nil {
		return nil, err
	}
	switch field {
	case "composition":
		return GetRecordingComposition(recording)
	case "label":
		return GetRecordingLabel(recording)
	case "performer":
		return GetRecordingPerformer(recording)
	case "producer":
		return GetRecordingProducer(recording)
	case "publishingLicense":
		return GetRecordingPublishingLicense(recording)
	case "rightHolders":
		return GetRecordingRightHolders(recording)
	default:
		return nil, ErrorAppend(ErrInvalidField, field)
	}
}

func GetRecordingComposition(recording Data) (Data, error) {
	compositionId := spec.GetRecordingComposition(recording)
	tx, err := bigchain.GetTx(compositionId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetRecordingLabel(recording Data) (Data, error) {
	labelId := spec.GetRecordingLabel(recording)
	tx, err := bigchain.GetTx(labelId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetRecordingPerformer(recording Data) (Data, error) {
	performerId := spec.GetRecordingPerformer(recording)
	tx, err := bigchain.GetTx(performerId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetRecordingProducer(recording Data) (Data, error) {
	producerId := spec.GetRecordingProducer(recording)
	tx, err := bigchain.GetTx(producerId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetRecordingPublishingLicense(recording Data) (Data, error) {
	licenseId := spec.GetRecordingPublishingLicense(recording)
	tx, err := bigchain.GetTx(licenseId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetRecordingRightHolders(recording Data) ([]Data, error) {
	return GetRightHolders(spec.GetRecordingRights(recording))
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
	rights := spec.GetCompositionRights(composition)
	rightHolderIds := make(map[string]struct{})
	for _, right := range rights {
		rightHolderIds[spec.GetRightHolder(right)] = struct{}{}
	}
	licenserId := spec.GetLicenseLicenser(license)
	if _, ok := rightHolderIds[licenserId]; !ok {
		return nil, ErrorAppend(ErrCriteriaNotMet, "licenser is not a composition right holder")
	}
	tx, err := bigchain.GetTx(licenserId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, err
	}
	if !pub.Equals(bigchain.GetTxPublicKey(tx)) {
		return nil, ErrInvalidKey
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
	rights := spec.GetRecordingRights(recording)
	rightHolderIds := make(map[string]struct{})
	for _, right := range rights {
		rightHolderIds[spec.GetRightHolder(right)] = struct{}{}
	}
	licenserId := spec.GetLicenseLicenser(license)
	if _, ok := rightHolderIds[licenserId]; !ok {
		return nil, ErrorAppend(ErrCriteriaNotMet, "licenser is not a recording right holder")
	}
	tx, err := bigchain.GetTx(licenserId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, err
	}
	if !pub.Equals(bigchain.GetTxPublicKey(tx)) {
		return nil, ErrInvalidKey
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

// Right holders

func GetRightHolders(rights []Data) ([]Data, error) {
	rightHolders := make([]Data, len(rights))
	for i, right := range rights {
		rightHolderId := spec.GetRightHolder(right)
		tx, err := bigchain.GetTx(rightHolderId)
		if err != nil {
			return nil, err
		}
		if !bigchain.FulfilledTx(tx) {
			return nil, ErrInvalidFulfillment
		}
		rightHolders[i] = bigchain.GetTxData(tx)
	}
	return rightHolders, nil
}
