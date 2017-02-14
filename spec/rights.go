package spec

import (
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/crypto"
	"time"
)

const (
	RIGHT      = "right"
	RIGHT_SIZE = 10
)

// Right

// What should context, usage be?..

func NewRight(context []string, issuerId, issuerType, musicId, percentageShares, recipientId string, sig crypto.Signature, usage []string, validFrom, validTo time.Time) Data {
	return Data{
		"context":           context,
		"instance":          NewInstance(RIGHT),
		"issuer_id":         issuerId,
		"issuer_type":       issuerType,
		"percentage_shares": percentageShares,
		"recipient_id":      recipientId,
		"signature":         NewSignature(musicId, issuerId, sig),
		"usage":             usage,
		"valid_from":        validFrom,
		"valid_to":          validTo,
	}
}

func IsRight(right Data) bool {
	return HasType(right, RIGHT)
}

func GetRightContext(right Data) []string {
	return right.GetStrSlice("context")
}

func GetRightIssuer(right Data) (string, string) {
	return right.GetStr("issuer_id"), right.GetStr("issuer_type")
}

func GetRightMusic(right Data) string {
	signature := GetRightSignature(right)
	if signature == nil {
		return ""
	}
	return GetSignatureModel(signature)
}

func GetRightRecipient(right Data) string {
	return right.GetStr("recipient_id")
}

func GetRightPercentageShares(right Data) int {
	return right.GetStrInt("percentage_shares")
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
	instance := GetInstance(right)
	if !ValidInstance(instance) {
		return false
	}
	if !HasType(right, RIGHT) {
		return false
	}
	issuerId, issuerType := GetRightIssuer(right)
	switch issuerType {
	case ARTIST, LABEL, ORGANIZATION, PUBLISHER:
		if !MatchString(ID_REGEX, issuerId) {
			return false
		}
	default:
		return false
	}
	// TODO: validate context
	musicId := GetRightMusic(right)
	if !MatchString(ID_REGEX, musicId) {
		return false
	}
	percentageShares := GetRightPercentageShares(right)
	if percentageShares <= 0 || percentageShares > 100 {
		return false
	}
	receipientId := GetRightRecipient(right)
	if !MatchString(ID_REGEX, receipientId) {
		return false
	}
	if issuerId == receipientId {
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
