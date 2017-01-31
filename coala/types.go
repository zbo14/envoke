package coala

import (
//. "github.com/zballs/envoke/util"
)

// Context URLs
const (
	SCHEMA  = "http://schema.org/"
	COALAIP = "<coalaip placeholder>"

	IPLD = "ipld"
	JSON = "json_ld"
)

type Data map[string]interface{}

// specs: github.com/COALAIP/specs/data-structure/README.md

// Geo coordinates

func NewGeo(impl, lat, long string) Data {
	var context interface{}
	if impl == IPLD {
		context = LinkIPLD(SCHEMA)
	}
	if impl == JSON {
		context = SCHEMA
	}
	data := make(Data)
	data["@context"] = context
	data["@type"] = "GeoCoordinates"
	data["latitude"] = lat
	data["longitude"] = long
	return data
}

// Localizable place

func NewPlace(impl, lat, long, name string) Data {
	var context interface{}
	if impl == IPLD {
		context = LinkIPLD(SCHEMA)
	}
	if impl == JSON {
		context = SCHEMA
	}
	data := make(Data)
	data["@context"] = context
	data["@type"] = "Place"
	data["geo"] = NewGeo(impl, lat, long)
	data["name"] = name
	return data
}

// Person

func NewPerson(impl, id, givenName, familyName, birthDate string) Data {
	var context interface{}
	if impl == IPLD {
		context = LinkIPLD(SCHEMA)
	}
	if impl == JSON {
		context = SCHEMA
	}
	data := make(Data)
	data["@context"] = context
	data["@type"] = "Person"
	if id != "" {
		data["@id"] = id
	}
	data["givenName"] = givenName
	data["familyName"] = familyName
	data["birthDate"] = birthDate
	return data
}

// Organization

func NewOrganization(impl, id, name string, founder interface{}, member ...interface{}) Data {
	var context interface{}
	if impl == IPLD {
		context = LinksIPLD(SCHEMA, COALAIP)
		founder = LinkIPLD(founder)
		member = LinksIPLD(member...)
	}
	if impl == JSON {
		context = []string{SCHEMA, COALAIP}
	}
	data := make(Data)
	data["@context"] = context
	data["@type"] = "Organization"
	if id != "" {
		data["@id"] = id
	}
	data["name"] = name
	data["founder"] = founder
	data["member"] = member
	return data
}

// Work

func NewWork(impl, id, name string, author interface{}) Data {
	var context interface{}
	if impl == IPLD {
		context = LinksIPLD(SCHEMA, COALAIP)
		author = LinkIPLD(author)
	}
	if impl == JSON {
		context = []string{SCHEMA, COALAIP}
	}
	data := make(Data)
	data["@context"] = context
	data["@type"] = "Work"
	if id != "" {
		data["@id"] = id
	}
	data["name"] = name
	data["author"] = author
	return data
}

// Digital manifestation (e.g. audio file)

func NewDigitalManifestation(impl, id, name string, example interface{}, isManifestation bool, project interface{}, datePublished string, locationCreated, url interface{}) Data {
	var context interface{}
	if impl == IPLD {
		context = LinksIPLD(SCHEMA, COALAIP)
		example = LinkIPLD(example)
		project = LinkIPLD(project)
		locationCreated = LinkIPLD(locationCreated)
		url = LinkIPLD(url)
	}
	if impl == JSON {
		context = []string{SCHEMA, COALAIP}
	}
	data := make(Data)
	data["@context"] = context
	data["@type"] = "Manifestation"
	if id != "" {
		data["@id"] = id
	}
	data["name"] = name
	data["exampleOfWork"] = example
	data["isManifestation"] = isManifestation
	data["isPartOf"] = project
	data["datePublished"] = datePublished
	data["locationCreated"] = locationCreated
	data["url"] = url
	return data
}

// Digital Fingerprint

func NewDigitalFingerprint(impl, id, creativeWork, fingerprint string) Data {
	var context interface{}
	if impl == IPLD {
		context = LinkIPLD(COALAIP)
	}
	if impl == JSON {
		context = COALAIP
	}
	data := make(Data)
	data["@context"] = context
	data["@type"] = "DigitalFingerprint"
	if id != "" {
		data["@id"] = id
	}
	data["fingerprintOf"] = creativeWork
	data["fingerprint"] = fingerprint
	return data
}

// Right

func NewRight(impl, id string, usageType []string, territory interface{}, rightContext []string, exclusive bool, numberOfUses, percentageShares int, validFrom, validTo string, creativeWork, license interface{}) Data {
	var context interface{}
	if impl == IPLD {
		context = LinksIPLD(SCHEMA, COALAIP)
		territory = LinkIPLD(territory)
		creativeWork = LinkIPLD(creativeWork)
		license = LinkIPLD(license)
	}
	if impl == JSON {
		context = []string{SCHEMA, COALAIP}
	}
	data := make(Data)
	data["@context"] = context
	data["@type"] = "Right"
	if id != "" {
		data["@id"] = id
	}
	data["usageType"] = usageType
	data["territory"] = territory
	data["rightContext"] = rightContext
	data["exclusive"] = exclusive
	data["numberOfUses"] = numberOfUses
	data["percentageShares"] = percentageShares
	data["validFrom"] = validFrom
	data["validTo"] = validTo
	data["creation"] = creativeWork
	data["license"] = license
	return data
}

// Rights assignment

func NewRightsAssignment(impl, id string, creativeWork interface{}) Data {
	var context interface{}
	if impl == IPLD {
		context = LinkIPLD(COALAIP)
		creativeWork = LinkIPLD(creativeWork)
	}
	if impl == JSON {
		context = COALAIP
	}
	data := make(Data)
	data["@context"] = context
	data["@type"] = "RightsTransferAction"
	if id != "" {
		data["@id"] = id
	}
	data["transferContract"] = creativeWork
	return data
}

// Rights assertion

func NewRightsAssertion(impl, id string, asserter interface{}, assertionTruth bool, assertionSubject interface{}, _error, validFrom, validThrough string) Data {
	var context interface{}
	if impl == IPLD {
		context = LinksIPLD(SCHEMA, COALAIP)
		asserter = LinkIPLD(asserter)
		assertionSubject = LinkIPLD(assertionSubject)
	}
	if impl == JSON {
		context = []string{SCHEMA, COALAIP}
	}
	data := make(Data)
	data["@context"] = context
	data["@type"] = "ReviewAction"
	if id != "" {
		data["@id"] = id
	}
	data["asserter"] = asserter
	data["assertionTruth"] = assertionTruth
	data["assertionSubject"] = assertionSubject
	data["error"] = _error
	data["validFrom"] = validFrom
	data["validThrough"] = validThrough
	return data
}

//-------------------------------------------

// Album

func NewAlbum(impl, id string, byArtist interface{}, name string, track ...interface{}) Data {
	var context interface{}
	if impl == IPLD {
		context = LinkIPLD(SCHEMA)
		byArtist = LinkIPLD(byArtist)
		track = LinksIPLD(track...)
	}
	if impl == JSON {
		context = SCHEMA
	}
	data := make(Data)
	data["@context"] = context
	data["@type"] = "MusicAlbum"
	if id != "" {
		data["@id"] = id
	}
	data["byArtist"] = byArtist
	data["numTracks"] = len(track)
	data["name"] = name
	data["track"] = track
	return data
}

func NewComposition(impl, id string, composer, lyrics, recordedAs interface{}) Data {
	var context interface{}
	if impl == IPLD {
		context = LinkIPLD(SCHEMA)
		composer = LinkIPLD(composer)
		lyrics = LinkIPLD(lyrics)
		recordedAs = LinkIPLD(recordedAs)
	}
	if impl == JSON {
		context = SCHEMA
	}
	data := make(Data)
	data["@context"] = context
	data["@type"] = "MusicComposition"
	if id != "" {
		data["@id"] = id
	}
	data["composer"] = composer
	data["lyrics"] = lyrics
	data["recordedAs"] = recordedAs
	return data
}

func NewRecording(impl, id string, byArtist, inAlbum, recordingOf interface{}) Data {
	var context interface{}
	if impl == IPLD {
		context = LinkIPLD(SCHEMA)
		byArtist = LinkIPLD(byArtist)
		inAlbum = LinkIPLD(inAlbum)
		recordingOf = LinkIPLD(recordingOf)
	}
	if impl == JSON {
		context = SCHEMA
	}
	data := make(Data)
	data["@context"] = context
	data["@type"] = "MusicRecording"
	if id != "" {
		data["@id"] = id
	}
	data["byArtist"] = byArtist
	data["inAlbum"] = inAlbum
	data["recordingOf"] = recordingOf
	return data
}
