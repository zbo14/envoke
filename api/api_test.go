package api

import (
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/spec"
	"net/url"
	"testing"
)

var path = "/Users/zach/Desktop/music/Allegro from Duet in C Major.mp3"

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
	// Register Record Label
	values.Set("email", "label@gmail.com")
	values.Set("name", "label_name")
	values.Set("password", "shhhh")
	values.Set("type", spec.LABEL)
	registerLabel, err := api.Register(values)
	if err != nil {
		t.Fatal(err)
	}
	labelId := registerLabel.AgentId
	PrintJSON(registerLabel)
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
	// New track by artist
	file, err := OpenFile(path)
	if err != nil {
		t.Fatal(err)
	}
	trackMessage, err := api.Track("", file, labelId, 0, publisherId)
	if err != nil {
		t.Fatal(err)
	}
	trackId := trackMessage.TrackId
	PrintJSON(trackMessage)
	// Artist issues right w/ 50% shares to publisher
	values.Set("context", "commercial_use")
	values.Set("issuer_id", artistId)
	values.Set("issuer_type", spec.ARTIST)
	values.Set("music_id", trackId)
	values.Set("percentage_shares", "50")
	values.Set("recipient_id", publisherId)
	values.Set("usage", "copy,play")
	values.Set("valid_from", DateStr(2018, 1, 1))
	values.Set("valid_to", DateStr(2020, 1, 1))
	publisherRightMessage, err := api.Right(values)
	if err != nil {
		t.Fatal(err)
	}
	publisherRightId := publisherRightMessage.RightId
	PrintJSON(publisherRightMessage)
	// Artist issues right w/ 20% shares to label
	values.Set("percentage_shares", "20")
	values.Set("recipient_id", labelId)
	labelRightMessage, err := api.Right(values)
	if err != nil {
		t.Fatal(err)
	}
	labelRightId := labelRightMessage.RightId
	PrintJSON(labelRightMessage)
	// Login Publisher
	privstr = registerPublisher.PrivKey
	loginPublisher, err := api.Login(publisherId, privstr, spec.PUBLISHER)
	if err != nil {
		t.Fatal(err)
	}
	Println(loginPublisher)
	// Rights
	rightsMessage, err := api.Rights(trackId, publisherRightId, labelRightId)
	if err != nil {
		t.Fatal(err)
	}
	rightId := rightsMessage.RightsId
	PrintJSON(rightsMessage)
	// Verify
	verifyMessage := api.Verify(rightId)
	PrintJSON(verifyMessage)
}
