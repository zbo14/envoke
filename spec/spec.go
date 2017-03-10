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

func NewPerson(birthDate, deathDate, familyName, givenName string) Data {
	person := Data{
		"@context":   SCHEMA,
		"@type":      "Person",
		"birthDate":  birthDate,
		"familyName": familyName,
		"givenName":  givenName,
	}
	if !EmptyStr(deathDate) {
		person.Set("deathDate", deathDate)
	}
	return person
}

/*
func NewOrganization(description, email, memberIds []string, name, sameAs string) Data {
	member := make([]Data, len(memberIds))
	for i, memberId := range memberIds {
		member[i] = Data{"@id": memberId}
	}
	return Data{
		"@context":    SCHEMA,
		"@type":       "Organization",
		"description": description,
		"email":       email,
		"member":      member,
		"name":        name,
		"sameAs":      sameAs,
	}
}
*/

func NewOrganization(description, email, memberNames []string, name, sameAs string) Data {
	org := Data{
		"@context":    SCHEMA,
		"@type":       "Organization",
		"description": description,
		"email":       email,
		"name":        name,
		"sameAs":      sameAs,
	}
	if memberNames != nil {
		member := make([]Data, len(memberNames))
		for i, name := range memberNames {
			member[i] = Data{
				"@type": "Person",
				"name":  name,
			}
		}
		org.Set("member", member)
	}
	return org
}

func NewMusicGroup(description, email, genre, memberNames []string, name, sameAs string) Data {
	musicGroup := NewOrganization(description, email, memberNames, name, sameAs)
	musicGroup.Set("genre", genre)
	return musicGroup
}

func NewComposition(composerId, hfa, ipi, iswc, name, proId string) Data {
	return Data{
		"@context": []string{ENVOKE, SCHEMA},
		"@type":    "MusicComposition",
		"composer": Data{
			"@id": composerId,
		},
		"hfaCode":   hfa,
		"ipiNumber": ipi,
		"iswcCode":  iswc,
		"name":      name,
		"pro": Data{
			"@id": proId,
		},
	}
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
		"@context": SCHEMA,
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

func NewRelease(mechanicalLicenseId, recordingId, recordLabelId string, recordingRightIds []string) Data {
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
		"@context": SCHEMA,
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

func NewCompositionRightTransfer(compositionRightId, compositionRightTransferId, publicationId, recipientId, senderId, txId string) Data {
	compositionRightTransfer := Data{
		"@context": []string{ENVOKE, SCHEMA},
		"@type":    ENVOKE + "/CompositionRightTransfer",
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
	if !EmptyStr(compositionRightId) {
		compositionRightTransfer.Set("compositionRight", Data{"@id": compositionRightId})
	} else if !EmptyStr(compositionRightTransferId) {
		compositionRightTransfer.Set("compositionRightTransfer", Data{"@id": compositionRightTransferId})
	} else {
		panic("Expected compositionRightId or compositionRightTransferId")
	}
	return compositionRightTransfer
}

func GetCompositionRightTransferId(data Data) string {
	compositionRightTransfer := data.GetData("compositionRightTransfer")
	return GetId(compositionRightTransfer)
}

func GetTxId(data Data) string {
	return data.GetStr("txId")
}

func NewRecordingRightTransfer(recipientId, recordingRightId, recordingRightTransferId, releaseId, senderId, txId string) Data {
	recordingRightTransfer := Data{
		"@context": []string{ENVOKE, SCHEMA},
		"@type":    ENVOKE + "/RecordingRightTransfer",
		"recipient": Data{
			"@id": recipientId,
		},
		"release": Data{
			"@id": releaseId,
		},
		"sender": Data{
			"@id": senderId,
		},
		"txId": txId,
	}
	if !EmptyStr(recordingRightId) {
		recordingRightTransfer.Set("recordingRight", Data{"@id": recordingRightId})
	} else if !EmptyStr(recordingRightTransferId) {
		recordingRightTransfer.Set("recordingRightTransfer", Data{"@id": recordingRightTransferId})
	} else {
		panic("Expected recordingRightId or recordingRightTransferId")
	}
	return recordingRightTransfer
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
