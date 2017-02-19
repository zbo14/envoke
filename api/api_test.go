package api

import (
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/spec"
	"testing"
)

var path = "/Users/zach/Desktop/music/Allegro from Duet in C Major.mp3"

func TestApi(t *testing.T) {
	api := NewApi()
	output := MustOpenWriteFile("output.json")
	composer, err := api.Register("composer@gmail.com", "composer", "itsasecret", "www.composer.com")
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteJSON(output, composer); err != nil {
		panic(err)
	}
	composerId := composer.AgentId
	label, err := api.Register("label@gmail.com", "label", "shhh", "www.record_label.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, label)
	labelId := label.AgentId
	performer, err := api.Register("performer@gmail.com", "performer", "canyouguess", "www.bandcamp_page.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, performer)
	performerId := performer.AgentId
	producer, err := api.Register("producer@gmail.com", "producer", "1234", "www.soundcloud_page.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, producer)
	producerId := producer.AgentId
	publisher, err := api.Register("publisher@gmail.com", "publisher", "password", "www.publisher.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, publisher)
	publisherId := publisher.AgentId
	radio, err := api.Register("radio@gmail.com", "radio", "waves", "www.radio_station.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, radio)
	radioId := radio.AgentId
	if err = api.Login(composer.AgentId, composer.PrivKey); err != nil {
		t.Fatal(err)
	}
	composerRight, err := api.Right("20", composerId, "2020-01-01", "3030-01-01")
	if err != nil {
		t.Fatal(err)
	}
	publisherRight, err := api.Right("80", publisherId, "2020-01-01", "2030-01-01")
	if err != nil {
		t.Fatal(err)
	}
	composition, err := api.Composition(composerId, publisherId, []Data{composerRight, publisherRight}, "untitled")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, composition)
	compositionId := composition.GetStr("id")
	if err = api.Login(publisherId, publisher.PrivKey); err != nil {
		t.Fatal(err)
	}
	publishingLicense, err := api.PublishingLicense(compositionId, labelId, spec.LICENSE_TYPE_MECHANICAL, "2020-01-01", "2025-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, publishingLicense)
	publishingLicenseId := publishingLicense.GetStr("id")
	if err = api.Login(labelId, label.PrivKey); err != nil {
		t.Fatal(err)
	}
	labelRight, err := api.Right("50", labelId, "2020-01-01", "2030-01-01")
	if err != nil {
		t.Fatal(err)
	}
	performerRight, err := api.Right("30", performerId, "2020-01-01", "3030-01-01")
	if err != nil {
		t.Fatal(err)
	}
	producerRight, err := api.Right("20", producerId, "2020-01-01", "3030-01-01")
	if err != nil {
		t.Fatal(err)
	}
	file, err := OpenFile(path)
	if err != nil {
		t.Fatal(err)
	}
	recording, err := api.Recording(compositionId, file, labelId, performerId, producerId, publishingLicenseId, []Data{labelRight, performerRight, producerRight})
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, recording)
	recordingId := recording.GetStr("id")
	recordingLicense, err := api.RecordingLicense(radioId, spec.LICENSE_TYPE_MASTER, recordingId, "2020-01-01", "2022-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, recordingLicense)
}
