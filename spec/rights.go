package spec

import (
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/crypto"
	"time"
)

const (
	RIGHT      = "right"
	RIGHT_SIZE = 8
)

// Right

// What should context, usage be?..

func NewRight(context []string, issuerId, musicId, recipientId string, sig crypto.Signature, usage []string, validFrom, validTo time.Time) Data {
	return Data{
		"context":      context,
		"entity":       NewEntity(RIGHT),
		"music_id":     musicId,
		"recipient_id": recipientId,
		"signature":    NewSignature(musicId, issuerId, sig),
		"usage":        usage,
		"valid_from":   validFrom,
		"valid_to":     validTo,
	}
}

func IsRight(right Data) bool {
	return HasType(right, RIGHT)
}

func GetRightContext(right Data) []string {
	return right.GetStrSlice("context")
}

func GetRightMusic(right Data) string {
	return right.GetStr("music_id")
}

func GetRightRecipient(right Data) string {
	return right.GetStr("recipient_id")
}

func GetRightSignature(right Data) Data {
	return right.GetData("signature")
}

func GetRightUsage(right Data) []string {
	return right.GetStrSlice("usage")
}

func GetRightValidFrom(right Data) time.Time {
	return right.GetTime("valid_from")
}

func GetRightValidTo(right Data) time.Time {
	return right.GetTime("valid_to")
}

func ValidRight(right Data) bool {
	entity := GetEntity(right)
	if !ValidEntity(entity) {
		return false
	}
	if !HasType(right, RIGHT) {
		return false
	}
	// TODO: validate context
	musicId := GetRightMusic(right)
	if !MatchString(ID_REGEX, musicId) {
		return false
	}
	receipientId := GetRightRecipient(right)
	if !MatchString(ID_REGEX, receipientId) {
		return false
	}
	signature := GetRightSignature(right)
	if !ValidSignature(signature) {
		return false
	}
	// TODO: validate usage
	validFrom := GetRightValidFrom(right)
	validTo := GetRightValidTo(right)
	if validFrom.After(validTo) {
		return false
	}
	return len(right) == RIGHT_SIZE
}
