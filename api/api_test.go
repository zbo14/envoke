package api

import (
	"testing"

	. "github.com/zbo14/envoke/common"
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
	if err = api.Login(composerId, composer.PrivKey); err != nil {
		t.Fatal(err)
	}
	composition, err := api.Compose(composerId, "B3107S", "T-034.524.680-1", publisherId, "untitled")
	if err != nil {
		t.Fatal(err)
	}
	compositionId := composition.GetStr("id")
	composerRight, err := api.CompositionRight(compositionId, "20", []string{"GB", "US"}, "2020-01-01", "3000-01-01")
	if err != nil {
		t.Fatal(err)
	}
	composerRightId := composerRight.GetStr("id")
	if err = api.Login(publisherId, publisher.PrivKey); err != nil {
		t.Fatal(err)
	}
	publisherRight, err := api.CompositionRight(compositionId, "80", []string{"GB", "US"}, "2020-01-01", "3000-01-01")
	if err != nil {
		t.Fatal(err)
	}
	publisherRightId := publisherRight.GetStr("id")
	publication, err := api.Publish(compositionId, []string{composerRightId, publisherRightId})
	if err != nil {
		t.Fatal(err)
	}
	publicationId := publication.GetStr("id")
	mechanicalLicense, err := api.MechanicalLicense(labelId, publicationId, []string{"US"}, "2020-01-01", "2025-01-01")
	if err != nil {
		t.Fatal(err)
	}
	mechanicalLicenseId := mechanicalLicense.GetStr("id")
	file, err := OpenFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = api.Login(performerId, performer.PrivKey); err != nil {
		t.Fatal(err)
	}
	recording, err := api.Record(file, "US-S1Z-99-00001", labelId, performerId, producerId, publicationId)
	if err != nil {
		t.Fatal(err)
	}
	recordingId := recording.GetStr("id")
	performerRight, err := api.RecordingRight("20", recordingId, []string{"GB", "US"}, "2020-01-01", "2080-01-01")
	if err != nil {
		t.Fatal(err)
	}
	performerRightId := performerRight.GetStr("id")
	if err = api.Login(producerId, producer.PrivKey); err != nil {
		t.Fatal(err)
	}
	producerRight, err := api.RecordingRight("10", recordingId, []string{"GB", "US"}, "2020-01-01", "2080-01-01")
	if err != nil {
		t.Fatal(err)
	}
	producerRightId := producerRight.GetStr("id")
	if err = api.Login(labelId, label.PrivKey); err != nil {
		t.Fatal(err)
	}
	labelRight, err := api.RecordingRight("70", recordingId, []string{"GB", "US"}, "2020-01-01", "2080-01-01")
	if err != nil {
		t.Fatal(err)
	}
	labelRightId := labelRight.GetStr("id")
	release, err := api.Release(mechanicalLicenseId, recordingId, []string{labelRightId, performerRightId, producerRightId})
	if err != nil {
		t.Fatal(err)
	}
	releaseId := release.GetStr("id")
	releaseLicense, err := api.MasterLicense(radioId, releaseId, []string{"US"}, "2020-01-01", "2022-01-01")
	if err != nil {
		t.Fatal(err)
	}
	releaseLicenseId := releaseLicense.GetStr("id")
	t.Log(releaseLicenseId)
}
