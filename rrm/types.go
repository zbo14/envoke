package rrm

import (
	"encoding/json"
)

// Context URLs
const (
	SCHEMA  = "http://schema.org/"
	COALAIP = "<coalaip placeholder>"
)

type uri string

// json-ld types for the coalaip rrm
// specs: github.com/COALAIP/specs/data-structure/README.md
// TODO: ipld

// Geo coordinates

func GeoData(lat, long string) []byte {
	data, _ := json.Marshal(struct {
		Context   string `json:"@context"`
		Type      string `json:"@type"`
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	}{
		Context:   SCHEMA,
		Type:      "GeoCoordinates",
		Latitude:  lat,
		Longitude: long,
	})
	return data
}

func GeoMap(lat, long string) map[string]interface{} {
	data := GeoData(lat, long)
	mp, _ := MapJSON(data)
	return mp
}

// Localizable place

func PlaceData(lat, long, name string) []byte {
	data, _ := json.Marshal(struct {
		Context string                 `json:"@context"`
		Type    string                 `json:"@type"`
		Geo     map[string]interface{} `json:"geo"`
		Name    string                 `json:"name"`
	}{
		Context: SCHEMA,
		Type:    "Place",
		Geo:     GeoMap(lat, long),
		Name:    name,
	})
	return data
}

// Person

func PersonData(id uri, givenName, familyName, birthDate, deathDate string) []byte {
	data, _ := json.Marshal(struct {
		Context    string `json: "@context"`
		Type       string `json:"@type"`
		Id         uri    `json:"@id"`
		GivenName  string `json:"givenName"`
		FamilyName string `json:"familyName"`
		BirthDate  string `json:"birthDate"`
		DeathDate  string `json:"deathDate"`
	}{
		Context:    SCHEMA,
		Type:       "Person",
		Id:         id,
		GivenName:  givenName,
		FamilyName: familyName,
		BirthDate:  birthDate,
		DeathDate:  deathDate,
	})
	return data
}

// Organization

func OrganizationData(id uri, name string, founder uri, members []uri) []byte {
	data, _ := json.Marshal(struct {
		Context string `json:"@context"`
		Type    string `json:"@type"`
		Id      uri    `json:"@id"`
		Name    string `json:"name"`
		Founder uri    `json:"founder"`
		Member  []uri  `json:"member"`
	}{
		Context: SCHEMA,
		Type:    "Organization",
		Id:      id,
		Name:    name,
		Founder: founder,
		Member:  members,
	})
	return data
}

// Creation

func WorkData(id uri, name string, author uri) []byte {
	data, _ := json.Marshal(struct {
		Context string `json:"@context"`
		Type    string `json:"@type"`
		Id      uri    `json:"@id"`
		Name    string `json:"name"`
		Author  uri    `json:"author"`
	}{
		Context: COALAIP,
		Type:    "Work",
		Id:      id,
		Name:    name,
		Author:  author,
	})
	return data
}

// Digital manifestation

func DigitalManifestationData(id uri, name string, example uri, isManifestation bool, project uri, datePublished string, locationCreated, url uri) []byte {
	data, _ := json.Marshal(struct {
		Context         string `json:"@context"`
		Type            string `json:"@type"`
		Id              uri    `json:"@id"`
		Name            string `json:"name"`
		ExampleOfWork   uri    `json:"exampleOfWork"`
		IsManifestation bool   `json:"isManifestation"`
		IsPartOf        uri    `json:"isPartOf"`
		DatePublished   string `json:"datePublished"`
		LocationCreated uri    `json:"locationCreated"`
		Url             uri    `json:"url"`
	}{
		Context:         COALAIP,
		Type:            "Manifestation",
		Id:              id,
		Name:            name,
		ExampleOfWork:   example,
		IsManifestation: isManifestation,
		IsPartOf:        project,
		DatePublished:   datePublished,
		LocationCreated: locationCreated,
		Url:             url,
	})
	return data
}

// Digital Fingerprint

func digitalFingerprint(id, creativeWork uri, fingerprint string) []byte {
	data, _ := json.Marshal(struct {
		Context       string `json:"@context"`
		Type          string `json:"@type"`
		Id            uri    `json:"@id"`
		FingerprintOf uri    `json:"fingerprintOf"`
		Fingerprint   string `json:"fingerprint"`
	}{
		Context:       COALAIP,
		Type:          "DigitalFingerprint",
		Id:            id,
		FingerprintOf: creativeWork,
		Fingerprint:   fingerprint,
	})
	return data
}

// Right

func RightData(id uri, usages []string, territory uri, rightContext []string, exclusive bool, numberOfUses, percentageShares int, validFrom, validTo string, creativeWork, license uri) []byte {
	data, _ := json.Marshal(struct {
		Context          string   `json:"@context"`
		Type             string   `json:"@type"`
		Id               uri      `json:"@id"`
		Usages           []string `json:"usages"`
		Territory        uri      `json:"territory"`
		RightContext     []string `json:"rightContext"`
		Exclusive        bool     `json:"exclusive"`
		NumberOfUses     int      `json:"numberOfUses"`
		PercentageShares int      `json:"share"`
		ValidFrom        string   `json:"validFrom"`
		ValidTo          string   `json:"validTo"`
		Creation         uri      `json:"creation"`
		License          uri      `json:"license"`
	}{
		Context:          COALAIP,
		Type:             "Right",
		Id:               id,
		Usages:           usages,
		Territory:        territory,
		RightContext:     rightContext,
		Exclusive:        exclusive,
		NumberOfUses:     numberOfUses,
		PercentageShares: percentageShares,
		ValidFrom:        validFrom,
		ValidTo:          validTo,
		Creation:         creativeWork,
		License:          license,
	})
	return data
}

// Rights assignment

func RightsAssignmentData(id, creativeWork uri) []byte {
	data, _ := json.Marshal(struct {
		Context          string `json:"@context"`
		Type             string `json:"@type"`
		Id               uri    `json:"@id"`
		TransferContract uri    `json:"transferContract"`
	}{
		Context:          COALAIP,
		Type:             "RightsTransferAction",
		Id:               id,
		TransferContract: creativeWork,
	})
	return data
}

// Rights assertion

func RightsAssertionData(id, asserter uri, assertionTruth bool, assertionSubject uri, _error, validFrom, validThrough string) []byte {
	data, _ := json.Marshal(struct {
		Context          string `json:"@context"`
		Type             string `json:"@type"`
		Id               uri    `json:"@id"`
		Asserter         uri    `json:"asserter"`
		AssertionTruth   bool   `json:"assertionTruth"`
		AssertionSubject uri    `json:"assertionSubject"`
		Error            string `json:"error"`
		ValidFrom        string `json:"validFrom"`
		ValidThrough     string `json:"validThrough"`
	}{
		Context:          SCHEMA,
		Type:             "ReviewAction",
		Id:               id,
		Asserter:         asserter,
		AssertionTruth:   assertionTruth,
		AssertionSubject: assertionSubject,
		Error:            _error,
		ValidFrom:        validFrom,
		ValidThrough:     validThrough,
	})
	return data
}
