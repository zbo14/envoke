package api

import (
	"net/url"
	"os"
	"testing"

	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/spec"
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
		t.Error(err)
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
		t.Error(err)
	}
	PrintJSON(registerArtist)
	// Login Artist
	artistId := registerArtist.AgentId
	privstr := registerArtist.PrivKey
	loginArtist, err := api.Login(artistId, privstr, spec.ARTIST)
	if err != nil {
		t.Error(err)
	}
	PrintJSON(loginArtist)
	// Track
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	trackMessage, err := api.Track("", file, 0, publisherId)
	if err != nil {
		t.Error(err)
	}
	trackId := trackMessage.TrackId
	PrintJSON(trackMessage)
	// Login Publisher
	privstr = registerPublisher.PrivKey
	loginPublisher, err := api.Login(publisherId, privstr, spec.PUBLISHER)
	if err != nil {
		t.Error(err)
	}
	PrintJSON(loginPublisher)
	// Sign
	signMessage, err := api.Sign(trackId)
	if err != nil {
		t.Error(err)
	}
	signatureId := signMessage.SignatureId
	PrintJSON(signMessage)
	// Verify
	verifyMessage := api.Verify(signatureId)
	PrintJSON(verifyMessage)
}
