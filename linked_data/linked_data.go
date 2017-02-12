package linked_data

import (
	"github.com/zbo14/envoke/bigchain"
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/spec"
)

func ValidateLdModelId(modelId string) (Data, error) {
	tx, err := bigchain.GetTx(modelId)
	if err != nil {
		return nil, err
	}
	model := bigchain.GetTxData(tx)
	if err = ValidateLdModel(model); err != nil {
		return nil, err
	}
	return model, nil
}

func ValidateLdMusicId(musicId string) (Data, error) {
	tx, err := bigchain.GetTx(musicId)
	if err != nil {
		return nil, err
	}
	music := bigchain.GetTxData(tx)
	if err = ValidateLdMusic(music); err != nil {
		return nil, err
	}
	return music, nil
}

func ValidateLdSignatureId(signatureId string) (Data, error) {
	tx, err := bigchain.GetTx(signatureId)
	if err != nil {
		return nil, err
	}
	signature := bigchain.GetTxData(tx)
	if err = ValidateLdSignature(signature); err != nil {
		return nil, err
	}
	return signature, nil
}

func ValidateLdModel(model Data) error {
	_type := spec.GetType(model)
	switch _type {
	case spec.ALBUM:
		return ValidateLdAlbum(model)
	case spec.TRACK:
		return ValidateLdTrack(model)
	case spec.SIGNATURE:
		return ValidateLdSignature(model)
	case spec.RIGHT:
		return ValidateLdRight(model)
	}
	return ErrInvalidType
}

func ValidateLdMusic(music Data) error {
	if _type := spec.GetType(music); _type == spec.ALBUM {
		return ValidateLdAlbum(music)
	} else if _type == spec.TRACK {
		return ValidateLdTrack(music)
	}
	return ErrInvalidType
}

func ValidateLdAlbum(album Data) error {
	if !spec.ValidAlbum(album) {
		return ErrorAppend(ErrInvalidModel, spec.ALBUM)
	}
	artistId := spec.GetMusicArtist(album)
	tx, err := bigchain.GetTx(artistId)
	if err != nil {
		return err
	}
	artist := bigchain.GetTxData(tx)
	if !spec.ValidArtist(artist) {
		return ErrorAppend(ErrInvalidModel, spec.ARTIST)
	}
	publisherId := spec.GetMusicPublisher(album)
	tx, err = bigchain.GetTx(publisherId)
	if err != nil {
		return err
	}
	publisher := bigchain.GetTxData(tx)
	if !spec.ValidPublisher(publisher) {
		return ErrorAppend(ErrInvalidModel, spec.PUBLISHER)
	}
	return nil
}

func ValidateLdTrack(track Data) error {
	if !spec.ValidTrack(track) {
		return ErrorAppend(ErrInvalidModel, spec.TRACK)
	}
	artistId := spec.GetMusicArtist(track)
	tx, err := bigchain.GetTx(artistId)
	if err != nil {
		return err
	}
	artist := bigchain.GetTxData(tx)
	if !spec.ValidArtist(artist) {
		return ErrorAppend(ErrInvalidModel, spec.ARTIST)
	}
	publisherId := spec.GetMusicPublisher(track)
	if publisherId != "" {
		tx, err = bigchain.GetTx(publisherId)
		if err != nil {
			return err
		}
		publisher := bigchain.GetTxData(tx)
		if !spec.ValidPublisher(publisher) {
			return ErrorAppend(ErrInvalidModel, spec.PUBLISHER)
		}
		return nil
	}
	albumId := spec.GetTrackAlbum(track)
	tx, err = bigchain.GetTx(albumId)
	if err != nil {
		return err
	}
	album := bigchain.GetTxData(tx)
	return ValidateLdAlbum(album)
}

func ValidateLdSignature(signature Data) error {
	if !spec.ValidSignature(signature) {
		return ErrorAppend(ErrInvalidModel, spec.SIGNATURE)
	}
	modelId := spec.GetSignatureModel(signature)
	tx, err := bigchain.GetTx(modelId)
	if err != nil {
		return err
	}
	model := bigchain.GetTxData(tx)
	if err := ValidateLdModel(model); err != nil {
		return err
	}
	signerId := spec.GetSignatureSigner(signature)
	tx, err = bigchain.GetTx(signerId)
	if err != nil {
		return err
	}
	signer := bigchain.GetTxData(tx)
	if !spec.ValidAgent(signer) {
		return err
	}
	pub := spec.GetAgentPublicKey(signer)
	sig := spec.GetSignatureValue(signature)
	if !pub.Verify(MustMarshalJSON(model), sig) {
		return ErrInvalidSignature
	}
	return nil
}

func ValidateLdRight(right Data) error {
	if !spec.ValidRight(right) {
		return ErrorAppend(ErrInvalidModel, spec.RIGHT)
	}
	musicId := spec.GetRightMusic(right)
	tx, err := bigchain.GetTx(musicId)
	if err != nil {
		return err
	}
	music := bigchain.GetTxData(tx)
	if err = ValidateLdMusic(music); err != nil {
		return err
	}
	artistId := spec.GetMusicArtist(music)
	tx, err = bigchain.GetTx(artistId)
	if err != nil {
		return err
	}
	artist := bigchain.GetTxData(tx)
	if !spec.ValidArtist(artist) {
		return ErrorAppend(ErrInvalidModel, spec.ARTIST)
	}
	pub := spec.GetAgentPublicKey(artist)
	recipientId := spec.GetRightRecipient(right)
	tx, err = bigchain.GetTx(recipientId)
	if err != nil {
		return err
	}
	recipient := bigchain.GetTxData(tx)
	if !spec.ValidAgent(recipient) {
		return ErrorAppend(ErrInvalidModel, spec.GetType(recipient))
	}
	signature := spec.GetRightSignature(right)
	if err = ValidateLdSignature(signature); err != nil {
		return err
	}
	sig := spec.GetSignatureValue(signature)
	if !pub.Verify(MustMarshalJSON(music), sig) {
		return ErrInvalidSignature
	}
	return nil
}

// Should have traversed graph to reach music publisher/signer
// Check if publisher_id equals signer_id
// publisherId := spec.GetMusicPublisher(music)
