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
	LICENSE     = "license"
	TRANSFER    = "transfer"

	AGENT_SIZE           = 4
	INSTANCE_SIZE        = 2
	LICENSE_SIZE         = 9
	MIN_COMPOSITION_SIZE = 3
	MIN_RECORDING_SIZE   = 4
	MIN_RELEASE_SIZE     = 4
	PUBLICATION_SIZE     = 4
	RIGHT_SIZE           = 8
	TRANSFER_SIZE        = 6

	EMAIL_REGEX           = `(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$)`
	FINGERPRINT_STD_REGEX = `^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$` // base64 std
	FINGERPRINT_URL_REGEX = `^(?:[A-Za-z0-9-_]{4})*(?:[A-Za-z0-9-_]{2}==|[A-Za-z0-9-_]{3})?$`  // base64 url-safe
	HFA_REGEX             = `^[A-Z0-9]{6}$`
	ID_REGEX              = `^[A-Fa-f0-9]{64}$` // hex
	IPI_REGEX             = `^[0-9]{9}$`
	ISRC_REGEX            = `^[A-Z]{2}-[A-Z0-9]{3}-[7890][0-9]-[0-9]{5}$`
	ISWC_REGEX            = `^T-[0-9]{3}\.[0-9]{3}\.[0-9]{3}-[0-9]$`
	PRO_REGEX             = `^ASCAP|BMI|SESAC$`
	PUBKEY_REGEX          = `^[1-9A-HJ-NP-Za-km-z]{43,44}$` // base58
	SIGNATURE_REGEX       = `^[1-9A-HJ-NP-Za-km-z]{87,88}$` // base58
	TERRITORY_REGEX       = `^[A-Z]{2}$`
)

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
		LICENSE,
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

func GetEmail(agent Data) string {
	return agent.GetStr("email")
}

func GetName(agent Data) string {
	return agent.GetStr("name")
}

func GetSocialMediaStr(agent Data) string {
	return agent.GetStr("socialMedia")
}

func GetSocialMedia(agent Data) *url.URL {
	return MustParseUrl(GetSocialMediaStr(agent))
}

func ValidAgent(agent Data) error {
	instance := GetInstance(agent)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(agent, AGENT) {
		return ErrorAppend(ErrInvalidType, GetType(agent))
	}
	email := GetEmail(agent)
	if !MatchStr(EMAIL_REGEX, email) {
		return ErrorAppend(ErrInvalidEmail, email)
	}
	name := GetName(agent)
	if EmptyStr(name) {
		return ErrorAppend(ErrEmptyStr, name)
	}
	socialMedia := GetSocialMediaStr(agent)
	if !MatchUrlRelaxed(socialMedia) {
		return ErrorAppend(ErrInvalidUrl, socialMedia)
	}
	if len(agent) != AGENT_SIZE {
		return ErrorAppend(ErrInvalidSize, AGENT)
	}
	return nil
}

// Composition

func NewComposition(composerId, hfa, ipi, iswc, pro, title string) Data {
	composition := Data{
		"composerId": composerId,
		"instance":   NewInstance(COMPOSITION),
		"title":      title,
	}
	if !EmptyStr(hfa) {
		composition.Set("hfa", hfa)
	}
	if !EmptyStr(ipi) {
		composition.Set("ipi", ipi)
	}
	if !EmptyStr(iswc) {
		composition.Set("iswc", iswc)
	}
	if !EmptyStr(pro) {
		composition.Set("pro", pro)
	}
	return composition
}

func GetComposerId(data Data) string {
	return data.GetStr("composerId")
}

func GetHFA(data Data) string {
	return data.GetStr("hfa")
}

func GetIPI(data Data) string {
	return data.GetStr("ipi")
}

func GetISWC(data Data) string {
	return data.GetStr("iswc")
}

func GetPRO(data Data) string {
	return data.GetStr("pro")
}

func GetTitle(data Data) string {
	return data.GetStr("title")
}

func ValidComposition(composition Data) error {
	instance := GetInstance(composition)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(composition, COMPOSITION) {
		return ErrorAppend(ErrInvalidType, GetType(composition))
	}
	size := MIN_COMPOSITION_SIZE
	composerId := GetComposerId(composition)
	if !MatchId(composerId) {
		return ErrorAppend(ErrInvalidId, composerId)
	}
	hfa := GetHFA(composition)
	if !EmptyStr(hfa) {
		if !MatchStr(HFA_REGEX, hfa) {
			return Error("Invalid HFA song code")
		}
		size++
	}
	ipi := GetIPI(composition)
	if !EmptyStr(ipi) {
		if !MatchStr(IPI_REGEX, ipi) {
			return Error("Invalid IPI number")
		}
		size++
	}
	iswc := GetISWC(composition)
	if !EmptyStr(iswc) {
		if !MatchStr(ISWC_REGEX, iswc) {
			return Error("Invalid ISWC code")
		}
		size++
	}
	pro := GetPRO(composition)
	if !EmptyStr(pro) {
		if !MatchStr(PRO_REGEX, pro) {
			return Error("Invalid PRO name")
		}
		size++
	}
	title := GetTitle(composition)
	if EmptyStr(title) {
		return ErrorAppend(ErrEmptyStr, "title")
	}
	if len(composition) != size {
		return ErrorAppend(ErrInvalidSize, COMPOSITION)
	}
	return nil
}

func NewPublication(compositionId string, compositionRightIds []string, publisherId string) Data {
	return Data{
		"compositionId":       compositionId,
		"compositionRightIds": compositionRightIds,
		"instance":            NewInstance(PUBLICATION),
		"publisherId":         publisherId,
	}
}

func GetCompositionId(data Data) string {
	return data.GetStr("compositionId")
}

func GetCompositionRightIds(data Data) []string {
	return data.GetStrSlice("compositionRightIds")
}

func GetPublisherId(data Data) string {
	return data.GetStr("publisherId")
}

func ValidPublication(publication Data) error {
	instance := GetInstance(publication)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(publication, PUBLICATION) {
		return ErrorAppend(ErrInvalidType, GetType(publication))
	}
	compositionRightIds := GetCompositionRightIds(publication)
	seen := make(map[string]struct{})
	for _, rightId := range compositionRightIds {
		if _, ok := seen[rightId]; ok {
			return ErrorAppend(ErrCriteriaNotMet, "multiple references to composition right")
		}
		if !MatchId(rightId) {
			return ErrorAppend(ErrInvalidId, rightId)
		}
		seen[rightId] = struct{}{}
	}
	compositionId := GetCompositionId(publication)
	if !MatchId(compositionId) {
		return ErrorAppend(ErrInvalidId, compositionId)
	}
	publisherId := GetPublisherId(publication)
	if !MatchId(publisherId) {
		return ErrorAppend(ErrInvalidId, publisherId)
	}
	if n := len(publication); n != PUBLICATION_SIZE {
		return ErrorAppend(ErrInvalidSize, Itoa(n))
	}
	return nil
}

// Recording

func NewRecording(compositionRightId, isrc, performerId, producerId, publicationId string) Data {
	recording := Data{
		"instance":      NewInstance(RECORDING),
		"performerId":   performerId,
		"producerId":    producerId,
		"publicationId": publicationId,
	}
	if !EmptyStr(compositionRightId) {
		recording.Set("compositionRightId", compositionRightId)
	}
	if !EmptyStr(isrc) {
		recording.Set("isrc", isrc)
	}
	return recording
}

func GetCompositionRightId(data Data) string {
	return data.GetStr("compositionRightId")
}

func GetISRC(data Data) string {
	return data.GetStr("isrc")
}

func GetPerformerId(data Data) string {
	return data.GetStr("performerId")
}

func GetProducerId(data Data) string {
	return data.GetStr("producerId")
}

func GetPublicationId(data Data) string {
	return data.GetStr("publicationId")
}

func ValidRecording(recording Data) error {
	instance := GetInstance(recording)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(recording, RECORDING) {
		return ErrorAppend(ErrInvalidType, GetType(recording))
	}
	size := MIN_RECORDING_SIZE
	compositionRightId := GetCompositionRightId(recording)
	if !EmptyStr(compositionRightId) {
		if !MatchId(compositionRightId) {
			return ErrorAppend(ErrInvalidId, compositionRightId)
		}
		size++
	}
	isrc := GetISRC(recording)
	if !EmptyStr(isrc) {
		if !MatchStr(ISRC_REGEX, isrc) {
			return Error("Invalid ISRC code")
		}
		size++
	}
	performerId := GetPerformerId(recording)
	if !MatchId(performerId) {
		return ErrorAppend(ErrInvalidId, performerId)
	}
	producerId := GetProducerId(recording)
	if !MatchId(producerId) {
		return ErrorAppend(ErrInvalidId, producerId)
	}
	publicationId := GetPublicationId(recording)
	if !MatchId(publicationId) {
		return ErrorAppend(ErrInvalidId, publicationId)
	}
	if n := len(recording); n != size {
		return ErrorAppend(ErrInvalidSize, Itoa(n))
	}
	return nil
}

func NewRelease(mechanicalLicenseId, recordingId string, recordingRightIds []string, recordLabelId string) Data {
	release := Data{
		"instance":          NewInstance(RELEASE),
		"recordingId":       recordingId,
		"recordingRightIds": recordingRightIds,
		"recordLabelId":     recordLabelId,
	}
	if !EmptyStr(mechanicalLicenseId) {
		release.Set("mechanicalLicenseId", mechanicalLicenseId)
	}
	return release
}

func GetMechanicalLicenseId(data Data) string {
	return data.GetStr("mechanicalLicenseId")
}

func GetRecordingId(data Data) string {
	return data.GetStr("recordingId")
}

func GetRecordingRightIds(data Data) []string {
	return data.GetStrSlice("recordingRightIds")
}

func GetRecordLabelId(data Data) string {
	return data.GetStr("recordLabelId")
}

func ValidRelease(release Data) error {
	instance := GetInstance(release)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(release, RELEASE) {
		return ErrorAppend(ErrInvalidType, GetType(release))
	}
	size := MIN_RELEASE_SIZE
	recordingId := GetRecordingId(release)
	if !MatchId(recordingId) {
		return ErrorAppend(ErrInvalidId, recordingId)
	}
	mechanicalLicenseId := GetMechanicalLicenseId(release)
	if !EmptyStr(mechanicalLicenseId) {
		if !MatchId(mechanicalLicenseId) {
			return ErrorAppend(ErrInvalidId, mechanicalLicenseId)
		}
		size++
	}
	recordingRightIds := GetRecordingRightIds(release)
	seen := make(map[string]struct{})
	for _, rightId := range recordingRightIds {
		if _, ok := seen[rightId]; ok {
			return ErrorAppend(ErrCriteriaNotMet, "multiple references to assignment")
		}
		if !MatchId(rightId) {
			return ErrorAppend(ErrInvalidId, rightId)
		}
		seen[rightId] = struct{}{}
	}
	recordLabelId := GetRecordLabelId(release)
	if !MatchId(recordLabelId) {
		return ErrorAppend(ErrInvalidId, recordLabelId)
	}
	if n := len(release); n != size {
		return ErrorAppend(ErrInvalidSize, Itoa(n))
	}
	return nil
}

// Right

func NewRight(recipientId, senderId string, territory, usage []string, validFrom, validThrough string) Data {
	return Data{
		"instance":     NewInstance(RIGHT),
		"recipientId":  recipientId,
		"senderId":     senderId,
		"territory":    territory,
		"usage":        usage,
		"validFrom":    validFrom,
		"validThrough": validThrough,
	}
}

func NewCompositionRight(compositionId string, recipientId, senderId string, territory, usage []string, validFrom, validThrough string) Data {
	right := NewRight(recipientId, senderId, territory, usage, validFrom, validThrough)
	right.Set("compositionId", compositionId)
	return right
}

func NewRecordingRight(recordingId string, recipientId, senderId string, territory, usage []string, validFrom, validThrough string) Data {
	right := NewRight(recipientId, senderId, territory, usage, validFrom, validThrough)
	right.Set("recordingId", recordingId)
	return right
}

func GetRecipientId(data Data) string {
	return data.GetStr("recipientId")
}

func GetSenderId(data Data) string {
	return data.GetStr("senderId")
}

func GetTerritory(right Data) []string {
	return right.GetStrSlice("territory")
}

func GetUsage(right Data) []string {
	return right.GetStrSlice("usage")
}

func GetValidFrom(right Data) time.Time {
	return MustParseDateStr(right.GetStr("validFrom"))
}

func GetValidThrough(right Data) time.Time {
	return MustParseDateStr(right.GetStr("validThrough"))
}

func ValidRight(right Data) error {
	instance := GetInstance(right)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(right, RIGHT) {
		return ErrorAppend(ErrInvalidType, GetType(right))
	}
	recipientId := GetRecipientId(right)
	if !MatchId(recipientId) {
		return ErrorAppend(ErrInvalidId, recipientId)
	}
	senderId := GetSenderId(right)
	if !MatchId(senderId) {
		return ErrorAppend(ErrInvalidId, senderId)
	}
	territory := GetTerritory(right)
	if len(territory) == 0 {
		return Error("no territory listed")
	}
	seen := make(map[string]struct{})
	for i := range territory {
		if !MatchStr(TERRITORY_REGEX, territory[i]) {
			return ErrInvalidTerritory
		}
		if _, ok := seen[territory[i]]; ok {
			return ErrorAppend(ErrCriteriaNotMet, "territory listed multiple times")
		}
		seen[territory[i]] = struct{}{}
	}
	validFrom := GetValidFrom(right)
	validThrough := GetValidThrough(right)
	if validFrom.After(validThrough) {
		return ErrorAppend(ErrInvalidTime, "range")
	}
	if validThrough.Before(Now()) {
		return ErrorAppend(ErrInvalidTime, "expired")
	}
	return nil
}

func ValidCompositionRight(right Data) error {
	if err := ValidRight(right); err != nil {
		return err
	}
	compositionId := GetCompositionId(right)
	if !MatchId(compositionId) {
		return ErrorAppend(ErrInvalidId, "compositionId")
	}
	if n := len(right); n != RIGHT_SIZE {
		return ErrorAppend(ErrInvalidSize, Itoa(n))
	}
	return nil
}

func ValidRecordingRight(right Data) error {
	if err := ValidRight(right); err != nil {
		return err
	}
	recordingId := GetRecordingId(right)
	if !MatchId(recordingId) {
		return ErrorAppend(ErrInvalidId, "recordingId")
	}
	if n := len(right); n != RIGHT_SIZE {
		return ErrorAppend(ErrInvalidSize, Itoa(n))
	}
	return nil
}

// License

func NewLicense(recipientId, senderId string, usage, territory []string, validFrom, validThrough string) Data {
	return Data{
		"instance":     NewInstance(LICENSE),
		"recipientId":  recipientId,
		"senderId":     senderId,
		"usage":        usage,
		"validFrom":    validFrom,
		"validThrough": validThrough,
		"territory":    territory,
	}
}

func NewMechanicalLicense(compositionRightId, compositionRightTransferId, publicationId, recipientId, senderId string, usage, territory []string, validFrom, validThrough string) Data {
	license := NewLicense(recipientId, senderId, usage, territory, validFrom, validThrough)
	license.Set("publicationId", publicationId)
	if !EmptyStr(compositionRightId) {
		license.Set("compositionRightId", compositionRightId)
	} else if !EmptyStr(compositionRightTransferId) {
		license.Set("compositionRightTransferId", compositionRightTransferId)
	} else {
		panic("Expected compositionRightId or compositionRightTransferId")
	}
	return license
}

func NewMasterLicense(recipientId, recordingRightId, recordingRightTransferId, releaseId, senderId string, usage, territory []string, validFrom, validThrough string) Data {
	license := NewLicense(recipientId, senderId, usage, territory, validFrom, validThrough)
	license.Set("releaseId", releaseId)
	if !EmptyStr(recordingRightId) {
		license.Set("recordingRightId", recordingRightId)
	} else if !EmptyStr(recordingRightTransferId) {
		license.Set("recordingRightTransferId", recordingRightTransferId)
	} else {
		panic("Expected recordingRightId or recordingRightTransferId")
	}
	return license
}

func GetReleaseId(data Data) string {
	return data.GetStr("releaseId")
}

func GetCompositionRightTransferId(data Data) string {
	return data.GetStr("compositionRightTransferId")
}

func GetRecordingRightTransferId(data Data) string {
	return data.GetStr("recordingRightTransferId")
}

func ValidMechanicalLicense(license Data) error {
	if err := ValidLicense(license); err != nil {
		return err
	}
	compositionRightId := GetCompositionRightId(license)
	if !EmptyStr(compositionRightId) {
		if !MatchId(compositionRightId) {
			return ErrorAppend(ErrInvalidId, compositionRightId)
		}
	} else {
		compositionRightTransferId := GetCompositionRightTransferId(license)
		if !MatchId(compositionRightTransferId) {
			return ErrorAppend(ErrInvalidId, compositionRightTransferId)
		}
	}
	publicationId := GetPublicationId(license)
	if !MatchId(publicationId) {
		return ErrorAppend(ErrInvalidId, publicationId)
	}
	if n := len(license); n != LICENSE_SIZE {
		return ErrorAppend(ErrInvalidSize, Itoa(n))
	}
	return nil
}

func ValidMasterLicense(license Data) error {
	if err := ValidLicense(license); err != nil {
		return err
	}
	recordingRightId := GetRecordingRightId(license)
	if !EmptyStr(recordingRightId) {
		if !MatchId(recordingRightId) {
			return ErrorAppend(ErrInvalidId, recordingRightId)
		}
	} else {
		recordingRightTransferId := GetRecordingRightTransferId(license)
		if !MatchId(recordingRightTransferId) {
			return ErrorAppend(ErrInvalidId, recordingRightTransferId)
		}
	}
	releaseId := GetReleaseId(license)
	if !MatchId(releaseId) {
		return ErrorAppend(ErrInvalidId, releaseId)
	}
	if n := len(license); n != LICENSE_SIZE {
		return ErrorAppend(ErrInvalidSize, Itoa(n))
	}
	return nil
}

func ValidLicense(license Data) error {
	instance := GetInstance(license)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(license, LICENSE) {
		return ErrorAppend(ErrInvalidType, GetType(license))
	}
	recipientId := GetRecipientId(license)
	if !MatchId(recipientId) {
		return ErrorAppend(ErrInvalidId, recipientId)
	}
	senderId := GetSenderId(license)
	if !MatchId(senderId) {
		return ErrorAppend(ErrInvalidId, senderId)
	}
	territory := GetTerritory(license)
	if len(territory) == 0 {
		return Error("no territory listed")
	}
	seen := make(map[string]struct{})
	for i := range territory {
		if !MatchStr(TERRITORY_REGEX, territory[i]) {
			return ErrInvalidTerritory
		}
		if _, ok := seen[territory[i]]; ok {
			return ErrorAppend(ErrCriteriaNotMet, "territory listed multiple times")
		}
		seen[territory[i]] = struct{}{}
	}
	// TODO: check usage
	validFrom := GetValidFrom(license)
	validThrough := GetValidThrough(license)
	if validFrom.After(validThrough) {
		return ErrInvalidTime
	}
	return nil
}

// Right Transfer

func NewRightTransfer(recipientId, senderId, txId string) Data {
	return Data{
		"instance":    NewInstance(TRANSFER),
		"recipientId": recipientId,
		"senderId":    senderId,
		"txId":        txId,
	}
}

func NewCompositionRightTransfer(compositionRightId, publicationId, recipientId, senderId, txId string) Data {
	transfer := NewRightTransfer(recipientId, senderId, txId)
	transfer.Set("compositionRightId", compositionRightId)
	transfer.Set("publicationId", publicationId)
	return transfer
}

func NewRecordingRightTransfer(recipientId, recordingRightId, releaseId, senderId, txId string) Data {
	transfer := NewRightTransfer(recipientId, senderId, txId)
	transfer.Set("recordingRightId", recordingRightId)
	transfer.Set("releaseId", releaseId)
	return transfer
}

func GetRecipientShares(data Data) int {
	return data.GetInt("recipientShares")
}

func GetRecordingRightId(data Data) string {
	return data.GetStr("recordingRightId")
}

func GetSenderShares(data Data) int {
	return data.GetInt("senderShares")
}

func GetTxId(data Data) string {
	return data.GetStr("txId")
}

func ValidCompositionRightTransfer(transfer Data) error {
	if err := ValidRightTransfer(transfer); err != nil {
		return err
	}
	compositionRightId := GetCompositionRightId(transfer)
	if !MatchId(compositionRightId) {
		return ErrorAppend(ErrInvalidId, compositionRightId)
	}
	publicationId := GetPublicationId(transfer)
	if !MatchId(publicationId) {
		return ErrorAppend(ErrInvalidId, publicationId)
	}
	if n := len(transfer); n != TRANSFER_SIZE {
		return ErrorAppend(ErrInvalidSize, Itoa(n))
	}
	return nil
}

func ValidRecordingRightTransfer(transfer Data) error {
	if err := ValidRightTransfer(transfer); err != nil {
		return err
	}
	recordingRightId := GetRecordingRightId(transfer)
	if !MatchId(recordingRightId) {
		return ErrorAppend(ErrInvalidId, recordingRightId)
	}
	releaseId := GetReleaseId(transfer)
	if !MatchId(releaseId) {
		return ErrorAppend(ErrInvalidId, releaseId)
	}
	if n := len(transfer); n != TRANSFER_SIZE {
		return ErrorAppend(ErrInvalidSize, Itoa(n))
	}
	return nil
}

func ValidRightTransfer(transfer Data) error {
	instance := GetInstance(transfer)
	if err := ValidInstance(instance); err != nil {
		return err
	}
	if !HasType(transfer, TRANSFER) {
		return ErrorAppend(ErrInvalidType, GetType(transfer))
	}
	recipientId := GetRecipientId(transfer)
	if !MatchId(recipientId) {
		return ErrorAppend(ErrInvalidId, recipientId)
	}
	senderId := GetSenderId(transfer)
	if !MatchId(senderId) {
		return ErrorAppend(ErrInvalidId, senderId)
	}
	if recipientId == senderId {
		return ErrorAppend(ErrCriteriaNotMet, "recipientId and senderId must be different")
	}
	txId := GetTxId(transfer)
	if !MatchId(txId) {
		return ErrorAppend(ErrInvalidId, txId)
	}
	return nil
}
