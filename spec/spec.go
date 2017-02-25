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
	PUBLICATION        = "publication"
	RELEASE            = "release"
	RIGHT              = "right"
	LICENSE_PUBLISHING = "publishing_license"
	LICENSE_RELEASE    = "release_license"
	// LICENSE_BLANKET

	LICENSE_TYPE_MASTER          = "master_license"
	LICENSE_TYPE_MECHANICAL      = "mechanical_license"
	LICENSE_TYPE_SYNCHRONIZATION = "synchronization_license"

	INSTANCE_SIZE    = 2
	AGENT_SIZE       = 4
	COMPOSITION_SIZE = 4
	RECORDING_SIZE   = 5
	PUBLICATION_SIZE = 3
	RELEASE_SIZE     = 3
	RIGHT_SIZE       = 5
	LICENSE_SIZE     = 7

	EMAIL_REGEX           = `(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$)`
	FINGERPRINT_STD_REGEX = `^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$` // base64 std
	FINGERPRINT_URL_REGEX = `^(?:[A-Za-z0-9-_]{4})*(?:[A-Za-z0-9-_]{2}==|[A-Za-z0-9-_]{3})?$`  // base64 url-safe
	HFA_REGEX             = `^[A-Z0-9]{6}$`
	ID_REGEX              = `^[A-Fa-f0-9]{64}$` // hex
	ISRC_REGEX            = `^[a-z]{2}[a-z0-9]{3}[7890][0-9][0-9]{5}$`
	ISWC_REGEX            = `^T-[0-9]{3}\.[0-9]{3}\.[0-9]{3}-[0-9]$`
	PUBKEY_REGEX          = `^[1-9A-HJ-NP-Za-km-z]{43,44}$` // base58
	SIGNATURE_REGEX       = `^[1-9A-HJ-NP-Za-km-z]{87,88}$` // base58
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
		PUBLICATION,
		RELEASE,
		RIGHT,
		LICENSE_PUBLISHING,
		LICENSE_RELEASE:
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

func NewComposition(composerId, hfa, iswc, publisherId, title string) Data {
	return Data{
		"composerId":  composerId,
		"hfa":         hfa,
		"instance":    NewInstance(COMPOSITION),
		"iswc":        iswc,
		"publisherId": publisherId,
		"title":       title,
	}
}

func GetCompositionComposerId(composition Data) string {
	return composition.GetStr("composerId")
}

func GetCompositionHFA(composition Data) string {
	return composition.GetStr("hfa")
}

func GetCompositionISWC(composition Data) string {
	return composition.GetStr("iswc")
}

func GetCompositionPublisherId(composition Data) string {
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
	composerId := GetCompositionComposerId(composition)
	if !MatchId(composerId) {
		return ErrorAppend(ErrInvalidId, "composerId")
	}
	hfa := GetCompositionHFA(composition)
	if !MatchStr(HFA_REGEX, hfa) {
		return Error("Invalid HFA song code")
	}
	iswc := GetCompositionISWC(composition)
	if !MatchStr(HFA_REGEX, iswc) {
		return Error("Invalid ISWC code")
	}
	publisherId := GetCompositionPublisherId(composition)
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

func NewPublication(compositionId string, rightIds []string) Data {
	return Data{
		"compositionId": compositionId,
		"instance":      NewInstance(PUBLICATION),
		"rightIds":      rightIds,
	}
}

func GetPublicationCompositionId(publication Data) string {
	return publication.GetStr("compositionId")
}

func GetPublicationRightIds(publication Data) []string {
	return publication.GetStrSlice("rightIds")
}

func ValidPublication(publication Data) error {
	instance := GetInstance(publication)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(publication, PUBLICATION) {
		return ErrorAppend(ErrInvalidType, GetType(publication))
	}
	compositionId := GetPublicationCompositionId(publication)
	if !MatchId(compositionId) {
		return ErrorAppend(ErrInvalidId, "compositionId")
	}
	rightIds := GetPublicationRightIds(publication)
	seen := make(map[string]struct{})
	for _, rightId := range rightIds {
		if _, ok := seen[rightId]; ok {
			return ErrorAppend(ErrCriteriaNotMet, "multiple references to right")
		}
		if !MatchId(rightId) {
			return ErrorAppend(ErrInvalidId, "rightId")
		}
		seen[rightId] = struct{}{}
	}
	if len(publication) != PUBLICATION_SIZE {
		return ErrorAppend(ErrInvalidSize, PUBLICATION)
	}
	return nil
}

// Recording

func NewRecording(isrc, labelId, performerId, producerId, publicationId, publishingLicenseId string) Data {
	recording := Data{
		"instance":      NewInstance(RECORDING),
		"isrc":          isrc,
		"labelId":       labelId,
		"performerId":   performerId,
		"producerId":    producerId,
		"publicationId": publicationId,
	}
	if publishingLicenseId != "" {
		recording.Set("publishingLicenseId", publishingLicenseId)
	}
	return recording
}

func GetRecordingISRC(recording Data) string {
	return recording.GetStr("isrc")
}

func GetRecordingLabelId(recording Data) string {
	return recording.GetStr("labelId")
}

func GetRecordingPublishingLicenseId(recording Data) string {
	return recording.GetStr("publishingLicenseId")
}

func GetRecordingPerformerId(recording Data) string {
	return recording.GetStr("performerId")
}

func GetRecordingProducerId(recording Data) string {
	return recording.GetStr("producerId")
}

func GetRecordingPublicationId(recording Data) string {
	return recording.GetStr("publicationId")
}

func ValidRecording(recording Data) error {
	instance := GetInstance(recording)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(recording, RECORDING) {
		return ErrorAppend(ErrInvalidType, GetType(recording))
	}
	isrc := GetRecordingISRC(recording)
	if !MatchStr(ISRC_REGEX, isrc) {
		return Error("Invalid ISRC code")
	}
	labelId := GetRecordingLabelId(recording)
	if !MatchId(labelId) {
		return ErrorAppend(ErrInvalidId, "labelId")
	}
	performerId := GetRecordingPerformerId(recording)
	if !MatchId(performerId) {
		return ErrorAppend(ErrInvalidId, "performerId")
	}
	producerId := GetRecordingProducerId(recording)
	if !MatchId(producerId) {
		return ErrorAppend(ErrInvalidId, "performerId")
	}
	publicationId := GetRecordingPublicationId(recording)
	if !MatchId(publicationId) {
		return ErrorAppend(ErrInvalidId, "publicationId")
	}
	publishingLicenseId := GetRecordingPublishingLicenseId(recording)
	if !EmptyStr(publishingLicenseId) {
		if !MatchId(publishingLicenseId) {
			return ErrorAppend(ErrInvalidId, "publishingLicenseId")
		}
		if len(recording) != RECORDING_SIZE+1 {
			return ErrorAppend(ErrInvalidSize, RECORDING)
		}
		return nil
	}
	if len(recording) != RECORDING_SIZE {
		return ErrorAppend(ErrInvalidSize, RECORDING)
	}
	return nil
}

func NewRelease(recordingId string, rightIds []string) Data {
	return Data{
		"instance":    NewInstance(RECORDING),
		"recordingId": recordingId,
		"rightIds":    rightIds,
	}
}

func GetReleaseRecordingId(release Data) string {
	return release.GetStr("recordingId")
}

func GetReleaseRightIds(release Data) []string {
	return release.GetStrSlice("rightIds")
}

func ValidRelease(release Data) error {
	instance := GetInstance(release)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(release, RELEASE) {
		return ErrorAppend(ErrInvalidType, GetType(release))
	}
	recordingId := GetReleaseRecordingId(release)
	if !MatchId(recordingId) {
		return ErrorAppend(ErrInvalidId, "recordingId")
	}
	rightIds := GetReleaseRightIds(release)
	seen := make(map[string]struct{})
	for _, rightId := range rightIds {
		if _, ok := seen[rightId]; ok {
			return ErrorAppend(ErrCriteriaNotMet, "multiple references to right")
		}
		if !MatchId(rightId) {
			return ErrorAppend(ErrInvalidId, "rightId")
		}
		seen[rightId] = struct{}{}
	}
	if len(release) != RELEASE_SIZE {
		return ErrorAppend(ErrInvalidSize, RELEASE)
	}
	return nil
}

// Right

func NewRight(percentageShares, validFrom, validTo string) Data {
	return Data{
		"instance":         NewInstance(RIGHT),
		"percentageShares": percentageShares,
		"validFrom":        validFrom,
		"validTo":          validTo,
	}
}

func NewCompositionRight(compositionId, percentageShares, validFrom, validTo string) Data {
	right := NewRight(percentageShares, validFrom, validTo)
	right.Set("compositionId", compositionId)
	return right
}

func NewRecordingRight(percentageShares, recordingId, validFrom, validTo string) Data {
	right := NewRight(percentageShares, validFrom, validTo)
	right.Set("recordingId", recordingId)
	return right
}

func GetRightCompositionId(right Data) string {
	return right.GetStr("compositionId")
}

func GetRightPercentageShares(right Data) int {
	return right.GetStrInt("percentageShares")
}

func GetRightRecordingId(right Data) string {
	return right.GetStr("recordingId")
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
	percentageShares := GetRightPercentageShares(right)
	if percentageShares <= 0 || percentageShares > 100 {
		return ErrorAppend(ErrCriteriaNotMet, "percentage shares must be greater than 0 and less than 100")
	}
	validFrom := GetRightValidFrom(right)
	validTo := GetRightValidTo(right)
	if !validFrom.Before(validTo) {
		return ErrorAppend(ErrInvalidTime, "range")
	}
	if validTo.Before(Now()) {
		return ErrorAppend(ErrInvalidTime, "expired")
	}
	return nil
}

func ValidCompositionRight(right Data) error {
	if err := ValidRight(right); err != nil {
		return err
	}
	compositionId := GetRightCompositionId(right)
	if !MatchId(compositionId) {
		return ErrorAppend(ErrInvalidId, "compositionId")
	}
	if len(right) != RIGHT_SIZE {
		return ErrorAppend(ErrInvalidSize, RIGHT)
	}
	return nil
}

func ValidRecordingRight(right Data) error {
	if err := ValidRight(right); err != nil {
		return err
	}
	recordingId := GetRightRecordingId(right)
	if !MatchId(recordingId) {
		return ErrorAppend(ErrInvalidId, "recordingId")
	}
	if len(right) != RIGHT_SIZE {
		return ErrorAppend(ErrInvalidSize, RIGHT)
	}
	return nil
}

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

func NewPublishingLicense(licenseeId, licenserId, licenseType, publicationId, validFrom, validTo string) Data {
	license := NewLicense(licenseeId, licenserId, licenseType, LICENSE_PUBLISHING, validFrom, validTo)
	license.Set("publicationId", publicationId)
	return license
}

func NewReleaseLicense(licenseeId, licenserId, licenseType, releaseId, validFrom, validTo string) Data {
	license := NewLicense(licenseeId, licenserId, licenseType, LICENSE_RELEASE, validFrom, validTo)
	license.Set("recordingId", releaseId)
	return license
}
func GetLicenseLicenseeId(license Data) string {
	return license.GetStr("licenseeId")
}

func GetLicenseLicenserId(license Data) string {
	return license.GetStr("licenserId")
}

func GetLicenseReleaseId(license Data) string {
	return license.GetStr("releaseId")
}

func GetLicensePublicationId(license Data) string {
	return license.GetStr("publicationId")
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
	licenseeId := GetLicenseLicenseeId(license)
	if !MatchId(licenseeId) {
		return ErrorAppend(ErrInvalidId, "licenseeId")
	}
	licenserId := GetLicenseLicenserId(license)
	if !MatchId(licenserId) {
		return ErrorAppend(ErrInvalidId, "licenserId")
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
	// if len(license) != LICENSE_SIZE {
	//	return ErrorAppend(ErrInvalidSize, "license")
	// }
	return nil
}

func ValidPublishingLicense(license Data) error {
	if err := ValidLicense(license, LICENSE_PUBLISHING); err != nil {
		return err
	}
	compositionId := GetLicensePublicationId(license)
	if !MatchId(compositionId) {
		return ErrorAppend(ErrInvalidId, "publicationId")
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
		return ErrorAppend(ErrInvalidSize, LICENSE_PUBLISHING)
	}
	return nil
}

func ValidReleaseLicense(license Data) error {
	if err := ValidLicense(license, LICENSE_RELEASE); err != nil {
		return err
	}
	recordingId := GetLicenseReleaseId(license)
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
		return ErrorAppend(ErrInvalidSize, LICENSE_RELEASE)
	}
	return nil
}
