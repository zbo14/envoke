package spec

import (
	"net/url"
	"time"

	. "github.com/zbo14/envoke/common"
)

const (
	AGENT              = "agent"
	COMPOSITION        = "composition"
	RECORDING          = "recording"
	RIGHT              = "right"
	LICENSE_PUBLISHING = "publishing_license"
	LICENSE_RECORDING  = "recording_license"
	// LICENSE_BLANKET

	LICENSE_TYPE_MASTER          = "master_license"
	LICENSE_TYPE_MECHANICAL      = "mechanical_license"
	LICENSE_TYPE_SYNCHRONIZATION = "synchronization_license"

	INSTANCE_SIZE    = 2
	AGENT_SIZE       = 4
	COMPOSITION_SIZE = 5
	RECORDING_SIZE   = 7
	RIGHT_SIZE       = 5
	LICENSE_SIZE     = 7

	EMAIL_REGEX           = `(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$)`
	FINGERPRINT_STD_REGEX = `^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$` // base64 std
	FINGERPRINT_URL_REGEX = `^(?:[A-Za-z0-9-_]{4})*(?:[A-Za-z0-9-_]{2}==|[A-Za-z0-9-_]{3})?$`  // base64 url-safe
	ID_REGEX              = `^[A-Fa-f0-9]{64}$`                                                // hex
	PUBKEY_REGEX          = `^[1-9A-HJ-NP-Za-km-z]{43,44}$`                                    // base58
	SIGNATURE_REGEX       = `^[1-9A-HJ-NP-Za-km-z]{87,88}$`                                    // base58
)

func MatchEmail(email string) bool {
	return MatchStr(EMAIL_REGEX, email)
}

func MatchFingerprint(fingerprint string) bool {
	return MatchStr(FINGERPRINT_URL_REGEX, fingerprint)
}

func MatchId(id string) bool {
	return MatchStr(ID_REGEX, id)
}

// Instance

func NewInstance(_type string) Data {
	return Data{
		"time": FormatInt64(Timestamp(), 10),
		"type": _type,
	}
}

func GetInstanceTime(instance Data) int64 {
	x, err := ParseInt64(instance.GetStr("time"), 10)
	if err != nil {
		return 0
	}
	return x
}

func GetInstanceType(instance Data) string {
	return instance.GetStr("type")
}

func GetInstance(thing Data) (instance Data) {
	if err := ValidInstance(thing); err == nil {
		return thing
	}
	if instance = thing.GetData("instance"); instance == nil {
		instance = AssertMapData(thing["instance"])
	}
	return instance
}

func GetType(thing Data) string {
	return GetInstanceType(GetInstance(thing))
}

func HasType(thing Data, Type string) bool {
	return GetType(thing) == Type
}

func ValidInstance(instance Data) error {
	time := GetInstanceTime(instance)
	if time > Timestamp() || time == 0 {
		return ErrInvalidTime
	}
	_type := GetInstanceType(instance)
	switch _type {
	case
		AGENT,
		COMPOSITION,
		RECORDING,
		RIGHT,
		LICENSE_PUBLISHING,
		LICENSE_RECORDING:
		// Ok..
	default:
		return ErrorAppend(ErrInvalidType, _type)
	}
	if len(instance) != INSTANCE_SIZE {
		return ErrorAppend(ErrInvalidSize, "instance")
	}
	return nil
}

// Agent

func NewAgent(email, name, socialMedia string) Data {
	return Data{
		"email":       email,
		"instance":    NewInstance(AGENT),
		"name":        name,
		"socialMedia": socialMedia,
	}
}

func GetAgentEmail(agent Data) string {
	return agent.GetStr("email")
}

func GetAgentName(agent Data) string {
	return agent.GetStr("name")
}

func GetAgentSocialMediaStr(agent Data) string {
	return agent.GetStr("socialMedia")
}

func GetAgentSocialMedia(agent Data) *url.URL {
	return MustParseUrl(GetAgentSocialMediaStr(agent))
}

func ValidAgent(agent Data) error {
	instance := GetInstance(agent)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(agent, AGENT) {
		return ErrorAppend(ErrInvalidType, GetType(agent))
	}
	email := GetAgentEmail(agent)
	if !MatchEmail(email) {
		return ErrorAppend(ErrInvalidEmail, email)
	}
	name := GetAgentName(agent)
	if EmptyStr(name) {
		return ErrorAppend(ErrEmptyStr, "name")
	}
	socialMedia := GetAgentSocialMediaStr(agent)
	if !MatchUrlRelaxed(socialMedia) {
		return ErrorAppend(ErrInvalidUrl, "social media")
	}
	if len(agent) != AGENT_SIZE {
		return ErrorAppend(ErrInvalidSize, AGENT)
	}
	return nil
}

// Composition

func NewComposition(composerId, publisherId string, rights []Data, title string) Data {
	return Data{
		"composerId":  composerId,
		"instance":    NewInstance(COMPOSITION),
		"publisherId": publisherId,
		"rights":      rights,
		"title":       title,
	}
}

func GetCompositionComposer(composition Data) string {
	return composition.GetStr("composerId")
}

func GetCompositionRights(composition Data) []Data {
	slice := composition.GetInterfaceSlice("rights")
	rights := make([]Data, len(slice))
	for i, s := range slice {
		rights[i] = AssertMapData(s)
	}
	return rights
}

func GetCompositionPublisher(composition Data) string {
	return composition.GetStr("publisherId")
}

func GetCompositionTitle(composition Data) string {
	return composition.GetStr("title")
}

func ValidComposition(composition Data) error {
	instance := GetInstance(composition)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(composition, COMPOSITION) {
		return ErrorAppend(ErrInvalidType, GetType(composition))
	}
	composerId := GetCompositionComposer(composition)
	if !MatchId(composerId) {
		return ErrorAppend(ErrInvalidId, "composerId")
	}
	percentageShares := 0
	rightHolderIds := make(map[string]struct{})
	rights := GetCompositionRights(composition)
	for _, right := range rights {
		if err := ValidRight(AssertData(right)); err != nil {
			return err
		}
		percentageShares += GetRightPercentageShares(AssertData(right))
		if percentageShares > 100 {
			return ErrorAppend(ErrCriteriaNotMet, "percentage shares cannot exceed 100")
		}
		rightHolderId := GetRightHolder(AssertData(right))
		if _, ok := rightHolderIds[rightHolderId]; ok {
			return ErrorAppend(ErrInvalidId, "agent cannot hold multiple rights to recording")
		}
		rightHolderIds[rightHolderId] = struct{}{}
	}
	publisherId := GetCompositionPublisher(composition)
	if !MatchId(publisherId) {
		return ErrorAppend(ErrInvalidId, "publisherId")
	}
	title := GetCompositionTitle(composition)
	if EmptyStr(title) {
		return ErrorAppend(ErrEmptyStr, "title")
	}
	if len(composition) != COMPOSITION_SIZE {
		return ErrorAppend(ErrInvalidSize, COMPOSITION)
	}
	return nil
}

func NewRecording(compositionId, labelId, performerId, producerId, publishingLicenseId string, rights []Data) Data {
	recording := Data{
		"compositionId": compositionId,
		"instance":      NewInstance(RECORDING),
		"labelId":       labelId,
		"performerId":   performerId,
		"producerId":    producerId,
		"rights":        rights,
	}
	if publishingLicenseId != "" {
		recording.Set("publishingLicenseId", publishingLicenseId)
	}
	return recording
}

func GetRecordingComposition(recording Data) string {
	return recording.GetStr("compositionId")
}

func GetRecordingLabel(recording Data) string {
	return recording.GetStr("labelId")
}

func GetRecordingPublishingLicense(recording Data) string {
	return recording.GetStr("publishingLicenseId")
}

func GetRecordingPerformer(recording Data) string {
	return recording.GetStr("performerId")
}

func GetRecordingProducer(recording Data) string {
	return recording.GetStr("producerId")
}

func GetRecordingRights(recording Data) []Data {
	slice := recording.GetInterfaceSlice("rights")
	rights := make([]Data, len(slice))
	for i, s := range slice {
		rights[i] = AssertMapData(s)
	}
	return rights
}

func ValidRecording(recording Data) error {
	instance := GetInstance(recording)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(recording, RECORDING) {
		return ErrorAppend(ErrInvalidType, GetType(recording))
	}
	compositionId := GetRecordingComposition(recording)
	if !MatchId(compositionId) {
		return ErrorAppend(ErrInvalidId, "compositionId")
	}
	// TODO: better fingerprint validation?
	// fingerprint := GetRecordingFingerprint(recording)
	// if !MatchFingerprint(fingerprint) {
	//	return ErrorAppend(ErrInvalidFingerprint, "fingerprint")
	// }
	labelId := GetRecordingLabel(recording)
	if !MatchId(labelId) {
		return ErrorAppend(ErrInvalidId, "labelId")
	}
	performerId := GetRecordingPerformer(recording)
	if !MatchId(performerId) {
		return ErrorAppend(ErrInvalidId, "performerId")
	}
	percentageShares := 0
	rightHolderIds := make(map[string]struct{})
	rights := GetRecordingRights(recording)
	for _, right := range rights {
		if err := ValidRight(right); err != nil {
			return err
		}
		percentageShares += GetRightPercentageShares(right)
		if percentageShares > 100 {
			return ErrorAppend(ErrCriteriaNotMet, "percentage shares cannot exceed 100")
		}
		rightHolderId := GetRightHolder(right)
		if _, ok := rightHolderIds[rightHolderId]; ok {
			return ErrorAppend(ErrInvalidId, "agent cannot hold multiple rights to recording")
		}
		rightHolderIds[rightHolderId] = struct{}{}
	}
	if percentageShares != 100 {
		return ErrorAppend(ErrCriteriaNotMet, "total percentage shares does not equal 100")
	}
	publishingLicenseId := GetRecordingPublishingLicense(recording)
	if !EmptyStr(publishingLicenseId) {
		if !MatchId(publishingLicenseId) {
			return ErrorAppend(ErrCriteriaNotMet, "Recording must have composition right or publishing license")
		}
	}
	if len(recording) != RECORDING_SIZE {
		return ErrorAppend(ErrInvalidSize, RECORDING)
	}
	return nil
}

// Right

func NewRight(percentageShares, rightHolderId, validFrom, validTo string) Data {
	return Data{
		// should we include context, usage?
		"instance":         NewInstance(RIGHT),
		"percentageShares": percentageShares,
		"rightHolderId":    rightHolderId,
		"validFrom":        validFrom,
		"validTo":          validTo,
	}
}

func GetRightPercentageShares(right Data) int {
	return right.GetStrInt("percentageShares")
}

func GetRightHolder(right Data) string {
	return right.GetStr("rightHolderId")
}

func GetRightValidFrom(right Data) time.Time {
	return MustParseDateStr(right.GetStr("validFrom"))
}

func GetRightValidTo(right Data) time.Time {
	return MustParseDateStr(right.GetStr("validTo"))
}

func ValidRight(right Data) error {
	instance := GetInstance(right)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(right, RIGHT) {
		return ErrorAppend(ErrInvalidType, GetType(right))
	}
	// TODO: validate context, exclusive
	percentageShares := GetRightPercentageShares(right)
	if percentageShares <= 0 || percentageShares > 100 {
		return ErrorAppend(ErrCriteriaNotMet, "percentage shares must be greater than 0 and less than 100")
	}
	rightHolderId := GetRightHolder(right)
	if !MatchId(rightHolderId) {
		return ErrorAppend(ErrInvalidId, "rightHolderId")
	}
	// TODO: validate usage
	validFrom := GetRightValidFrom(right)
	validTo := GetRightValidTo(right)
	if !validFrom.Before(validTo) {
		return ErrorAppend(ErrInvalidTime, "range")
	}
	if validTo.Before(Now()) {
		return ErrorAppend(ErrInvalidTime, "expired")
	}
	if len(right) != RIGHT_SIZE {
		return ErrorAppend(ErrInvalidSize, RIGHT)
	}
	return nil
}

// TODO: right transfers

// License

func NewLicense(licenseeId, licenserId, licenseType, _type, validFrom, validTo string) Data {
	return Data{
		"instance":    NewInstance(_type),
		"licenseeId":  licenseeId,
		"licenserId":  licenserId,
		"licenseType": licenseType,
		"validFrom":   validFrom,
		"validTo":     validTo,
	}
}

func NewPublishingLicense(compositionId, licenseeId, licenserId, licenseType, validFrom, validTo string) Data {
	license := NewLicense(licenseeId, licenserId, licenseType, LICENSE_PUBLISHING, validFrom, validTo)
	license.Set("compositionId", compositionId)
	return license
}

func NewRecordingLicense(licenseeId, licenserId, licenseType, recordingId, validFrom, validTo string) Data {
	license := NewLicense(licenseeId, licenserId, licenseType, LICENSE_RECORDING, validFrom, validTo)
	license.Set("recordingId", recordingId)
	return license
}

func GetLicenseComposition(license Data) string {
	return license.GetStr("compositionId")
}

func GetLicenseLicensee(license Data) string {
	return license.GetStr("licenseeId")
}

func GetLicenseLicenser(license Data) string {
	return license.GetStr("licenserId")
}

func GetLicenseRecording(license Data) string {
	return license.GetStr("recordingId")
}

func GetLicenseType(license Data) string {
	return license.GetStr("licenseType")
}

func GetLicenseValidFrom(license Data) time.Time {
	return MustParseDateStr(license.GetStr("validFrom"))
}

func GetLicenseValidTo(license Data) time.Time {
	return MustParseDateStr(license.GetStr("validTo"))
}

func ValidLicense(license Data, _type string) error {
	instance := GetInstance(license)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(license, _type) {
		return ErrorAppend(ErrInvalidType, _type)
	}
	licenseeId := GetLicenseLicensee(license)
	if !MatchId(licenseeId) {
		return ErrorAppend(ErrInvalidId, "licenseeId")
	}
	licenseType := GetLicenseType(license)
	switch licenseType {
	case
		LICENSE_TYPE_MASTER,
		LICENSE_TYPE_MECHANICAL,
		LICENSE_TYPE_SYNCHRONIZATION:
		//..
	default:
		return ErrorAppend(ErrInvalidType, licenseType)
	}
	validFrom := GetLicenseValidFrom(license)
	validTo := GetLicenseValidTo(license)
	if validFrom.After(validTo) {
		return ErrInvalidTime
	}
	if len(license) != LICENSE_SIZE-1 {
		return ErrorAppend(ErrInvalidSize, "license")
	}
	return nil
}

func ValidPublishingLicense(license Data) error {
	if err := ValidLicense(license, LICENSE_PUBLISHING); err != nil {
		return err
	}
	compositionId := GetLicenseComposition(license)
	if !MatchId(compositionId) {
		return ErrorAppend(ErrInvalidId, "compositionId")
	}
	licenseType := GetLicenseType(license)
	switch licenseType {
	case
		LICENSE_TYPE_MECHANICAL,
		LICENSE_TYPE_SYNCHRONIZATION:
		//..
	default:
		return ErrorAppend(ErrInvalidType, licenseType)
	}
	if len(license) != LICENSE_SIZE {
		return ErrorAppend(ErrInvalidSize, "license")
	}
	return nil
}

func ValidRecordingLicense(license Data) error {
	if err := ValidLicense(license, LICENSE_RECORDING); err != nil {
		return err
	}
	recordingId := GetLicenseRecording(license)
	if !MatchId(recordingId) {
		return ErrorAppend(ErrInvalidId, "recordingId")
	}
	licenseType := GetLicenseType(license)
	switch licenseType {
	case
		LICENSE_TYPE_MASTER:
		//..
	default:
		return ErrorAppend(ErrInvalidType, licenseType)
	}
	if len(license) != LICENSE_SIZE {
		return ErrorAppend(ErrInvalidSize, "license")
	}
	return nil
}
