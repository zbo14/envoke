package api

import (
	"github.com/zballs/envoke/crypto/ed25519"
)

type AlbumInfo struct {
	AlbumId  string   `json:"album_id"`
	TrackIds []string `json:"track_ids"`
}

func NewAlbumInfo(albumId string, songIds []string) *AlbumInfo {
	return &AlbumInfo{
		AlbumId:  albumId,
		TrackIds: songIds,
	}
}

type PartnerInfo struct {
	UserId  string              `json:"user_id"`
	Privkey *ed25519.PrivateKey `json:"private_key"`
	Pubkey  *ed25519.PublicKey  `json:"public_key"`
}

func NewPartnerInfo(userId string, priv *ed25519.PrivateKey, pub *ed25519.PublicKey) *PartnerInfo {
	return &PartnerInfo{
		UserId:  userId,
		Privkey: priv,
		Pubkey:  pub,
	}
}

type Stream struct {
	Artist      string `json:"artist"`
	ProjecTitle string `json:"album_title"`
	TrackTitle  string `json:"track_title"`
	URL         string `json:"url"`
}

func NewStream(artist, albumTitle, trackTitle, url string) *Stream {
	return &Stream{
		Artist:      artist,
		ProjecTitle: albumTitle,
		TrackTitle:  trackTitle,
		URL:         url,
	}
}
