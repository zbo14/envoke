package music_ontology

import (
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/spec"
)

const (
	LABEL     = "label"
	LC_REGEX  = `^LC-\d{4,5}$`
	PUBLISHER = "publisher"
)

var CONTEXT = spec.Data{
	"mo":    "http://purl.org/ontology/mo/",
	"dc":    "http://purl.org/dc/elements/1.1/",
	"xsd":   "http://www.w3.org/2001/XMLSchema#",
	"tl":    "http://purl.org/NET/c4dm/timeline.owl#",
	"event": "http://purl.org/NET/c4dm/event.owl#",
	"foaf":  "http://xmlns.com/foaf/0.1/",
	"rdfs":  "http://www.w3.org/2000/01/rdf-schema#",
	"sec":   "https://w3id.org/security#",
}

func SetId(data spec.Data, id string) {
	data["@id"] = id
}

func NewTrack(impl string, artist interface{}, number int, recordId interface{}, title string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(context)
		artist = spec.LinkIPLD(artist)
		recordId = spec.LinkIPLD(recordId)
	}
	if impl == spec.JSON {
		context = CONTEXT
	}
	return spec.Data{
		"@context":        context,
		"@type":           "mo:Track",
		"dc:title":        title,
		"foaf:maker":      artist,
		"mo:track_number": number,
		"dc:isPartOf":     recordId,
	}
}

func NewArtist(impl, name, openId string, partnerId interface{}) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(CONTEXT)
		partnerId = spec.LinkIPLD(partnerId)
	}
	if impl == spec.JSON {
		context = CONTEXT
	}
	return spec.Data{
		"@context":    context,
		"@type":       "mo:MusicArtist",
		"foaf:name":   name,
		"foaf:openid": openId,
		"foaf:member": partnerId,
	}
}

func NewPublicKey(impl, pem string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(CONTEXT)
	}
	if impl == spec.JSON {
		context = CONTEXT
	}
	return spec.Data{
		"@context":     context,
		"@type":        "sec:publicKey",
		"publicKeyPem": pem,
	}
}

func AddOwner(impl string, ownerId interface{}, pub spec.Data) {
	if impl == spec.IPLD {
		ownerId = spec.LinkIPLD(ownerId)
	}
	pub["sec:owner"] = ownerId
}

func NewPublisher(impl, login, name, openId string) spec.Data {
	data := NewOrganization(impl, login, name, openId)
	// How to differentiate publisher from general org?
	return data
}

func NewOrganization(impl, login, name, openId string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(CONTEXT)
	}
	if impl == spec.JSON {
		context = CONTEXT
	}
	return spec.Data{
		"@context":    context,
		"@type":       "foaf:Organization",
		"foaf:page":   login,
		"foaf:name":   name,
		"foaf:openid": openId,
	}
}

func NewLabel(impl, lc, login, name, openId string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(CONTEXT)
	}
	if impl == spec.JSON {
		context = CONTEXT
	}
	data := spec.Data{
		"@context":    context,
		"@type":       "mo:Label",
		"foaf:page":   login,
		"foaf:name":   name,
		"foaf:openid": openId,
	}
	if lc != "" {
		if !MatchString(LC_REGEX, lc) {
			panic("Label code does not match regex")
		}
		data["mo:lc"] = lc
	}
	return data
}

func AddPublicKey(impl string, agent spec.Data, pubId interface{}) {
	if impl == spec.IPLD {
		pubId = spec.LinkIPLD(pubId)
	}
	agent["sec:publicKey"] = pubId
}

func NewRecord(impl string, number int, publisherId interface{}, title string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(CONTEXT)
		publisherId = spec.LinkIPLD(publisherId)
	}
	if impl == spec.JSON {
		context = CONTEXT
	}
	data := spec.Data{
		"@context":     context,
		"@type":        "mo:Record",
		"mo:publisher": publisherId,
		"dc:title":     title,
	}
	if number > 0 {
		data["record_number"] = number
	}
	return data
}

func AddTracks(impl string, record spec.Data, trackIds []interface{}) {
	if impl == spec.IPLD {
		for i, trackId := range trackIds {
			trackIds[i] = spec.LinkIPLD(trackId)
		}
	}
	record["track"] = trackIds
}
