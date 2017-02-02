package coala

import . "github.com/zballs/envoke/util"

// Context URLs
const (
	SCHEMA  = "http://schema.org/"
	COALAIP = "<coalaip placeholder>"

	IPLD = "ipld"
	JSON = "json_ld"

	LC_REGEX = `^LC-\d{4,5}$`
)

type Data map[string]interface{}

// Coala IP specs: github.com/COALAIP/specs/data-structure/README.md

func NewGeo(impl, lat, long string) Data {
	var context interface{}
	if impl == IPLD {
		context = LinkIPLD(SCHEMA)
	}
	if impl == JSON {
		context = SCHEMA
	}
	data := Data{
		"@context":  context,
		"@type":     "GeoCoordinates",
		"latitude":  lat,
		"longitude": long,
	}
	return data
}

func NewPlace(impl, lat, long, name string) Data {
	var context interface{}
	if impl == IPLD {
		context = LinkIPLD(SCHEMA)
	}
	if impl == JSON {
		context = SCHEMA
	}
	data := Data{
		"@context": context,
		"@type":    "Place",
		"geo":      NewGeo(impl, lat, long),
		"name":     name,
	}
	return data
}

func NewPerson(impl, id, givenName, familyName, birthDate string) Data {
	var context interface{}
	if impl == IPLD {
		context = LinkIPLD(SCHEMA)
	}
	if impl == JSON {
		context = SCHEMA
	}
	data := Data{
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

func NewOrganization(impl, id, email, name, _type string) Data {
	var context interface{}
	if impl == IPLD {
		context = LinksIPLD(SCHEMA)
	}
	if impl == JSON {
		context = SCHEMA
	}
	data := Data{
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

func NewWork(impl, id, name string, author interface{}) Data {
	var context interface{}
	if impl == IPLD {
		context = LinksIPLD(SCHEMA, COALAIP)
		author = LinkIPLD(author)
	}
	if impl == JSON {
		context = []string{SCHEMA, COALAIP}
	}
	data := Data{
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
	data := Data{
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

func NewDigitalFingerprint(impl, id, creativeWork, fingerprint string) Data {
	var context interface{}
	if impl == IPLD {
		context = LinkIPLD(COALAIP)
	}
	if impl == JSON {
		context = COALAIP
	}
	data := Data{
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
	data := Data{
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

func NewRightsAssignment(impl, id string, creativeWork interface{}) Data {
	var context interface{}
	if impl == IPLD {
		context = LinkIPLD(COALAIP)
		creativeWork = LinkIPLD(creativeWork)
	}
	if impl == JSON {
		context = COALAIP
	}
	data := Data{
		"@context":         context,
		"@type":            "RightsTransferAction",
		"transferContract": creativeWork,
	}
	if id != "" {
		data["@id"] = id
	}
	return data
}

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
	data := Data{
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

func NewAlbum(impl, id string, byArtist interface{}, name string, release interface{}, track ...interface{}) Data {
	var context interface{}
	if impl == IPLD {
		context = LinkIPLD(SCHEMA)
		byArtist = LinkIPLD(byArtist)
		release = LinkIPLD(release)
		track = LinksIPLD(track...)
	}
	if impl == JSON {
		context = SCHEMA
	}
	data := Data{
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
	data := Data{
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
	data := Data{
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

func NewRelease(impl, id, datePublished string, recordLabel, publisher, releaseOf interface{}) Data {
	var context interface{}
	if impl == IPLD {
		context = LinkIPLD(SCHEMA)
		recordLabel = LinkIPLD(recordLabel)
		publisher = LinkIPLD(publisher)
		releaseOf = LinkIPLD(releaseOf)
	}
	if impl == JSON {
		context = SCHEMA
	}
	data := Data{
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

// Music Ontology (json-ld for now)

// should we just pass URIs?

var CONTEXT = Data{
	"mo":    "http://purl.org/ontology/mo/",
	"dc":    "http://purl.org/dc/elements/1.1/",
	"xsd":   "http://www.w3.org/2001/XMLSchema#",
	"tl":    "http://purl.org/NET/c4dm/timeline.owl#",
	"event": "http://purl.org/NET/c4dm/event.owl#",
	"foaf":  "http://xmlns.com/foaf/0.1/",
	"rdfs":  "http://www.w3.org/2000/01/rdf-schema#",
}

func NewTrack(id, title, artistId string) Data {
	return Data{
		"@context": CONTEXT,
		"@type":    "mo:Track",
		"dc:title": title,
		"foaf:maker": Data{
			"@id":   artistId,
			"@type": "mo:MusicArtist",
			// need artist name?
		},
	}
}

func NewArtist(id, name string) Data {
	return Data{
		"@context":  CONTEXT,
		"@id":       id,
		"@type":     "mo:MusicArtist",
		"foaf:name": name,
	}
}

func NewPublisher(id, name string) Data {
	return Data{
		"@context":  CONTEXT,
		"@id":       id,
		"@type":     "foaf:Organization",
		"foaf:name": name,
	}
}

func NewLabel(id, lc string) Data {
	if !MatchString(LC_REGEX, lc) {
		panic("Label code does not match regex pattern")
	}
	return Data{
		"@context": CONTEXT,
		"@id":      id,
		"@type":    "mo:Label",
		"mo:lc":    lc,
	}
}

func NewRecord(id string, number int, publisherId string, trackIds []string) Data {
	if number <= 0 {
		panic("Record number must be positive")
	}
	tracks := make([]Data, len(trackIds))
	for i, trackId := range trackIds {
		tracks[i] = Data{
			"@id":   trackId,
			"@type": "mo:Track",
		}
	}
	return Data{
		"@context":         CONTEXT,
		"@id":              id,
		"@type":            "mo:Record",
		"mo:record_number": number,
		"mo:publisher": Data{
			"@id":   publisherId,
			"@type": "foaf:Agent",
		},
		"mo:track": tracks,
	}
}
