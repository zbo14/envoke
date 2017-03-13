package linked_data

import (
	jsonschema "github.com/xeipuuv/gojsonschema"
	"github.com/zbo14/balloon"
	"github.com/zbo14/envoke/bigchain"
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/crypto"
	"github.com/zbo14/envoke/spec"
)

const SPEC = "http://localhost:8888/spec"

func QueryAndValidateModel(id string) (Data, error) {
	tx, err := bigchain.GetTx(id)
	if err != nil {
		return nil, err
	}
	// model := bigchain.GetTxData(tx)
	// if err = ValidateModel(model); err != nil {
	//	return nil, err
	// }
	return tx, nil
}

func ValidateModel(model Data) error {
	source := spec.GetType(model)
	schemaLoader, err := GetSchemaLoader(source)
	if err != nil {
		return err
	}
	goLoader := jsonschema.NewGoLoader(model)
	result, err := jsonschema.Validate(schemaLoader, goLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return Error("Validation failed")
	}
	return nil
}

func GetSchemaLoader(source string) (jsonschema.JSONLoader, error) {
	response, err := HttpGet(SPEC)
	if err != nil {
		return nil, err
	}
	p, err := ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	bytes := Submatch(Sprintf(`<script id="%s".*?>([\s\S]*?)</script>`, source), p)[1]
	return jsonschema.NewBytesLoader(bytes), nil
}

var SALT = balloon.GenerateSalt()

func DefaultBalloonHash(challenge string) ([]byte, error) {
	p, err := Base64UrlDecode(challenge)
	if err != nil {
		return nil, err
	}
	// TODO: adjust params
	return balloon.BalloonHash(p, SALT, 256, 32, 2), nil
}

func ValidateComposition(compositionId string) (Data, error) {
	tx, err := QueryAndValidateModel(compositionId)
	if err != nil {
		return nil, err
	}
	composition := bigchain.GetTxData(tx)
	senderPub := bigchain.DefaultGetTxSender(tx)
	composerId := spec.GetComposerId(composition)
	tx, err = QueryAndValidateModel(composerId)
	if err != nil {
		return nil, err
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, ErrorAppend(ErrInvalidKey, senderPub.String())
	}
	return composition, nil
}

func QueryCompositionField(compositionId, field string) (interface{}, error) {
	composition, err := ValidateComposition(compositionId)
	if err != nil {
		return nil, err
	}
	switch field {
	case "composer":
		return GetComposer(composition)
		//..
	}
	return nil, ErrorAppend(ErrInvalidField, field)
}

func GetComposer(data Data) (Data, error) {
	composerId := spec.GetComposerId(data)
	tx, err := bigchain.GetTx(composerId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func ProveComposer(challenge, compositionId string, priv crypto.PrivateKey) (crypto.Signature, error) {
	composition, err := ValidateComposition(compositionId)
	if err != nil {
		return nil, err
	}
	composerId := spec.GetComposerId(composition)
	tx, err := bigchain.GetTx(composerId)
	if err != nil {
		return nil, err
	}
	senderPub := bigchain.DefaultGetTxSender(tx)
	if !senderPub.Equals(priv.Public()) {
		return nil, ErrorAppend(ErrInvalidKey, priv.Public().String())
	}
	hash, err := DefaultBalloonHash(challenge)
	if err != nil {
		return nil, err
	}
	return priv.Sign(hash), nil
}

func VerifyComposer(challenge, compositionId string, sig crypto.Signature) error {
	composition, err := ValidateComposition(compositionId)
	if err != nil {
		return err
	}
	composerId := spec.GetComposerId(composition)
	tx, err := bigchain.GetTx(composerId)
	if err != nil {
		return err
	}
	hash, err := DefaultBalloonHash(challenge)
	if err != nil {
		return err
	}
	senderPub := bigchain.DefaultGetTxSender(tx)
	if !senderPub.Verify(hash, sig) {
		return ErrorAppend(ErrInvalidSignature, sig.String())
	}
	return nil
}

func ValidateRight(rightId string) (Data, crypto.PublicKey, crypto.PublicKey, error) {
	tx, err := QueryAndValidateModel(rightId)
	if err != nil {
		return nil, nil, nil, err
	}
	right := bigchain.GetTxData(tx)
	recipientId := spec.GetRecipientId(right)
	recipientPub := bigchain.DefaultGetTxRecipient(tx)
	recipientShares := bigchain.GetTxShares(tx)
	senderId := spec.GetSenderId(right)
	senderPub := bigchain.DefaultGetTxSender(tx)
	tx, err = QueryAndValidateModel(recipientId)
	if err != nil {
		return nil, nil, nil, err
	}
	if !recipientPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, nil, nil, ErrorAppend(ErrInvalidKey, recipientPub.String())
	}
	tx, err = QueryAndValidateModel(senderId)
	if err != nil {
		return nil, nil, nil, err
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, nil, nil, ErrorAppend(ErrInvalidKey, senderPub.String())
	}
	right.Set("recipientShares", recipientShares)
	return right, recipientPub, senderPub, nil
}

func ProveCompositionRightHolder(challenge, compositionRightId string, priv crypto.PrivateKey, publicationId string) (crypto.Signature, error) {
	_, _, compositionRights, err := ValidatePublication(publicationId)
	if err != nil {
		return nil, err
	}
	for _, compositionRight := range compositionRights {
		if compositionRightId == spec.GetId(compositionRight) {
			recipientId := spec.GetRecipientId(compositionRight)
			tx, err := bigchain.GetTx(recipientId)
			if err != nil {
				return nil, err
			}
			recipientPub := bigchain.DefaultGetTxSender(tx)
			if pub := priv.Public(); !recipientPub.Equals(pub) {
				return nil, ErrorAppend(ErrInvalidKey, pub.String())
			}
			hash, err := DefaultBalloonHash(challenge)
			if err != nil {
				return nil, err
			}
			return priv.Sign(hash), nil
		}
	}
	return nil, ErrorAppend(ErrCriteriaNotMet, "publication does not link to composition right")
}

func VerifyCompositionRightHolder(challenge, compositionRightId, publicationId string, sig crypto.Signature) error {
	_, _, compositionRights, err := ValidatePublication(publicationId)
	if err != nil {
		return err
	}
	for _, compositionRight := range compositionRights {
		if compositionRightId == spec.GetId(compositionRight) {
			recipientId := spec.GetRecipientId(compositionRight)
			tx, err := bigchain.GetTx(recipientId)
			if err != nil {
				return err
			}
			hash, err := DefaultBalloonHash(challenge)
			if err != nil {
				return err
			}
			recipientPub := bigchain.DefaultGetTxSender(tx)
			if !recipientPub.Verify(hash, sig) {
				return ErrorAppend(ErrInvalidSignature, sig.String())
			}
			return nil
		}
	}
	return ErrorAppend(ErrCriteriaNotMet, "publication does not link to composition right")
}

func GetRecipient(data Data) (Data, error) {
	recipientId := spec.GetRecipientId(data)
	tx, err := bigchain.GetTx(recipientId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetSender(data Data) (Data, error) {
	senderId := spec.GetSenderId(data)
	tx, err := bigchain.GetTx(senderId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func QueryPublicationField(field, publicationId string) (interface{}, error) {
	publication, compositions, compositionRights, err := ValidatePublication(publicationId)
	if err != nil {
		return nil, err
	}
	switch field {
	case "compositions":
		return compositions, nil
	case "composition_rights":
		return compositionRights, nil
	case "publisher":
		return GetPublisher(publication)
	}
	return nil, ErrorAppend(ErrInvalidField, field)
}

func ValidatePublication(publicationId string) (Data, []Data, []Data, error) {
	tx, err := QueryAndValidateModel(publicationId)
	if err != nil {
		return nil, nil, nil, err
	}
	publication := bigchain.GetTxData(tx)
	senderPub := bigchain.DefaultGetTxSender(tx)
	var composerId string
	compositionIds := spec.GetCompositionIds(publication)
	compositions := make([]Data, len(compositionIds))
	for i, compositionId := range compositionIds {
		composition, err := ValidateComposition(compositionId)
		if err != nil {
			return nil, nil, nil, err
		}
		if i == 0 {
			composerId = spec.GetComposerId(composition)
			// TODO: check composerId
		} else if composerId != spec.GetComposerId(composition) {
			return nil, nil, nil, ErrorAppend(ErrCriteriaNotMet, "publication cannot link to compositions with different composers")
		}
		composition.Set("@id", compositionId)
		compositions[i] = composition
	}
	compositionRightIds := spec.GetCompositionRightIds(publication)
	compositionRights := make([]Data, len(compositionRightIds))
	publisherId := spec.GetPublisherId(publication)
	recipientIds := make(map[string]struct{})
	rightHolder := false
	totalShares := 0
	for i, compositionRightId := range compositionRightIds {
		compositionRight, recipientPub, _, err := ValidateRight(compositionRightId)
		if err != nil {
			return nil, nil, nil, err
		}
		if composerId != spec.GetSenderId(compositionRight) {
			return nil, nil, nil, ErrorAppend(ErrCriteriaNotMet, "composer must be right sender")
		}
		recipientId := spec.GetRecipientId(compositionRight)
		if _, ok := recipientIds[recipientId]; ok {
			return nil, nil, nil, ErrorAppend(ErrCriteriaNotMet, "recipient cannot hold multiple composition rights")
		}
		if !rightHolder && publisherId == recipientId {
			if !senderPub.Equals(recipientPub) {
				return nil, nil, nil, ErrorAppend(ErrCriteriaNotMet, "publisher is not publication sender")
			}
			rightHolder = true
		}
		recipientIds[recipientId] = struct{}{}
		shares := spec.GetRecipientShares(compositionRight)
		if shares <= 0 {
			return nil, nil, nil, ErrorAppend(ErrCriteriaNotMet, "percentage shares must be greater than 0")
		}
		if totalShares += shares; totalShares > 100 {
			return nil, nil, nil, ErrorAppend(ErrCriteriaNotMet, "total percentage shares cannot exceed 100")
		}
		compositionRight.Set("@id", compositionRightId)
		compositionRights[i] = compositionRight
	}
	if !rightHolder {
		return nil, nil, nil, ErrorAppend(ErrCriteriaNotMet, "publisher is not composition right holder")
	}
	if totalShares != 100 {
		return nil, nil, nil, ErrorAppend(ErrCriteriaNotMet, "total percentage shares do not equal 100")
	}
	return publication, compositions, compositionRights, nil
}

func ProvePublisher(challenge string, priv crypto.PrivateKey, publicationId string) (crypto.Signature, error) {
	publication, _, _, err := ValidatePublication(publicationId)
	if err != nil {
		return nil, err
	}
	publisherId := spec.GetPublisherId(publication)
	tx, err := bigchain.GetTx(publisherId)
	if err != nil {
		return nil, err
	}
	senderPub := bigchain.DefaultGetTxSender(tx)
	if pub := priv.Public(); !senderPub.Equals(pub) {
		return nil, ErrorAppend(ErrInvalidKey, pub.String())
	}
	hash, err := DefaultBalloonHash(challenge)
	if err != nil {
		return nil, err
	}
	return priv.Sign(hash), nil
}

func VerifyPublisher(challenge, publicationId string, sig crypto.Signature) error {
	publication, _, _, err := ValidatePublication(publicationId)
	if err != nil {
		return err
	}
	publisherId := spec.GetPublisherId(publication)
	tx, err := bigchain.GetTx(publisherId)
	if err != nil {
		return err
	}
	senderPub := bigchain.DefaultGetTxSender(tx)
	if !senderPub.Verify(MustMarshalJSON(publication), sig) {
		return ErrorAppend(ErrInvalidSignature, sig.String())
	}
	return nil
}

func GetPublisher(data Data) (Data, error) {
	publisherId := spec.GetPublisherId(data)
	tx, err := bigchain.GetTx(publisherId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func ValidateCompositionRightTransfer(compositionRightTransferId string) (Data, error) {
	tx, err := QueryAndValidateModel(compositionRightTransferId)
	if err != nil {
		return nil, err
	}
	compositionRightTransfer := bigchain.GetTxData(tx)
	senderPub := bigchain.DefaultGetTxSender(tx)
	recipientId := spec.GetRecipientId(compositionRightTransfer)
	tx, err = QueryAndValidateModel(recipientId)
	if err != nil {
		return nil, err
	}
	recipientPub := bigchain.DefaultGetTxSender(tx)
	if senderPub.Equals(recipientPub) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "recipient and sender keys must be different")
	}
	senderId := spec.GetSenderId(compositionRightTransfer)
	tx, err = QueryAndValidateModel(senderId)
	if err != nil {
		return nil, err
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, ErrorAppend(ErrInvalidKey, senderPub.String())
	}
	publicationId := spec.GetPublicationId(compositionRightTransfer)
	_, _, compositionRights, err := ValidatePublication(publicationId)
	if err != nil {
		return nil, err
	}
	txId := spec.GetTxId(compositionRightTransfer)
	tx, err = bigchain.GetTx(txId)
	if err != nil {
		return nil, err
	}
	if bigchain.TRANSFER != bigchain.GetTxOperation(tx) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "expected TRANSFER tx")
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "sender is not signer of TRANSFER tx")
	}
	n := len(bigchain.GetTxOutputs(tx))
	if n != 1 && n != 2 {
		return nil, ErrorAppend(ErrInvalidSize, "tx outputs must have size 1 or 2")
	}
	if !recipientPub.Equals(bigchain.GetTxRecipient(tx, 0)) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "recipient does not hold primary output of TRANSFER tx")
	}
	recipientShares := bigchain.GetTxOutputAmount(tx, 0)
	if recipientShares <= 0 || recipientShares > 100 {
		return nil, ErrorAppend(ErrCriteriaNotMet, "recipient shares must be greater than 0 and less than/equal to 100")
	}
	compositionRightTransfer.Set("recipientShares", recipientShares)
	if n == 2 {
		if !senderPub.Equals(bigchain.GetTxRecipient(tx, 1)) {
			return nil, ErrorAppend(ErrCriteriaNotMet, "sender does not hold secondary output of TRANSFER tx")
		}
		senderShares := bigchain.GetTxOutputAmount(tx, 1)
		if senderShares < 0 || senderShares > 100 {
			return nil, ErrorAppend(ErrCriteriaNotMet, "sender shares cannot be less than 0 or greater than 100")
		}
		compositionRightTransfer.Set("senderShares", senderShares)
	}
	found := false
	compositionRightId := spec.GetCompositionRightId(compositionRightTransfer)
	if compositionRightId != bigchain.GetTxAssetId(tx) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "TRANSFER tx does not link to correct composition right")
	}
	for _, compositionRight := range compositionRights {
		if compositionRightId == spec.GetId(compositionRight) {
			found = true
			break
		}
	}
	if !found {
		return nil, ErrorAppend(ErrCriteriaNotMet, "publication does not link to underlying composition right")
	}
	return compositionRightTransfer, nil
}

func ProveCompositionRightTransferHolder(challenge, compositionRightTransferId, holderId string, priv crypto.PrivateKey, publicationId string) (crypto.Signature, error) {
	compositionRightTransfer, err := ValidateCompositionRightTransfer(compositionRightTransferId)
	if err != nil {
		return nil, err
	}
	compositionRightId := spec.GetCompositionRightId(compositionRightTransfer)
	if holderId == spec.GetRecipientId(compositionRightTransfer) {
		//..
	} else if holderId == spec.GetSenderId(compositionRightTransfer) {
		if spec.GetSenderShares(compositionRightTransfer) == 0 {
			return nil, ErrorAppend(ErrCriteriaNotMet, "sender does not have any shares")
		}
	} else {
		return nil, ErrorAppend(ErrCriteriaNotMet, "holder is not recipient or sender of transfer")
	}
	_, _, compositionRights, err := ValidatePublication(publicationId)
	if err != nil {
		return nil, err
	}
	for _, compositionRight := range compositionRights {
		if compositionRightId == spec.GetId(compositionRight) {
			tx, err := bigchain.GetTx(holderId)
			if err != nil {
				return nil, err
			}
			holderPub := bigchain.DefaultGetTxSender(tx)
			if pub := priv.Public(); !holderPub.Equals(pub) {
				return nil, ErrorAppend(ErrInvalidKey, pub.String())
			}
			hash, err := DefaultBalloonHash(challenge)
			if err != nil {
				return nil, err
			}
			return priv.Sign(hash), nil
		}
	}
	return nil, ErrorAppend(ErrCriteriaNotMet, "publication does not link to underlying composition right")
}

func VerifyCompositionRightTransferHolder(challenge, compositionRightTransferId, holderId, publicationId string, sig crypto.Signature) error {
	compositionRightTransfer, err := ValidateCompositionRightTransfer(compositionRightTransferId)
	if err != nil {
		return err
	}
	compositionRightId := spec.GetCompositionRightId(compositionRightTransfer)
	if holderId == spec.GetRecipientId(compositionRightTransfer) {
		//..
	} else if holderId == spec.GetSenderId(compositionRightTransfer) {
		if spec.GetSenderShares(compositionRightTransfer) == 0 {
			return ErrorAppend(ErrCriteriaNotMet, "sender does not have any shares")
		}
	} else {
		return ErrorAppend(ErrCriteriaNotMet, "holder is not recipient or sender of transfer")
	}
	_, _, compositionRights, err := ValidatePublication(publicationId)
	if err != nil {
		return err
	}
	for _, compositionRight := range compositionRights {
		if compositionRightId == spec.GetId(compositionRight) {
			tx, err := bigchain.GetTx(holderId)
			if err != nil {
				return err
			}
			hash, err := DefaultBalloonHash(challenge)
			if err != nil {
				return err
			}
			holderPub := bigchain.DefaultGetTxSender(tx)
			if !holderPub.Verify(hash, sig) {
				return ErrorAppend(ErrInvalidSignature, sig.String())
			}
			return nil
		}
	}
	return ErrorAppend(ErrCriteriaNotMet, "publication does not link to underlying composition right")
}

func GetCompositionRight(data Data) (Data, error) {
	compositionRightId := spec.GetCompositionRightId(data)
	tx, err := bigchain.GetTx(compositionRightId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetPublication(data Data) (Data, error) {
	publicationId := spec.GetPublicationId(data)
	tx, err := bigchain.GetTx(publicationId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func QueryMechanicalLicenseField(field, mechanicalLicenseId string) (interface{}, error) {
	mechanicalLicense, compositions, err := ValidateMechanicalLicense(mechanicalLicenseId)
	if err != nil {
		return nil, err
	}
	switch field {
	case "compositions":
		return compositions, nil
	case "recipient":
		return GetRecipient(mechanicalLicense)
	case "sender":
		return GetSender(mechanicalLicense)
	}
	return nil, ErrorAppend(ErrInvalidField, field)
}

func ValidateMechanicalLicense(mechanicalLicenseId string) (Data, []Data, error) {
	tx, err := QueryAndValidateModel(mechanicalLicenseId)
	if err != nil {
		return nil, nil, err
	}
	mechanicalLicense := bigchain.GetTxData(tx)
	senderPub := bigchain.DefaultGetTxSender(tx)
	senderId := spec.GetSenderId(mechanicalLicense)
	tx, err = bigchain.GetTx(senderId)
	if err != nil {
		return nil, nil, err
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, nil, ErrorAppend(ErrInvalidKey, senderPub.String())
	}
	var compositions []Data
	compositionIds := spec.GetCompositionIds(mechanicalLicense)
	seen := make(map[string]struct{})
	if n := len(compositionIds); n > 0 {
		compositions = make([]Data, n)
		for i, compositionId := range compositionIds {
			if _, ok := seen[compositionId]; ok {
				return nil, nil, ErrorAppend(ErrCriteriaNotMet, "cannot license composition multiple times")
			}
			composition, err := ValidateComposition(compositionId)
			if err != nil {
				return nil, nil, err
			}
			if senderId != spec.GetComposerId(composition) {
				return nil, nil, ErrorAppend(ErrCriteriaNotMet, "cannot license composition by another composer")
			}
			seen[compositionId] = struct{}{}
			composition.Set("@id", compositionId)
			compositions[i] = composition
		}
	}
	publicationId := spec.GetPublicationId(mechanicalLicense)
	if !EmptyStr(publicationId) {
		_, moreCompositions, compositionRights, err := ValidatePublication(publicationId)
		if err != nil {
			return nil, nil, err
		}
		compositionRightId := spec.GetCompositionRightId(mechanicalLicense)
		compositionRightTransferHolder := false
		if EmptyStr(compositionRightId) {
			compositionRightTransferId := spec.GetCompositionRightTransferId(mechanicalLicense)
			compositionRightTransfer, err := ValidateCompositionRightTransfer(compositionRightTransferId)
			if err != nil {
				return nil, nil, err
			}
			if publicationId != spec.GetPublicationId(compositionRightTransfer) {
				return nil, nil, ErrorAppend(ErrCriteriaNotMet, "compositionRightTransfer links to wrong publication")
			}
			if senderId == spec.GetRecipientId(compositionRightTransfer) {
				//..
			} else if senderId == spec.GetSenderId(compositionRightTransfer) {
				if spec.GetSenderShares(compositionRightTransfer) == 0 {
					return nil, nil, ErrorAppend(ErrCriteriaNotMet, "sender does not have shares in compositionRightTransfer")
				}
			} else {
				return nil, nil, ErrorAppend(ErrCriteriaNotMet, "sender does not have compositionRightTransfer")
			}
			compositionRightId = spec.GetCompositionRightId(compositionRightTransfer)
			compositionRightTransferHolder = true
		}
		var compositionRight Data = nil
		for _, right := range compositionRights {
			if compositionRightId == spec.GetId(right) {
				if !compositionRightTransferHolder {
					if senderId != spec.GetRecipientId(right) {
						return nil, nil, ErrorAppend(ErrCriteriaNotMet, "sender does not hold composition right")
					}
				}
				compositionRight = right
				break
			}
		}
		if compositionRight == nil {
			return nil, nil, ErrorAppend(ErrCriteriaNotMet, "could not find composition right")
		}
		licenseTerritory := spec.GetTerritory(mechanicalLicense)
		rightTerritory := spec.GetTerritory(compositionRight)
	OUTER:
		for i := range licenseTerritory {
			for j := range rightTerritory {
				if licenseTerritory[i] == rightTerritory[j] {
					rightTerritory = append(rightTerritory[:j], rightTerritory[j+1:]...)
					continue OUTER
				}
			}
			return nil, nil, ErrorAppend(ErrCriteriaNotMet, "license territory not part of right territory")
		}
		for _, composition := range moreCompositions {
			compositionId := spec.GetId(composition)
			if _, ok := seen[compositionId]; ok {
				return nil, nil, ErrorAppend(ErrCriteriaNotMet, "cannot license composition multiple times")
			}
			seen[compositionId] = struct{}{}
		}
		compositions = append(compositions, moreCompositions...)
	}
	if len(compositions) == 0 {
		return nil, nil, ErrorAppend(ErrCriteriaNotMet, "empty mechanical license; no compositions")
	}
	recipientId := spec.GetRecipientId(mechanicalLicense)
	if _, err = QueryAndValidateModel(recipientId); err != nil {
		return nil, nil, err
	}
	return mechanicalLicense, compositions, nil
}

func ProveMechanicalLicenseHolder(challenge, mechanicalLicenseId string, priv crypto.PrivateKey) (crypto.Signature, error) {
	mechanicalLicense, _, err := ValidateMechanicalLicense(mechanicalLicenseId)
	if err != nil {
		return nil, err
	}
	recipientId := spec.GetRecipientId(mechanicalLicense)
	tx, err := bigchain.GetTx(recipientId)
	if err != nil {
		return nil, err
	}
	recipientPub := bigchain.DefaultGetTxSender(tx)
	if pub := priv.Public(); !recipientPub.Equals(pub) {
		return nil, ErrorAppend(ErrInvalidKey, pub.String())
	}
	hash, err := DefaultBalloonHash(challenge)
	if err != nil {
		return nil, err
	}
	return priv.Sign(hash), nil
}

func VerifyMechanicalLicenseHolder(challenge, mechanicalLicenseId string, sig crypto.Signature) error {
	mechanicalLicense, _, err := ValidateMechanicalLicense(mechanicalLicenseId)
	if err != nil {
		return err
	}
	recipientId := spec.GetRecipientId(mechanicalLicense)
	tx, err := bigchain.GetTx(recipientId)
	if err != nil {
		return err
	}
	recipientPub := bigchain.DefaultGetTxSender(tx)
	hash, err := DefaultBalloonHash(challenge)
	if err != nil {
		return err
	}
	if !recipientPub.Verify(hash, sig) {
		return ErrorAppend(ErrInvalidSignature, sig.String())
	}
	return nil
}

func ValidateRecording(recordingId string) (Data, error) {
	tx, err := QueryAndValidateModel(recordingId)
	if err != nil {
		return nil, err
	}
	recording := bigchain.GetTxData(tx)
	senderPub := bigchain.DefaultGetTxSender(tx)
	performerId := spec.GetPerformerId(recording)
	tx, err = QueryAndValidateModel(performerId)
	if err != nil {
		return nil, err
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "performer is not recording sender")
	}
	compositionId := spec.GetRecordingOfId(recording)
	composition, err := ValidateComposition(compositionId)
	if err != nil {
		return nil, err
	}
	if performerId == spec.GetComposerId(composition) {
		return recording, nil
		// what if composer is no longer composition right-holder?
	}
	compositionRightId := spec.GetCompositionRightId(recording)
	if !EmptyStr(compositionRightId) {
		publicationId := spec.GetPublicationId(recording)
		_, compositions, compositionRights, err := ValidatePublication(publicationId)
		if err != nil {
			return nil, err
		}
		found := false
		for _, composition := range compositions {
			if compositionId == spec.GetId(composition) {
				found = true
				break
			}
		}
		if !found {
			return nil, ErrorAppend(ErrCriteriaNotMet, "publication does not link to composition")
		}
		rightHolder := false
		for _, compositionRight := range compositionRights {
			if compositionRightId == spec.GetId(compositionRight) {
				if performerId != spec.GetRecipientId(compositionRight) {
					return nil, ErrorAppend(ErrCriteriaNotMet, "sender does not hold composition right")
				}
				rightHolder = true
				break
			}
		}
		if !rightHolder {
			return nil, ErrorAppend(ErrCriteriaNotMet, "sender does not hold composition right")
		}
	}
	mechanicalLicenseId := spec.GetMechanicalLicenseId(recording)
	mechanicalLicense, compositions, err := ValidateMechanicalLicense(mechanicalLicenseId)
	if err != nil {
		return nil, err
	}
	if performerId != spec.GetRecipientId(mechanicalLicense) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "perfomer is not mechanical license holder")
	}
	for _, composition := range compositions {
		if compositionId == spec.GetId(composition) {
			return recording, nil
		}
	}
	return nil, ErrorAppend(ErrCriteriaNotMet, "mechanical license does not cover composition")
}

func ProvePerformer(challenge string, priv crypto.PrivateKey, recordingId string) (crypto.Signature, error) {
	recording, err := ValidateRecording(recordingId)
	if err != nil {
		return nil, err
	}
	performerId := spec.GetPerformerId(recording)
	tx, err := bigchain.GetTx(performerId)
	if err != nil {
		return nil, err
	}
	senderPub := bigchain.DefaultGetTxSender(tx)
	if pub := priv.Public(); !senderPub.Equals(pub) {
		return nil, ErrorAppend(ErrInvalidKey, pub.String())
	}
	hash, err := DefaultBalloonHash(challenge)
	if err != nil {
		return nil, err
	}
	return priv.Sign(hash), nil
}

func VerifyPerformer(challenge, recordingId string, sig crypto.Signature) error {
	recording, err := ValidateRecording(recordingId)
	if err != nil {
		return err
	}
	performerId := spec.GetPerformerId(recording)
	tx, err := bigchain.GetTx(performerId)
	if err != nil {
		return err
	}
	hash, err := DefaultBalloonHash(challenge)
	if err != nil {
		return err
	}
	senderPub := bigchain.DefaultGetTxSender(tx)
	if !senderPub.Verify(hash, sig) {
		return ErrorAppend(ErrInvalidSignature, sig.String())
	}
	return nil
}

func GetComposition(data Data) (Data, error) {
	compositionId := spec.GetRecordingOfId(data)
	tx, err := bigchain.GetTx(compositionId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetMechanicalLicense(data Data) (Data, error) {
	mechanicalLicenseId := spec.GetMechanicalLicenseId(data)
	tx, err := bigchain.GetTx(mechanicalLicenseId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func GetPerformer(data Data) (Data, error) {
	performerId := spec.GetPerformerId(data)
	tx, err := bigchain.GetTx(performerId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func QueryRecordingField(field, recordingId string) (interface{}, error) {
	recording, err := ValidateRecording(recordingId)
	if err != nil {
		return nil, err
	}
	switch field {
	case "composition":
		return GetComposition(recording)
	case "composition_right":
		return GetCompositionRight(recording)
	case "mechanical_license":
		return GetMechanicalLicense(recording)
	case "performer":
		return GetPerformer(recording)
	}
	return nil, ErrorAppend(ErrInvalidField, field)
}

func ProveRecordingRightHolder(challenge string, priv crypto.PrivateKey, recordingRightId, releaseId string) (crypto.Signature, error) {
	_, _, recordingRights, err := ValidateRelease(releaseId)
	if err != nil {
		return nil, err
	}
	for _, recordingRight := range recordingRights {
		if recordingRightId == spec.GetId(recordingRight) {
			recipientId := spec.GetRecipientId(recordingRight)
			tx, err := bigchain.GetTx(recipientId)
			if err != nil {
				return nil, err
			}
			recipientPub := bigchain.DefaultGetTxSender(tx)
			if pub := priv.Public(); !recipientPub.Equals(pub) {
				return nil, ErrorAppend(ErrInvalidKey, pub.String())
			}
			hash, err := DefaultBalloonHash(challenge)
			if err != nil {
				return nil, err
			}
			return priv.Sign(hash), nil
		}
	}
	return nil, ErrorAppend(ErrCriteriaNotMet, "release does not link to recording right")
}

func VerifyRecordingRightHolder(challenge, recordingRightId, releaseId string, sig crypto.Signature) error {
	_, _, recordingRights, err := ValidateRelease(releaseId)
	if err != nil {
		return err
	}
	for _, recordingRight := range recordingRights {
		if recordingRightId == spec.GetId(recordingRight) {
			recipientId := spec.GetRecipientId(recordingRight)
			tx, err := bigchain.GetTx(recipientId)
			if err != nil {
				return err
			}
			hash, err := DefaultBalloonHash(challenge)
			if err != nil {
				return err
			}
			recipientPub := bigchain.DefaultGetTxSender(tx)
			if !recipientPub.Verify(hash, sig) {
				return ErrorAppend(ErrInvalidSignature, sig.String())
			}
			return nil
		}
	}
	return ErrorAppend(ErrCriteriaNotMet, "release does not link to recording right")
}

func ValidateRelease(releaseId string) (Data, []Data, []Data, error) {
	tx, err := QueryAndValidateModel(releaseId)
	if err != nil {
		return nil, nil, nil, err
	}
	var performerId string
	release := bigchain.GetTxData(tx)
	senderPub := bigchain.DefaultGetTxSender(tx)
	recordLabelId := spec.GetRecordLabelId(release)
	recordingIds := spec.GetRecordingIds(release)
	recordings := make([]Data, len(recordingIds))
	for i, recordingId := range recordingIds {
		recording, err := ValidateRecording(recordingId)
		if err != nil {
			return nil, nil, nil, err
		}
		if i == 0 {
			performerId = spec.GetPerformerId(recording)
			// TODO: check performerId
		} else if performerId != spec.GetPerformerId(recording) {
			return nil, nil, nil, ErrorAppend(ErrCriteriaNotMet, "release cannot link to recording with different performers")
		}
		recording.Set("@id", recordingId)
		recordings[i] = recording
	}
	recipientIds := make(map[string]struct{})
	recordingRightIds := spec.GetRecordingRightIds(release)
	recordingRights := make([]Data, len(recordingRightIds))
	rightHolder := false
	totalShares := 0
	for i, rightId := range recordingRightIds {
		recordingRight, recipientPub, _, err := ValidateRight(rightId)
		if err != nil {
			return nil, nil, nil, err
		}
		if performerId != spec.GetSenderId(recordingRight) {
			return nil, nil, nil, ErrorAppend(ErrCriteriaNotMet, "performer must be right sender")
		}
		recipientId := spec.GetRecipientId(recordingRight)
		if _, ok := recipientIds[recipientId]; ok {
			return nil, nil, nil, ErrorAppend(ErrCriteriaNotMet, "holder cannot have multiple assignments")
		}
		if !rightHolder && recordLabelId == recipientId {
			if !senderPub.Equals(recipientPub) {
				return nil, nil, nil, ErrorAppend(ErrCriteriaNotMet, "record label is not release sender")
			}
			rightHolder = true
		}
		recipientIds[recipientId] = struct{}{}
		shares := spec.GetRecipientShares(recordingRight)
		if shares <= 0 {
			return nil, nil, nil, ErrorAppend(ErrCriteriaNotMet, "percentage shares must be greater than 0")
		}
		if totalShares += shares; totalShares > 100 {
			return nil, nil, nil, ErrorAppend(ErrCriteriaNotMet, "total percentage shares cannot exceed 100")
		}
		recordingRight.Set("@id", rightId)
		recordingRights[i] = recordingRight
	}
	if !rightHolder {
		return nil, nil, nil, ErrorAppend(ErrCriteriaNotMet, "record label is not recording right holder")
	}
	if totalShares != 100 {
		return nil, nil, nil, ErrorAppend(ErrCriteriaNotMet, "total percentage shares do not equal 100")
	}
	return release, recordings, recordingRights, nil
}

func ProveRecordLabel(challenge string, priv crypto.PrivateKey, releaseId string) (crypto.Signature, error) {
	release, _, _, err := ValidateRelease(releaseId)
	if err != nil {
		return nil, err
	}
	recordLabelId := spec.GetRecordLabelId(release)
	tx, err := bigchain.GetTx(recordLabelId)
	if err != nil {
		return nil, err
	}
	senderPub := bigchain.DefaultGetTxSender(tx)
	if pub := priv.Public(); !senderPub.Equals(pub) {
		return nil, ErrorAppend(ErrInvalidKey, pub.String())
	}
	hash, err := DefaultBalloonHash(challenge)
	if err != nil {
		return nil, err
	}
	return priv.Sign(hash), nil
}

func VerifyRecordLabel(challenge, releaseId string, sig crypto.Signature) error {
	release, _, _, err := ValidateRelease(releaseId)
	if err != nil {
		return err
	}
	recordLabelId := spec.GetRecordLabelId(release)
	tx, err := bigchain.GetTx(recordLabelId)
	if err != nil {
		return err
	}
	hash, err := DefaultBalloonHash(challenge)
	if err != nil {
		return err
	}
	senderPub := bigchain.DefaultGetTxSender(tx)
	if !senderPub.Verify(hash, sig) {
		return ErrorAppend(ErrInvalidSignature, sig.String())
	}
	return nil
}

func GetRecordLabel(data Data) (Data, error) {
	recordLabelId := spec.GetRecordLabelId(data)
	tx, err := bigchain.GetTx(recordLabelId)
	if err != nil {
		return nil, err
	}
	return bigchain.GetTxData(tx), nil
}

func QueryReleaseField(field, releaseId string) (interface{}, error) {
	release, recordings, recordingRights, err := ValidateRelease(releaseId)
	if err != nil {
		return nil, err
	}
	switch field {
	case "recordings":
		return recordings, nil
	case "recording_rights":
		return recordingRights, nil
	case "record_label":
		return GetRecordLabel(release)
	}
	return nil, ErrorAppend(ErrInvalidField, field)
}

func ValidateRecordingRightTransfer(recordingRightTransferId string) (Data, error) {
	tx, err := QueryAndValidateModel(recordingRightTransferId)
	if err != nil {
		return nil, err
	}
	recordingRightTransfer := bigchain.GetTxData(tx)
	senderPub := bigchain.DefaultGetTxSender(tx)
	recipientId := spec.GetRecipientId(recordingRightTransfer)
	tx, err = QueryAndValidateModel(recipientId)
	if err != nil {
		return nil, err
	}
	recipientPub := bigchain.DefaultGetTxSender(tx)
	if senderPub.Equals(recipientPub) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "recipient and sender keys must be different")
	}
	senderId := spec.GetSenderId(recordingRightTransfer)
	tx, err = QueryAndValidateModel(senderId)
	if err != nil {
		return nil, err
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, ErrorAppend(ErrInvalidKey, senderPub.String())
	}
	releaseId := spec.GetReleaseId(recordingRightTransfer)
	_, _, recordingRights, err := ValidateRelease(releaseId)
	if err != nil {
		return nil, err
	}
	txId := spec.GetTxId(recordingRightTransfer)
	tx, err = bigchain.GetTx(txId)
	if err != nil {
		return nil, err
	}
	if bigchain.TRANSFER != bigchain.GetTxOperation(tx) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "expected TRANSFER tx")
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "sender is not signer of TRANSFER tx")
	}
	n := len(bigchain.GetTxOutputs(tx))
	if n != 1 && n != 2 {
		return nil, ErrorAppend(ErrInvalidSize, "tx outputs must have size 1 or 2")
	}
	if !recipientPub.Equals(bigchain.GetTxRecipient(tx, 0)) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "recipient does not hold primary output of TRANSFER tx")
	}
	recipientShares := bigchain.GetTxOutputAmount(tx, 0)
	if recipientShares <= 0 || recipientShares > 100 {
		return nil, ErrorAppend(ErrCriteriaNotMet, "recipient shares must be greater than 0 and less than/equal to 100")
	}
	recordingRightTransfer.Set("recipientShares", recipientShares)
	if n == 2 {
		if !senderPub.Equals(bigchain.GetTxRecipient(tx, 1)) {
			return nil, ErrorAppend(ErrCriteriaNotMet, "sender does not hold secondary output of TRANSFER tx")
		}
		senderShares := bigchain.GetTxOutputAmount(tx, 1)
		if senderShares < 0 || senderShares > 100 {
			return nil, ErrorAppend(ErrCriteriaNotMet, "sender shares cannot be less than 0 or greater than 100")
		}
		recordingRightTransfer.Set("senderShares", senderShares)
	}
	found := false
	recordingRightId := spec.GetRecordingRightId(recordingRightTransfer)
	if recordingRightId != bigchain.GetTxAssetId(tx) {
		return nil, ErrorAppend(ErrCriteriaNotMet, "TRANSFER tx does not link to correct recording right")
	}
	for _, recordingRight := range recordingRights {
		if recordingRightId == bigchain.GetId(recordingRight) {
			found = true
			break
		}
	}
	if !found {
		return nil, ErrorAppend(ErrCriteriaNotMet, "release does not link to recording right")
	}
	return recordingRightTransfer, nil
}

func ProveRecordingRightTransferHolder(challenge, holderId string, priv crypto.PrivateKey, recordingRightTransferId, releaseId string) (crypto.Signature, error) {
	recordingRightTransfer, err := ValidateRecordingRightTransfer(recordingRightTransferId)
	if err != nil {
		return nil, err
	}
	recordingRightId := spec.GetRecordingRightId(recordingRightTransfer)
	if holderId == spec.GetRecipientId(recordingRightTransfer) {
		//..
	} else if holderId == spec.GetSenderId(recordingRightTransfer) {
		if spec.GetSenderShares(recordingRightTransfer) == 0 {
			return nil, ErrorAppend(ErrCriteriaNotMet, "sender does not have any shares")
		}
	} else {
		return nil, ErrorAppend(ErrCriteriaNotMet, "holder is not recipient or sender of transfer")
	}
	_, _, recordingRights, err := ValidateRelease(releaseId)
	if err != nil {
		return nil, err
	}
	for _, recordingRight := range recordingRights {
		if recordingRightId == spec.GetId(recordingRight) {
			tx, err := bigchain.GetTx(holderId)
			if err != nil {
				return nil, err
			}
			holderPub := bigchain.DefaultGetTxSender(tx)
			if pub := priv.Public(); !holderPub.Equals(pub) {
				return nil, ErrorAppend(ErrInvalidKey, pub.String())
			}
			hash, err := DefaultBalloonHash(challenge)
			if err != nil {
				return nil, err
			}
			return priv.Sign(hash), nil
		}
	}
	return nil, ErrorAppend(ErrCriteriaNotMet, "release does not link to underlying recording right")
}

func VerifyRecordingRightTransferHolder(challenge, holderId, recordingRightTransferId, releaseId string, sig crypto.Signature) error {
	recordingRightTransfer, err := ValidateRecordingRightTransfer(recordingRightTransferId)
	if err != nil {
		return err
	}
	recordingRightId := spec.GetRecordingRightId(recordingRightTransfer)
	if holderId == spec.GetRecipientId(recordingRightTransfer) {
		//..
	} else if holderId == spec.GetSenderId(recordingRightTransfer) {
		if spec.GetSenderShares(recordingRightTransfer) == 0 {
			return ErrorAppend(ErrCriteriaNotMet, "sender does not have any shares")
		}
	} else {
		return ErrorAppend(ErrCriteriaNotMet, "holder is not recipient or sender of transfer")
	}
	_, _, recordingRights, err := ValidateRelease(releaseId)
	if err != nil {
		return err
	}
	for _, recordingRight := range recordingRights {
		if recordingRightId == spec.GetId(recordingRight) {
			tx, err := bigchain.GetTx(holderId)
			if err != nil {
				return err
			}
			hash, err := DefaultBalloonHash(challenge)
			if err != nil {
				return err
			}
			holderPub := bigchain.DefaultGetTxSender(tx)
			if !holderPub.Verify(hash, sig) {
				return ErrorAppend(ErrInvalidSignature, sig.String())
			}
			return nil
		}
	}
	return ErrorAppend(ErrCriteriaNotMet, "release does not link to underlying recording right")
}

func QueryMasterLicenseField(field, masterLicenseId string) (interface{}, error) {
	masterLicense, recordings, err := ValidateMasterLicense(masterLicenseId)
	if err != nil {
		return nil, err
	}
	switch field {
	case "recipient":
		return GetRecipient(masterLicense)
	case "recordings":
		return recordings, nil
	case "sender":
		return GetSender(masterLicense)
	}
	return nil, ErrorAppend(ErrInvalidField, field)
}

func ValidateMasterLicense(masterLicenseId string) (Data, []Data, error) {
	tx, err := QueryAndValidateModel(masterLicenseId)
	if err != nil {
		return nil, nil, err
	}
	masterLicense := bigchain.GetTxData(tx)
	senderPub := bigchain.DefaultGetTxSender(tx)
	senderId := spec.GetSenderId(masterLicense)
	tx, err = QueryAndValidateModel(senderId)
	if err != nil {
		return nil, nil, err
	}
	if !senderPub.Equals(bigchain.DefaultGetTxSender(tx)) {
		return nil, nil, ErrorAppend(ErrInvalidKey, senderPub.String())
	}
	var recordings []Data
	recordingIds := spec.GetRecordingIds(masterLicense)
	seen := make(map[string]struct{})
	if n := len(recordingIds); n > 0 {
		recordings = make([]Data, n)
		for i, recordingId := range recordingIds {
			if _, ok := seen[recordingId]; ok {
				return nil, nil, ErrorAppend(ErrCriteriaNotMet, "cannot license recording multiple times")
			}
			recording, err := ValidateRecording(recordingId)
			if err != nil {
				return nil, nil, err
			}
			if senderId != spec.GetPerformerId(recording) {
				return nil, nil, ErrorAppend(ErrCriteriaNotMet, "cannot license recording by another performer")
			}
			seen[recordingId] = struct{}{}
			spec.SetId(recording, recordingId)
			recordings[i] = recording
		}
	}
	releaseId := spec.GetReleaseId(masterLicense)
	if !EmptyStr(releaseId) {
		_, moreRecordings, recordingRights, err := ValidateRelease(releaseId)
		if err != nil {
			return nil, nil, err
		}
		recordingRightId := spec.GetRecordingRightId(masterLicense)
		recordingRightTransferHolder := false
		if EmptyStr(recordingRightId) {
			recordingRightTransferId := spec.GetRecordingRightTransferId(masterLicense)
			recordingRightTransfer, err := ValidateRecordingRightTransfer(recordingRightTransferId)
			if err != nil {
				return nil, nil, err
			}
			if releaseId != spec.GetReleaseId(recordingRightTransfer) {
				return nil, nil, ErrorAppend(ErrCriteriaNotMet, "transfer links to wrong release")
			}
			if senderId == spec.GetRecipientId(recordingRightTransfer) {
				//..
			} else if senderId == spec.GetSenderId(recordingRightTransfer) {
				if spec.GetSenderShares(recordingRightTransfer) == 0 {
					return nil, nil, ErrorAppend(ErrCriteriaNotMet, "sender does not have shares in transfer")
				}
			} else {
				return nil, nil, ErrorAppend(ErrCriteriaNotMet, "sender does not have transfer")
			}
			recordingRightId = spec.GetRecordingRightId(recordingRightTransfer)
			recordingRightTransferHolder = true
		}
		var recordingRight Data = nil
		for _, right := range recordingRights {
			if recordingRightId == spec.GetId(right) {
				if !recordingRightTransferHolder {
					if senderId != spec.GetRecipientId(right) {
						return nil, nil, ErrorAppend(ErrCriteriaNotMet, "sender does not hold recording right")
					}
				}
				recordingRight = right
				break
			}
		}
		if recordingRight == nil {
			return nil, nil, ErrorAppend(ErrCriteriaNotMet, "could not find recording right")
		}
		licenseTerritory := spec.GetTerritory(masterLicense)
		rightTerritory := spec.GetTerritory(recordingRight)
	OUTER:
		for i := range licenseTerritory {
			for j := range rightTerritory {
				if licenseTerritory[i] == rightTerritory[j] {
					rightTerritory = append(rightTerritory[:j], rightTerritory[j+1:]...)
					continue OUTER
				}
			}
			return nil, nil, ErrorAppend(ErrCriteriaNotMet, "license territory not part of right territory")
		}
		for _, recording := range moreRecordings {
			recordingId := spec.GetId(recording)
			if _, ok := seen[recordingId]; ok {
				return nil, nil, ErrorAppend(ErrCriteriaNotMet, "cannot license recording multiple times")
			}
			seen[recordingId] = struct{}{}
		}
		recordings = append(recordings, moreRecordings...)
	}
	if len(recordings) == 0 {
		return nil, nil, ErrorAppend(ErrCriteriaNotMet, "empty master license; no recordings")
	}
	recipientId := spec.GetRecipientId(masterLicense)
	tx, err = QueryAndValidateModel(recipientId)
	if err != nil {
		return nil, nil, err
	}
	return masterLicense, recordings, nil
}

func ProveMasterLicenseHolder(challenge, masterLicenseId string, priv crypto.PrivateKey) (crypto.Signature, error) {
	masterLicense, _, err := ValidateMasterLicense(masterLicenseId)
	if err != nil {
		return nil, err
	}
	recipientId := spec.GetRecipientId(masterLicense)
	tx, err := bigchain.GetTx(recipientId)
	if err != nil {
		return nil, err
	}
	recipientPub := bigchain.DefaultGetTxSender(tx)
	if pub := priv.Public(); !recipientPub.Equals(pub) {
		return nil, ErrorAppend(ErrInvalidKey, pub.String())
	}
	hash, err := DefaultBalloonHash(challenge)
	if err != nil {
		return nil, err
	}
	return priv.Sign(hash), nil
}

func VerifyMasterLicenseHolder(challenge, masterLicenseId string, sig crypto.Signature) error {
	masterLicense, _, err := ValidateMasterLicense(masterLicenseId)
	if err != nil {
		return err
	}
	recipientId := spec.GetRecipientId(masterLicense)
	tx, err := bigchain.GetTx(recipientId)
	if err != nil {
		return err
	}
	hash, err := DefaultBalloonHash(challenge)
	if err != nil {
		return err
	}
	recipientPub := bigchain.DefaultGetTxSender(tx)
	if !recipientPub.Verify(hash, sig) {
		return ErrorAppend(ErrInvalidSignature, sig.String())
	}
	return nil
}
