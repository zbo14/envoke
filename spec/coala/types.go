package coala

import (
	// . "github.com/zballs/envoke/util"
	"github.com/zballs/envoke/spec"
)

// Context URLs
const (
	SCHEMA  = "http://schema.org/"
	COALAIP = "<coalaip placeholder>"
)

// Coala IP spec
// github.com/COALAIP/specs/data-structure/README.md

func NewGeo(impl, lat, long string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = LinkIPLD(SCHEMA)
	}
	if impl == spec.JSON {
		context = SCHEMA
	}
	data := spec.Data{
		"@context":  context,
		"@type":     "GeoCoordinates",
		"latitude":  lat,
		"longitude": long,
	}
	return data
}

func NewPlace(impl, lat, long, name string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = LinkIPLD(SCHEMA)
	}
	if impl == spec.JSON {
		context = SCHEMA
	}
	data := spec.Data{
		"@context": context,
		"@type":    "Place",
		"geo":      NewGeo(impl, lat, long),
		"name":     name,
	}
	return data
}

func NewPerson(impl, id, givenName, familyName, birthDate string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = LinkIPLD(SCHEMA)
	}
	if impl == spec.JSON {
		context = SCHEMA
	}
	data := spec.Data{
		"@context":   context,
		"@type":      "Person",
		"givenName":  givenName,
		"familyName": familyName,
		"birthDate":  birthDate,
	}
	if id != "" {
		data["@id"] = id
	}
	return data
}

func NewOrganization(impl, id, email, name, _type string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = LinksIPLD(SCHEMA)
	}
	if impl == spec.JSON {
		context = SCHEMA
	}
	data := spec.Data{
		"@context":    context,
		"@type":       "Organization",
		"email":       email,
		"name":        name,
		"description": _type,
	}
	if id != "" {
		data["@id"] = id
	}
	return data
}

func NewWork(impl, id, name string, author interface{}) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = LinksIPLD(SCHEMA, COALAIP)
		author = LinkIPLD(author)
	}
	if impl == spec.JSON {
		context = []string{SCHEMA, COALAIP}
	}
	data := spec.Data{
		"@context": context,
		"@type":    "Work",
		"name":     name,
		"author":   author,
	}
	if id != "" {
		data["@id"] = id
	}
	return data
}

func NewDigitalManifestation(impl, id, name string, example interface{}, isManifestation bool, project interface{}, datePublished string, locationCreated, url interface{}) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = LinksIPLD(SCHEMA, COALAIP)
		example = LinkIPLD(example)
		project = LinkIPLD(project)
		locationCreated = LinkIPLD(locationCreated)
		url = LinkIPLD(url)
	}
	if impl == spec.JSON {
		context = []string{SCHEMA, COALAIP}
	}
	data := spec.Data{
		"@context":        context,
		"@type":           "Manifestation",
		"name":            name,
		"exampleOfWork":   example,
		"isManifestation": isManifestation,
		"isPartOf":        project,
		"datePublished":   datePublished,
		"locationCreated": locationCreated,
		"url":             url,
	}
	if id != "" {
		data["@id"] = id
	}
	return data
}

func NewDigitalFingerprint(impl, id, creativeWork, fingerprint string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = LinkIPLD(COALAIP)
	}
	if impl == spec.JSON {
		context = COALAIP
	}
	data := spec.Data{
		"@context":      context,
		"@type":         "DigitalFingerprint",
		"fingerprintOf": creativeWork,
		"fingerprint":   fingerprint,
	}
	if id != "" {
		data["@id"] = id
	}
	return data
}

func NewRight(impl, id string, usageType []string, territory interface{}, rightContext []string, exclusive bool, numberOfUses, percentageShares int, validFrom, validTo string, creativeWork, license interface{}) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = LinksIPLD(SCHEMA, COALAIP)
		territory = LinkIPLD(territory)
		creativeWork = LinkIPLD(creativeWork)
		license = LinkIPLD(license)
	}
	if impl == spec.JSON {
		context = []string{SCHEMA, COALAIP}
	}
	data := spec.Data{
		"@context":         context,
		"@type":            "Right",
		"usageType":        usageType,
		"territory":        territory,
		"rightContext":     rightContext,
		"exclusive":        exclusive,
		"numberOfUses":     numberOfUses,
		"percentageShares": percentageShares,
		"validFrom":        validFrom,
		"validTo":          validTo,
		"creation":         creativeWork,
		"license":          license,
	}
	if id != "" {
		data["@id"] = id
	}
	return data
}

func NewRightsAssignment(impl, id string, creativeWork interface{}) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = LinkIPLD(COALAIP)
		creativeWork = LinkIPLD(creativeWork)
	}
	if impl == spec.JSON {
		context = COALAIP
	}
	data := spec.Data{
		"@context":         context,
		"@type":            "RightsTransferAction",
		"transferContract": creativeWork,
	}
	if id != "" {
		data["@id"] = id
	}
	return data
}

func NewRightsAssertion(impl, id string, asserter interface{}, assertionTruth bool, assertionSubject interface{}, _error, validFrom, validThrough string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = LinksIPLD(SCHEMA, COALAIP)
		asserter = LinkIPLD(asserter)
		assertionSubject = LinkIPLD(assertionSubject)
	}
	if impl == spec.JSON {
		context = []string{SCHEMA, COALAIP}
	}
	data := spec.Data{
		"@context":         context,
		"@type":            "ReviewAction",
		"asserter":         asserter,
		"assertionTruth":   assertionTruth,
		"assertionSubject": assertionSubject,
		"error":            _error,
		"validFrom":        validFrom,
		"validThrough":     validThrough,
	}
	if id != "" {
		data["@id"] = id
	}
	return data
}

// Schema

func NewAlbum(impl, id string, byArtist interface{}, name string, release interface{}, track ...interface{}) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = LinkIPLD(SCHEMA)
		byArtist = LinkIPLD(byArtist)
		release = LinkIPLD(release)
		track = LinksIPLD(track...)
	}
	if impl == spec.JSON {
		context = SCHEMA
	}
	data := spec.Data{
		"@context":  context,
		"@type":     "MusicAlbum",
		"byArtist":  byArtist,
		"numTracks": len(track),
		"name":      name,
		"release":   release,
		"track":     track,
	}
	if id != "" {
		data["@id"] = id
	}
	return data
}

func NewComposition(impl, id string, composer, lyrics, recordedAs interface{}) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = LinkIPLD(SCHEMA)
		composer = LinkIPLD(composer)
		lyrics = LinkIPLD(lyrics)
		recordedAs = LinkIPLD(recordedAs)
	}
	if impl == spec.JSON {
		context = SCHEMA
	}
	data := spec.Data{
		"@context":   context,
		"@type":      "MusicComposition",
		"composer":   composer,
		"lyrics":     lyrics,
		"recordedAs": recordedAs,
	}
	if id != "" {
		data["@id"] = id
	}
	return data
}

func NewRecording(impl, id string, byArtist, inAlbum, recordingOf interface{}) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = LinkIPLD(SCHEMA)
		byArtist = LinkIPLD(byArtist)
		inAlbum = LinkIPLD(inAlbum)
		recordingOf = LinkIPLD(recordingOf)
	}
	if impl == spec.JSON {
		context = SCHEMA
	}
	data := spec.Data{
		"@context":    context,
		"@type":       "MusicRecording",
		"byArtist":    byArtist,
		"inAlbum":     inAlbum,
		"recordingOf": recordingOf,
	}
	if id != "" {
		data["@id"] = id
	}
	return data
}

func NewRelease(impl, id, datePublished string, recordLabel, publisher, releaseOf interface{}) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = LinkIPLD(SCHEMA)
		recordLabel = LinkIPLD(recordLabel)
		publisher = LinkIPLD(publisher)
		releaseOf = LinkIPLD(releaseOf)
	}
	if impl == spec.JSON {
		context = SCHEMA
	}
	data := spec.Data{
		"@context":      context,
		"@type":         "MusicRelease",
		"datePublished": datePublished,
		"recordLabel":   recordLabel,
		"publisher":     publisher,
		"releaseOf":     releaseOf,
	}
	if id != "" {
		data["@id"] = id
	}
	return data
}
