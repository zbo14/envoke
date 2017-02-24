package api

import (
	. "github.com/zbo14/envoke/common"
	// "github.com/zbo14/envoke/spec"
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
	// composerId := composer.AgentId
	label, err := api.Register("label@gmail.com", "label", "shhh", "www.record_label.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, label)
	// labelId := label.AgentId
	performer, err := api.Register("performer@gmail.com", "performer", "canyouguess", "www.bandcamp_page.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, performer)
	// performerId := performer.AgentId
	producer, err := api.Register("producer@gmail.com", "producer", "1234", "www.soundcloud_page.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, producer)
	// producerId := producer.AgentId
	publisher, err := api.Register("publisher@gmail.com", "publisher", "password", "www.publisher.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, publisher)
	// publisherId := publisher.AgentId
	radio, err := api.Register("radio@gmail.com", "radio", "waves", "www.radio_station.com")
	if err != nil {
		t.Fatal(err)
	}
	WriteJSON(output, radio)
	// radioId := radio.AgentId
	/*
		if err = api.Login(composerId, composer.PrivKey); err != nil {
			t.Fatal(err)
		}
		compositionInfo, err := api.CompositionInfo(composerId, publisherId, "untitled")
		if err != nil {
			t.Fatal(err)
		}
		compositionInfoId := compositionInfo.GetStr("id")
		composerRight, err := api.Right(compositionInfoId, "20", "2020-01-01", "3000-01-01")
		if err != nil {
			t.Fatal(err)
		}
		composerRightId := composerRight.GetStr("id")
		if err = api.Login(publisherId, publisher.PrivKey); err != nil {
			t.Fatal(err)
		}
		publisherRight, err := api.Right(compositionInfoId, "80", "2020-01-01", "2030-01-01")
		if err != nil {
			t.Fatal(err)
		}
		publisherRightId := publisherRight.GetStr("id")
		if err = api.Login(composerId, composer.PrivKey); err != nil {
			t.Fatal(err)
		}
		composition, err := api.Composition(compositionInfoId, []string{composerRightId, publisherRightId})
		if err != nil {
			t.Fatal(err)
		}
		compositionId := composition.GetStr("id")
		if err = api.Login(publisherId, publisher.PrivKey); err != nil {
			t.Fatal(err)
		}
		Println(compositionId)
		publishingLicense, err := api.PublishingLicense(compositionId, labelId, spec.LICENSE_TYPE_MECHANICAL, "2020-01-01", "2025-01-01")
		if err != nil {
			t.Fatal(err)
		}
		publishingLicenseId := publishingLicense.GetStr("id")
		if err = api.Login(labelId, label.PrivKey); err != nil {
			t.Fatal(err)
		}
		file, err := OpenFile(path)
		if err != nil {
			t.Fatal(err)
		}
		recordingInfo, err := api.RecordingInfo(compositionId, file, labelId, performerId, producerId, publishingLicenseId)
		if err != nil {
			t.Fatal(err)
		}
		recordingInfoId := recordingInfo.GetStr("id")
		labelRight, err := api.Right(recordingInfoId, "70", "2020-01-01", "2080-01-01")
		if err != nil {
			t.Fatal(err)
		}
		labelRightId := labelRight.GetStr("id")
		if err = api.Login(performerId, performer.PrivKey); err != nil {
			t.Fatal(err)
		}
		performerRight, err := api.Right(recordingInfoId, "20", "2020-01-01", "2080-01-01")
		if err != nil {
			t.Fatal(err)
		}
		performerRightId := performerRight.GetStr("id")
		if err = api.Login(producerId, producer.PrivKey); err != nil {
			t.Fatal(err)
		}
		producerRight, err := api.Right(recordingInfoId, "10", "2020-01-01", "2080-01-01")
		if err != nil {
			t.Fatal(err)
		}
		producerRightId := producerRight.GetStr("id")
		if err = api.Login(labelId, label.PrivKey); err != nil {
			t.Fatal(err)
		}
		recording, err := api.Recording(recordingInfoId, []string{labelRightId, performerRightId, producerRightId})
		if err != nil {
			t.Fatal(err)
		}
		recordingId := recording.GetStr("id")
		recordingLicense, err := api.RecordingLicense(radioId, spec.LICENSE_TYPE_MASTER, recordingId, "2020-01-01", "2022")
		if err != nil {
			t.Fatal(err)
		}
		recordingLicenseId := recordingLicense.GetStr("id")
		t.Log(recordingLicenseId)
	*/
}
