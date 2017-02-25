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
	case spec.PUBLICATION:
		model, err = ValidatePublication(model, pub)
	case spec.RELEASE:
		model, err = ValidateRelease(model, pub)
	case spec.LICENSE_PUBLISHING:
		model, err = ValidatePublishingLicense(model, pub)
	case spec.LICENSE_RELEASE:
		model, err = ValidateReleaseLicense(model, pub)
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
	case spec.PUBLICATION:
		return QueryPublicationField(field, model, pub)
	case spec.RELEASE:
		return QueryReleaseField(field, model, pub)
	case spec.LICENSE_PUBLISHING:
		return QueryPublishingLicenseField(field, model, pub)
	case spec.LICENSE_RELEASE:
		return QueryReleaseLicenseField(field, model, pub)
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
	composition := bigchain.GetTxData(tx)
	if err := spec.ValidComposition(composition); err != nil {
		return nil, err
	}
	pub := bigchain.GetTxPublicKey(tx)
	return ValidateComposition(composition, pub)
}

func ValidateComposition(composition Data, pub crypto.PublicKey) (Data, error) {
	composerId := spec.GetCompositionComposerId(composition)
	tx, err := bigchain.GetTx(composerId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	if !pub.Equals(bigchain.GetTxPublicKey(tx)) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "composition must be signed by composer")
	}
	composer := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(composer); err != nil {
		return nil, err
	}
	publisherId := spec.GetCompositionPublisherId(composition)
	tx, err = bigchain.GetTx(publisherId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	publisher := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(publisher); err != nil {
		return nil, err
	}
	return composition, nil
}

func GetCompositionComposer(composition Data) (Data, error) {
	composerId := spec.GetCompositionComposerId(composition)
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
	publisherId := spec.GetCompositionPublisherId(composition)
	tx, err := bigchain.GetTx(publisherId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func ValidatePublicationById(publicationId string) (Data, error) {
	tx, err := bigchain.GetTx(publicationId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	publication := bigchain.GetTxData(tx)
	if err := spec.ValidPublication(publication); err != nil {
		return nil, err
	}
	pub := bigchain.GetTxPublicKey(tx)
	return ValidatePublication(publication, pub)
}

func ValidatePublication(publication Data, pub crypto.PublicKey) (Data, error) {
	compositionId := spec.GetPublicationCompositionId(publication)
	composition, err := ValidateCompositionById(compositionId)
	if err != nil {
		return nil, err
	}
	publisherId := spec.GetCompositionPublisherId(composition)
	tx, err := bigchain.GetTx(publisherId)
	if err != nil {
		return nil, err
	}
	if !pub.Equals(bigchain.GetTxPublicKey(tx)) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "publication must be signed by publisher")
	}
	percentageShares := 0
	rightHolders := make(map[string]struct{})
	rightIds := spec.GetPublicationRightIds(publication)
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
		if compositionId != spec.GetRightCompositionId(right) {
			return nil, ErrorAppend(ErrInvalidId, "right compositionId")
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

func QueryPublicationField(field string, publication Data, pub crypto.PublicKey) (interface{}, error) {
	if _, err := ValidatePublication(publication, pub); err != nil {
		return nil, err
	}
	switch field {
	case "composer":
		composition, err := GetPublicationComposition(publication)
		if err != nil {
			return nil, err
		}
		return GetCompositionComposer(composition)
	case "composition":
		return GetPublicationComposition(publication)
	case "publisher":
		composition, err := GetPublicationComposition(publication)
		if err != nil {
			return nil, err
		}
		return GetCompositionPublisher(composition)
	case "rights":
		return GetPublicationRights(publication)
	case "title":
		composition, err := GetPublicationComposition(publication)
		if err != nil {
			return nil, err
		}
		return spec.GetCompositionTitle(composition), nil
	default:
		return nil, ErrorAppend(ErrInvalidField, field)
	}
}

func GetPublicationComposition(publication Data) (Data, error) {
	compositionId := spec.GetPublicationCompositionId(publication)
	tx, err := bigchain.GetTx(compositionId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetPublicationRights(publication Data) ([]Data, error) {
	rightIds := spec.GetPublicationRightIds(publication)
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

func ValidateRecordingById(recordingId string) (Data, error) {
	tx, err := bigchain.GetTx(recordingId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	recording := bigchain.GetTxData(tx)
	if err := spec.ValidRecording(recording); err != nil {
		return nil, err
	}
	pub := bigchain.GetTxPublicKey(tx)
	return ValidateRecording(recording, pub)
}

func ValidateRecording(recording Data, pub crypto.PublicKey) (Data, error) {
	labelId := spec.GetRecordingLabelId(recording)
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
	performerId := spec.GetRecordingPerformerId(recording)
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
	signed := false
	var signerId string
	if pub.Equals(bigchain.GetTxPublicKey(tx)) {
		signed = true
		signerId = performerId
	}
	producerId := spec.GetRecordingProducerId(recording)
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
		return nil, ErrorAppend(ErrCriteriaNotMet, "recording must be signed by performer or producer")
	}
	publicationId := spec.GetRecordingPublicationId(recording)
	publication, err := ValidatePublicationById(publicationId)
	if err != nil {
		return nil, err
	}
	rightHolder := false
	rightIds := spec.GetPublicationRightIds(publication)
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
		licenseId := spec.GetRecordingPublishingLicenseId(recording)
		license, err := ValidatePublishingLicenseById(licenseId)
		if err != nil {
			return nil, err
		}
		if publicationId != spec.GetLicensePublicationId(license) {
			return nil, ErrorAppend(ErrCriteriaNotMet, "wrong publishing license")
		}
		if signerId != spec.GetLicenseLicenseeId(license) {
			return nil, ErrorAppend(ErrCriteriaNotMet, "signer is not licensee of publishing license")
		}
	}
	return recording, nil
}
func GetRecordingPublication(recording Data) (Data, error) {
	publicationId := spec.GetRecordingPublicationId(recording)
	tx, err := bigchain.GetTx(publicationId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetRecordingLabel(recording Data) (Data, error) {
	labelId := spec.GetRecordingLabelId(recording)
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
	performerId := spec.GetRecordingPerformerId(recording)
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
	producerId := spec.GetRecordingProducerId(recording)
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
	licenseId := spec.GetRecordingPublishingLicenseId(recording)
	tx, err := bigchain.GetTx(licenseId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func ValidateReleaseById(releaseId string) (Data, error) {
	tx, err := bigchain.GetTx(releaseId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	release := bigchain.GetTxData(tx)
	if err := spec.ValidRelease(release); err != nil {
		return nil, err
	}
	pub := bigchain.GetTxPublicKey(tx)
	return ValidateRelease(release, pub)
}

func ValidateRelease(release Data, pub crypto.PublicKey) (Data, error) {
	recordingId := spec.GetReleaseRecordingId(release)
	recording, err := ValidateRecordingById(recordingId)
	if err != nil {
		return nil, err
	}
	labelId := spec.GetRecordingLabelId(recording)
	tx, err := bigchain.GetTx(labelId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	if !pub.Equals(bigchain.GetTxPublicKey(tx)) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "release must be signed by record label")
	}
	percentageShares := 0
	rightHolders := make(map[string]struct{})
	rightIds := spec.GetReleaseRightIds(release)
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
	return release, nil
}

func QueryReleaseField(field string, release Data, pub crypto.PublicKey) (interface{}, error) {
	if _, err := ValidateRelease(release, pub); err != nil {
		return nil, err
	}
	switch field {
	case "label":
		recording, err := GetReleaseRecording(release)
		if err != nil {
			return nil, err
		}
		return GetRecordingLabel(recording)
	case "performer":
		recording, err := GetReleaseRecording(release)
		if err != nil {
			return nil, err
		}
		return GetRecordingPerformer(recording)
	case "producer":
		recording, err := GetReleaseRecording(release)
		if err != nil {
			return nil, err
		}
		return GetRecordingProducer(recording)
	case "publication":
		recording, err := GetReleaseRecording(release)
		if err != nil {
			return nil, err
		}
		return GetRecordingPublication(recording)
	case "publishingLicense":
		recording, err := GetReleaseRecording(release)
		if err != nil {
			return nil, err
		}
		return GetRecordingPublishingLicense(recording)
	case "recording":
		return GetReleaseRecording(release)
	case "rights":
		return GetReleaseRights(release)
	default:
		return nil, ErrorAppend(ErrInvalidField, field)
	}
}

func GetReleaseRecording(release Data) (Data, error) {
	recordingId := spec.GetReleaseRecordingId(release)
	tx, err := bigchain.GetTx(recordingId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetReleaseRights(release Data) ([]Data, error) {
	rightIds := spec.GetReleaseRightIds(release)
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

// License

func GetLicenseLicensee(license Data) (Data, error) {
	licenseeId := spec.GetLicenseLicenseeId(license)
	tx, err := bigchain.GetTx(licenseeId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetLicenseLicenser(license Data) (Data, error) {
	licenserId := spec.GetLicenseLicenserId(license)
	tx, err := bigchain.GetTx(licenserId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetLicensePublication(license Data) (Data, error) {
	publicationId := spec.GetLicensePublicationId(license)
	tx, err := bigchain.GetTx(publicationId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
}

func GetLicenseRelease(license Data) (Data, error) {
	releaseId := spec.GetLicenseReleaseId(license)
	tx, err := bigchain.GetTx(releaseId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
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
	license := bigchain.GetTxData(tx)
	if err := spec.ValidPublishingLicense(license); err != nil {
		return nil, err
	}
	pub := bigchain.GetTxPublicKey(tx)
	return ValidatePublishingLicense(license, pub)
}

func ValidatePublishingLicense(license Data, pub crypto.PublicKey) (Data, error) {
	publicationId := spec.GetLicensePublicationId(license)
	publication, err := ValidatePublicationById(publicationId)
	if err != nil {
		return nil, err
	}
	rightHolder := false
	rightIds := spec.GetPublicationRightIds(publication)
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
		return nil, ErrorAppend(ErrCriteriaNotMet, "signer is not a publication right holder")
	}
	licenserId := spec.GetLicenseLicenserId(license)
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
	licenseeId := spec.GetLicenseLicenseeId(license)
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
		return GetLicenseLicensee(license)
	case "licenser":
		return GetLicenseLicenser(license)
	case "publication":
		return GetLicensePublication(license)
	default:
		return nil, ErrorAppend(ErrInvalidField, field)
	}
}

// Release license

func ValidateReleaseLicenseById(licenseId string) (Data, error) {
	tx, err := bigchain.GetTx(licenseId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	license := bigchain.GetTxData(tx)
	if err := spec.ValidReleaseLicense(license); err != nil {
		return nil, err
	}
	pub := bigchain.GetTxPublicKey(tx)
	return ValidateReleaseLicense(license, pub)
}

func ValidateReleaseLicense(license Data, pub crypto.PublicKey) (Data, error) {
	releaseId := spec.GetLicenseReleaseId(license)
	release, err := ValidateReleaseById(releaseId)
	if err != nil {
		return nil, err
	}
	rightHolder := false
	rightIds := spec.GetReleaseRightIds(release)
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
		return nil, ErrorAppend(ErrCriteriaNotMet, "signer is not a release right holder")
	}
	licenserId := spec.GetLicenseLicenserId(license)
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
	licenseeId := spec.GetLicenseLicenseeId(license)
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

func QueryReleaseLicenseField(field string, license Data, pub crypto.PublicKey) (interface{}, error) {
	if _, err := ValidateReleaseLicense(license, pub); err != nil {
		return nil, err
	}
	switch field {
	case "licensee":
		return GetLicenseLicensee(license)
	case "licenser":
		return GetLicenseLicenser(license)
	case "release":
		return GetLicenseRelease(license)
	default:
		return nil, ErrorAppend(ErrInvalidField, field)
	}
}
