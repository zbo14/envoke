package chroma

import (
	// . "github.com/zbo14/envoke/common"
	"testing"
)

var dir = "/Users/zach/Desktop/music/"

func TestChroma(t *testing.T) {
	jude1, err := NewFingerprint(120, dir+"hey_jude_1.mp3")
	if err != nil {
		t.Fatal(err)
	}
	rhapsody1, err := NewFingerprint(120, dir+"rhapsody_1.mp3")
	if err != nil {
		t.Fatal(err)
	}
	rhapsody2, err := NewFingerprint(120, dir+"rhapsody_2.mp3")
	if err != nil {
		t.Fatal(err)
	}
	same, err := CompareFingerprints(rhapsody1, rhapsody2)
	if err != nil {
		t.Fatal(err)
	}
	if !same {
		t.Error("Did not recognize match")
	}
	same, err = CompareFingerprints(jude1, rhapsody1)
	if err != nil {
		t.Fatal(err)
	}
	if same {
		t.Error("different songs")
	}
}
