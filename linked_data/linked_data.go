package linked_data

import (
	"github.com/zbo14/envoke/bigchain"
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/spec"
)

func ValidateModelId(modelId string) (Data, error) {
	tx, err := bigchain.GetTx(modelId)
	if err != nil {
		return nil, err
	}
	model := bigchain.GetTxData(tx)
	if err = ValidateModel(model); err != nil {
		return nil, err
	}
	return model, nil
}

func ValidateMusicId(musicId string) (Data, error) {
	tx, err := bigchain.GetTx(musicId)
	if err != nil {
		return nil, err
	}
	music := bigchain.GetTxData(tx)
	if err = ValidateMusic(music); err != nil {
		return nil, err
	}
	return music, nil
}

func ValidateSignatureId(signatureId string) (Data, error) {
	tx, err := bigchain.GetTx(signatureId)
	if err != nil {
		return nil, err
	}
	signature := bigchain.GetTxData(tx)
	if err = ValidateSignature(signature); err != nil {
		return nil, err
	}
	return signature, nil
}

func ValidateRightId(rightId string) (Data, error) {
	tx, err := bigchain.GetTx(rightId)
	if err != nil {
		return nil, err
	}
	right := bigchain.GetTxData(tx)
	if err := ValidateRight(right); err != nil {
		return nil, err
	}
	return right, nil
}

func ValidateRightsId(rightsId string) (Data, error) {
	tx, err := bigchain.GetTx(rightsId)
	if err != nil {
		return nil, err
	}
	rights := bigchain.GetTxData(tx)
	if err := ValidateRights(rights); err != nil {
		return nil, err
	}
	return rights, nil
}

func ValidateModel(model Data) error {
	_type := spec.GetType(model)
	switch _type {
	case spec.ALBUM:
		return ValidateAlbum(model)
	case spec.TRACK:
		return ValidateTrack(model)
	case spec.SIGNATURE:
		return ValidateSignature(model)
	case spec.RIGHT:
		return ValidateRight(model)
	case spec.RIGHTS:
		return ValidateRights(model)
	}
	return ErrorAppend(ErrInvalidType, _type)
}

func ValidateMusic(music Data) error {
	_type := spec.GetType(music)
	if _type == spec.ALBUM {
		return ValidateAlbum(music)
	}
	if _type == spec.TRACK {
		return ValidateTrack(music)
	}
	return ErrorAppend(ErrInvalidType, _type)
}

func ValidateAlbum(album Data) error {
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
	labelId := spec.GetMusicLabel(album)
	if !EmptyStr(labelId) {
		tx, err = bigchain.GetTx(labelId)
		if err != nil {
			return err
		}
		label := bigchain.GetTxData(tx)
		if !spec.ValidLabel(label) {
			return ErrorAppend(ErrInvalidModel, spec.LABEL)
		}
	}
	publisherId := spec.GetMusicPublisher(album)
	if !EmptyStr(publisherId) {
		tx, err = bigchain.GetTx(publisherId)
		if err != nil {
			return err
		}
		publisher := bigchain.GetTxData(tx)
		if !spec.ValidPublisher(publisher) {
			return ErrorAppend(ErrInvalidModel, spec.PUBLISHER)
		}
	}
	return nil
}

func ValidateTrack(track Data) error {
	if !spec.ValidTrack(track) {
		return ErrorAppend(ErrInvalidModel, spec.TRACK)
	}
	labelId := spec.GetMusicLabel(track)
	if !EmptyStr(labelId) {
		tx, err := bigchain.GetTx(labelId)
		if err != nil {
			return err
		}
		label := bigchain.GetTxData(tx)
		if !spec.ValidLabel(label) {
			return ErrorAppend(ErrInvalidModel, spec.LABEL)
		}
	}
	publisherId := spec.GetMusicPublisher(track)
	if !EmptyStr(publisherId) {
		tx, err := bigchain.GetTx(publisherId)
		if err != nil {
			return err
		}
		publisher := bigchain.GetTxData(tx)
		if !spec.ValidPublisher(publisher) {
			return ErrorAppend(ErrInvalidModel, spec.PUBLISHER)
		}
	}
	albumId := spec.GetTrackAlbum(track)
	if !EmptyStr(albumId) {
		tx, err := bigchain.GetTx(albumId)
		if err != nil {
			return err
		}
		album := bigchain.GetTxData(tx)
		return ValidateAlbum(album)
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
	return nil
}

func ValidateSignature(signature Data) error {
	if !spec.ValidSignature(signature) {
		return ErrorAppend(ErrInvalidModel, spec.SIGNATURE)
	}
	modelId := spec.GetSignatureModel(signature)
	tx, err := bigchain.GetTx(modelId)
	if err != nil {
		return err
	}
	model := bigchain.GetTxData(tx)
	if err := ValidateModel(model); err != nil {
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

func ValidateRight(right Data) error {
	if !spec.ValidRight(right) {
		return ErrorAppend(ErrInvalidModel, spec.RIGHT)
	}
	issuerId, issuerType := spec.GetRightIssuer(right)
	musicId := spec.GetRightMusic(right)
	tx, err := bigchain.GetTx(musicId)
	if err != nil {
		return err
	}
	music := bigchain.GetTxData(tx)
	if err = ValidateMusic(music); err != nil {
		return err
	}
	agentId := spec.GetMusicAgent(music, issuerType)
	if EmptyStr(agentId) {
		albumId := spec.GetTrackAlbum(music)
		tx, err := bigchain.GetTx(albumId)
		if err != nil {
			return err
		}
		album := bigchain.GetTxData(tx)
		if err := ValidateAlbum(album); err != nil {
			return err
		}
		agentId = spec.GetMusicAgent(album, issuerType)
	}
	if agentId != issuerId {
		return ErrorAppend(ErrInvalidId, "issuer")
	}
	tx, err = bigchain.GetTx(agentId)
	if err != nil {
		return err
	}
	agent := bigchain.GetTxData(tx)
	if !spec.ValidAgentWithType(agent, issuerType) {
		return ErrorAppend(ErrInvalidModel, issuerType)
	}
	pub := spec.GetAgentPublicKey(agent)
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
	if err = ValidateSignature(signature); err != nil {
		return err
	}
	sig := spec.GetSignatureValue(signature)
	if !pub.Verify(MustMarshalJSON(music), sig) {
		return ErrInvalidSignature
	}
	return nil
}

// Criteria
// - Rights must point to same music
// - Total percentage shares must equal 100
// - An agent cannot receive multiple rights

func ValidateRights(rights Data) error {
	if !spec.ValidRights(rights) {
		return ErrorAppend(ErrInvalidModel, spec.RIGHTS)
	}
	musicId := spec.GetRightsMusic(rights)
	percentageShares := 0
	recipientIds := make(map[string]struct{})
	rightId := spec.GetMyRight(rights)
	tx, err := bigchain.GetTx(rightId)
	if err != nil {
		return err
	}
	right := bigchain.GetTxData(tx)
	if err := ValidateRight(right); err != nil {
		return err
	}
	if musicId != spec.GetRightMusic(right) {
		return ErrorAppend(ErrCriteriaNotMet, "right does not point to same music")
	}
	percentageShares += spec.GetRightPercentageShares(right)
	recipientIds[spec.GetRightRecipient(right)] = struct{}{}
	rightIds := spec.GetOtherRights(rights)
	for _, rightId = range rightIds {
		tx, err := bigchain.GetTx(rightId)
		if err != nil {
			return err
		}
		right := bigchain.GetTxData(tx)
		if err := ValidateRight(right); err != nil {
			return err
		}
		if musicId != spec.GetRightMusic(right) {
			return ErrorAppend(ErrCriteriaNotMet, "right does not point to same music")
		}
		percentageShares += spec.GetRightPercentageShares(right)
		if percentageShares > 100 {
			return ErrorAppend(ErrCriteriaNotMet, "percentage shares cannot exceed 100")
		}
		recipientId := spec.GetRightRecipient(right)
		if _, ok := recipientIds[recipientId]; ok {
			return ErrorAppend(ErrCriteriaNotMet, "agent cannot receive multiple rights for music")
		}
		recipientIds[recipientId] = struct{}{}
	}
	// if percentageShares != 100 {
	// 	return ErrorAppend(ErrCriteriaNotMet, "total percentage shares do not equal 100")
	// }
	signature := spec.GetRightsSignature(rights)
	return ValidateSignature(signature)
}
