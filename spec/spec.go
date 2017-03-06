package spec

import (
	"net/url"
	"time"

	. "github.com/zbo14/envoke/common"
)

const (
	AGENT       = "agent"
	COMPOSITION = "composition"
	RECORDING   = "recording"
	PUBLICATION = "publication"
	RELEASE     = "release"
	RIGHT       = "right"
	ASSIGNMENT  = "assignment"
	TRANSFER    = "transfer"

	LICENSE_MASTER     = "master_license"
	LICENSE_MECHANICAL = "mechanical_license"
	// LICENSE_SYNCHRONIZATION
	// LICENSE_BLANKET

	INSTANCE_SIZE    = 2
	AGENT_SIZE       = 4
	COMPOSITION_SIZE = 6
	RECORDING_SIZE   = 6
	PUBLICATION_SIZE = 3
	RELEASE_SIZE     = 3
	RIGHT_SIZE       = 5
	LICENSE_SIZE     = 8
	ASSIGNMENT_SIZE  = 5
	TRANSFER_SIZE    = 5

	EMAIL_REGEX           = `(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$)`
	FINGERPRINT_STD_REGEX = `^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$` // base64 std
	FINGERPRINT_URL_REGEX = `^(?:[A-Za-z0-9-_]{4})*(?:[A-Za-z0-9-_]{2}==|[A-Za-z0-9-_]{3})?$`  // base64 url-safe
	HFA_REGEX             = `^[A-Z0-9]{6}$`
	ID_REGEX              = `^[A-Fa-f0-9]{64}$` // hex
	ISRC_REGEX            = `^[A-Z]{2}-[A-Z0-9]{3}-[7890][0-9]-[0-9]{5}$`
	ISWC_REGEX            = `^T-[0-9]{3}\.[0-9]{3}\.[0-9]{3}-[0-9]$`
	PUBKEY_REGEX          = `^[1-9A-HJ-NP-Za-km-z]{43,44}$` // base58
	SIGNATURE_REGEX       = `^[1-9A-HJ-NP-Za-km-z]{87,88}$` // base58
	TERRITORY_REGEX       = `^[A-Z]{2}$`
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

func MatchTerritory(territory string) bool {
	return MatchStr(TERRITORY_REGEX, territory)
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
		instance = thing.GetMapData("instance")
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
		ASSIGNMENT,
		LICENSE_MECHANICAL,
		LICENSE_MASTER,
		TRANSFER:
		//..
	default:
		return ErrorAppend(ErrInvalidType, _type)
	}
	if len(instance) != INSTANCE_SIZE {
		return ErrInvalidSize
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
		return ErrorAppend(ErrEmptyStr, name)
	}
	socialMedia := GetAgentSocialMediaStr(agent)
	if !MatchUrlRelaxed(socialMedia) {
		return ErrorAppend(ErrInvalidUrl, socialMedia)
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
		return ErrorAppend(ErrInvalidId, composerId)
	}
	hfa := GetCompositionHFA(composition)
	if !MatchStr(HFA_REGEX, hfa) {
		return Error("Invalid HFA song code")
	}
	iswc := GetCompositionISWC(composition)
	if !MatchStr(ISWC_REGEX, iswc) {
		return Error("Invalid ISWC code")
	}
	publisherId := GetCompositionPublisherId(composition)
	if !MatchId(publisherId) {
		return ErrorAppend(ErrInvalidId, publisherId)
	}
	title := GetCompositionTitle(composition)
	if EmptyStr(title) {
		return ErrEmptyStr
	}
	if len(composition) != COMPOSITION_SIZE {
		return ErrorAppend(ErrInvalidSize, COMPOSITION)
	}
	return nil
}

func NewPublication(assignmentIds []string, compositionId string) Data {
	return Data{
		"assignmentIds": assignmentIds,
		"compositionId": compositionId,
		"instance":      NewInstance(PUBLICATION),
	}
}

func GetPublicationAssignmentIds(publication Data) []string {
	return publication.GetStrSlice("assignmentIds")
}

func GetPublicationCompositionId(publication Data) string {
	return publication.GetStr("compositionId")
}

func ValidPublication(publication Data) error {
	instance := GetInstance(publication)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(publication, PUBLICATION) {
		return ErrorAppend(ErrInvalidType, GetType(publication))
	}
	assignmentIds := GetPublicationAssignmentIds(publication)
	seen := make(map[string]struct{})
	for _, assignmentId := range assignmentIds {
		if _, ok := seen[assignmentId]; ok {
			return ErrorAppend(ErrCriteriaNotMet, "multiple references to assignment")
		}
		if !MatchId(assignmentId) {
			return ErrorAppend(ErrInvalidId, assignmentId)
		}
		seen[assignmentId] = struct{}{}
	}
	compositionId := GetPublicationCompositionId(publication)
	if !MatchId(compositionId) {
		return ErrorAppend(ErrInvalidId, compositionId)
	}
	if len(publication) != PUBLICATION_SIZE {
		return ErrorAppend(ErrInvalidSize, PUBLICATION)
	}
	return nil
}

// Recording

func NewRecording(assignmentId, isrc, labelId, performerId, producerId, publicationId string) Data {
	recording := Data{
		"instance":      NewInstance(RECORDING),
		"isrc":          isrc,
		"labelId":       labelId,
		"performerId":   performerId,
		"producerId":    producerId,
		"publicationId": publicationId,
	}
	if !EmptyStr(assignmentId) {
		recording.Set("assignmentId", assignmentId)
	}
	return recording
}

func GetRecordingAssignmentId(recording Data) string {
	return recording.GetStr("assignmentId")
}

func GetRecordingISRC(recording Data) string {
	return recording.GetStr("isrc")
}

func GetRecordingLabelId(recording Data) string {
	return recording.GetStr("labelId")
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
		return ErrorAppend(ErrInvalidId, labelId)
	}
	performerId := GetRecordingPerformerId(recording)
	if !MatchId(performerId) {
		return ErrorAppend(ErrInvalidId, performerId)
	}
	producerId := GetRecordingProducerId(recording)
	if !MatchId(producerId) {
		return ErrorAppend(ErrInvalidId, producerId)
	}
	publicationId := GetRecordingPublicationId(recording)
	if !MatchId(publicationId) {
		return ErrorAppend(ErrInvalidId, publicationId)
	}
	assignmentId := GetRecordingAssignmentId(recording)
	if !EmptyStr(assignmentId) {
		if !MatchId(assignmentId) {
			return ErrorAppend(ErrInvalidId, assignmentId)
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

func NewRelease(assignmentIds []string, licenseId, recordingId string) Data {
	release := Data{
		"assignmentIds": assignmentIds,
		"instance":      NewInstance(RELEASE),
		"recordingId":   recordingId,
	}
	if licenseId != "" {
		release.Set("licenseId", licenseId)
	}
	return release
}

func GetReleaseAssignmentIds(release Data) []string {
	return release.GetStrSlice("assignmentIds")
}

func GetReleaseLicenseId(release Data) string {
	return release.GetStr("licenseId")
}

func GetReleaseRecordingId(release Data) string {
	return release.GetStr("recordingId")
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
		return ErrorAppend(ErrInvalidId, recordingId)
	}
	licenseId := GetReleaseLicenseId(release)
	if !EmptyStr(licenseId) {
		if !MatchId(licenseId) {
			return ErrorAppend(ErrInvalidId, licenseId)
		}
		if len(release) != RELEASE_SIZE+1 {
			return ErrorAppend(ErrInvalidSize, RELEASE)
		}
		return nil
	}
	assignmentIds := GetReleaseAssignmentIds(release)
	seen := make(map[string]struct{})
	for _, assignmentId := range assignmentIds {
		if _, ok := seen[assignmentId]; ok {
			return ErrorAppend(ErrCriteriaNotMet, "multiple references to assignment")
		}
		if !MatchId(assignmentId) {
			return ErrorAppend(ErrInvalidId, assignmentId)
		}
		seen[assignmentId] = struct{}{}
	}
	if len(release) != RELEASE_SIZE {
		return ErrorAppend(ErrInvalidSize, RELEASE)
	}
	return nil
}

// Assignment
func NewAssignment(holderId, rightId, signerId string) Data {
	return Data{
		"holderId": holderId,
		"instance": NewInstance(ASSIGNMENT),
		"rightId":  rightId,
		"signerId": signerId,
	}
}

func GetAssignmentHolderId(assignment Data) string {
	return assignment.GetStr("holderId")
}

func GetAssignmentRight(assignment Data) Data {
	return assignment.GetData("right")
}

func GetAssignmentRightId(assignment Data) string {
	return assignment.GetStr("rightId")
}

func GetAssignmentSignerId(assignment Data) string {
	return assignment.GetStr("signerId")
}

func ValidAssignment(assignment Data) error {
	instance := GetInstance(assignment)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(assignment, ASSIGNMENT) {
		return ErrorAppend(ErrInvalidType, GetType(assignment))
	}
	holderId := GetAssignmentHolderId(assignment)
	if !MatchId(holderId) {
		return ErrorAppend(ErrInvalidId, holderId)
	}
	rightId := GetAssignmentRightId(assignment)
	if !MatchId(rightId) {
		return ErrorAppend(ErrInvalidId, rightId)
	}
	signerId := GetAssignmentSignerId(assignment)
	if !MatchId(signerId) {
		return ErrorAppend(ErrInvalidId, signerId)
	}
	if len(assignment) != ASSIGNMENT_SIZE {
		return ErrorAppend(ErrInvalidSize, ASSIGNMENT)
	}
	return nil
}

// Right

func NewRight(territory []string, validFrom, validTo string) Data {
	return Data{
		"instance":  NewInstance(RIGHT),
		"territory": territory,
		"validFrom": validFrom,
		"validTo":   validTo,
	}
}

func NewCompositionRight(compositionId string, territory []string, validFrom, validTo string) Data {
	right := NewRight(territory, validFrom, validTo)
	right.Set("compositionId", compositionId)
	return right
}

func NewRecordingRight(recordingId string, territory []string, validFrom, validTo string) Data {
	right := NewRight(territory, validFrom, validTo)
	right.Set("recordingId", recordingId)
	return right
}

func GetRightCompositionId(right Data) string {
	return right.GetStr("compositionId")
}

func GetRightPercentageShares(right Data) int {
	return right.GetInt("percentageShares")
}

func GetRightRecordingId(right Data) string {
	return right.GetStr("recordingId")
}

func GetTerritory(right Data) []string {
	return right.GetStrSlice("territory")
}

func GetValidFrom(right Data) time.Time {
	return MustParseDateStr(right.GetStr("validFrom"))
}

func GetValidTo(right Data) time.Time {
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
	seen := make(map[string]struct{})
	for _, territory := range GetTerritory(right) {
		if !MatchTerritory(territory) {
			return ErrInvalidTerritory
		}
		if _, ok := seen[territory]; ok {
			return ErrorAppend(ErrCriteriaNotMet, "territory listed multiple times")
		}
		seen[territory] = struct{}{}
	}
	validFrom := GetValidFrom(right)
	validTo := GetValidTo(right)
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

func NewLicense(assignmentId, licenseeId, licenserId, publicationId, releaseId string, territory []string, _type, validFrom, validTo string) Data {
	license := Data{
		"assignmentId": assignmentId,
		"instance":     NewInstance(_type),
		"licenseeId":   licenseeId,
		"licenserId":   licenserId,
		"territory":    territory,
		"validFrom":    validFrom,
		"validTo":      validTo,
	}
	switch _type {
	case LICENSE_MECHANICAL:
		license.Set("publicationId", publicationId)
	case LICENSE_MASTER:
		license.Set("releaseId", releaseId)
	default:
		panic(ErrorAppend(ErrInvalidType, _type))
	}
	return license
}

func GetLicenseAssignmentId(license Data) string {
	return license.GetStr("assignmentId")
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

func ValidLicense(license Data) error {
	instance := GetInstance(license)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	_type := GetType(license)
	switch _type {
	case LICENSE_MECHANICAL:
		publicationId := GetLicensePublicationId(license)
		if !MatchId(publicationId) {
			return ErrorAppend(ErrInvalidId, publicationId)
		}
	case LICENSE_MASTER:
		releaseId := GetLicenseReleaseId(license)
		if !MatchId(releaseId) {
			return ErrorAppend(ErrInvalidId, releaseId)
		}
	default:
		return ErrorAppend(ErrInvalidType, _type)
	}
	assignmentId := GetLicenseAssignmentId(license)
	if !MatchId(assignmentId) {
		return ErrorAppend(ErrInvalidId, assignmentId)
	}
	licenseeId := GetLicenseLicenseeId(license)
	if !MatchId(licenseeId) {
		return ErrorAppend(ErrInvalidId, licenseeId)
	}
	licenserId := GetLicenseLicenserId(license)
	if !MatchId(licenserId) {
		return ErrorAppend(ErrInvalidId, licenserId)
	}
	seen := make(map[string]struct{})
	for _, territory := range GetTerritory(license) {
		if !MatchTerritory(territory) {
			return ErrInvalidTerritory
		}
		if _, ok := seen[territory]; ok {
			return ErrorAppend(ErrCriteriaNotMet, "territory listed multiple times")
		}
		seen[territory] = struct{}{}
	}
	validFrom := GetValidFrom(license)
	validTo := GetValidTo(license)
	if validFrom.After(validTo) {
		return ErrInvalidTime
	}
	if len(license) != LICENSE_SIZE {
		return ErrorAppend(ErrInvalidSize, "license")
	}
	return nil
}

// Transfer

func NewTransfer(recipientId, senderId, txId string) Data {
	return Data{
		"instance":    NewInstance(TRANSFER),
		"recipientId": recipientId,
		"senderId":    senderId,
		"txId":        txId,
	}
}

func NewCompositionRightTransfer(publicationId, recipientId, senderId, txId string) Data {
	transfer := NewTransfer(recipientId, senderId, txId)
	transfer.Set("publicationId", publicationId)
	return transfer
}

func NewRecordingRightTransfer(recipientId, releaseId, senderId, txId string) Data {
	transfer := NewTransfer(recipientId, senderId, txId)
	transfer.Set("releaseId", releaseId)
	return transfer
}

func GetTransferAssignmentId(transfer Data) string {
	return transfer.GetStr("assignmentId")
}

func GetTransferRecipientShares(transfer Data) int {
	return transfer.GetInt("recipientShares")
}

func GetTransferSenderShares(transfer Data) int {
	return transfer.GetInt("senderShares")
}

func GetTransferPublicationId(transfer Data) string {
	return transfer.GetStr("publicationId")
}

func GetTransferRecipientId(transfer Data) string {
	return transfer.GetStr("recipientId")
}

func GetTransferReleaseId(transfer Data) string {
	return transfer.GetStr("releaseId")
}

func GetTransferSenderId(transfer Data) string {
	return transfer.GetStr("senderId")
}

func GetTransferTxId(transfer Data) string {
	return transfer.GetStr("txId")
}

func ValidCompositionRightTransfer(transfer Data) error {
	if err := ValidTransfer(transfer); err != nil {
		return err
	}
	publicationId := GetTransferPublicationId(transfer)
	if !MatchId(publicationId) {
		return ErrorAppend(ErrInvalidId, publicationId)
	}
	if len(transfer) != TRANSFER_SIZE {
		return ErrorAppend(ErrInvalidSize, TRANSFER)
	}
	return nil
}

func ValidRecordingRightTransfer(transfer Data) error {
	if err := ValidTransfer(transfer); err != nil {
		return err
	}
	releaseId := GetTransferReleaseId(transfer)
	if !MatchId(releaseId) {
		return ErrorAppend(ErrInvalidId, releaseId)
	}
	if len(transfer) != TRANSFER_SIZE {
		return ErrorAppend(ErrInvalidSize, TRANSFER)
	}
	return nil
}

func ValidTransfer(transfer Data) error {
	instance := GetInstance(transfer)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(transfer, TRANSFER) {
		return ErrorAppend(ErrInvalidType, GetType(transfer))
	}
	recipientId := GetTransferRecipientId(transfer)
	if !MatchId(recipientId) {
		return ErrorAppend(ErrInvalidId, recipientId)
	}
	senderId := GetTransferSenderId(transfer)
	if !MatchId(senderId) {
		return ErrorAppend(ErrInvalidId, senderId)
	}
	if recipientId == senderId {
		return ErrorAppend(ErrCriteriaNotMet, "recipientId and senderId must be different")
	}
	txId := GetTransferTxId(transfer)
	if !MatchId(txId) {
		return ErrorAppend(ErrInvalidId, txId)
	}
	return nil
}
