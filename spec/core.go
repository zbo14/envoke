package spec

import (
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/crypto"
	"github.com/zbo14/envoke/crypto/ed25519"
)

const (
	ARTIST       = "artist"
	LABEL        = "label"
	ORGANIZATION = "organization"
	PUBLISHER    = "publisher"

	ALBUM = "album"
	TRACK = "track"

	SIGNATURE = "signature"

	ENTITY_SIZE       = 2
	AGENT_SIZE        = 4
	ALBUM_SIZE        = 4
	TRACK_ALBUM_SIZE  = 6
	TRACK_SINGLE_SIZE = 5
	SIGNATURE_SIZE    = 4

	EMAIL_REGEX           = `(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$)`
	FINGERPRINT_STD_REGEX = `^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$` // base64 std
	FINGERPRINT_URL_REGEX = `^(?:[A-Za-z0-9-_]{4})*(?:[A-Za-z0-9-_]{2}==|[A-Za-z0-9-_]{3})?$`  // base64 url-safe
	ID_REGEX              = `^[A-Fa-f0-9]{64}$`                                                // hex
	PUBKEY_REGEX          = `^[1-9A-HJ-NP-Za-km-z]{43,44}$`                                    // base58
	SIGNATURE_REGEX       = `^[1-9A-HJ-NP-Za-km-z]{87,88}$`                                    // base58
)

// Entity

func NewEntity(_type string) Data {
	entity := Data{
		"time": Timestamp(),
		"type": _type,
	}
	if !ValidEntity(entity) {
		panic("Invalid entity")
	}
	return entity
}

func GetEntityTime(entity Data) int64 {
	return entity.GetInt64("type")
}

func GetEntityType(entity Data) string {
	return entity.GetStr("type")
}

func GetEntity(thing Data) Data {
	if ValidEntity(thing) {
		return thing
	}
	return thing.GetData("entity")
}

func GetType(thing Data) string {
	entity := GetEntity(thing)
	return GetEntityType(entity)
}

func HasType(thing Data, _type string) bool {
	return GetType(thing) == _type
}

func ValidEntity(entity Data) bool {
	time := GetEntityTime(entity)
	if time > Timestamp() {
		return false
	}
	_type := GetEntityType(entity)
	switch _type {
	case
		ARTIST, LABEL, ORGANIZATION, PUBLISHER,
		ALBUM, TRACK,
		SIGNATURE,
		RIGHT:
		// Ok..
	default:
		return false
	}
	return len(entity) == ENTITY_SIZE
}

// Agent

func NewAgent(email, name string, pub crypto.PublicKey, _type string) Data {
	agent := Data{
		"email":      email,
		"entity":     NewEntity(_type),
		"name":       name,
		"public_key": pub.String(),
	}
	if !ValidAgent(agent) {
		panic("Invalid agent")
	}
	return agent
}

func GetAgentEmail(agent Data) string {
	return agent.GetStr("email")
}

func GetAgentName(agent Data) string {
	return agent.GetStr("name")
}

func GetAgentPublicKey(agent Data) crypto.PublicKey {
	pubstr := GetAgentPublicKeyStr(agent)
	if pubstr == "" {
		return nil
	}
	pub := new(ed25519.PublicKey)
	err := pub.FromString(pubstr)
	Check(err)
	return pub
}

func GetAgentPublicKeyStr(agent Data) string {
	return agent.GetStr("public_key")
}

func NewArtist(email, name string, pub crypto.PublicKey) Data {
	return NewAgent(email, name, pub, ARTIST)
}

func NewLabel(email, name string, pub crypto.PublicKey) Data {
	return NewAgent(email, name, pub, LABEL)
}

func NewOrganization(email, name string, pub crypto.PublicKey) Data {
	return NewAgent(email, name, pub, ORGANIZATION)
}

func NewPublisher(email, name string, pub crypto.PublicKey) Data {
	return NewAgent(email, name, pub, PUBLISHER)
}

func ValidArtist(agent Data) bool {
	if !ValidAgent(agent) {
		return false
	}
	return HasType(agent, ARTIST)
}

func ValidLabel(agent Data) bool {
	if !ValidAgent(agent) {
		return false
	}
	return HasType(agent, LABEL)
}

func ValidOrganization(agent Data) bool {
	if !ValidAgent(agent) {
		return false
	}
	return HasType(agent, ORGANIZATION)
}

func ValidPublisher(agent Data) bool {
	if !ValidAgent(agent) {
		return false
	}
	return HasType(agent, PUBLISHER)
}

func ValidAgent(agent Data) bool {
	entity := GetEntity(agent)
	if !ValidEntity(entity) {
		return false
	}
	_type := GetEntityType(entity)
	switch _type {
	case ARTIST, LABEL, ORGANIZATION, PUBLISHER:
		// Ok..
	default:
		return false
	}
	email := GetAgentEmail(agent)
	if !MatchString(EMAIL_REGEX, email) {
		return false
	}
	name := GetAgentName(agent)
	if name == "" {
		return false
	}
	pubstr := GetAgentPublicKeyStr(agent)
	if !MatchString(PUBKEY_REGEX, pubstr) {
		return false
	}
	return len(agent) == AGENT_SIZE
}

func ValidAgentWithType(agent Data, _type string) bool {
	if !ValidAgent(agent) {
		return false
	}
	return HasType(agent, _type)
}

// Music

func ValidMusic(music Data) bool {
	_type := GetType(music)
	if _type == ALBUM {
		return ValidAlbum(music)
	}
	if _type == TRACK {
		return ValidTrack(music)
	}
	return false
}

func NewAlbum(artistId, publisherId, title string) Data {
	album := Data{
		"artist_id":    artistId,
		"entity":       NewEntity(ALBUM),
		"publisher_id": publisherId,
		"title":        title,
	}
	if !ValidAlbum(album) {
		panic("Invalid album")
	}
	return album
}

func GetMusicArtist(music Data) string {
	return music.GetStr("artist_id")
}

func GetMusicPublisher(music Data) string {
	return music.GetStr("publisher_id")
}

func GetMusicTitle(music Data) string {
	return music.GetStr("title")
}

func ValidAlbum(album Data) bool {
	entity := GetEntity(album)
	if !ValidEntity(entity) {
		return false
	}
	if GetEntityType(entity) != ALBUM {
		return false
	}
	artistId := GetMusicArtist(album)
	if !MatchString(ID_REGEX, artistId) {
		return false
	}
	publisherId := GetMusicPublisher(album)
	if !MatchString(ID_REGEX, publisherId) {
		return false
	}
	title := GetMusicTitle(album)
	if title == "" {
		return false
	}
	return len(album) == ALBUM_SIZE
}

func NewTrack(albumId, artistId, fingerprint string, number int, publisherId, title string) Data {
	track := Data{
		"artist_id":   artistId,
		"entity":      NewEntity(TRACK),
		"fingerprint": fingerprint,
		"title":       title,
	}
	if publisherId == "" {
		track["album_id"] = albumId
		track["track_number"] = number
	} else {
		track["publisher_id"] = publisherId
	}
	if !ValidTrack(track) {
		panic("Invalid track")
	}
	return track
}

func GetTrackAlbum(track Data) string {
	return track.GetStr("album_id")
}

func GetTrackFingerprint(track Data) string {
	return track.GetStr("fingerprint")
}

func GetTrackNumber(track Data) int {
	return track.GetInt("track_number")
}

func ValidTrack(track Data) bool {
	entity := GetEntity(track)
	if !ValidEntity(entity) {
		return false
	}
	if !HasType(track, TRACK) {
		return false
	}
	artistId := GetMusicArtist(track)
	if !MatchString(ID_REGEX, artistId) {
		return false
	}
	// TODO: better fingerprint validation?
	fingerprint := GetTrackFingerprint(track)
	if !MatchString(FINGERPRINT_URL_REGEX, fingerprint) {
		return false
	}
	title := GetMusicTitle(track)
	if title == "" {
		return false
	}
	publisherId := GetMusicPublisher(track)
	if MatchString(ID_REGEX, publisherId) {
		return len(track) == TRACK_SINGLE_SIZE
	}
	albumId := GetTrackAlbum(track)
	if !MatchString(ID_REGEX, albumId) {
		Println(6)
		return false
	}
	trackNumber := GetTrackNumber(track)
	if trackNumber <= 0 {
		Println(7)
		return false
	}
	return len(track) == TRACK_ALBUM_SIZE
}

// Signature

func NewSignature(modelId, signerId string, sig crypto.Signature) Data {
	signature := Data{
		"entity":    NewEntity(SIGNATURE),
		"model_id":  modelId,
		"signer_id": signerId,
		"value":     sig.String(),
	}
	if !ValidSignature(signature) {
		panic("Invalid signature")
	}
	return signature
}

func GetSignatureModel(signature Data) string {
	return signature.GetStr("model_id")
}

func GetSignatureSigner(signature Data) string {
	return signature.GetStr("signer_id")
}

func GetSignatureValue(signature Data) crypto.Signature {
	sigstr := GetSignatureValueStr(signature)
	if sigstr == "" {
		return nil
	}
	sig := new(ed25519.Signature)
	err := sig.FromString(sigstr)
	Check(err)
	return sig
}

func GetSignatureValueStr(signature Data) string {
	return signature.GetStr("value")
}

func ValidSignature(signature Data) bool {
	entity := GetEntity(signature)
	if !ValidEntity(entity) {
		return false
	}
	if GetEntityType(entity) != SIGNATURE {
		return false
	}
	signerId := GetSignatureSigner(signature)
	if !MatchString(ID_REGEX, signerId) {
		return false
	}
	modelId := GetSignatureModel(signature)
	if !MatchString(ID_REGEX, modelId) {
		return false
	}
	valueStr := GetSignatureValueStr(signature)
	if !MatchString(SIGNATURE_REGEX, valueStr) {
		return false
	}
	return len(signature) == SIGNATURE_SIZE
}
