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
	composerAssignment, err := api.AssignCompositionRight(compositionId, composerId, 20, []string{"GB", "US"}, "2020-01-01", "3000-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, composerAssignment)
	composerAssignmentId := bigchain.GetId(composerAssignment)
	publisherAssignment, err := api.AssignCompositionRight(compositionId, publisherId, 80, []string{"GB", "US"}, "2020-01-01", "3000-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, publisherAssignment)
	publisherAssignmentId := bigchain.GetId(publisherAssignment)
	if err = api.Login(publisherId, publisherPriv); err != nil {
		t.Fatal(err)
	}
	publication, err := api.Publish([]string{composerAssignmentId, publisherAssignmentId}, compositionId)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, publication)
	publicationId := bigchain.GetId(publication)
	mechanicalLicense, err := api.MechanicalLicense(publisherAssignmentId, labelId, publicationId, []string{"US"}, "2020-01-01", "2025-01-01")
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
	performerAssignment, err := api.AssignRecordingRight(performerId, 20, recordingId, []string{"GB", "US"}, "2020-01-01", "2080-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, performerAssignment)
	performerAssignmentId := bigchain.GetId(performerAssignment)
	producerAssignment, err := api.AssignRecordingRight(producerId, 10, recordingId, []string{"GB", "US"}, "2020-01-01", "2080-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, producerAssignment)
	producerAssignmentId := bigchain.GetId(producerAssignment)
	labelAssignment, err := api.AssignRecordingRight(labelId, 70, recordingId, []string{"GB", "US"}, "2020-01-01", "2080-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, labelAssignment)
	labelAssignmentId := bigchain.GetId(labelAssignment)
	if err = api.Login(labelId, labelPriv); err != nil {
		t.Fatal(err)
	}
	release, err := api.Release([]string{labelAssignmentId, performerAssignmentId, producerAssignmentId}, mechanicalLicenseId, recordingId)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, release)
	releaseId := bigchain.GetId(release)
	masterLicense, err := api.MasterLicense(labelAssignmentId, radioId, releaseId, []string{"US"}, "2020-01-01", "2022-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, masterLicense)
	// masterLicenseId := bigchain.GetId(masterLicense)
	if err = api.Login(composerId, composerPriv); err != nil {
		t.Fatal(err)
	}
	compositionRightTransfer, err := api.TransferCompositionRight(composerAssignmentId, publicationId, publisherId, 10, "")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, compositionRightTransfer)
	compositionRightTransferId := bigchain.GetId(compositionRightTransfer)
	if err = api.Login(publisherId, publisherPriv); err != nil {
		t.Fatal(err)
	}
	compositionRightTransfer, err = api.TransferCompositionRight(composerAssignmentId, publicationId, composerId, 5, compositionRightTransferId)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, compositionRightTransfer)
	if err = api.Login(performerId, performerPriv); err != nil {
		t.Fatal(err)
	}
	recordingRightTransfer, err := api.TransferRecordingRight(performerAssignmentId, labelId, 10, releaseId, "")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, recordingRightTransfer)
	recordingRightTransferId := bigchain.GetId(recordingRightTransfer)
	if err = api.Login(labelId, labelPriv); err != nil {
		t.Fatal(err)
	}
	recordingRightTransfer, err = api.TransferRecordingRight(performerAssignmentId, performerId, 5, releaseId, recordingRightTransferId)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, recordingRightTransfer)
}
