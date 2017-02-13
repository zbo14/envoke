package api

import (
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/spec"
	"net/url"
	"testing"
)

var path = ""

func TestApi(t *testing.T) {
	api := NewApi()
	// Register Publisher
	values := make(url.Values)
	values.Set("email", "publisher@gmail.com")
	values.Set("name", "publisher_name")
	values.Set("password", "canyouguess")
	values.Set("type", spec.PUBLISHER)
	registerPublisher, err := api.Register(values)
	if err != nil {
		t.Fatal(err)
	}
	publisherId := registerPublisher.AgentId
	PrintJSON(registerPublisher)
	// Register Artist
	values.Set("email", "artist@gmail.com")
	values.Set("name", "artist_name")
	values.Set("password", "itsasecret")
	values.Set("type", spec.ARTIST)
	registerArtist, err := api.Register(values)
	if err != nil {
		t.Fatal(err)
	}
	PrintJSON(registerArtist)
	// Login Artist
	artistId := registerArtist.AgentId
	privstr := registerArtist.PrivKey
	loginArtist, err := api.Login(artistId, privstr, spec.ARTIST)
	if err != nil {
		t.Fatal(err)
	}
	PrintJSON(loginArtist)
	// Track
	file, err := OpenFile(path)
	if err != nil {
		t.Fatal(err)
	}
	trackMessage, err := api.Track("", file, 0, publisherId)
	if err != nil {
		t.Fatal(err)
	}
	trackId := trackMessage.TrackId
	PrintJSON(trackMessage)
	// Right
	values.Set("context", "commercialuse")
	values.Set("issuer_id", artistId)
	values.Set("music_id", trackId)
	values.Set("recipient_id", publisherId)
	values.Set("usage", "copy,sell")
	values.Set("valid_from", DateStr(2018, 1, 1))
	values.Set("valid_to", DateStr(2020, 1, 1))
	rightMessage, err := api.Right(values)
	if err != nil {
		t.Fatal(err)
	}
	rightId := rightMessage.RightId
	PrintJSON(rightMessage)
	// Login Publisher
	privstr = registerPublisher.PrivKey
	loginPublisher, err := api.Login(publisherId, privstr, spec.PUBLISHER)
	if err != nil {
		t.Fatal(err)
	}
	Println(loginPublisher)
	// Sign
	signMessage, err := api.Sign(rightId)
	if err != nil {
		t.Fatal(err)
	}
	signatureId := signMessage.SignatureId
	PrintJSON(signMessage)
	// Verify
	verifyMessage, err := api.Verify(signatureId)
	if err != nil {
		t.Fatal(err)
	}
	PrintJSON(verifyMessage)
}
