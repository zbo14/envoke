package spec

import (
	"net/url"
	"time"

	. "github.com/zbo14/envoke/common"
)

const (
	AGENT              = "agent"
	INFO_COMPOSITION   = "composition_info"
	INFO_RECORDING     = "recording_info"
	COMPOSITION        = "composition"
	RECORDING          = "recording"
	RIGHT              = "right"
	LICENSE_PUBLISHING = "publishing_license"
	LICENSE_RECORDING  = "recording_license"
	// LICENSE_BLANKET

	LICENSE_TYPE_MASTER          = "master_license"
	LICENSE_TYPE_MECHANICAL      = "mechanical_license"
	LICENSE_TYPE_SYNCHRONIZATION = "synchronization_license"

	INSTANCE_SIZE         = 2
	AGENT_SIZE            = 4
	INFO_COMPOSITION_SIZE = 4
	INFO_RECORDING_SIZE   = 5
	COMPOSITION_SIZE      = 3
	RECORDING_SIZE        = 3
	RIGHT_SIZE            = 5
	LICENSE_SIZE          = 7

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

func NewCompositionInfo(composerId, publisherId, title string) Data {
	return Data{
		"composerId":  composerId,
		"instance":    NewInstance(INFO_COMPOSITION),
		"publisherId": publisherId,
		"title":       title,
	}
}

func NewComposition(infoId string, rightIds []string) Data {
	return Data{
		"infoId":   infoId,
		"instance": NewInstance(COMPOSITION),
		"rightIds": rightIds,
	}
}

func GetInfoComposer(info Data) string {
	return info.GetStr("composerId")
}

func GetInfoPublisher(info Data) string {
	return info.GetStr("publisherId")
}

func GetInfoTitle(info Data) string {
	return info.GetStr("title")
}

func GetCompositionRights(composition Data) []string {
	return composition.GetStrSlice("rightIds")
}

func GetCompositionInfo(composition Data) string {
	return composition.GetStr("infoId")
}

func ValidCompositionInfo(info Data) error {
	instance := GetInstance(info)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(info, INFO_COMPOSITION) {
		return ErrorAppend(ErrInvalidType, GetType(info))
	}
	composerId := GetInfoComposer(info)
	if !MatchId(composerId) {
		return ErrorAppend(ErrInvalidId, "composerId")
	}
	publisherId := GetInfoPublisher(info)
	if !MatchId(publisherId) {
		return ErrorAppend(ErrInvalidId, "publisherId")
	}
	title := GetInfoTitle(info)
	if EmptyStr(title) {
		return ErrorAppend(ErrEmptyStr, "title")
	}
	if len(info) != INFO_COMPOSITION_SIZE {
		return ErrorAppend(ErrInvalidSize, INFO_COMPOSITION)
	}
	return nil
}

func ValidComposition(composition Data) error {
	instance := GetInstance(composition)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(composition, COMPOSITION) {
		return ErrorAppend(ErrInvalidType, GetType(composition))
	}
	infoId := GetCompositionInfo(composition)
	if !MatchId(infoId) {
		return ErrorAppend(ErrInvalidId, "infoId")
	}
	rightIds := GetCompositionRights(composition)
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
	if len(composition) != COMPOSITION_SIZE {
		return ErrorAppend(ErrInvalidSize, COMPOSITION)
	}
	return nil
}

// Recording

func NewRecordingInfo(compositionId, labelId, performerId, producerId, publishingLicenseId string) Data {
	info := Data{
		"compositionId": compositionId,
		"instance":      NewInstance(INFO_RECORDING),
		"labelId":       labelId,
		"performerId":   performerId,
		"producerId":    producerId,
	}
	if publishingLicenseId != "" {
		info.Set("publishingLicenseId", publishingLicenseId)
	}
	return info
}

func NewRecording(infoId string, rightIds []string) Data {
	return Data{
		"infoId":   infoId,
		"instance": NewInstance(RECORDING),
		"rightIds": rightIds,
	}
}

func GetInfoComposition(info Data) string {
	return info.GetStr("compositionId")
}

func GetInfoLabel(info Data) string {
	return info.GetStr("labelId")
}

func GetInfoPublishingLicense(info Data) string {
	return info.GetStr("publishingLicenseId")
}

func GetInfoPerformer(info Data) string {
	return info.GetStr("performerId")
}

func GetInfoProducer(info Data) string {
	return info.GetStr("producerId")
}

func GetRecordingRights(recording Data) []string {
	return recording.GetStrSlice("rightIds")
}

func GetRecordingInfo(recording Data) string {
	return recording.GetStr("infoId")
}

func ValidRecordingInfo(info Data) error {
	instance := GetInstance(info)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(info, INFO_RECORDING) {
		return ErrorAppend(ErrInvalidType, GetType(info))
	}
	compositionId := GetInfoComposition(info)
	if !MatchId(compositionId) {
		return ErrorAppend(ErrInvalidId, "compositionId")
	}
	labelId := GetInfoLabel(info)
	if !MatchId(labelId) {
		return ErrorAppend(ErrInvalidId, "labelId")
	}
	performerId := GetInfoPerformer(info)
	if !MatchId(performerId) {
		return ErrorAppend(ErrInvalidId, "performerId")
	}
	producerId := GetInfoProducer(info)
	if !MatchId(producerId) {
		return ErrorAppend(ErrInvalidId, "performerId")
	}
	publishingLicenseId := GetInfoPublishingLicense(info)
	if !EmptyStr(publishingLicenseId) {
		if !MatchId(publishingLicenseId) {
			return ErrorAppend(ErrInvalidId, "publishingLicenseId")
		}
		if len(info) != INFO_RECORDING_SIZE+1 {
			return ErrorAppend(ErrInvalidSize, INFO_RECORDING)
		}
		return nil
	}
	if len(info) != INFO_RECORDING_SIZE {
		return ErrorAppend(ErrInvalidSize, INFO_RECORDING)
	}
	return nil
}

func ValidRecording(recording Data) error {
	instance := GetInstance(recording)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(recording, RECORDING) {
		return ErrorAppend(ErrInvalidType, GetType(recording))
	}
	infoId := GetRecordingInfo(recording)
	if !MatchId(infoId) {
		return ErrorAppend(ErrInvalidId, "infoId")
	}
	rightIds := GetRecordingRights(recording)
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
	if len(recording) != RECORDING_SIZE {
		return ErrorAppend(ErrInvalidSize, RECORDING)
	}
	return nil
}

// Right

func NewRight(infoId, percentageShares, validFrom, validTo string) Data {
	return Data{
		"infoId":           infoId,
		"instance":         NewInstance(RIGHT),
		"percentageShares": percentageShares,
		"validFrom":        validFrom,
		"validTo":          validTo,
	}
}

func GetRightInfo(right Data) string {
	return right.GetStr("infoId")
}

func GetRightPercentageShares(right Data) int {
	return right.GetStrInt("percentageShares")
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
	infoId := GetRightInfo(right)
	if !MatchId(infoId) {
		return ErrorAppend(ErrInvalidId, "infoId")
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
	if len(right) != RIGHT_SIZE {
		return ErrorAppend(ErrInvalidSize, RIGHT)
	}
	return nil
}

// License

func NewLicense(licenseeId, licenserId, licenseType, rightId, _type, validFrom, validTo string) Data {
	return Data{
		"instance":    NewInstance(_type),
		"licenseeId":  licenseeId,
		"licenserId":  licenserId,
		"licenseType": licenseType,
		"validFrom":   validFrom,
		"validTo":     validTo,
	}
}

func NewPublishingLicense(compositionId, licenseeId, licenserId, licenseType, rightId, validFrom, validTo string) Data {
	license := NewLicense(licenseeId, licenserId, licenseType, rightId, LICENSE_PUBLISHING, validFrom, validTo)
	license.Set("compositionId", compositionId)
	return license
}

func NewRecordingLicense(licenseeId, licenserId, licenseType, recordingId, rightId, validFrom, validTo string) Data {
	license := NewLicense(licenseeId, licenserId, licenseType, rightId, LICENSE_RECORDING, validFrom, validTo)
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
	// if len(license) != LICENSE_SIZE {
	//	return ErrorAppend(ErrInvalidSize, "license")
	// }
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
