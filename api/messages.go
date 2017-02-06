package api

import (
	"github.com/zbo14/envoke/crypto/crypto"
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
	UserId  string            `json:"partner_id"`
	Privkey crypto.PrivateKey `json:"private_key"`
	Pubkey  crypto.PublicKey  `json:"public_key"`
}

func NewPartnerInfo(userId string, priv crypto.PrivateKey, pub crypto.PublicKey) *PartnerInfo {
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
