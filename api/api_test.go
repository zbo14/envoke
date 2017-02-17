package api

import (
	. "github.com/zbo14/envoke/common"
	ld "github.com/zbo14/envoke/linked_data"
	"github.com/zbo14/envoke/spec"
	"testing"
)

var path = "/Users/zach/Desktop/music/Allegro from Duet in C Major.mp3"

func TestApi(t *testing.T) {
	api := NewApi()
	// Register publisher
	registerPublisher, err := api.Register("publisher@gmail.com", "publisher", "password", "www.publisher.com")
	if err != nil {
		t.Fatal(err)
	}
	publisherId := registerPublisher.AgentId
	// Register record label
	registerLabel, err := api.Register("label@gmail.com", "label", "shhh", "www.record_label.com")
	if err != nil {
		t.Fatal(err)
	}
	labelId := registerLabel.AgentId
	// Register artist
	registerArtist, err := api.Register("artist@gmail.com", "artist", "itsasecret", "www.artist.com")
	if err != nil {
		t.Fatal(err)
	}
	// Register radio station
	registerRadio, err := api.Register("radio@gmail.com", "radio", "waves", "www.radio.com")
	if err != nil {
		t.Fatal(err)
	}
	radioId := registerRadio.AgentId
	// Login artist
	artistId := registerArtist.AgentId
	privstr := registerArtist.PrivKey
	if err := api.Login(artistId, privstr); err != nil {
		t.Fatal(err)
	}
	// Composition Rights
	composerRight, err := api.Right([]string{"commercial_use"}, false, "30", artistId, []string{"copy", "play"}, MustParseDateStr("2018-01-01"), MustParseDateStr("2020-01-01"))
	if err != nil {
		t.Fatal(err)
	}
	publisherRight, err := api.Right([]string{"commercial_use"}, false, "70", publisherId, []string{"copy", "play"}, MustParseDateStr("2018-01-01"), MustParseDateStr("2020-01-01"))
	if err != nil {
		t.Fatal(err)
	}
	// Composition
	composition, err := api.Composition(artistId, publisherId, []Data{composerRight, publisherRight}, "composition_title")
	if err != nil {
		t.Fatal(err)
	}
	// Recording Rights
	labelRight, err := api.Right([]string{"commercial_use"}, false, "60", labelId, []string{"copy", "play"}, MustParseDateStr("2018-01-01"), MustParseDateStr("2022-01-01"))
	if err != nil {
		t.Fatal(err)
	}
	performerRight, err := api.Right([]string{"commercial_use"}, false, "40", artistId, []string{"copy", "play"}, MustParseDateStr("2018-01-01"), MustParseDateStr("2023-01-01"))
	if err != nil {
		t.Fatal(err)
	}
	// Recording
	file, err := OpenFile(path)
	if err != nil {
		t.Fatal(err)
	}
	recording, err := api.Recording(composition.GetStr("id"), file, labelId, artistId, artistId, "", []Data{labelRight, performerRight})
	if err != nil {
		t.Fatal(err)
	}
	// Login label
	privstr = registerLabel.PrivKey
	if err := api.Login(labelId, privstr); err != nil {
		t.Fatal(err)
	}
	// Recording license
	license, err := api.RecordingLicense(radioId, spec.LICENSE_TYPE_MASTER, recording.GetStr("id"), MustParseDateStr("2018-01-01"), MustParseDateStr("2019-01-01"))
	if err != nil {
		t.Fatal(err)
	}
	// Verify
	if _, err = ld.ValidateRecordingLicenseById(license.GetStr("id")); err != nil {
		t.Fatal(err)
	}
}
