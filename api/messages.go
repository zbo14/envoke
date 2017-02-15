package api

import (
	. "github.com/zbo14/envoke/common"
)

type RegisterMessage struct {
	AgentId string `json:"agent_id"`
	PrivKey string `json:"private_key"`
	PubKey  string `json:"public_key"`
}

func NewRegisterMessage(agentId, priv, pub string) *RegisterMessage {
	return &RegisterMessage{
		AgentId: agentId,
		PrivKey: priv,
		PubKey:  pub,
	}
}

func NewLoginMessage(agentName string) string {
	return Sprintf("Welcome %s!", agentName)
}

type AlbumMessage struct {
	AlbumId  string   `json:"album_id"`
	TrackIds []string `json:"track_ids"`
}

func NewAlbumMessage(albumId string, trackIds []string) *AlbumMessage {
	return &AlbumMessage{
		AlbumId:  albumId,
		TrackIds: trackIds,
	}
}

type TrackMessage struct {
	TrackId string `json:"track_id"`
}

func NewTrackMessage(trackId string) *TrackMessage {
	return &TrackMessage{
		TrackId: trackId,
	}
}

type RightMessage struct {
	RightId string `json:"right_id"`
}

func NewRightMessage(rightId string) *RightMessage {
	return &RightMessage{
		RightId: rightId,
	}
}

type RightsMessage struct {
	RightsId string `json:"rights_id"`
}

func NewRightsMessage(rightsId string) *RightsMessage {
	return &RightsMessage{
		RightsId: rightsId,
	}
}

type VerifyMessage struct {
	Data  Data
	Log   string `json:"log"`
	Valid bool   `json:"valid"`
}

func NewVerifyMessage(data Data, log string, valid bool) *VerifyMessage {
	return &VerifyMessage{
		Data:  data,
		Log:   log,
		Valid: valid,
	}
}
