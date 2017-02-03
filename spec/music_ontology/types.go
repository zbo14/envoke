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

func NewTrack(id, title, artistId string) spec.Data {
	return spec.Data{
		"@context": CONTEXT,
		"@type":    "mo:Track",
		"dc:title": title,
		"foaf:maker": spec.Data{
			"@id":   artistId,
			"@type": "mo:MusicArtist",
			// need artist name?
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

func NewPublisher(id, name string) spec.Data {
	return spec.Data{
		"@context":  CONTEXT,
		"@id":       id,
		"@type":     "foaf:Organization",
		"foaf:name": name,
	}
}

func NewLabel(id, name string) spec.Data {
	/*
		if !MatchString(LC_REGEX, lc) {
			panic("Label code does not match regex pattern")
		}
	*/
	return spec.Data{
		"@context":  CONTEXT,
		"@id":       id,
		"@type":     "mo:Label",
		"foaf:name": name,
		// "mo:lc":    lc,
	}
}

func NewRecord(id string, number int, publisherId string, trackIds []string) spec.Data {
	if number <= 0 {
		panic("Record number must be positive")
	}
	tracks := make([]spec.Data, len(trackIds))
	for i, trackId := range trackIds {
		tracks[i] = spec.Data{
			"@id":   trackId,
			"@type": "mo:Track",
		}
	}
	return spec.Data{
		"@context":         CONTEXT,
		"@id":              id,
		"@type":            "mo:Record",
		"mo:record_number": number,
		"mo:publisher": spec.Data{
			"@id":   publisherId,
			"@type": "foaf:Agent",
		},
		"mo:track": tracks,
	}
}
