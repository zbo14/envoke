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
	publisher1, err := api.Register("publisher1@gmail.com", "publisher1", "password1", "www.publisher1.com")
	if err != nil {
		t.Fatal(err)
	}
	publisher1Id := publisher1.AgentId
	publisher2, err := api.Register("publisher2@gmail.com", "publisher2", "password2", "www.publisher2.com")
	if err != nil {
		t.Fatal(err)
	}
	publisher2Id := publisher2.AgentId
	label, err := api.Register("label@gmail.com", "label", "shhh", "www.record_label.com")
	if err != nil {
		t.Fatal(err)
	}
	labelId := label.AgentId
	artist, err := api.Register("artist@gmail.com", "artist", "itsasecret", "www.soundcloud_profile.com")
	if err != nil {
		t.Fatal(err)
	}
	radio, err := api.Register("radio@gmail.com", "radio", "waves", "www.radio_station.com")
	if err != nil {
		t.Fatal(err)
	}
	radioId := radio.AgentId
	artistId := artist.AgentId
	privstr := artist.PrivKey
	if err = api.Login(artistId, privstr); err != nil {
		t.Fatal(err)
	}
	composerRight, err := api.Right("30", artistId, "2018-01-01", "2020-01-01")
	if err != nil {
		t.Fatal(err)
	}
	publisher1Right, err := api.Right("70", publisher1Id, "2018-01-01", "2020-01-01")
	if err != nil {
		t.Fatal(err)
	}
	composition, err := api.Composition(artistId, publisher1Id, []Data{composerRight, publisher1Right}, "untitled")
	if err != nil {
		t.Fatal(err)
	}
	compositionId := composition.GetStr("id")
	labelRight, err := api.Right("60", labelId, "2018-01-01", "2022-01-01")
	if err != nil {
		t.Fatal(err)
	}
	performerRight, err := api.Right("40", artistId, "2018-01-01", "2023-01-01")
	if err != nil {
		t.Fatal(err)
	}
	file, err := OpenFile(path)
	if err != nil {
		t.Fatal(err)
	}
	recording, err := api.Recording(compositionId, file, labelId, artistId, artistId, "", []Data{labelRight, performerRight})
	if err != nil {
		t.Fatal(err)
	}
	recordingId := recording.GetStr("id")
	privstr = label.PrivKey
	if err = api.Login(labelId, privstr); err != nil {
		t.Fatal(err)
	}
	license, err := api.RecordingLicense(radioId, spec.LICENSE_TYPE_MASTER, recordingId, "2018-01-01", "2019-01-01")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = ld.ValidateRecordingLicenseById(license.GetStr("id")); err != nil {
		t.Fatal(err)
	}
	privstr = publisher1.PrivKey
	if err = api.Login(publisher1Id, privstr); err != nil {
		t.Fatal(err)
	}
	license, err = api.PublishingLicense(compositionId, publisher2Id, spec.LICENSE_TYPE_MECHANICAL, "2018-01-01", "2025-01-01")
	if _, err = ld.ValidatePublishingLicenseById(license.GetStr("id")); err != nil {
		t.Fatal(err)
	}
}
