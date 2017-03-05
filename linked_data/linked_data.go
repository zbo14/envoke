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
	pub := bigchain.DefaultGetTxSigner(tx)
	if err := ValidateModel(model, pub); err != nil {
		return nil, err
	}
	return model, nil
}

func ValidateModel(model Data, pub crypto.PublicKey) error {
	_type := spec.GetType(model)
	switch _type {
	case spec.AGENT:
		return spec.ValidAgent(model)
	case spec.COMPOSITION:
		return ValidateComposition(model, pub)
	case spec.RECORDING:
		_, err := ValidateRecording(model, pub)
		return err
	case spec.PUBLICATION:
		return ValidatePublication(model, pub)
	case spec.RELEASE:
		return ValidateRelease(model, pub)
	case spec.LICENSE_MECHANICAL:
		return ValidateMechanicalLicense(model, pub)
	case spec.LICENSE_MASTER:
		return ValidateMasterLicense(model, pub)
	default:
		return ErrorAppend(ErrInvalidType, _type)
	}
}

func QueryModelIdField(field, modelId string) (interface{}, error) {
	tx, err := bigchain.GetTx(modelId)
	if err != nil {
		return nil, err
	}
	model := bigchain.GetTxData(tx)
	pub := bigchain.DefaultGetTxSigner(tx)
	result, err := QueryModelField(field, model, pub)
	if err != nil {
		return nil, err
	}
	return result, nil
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
	if err = spec.ValidComposition(composition); err != nil {
		return nil, err
	}
	pub := bigchain.DefaultGetTxSigner(tx)
	if err = ValidateComposition(composition, pub); err != nil {
		return nil, err
	}
	return composition, nil
}

func ValidateComposition(composition Data, pub crypto.PublicKey) error {
	composerId := spec.GetCompositionComposerId(composition)
	tx, err := bigchain.GetTx(composerId)
	if err != nil {
		return err
	}
	if !pub.Equals(bigchain.DefaultGetTxSigner(tx)) {
		return ErrorAppend(ErrCriteriaNotMet, "composition must be signed by composer")
	}
	composer := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(composer); err != nil {
		return err
	}
	publisherId := spec.GetCompositionPublisherId(composition)
	tx, err = bigchain.GetTx(publisherId)
	if err != nil {
		return err
	}
	publisher := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(publisher); err != nil {
		return err
	}
	return nil
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
	if err = spec.ValidPublication(publication); err != nil {
		return nil, err
	}
	pub := bigchain.DefaultGetTxSigner(tx)
	if err = ValidatePublication(publication, pub); err != nil {
		return nil, err
	}
	return publication, nil
}

func ValidatePublication(publication Data, pub crypto.PublicKey) error {
	compositionId := spec.GetPublicationCompositionId(publication)
	composition, err := ValidateCompositionById(compositionId)
	if err != nil {
		return err
	}
	composerId := spec.GetCompositionComposerId(composition)
	publisherId := spec.GetCompositionPublisherId(composition)
	tx, err := bigchain.GetTx(publisherId)
	if err != nil {
		return err
	}
	if !pub.Equals(bigchain.DefaultGetTxSigner(tx)) {
		return ErrorAppend(ErrCriteriaNotMet, "publication must be signed by publisher")
	}
	totalShares := 0
	holderIds := make(map[string]struct{})
	rightIds := make(map[string]struct{})
	assignmentIds := spec.GetPublicationAssignmentIds(publication)
	for _, assignmentId := range assignmentIds {
		assignment, err := ValidateCompositionRightAssignmentById(assignmentId)
		if err != nil {
			return err
		}
		right := spec.GetAssignmentRight(assignment)
		if compositionId != spec.GetRightCompositionId(right) {
			return ErrorAppend(ErrInvalidId, "right has wrong compositionId")
		}
		if composerId != spec.GetAssignmentIssuerId(assignment) {
			return ErrorAppend(ErrCriteriaNotMet, "composer must be assignment issuer")
		}
		holderId := spec.GetAssignmentHolderId(assignment)
		if _, ok := holderIds[holderId]; ok {
			return ErrorAppend(ErrCriteriaNotMet, "holder cannot have multiple assignments")
		}
		holderIds[holderId] = struct{}{}
		rightId := spec.GetAssignmentRightId(assignment)
		if _, ok := rightIds[rightId]; ok {
			return ErrorAppend(ErrCriteriaNotMet, "multiple assignments linked to same right")
		}
		rightIds[rightId] = struct{}{}
		shares := spec.GetRightPercentageShares(right)
		if shares <= 0 {
			return ErrorAppend(ErrCriteriaNotMet, "percentage shares must be greater than 0")
		}
		if totalShares += shares; totalShares > 100 {
			return ErrorAppend(ErrCriteriaNotMet, "total percentage shares cannot exceed 100")
		}
	}
	if totalShares != 100 {
		return ErrorAppend(ErrCriteriaNotMet, "total percentage shares do not equal 100")
	}
	return nil
}

func QueryPublicationField(field string, publication Data, pub crypto.PublicKey) (interface{}, error) {
	if err := ValidatePublication(publication, pub); err != nil {
		return nil, err
	}
	switch field {
	case "assignments":
		return GetPublicationAssignments(publication)
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
	return bigchain.GetTxData(tx), nil
}

func GetPublicationAssignments(publication Data) ([]Data, error) {
	assignmentIds := spec.GetPublicationAssignmentIds(publication)
	assignments := make([]Data, len(assignmentIds))
	for i, assignmentId := range assignmentIds {
		tx, err := bigchain.GetTx(assignmentId)
		if err != nil {
			return nil, err
		}
		assignments[i] = bigchain.GetTxData(tx)
	}
	return assignments, nil
}

// Recording

func ValidateRecordingById(recordingId string) (Data, string, error) {
	tx, err := bigchain.GetTx(recordingId)
	if err != nil {
		return nil, "", err
	}
	recording := bigchain.GetTxData(tx)
	if err := spec.ValidRecording(recording); err != nil {
		return nil, "", err
	}
	pub := bigchain.DefaultGetTxSigner(tx)
	signerId, err := ValidateRecording(recording, pub)
	if err != nil {
		return nil, "", err
	}
	return recording, signerId, nil
}

func ValidateRecording(recording Data, pub crypto.PublicKey) (string, error) {
	labelId := spec.GetRecordingLabelId(recording)
	tx, err := bigchain.GetTx(labelId)
	if err != nil {
		return "", err
	}
	label := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(label); err != nil {
		return "", err
	}
	performerId := spec.GetRecordingPerformerId(recording)
	tx, err = bigchain.GetTx(performerId)
	if err != nil {
		return "", err
	}
	performer := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(performer); err != nil {
		return "", err
	}
	var signed bool
	var signerId string
	if pub.Equals(bigchain.DefaultGetTxSigner(tx)) {
		signed = true
		signerId = performerId
	}
	producerId := spec.GetRecordingProducerId(recording)
	tx, err = bigchain.GetTx(producerId)
	if err != nil {
		return "", err
	}
	producer := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(producer); err != nil {
		return "", err
	}
	if !signed && pub.Equals(bigchain.DefaultGetTxSigner(tx)) {
		signed = true
		signerId = producerId
	}
	if !signed {
		return "", ErrorAppend(ErrCriteriaNotMet, "recording must be signed by performer or producer")
	}
	publicationId := spec.GetRecordingPublicationId(recording)
	publication, err := ValidatePublicationById(publicationId)
	if err != nil {
		return "", err
	}
	assignmentId := spec.GetRecordingAssignmentId(recording)
	if !EmptyStr(assignmentId) {
		holder := false
		assignmentIds := spec.GetPublicationAssignmentIds(publication)
		for i := range assignmentIds {
			if assignmentId == assignmentIds[i] {
				tx, err := bigchain.GetTx(assignmentId)
				if err != nil {
					return "", err
				}
				assignment := bigchain.GetTxData(tx)
				if signerId != spec.GetAssignmentHolderId(assignment) {
					return "", ErrorAppend(ErrCriteriaNotMet, "signer referenced assignment but is not holder")
				}
				holder = true
				break
			}
		}
		if !holder {
			return "", ErrorAppend(ErrCriteriaNotMet, "signer referenced assignment but is not holder")
		}
	}
	return signerId, nil
}

func GetRecordingAssignment(recording Data) (Data, error) {
	assignmentId := spec.GetRecordingAssignmentId(recording)
	tx, err := bigchain.GetTx(assignmentId)
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
	if err = spec.ValidRelease(release); err != nil {
		return nil, err
	}
	pub := bigchain.DefaultGetTxSigner(tx)
	if err = ValidateRelease(release, pub); err != nil {
		return nil, err
	}
	return release, nil
}

func ValidateRelease(release Data, pub crypto.PublicKey) error {
	recordingId := spec.GetReleaseRecordingId(release)
	recording, signerId, err := ValidateRecordingById(recordingId)
	if err != nil {
		return err
	}
	labelId := spec.GetRecordingLabelId(recording)
	tx, err := bigchain.GetTx(labelId)
	if err != nil {
		return err
	}
	if !pub.Equals(bigchain.DefaultGetTxSigner(tx)) {
		return ErrorAppend(ErrCriteriaNotMet, "release must be signed by record label")
	}
	publicationId := spec.GetRecordingPublicationId(recording)
	assignmentId := spec.GetRecordingAssignmentId(recording)
	if EmptyStr(assignmentId) {
		licenseId := spec.GetReleaseLicenseId(release)
		license, err := ValidateMechanicalLicenseById(licenseId)
		if err != nil {
			return err
		}
		if publicationId != spec.GetLicensePublicationId(license) {
			return ErrorAppend(ErrCriteriaNotMet, "wrong mechanical license")
		}
		if labelId != spec.GetLicenseLicenseeId(license) {
			return ErrorAppend(ErrCriteriaNotMet, "label is not licensee of mechanical license")
		}
	}
	totalShares := 0
	holderIds := make(map[string]struct{})
	rightIds := make(map[string]struct{})
	assignmentIds := spec.GetReleaseAssignmentIds(release)
	for _, assignmentId := range assignmentIds {
		assignment, err := ValidateRecordingRightAssignmentById(assignmentId)
		if err != nil {
			return err
		}
		right := spec.GetAssignmentRight(assignment)
		if recordingId != spec.GetRightRecordingId(right) {
			return ErrorAppend(ErrInvalidId, "right has wrong recordingId")
		}
		if signerId != spec.GetAssignmentIssuerId(assignment) {
			return ErrorAppend(ErrCriteriaNotMet, "signer must be assignment issuer")
		}
		holderId := spec.GetAssignmentHolderId(assignment)
		if _, ok := holderIds[holderId]; ok {
			return ErrorAppend(ErrCriteriaNotMet, "holder cannot have multiple assignments")
		}
		holderIds[holderId] = struct{}{}
		rightId := spec.GetAssignmentRightId(assignment)
		if _, ok := rightIds[rightId]; ok {
			return ErrorAppend(ErrCriteriaNotMet, "multiple assignments link to same right")
		}
		rightIds[rightId] = struct{}{}
		shares := spec.GetRightPercentageShares(right)
		if shares <= 0 {
			return ErrorAppend(ErrCriteriaNotMet, "percentage shares must be greater than 0")
		}
		if totalShares += shares; totalShares > 100 {
			return ErrorAppend(ErrCriteriaNotMet, "total percentage shares cannot exceed 100")
		}
	}
	if totalShares != 100 {
		return ErrorAppend(ErrCriteriaNotMet, "total percentage shares do not equal 100")
	}
	return nil
}

func QueryReleaseField(field string, release Data, pub crypto.PublicKey) (interface{}, error) {
	if err := ValidateRelease(release, pub); err != nil {
		return nil, err
	}
	switch field {
	case "assignment":
		recording, err := GetReleaseRecording(release)
		if err != nil {
			return nil, err
		}
		return GetRecordingAssignment(recording)
	case "assignments":
		return GetReleaseAssignments(release)
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
	default:
		return nil, ErrorAppend(ErrInvalidField, field)
	}
}

func GetReleaseAssignments(release Data) ([]Data, error) {
	assignmentIds := spec.GetReleaseAssignmentIds(release)
	assignments := make([]Data, len(assignmentIds))
	for i, assignmentId := range assignmentIds {
		tx, err := bigchain.GetTx(assignmentId)
		if err != nil {
			return nil, err
		}
		assignments[i] = bigchain.GetTxData(tx)
	}
	return assignments, nil
}

func GetReleaseMechanicalLicense(release Data) (Data, error) {
	licenseId := spec.GetReleaseLicenseId(release)
	tx, err := bigchain.GetTx(licenseId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetReleaseRecording(release Data) (Data, error) {
	recordingId := spec.GetReleaseRecordingId(release)
	tx, err := bigchain.GetTx(recordingId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

// License

func GetLicenseAssignment(license Data) (Data, error) {
	assignmentId := spec.GetLicenseAssignmentId(license)
	tx, err := bigchain.GetTx(assignmentId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

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

// Publishing License

func ValidateMechanicalLicenseById(licenseId string) (Data, error) {
	tx, err := bigchain.GetTx(licenseId)
	if err != nil {
		return nil, err
	}
	license := bigchain.GetTxData(tx)
	if err = spec.ValidLicense(license); err != nil {
		return nil, err
	}
	pub := bigchain.DefaultGetTxSigner(tx)
	if err = ValidateMechanicalLicense(license, pub); err != nil {
		return nil, err
	}
	return license, nil
}

func ValidateMechanicalLicense(license Data, pub crypto.PublicKey) error {
	licenserId := spec.GetLicenseLicenserId(license)
	tx, err := bigchain.GetTx(licenserId)
	if err != nil {
		return err
	}
	if !pub.Equals(bigchain.DefaultGetTxSigner(tx)) {
		return ErrorAppend(ErrCriteriaNotMet, "signer must be licenser")
	}
	licenser := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(licenser); err != nil {
		return err
	}
	publicationId := spec.GetLicensePublicationId(license)
	publication, err := ValidatePublicationById(publicationId)
	if err != nil {
		return err
	}
	holder := false
	licenseTerritory := spec.GetTerritory(license)
	assignmentId := spec.GetLicenseAssignmentId(license)
	assignmentIds := spec.GetPublicationAssignmentIds(publication)
	for i := range assignmentIds {
		if assignmentId == assignmentIds[i] {
			tx, err := bigchain.GetTx(assignmentId)
			if err != nil {
				return err
			}
			assignment := bigchain.GetTxData(tx)
			if licenserId != spec.GetAssignmentHolderId(assignment) {
				return ErrorAppend(ErrCriteriaNotMet, "licenser does not hold assignment")
			}
			holder = true
			right := spec.GetAssignmentRight(assignment)
			rightTerritory := spec.GetTerritory(right)
		OUTER:
			for i := range licenseTerritory {
				for j := range rightTerritory {
					if licenseTerritory[i] == rightTerritory[j] {
						rightTerritory = append(rightTerritory[:j], rightTerritory[j+1:]...)
						continue OUTER
					}
				}
				return ErrorAppend(ErrCriteriaNotMet, "license territory not part of right territory")
			}
			break
		}
	}
	if !holder {
		return ErrorAppend(ErrCriteriaNotMet, "licenser does not hold assignment")
	}
	licenseeId := spec.GetLicenseLicenseeId(license)
	tx, err = bigchain.GetTx(licenseeId)
	if err != nil {
		return err
	}
	licensee := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(licensee); err != nil {
		return err
	}
	return nil
}

func QueryMechanicalLicenseField(field string, license Data, pub crypto.PublicKey) (interface{}, error) {
	if err := ValidateMechanicalLicense(license, pub); err != nil {
		return nil, err
	}
	switch field {
	case "assignment":
		return GetLicenseAssignment(license)
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

func ValidateMasterLicenseById(licenseId string) (Data, error) {
	tx, err := bigchain.GetTx(licenseId)
	if err != nil {
		return nil, err
	}
	license := bigchain.GetTxData(tx)
	if err = spec.ValidLicense(license); err != nil {
		return nil, err
	}
	pub := bigchain.DefaultGetTxSigner(tx)
	if err = ValidateMasterLicense(license, pub); err != nil {
		return nil, err
	}
	return license, nil
}

func ValidateMasterLicense(license Data, pub crypto.PublicKey) error {
	licenserId := spec.GetLicenseLicenserId(license)
	tx, err := bigchain.GetTx(licenserId)
	if err != nil {
		return err
	}
	if !pub.Equals(bigchain.DefaultGetTxSigner(tx)) {
		return ErrorAppend(ErrCriteriaNotMet, "signer must be licenser")
	}
	licenser := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(licenser); err != nil {
		return err
	}
	releaseId := spec.GetLicenseReleaseId(license)
	release, err := ValidateReleaseById(releaseId)
	if err != nil {
		return err
	}
	holder := false
	licenseTerritory := spec.GetTerritory(license)
	assignmentId := spec.GetLicenseAssignmentId(license)
	assignmentIds := spec.GetPublicationAssignmentIds(release)
	for i := range assignmentIds {
		if assignmentId == assignmentIds[i] {
			tx, err := bigchain.GetTx(assignmentId)
			if err != nil {
				return err
			}
			assignment := bigchain.GetTxData(tx)
			if licenserId != spec.GetAssignmentHolderId(assignment) {
				return ErrorAppend(ErrCriteriaNotMet, "licenser does not hold assignment")
			}
			holder = true
			right := bigchain.GetTxData(tx)
			rightTerritory := spec.GetTerritory(right)
		OUTER:
			for i := range licenseTerritory {
				for j := range rightTerritory {
					if licenseTerritory[i] == rightTerritory[j] {
						rightTerritory = append(rightTerritory[:j], rightTerritory[j+1:]...)
						continue OUTER
					}
				}
				return ErrorAppend(ErrCriteriaNotMet, "license territory not part of right territory")
			}
			break
		}
	}
	if !holder {
		return ErrorAppend(ErrCriteriaNotMet, "signer does not hold assignment")
	}
	licenseeId := spec.GetLicenseLicenseeId(license)
	tx, err = bigchain.GetTx(licenseeId)
	if err != nil {
		return err
	}
	licensee := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(licensee); err != nil {
		return err
	}
	return nil
}

func QueryMasterLicenseField(field string, license Data, pub crypto.PublicKey) (interface{}, error) {
	if err := ValidateMasterLicense(license, pub); err != nil {
		return nil, err
	}
	switch field {
	case "assignment":
		return GetLicenseAssignment(license)
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

// Assignment

func ValidateCompositionRightAssignmentById(assignmentId string) (Data, error) {
	tx, err := bigchain.GetTx(assignmentId)
	if err != nil {
		return nil, err
	}
	assignment := bigchain.GetTxData(tx)
	if err = spec.ValidAssignment(assignment); err != nil {
		return nil, err
	}
	pub := bigchain.DefaultGetTxSigner(tx)
	if err = ValidateCompositionRightAssignment(assignment, pub); err != nil {
		return nil, err
	}
	return assignment, nil
}

func ValidateCompositionRightAssignment(assignment Data, pub crypto.PublicKey) error {
	holderId := spec.GetAssignmentHolderId(assignment)
	tx, err := bigchain.GetTx(holderId)
	if err != nil {
		return err
	}
	holder := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(holder); err != nil {
		return err
	}
	if !pub.Equals(bigchain.DefaultGetTxSigner(tx)) {
		return ErrInvalidKey
	}
	issuerId := spec.GetAssignmentIssuerId(assignment)
	tx, err = bigchain.GetTx(issuerId)
	if err != nil {
		return err
	}
	issuer := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(issuer); err != nil {
		return err
	}
	issuerPub := bigchain.DefaultGetTxSigner(tx)
	rightId := spec.GetAssignmentRightId(assignment)
	tx, err = bigchain.GetTx(rightId)
	if err != nil {
		return err
	}
	percentageShares := bigchain.GetTxShares(tx)
	right := bigchain.GetTxData(tx)
	if err = spec.ValidCompositionRight(right); err != nil {
		return err
	}
	if !pub.Equals(bigchain.DefaultGetTxRecipient(tx)) {
		return ErrInvalidKey
	}
	if !issuerPub.Equals(bigchain.DefaultGetTxSigner(tx)) {
		return ErrInvalidKey
	}
	right.Set("percentageShares", percentageShares)
	assignment.Set("right", right)
	return nil
}

func ValidateRecordingRightAssignmentById(assignmentId string) (Data, error) {
	tx, err := bigchain.GetTx(assignmentId)
	if err != nil {
		return nil, err
	}
	assignment := bigchain.GetTxData(tx)
	if err = spec.ValidAssignment(assignment); err != nil {
		return nil, err
	}
	pub := bigchain.DefaultGetTxSigner(tx)
	if err = ValidateRecordingRightAssignment(assignment, pub); err != nil {
		return nil, err
	}
	return assignment, nil
}

func ValidateRecordingRightAssignment(assignment Data, pub crypto.PublicKey) error {
	holderId := spec.GetAssignmentHolderId(assignment)
	tx, err := bigchain.GetTx(holderId)
	if err != nil {
		return err
	}
	holder := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(holder); err != nil {
		return err
	}
	if !pub.Equals(bigchain.DefaultGetTxSigner(tx)) {
		return ErrInvalidKey
	}
	issuerId := spec.GetAssignmentIssuerId(assignment)
	tx, err = bigchain.GetTx(issuerId)
	if err != nil {
		return err
	}
	issuer := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(issuer); err != nil {
		return err
	}
	issuerPub := bigchain.DefaultGetTxSigner(tx)
	rightId := spec.GetAssignmentRightId(assignment)
	tx, err = bigchain.GetTx(rightId)
	if err != nil {
		return err
	}
	percentageShares := bigchain.GetTxShares(tx)
	right := bigchain.GetTxData(tx)
	if err = spec.ValidRecordingRight(right); err != nil {
		return err
	}
	if !pub.Equals(bigchain.DefaultGetTxRecipient(tx)) {
		return ErrInvalidKey
	}
	if !issuerPub.Equals(bigchain.DefaultGetTxSigner(tx)) {
		return ErrInvalidKey
	}
	right.Set("percentageShares", percentageShares)
	assignment.Set("right", right)
	return nil
}

// Assignment Holders

func ValidatePublicationAssignmentHolder(assignmentId, holderId, publicationId string) (Data, error) {
	publication, err := ValidatePublicationById(publicationId)
	if err != nil {
		return nil, err
	}
	var assignment Data = nil
	assignmentIds := spec.GetPublicationAssignmentIds(publication)
	for i := range assignmentIds {
		if assignmentId == assignmentIds[i] {
			tx, err := bigchain.GetTx(assignmentId)
			if err != nil {
				return nil, err
			}
			assignment = bigchain.GetTxData(tx)
			if holderId != spec.GetAssignmentHolderId(assignment) {
				return nil, ErrorAppend(ErrCriteriaNotMet, "agent does not hold assignment")
			}
			right := spec.GetAssignmentRight(assignment)
			rightId := spec.GetAssignmentRightId(assignment)
			tx, err = bigchain.GetTx(rightId)
			if err != nil {
				return nil, err
			}
			right.Set("percentageShares", bigchain.GetTxShares(tx))
			assignment.Set("right", right)
			break
		}
	}
	if assignment == nil {
		return nil, ErrorAppend(ErrCriteriaNotMet, "publication does not link to assignment")
	}
	return assignment, nil
}

func ValidateReleaseAssignmentHolder(assignmentId, holderId, releaseId string) (Data, error) {
	release, err := ValidateReleaseById(releaseId)
	if err != nil {
		return nil, err
	}
	var assignment Data = nil
	assignmentIds := spec.GetReleaseAssignmentIds(release)
	for i := range assignmentIds {
		if assignmentId == assignmentIds[i] {
			tx, err := bigchain.GetTx(assignmentId)
			if err != nil {
				return nil, err
			}
			assignment = bigchain.GetTxData(tx)
			if holderId != spec.GetAssignmentHolderId(assignment) {
				return nil, ErrorAppend(ErrCriteriaNotMet, "agent does not hold assignment")
			}
			right := spec.GetAssignmentRight(assignment)
			rightId := spec.GetAssignmentRightId(assignment)
			tx, err = bigchain.GetTx(rightId)
			if err != nil {
				return nil, err
			}
			right.Set("percentageShares", bigchain.GetTxShares(tx))
			assignment.Set("right", right)
			break
		}
	}
	if assignment == nil {
		return nil, ErrorAppend(ErrCriteriaNotMet, "release does not link to assignment")
	}
	return assignment, nil
}

// Transfer

func ValidateCompositionRightTransferById(transferId string) (Data, error) {
	tx, err := bigchain.GetTx(transferId)
	if err != nil {
		return nil, err
	}
	transfer := bigchain.GetTxData(tx)
	if err = spec.ValidCompositionRightTransfer(transfer); err != nil {
		return nil, err
	}
	pub := bigchain.DefaultGetTxSigner(tx)
	if err = ValidateCompositionRightTransfer(pub, transfer); err != nil {
		return nil, err
	}
	return transfer, nil
}

func ValidateCompositionRightTransfer(pub crypto.PublicKey, transfer Data) error {
	recipientId := spec.GetTransferRecipientId(transfer)
	tx, err := bigchain.GetTx(recipientId)
	if err != nil {
		return err
	}
	recipient := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(recipient); err != nil {
		return err
	}
	recipientPub := bigchain.DefaultGetTxSigner(tx)
	senderId := spec.GetTransferSenderId(transfer)
	tx, err = bigchain.GetTx(senderId)
	if err != nil {
		return err
	}
	sender := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(sender); err != nil {
		return err
	}
	senderPub := bigchain.DefaultGetTxSigner(tx)
	if recipientPub.Equals(senderPub) {
		return ErrorAppend(ErrCriteriaNotMet, "recipient key and sender key must be different")
	}
	publicationId := spec.GetTransferPublicationId(transfer)
	publication, err := ValidatePublicationById(publicationId)
	if err != nil {
		return err
	}
	txId := spec.GetTransferTxId(transfer)
	tx, err = bigchain.GetTx(txId)
	if err != nil {
		return err
	}
	if bigchain.TRANSFER != bigchain.GetTxOperation(tx) {
		return ErrorAppend(ErrCriteriaNotMet, "expected TRANSFER tx")
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSigner(tx)) {
		return ErrorAppend(ErrCriteriaNotMet, "sender is not signer of TRANSFER tx")
	}
	n := len(bigchain.GetTxOutputs(tx))
	if n != 1 && n != 2 {
		return ErrorAppend(ErrInvalidSize, "tx outputs must have size 1 or 2")
	}
	if !recipientPub.Equals(bigchain.GetTxRecipient(tx, 0)) {
		return ErrorAppend(ErrCriteriaNotMet, "recipient does not hold primary output of TRANSFER tx")
	}
	if n == 2 {
		if !senderPub.Equals(bigchain.GetTxRecipient(tx, 1)) {
			return ErrorAppend(ErrCriteriaNotMet, "sender does not hold secondary output of TRANSFER tx")
		}
	}
	found := false
	assignmentId := bigchain.GetTxAssetId(tx)
	assignmentIds := spec.GetPublicationAssignmentIds(publication)
	for i := range assignmentIds {
		if assignmentId == assignmentIds[i] {
			found = true
			break
		}
	}
	if !found {
		return ErrorAppend(ErrCriteriaNotMet, "publication does not link to assignment")
	}
	transfer.Set("assignmentId", assignmentId)
	transfer.Set("recipientShares", bigchain.GetTxOutputAmount(tx, 0))
	if n == 2 {
		transfer.Set("senderShares", bigchain.GetTxOutputAmount(tx, 1))
	}
	return nil
}

func ValidateRecordingRightTransferById(transferId string) (Data, error) {
	tx, err := bigchain.GetTx(transferId)
	if err != nil {
		return nil, err
	}
	transfer := bigchain.GetTxData(tx)
	if err = spec.ValidRecordingRightTransfer(transfer); err != nil {
		return nil, err
	}
	pub := bigchain.DefaultGetTxSigner(tx)
	if err = ValidateRecordingRightTransfer(pub, transfer); err != nil {
		return nil, err
	}
	return transfer, nil
}

func ValidateRecordingRightTransfer(pub crypto.PublicKey, transfer Data) error {
	recipientId := spec.GetTransferRecipientId(transfer)
	tx, err := bigchain.GetTx(recipientId)
	if err != nil {
		return err
	}
	recipient := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(recipient); err != nil {
		return err
	}
	recipientPub := bigchain.DefaultGetTxSigner(tx)
	senderId := spec.GetTransferSenderId(transfer)
	tx, err = bigchain.GetTx(senderId)
	if err != nil {
		return err
	}
	sender := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(sender); err != nil {
		return err
	}
	senderPub := bigchain.DefaultGetTxSigner(tx)
	if recipientPub.Equals(senderPub) {
		return ErrorAppend(ErrCriteriaNotMet, "recipient key and sender key must be different")
	}
	releaseId := spec.GetTransferReleaseId(transfer)
	release, err := ValidateReleaseById(releaseId)
	if err != nil {
		return err
	}
	txId := spec.GetTransferTxId(transfer)
	tx, err = bigchain.GetTx(txId)
	if err != nil {
		return err
	}
	if bigchain.TRANSFER != bigchain.GetTxOperation(tx) {
		return ErrorAppend(ErrCriteriaNotMet, "expected TRANSFER tx")
	}
	n := len(bigchain.GetTxOutputs(tx))
	if n != 1 && n != 2 {
		return ErrorAppend(ErrInvalidSize, "outputs must have size 1 or 2")
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSigner(tx)) {
		return ErrorAppend(ErrCriteriaNotMet, "sender is not signer of TRANSFER tx")
	}
	if !recipientPub.Equals(bigchain.GetTxRecipient(tx, 0)) {
		return ErrorAppend(ErrCriteriaNotMet, "recipient does not hold primary output of TRANSFER tx")
	}
	if n == 2 {
		if !senderPub.Equals(bigchain.GetTxRecipient(tx, 1)) {
			return ErrorAppend(ErrCriteriaNotMet, "sender does not hold secondary output of TRANSFER tx")
		}
	}
	found := false
	assignmentId := bigchain.GetTxAssetId(tx)
	assignmentIds := spec.GetReleaseAssignmentIds(release)
	for i := range assignmentIds {
		if assignmentId == assignmentIds[i] {
			found = true
			break
		}
	}
	if !found {
		return ErrorAppend(ErrCriteriaNotMet, "release does not link to assignment")
	}
	transfer.Set("assignmentId", assignmentId)
	transfer.Set("recipientShares", bigchain.GetTxOutputAmount(tx, 0))
	if n == 2 {
		transfer.Set("senderShares", bigchain.GetTxOutputAmount(tx, 1))
	}
	return nil
}
