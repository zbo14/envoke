package spec

import . "github.com/zbo14/envoke/common"

const CONTEXT = "http://localhost:8888/spec#Context"

func GetId(data Data) string {
	return data.GetStr("@id")
}

func SetId(data Data, id string) {
	data.Set("@id", id)
}

func GetType(data Data) string {
	return data.GetStr("@type")
}

func NewParty(email, ipi, isni string, memberIds []string, name, proId, sameAs, _type string) Data {
	party := Data{
		"@context": CONTEXT,
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

func GetIPI(data Data) string {
	return data.GetStr("ipiNumber")
}

func GetISNI(data Data) string {
	return data.GetStr("isniNumber")
}

func GetName(data Data) string {
	return data.GetStr("name")
}

func GetPROId(data Data) string {
	pro := data.GetData("pro")
	return GetId(pro)
}

func GetSameAs(data Data) string {
	return data.GetStr("sameAs")
}

// TODO: add lyricist

func NewComposition(composerId, hfa, iswc, name string) Data {
	composition := Data{
		"@context": CONTEXT,
		"@type":    "MusicComposition",
		"composer": Data{"@id": composerId},
		"name":     name,
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

func GetHFA(data Data) string {
	return data.GetStr("hfaCode")
}

func GetISWC(data Data) string {
	return data.GetStr("iswcCode")
}

func NewPublication(compositionIds []string, compositionRightIds []string, publisherId string) Data {
	n := len(compositionIds)
	if n == 0 {
		panic("No compositionIds")
	}
	compositions := make([]Data, n)
	for i, compositionId := range compositionIds {
		compositions[i] = Data{
			"@type":    "schema:ListItem",
			"position": i + 1,
			"item": Data{
				"@type": "MusicComposition",
				"@id":   compositionId,
			},
		}
	}
	if n = len(compositionRightIds); n == 0 {
		panic("No compositionRightIds")
	}
	compositionRights := make([]Data, n)
	for i, compositionRightId := range compositionRightIds {
		compositionRights[i] = Data{
			"@type":    "schema:ListItem",
			"position": i + 1,
			"item": Data{
				"@type": "CompositionRight",
				"@id":   compositionRightId,
			},
		}
	}
	return Data{
		"@context": CONTEXT,
		"@type":    "MusicPublication",
		"composition": Data{
			"@type":           "schema:ItemList",
			"numberOfItems":   n,
			"itemListElement": compositions,
		},
		"compositionRight": Data{
			"@type":           "schema:ItemList",
			"numberOfItems":   n,
			"itemListElement": compositionRights,
		},
		"publisher": Data{"@id": publisherId},
	}
}

func GetCompositionIds(data Data) []string {
	compositions := data.GetData("composition")
	n := compositions.GetInt("numberOfItems")
	compositionIds := make([]string, n)
	itemListElement := compositions.GetInterfaceSlice("itemListElement")
	for i, elem := range itemListElement {
		item := AssertData(elem).GetData("item")
		compositionIds[i] = GetId(item)
	}
	return compositionIds
}

func GetCompositionRightIds(data Data) []string {
	compositionRights := data.GetData("compositionRight")
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

// TODO: add producer

func NewRecording(compositionId, compositionRightId, duration, isrc, mechanicalLicenseId, performerId, publicationId string) Data {
	recording := Data{
		"@context": CONTEXT,
		"@type":    "MusicRecording",
		"byArtist": Data{
			"@id": performerId,
		},
		"duration": duration,
		"recordingOf": Data{
			"@id": compositionId,
		},
	}
	if !EmptyStr(compositionRightId) {
		if EmptyStr(publicationId) {
			panic("must have compositionRightId and publicationId")
		}
		recording.Set("compositionRight", Data{"@id": compositionRightId})
		recording.Set("publication", Data{"@id": publicationId})
	} else if !EmptyStr(mechanicalLicenseId) {
		recording.Set("mechanicalLicense", Data{"@id": mechanicalLicenseId})
	} else {
		// performer should be composer
	}
	if !EmptyStr(isrc) {
		recording.Set("isrc", isrc)
	}
	return recording
}

func GetCompositionRightId(data Data) string {
	compositionRight := data.GetData("compositionRight")
	return GetId(compositionRight)
}

func GetMechanicalLicenseId(data Data) string {
	mechanicalLicense := data.GetData("mechanicalLicense")
	return GetId(mechanicalLicense)
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

func GetRecordingOfId(data Data) string {
	composition := data.GetData("recordingOf")
	return GetId(composition)
}

func NewRelease(mechanicalLicenseId, name string, recordingIds, recordingRightIds []string, recordLabelId string) Data {
	n := len(recordingIds)
	if n == 0 {
		panic("No recordingIds")
	}
	recordings := make([]Data, n)
	for i, recordingId := range recordingIds {
		recordings[i] = Data{
			"@type":    "schema:ListItem",
			"position": i + 1,
			"item": Data{
				"@type": "MusicRecording",
				"@id":   recordingId,
			},
		}
	}
	if n = len(recordingRightIds); n == 0 {
		panic("No recordingRightIds")
	}
	recordingRights := make([]Data, n)
	for i, recordingRightId := range recordingRightIds {
		recordingRights[i] = Data{
			"@type":    "schema:ListItem",
			"position": i + 1,
			"item": Data{
				"@type": "RecordingRight",
				"@id":   recordingRightId,
			},
		}
	}
	release := Data{
		"@context": CONTEXT,
		"@type":    "MusicRelease",
		"name":     name,
		"recording": Data{
			"@type":           "schema:ItemList",
			"numberOfItems":   n,
			"itemListElement": recordings,
		},
		"recordingRights": Data{
			"@type":           "schema:ItemList",
			"numberOfItems":   n,
			"itemListElement": recordingRights,
		},
		"recordLabel": Data{"@id": recordLabelId},
	}
	if !EmptyStr(mechanicalLicenseId) {
		release.Set("mechanicalLicenseId", Data{"@id": mechanicalLicenseId})
	}
	return release
}

func GetRecordingIds(data Data) []string {
	recordings := data.GetData("recordings")
	n := recordings.GetInt("numberOfItems")
	recordingIds := make([]string, n)
	itemListElement := recordings.GetInterfaceSlice("itemListElement")
	for i, elem := range itemListElement {
		item := AssertData(elem).GetData("item")
		recordingIds[i] = GetId(item)
	}
	return recordingIds
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

func NewCompositionRight(recipientId, senderId string, territory []string, validFrom, validThrough string) Data {
	return NewRight(recipientId, senderId, territory, "CompositionRight", validFrom, validThrough)
}

func NewRecordingRight(recipientId, senderId string, territory []string, validFrom, validThrough string) Data {
	return NewRight(recipientId, senderId, territory, "RecordingRight", validFrom, validThrough)
}

func NewRight(recipientId, senderId string, territory []string, _type, validFrom, validThrough string) Data {
	return Data{
		"@context": CONTEXT,
		"@type":    _type,
		"recipient": Data{
			"@id": recipientId,
		},
		"sender": Data{
			"@id": senderId,
		},
		"territory":    territory,
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

// Note: txId is the hex id of a TRANSFER tx in Bigchain/IPDB
// the output amount(s) will specify shares transferred/kept

func NewCompositionRightTransfer(compositionRightId, publicationId, recipientId, senderId, txId string) Data {
	return Data{
		"@context": CONTEXT,
		"@type":    "CompositionRightTransfer",
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
		"@context": CONTEXT,
		"@type":    "RecordingRightTransfer",
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

func NewMechanicalLicense(compositionIds []string, compositionRightId, compositionRightTransferId, publicationId, recipientId, senderId string, territory, usage []string, validFrom, validThrough string) Data {
	mechanicalLicense := Data{
		"@context": CONTEXT,
		"@type":    "MechanicalLicense",
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
	n := len(compositionIds)
	if n > 0 {
		compositions := make([]Data, n)
		for i, compositionId := range compositionIds {
			compositions[i] = Data{
				"@type":    "schema:ListItem",
				"position": i + 1,
				"item": Data{
					"@type": "MusicComposition",
					"@id":   compositionId,
				},
			}
		}
		mechanicalLicense.Set("composition", Data{
			"@type":           "schema:ItemList",
			"numberOfItems":   n,
			"itemListElement": compositions,
		})
	} else if EmptyStr(publicationId) {
		panic("Expected compositionIds or publicationId")
	}
	if !EmptyStr(publicationId) {
		mechanicalLicense.Set("publication", Data{"@id": publicationId})
		if !EmptyStr(compositionRightId) {
			mechanicalLicense.Set("compositionRight", Data{"@id": compositionRightId})
		} else if !EmptyStr(compositionRightTransferId) {
			mechanicalLicense.Set("compositionRightTransfer", Data{"@id": compositionRightTransferId})
		} else {
			panic("Expected compositionRightId or compositionRightTransferId")
		}
	}
	return mechanicalLicense
}

func NewMasterLicense(recipientId string, recordingIds []string, recordingRightId, recordingRightTransferId, releaseId, senderId string, territory, usage []string, validFrom, validThrough string) Data {
	masterLicense := Data{
		"@context": CONTEXT,
		"@type":    "MasterLicense",
		"recipient": Data{
			"@id": recipientId,
		},
		"sender": Data{
			"@id": senderId,
		},
		"usage":        usage,
		"territory":    territory,
		"validFrom":    validFrom,
		"validThrough": validThrough,
	}
	n := len(recordingIds)
	if n > 0 {
		recordings := make([]Data, n)
		for i, recordingId := range recordingIds {
			recordings[i] = Data{
				"@type":    "schema:ListItem",
				"position": i + 1,
				"item": Data{
					"@type": "MusicRecording",
					"@id":   recordingId,
				},
			}
		}
		masterLicense.Set("recording", Data{
			"@type":           "schema:ItemList",
			"numberOfItems":   n,
			"itemListElement": recordings,
		})
	} else if EmptyStr(releaseId) {
		panic("Expected recordingIds or releaseId")
	}
	if !EmptyStr(releaseId) {
		masterLicense.Set("release", Data{"@id": releaseId})
		if !EmptyStr(recordingRightId) {
			masterLicense.Set("recordingRight", Data{"@id": recordingRightId})
		} else if !EmptyStr(recordingRightTransferId) {
			masterLicense.Set("recordingRightTransfer", Data{"@id": recordingRightTransferId})
		} else {
			panic("Expected recordingRightId or recordingRightTransferId")
		}
	}
	return masterLicense
}

func GetRecordingRightId(data Data) string {
	recordingRight := data.GetData("recordingRight")
	return GetId(recordingRight)
}
