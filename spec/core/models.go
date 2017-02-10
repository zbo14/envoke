package core

import (
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/crypto"
)

const (
	ARTIST       = "artist"
	LABEL        = "label"
	ORGANIZATION = "organization"
	PUBLISHER    = "publisher"

	ALBUM = "album"
	TRACK = "track"

	SIGNATURE = "signature"

	EMAIL_REGEX     = `(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$)`
	ID_REGEX        = `[A-Fa-f0-9]{64}`
	PUBKEY_REGEX    = `^[1-9A-HJ-NP-Za-km-z]{43,44}$`
	SIGNATURE_REGEX = `^[1-9A-HJ-NP-Za-km-z]{87}$`
)

// TODO: replace regex validation with id queries and data model validation

// Entity

type Entity struct {
	Time int64  `json:"time"`
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

// Agent

type Agent struct {
	Email string `json:"email"`
	*Entity
	Name   string           `json:"name"`
	PubKey crypto.PublicKey `json:"public_key"`
}

func NewAgent(email, name string, pub crypto.PublicKey, _type string) *Agent {
	return &Agent{
		Email:  email,
		Entity: NewEntity(_type),
		Name:   name,
		PubKey: pub,
	}
}

func NewArtist(email, name string, pub crypto.PublicKey) *Agent {
	return NewAgent(email, name, pub, ARTIST)
}

func NewLabel(email, name string, pub crypto.PublicKey) *Agent {
	return NewAgent(email, name, pub, LABEL)
}

func NewOrganization(email, name string, pub crypto.PublicKey) *Agent {
	return NewAgent(email, name, pub, ORGANIZATION)
}

func NewPublisher(email, name string, pub crypto.PublicKey) *Agent {
	return NewAgent(email, name, pub, PUBLISHER)
}

func (a *Agent) IsArtist() bool {
	return a.Type == ARTIST
}

func (a *Agent) IsPartner() bool {
	switch a.Type {
	case LABEL, ORGANIZATION, PUBLISHER:
		return true
	default:
		return false
	}
}

func (a *Agent) Valid() bool {
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
	if !MatchString(PUBKEY_REGEX, a.PubKey.String()) {
		return false
	}
	return true
}

// Music

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
	Number      int    `json:"track_number,omitempty"`
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

// Signature

type Signature struct {
	*Entity
	SignerId string
	Value    crypto.Signature
}

func NewSignature(signerId string, value crypto.Signature) *Signature {
	return &Signature{
		Entity:   NewEntity(SIGNATURE),
		SignerId: signerId,
		Value:    value,
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
	if !MatchString(SIGNATURE_REGEX, s.Value.String()) {
		return false
	}
	return true
}
