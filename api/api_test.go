package api

import (
	"testing"

	"github.com/zbo14/envoke/bigchain"
	. "github.com/zbo14/envoke/common"
	// ld "github.com/zbo14/envoke/linked_data"
)

var path = "/Users/zach/Desktop/music/Allegro from Duet in C Major.mp3"

func TestApi(t *testing.T) {
	api := NewApi()
	output := MustOpenWriteFile("output.json")
	composer, err := api.Register("composer@email.com", "composer", "itsasecret", "www.composer.com")
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteJSON(output, composer); err != nil {
		panic(err)
	}
	composerId := composer.AgentId
	label, err := api.Register("label@email.com", "label", "shhh", "www.record_label.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, label)
	labelId := label.AgentId
	performer, err := api.Register("performer@email.com", "performer", "canyouguess", "www.bandcamp_page.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, performer)
	performerId := performer.AgentId
	producer, err := api.Register("producer@email.com", "producer", "1234", "www.soundcloud_page.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, producer)
	producerId := producer.AgentId
	publisher, err := api.Register("publisher@email.com", "publisher", "password", "www.publisher.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, publisher)
	publisherId := publisher.AgentId
	radio, err := api.Register("radio@email.com", "radio", "waves", "www.radio_station.com")
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
	compositionId := bigchain.GetId(composition)
	composerRight, err := api.CompositionRight(compositionId, 20, composerId, []string{"GB", "US"}, "2020-01-01", "3000-01-01")
	if err != nil {
		t.Fatal(err)
	}
	composerRightId := bigchain.GetId(composerRight)
	publisherRight, err := api.CompositionRight(compositionId, 80, publisherId, []string{"GB", "US"}, "2020-01-01", "3000-01-01")
	if err != nil {
		t.Fatal(err)
	}
	publisherRightId := bigchain.GetId(publisherRight)
	if err = api.Login(publisherId, publisher.PrivKey); err != nil {
		t.Fatal(err)
	}
	publication, err := api.Publish(compositionId, []string{composerRightId, publisherRightId})
	if err != nil {
		t.Fatal(err)
	}
	publicationId := bigchain.GetId(publication)
	mechanicalLicense, err := api.MechanicalLicense(labelId, publicationId, publisherRightId, []string{"US"}, "2020-01-01", "2025-01-01")
	if err != nil {
		t.Fatal(err)
	}
	mechanicalLicenseId := bigchain.GetId(mechanicalLicense)
	file, err := OpenFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = api.Login(performerId, performer.PrivKey); err != nil {
		t.Fatal(err)
	}
	recording, err := api.Record("", file, "US-S1Z-99-00001", labelId, performerId, producerId, publicationId)
	if err != nil {
		t.Fatal(err)
	}
	recordingId := bigchain.GetId(recording)
	performerRight, err := api.RecordingRight(20, recordingId, performerId, []string{"GB", "US"}, "2020-01-01", "2080-01-01")
	if err != nil {
		t.Fatal(err)
	}
	performerRightId := bigchain.GetId(performerRight)
	producerRight, err := api.RecordingRight(10, recordingId, producerId, []string{"GB", "US"}, "2020-01-01", "2080-01-01")
	if err != nil {
		t.Fatal(err)
	}
	producerRightId := bigchain.GetId(producerRight)
	labelRight, err := api.RecordingRight(70, recordingId, labelId, []string{"GB", "US"}, "2020-01-01", "2080-01-01")
	if err != nil {
		t.Fatal(err)
	}
	labelRightId := bigchain.GetId(labelRight)
	if err = api.Login(labelId, label.PrivKey); err != nil {
		t.Fatal(err)
	}
	release, err := api.Release(mechanicalLicenseId, recordingId, []string{labelRightId, performerRightId, producerRightId})
	if err != nil {
		t.Fatal(err)
	}
	releaseId := bigchain.GetId(release)
	releaseLicense, err := api.MasterLicense(radioId, releaseId, labelRightId, []string{"US"}, "2020-01-01", "2022-01-01")
	if err != nil {
		t.Fatal(err)
	}
	releaseLicenseId := bigchain.GetId(releaseLicense)
	t.Log(releaseLicenseId)
	if err = api.Login(composerId, composer.PrivKey); err != nil {
		t.Fatal(err)
	}
	compositionRightTransfer1, err := api.TransferCompositionRight(0, 10, publicationId, publisherId, composerRightId)
	if err != nil {
		t.Fatal(err)
	}
	compositionRightTransfer1Id := bigchain.GetId(compositionRightTransfer1)
	if err = api.Login(publisherId, publisher.PrivKey); err != nil {
		t.Fatal(err)
	}
	compositionRightTransfer2, err := api.TransferCompositionRight(1, 5, publicationId, composerId, compositionRightTransfer1Id)
	if err != nil {
		t.Fatal(err)
	}
	compositionRightTransfer2Id := bigchain.GetId(compositionRightTransfer2)
	t.Log(compositionRightTransfer2Id)
	if err = api.Login(performerId, performer.PrivKey); err != nil {
		t.Fatal(err)
	}
	recordingRightTransfer1, err := api.TransferRecordingRight(0, 10, labelId, releaseId, performerRightId)
	if err != nil {
		t.Fatal(err)
	}
	recordingRightTransfer1Id := bigchain.GetId(recordingRightTransfer1)
	if err = api.Login(labelId, label.PrivKey); err != nil {
		t.Fatal(err)
	}
	recordingRightTransfer2, err := api.TransferRecordingRight(1, 5, performerId, releaseId, recordingRightTransfer1Id)
	if err != nil {
		t.Fatal(err)
	}
	recordingRightTransfer2Id := bigchain.GetId(recordingRightTransfer2)
	t.Log(recordingRightTransfer2Id)
}
