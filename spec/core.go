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

	ENTITY_SIZE     = 2
	AGENT_SIZE      = 4
	BASE_ALBUM_SIZE = 3
	BASE_TRACK_SIZE = 3
	SIGNATURE_SIZE  = 4

	EMAIL_REGEX           = `(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$)`
	FINGERPRINT_STD_REGEX = `^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$` // base64 std
	FINGERPRINT_URL_REGEX = `^(?:[A-Za-z0-9-_]{4})*(?:[A-Za-z0-9-_]{2}==|[A-Za-z0-9-_]{3})?$`  // base64 url-safe
	ID_REGEX              = `^[A-Fa-f0-9]{64}$`                                                // hex
	PUBKEY_REGEX          = `^[1-9A-HJ-NP-Za-km-z]{43,44}$`                                    // base58
	SIGNATURE_REGEX       = `^[1-9A-HJ-NP-Za-km-z]{87,88}$`                                    // base58
)

// Instance

func NewInstance(_type string) Data {
	return Data{
		"time": Timestamp(),
		"type": _type,
	}
}

func GetInstanceTime(instance Data) int64 {
	return instance.GetInt64("type")
}

func GetInstanceType(instance Data) string {
	return instance.GetStr("type")
}

func GetInstance(thing Data) Data {
	if ValidInstance(thing) {
		return thing
	}
	return thing.GetData("instance")
}

func GetType(thing Data) string {
	instance := GetInstance(thing)
	return GetInstanceType(instance)
}

func HasType(thing Data, _type string) bool {
	return GetType(thing) == _type
}

func ValidInstance(instance Data) bool {
	time := GetInstanceTime(instance)
	if time > Timestamp() {
		return false
	}
	_type := GetInstanceType(instance)
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
	return len(instance) == 2
}

// Agent

func NewAgent(email, name string, pub crypto.PublicKey, _type string) Data {
	return Data{
		"email":      email,
		"instance":   NewInstance(_type),
		"name":       name,
		"public_key": pub.String(),
	}
}

func GetAgentEmail(agent Data) string {
	return agent.GetStr("email")
}

func GetAgentName(agent Data) string {
	return agent.GetStr("name")
}

func GetAgentPublicKey(agent Data) crypto.PublicKey {
	pubstr := GetAgentPublicKeyStr(agent)
	if Empty(pubstr) {
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
	instance := GetInstance(agent)
	if !ValidInstance(instance) {
		return false
	}
	_type := GetInstanceType(instance)
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
	if Empty(name) {
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

func NewAlbum(artistId, labelId, publisherId, title string) Data {
	album := Data{
		"artist_id": artistId,
		"instance":  NewInstance(ALBUM),
		"title":     title,
	}
	if !Empty(labelId) {
		album.Set("label_id", labelId)
	}
	if !Empty(publisherId) {
		album.Set("publisher_id", publisherId)
	}
	return album
}

func GetMusicAgent(music Data, _type string) string {
	switch _type {
	case ARTIST:
		return GetMusicArtist(music)
	case LABEL:
		return GetMusicLabel(music)
	case PUBLISHER:
		return GetMusicPublisher(music)
	default:
		return ""
	}
}

func GetMusicArtist(music Data) string {
	return music.GetStr("artist_id")
}

func GetMusicLabel(music Data) string {
	return music.GetStr("label_id")
}

func GetMusicPublisher(music Data) string {
	return music.GetStr("publisher_id")
}

func GetMusicTitle(music Data) string {
	return music.GetStr("title")
}

func ValidAlbum(album Data) bool {
	count := BASE_ALBUM_SIZE
	instance := GetInstance(album)
	if !ValidInstance(instance) {
		return false
	}
	if GetInstanceType(instance) != ALBUM {
		return false
	}
	artistId := GetMusicArtist(album)
	if !MatchString(ID_REGEX, artistId) {
		return false
	}
	labelId := GetMusicLabel(album)
	if !Empty(labelId) {
		if !MatchString(ID_REGEX, labelId) {
			return false
		}
		count++
	}
	publisherId := GetMusicPublisher(album)
	if !Empty(publisherId) {
		if !MatchString(ID_REGEX, publisherId) {
			return false
		}
		count++
	}
	title := GetMusicTitle(album)
	if Empty(title) {
		return false
	}
	return len(album) != count
}

func NewTrack(albumId, artistId, fingerprint, labelId string, number int, publisherId, title string) Data {
	track := Data{
		"instance":    NewInstance(TRACK),
		"fingerprint": fingerprint,
		"title":       title,
	}
	// Track must have album_id and track_number or artist_id
	if albumId != "" && number > 0 {
		track.Set("album_id", albumId)
		track.Set("track_number", number)
	} else {
		track.Set("artist_id", artistId)
	}
	if !Empty(labelId) {
		track.Set("label_id", labelId)
	}
	if !Empty(publisherId) {
		track.Set("publisher_id", publisherId)
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
	count := BASE_TRACK_SIZE
	instance := GetInstance(track)
	if !ValidInstance(instance) {
		return false
	}
	if !HasType(track, TRACK) {
		return false
	}
	// TODO: better fingerprint validation?
	fingerprint := GetTrackFingerprint(track)
	if !MatchString(FINGERPRINT_URL_REGEX, fingerprint) {
		return false
	}
	labelId := GetMusicLabel(track)
	if !Empty(labelId) {
		if !MatchString(ID_REGEX, labelId) {
			return false
		}
		count++
	}
	publisherId := GetMusicPublisher(track)
	if !Empty(publisherId) {
		if !MatchString(ID_REGEX, publisherId) {
			return false
		}
		count++
	}
	title := GetMusicTitle(track)
	if Empty(title) {
		return false
	}
	artistId := GetMusicArtist(track)
	if MatchString(ID_REGEX, artistId) {
		count++
		return len(track) == count
	}
	albumId := GetTrackAlbum(track)
	if !MatchString(ID_REGEX, albumId) {
		return false
	}
	count++
	trackNumber := GetTrackNumber(track)
	if trackNumber <= 0 {
		return false
	}
	count++
	return len(track) == count
}

// Signature

func NewSignature(modelId, signerId string, sig crypto.Signature) Data {
	return Data{
		"instance":  NewInstance(SIGNATURE),
		"model_id":  modelId,
		"signer_id": signerId,
		"value":     sig.String(),
	}
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
	instance := GetInstance(signature)
	if !ValidInstance(instance) {
		return false
	}
	if GetInstanceType(instance) != SIGNATURE {
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
