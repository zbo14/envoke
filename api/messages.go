package api

import (
	"github.com/zballs/envoke/crypto/ed25519"
)

type UserInfo struct {
	Privkey *ed25519.PrivateKey `json:"private_key"`
	Pubkey  *ed25519.PublicKey  `json:"public_key"`
	UserId  string              `json:"user_id"`
}

func NewUserInfo(userId string, priv *ed25519.PrivateKey, pub *ed25519.PublicKey) *UserInfo {
	return &UserInfo{
		UserId:  userId,
		Privkey: priv,
		Pubkey:  pub,
	}
}

type Login struct {
	Message  string `json:"message"`
	UserType string `json:"user_type"`
}

func NewLogin(_type string) *Login {
	return &Login{
		Message:  "Logged in!",
		UserType: _type,
	}
}

type ProjectInfo struct {
	ProjectId string   `json:"project_id"`
	SongIds   []string `json:"track_ids"`
}

func NewProjectInfo(projectId string, songIds []string) *ProjectInfo {
	return &ProjectInfo{
		ProjectId: projectId,
		SongIds:   songIds,
	}
}

type Stream struct {
	Artist      string `json:"artist"`
	ProjecTitle string `json:"project_title"`
	TrackTitle  string `json:"track_title"`
	URL         string `json:"url"`
}

func NewStream(artist, projectTitle, trackTitle, url string) *Stream {
	return &Stream{
		Artist:      artist,
		ProjecTitle: projectTitle,
		TrackTitle:  trackTitle,
		URL:         url,
	}
}
