package api

import (
	"testing"

	"github.com/zbo14/envoke/bigchain"
	. "github.com/zbo14/envoke/common"
	// ld "github.com/zbo14/envoke/linked_data"
)

func GetAgentPrivateKey(data Data) string {
	return data.GetData("agent").GetStr("privateKey")
}

func TestApi(t *testing.T) {
	api := NewApi()
	output := MustOpenWriteFile("output.json")
	composer, err := api.Register("composer@email.com", "composer", "itsasecret", "www.composer.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, composer)
	composerId := bigchain.GetId(composer)
	composerPriv := GetAgentPrivateKey(composer)
	label, err := api.Register("label@email.com", "label", "shhh", "www.record_label.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, label)
	labelId := bigchain.GetId(label)
	labelPriv := GetAgentPrivateKey(label)
	performer, err := api.Register("performer@email.com", "performer", "canyouguess", "www.bandcamp_page.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, performer)
	performerId := bigchain.GetId(performer)
	performerPriv := GetAgentPrivateKey(performer)
	producer, err := api.Register("producer@email.com", "producer", "1234", "www.soundcloud_page.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, producer)
	producerId := bigchain.GetId(producer)
	publisher, err := api.Register("publisher@email.com", "publisher", "password", "www.publisher.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, publisher)
	publisherId := bigchain.GetId(publisher)
	publisherPriv := GetAgentPrivateKey(publisher)
	radio, err := api.Register("radio@email.com", "radio", "waves", "www.radio_station.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, radio)
	radioId := bigchain.GetId(radio)
	if err = api.Login(composerId, composerPriv); err != nil {
		t.Fatal(err)
	}
	composition, err := api.Compose("B3107S", "T-034.524.680-1", publisherId, "untitled")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, composition)
	compositionId := bigchain.GetId(composition)
	composerRight, err := api.CompositionRight(compositionId, 20, composerId, []string{"GB", "US"}, "2020-01-01", "3000-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, composerRight)
	composerRightId := bigchain.GetId(composerRight)
	publisherRight, err := api.CompositionRight(compositionId, 80, publisherId, []string{"GB", "US"}, "2020-01-01", "3000-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, publisherRight)
	publisherRightId := bigchain.GetId(publisherRight)
	if err = api.Login(publisherId, publisherPriv); err != nil {
		t.Fatal(err)
	}
	publication, err := api.Publish(compositionId, []string{composerRightId, publisherRightId})
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, publication)
	publicationId := bigchain.GetId(publication)
	mechanicalLicense, err := api.MechanicalLicense(labelId, publicationId, publisherRightId, []string{"US"}, "2020-01-01", "2025-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, mechanicalLicense)
	mechanicalLicenseId := bigchain.GetId(mechanicalLicense)
	file, err := OpenFile(Getenv("PATH_TO_AUDIO_FILE"))
	if err != nil {
		t.Fatal(err)
	}
	if err = api.Login(performerId, performerPriv); err != nil {
		t.Fatal(err)
	}
	recording, err := api.Record("", file, "US-S1Z-99-00001", labelId, performerId, producerId, publicationId)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, recording)
	recordingId := bigchain.GetId(recording)
	performerRight, err := api.RecordingRight(20, recordingId, performerId, []string{"GB", "US"}, "2020-01-01", "2080-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, performerRight)
	performerRightId := bigchain.GetId(performerRight)
	producerRight, err := api.RecordingRight(10, recordingId, producerId, []string{"GB", "US"}, "2020-01-01", "2080-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, producerRight)
	producerRightId := bigchain.GetId(producerRight)
	labelRight, err := api.RecordingRight(70, recordingId, labelId, []string{"GB", "US"}, "2020-01-01", "2080-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, labelRight)
	labelRightId := bigchain.GetId(labelRight)
	if err = api.Login(labelId, labelPriv); err != nil {
		t.Fatal(err)
	}
	release, err := api.Release(mechanicalLicenseId, recordingId, []string{labelRightId, performerRightId, producerRightId})
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, release)
	releaseId := bigchain.GetId(release)
	masterLicense, err := api.MasterLicense(radioId, releaseId, labelRightId, []string{"US"}, "2020-01-01", "2022-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, masterLicense)
	// masterLicenseId := bigchain.GetId(masterLicense)
	if err = api.Login(composerId, composerPriv); err != nil {
		t.Fatal(err)
	}
	compositionRightTransfer, err := api.TransferCompositionRight(publicationId, publisherId, 10, composerRightId, "")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, compositionRightTransfer)
	compositionRightTransferId := bigchain.GetId(compositionRightTransfer)
	if err = api.Login(publisherId, publisherPriv); err != nil {
		t.Fatal(err)
	}
	compositionRightTransfer, err = api.TransferCompositionRight(publicationId, composerId, 5, composerRightId, compositionRightTransferId)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, compositionRightTransfer)
	if err = api.Login(performerId, performerPriv); err != nil {
		t.Fatal(err)
	}
	recordingRightTransfer, err := api.TransferRecordingRight(labelId, 10, releaseId, performerRightId, "")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, recordingRightTransfer)
	recordingRightTransferId := bigchain.GetId(recordingRightTransfer)
	if err = api.Login(labelId, labelPriv); err != nil {
		t.Fatal(err)
	}
	recordingRightTransfer, err = api.TransferRecordingRight(performerId, 5, releaseId, performerRightId, recordingRightTransferId)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, recordingRightTransfer)
}
