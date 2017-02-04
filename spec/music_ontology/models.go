package music_ontology

import (
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/spec"
)

const LC_REGEX = `^LC-\d{4,5}$`

var CONTEXT = spec.Data{
	"mo":    "http://purl.org/ontology/mo/",
	"dc":    "http://purl.org/dc/elements/1.1/",
	"xsd":   "http://www.w3.org/2001/XMLSchema#",
	"tl":    "http://purl.org/NET/c4dm/timeline.owl#",
	"event": "http://purl.org/NET/c4dm/event.owl#",
	"foaf":  "http://xmlns.com/foaf/0.1/",
	"rdfs":  "http://www.w3.org/2000/01/rdf-schema#",
}

func NewTrack(impl, id, title string, artist interface{}, number int, recordId interface{}) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(context)
		artist = spec.LinkIPLD(artist)
		recordId = spec.LinkIPLD(recordId)
	}
	if impl == spec.JSON {
		context = CONTEXT
	}
	data := spec.Data{
		"@context":        context,
		"@type":           "mo:Track",
		"dc:title":        title,
		"foaf:maker":      artist,
		"mo:track_number": number,
		"dc:isPartOf": spec.Data{
			"@id":   recordId,
			"@type": "mo:Record",
		},
	}
	if impl == spec.JSON {
		if id != "" {
			data["@id"] = id
		}
	}
	return data
}

func NewArtist(impl, id, name, openId string, partnerId interface{}) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(CONTEXT)
		partnerId = spec.LinkIPLD(partnerId)
	}
	if impl == spec.JSON {
		context = CONTEXT
	}
	data := spec.Data{
		"@context":    context,
		"@type":       "mo:MusicArtist",
		"foaf:name":   name,
		"foaf:openid": openId,
		"foaf:member": spec.Data{
			"@id":   partnerId,
			"@type": "foaf:Agent",
		},
	}
	if impl == spec.JSON {
		if id != "" {
			data["@id"] = id
		}
	}
	return data
}

func NewPublisher(impl, id, name, login string) spec.Data {
	data := NewOrganization(impl, id, name, login)
	// How to differentiate publisher from general org?
	return data
}

func NewOrganization(impl, id, name, login string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(CONTEXT)
	}
	if impl == spec.JSON {
		context = CONTEXT
	}
	data := spec.Data{
		"@context":  context,
		"@type":     "foaf:Organization",
		"foaf:name": name,
		"foaf:page": login,
	}
	if impl == spec.JSON {
		if id != "" {
			data["@id"] = id
		}
	}
	return data
}

func NewLabel(impl, id, name, lc, login string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(CONTEXT)
	}
	if impl == spec.JSON {
		context = CONTEXT
	}
	data := spec.Data{
		"@context":  context,
		"@type":     "mo:Label",
		"foaf:name": name,
		"foaf:page": login,
	}
	if impl == spec.JSON {
		if id != "" {
			data["@id"] = id
		}
	}
	if lc != "" {
		if !MatchString(LC_REGEX, lc) {
			panic("Label code does not match regex")
		}
		data["mo:lc"] = lc
	}
	return data
}

func NewRecord(impl, id string, title string, number int, publisherId interface{}) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(CONTEXT)
		publisherId = spec.LinkIPLD(publisherId)
	}
	if impl == spec.JSON {
		context = CONTEXT
	}
	data := spec.Data{
		"@context": context,
		"@type":    "mo:Record",
		"dc:title": title,
		"mo:publisher": spec.Data{
			"@id":   publisherId,
			"@type": "foaf:Agent",
		},
	}
	if impl == spec.JSON {
		if id != "" {
			data["@id"] = id
		}
	}
	if number > 0 {
		data["record_number"] = number
	}
	return data
}

func AddTracks(impl string, record spec.Data, trackIds []interface{}) {
	track := make([]spec.Data, len(trackIds))
	for i, trackId := range trackIds {
		if impl == spec.IPLD {
			trackId = spec.LinkIPLD(trackId)
		}
		track[i] = spec.Data{
			"@id":   trackId,
			"@type": "mo:Track",
		}
	}
	record["track"] = track
}
