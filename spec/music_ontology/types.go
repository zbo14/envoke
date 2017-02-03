package music_ontology

import (
	"github.com/zballs/envoke/spec"
)

const LC_REGEX = `^LC-\d{4,5}$`

// json-ld for now
// should we just pass URIs?

var CONTEXT = spec.Data{
	"mo":    "http://purl.org/ontology/mo/",
	"dc":    "http://purl.org/dc/elements/1.1/",
	"xsd":   "http://www.w3.org/2001/XMLSchema#",
	"tl":    "http://purl.org/NET/c4dm/timeline.owl#",
	"event": "http://purl.org/NET/c4dm/event.owl#",
	"foaf":  "http://xmlns.com/foaf/0.1/",
	"rdfs":  "http://www.w3.org/2000/01/rdf-schema#",
}

func NewTrack(id, title, artistId, recordId string) spec.Data {
	return spec.Data{
		"@context": CONTEXT,
		"@id":      id,
		"@type":    "mo:Track",
		"dc:title": title,
		"foaf:maker": spec.Data{
			"@id":   artistId,
			"@type": "mo:MusicArtist",
			//"foaf:name": artistName,
		},
		"dc:isPartOf": spec.Data{
			"@id":   recordId,
			"@type": "mo:Record",
			//"foaf:name": recordName,
		},
	}
}

func NewArtist(id, name string) spec.Data {
	return spec.Data{
		"@context":  CONTEXT,
		"@id":       id,
		"@type":     "mo:MusicArtist",
		"foaf:name": name,
	}
}

func NewOrganization(id, name string) spec.Data {
	return spec.Data{
		"@context":  CONTEXT,
		"@id":       id,
		"@type":     "foaf:Organization",
		"foaf:name": name,
	}
}

// TODO: NewPublisher

func NewLabel(id, name string) spec.Data {
	// TODO: add label code
	return spec.Data{
		"@context":  CONTEXT,
		"@id":       id,
		"@type":     "mo:Label",
		"foaf:name": name,
	}
}

func NewRecord(id string, publisherId string, trackIds []string) spec.Data {
	// TODO: add record number
	tracks := make([]spec.Data, len(trackIds))
	for i, trackId := range trackIds {
		tracks[i] = spec.Data{
			"@id":   trackId,
			"@type": "mo:Track",
		}
	}
	return spec.Data{
		"@context": CONTEXT,
		"@id":      id,
		"@type":    "mo:Record",
		"mo:publisher": spec.Data{
			"@id":   publisherId,
			"@type": "foaf:Agent",
			//"foaf:name": publisherName,
		},
		"mo:track": tracks,
	}
}
