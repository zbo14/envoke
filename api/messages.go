package api

import "github.com/zballs/go_resonate/crypto"

type UserInfo struct {
	Privkey *crypto.PrivateKey `json:"private_key"`
	Pubkey  *crypto.PublicKey  `json:"public_key"`
	UserId  string             `json:"user_id"`
}

func NewUserInfo(userId string, priv *crypto.PrivateKey, pub *crypto.PublicKey) *UserInfo {
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
