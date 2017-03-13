package api

import (
	"testing"

	. "github.com/zbo14/envoke/common"
)

func GetPartyId(data Data) string {
	return data.GetStr("partyId")
}

func GetPartyPrivateKey(data Data) string {
	return data.GetStr("privateKey")
}

func TestApi(t *testing.T) {
	api := NewApi()
	output := MustOpenWriteFile("output.json")
	composer, err := api.Register("composer@email.com", "", "", nil, "composer", "itsasecret", "/Users/zach/Desktop/envoke/composer", "", "www.composer.com", "Person")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, composer)
	composerId := GetPartyId(composer)
	composerPriv := GetPartyPrivateKey(composer)
	recordLabel, err := api.Register("record_label@email.com", "", "", nil, "record_label", "shhhh", "/Users/zach/Desktop/envoke/record_label", "", "www.record_label.com", "Organization")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, recordLabel)
	recordLabelId := GetPartyId(recordLabel)
	recordLabelPriv := GetPartyPrivateKey(recordLabel)
	performer, err := api.Register("performer@email.com", "123456789", "", nil, "performer", "makeitup", "/Users/zach/Desktop/envoke/performer", "ASCAP", "www.performer.com", "MusicGroup")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, performer)
	performerId := GetPartyId(performer)
	performerPriv := GetPartyPrivateKey(performer)
	// producer, err := api.Register("producer@email.com", "producer", "1234", "www.soundcloud_page.com")
	// if err != nil {
	//	t.Fatal(err)
	// }
	// WriteJSON(output, producer)
	// producerId := GetPartyId(producer)
	publisher, err := api.Register("publisher@email.com", "", "", nil, "publisher", "didyousaysomething?", "/Users/zach/Desktop/envoke/publisher", "", "www.soundcloud_page.com", "MusicGroup")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, publisher)
	publisherId := GetPartyId(publisher)
	publisherPriv := GetPartyPrivateKey(publisher)
	radio, err := api.Register("radio@email.com", "", "", nil, "radio", "waves", "/Users/zach/Desktop/envoke/radio", "", "www.radio_station.com", "Organization")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, radio)
	radioId := GetPartyId(radio)
	if err = api.Login(composerId, composerPriv); err != nil {
		t.Fatal(err)
	}
	composition, err := api.Compose("B3107S", "T-034.524.680-1", "EN", "", "www.url_to_composition.com", "untitled")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, composition)
	compositionId := GetPartyId(composition)
	composerRight, err := api.CompositionRight(composerId, 20, []string{"GB", "US"}, "2020-01-01", "3000-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, composerRight)
	composerRightId := GetPartyId(composerRight)
	publisherRight, err := api.CompositionRight(publisherId, 80, []string{"GB", "US"}, "2020-01-01", "3000-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, publisherRight)
	publisherRightId := GetPartyId(publisherRight)
	if err = api.Login(publisherId, publisherPriv); err != nil {
		t.Fatal(err)
	}
	publication, err := api.Publish([]string{compositionId}, []string{composerRightId, publisherRightId}, "publication_title")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, publication)
	publicationId := GetPartyId(publication)
	mechanicalLicense, err := api.MechanicalLicense(nil, publisherRightId, "", publicationId, recordLabelId, []string{"US"}, nil, "2020-01-01", "2025-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, mechanicalLicense)
	mechanicalLicenseId := GetPartyId(mechanicalLicense)
	file, err := OpenFile(Getenv("PATH_TO_AUDIO_FILE"))
	if err != nil {
		t.Fatal(err)
	}
	if err = api.Login(performerId, performerPriv); err != nil {
		t.Fatal(err)
	}
	recording, err := api.Record(compositionId, "", "PT2M43S", file, "US-S1Z-99-00001", mechanicalLicenseId, performerId, "")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, recording)
	recordingId := GetPartyId(recording)
	performerRight, err := api.RecordingRight(performerId, 30, []string{"GB", "US"}, "2020-01-01", "2080-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, performerRight)
	performerRightId := GetPartyId(performerRight)
	recordLabelRight, err := api.RecordingRight(recordLabelId, 70, []string{"GB", "US"}, "2020-01-01", "2080-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, recordLabelRight)
	recordLabelRightId := GetPartyId(recordLabelRight)
	if err = api.Login(recordLabelId, recordLabelPriv); err != nil {
		t.Fatal(err)
	}
	release, err := api.Release([]string{recordingId}, []string{performerRightId, recordLabelRightId}, "release_title")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, release)
	releaseId := GetPartyId(release)
	masterLicense, err := api.MasterLicense(radioId, nil, recordLabelRightId, "", releaseId, []string{"US"}, nil, "2020-01-01", "2022-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, masterLicense)
	// masterLicenseId := GetPartyId(masterLicense)
	if err = api.Login(composerId, composerPriv); err != nil {
		t.Fatal(err)
	}
	compositionRightTransfer, err := api.TransferCompositionRight(composerRightId, "", publicationId, publisherId, 10)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, compositionRightTransfer)
	compositionRightTransferId := GetPartyId(compositionRightTransfer)
	if err = api.Login(publisherId, publisherPriv); err != nil {
		t.Fatal(err)
	}
	compositionRightTransfer, err = api.TransferCompositionRight(composerRightId, compositionRightTransferId, publicationId, composerId, 5)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, compositionRightTransfer)
	compositionRightTransferId = GetPartyId(compositionRightTransfer)
	if err = api.Login(performerId, performerPriv); err != nil {
		t.Fatal(err)
	}
	recordingRightTransfer, err := api.TransferRecordingRight(performerRightId, "", recordLabelId, 10, releaseId)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, recordingRightTransfer)
	recordingRightTransferId := GetPartyId(recordingRightTransfer)
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
	mechanicalLicenseFromTransfer, err := api.MechanicalLicense(nil, "", compositionRightTransferId, publicationId, radioId, []string{"US"}, nil, "2020-01-01", "2030-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, mechanicalLicenseFromTransfer)
}
