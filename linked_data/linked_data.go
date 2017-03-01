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
	model := bigchain.GetTxData(tx)
	pub := bigchain.GetTxSigner(tx)
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
	case spec.LICENSE_MECHANICAL:
		model, err = ValidateMechanicalLicense(model, pub)
	case spec.LICENSE_MASTER:
		model, err = ValidateMasterLicense(model, pub)
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
	model := bigchain.GetTxData(tx)
	pub := bigchain.GetTxSigner(tx)
	return QueryModelField(field, model, pub)
}

func QueryModelField(field string, model Data, pub crypto.PublicKey) (interface{}, error) {
	_type := spec.GetType(model)
	switch _type {
	case spec.PUBLICATION:
		return QueryPublicationField(field, model, pub)
	case spec.RELEASE:
		return QueryReleaseField(field, model, pub)
	case spec.LICENSE_MECHANICAL:
		return QueryMechanicalLicenseField(field, model, pub)
	case spec.LICENSE_MASTER:
		return QueryMasterLicenseField(field, model, pub)
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
	composition := bigchain.GetTxData(tx)
	if err := spec.ValidComposition(composition); err != nil {
		return nil, err
	}
	pub := bigchain.GetTxSigner(tx)
	return ValidateComposition(composition, pub)
}

func ValidateComposition(composition Data, pub crypto.PublicKey) (Data, error) {
	composerId := spec.GetCompositionComposerId(composition)
	tx, err := bigchain.GetTx(composerId)
	if err != nil {
		return nil, err
	}
	if !pub.Equals(bigchain.GetTxSigner(tx)) {
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
	return bigchain.GetTxData(tx), nil
}

func GetCompositionPublisher(composition Data) (Data, error) {
	publisherId := spec.GetCompositionPublisherId(composition)
	tx, err := bigchain.GetTx(publisherId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func ValidatePublicationById(publicationId string) (Data, error) {
	tx, err := bigchain.GetTx(publicationId)
	if err != nil {
		return nil, err
	}
	publication := bigchain.GetTxData(tx)
	if err := spec.ValidPublication(publication); err != nil {
		return nil, err
	}
	pub := bigchain.GetTxSigner(tx)
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
	if !pub.Equals(bigchain.GetTxSigner(tx)) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "publication must be signed by publisher")
	}
	composerId := spec.GetCompositionComposerId(composition)
	tx, err = bigchain.GetTx(composerId)
	if err != nil {
		return nil, err
	}
	composerPub := bigchain.GetTxSigner(tx)
	percentageShares := 0
	rightHolders := make(map[string]struct{})
	rightIds := spec.GetCompositionRightIds(publication)
	for _, rightId := range rightIds {
		tx, err = bigchain.GetTx(rightId)
		if err != nil {
			return nil, err
		}
		if !composerPub.Equals(bigchain.GetTxSigner(tx)) {
			return nil, ErrorAppend(ErrCriteriaNotMet, "composition right must be signed by composer")
		}
		if _, ok := rightHolders[bigchain.GetTxRecipient(tx).String()]; ok {
			return nil, ErrorAppend(ErrCriteriaNotMet, "rightHolder cannot have multiple rights to composition")
		}
		right := bigchain.GetTxData(tx)
		if err = spec.ValidRight(right); err != nil {
			return nil, err
		}
		if compositionId != spec.GetRightCompositionId(right) {
			return nil, ErrorAppend(ErrInvalidId, "wrong compositionId")
		}
		percentageShares += bigchain.GetTxShares(tx)
		if percentageShares > 100 {
			return nil, ErrorAppend(ErrCriteriaNotMet, "percentage shares cannot exceed 100")
		}
	}
	if percentageShares != 100 {
		return nil, ErrorAppend(ErrCriteriaNotMet, "total percentage shares do not equal 100")
	}
	return publication, nil
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
	rightIds := spec.GetCompositionRightIds(publication)
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
	pub := bigchain.GetTxSigner(tx)
	return ValidateRecording(recording, pub)
}

func ValidateRecording(recording Data, pub crypto.PublicKey) (Data, error) {
	labelId := spec.GetRecordingLabelId(recording)
	tx, err := bigchain.GetTx(labelId)
	if err != nil {
		return nil, err
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
	performer := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(performer); err != nil {
		return nil, err
	}
	signed := false
	if pub.Equals(bigchain.GetTxSigner(tx)) {
		signed = true
	}
	producerId := spec.GetRecordingProducerId(recording)
	tx, err = bigchain.GetTx(producerId)
	if err != nil {
		return nil, err
	}
	producer := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(producer); err != nil {
		return nil, err
	}
	if !signed && pub.Equals(bigchain.GetTxSigner(tx)) {
		signed = true
	}
	if !signed {
		return nil, ErrorAppend(ErrCriteriaNotMet, "recording must be signed by performer or producer")
	}
	publicationId := spec.GetRecordingPublicationId(recording)
	publication, err := ValidatePublicationById(publicationId)
	if err != nil {
		return nil, err
	}
	rightId := spec.GetRecordingCompositionRightId(recording)
	if !EmptyStr(rightId) {
		rightHolder := false
		rightIds := spec.GetCompositionRightIds(publication)
		for i := range rightIds {
			if rightId == rightIds[i] {
				if !pub.Equals(bigchain.GetTxRecipient(tx)) {
					return nil, ErrorAppend(ErrCriteriaNotMet, "signer referenced composition right but is not rightHolder")
				}
				rightHolder = true
				break
			}
		}
		if !rightHolder {
			return nil, ErrorAppend(ErrCriteriaNotMet, "signer referenced composition right but is not rightHolder")
		}
	}
	return recording, nil
}

func GetRecordingCompositionRight(recording Data) (Data, error) {
	compositionRightId := spec.GetRecordingCompositionRightId(recording)
	tx, err := bigchain.GetTx(compositionRightId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetRecordingLabel(recording Data) (Data, error) {
	labelId := spec.GetRecordingLabelId(recording)
	tx, err := bigchain.GetTx(labelId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetRecordingPerformer(recording Data) (Data, error) {
	performerId := spec.GetRecordingPerformerId(recording)
	tx, err := bigchain.GetTx(performerId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetRecordingProducer(recording Data) (Data, error) {
	producerId := spec.GetRecordingProducerId(recording)
	tx, err := bigchain.GetTx(producerId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetRecordingPublication(recording Data) (Data, error) {
	publicationId := spec.GetRecordingPublicationId(recording)
	tx, err := bigchain.GetTx(publicationId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func ValidateReleaseById(releaseId string) (Data, error) {
	tx, err := bigchain.GetTx(releaseId)
	if err != nil {
		return nil, err
	}
	release := bigchain.GetTxData(tx)
	if err := spec.ValidRelease(release); err != nil {
		return nil, err
	}
	pub := bigchain.GetTxSigner(tx)
	return ValidateRelease(release, pub)
}

func ValidateRelease(release Data, pub crypto.PublicKey) (Data, error) {
	recordingId := spec.GetReleaseRecordingId(release)
	recording, err := ValidateRecordingById(recordingId)
	if err != nil {
		return nil, err
	}
	recordingSignerPub := bigchain.GetTxSigner(recording)
	labelId := spec.GetRecordingLabelId(recording)
	tx, err := bigchain.GetTx(labelId)
	if err != nil {
		return nil, err
	}
	if !pub.Equals(bigchain.GetTxSigner(tx)) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "release must be signed by record label")
	}
	publicationId := spec.GetRecordingPublicationId(recording)
	rightId := spec.GetRecordingCompositionRightId(recording)
	if EmptyStr(rightId) {
		licenseId := spec.GetReleaseLicenseId(release)
		license, err := ValidateMechanicalLicenseById(licenseId)
		if err != nil {
			return nil, err
		}
		if publicationId != spec.GetLicensePublicationId(license) {
			return nil, ErrorAppend(ErrCriteriaNotMet, "wrong publishing license")
		}
		if labelId != spec.GetLicenseLicenseeId(license) {
			return nil, ErrorAppend(ErrCriteriaNotMet, "label is not licensee of mechanical license")
		}
	}
	percentageShares := 0
	rightHolders := make(map[string]struct{})
	rightIds := spec.GetRecordingRightIds(release)
	for _, rightId := range rightIds {
		tx, err = bigchain.GetTx(rightId)
		if err != nil {
			return nil, err
		}
		if !recordingSignerPub.Equals(bigchain.GetTxSigner(tx)) {
			return nil, ErrorAppend(ErrCriteriaNotMet, "recording right must be signed by recording signer")
		}
		if _, ok := rightHolders[bigchain.GetTxRecipient(tx).String()]; ok {
			return nil, ErrorAppend(ErrCriteriaNotMet, "rightHolder cannot have multiple rights to recording")
		}
		right := bigchain.GetTxData(tx)
		if err = spec.ValidRight(right); err != nil {
			return nil, err
		}
		if recordingId != spec.GetRightRecordingId(right) {
			return nil, ErrorAppend(ErrInvalidId, "wrong recordingId")
		}
		shares := bigchain.GetTxShares(right)
		if shares <= 0 {
			return nil, ErrorAppend(ErrCriteriaNotMet, "percentage shares must be greater than 0")
		}
		if percentageShares += shares; percentageShares > 100 {
			return nil, ErrorAppend(ErrCriteriaNotMet, "total percentage shares cannot exceed 100")
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
	case "composition_right":
		recording, err := GetReleaseRecording(release)
		if err != nil {
			return nil, err
		}
		return GetRecordingCompositionRight(recording)
	case "label":
		recording, err := GetReleaseRecording(release)
		if err != nil {
			return nil, err
		}
		return GetRecordingLabel(recording)
	case "mechanical_license":
		return GetReleaseMechanicalLicense(release)
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
	case "recording":
		return GetReleaseRecording(release)
	case "rights":
		return GetReleaseRights(release)
	default:
		return nil, ErrorAppend(ErrInvalidField, field)
	}
}

func GetReleaseMechanicalLicense(release Data) (Data, error) {
	licenseId := spec.GetReleaseLicenseId(release)
	tx, err := bigchain.GetTx(licenseId)
	if err != nil {
		return nil, err
	}
	if !bigchain.FulfilledTx(tx) {
		return nil, ErrInvalidFulfillment
	}
	return bigchain.GetTxData(tx), nil
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
	rightIds := spec.GetRecordingRightIds(release)
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
	return bigchain.GetTxData(tx), nil
}

func GetLicenseLicenser(license Data) (Data, error) {
	licenserId := spec.GetLicenseLicenserId(license)
	tx, err := bigchain.GetTx(licenserId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetLicensePublication(license Data) (Data, error) {
	publicationId := spec.GetLicensePublicationId(license)
	tx, err := bigchain.GetTx(publicationId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetLicenseRelease(license Data) (Data, error) {
	releaseId := spec.GetLicenseReleaseId(license)
	tx, err := bigchain.GetTx(releaseId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetLicenseRight(license Data) (Data, error) {
	rightId := spec.GetLicenseRightId(license)
	tx, err := bigchain.GetTx(rightId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

// Right Holders

func ValidateCompositionRightHolderById(publicationId, rightHolderId, rightId string) (Data, error) {
	tx, err := bigchain.GetTx(rightHolderId)
	if err != nil {
		return nil, err
	}
	rightHolder := bigchain.GetTxData(tx)
	if err := spec.ValidAgent(rightHolder); err != nil {
		return nil, err
	}
	pub := bigchain.GetTxSigner(tx)
	return ValidateCompositionRightHolder(pub, publicationId, rightId)
}

func ValidateCompositionRightHolder(pub crypto.PublicKey, publicationId, rightId string) (Data, error) {
	publication, err := ValidatePublicationById(publicationId)
	if err != nil {
		return nil, err
	}
	rightIds := make(map[string]struct{})
	for _, id := range spec.GetCompositionRightIds(publication) {
		rightIds[id] = struct{}{}
	}
	tx, err := bigchain.GetTx(rightId)
	if err != nil {
		return nil, err
	}
	if !pub.Equals(bigchain.GetTxRecipient(tx)) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "agent does not hold this composition right")
	}
	if id := bigchain.GetTxAssetId(tx); spec.MatchId(id) {
		tx, err = bigchain.GetTx(id)
		if err != nil {
			return nil, err
		}
		rightId = id
	}
	if _, ok := rightIds[rightId]; !ok {
		return nil, ErrorAppend(ErrCriteriaNotMet, "publication does not reference this composition right")
	}
	return bigchain.GetTxData(tx), nil
}

func ValidateRecordingRightHolderById(releaseId, rightHolderId, rightId string) (Data, error) {
	tx, err := bigchain.GetTx(rightHolderId)
	if err != nil {
		return nil, err
	}
	rightHolder := bigchain.GetTxData(tx)
	if err := spec.ValidAgent(rightHolder); err != nil {
		return nil, err
	}
	pub := bigchain.GetTxSigner(tx)
	return ValidateRecordingRightHolder(pub, releaseId, rightId)
}

func ValidateRecordingRightHolder(pub crypto.PublicKey, releaseId, rightId string) (Data, error) {
	release, err := ValidateReleaseById(releaseId)
	if err != nil {
		return nil, err
	}
	rightIds := make(map[string]struct{})
	for _, id := range spec.GetRecordingRightIds(release) {
		rightIds[id] = struct{}{}
	}
	tx, err := bigchain.GetTx(rightId)
	if err != nil {
		return nil, err
	}
	if !pub.Equals(bigchain.GetTxRecipient(tx)) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "agent does not hold this recording right")
	}
	if id := bigchain.GetTxAssetId(tx); spec.MatchId(id) {
		tx, err = bigchain.GetTx(id)
		if err != nil {
			return nil, err
		}
		rightId = id
	}
	if _, ok := rightIds[rightId]; !ok {
		return nil, ErrorAppend(ErrCriteriaNotMet, "release does not reference this recording right")
	}
	return bigchain.GetTxData(tx), nil
}

// Publishing License

func ValidateMechanicalLicenseById(licenseId string) (Data, error) {
	tx, err := bigchain.GetTx(licenseId)
	if err != nil {
		return nil, err
	}
	license := bigchain.GetTxData(tx)
	if err := spec.ValidLicense(license); err != nil {
		return nil, err
	}
	pub := bigchain.GetTxSigner(tx)
	return ValidateMechanicalLicense(license, pub)
}

func ValidateMechanicalLicense(license Data, pub crypto.PublicKey) (Data, error) {
	licenserId := spec.GetLicenseLicenserId(license)
	tx, err := bigchain.GetTx(licenserId)
	if err != nil {
		return nil, err
	}
	if !pub.Equals(bigchain.GetTxSigner(tx)) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "signer must be licenser")
	}
	licenser := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(licenser); err != nil {
		return nil, err
	}
	publicationId := spec.GetLicensePublicationId(license)
	publication, err := ValidatePublicationById(publicationId)
	if err != nil {
		return nil, err
	}
	licenseTerritory := spec.GetTerritory(license)
	rightHolder := false
	rightId := spec.GetLicenseRightId(license)
	rightIds := spec.GetCompositionRightIds(publication)
	for i := range rightIds {
		if rightId == rightIds[i] {
			tx, err := bigchain.GetTx(rightId)
			if err != nil {
				return nil, err
			}
			if !pub.Equals(bigchain.GetTxRecipient(tx)) {
				return nil, ErrorAppend(ErrCriteriaNotMet, "licenser does not hold this composition right")
			}
			right := bigchain.GetTxData(tx)
			rightHolder = true
			rightTerritory := spec.GetTerritory(right)
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
			break
		}
	}
	if !rightHolder {
		return nil, ErrorAppend(ErrCriteriaNotMet, "licenser does not hold a composition right")
	}
	licenseeId := spec.GetLicenseLicenseeId(license)
	tx, err = bigchain.GetTx(licenseeId)
	if err != nil {
		return nil, err
	}
	licensee := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(licensee); err != nil {
		return nil, err
	}
	return license, nil
}

func QueryMechanicalLicenseField(field string, license Data, pub crypto.PublicKey) (interface{}, error) {
	if _, err := ValidateMechanicalLicense(license, pub); err != nil {
		return nil, err
	}
	switch field {
	case "licensee":
		return GetLicenseLicensee(license)
	case "licenser":
		return GetLicenseLicenser(license)
	case "publication":
		return GetLicensePublication(license)
	case "right":
		return GetLicenseRight(license)
	default:
		return nil, ErrorAppend(ErrInvalidField, field)
	}
}

// Release license

func ValidateMasterLicenseById(licenseId string) (Data, error) {
	tx, err := bigchain.GetTx(licenseId)
	if err != nil {
		return nil, err
	}
	license := bigchain.GetTxData(tx)
	if err := spec.ValidLicense(license); err != nil {
		return nil, err
	}
	pub := bigchain.GetTxSigner(tx)
	return ValidateMasterLicense(license, pub)
}

func ValidateMasterLicense(license Data, pub crypto.PublicKey) (Data, error) {
	licenserId := spec.GetLicenseLicenserId(license)
	tx, err := bigchain.GetTx(licenserId)
	if err != nil {
		return nil, err
	}
	if !pub.Equals(bigchain.GetTxSigner(tx)) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "signer must be licenser")
	}
	licenser := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(licenser); err != nil {
		return nil, err
	}
	releaseId := spec.GetLicenseReleaseId(license)
	release, err := ValidateReleaseById(releaseId)
	if err != nil {
		return nil, err
	}
	licenseTerritory := spec.GetTerritory(license)
	rightHolder := false
	rightId := spec.GetLicenseRightId(license)
	rightIds := spec.GetRecordingRightIds(release)
	for i := range rightIds {
		if rightId == rightIds[i] {
			tx, err := bigchain.GetTx(rightId)
			if err != nil {
				return nil, err
			}
			if !pub.Equals(bigchain.GetTxRecipient(tx)) {
				return nil, ErrorAppend(ErrCriteriaNotMet, "licenser does not hold this recording right")
			}
			right := bigchain.GetTxData(tx)
			rightHolder = true
			rightTerritory := spec.GetTerritory(right)
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
			break
		}
	}
	if !rightHolder {
		return nil, ErrorAppend(ErrCriteriaNotMet, "signer is not a recording right holder")
	}
	licenseeId := spec.GetLicenseLicenseeId(license)
	tx, err = bigchain.GetTx(licenseeId)
	if err != nil {
		return nil, err
	}
	licensee := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(licensee); err != nil {
		return nil, err
	}
	return license, nil
}

func QueryMasterLicenseField(field string, license Data, pub crypto.PublicKey) (interface{}, error) {
	if _, err := ValidateMasterLicense(license, pub); err != nil {
		return nil, err
	}
	switch field {
	case "licensee":
		return GetLicenseLicensee(license)
	case "licenser":
		return GetLicenseLicenser(license)
	case "release":
		return GetLicenseRelease(license)
	case "right":
		return GetLicenseRight(license)
	default:
		return nil, ErrorAppend(ErrInvalidField, field)
	}
}
