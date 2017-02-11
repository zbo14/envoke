package core

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
	TRACK_ALBUM_SIZE  = 5
	TRACK_SINGLE_SIZE = 4
	SIGNATURE_SIZE    = 4

	EMAIL_REGEX     = `(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$)`
	ID_REGEX        = `[A-Fa-f0-9]{64}`
	PUBKEY_REGEX    = `^[1-9A-HJ-NP-Za-km-z]{43,44}$`
	SIGNATURE_REGEX = `^[1-9A-HJ-NP-Za-km-z]{87}$`
)

type Data map[string]interface{}

func AssertData(v interface{}) Data {
	if data, ok := v.(Data); ok {
		return data
	}
	return nil
}

func AssertInt(v interface{}) int {
	if n, ok := v.(int); ok {
		return n
	}
	return 0
}

func AssertInt64(v interface{}) int64 {
	if n, ok := v.(int64); ok {
		return n
	}
	return 0
}

func AssertStr(v interface{}) string {
	if str, ok := v.(string); ok {
		return str
	}
	return ""
}

// Entity

func NewEntity(_type string) Data {
	return Data{
		"time": Timestamp(),
		"type": _type,
	}
}

func GetEntityTime(entity Data) int64 {
	return AssertInt64(entity["time"])
}

func GetEntityType(entity Data) string {
	return AssertStr(entity["type"])
}

func GetEntity(thing Data) Data {
	if ValidEntity(thing) {
		return thing
	}
	return AssertData(thing["entity"])
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
	if time >= Timestamp() {
		return false
	}
	_type := GetEntityType(entity)
	switch _type {
	case
		ARTIST, LABEL, ORGANIZATION, PUBLISHER,
		ALBUM, TRACK, SIGNATURE:
		// Ok..
	default:
		return false
	}
	return len(entity) == ENTITY_SIZE
}

// Agent

func NewAgent(email, name string, pub crypto.PublicKey, _type string) Data {
	return Data{
		"email":      email,
		"entity":     NewEntity(_type),
		"name":       name,
		"public_key": pub.String(),
	}
}

func GetAgentEmail(agent Data) string {
	return AssertStr(agent["email"])
}

func GetAgentName(agent Data) string {
	return AssertStr(agent["name"])
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
	return AssertStr(agent["public_key"])
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
	switch GetType(music) {
	case ALBUM:
		return ValidAlbum(music)
	case TRACK:
		return ValidTrack(music)
	default:
		return false
	}
}

func NewAlbum(artistId, publisherId, title string) Data {
	return Data{
		"artist_id":    artistId,
		"entity":       NewEntity(ALBUM),
		"publisher_id": publisherId,
		"title":        title,
	}
}

func GetMusicArtist(music Data) string {
	return AssertStr(music["artist_id"])
}

func GetMusicPublisher(music Data) string {
	return AssertStr(music["publisher_id"])
}

func GetMusicTitle(music Data) string {
	return AssertStr(music["title"])
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

func NewTrack(albumId, artistId string, number int, publisherId, title string) Data {
	if publisherId == "" {
		if albumId == "" || number <= 0 {
			panic("")
		}
	}
	data := Data{
		"artist_id": artistId,
		"entity":    NewEntity(TRACK),
		"title":     title,
	}
	if publisherId != "" {
		data["publisher_id"] = publisherId
	} else if albumId != "" && number > 0 {
		data["album_id"] = albumId
		data["track_number"] = number
	}
	return data
}

func GetTrackAlbum(track Data) string {
	return AssertStr(track["album_id"])
}

func GetTrackNumber(track Data) int {
	return AssertInt(track["track_number"])
}

func ValidTrack(track Data) bool {
	entity := GetEntity(track)
	if !ValidEntity(entity) {
		return false
	}
	if GetEntityType(entity) != TRACK {
		return false
	}
	artistId := GetMusicArtist(track)
	if !MatchString(ID_REGEX, artistId) {
		return false
	}
	title := GetMusicTitle(track)
	if !MatchString(ID_REGEX, title) {
		return false
	}
	publisherId := GetMusicPublisher(track)
	if MatchString(ID_REGEX, publisherId) {
		return len(track) == TRACK_SINGLE_SIZE
	}
	albumId := GetTrackAlbum(track)
	if !MatchString(ID_REGEX, albumId) {
		return false
	}
	trackNumber := GetTrackNumber(track)
	if trackNumber <= 0 {
		return false
	}
	return len(track) == TRACK_ALBUM_SIZE
}

// Signature

func NewSignature(agentId string, musicId string, sig crypto.Signature) Data {
	return Data{
		"agent_id": agentId,
		"entity":   NewEntity(SIGNATURE),
		"music_id": musicId,
		"value":    sig.String(),
	}
}

func GetSignatureAgent(signature Data) string {
	return AssertStr(signature["agent_id"])
}

func GetSignatureMusic(signature Data) string {
	return AssertStr(signature["music_id"])
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
	return AssertStr(signature["value"])
}

func ValidSignature(signature Data) bool {
	entity := GetEntity(signature)
	if !ValidEntity(entity) {
		return false
	}
	if GetEntityType(entity) != SIGNATURE {
		return false
	}
	agentId := GetSignatureAgent(signature)
	if !MatchString(ID_REGEX, agentId) {
		return false
	}
	valueStr := GetSignatureValueStr(signature)
	if !MatchString(SIGNATURE_REGEX, valueStr) {
		return false
	}
	return len(signature) == SIGNATURE_SIZE
}

/*
type Signature struct {
	*Entity
	SignerId string
	Value    string
}

func NewSignature(sig crypto.Signature, signerId string) *Signature {
	return &Signature{
		Entity:   NewEntity(SIGNATURE),
		SignerId: signerId,
		Value:    sig.String(),
	}
}

func (s *Signature) Valid() bool {
	if !s.Entity.Valid() {
		return false
	}
	if s.Type != SIGNATURE {
		return false
	}
	if !MatchString(ID_REGEX, s.SignerId) {
		return false
	}
	if !MatchString(SIGNATURE_REGEX, s.Value) {
		return false
	}
	return true
}
*/

/*
type Entity struct {
	Time int64  `json:"time,string"`
	Type string `json:"type"`
}

func NewEntity(_type string) *Entity {
	return &Entity{
		Time: Timestamp(),
		Type: _type,
	}
}

func (e *Entity) Valid() bool {
	if e.Time >= Timestamp() {
		return false
	}
	switch e.Type {
	case
		ARTIST, LABEL, ORGANIZATION, PUBLISHER,
		ALBUM, TRACK, SIGNATURE:
		return true
	default:
		return false
	}
}

type Agent struct {
	Email string `json:"email"`
	*Entity
	Name   string `json:"name"`
	PubKey string `json:"public_key"`
}

func NewAgent(email, name string, pub crypto.PublicKey, _type string) Data {
	return &Agent{
		Email:  email,
		Entity: NewEntity(_type),
		Name:   name,
		PubKey: pub.String(),
	}
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

func (a Data) IsArtist() bool {
	return a.Type == ARTIST
}

func (a Data) IsPartner() bool {
	switch a.Type {
	case LABEL, ORGANIZATION, PUBLISHER:
		return true
	default:
		return false
	}
}

func (a Data) Valid() bool {
	if !a.Entity.Valid() {
		return false
	}
	switch a.Type {
	case ARTIST, LABEL, ORGANIZATION, PUBLISHER:
		// Ok..
	default:
		return false
	}
	if !MatchString(EMAIL_REGEX, a.Email) {
		return false
	}
	if a.Name == "" {
		return false
	}
	if !MatchString(PUBKEY_REGEX, a.PubKey) {
		return false
	}
	return true
}

func GetMusicPublisher(music interface{}) (string, error) {
	switch music.(type) {
	case *Album:
		return music.(*Album).PublisherId, nil
	case *Track:
		return music.(*Track).PublisherId, nil
	default:
		return "", ErrInvalidType
	}
}

type Album struct {
	ArtistId string `json:"artist_id"`
	*Entity
	PublisherId string `json:"publisher_id"`
	Title       string `json:"title"`
}

func NewAlbum(artistId, publisherId, title string) *Album {
	return &Album{
		ArtistId:    artistId,
		Entity:      NewEntity(ALBUM),
		PublisherId: publisherId,
		Title:       title,
	}
}

func (a *Album) Valid() bool {
	if !a.Entity.Valid() {
		return false
	}
	if a.Type != ALBUM {
		return false
	}
	if !MatchString(ID_REGEX, a.ArtistId) {
		return false
	}
	if !MatchString(ID_REGEX, a.PublisherId) {
		return false
	}
	if a.Title == "" {
		return false
	}
	return true
}

type Track struct {
	AlbumId  string `json:"album_id,omitempty"`
	ArtistId string `json:"artist_id"`
	*Entity
	Number      int    `json:"track_number,string,omitempty"`
	PublisherId string `json:"publisher_id,omitempty"`
	Title       string `json:"title"`
}

func NewTrack(albumId, artistId string, number int, publisherId, title string) *Track {
	return &Track{
		AlbumId:     albumId,
		ArtistId:    artistId,
		Entity:      NewEntity(TRACK),
		Number:      number,
		PublisherId: publisherId,
		Title:       title,
	}
}

func (t *Track) Valid() bool {
	if !t.Entity.Valid() {
		return false
	}
	if t.Type != TRACK {
		return false
	}
	if !MatchString(ID_REGEX, t.ArtistId) {
		return false
	}
	if t.Title == "" {
		return false
	}
	if MatchString(ID_REGEX, t.PublisherId) {
		return true
	}
	if !MatchString(ID_REGEX, t.AlbumId) {
		return false
	}
	if t.Number <= 0 {
		return false
	}
	return true
}
*/
