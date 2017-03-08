package schema

import (
	. "github.com/zbo14/envoke/common"
)

const (
	ENVOKE = "<envoke placeholder>"
	COALA  = "<coalaip placeholder>"
	SCHEMA = "http://schema.org/"
)

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

// publisher

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

func NewPublication(compositionId string, compositionRightIds []string, publisherId string) Data {
	n := len(compositionRightIds)
	compositionRights := make([]Data, n)
	for i, recordingRightId := range compositionRightIds {
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

func NewRecording(isrc, performerId, producerId, publicationId) Data {
	return Data{
		"@context": SCHEMA,
		"@type":    "MusicRecording",
		"byArtist": Data{
			"@id": performerId,
		},
		"isrcCode": isrc,
		"producer": Data{
			"@id": producerId,
		},
		"recordingOf": Data{
			"@id": compositionId,
		},
		"publication": Data{
			"@id": publicationId,
		},
	}
}

// genre
// name

func NewRelease(recordingId, recordLabelId string, recordingRightIds []string) {
	n := len(recordingIds)
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
		"recordingRights": Data{
			"@type":           "ItemList",
			"numberOfItems":   n,
			"itemListElement": recordingRights,
		},
		"recordLabel": Data{
			"@id": recordLabelId,
		},
		"track": Data{
			"@id": recordingId,
		},
	}
}

// Note: percentageShares is taken from the tx output amount so it's not included in the data model

func NewCompositionRight(compositionId, recipientId, senderId string, territory []string, validFrom, validThrough string) Data {
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
		"usage":        usage,
		"territory":    territory,
		"validFrom":    validFrom,
		"validThrough": validThrough,
	}
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
		"usage":        usage,
		"territory":    territory,
		"validFrom":    validFrom,
		"validThrough": validThrough,
	}
}

// Note: txId is the hex id of a TRANSFER tx in Bigchain/IPDB
// the output amount(s) will specify shares transferred/kept

func NewCompositionRightTransfer(compositionRightId, recipientId, senderId, txId string) Data {
	return Data{
		"@context": []string{ENVOKE, SCHEMA},
		"@type":    ENVOKE + "/CompositionRightTransfer",
		"compositionRight": Data{
			"@id": compositionRightId,
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

func NewRecordingRightTransfer(recipientId, recordingRightId, senderId, txId string) Data {
	return Data{
		"@context": []string{ENVOKE, SCHEMA},
		"@type":    ENVOKE + "/RecordingRightTransfer",
		"recipient": Data{
			"@id": recipientId,
		},
		"recordingRight": Data{
			"@id": recordingRightId,
		},
		"sender": Data{
			"@id": senderId,
		},
		"txId": txId,
	}
}

func NewMechanicalLicense(compositionRightId, compositionRightTransferId, recipientId, senderId string, usage, territory []string, validFrom, validThrough string) Data {
	license := Data{
		"@context": []string{ENVOKE, SCHEMA},
		"@type":    ENVOKE + "/MechanicalLicense",
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
	if !EmptyStr(compositionRightId) {
		license.Set("compositionRight", Data{"@id": compositionRightId})
	} else if !EmptyStr(compositionRightTransferId) {
		license.Set("compositionRightTransfer", Data{"@id": compositionRightTransferId})
	} else {
		panic("Expected compositionRightId or compositionRightTransferId")
	}
	return license
}

func NewSynchronizationLicense(compositionRightId, compositionRightTransferId, recipientId, senderId string, usage, territory []string, validFrom, validThrough string) Data {
	license := Data{
		"@context": []string{ENVOKE, SCHEMA},
		"@type":    ENVOKE + "/SynchronizationLicense",
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
	if !EmptyStr(compositionRightId) {
		license.Set("compositionRight", Data{"@id": compositionRightId})
	} else if !EmptyStr(compositionRightTransferId) {
		license.Set("compositionRightTransfer", Data{"@id": compositionRightTransferId})
	} else {
		panic("Expected compositionRightId or compositionRightTransferId")
	}
	return license
}

func NewMasterLicense(recipientId, recordingRightId, recordingRightTransferId, senderId string, usage, territory []string, validFrom, validThrough string) Data {
	license := Data{
		"@context": []string{ENVOKE, SCHEMA},
		"@type":    ENVOKE + "/SynchronizationLicense",
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
	if !EmptyStr(recordingRightId) {
		license.Set("recordingRight", Data{"@id": recordingRightId})
	} else if !EmptyStr(recordingRightTransferId) {
		license.Set("recordingRightTransfer", Data{"@id": recordingRightTransferId})
	} else {
		panic("Expected recordingRightId or recordingRightTransferId")
	}
	return license
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
