package api

import (
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/spec/core"
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

type SignMessage struct {
	SignatureId string `json:"signature_id"`
}

func NewSignMessage(signatureId string) *SignMessage {
	return &SignMessage{
		SignatureId: signatureId,
	}
}

type VerifyMessage struct {
	Log       string `json:"log"`
	Signature core.Data
	Valid     bool `json:"valid"`
}

func NewVerifyMessage(log string, signature core.Data, valid bool) *VerifyMessage {
	return &VerifyMessage{
		Log:       log,
		Signature: signature,
		Valid:     valid,
	}
}
