package spec

import . "github.com/zbo14/envoke/common"

const (
	ENVOKE = "<envoke placeholder>"
	COALA  = "<coalaip placeholder>"
	SCHEMA = "http://schema.org/"
)

func GetId(data Data) string {
	return data.GetStr("@id")
}

func GetType(data Data) string {
	return data.GetStr("@type")
}

func NewParty(email, ipi, isni string, memberIds []string, name, proId, sameAs, _type string) Data {
	party := Data{
		"@context": []string{ENVOKE, SCHEMA},
		"@type":    _type,
		"email":    email,
		"name":     name,
		"sameAs":   sameAs,
	}
	switch _type {
	case "MusicGroup", "Organization":
		if len(memberIds) > 0 {
			member := make([]Data, len(memberIds))
			for i, memberId := range memberIds {
				member[i] = Data{"@id": memberId}
			}
			party.Set("member", member)
		}
	case "Person":
		//..
	default:
		panic(ErrorAppend(ErrInvalidType, _type))
	}
	if !EmptyStr(ipi) {
		party.Set("ipiNumber", ipi)
	}
	if !EmptyStr(isni) {
		party.Set("isniNumber", isni)
	}
	if !EmptyStr(proId) {
		party.Set("pro", Data{"@id": proId})
	}
	return party
}

func GetDescription(data Data) string {
	return data.GetStr("description")
}

func GetEmail(data Data) string {
	return data.GetStr("email")
}

func GetName(data Data) string {
	return data.GetStr("name")
}

func GetSameAs(data Data) string {
	return data.GetStr("sameAs")
}

func NewComposition(composerId, hfa, iswc, name string) Data {
	composition := Data{
		"@context": []string{ENVOKE, SCHEMA},
		"@type":    "MusicComposition",
		"composer": Data{
			"@id": composerId,
		},
		"name": name,
	}
	if !EmptyStr(hfa) {
		composition.Set("hfaCode", hfa)
	}
	if !EmptyStr(iswc) {
		composition.Set("iswcCode", iswc)
	}
	return composition
}

func GetComposerId(data Data) string {
	composer := data.GetData("composer")
	return GetId(composer)
}

func GetProId(data Data) string {
	pro := data.GetData("pro")
	return GetId(pro)
}

func NewPublication(compositionId string, compositionRightIds []string, publisherId string) Data {
	n := len(compositionRightIds)
	compositionRights := make([]Data, n)
	for i, compositionRightId := range compositionRightIds {
		compositionRights[i] = Data{
			"@type":    "ListItem",
			"position": i + 1,
			"item": Data{
				"@type": ENVOKE + "/CompositionRight",
				"@id":   compositionRightId,
			},
		}
	}
	return Data{
		"@context": []string{ENVOKE, SCHEMA},
		"@type":    ENVOKE + "/MusicPublication",
		"composition": Data{
			"@id": compositionId,
		},
		"compositionRights": Data{
			"@type":           "ItemList",
			"numberOfItems":   n,
			"itemListElement": compositionRights,
		},
		"publisher": Data{
			"@id": publisherId,
		},
	}
}

func GetCompositionId(data Data) string {
	compositon := data.GetData("composition")
	return GetId(compositon)
}

func GetCompositionRightIds(data Data) []string {
	compositionRights := data.GetData("compositionRights")
	n := compositionRights.GetInt("numberOfItems")
	compositionRightIds := make([]string, n)
	itemListElement := compositionRights.GetInterfaceSlice("itemListElement")
	for i, elem := range itemListElement {
		item := AssertData(elem).GetData("item")
		compositionRightIds[i] = GetId(item)
	}
	return compositionRightIds
}

func GetPublisherId(data Data) string {
	publisher := data.GetData("publisher")
	return GetId(publisher)
}

func NewRecording(compositionRightId, isrc, performerId, producerId, publicationId string) Data {
	return Data{
		"@context": []string{ENVOKE, SCHEMA},
		"@type":    "MusicRecording",
		"byArtist": Data{
			"@id": performerId,
		},
		"compositionRight": Data{
			"@id": compositionRightId,
		},
		"isrcCode": isrc,
		"producer": Data{
			"@id": producerId,
		},
		"publication": Data{
			"@id": publicationId,
		},
	}
}

func GetCompositionRightId(data Data) string {
	compositionRight := data.GetData("compositionRight")
	return GetId(compositionRight)
}

func GetPerformerId(data Data) string {
	performer := data.GetData("byArtist")
	return GetId(performer)
}

func GetProducerId(data Data) string {
	producer := data.GetData("producer")
	return GetId(producer)
}

func GetPublicationId(data Data) string {
	publication := data.GetData("publication")
	return GetId(publication)
}

func NewRelease(mechanicalLicenseId, recordingId string, recordingRightIds []string, recordLabelId string) Data {
	n := len(recordingRightIds)
	recordingRights := make([]Data, n)
	for i, recordingRightId := range recordingRightIds {
		recordingRights[i] = Data{
			"@type":    "ListItem",
			"position": i + 1,
			"item": Data{
				"@type": ENVOKE + "/RecordingRight",
				"@id":   recordingRightId,
			},
		}
	}
	return Data{
		"@context": []string{ENVOKE, SCHEMA},
		"@type":    "MusicRelease",
		"mechanicalLicense": Data{
			"@id": mechanicalLicenseId,
		},
		"recordingRights": Data{
			"@type":           "ItemList",
			"numberOfItems":   n,
			"itemListElement": recordingRights,
		},
		"recording": Data{
			"@id": recordingId,
		},
		"recordLabel": Data{
			"@id": recordLabelId,
		},
	}
}

func GetMechanicalLicenseId(data Data) string {
	mechanicalLicense := data.GetData("mechanicalLicense")
	return GetId(mechanicalLicense)
}

func GetRecordingId(data Data) string {
	recording := data.GetData("recording")
	return GetId(recording)
}

func GetRecordingRightIds(data Data) []string {
	recordingRights := data.GetData("recordingRights")
	n := recordingRights.GetInt("numberOfItems")
	recordingRightIds := make([]string, n)
	itemListElement := recordingRights.GetInterfaceSlice("itemListElement")
	for i, elem := range itemListElement {
		item := AssertData(elem).GetData("item")
		recordingRightIds[i] = GetId(item)
	}
	return recordingRightIds
}

func GetRecordLabelId(data Data) string {
	recordLabel := data.GetData("recordLabel")
	return GetId(recordLabel)
}

// Note: percentageShares is taken from the tx output amount so it's not included in the data model

func NewCompositionRight(compositionId, recipientId, senderId string, territory, usage []string, validFrom, validThrough string) Data {
	return Data{
		"@context": []string{ENVOKE, SCHEMA},
		"@type":    ENVOKE + "/CompositionRight",
		"composition": Data{
			"@id": compositionId,
		},
		"recipient": Data{
			"@id": recipientId,
		},
		"sender": Data{
			"@id": senderId,
		},
		"territory":    territory,
		"usage":        usage,
		"validFrom":    validFrom,
		"validThrough": validThrough,
	}
}

func GetRecipientId(data Data) string {
	recipient := data.GetData("recipient")
	return GetId(recipient)
}

func GetRecipientShares(data Data) int {
	return data.GetInt("recipientShares")
}

func GetSenderId(data Data) string {
	sender := data.GetData("sender")
	return GetId(sender)
}

func GetSenderShares(data Data) int {
	return data.GetInt("senderShares")
}

func GetTerritory(data Data) []string {
	return data.GetStrSlice("territory")
}

func NewRecordingRight(recipientId, recordingId, senderId string, territory, usage []string, validFrom, validThrough string) Data {
	return Data{
		"@context": []string{ENVOKE, SCHEMA},
		"@type":    ENVOKE + "/RecordingRight",
		"recipient": Data{
			"@id": recipientId,
		},
		"recording": Data{
			"@id": recordingId,
		},
		"sender": Data{
			"@id": senderId,
		},
		"territory":    territory,
		"usage":        usage,
		"validFrom":    validFrom,
		"validThrough": validThrough,
	}
}

// Note: txId is the hex id of a TRANSFER tx in Bigchain/IPDB
// the output amount(s) will specify shares transferred/kept

func NewCompositionRightTransfer(compositionRightId, publicationId, recipientId, senderId, txId string) Data {
	return Data{
		"@context": []string{ENVOKE, SCHEMA},
		"@type":    ENVOKE + "/CompositionRightTransfer",
		"compositionRight": Data{
			"@id": compositionRightId,
		},
		"publication": Data{
			"@id": publicationId,
		},
		"recipient": Data{
			"@id": recipientId,
		},
		"sender": Data{
			"@id": senderId,
		},
		"txId": txId,
	}
}

func GetCompositionRightTransferId(data Data) string {
	compositionRightTransfer := data.GetData("compositionRightTransfer")
	return GetId(compositionRightTransfer)
}

func GetTxId(data Data) string {
	return data.GetStr("txId")
}

func NewRecordingRightTransfer(recipientId, recordingRightId, releaseId, senderId, txId string) Data {
	return Data{
		"@context": []string{ENVOKE, SCHEMA},
		"@type":    ENVOKE + "/RecordingRightTransfer",
		"recipient": Data{
			"@id": recipientId,
		},
		"recordingRight": Data{
			"@id": recordingRightId,
		},
		"release": Data{
			"@id": releaseId,
		},
		"sender": Data{
			"@id": senderId,
		},
		"txId": txId,
	}
}

func GetReleaseId(data Data) string {
	release := data.GetData("release")
	return GetId(release)
}

func GetRecordingRightTransferId(data Data) string {
	recordingRightTransfer := data.GetData("recordingRightTransfer")
	return GetId(recordingRightTransfer)
}

func NewMechanicalLicense(compositionRightId, compositionRightTransferId, publicationId, recipientId, senderId string, territory, usage []string, validFrom, validThrough string) Data {
	license := Data{
		"@context": []string{ENVOKE, SCHEMA},
		"@type":    ENVOKE + "/MechanicalLicense",
		"publication": Data{
			"@id": publicationId,
		},
		"recipient": Data{
			"@id": recipientId,
		},
		"sender": Data{
			"@id": senderId,
		},
		"territory":    territory,
		"usage":        usage,
		"validFrom":    validFrom,
		"validThrough": validThrough,
	}
	if !EmptyStr(compositionRightId) {
		license.Set("compositionRight", Data{"@id": compositionRightId})
	} else if !EmptyStr(compositionRightTransferId) {
		license.Set("compositionRightTransfer", Data{"@id": compositionRightTransferId})
	} else {
		panic("Expected compositionRightId or compositionRightTransferId")
	}
	return license
}

func NewMasterLicense(recipientId, recordingRightId, recordingRightTransferId, releaseId, senderId string, territory, usage []string, validFrom, validThrough string) Data {
	license := Data{
		"@context": []string{ENVOKE, SCHEMA},
		"@type":    ENVOKE + "/SynchronizationLicense",
		"recipient": Data{
			"@id": recipientId,
		},
		"release": Data{
			"@id": releaseId,
		},
		"sender": Data{
			"@id": senderId,
		},
		"usage":        usage,
		"territory":    territory,
		"validFrom":    validFrom,
		"validThrough": validThrough,
	}
	if !EmptyStr(recordingRightId) {
		license.Set("recordingRight", Data{"@id": recordingRightId})
	} else if !EmptyStr(recordingRightTransferId) {
		license.Set("recordingRightTransfer", Data{"@id": recordingRightTransferId})
	} else {
		panic("Expected recordingRightId or recordingRightTransferId")
	}
	return license
}

func GetRecordingRightId(data Data) string {
	recordingRight := data.GetData("recordingRight")
	return GetId(recordingRight)
}

//----------------------COALA------------------------

func NewCopyright(manifestationId string, territoryNames []string, validFrom, validThrough string) Data {
	territory := make([]Data, len(territoryNames))
	for i, name := range territoryNames {
		territory[i] = Data{
			"@type": "Place",
			"name":  name,
		}
	}
	return Data{
		"@context": []string{COALA, SCHEMA},
		"@type":    COALA + "/Copyright",
		"rightsOf": Data{
			"@id": manifestationId,
		},
		"territory":    territory,
		"validFrom":    validFrom,
		"validThrough": validThrough,
	}
}

func NewRight(copyrightId string, percentageShares int, territoryNames, usageType []string, validFrom, validThrough string) Data {
	territory := make([]Data, len(territoryNames))
	for i, name := range territoryNames {
		territory[i] = Data{
			"@type": "Place",
			"name":  name,
		}
	}
	return Data{
		"@context":         []string{COALA, SCHEMA},
		"@type":            COALA + "/Right",
		"percentageShares": percentageShares,
		"source": Data{
			"@id": copyrightId,
		},
		"territory":    territory,
		"usageType":    usageType,
		"validFrom":    validFrom,
		"validThrough": validThrough,
	}
}
