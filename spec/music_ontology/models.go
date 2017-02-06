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

// json-ld for now

func SetId(data spec.Data, id string) {
	data["@id"] = id
}

func NewTrack(impl string, artist spec.Data, number int, recordId interface{}, title string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(context)
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

// PubKey

func NewPublicKey(impl, pem string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(CONTEXT)
	}
	if impl == spec.JSON {
		context = CONTEXT
	}
	return spec.Data{
		"@context":         context,
		"@type":            "sec:publicKey",
		"sec:publicKeyPem": pem,
	}
}

func AddOwner(impl string, ownerId interface{}, pub spec.Data) {
	if impl == spec.IPLD {
		ownerId = spec.LinkIPLD(ownerId)
	}
	pub["sec:owner"] = ownerId
}

func GetOwner(pub spec.Data) string {
	ownerId := pub["sec:owner"]
	if ownerData, ok := ownerId.(spec.Data); ok {
		ownerId = ownerData[spec.LINK_SYMBOL]
	}
	return ownerId.(string)
}

func GetPEM(pub spec.Data) string {
	pem := pub["sec:publicKeyPem"]
	if pemData, ok := pem.(spec.Data); ok {
		pem = pemData[spec.LINK_SYMBOL]
	}
	return pem.(string)
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

func AddPublicKey(impl string, org spec.Data, pubId interface{}) {
	if impl == spec.IPLD {
		pubId = spec.LinkIPLD(pubId)
	}
	org["sec:publicKey"] = pubId
}

func GetPublicKey(org spec.Data) string {
	pubId := org["sec:publicKey"]
	if pubIdData, ok := pubId.(spec.Data); ok {
		pubId = pubIdData[spec.LINK_SYMBOL]
	}
	return pubId.(string)
}

func GetLogin(org spec.Data) string {
	login := org["foaf:page"]
	if loginData, ok := login.(spec.Data); ok {
		login = loginData[spec.LINK_SYMBOL]
	}
	return login.(string)
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
