package linked_data

import (
	"github.com/zbo14/envoke/bigchain"
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/spec/core"
)

func ValidLinkedMusicId(musicId string) (core.Data, error) {
	tx, err := bigchain.GetTx(musicId)
	if err != nil {
		return nil, err
	}
	music := bigchain.GetTxData(tx)
	if err = ValidLinkedMusic(music); err != nil {
		return nil, err
	}
	return music, nil
}

func ValidLinkedSignatureId(signatureId string) (core.Data, error) {
	tx, err := bigchain.GetTx(signatureId)
	if err != nil {
		return nil, err
	}
	signature := bigchain.GetTxData(tx)
	if err = ValidLinkedSignature(signature); err != nil {
		return nil, err
	}
	return signature, nil
}

func ValidLinkedMusic(music core.Data) error {
	if _type := core.GetType(music); _type == core.ALBUM {
		return ValidLinkedAlbum(music)
	} else if _type == core.TRACK {
		return ValidLinkedTrack(music)
	}
	return ErrInvalidType
}

func ValidLinkedAlbum(album core.Data) error {
	if !core.ValidAlbum(album) {
		return ErrorAppend(ErrInvalidModel, core.ALBUM)
	}
	artistId := core.GetMusicArtist(album)
	tx, err := bigchain.GetTx(artistId)
	if err != nil {
		return err
	}
	artist := bigchain.GetTxData(tx)
	if !core.ValidArtist(artist) {
		return ErrorAppend(ErrInvalidModel, core.ARTIST)
	}
	publisherId := core.GetMusicPublisher(album)
	tx, err = bigchain.GetTx(publisherId)
	if err != nil {
		return err
	}
	publisher := bigchain.GetTxData(tx)
	if !core.ValidPublisher(publisher) {
		return ErrorAppend(ErrInvalidModel, core.PUBLISHER)
	}
	return nil
}

func ValidLinkedTrack(track core.Data) error {
	if !core.ValidTrack(track) {
		return ErrorAppend(ErrInvalidModel, core.TRACK)
	}
	artistId := core.GetMusicArtist(track)
	tx, err := bigchain.GetTx(artistId)
	if err != nil {
		return err
	}
	artist := bigchain.GetTxData(tx)
	if !core.ValidArtist(artist) {
		return ErrorAppend(ErrInvalidModel, core.ARTIST)
	}
	publisherId := core.GetMusicPublisher(track)
	if publisherId != "" {
		tx, err = bigchain.GetTx(publisherId)
		if err != nil {
			return err
		}
		publisher := bigchain.GetTxData(tx)
		if !core.ValidPublisher(publisher) {
			return ErrorAppend(ErrInvalidModel, core.PUBLISHER)
		}
		return nil
	}
	albumId := core.GetTrackAlbum(track)
	tx, err = bigchain.GetTx(albumId)
	if err != nil {
		return err
	}
	album := bigchain.GetTxData(tx)
	return ValidLinkedAlbum(album)
}

func ValidLinkedSignature(signature core.Data) error {
	if !core.ValidSignature(signature) {
		return ErrorAppend(ErrInvalidModel, core.SIGNATURE)
	}
	musicId := core.GetSignatureMusic(signature)
	tx, err := bigchain.GetTx(musicId)
	if err != nil {
		return err
	}
	music := bigchain.GetTxData(tx)
	if err := ValidLinkedMusic(music); err != nil {
		return err
	}
	// Should have traversed graph to reach music publisher/signature agent
	// Check if music publisher_id == signature agent_id
	// signature agents are just publishers for now
	publisherId := core.GetMusicPublisher(music)
	if publisherId != core.GetSignatureAgent(signature) {
		return ErrInvalidId
	}
	tx, err = bigchain.GetTx(publisherId)
	if err != nil {
		return err
	}
	publisher := bigchain.GetTxData(tx)
	pub := core.GetAgentPublicKey(publisher)
	sig := core.GetSignatureValue(signature)
	p := MustMarshalJSON(music)
	if !pub.Verify(p, sig) {
		return ErrInvalidSignature
	}
	return nil
}
