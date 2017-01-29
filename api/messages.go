package api

import "github.com/zballs/go_resonate/crypto/ed25519"

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

type ProjectInfo struct {
	ProjectId string   `json:"project_id"`
	SongIds   []string `json:"song_ids"`
}

func NewProjectInfo(projectId string, songIds []string) *ProjectInfo {
	return &ProjectInfo{
		ProjectId: projectId,
		SongIds:   songIds,
	}
}
