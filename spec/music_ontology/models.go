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

// Signature

func NewSignature(impl string, pubId interface{}, sig string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(CONTEXT)
		pubId = spec.LinkIPLD(pubId)
	}
	if impl == spec.JSON {
		context = CONTEXT
	}
	return spec.Data{
		"@context":           context,
		"@type":              "sec:LinkedDataSignature2015",
		"sec:created":        Timestr(),
		"sec:creator":        pubId,
		"sec:signatureValue": sig,
	}
}

func SignData(impl string, data spec.Data, sigId interface{}) {
	if impl == spec.IPLD {
		sigId = spec.LinkIPLD(sigId)
	}
	data["sec:signature"] = sigId
}

// Agents

func NewArtist(impl, name, openId, pub string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(CONTEXT)
	}
	if impl == spec.JSON {
		context = CONTEXT
	}
	return spec.Data{
		"@context":    context,
		"@type":       "mo:MusicArtist",
		"foaf:name":   name,
		"foaf:openid": openId,
		"publicKey":   pub,
	}
}

func NewOrganization(impl, name, openId, pub string) spec.Data {
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
		"foaf:name":   name,
		"foaf:openid": openId,
		"publicKey":   pub,
	}
}

func NewPublisher(impl, name, openId, pub string) spec.Data {
	data := NewOrganization(impl, name, openId, pub)
	// How to differentiate publisher from general org?
	return data
}

func NewLabel(impl, lc, name, openId, pub string) spec.Data {
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
		"foaf:name":   name,
		"foaf:openid": openId,
		"publicKey":   pub,
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

func GetPublicKey(agent spec.Data) string {
	pubId := agent["publicKey"]
	if pubData, ok := pubId.(spec.Data); ok {
		pubId = pubData[spec.LINK_SYMBOL]
	}
	return pubId.(string)
}

/*
func GetPublicKey(agent spec.Data) string {
	pubId := agent["sec:publicKey"]
	if pubIdData, ok := pubId.(spec.Data); ok {
		pubId = pubIdData[spec.LINK_SYMBOL]
	}
	return pubId.(string)
}
*/

func GetLogin(agent spec.Data) string {
	login := agent["foaf:page"]
	if loginData, ok := login.(spec.Data); ok {
		login = loginData[spec.LINK_SYMBOL]
	}
	return login.(string)
}

// Track, Record

func NewTrack(impl string, artistId interface{}, number int, recordId interface{}, title string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(context)
		artistId = spec.LinksIPLD(artistId)
		recordId = spec.LinkIPLD(recordId)
	}
	if impl == spec.JSON {
		context = CONTEXT
	}
	return spec.Data{
		"@context":        context,
		"@type":           "mo:Track",
		"dc:isPartOf":     recordId,
		"dc:title":        title,
		"foaf:maker":      artistId,
		"mo:track_number": number,
	}
}

func NewRecord(impl string, artistId interface{}, number int, publisherId interface{}, title string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(CONTEXT)
		artistId = spec.LinksIPLD(artistId)
		publisherId = spec.LinkIPLD(publisherId)
	}
	if impl == spec.JSON {
		context = CONTEXT
	}
	data := spec.Data{
		"@context":     context,
		"@type":        "mo:Record",
		"dc:title":     title,
		"foaf:maker":   artistId,
		"mo:publisher": publisherId,
	}
	if number > 0 {
		data["mo:record_number"] = number
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

func GetPublisher(album spec.Data) string {
	publisherId := album["mo:publisher"]
	if publisherData, ok := publisherId.(spec.Data); ok {
		publisherId = publisherData[spec.LINK_SYMBOL]
	}
	return publisherId.(string)
}
