package api

import (
	"testing"

	. "github.com/zbo14/envoke/common"
)

func GetId(data Data) string {
	return data.GetStr("id")
}

func GetPrivateKey(data Data) string {
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
	composerId := GetId(composer)
	composerPriv := GetPrivateKey(composer)
	recordLabel, err := api.Register("record_label@email.com", "", "", nil, "record_label", "shhhh", "/Users/zach/Desktop/envoke/record_label", "", "www.record_label.com", "Organization")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, recordLabel)
	recordLabelId := GetId(recordLabel)
	recordLabelPriv := GetPrivateKey(recordLabel)
	performer, err := api.Register("performer@email.com", "123456789", "", nil, "performer", "makeitup", "/Users/zach/Desktop/envoke/performer", "ASCAP", "www.performer.com", "MusicGroup")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, performer)
	performerId := GetId(performer)
	performerPriv := GetPrivateKey(performer)
	// producer, err := api.Register("producer@email.com", "producer", "1234", "www.soundcloud_page.com")
	// if err != nil {
	//	t.Fatal(err)
	// }
	// WriteJSON(output, producer)
	// producerId := GetId(producer)
	publisher, err := api.Register("publisher@email.com", "", "", nil, "publisher", "didyousaysomething?", "/Users/zach/Desktop/envoke/publisher", "", "www.soundcloud_page.com", "MusicGroup")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, publisher)
	publisherId := GetId(publisher)
	publisherPriv := GetPrivateKey(publisher)
	radio, err := api.Register("radio@email.com", "", "", nil, "radio", "waves", "/Users/zach/Desktop/envoke/radio", "", "www.radio_station.com", "Organization")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, radio)
	radioId := GetId(radio)
	if err = api.Login(composerId, composerPriv); err != nil {
		t.Fatal(err)
	}
	composition, err := api.Compose("B3107S", "T-034.524.680-1", "EN", "www.url_to_composition.com", "untitled")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, composition)
	compositionId := GetId(composition)
	composerRight, err := api.CompositionRight(composerId, 20, []string{"GB", "US"}, "2020-01-01", "2096-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, composerRight)
	composerRightId := GetId(composerRight)
	publisherRight, err := api.CompositionRight(publisherId, 80, []string{"GB", "US"}, "2020-01-01", "2096-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, publisherRight)
	publisherRightId := GetId(publisherRight)
	publication, err := api.Publish([]string{compositionId}, []string{composerRightId, publisherRightId}, publisherId, "publication_title")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, publication)
	publicationId := GetId(publication)
	if err = api.Login(publisherId, publisherPriv); err != nil {
		t.Fatal(err)
	}
	mechanicalLicense, err := api.MechanicalLicense(nil, publisherRightId, "", publicationId, performerId, []string{"US"}, nil, "2020-01-01", "2024-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, mechanicalLicense)
	mechanicalLicenseId := GetId(mechanicalLicense)
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
	recordingId := GetId(recording)
	performerRight, err := api.RecordingRight(performerId, 30, []string{"GB", "US"}, "2020-01-01", "2080-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, performerRight)
	performerRightId := GetId(performerRight)
	recordLabelRight, err := api.RecordingRight(recordLabelId, 70, []string{"GB", "US"}, "2020-01-01", "2080-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, recordLabelRight)
	recordLabelRightId := GetId(recordLabelRight)
	release, err := api.Release([]string{recordingId}, []string{performerRightId, recordLabelRightId}, recordLabelId, "release_title")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, release)
	releaseId := GetId(release)
	if err = api.Login(recordLabelId, recordLabelPriv); err != nil {
		t.Fatal(err)
	}
	masterLicense, err := api.MasterLicense(radioId, nil, recordLabelRightId, "", releaseId, []string{"US"}, nil, "2020-01-01", "2022-01-01")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, masterLicense)
	// masterLicenseId := GetId(masterLicense)
	if err = api.Login(composerId, composerPriv); err != nil {
		t.Fatal(err)
	}
	SleepSeconds(2)
	compositionRightTransfer, err := api.TransferCompositionRight(composerRightId, "", publicationId, publisherId, 10)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, compositionRightTransfer)
	compositionRightTransferId := GetId(compositionRightTransfer)
	if err = api.Login(publisherId, publisherPriv); err != nil {
		t.Fatal(err)
	}
	SleepSeconds(2)
	compositionRightTransfer, err = api.TransferCompositionRight(composerRightId, compositionRightTransferId, publicationId, composerId, 5)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, compositionRightTransfer)
	compositionRightTransferId = GetId(compositionRightTransfer)
	if err = api.Login(performerId, performerPriv); err != nil {
		t.Fatal(err)
	}
	SleepSeconds(2)
	recordingRightTransfer, err := api.TransferRecordingRight(recordLabelId, 10, performerRightId, "", releaseId)
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, recordingRightTransfer)
	recordingRightTransferId := GetId(recordingRightTransfer)
	if err = api.Login(recordLabelId, recordLabelPriv); err != nil {
		t.Fatal(err)
	}
	SleepSeconds(2)
	recordingRightTransfer, err = api.TransferRecordingRight(performerId, 5, performerRightId, recordingRightTransferId, releaseId)
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
