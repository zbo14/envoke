package coala

import (
	// . "github.com/zbo14/envoke/util"
	"github.com/zbo14/envoke/spec"
)

// Context URLs
const (
	SCHEMA  = "http://schema.org/"
	COALAIP = "<coalaip placeholder>"

	LABEL     = "label"
	PUBLISHER = "publisher"
	TRACK     = "MusicRecording"
)

// Coala IP spec
// github.com/COALAIP/specs/data-structure/README.md

func NewGeo(impl, lat, long string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(SCHEMA)
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
		context = spec.LinkIPLD(SCHEMA)
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

func NewPerson(impl, id, givenName, familyName, birthDate, deathDate string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(SCHEMA)
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
		"deathDate":  deathDate,
	}
	if impl == spec.JSON {
		if id != "" {
			data["@id"] = id
		}
	}
	return data
}

func NewOrganization(impl, id, email, login, name, _type string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinksIPLD(SCHEMA)
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
		"url":         login,
	}
	if impl == spec.JSON {
		if id != "" {
			data["@id"] = id
		}
	}
	return data
}

func NewCreativeWork(impl, id, name string, author interface{}) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinksIPLD(SCHEMA, COALAIP)
		author = spec.LinkIPLD(author)
	}
	if impl == spec.JSON {
		context = []string{SCHEMA, COALAIP}
	}
	data := spec.Data{
		"@context": context,
		"@type":    "CreativeWork",
		"name":     name,
		"author":   author,
	}
	if impl == spec.JSON {
		if id != "" {
			data["@id"] = id
		}
	}
	return data
}

func NewDigitalManifestation(impl, id, _type, name string, example interface{}, isManifestation bool, isPartOf interface{}, datePublished, locationCreated string, url interface{}) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinksIPLD(SCHEMA, COALAIP)
		example = spec.LinkIPLD(example)
		isPartOf = spec.LinkIPLD(isPartOf)
		url = spec.LinkIPLD(url)
	}
	if impl == spec.JSON {
		context = []string{SCHEMA, COALAIP}
	}
	data := spec.Data{
		"@context":        context,
		"@type":           _type,
		"name":            name,
		"exampleOfWork":   example,
		"isManifestation": isManifestation,
		"isPartOf":        isPartOf,
		"datePublished":   datePublished,
		"locationCreated": spec.Data{
			"@type": "Place",
			"name":  locationCreated,
		},
		"url": url,
	}
	if impl == spec.JSON {
		if id != "" {
			data["@id"] = id
		}
	}
	return data
}

func NewDigitalFingerprint(impl, id, creativeWork, fingerprint string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(COALAIP)
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
	if impl == spec.JSON {
		if id != "" {
			data["@id"] = id
		}
	}
	return data
}

func NewRight(impl, id string, usageType []string, territory interface{}, rightContext []string, exclusive bool, numberOfUses, percentageShares int, validFrom, validTo string, creativeWork, license interface{}) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinksIPLD(SCHEMA, COALAIP)
		territory = spec.LinkIPLD(territory)
		creativeWork = spec.LinkIPLD(creativeWork)
		license = spec.LinkIPLD(license)
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
	if impl == spec.JSON {
		if id != "" {
			data["@id"] = id
		}
	}
	return data
}

func NewRightsAssignment(impl, id string, creativeWork interface{}) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(COALAIP)
		creativeWork = spec.LinkIPLD(creativeWork)
	}
	if impl == spec.JSON {
		context = COALAIP
	}
	data := spec.Data{
		"@context":         context,
		"@type":            "RightsTransferAction",
		"transferContract": creativeWork,
	}
	if impl == spec.JSON {
		if id != "" {
			data["@id"] = id
		}
	}
	return data
}

func NewRightsAssertion(impl, id string, asserter interface{}, assertionTruth bool, assertionSubject interface{}, _error, validFrom, validThrough string) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinksIPLD(SCHEMA, COALAIP)
		asserter = spec.LinkIPLD(asserter)
		assertionSubject = spec.LinkIPLD(assertionSubject)
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
	if impl == spec.JSON {
		if id != "" {
			data["@id"] = id
		}
	}
	return data
}

//------------------------------------------------------------

func NewAlbum(impl, id, name string, artist interface{}) spec.Data {
	return NewCreativeWork(impl, id, name, artist)
}

func NewTrack(impl, id, name string, example interface{}, album, artist interface{}, datePublished, locationCreated string, url interface{}) spec.Data {
	// TODO: isrc code
	// what should example be?
	data := NewDigitalManifestation(impl, id, TRACK, name, example, true, album, datePublished, locationCreated, url)
	data["byArtist"] = artist
	return data
}

func NewArtist(email, name, openId, partnerId string) spec.Data {
	return spec.Data{
		"@type": "MusicGroup",
		"email": email,
		"name":  "name",
		"memberOf": spec.Data{
			"@id":   partnerId,
			"@type": "Organization",
		},
		"sameAs": openId,
	}
}

func NewLabel(impl, id, email, login, name string) spec.Data {
	return NewOrganization(impl, id, email, login, name, LABEL)
}

func NewPublisher(impl, id, email, login, name string) spec.Data {
	return NewOrganization(impl, id, email, login, name, PUBLISHER)
}

/*
func NewArtist(impl, id, email, name string, members []string, partnerId string) spec.Data {
	// TODO: add roles
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinksIPLD(SCHEMA, COALAIP)
	}
	if impl == spec.JSON {
		context = []string{SCHEMA, COALAIP}
	}
	member := make([]spec.Data, len(members))
	for i, name := range members {
		member[i] = spec.Data{
			"@type": "Person",
			"name":  name,
		}
	}
	data := spec.Data{
		"@context": context,
		"@id":      "MusicGroup",
		"email":    email,
		"name":     name,
		"member":   member,
		"memberOf": spec.Data{
			"@id":   partnerId,
			"@type": "Organization",
		},
	}
	if id != "" {
		data["@id"] = id
	}
	return data
}
*/

/*
// Schema

func NewAlbum(impl, id string, byArtist interface{}, name string, release interface{}, track ...interface{}) spec.Data {
	var context interface{}
	if impl == spec.IPLD {
		context = spec.LinkIPLD(SCHEMA)
		byArtist = spec.LinkIPLD(byArtist)
		release = spec.LinkIPLD(release)
		track = spec.LinksIPLD(track...)
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
		context = spec.LinkIPLD(SCHEMA)
		composer = spec.LinkIPLD(composer)
		lyrics = spec.LinkIPLD(lyrics)
		recordedAs = spec.LinkIPLD(recordedAs)
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
		context = spec.LinkIPLD(SCHEMA)
		byArtist = spec.LinkIPLD(byArtist)
		inAlbum = spec.LinkIPLD(inAlbum)
		recordingOf = spec.LinkIPLD(recordingOf)
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
		context = spec.LinkIPLD(SCHEMA)
		recordLabel = spec.LinkIPLD(recordLabel)
		publisher = spec.LinkIPLD(publisher)
		releaseOf = spec.LinkIPLD(releaseOf)
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
*/
