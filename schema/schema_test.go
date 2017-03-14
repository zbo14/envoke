package schema

import (
	"testing"

	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/spec"
)

func TestSchema(t *testing.T) {
	composerId := BytesToHex(Checksum256(nil))
	composition := spec.NewComposition(composerId, "B3107S", "T-034.524.680-1", "EN", "untitled", "http://www.composition.com")
	if err := ValidateModel(composition, "composition"); err != nil {
		t.Error(err)
	}
	compositionId := BytesToHex(Checksum256(MustMarshalJSON(composition)))
	compositionRight := spec.NewCompositionRight(composerId, composerId, []string{"US"}, "2018-01-01", "2088-01-01")
	if err := ValidateModel(compositionRight, "right"); err != nil {
		t.Error(err)
	}
	compositionRightId := BytesToHex(Checksum256(MustMarshalJSON(compositionRight)))
	publisherId := BytesToHex(Checksum256([]byte{1, 2, 3}))
	publication := spec.NewPublication([]string{compositionId}, []string{compositionRightId}, "publication_name", publisherId)
	if err := ValidateModel(publication, "publication"); err != nil {
		t.Error(err)
	}
	publicationId := BytesToHex(Checksum256(MustMarshalJSON(publication)))
	licenseeId := BytesToHex(Checksum256([]byte{4, 5, 6}))
	mechanicalLicense := spec.NewMechanicalLicense(nil, compositionRightId, "", publicationId, licenseeId, publisherId, []string{"US"}, []string{"USAGE"}, "2018-01-01", "2024-01-01")
	if err := ValidateModel(mechanicalLicense, "mechanical_license"); err != nil {
		PrintJSON(mechanicalLicense)
		t.Error(err)
	}
}
