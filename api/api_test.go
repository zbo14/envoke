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
	recordLabel, err := api.Register("record_label@email.com", "record_label", "shhh", "www.record_label.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, recordLabel)
	recordLabelId := bigchain.GetId(recordLabel)
	recordLabelPriv := GetAgentPrivateKey(recordLabel)
	radio, err := api.Register("radio@email.com", "radio", "waves", "www.radio.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, radio)
	radioId := bigchain.GetId(radio)
	if err = api.Login(composerId, composerPriv); err != nil {
		t.Fatal(err)
	}
	composition, err := api.Compose("B3107S", "123456789", "T-034.524.680-1", "ASCAP", "untitled")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, composition)
	compositionId := bigchain.GetId(composition)
	composerRight, err := api.CompositionRight(compositionId, composerId, 20, []string{"GB", "US"}, nil, "2020-01-01", "3000-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, composerRight)
	composerRightId := bigchain.GetId(composerRight)
	publisherRight, err := api.CompositionRight(compositionId, publisherId, 80, []string{"GB", "US"}, nil, "2020-01-01", "3000-01-01")
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
	mechanicalLicense, err := api.MechanicalLicense(publisherRightId, "", publicationId, recordLabelId, []string{"US"}, nil, "2020-01-01", "2025-01-01")
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
	recording, err := api.Record("", file, "US-S1Z-99-00001", performerId, producerId, publicationId)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, recording)
	recordingId := bigchain.GetId(recording)
	performerRight, err := api.RecordingRight(performerId, 20, recordingId, []string{"GB", "US"}, nil, "2020-01-01", "2080-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, performerRight)
	performerRightId := bigchain.GetId(performerRight)
	producerRight, err := api.RecordingRight(producerId, 10, recordingId, []string{"GB", "US"}, nil, "2020-01-01", "2080-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, producerRight)
	producerRightId := bigchain.GetId(producerRight)
	recordLabelRight, err := api.RecordingRight(recordLabelId, 70, recordingId, []string{"GB", "US"}, nil, "2020-01-01", "2080-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, recordLabelRight)
	recordLabelRightId := bigchain.GetId(recordLabelRight)
	if err = api.Login(recordLabelId, recordLabelPriv); err != nil {
		t.Fatal(err)
	}
	release, err := api.Release(mechanicalLicenseId, recordingId, []string{performerRightId, producerRightId, recordLabelRightId})
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, release)
	releaseId := bigchain.GetId(release)
	masterLicense, err := api.MasterLicense(radioId, recordLabelRightId, "", releaseId, []string{"US"}, nil, "2020-01-01", "2022-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, masterLicense)
	// masterLicenseId := bigchain.GetId(masterLicense)
	if err = api.Login(composerId, composerPriv); err != nil {
		t.Fatal(err)
	}
	compositionRightTransfer, err := api.TransferCompositionRight(composerRightId, "", publicationId, publisherId, 10)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, compositionRightTransfer)
	compositionRightTransferId := bigchain.GetId(compositionRightTransfer)
	if err = api.Login(publisherId, publisherPriv); err != nil {
		t.Fatal(err)
	}
	compositionRightTransfer, err = api.TransferCompositionRight(composerRightId, compositionRightTransferId, publicationId, composerId, 5)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, compositionRightTransfer)
	compositionRightTransferId = bigchain.GetId(compositionRightTransfer)
	if err = api.Login(performerId, performerPriv); err != nil {
		t.Fatal(err)
	}
	recordingRightTransfer, err := api.TransferRecordingRight(performerRightId, "", recordLabelId, 10, releaseId)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, recordingRightTransfer)
	recordingRightTransferId := bigchain.GetId(recordingRightTransfer)
	if err = api.Login(recordLabelId, recordLabelPriv); err != nil {
		t.Fatal(err)
	}
	recordingRightTransfer, err = api.TransferRecordingRight(performerRightId, recordingRightTransferId, performerId, 5, releaseId)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, recordingRightTransfer)
	if err = api.Login(composerId, composerPriv); err != nil {
		t.Fatal(err)
	}
	mechanicalLicenseFromTransfer, err := api.MechanicalLicense("", compositionRightTransferId, producerId, publicationId, []string{"US"}, nil, "2020-01-01", "2030-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, mechanicalLicenseFromTransfer)
}
