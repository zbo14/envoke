package api

import (
	"io"
	"net/http"

	// "github.com/dhowden/tag"
	"github.com/zbo14/envoke/bigchain"
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/crypto"
	"github.com/zbo14/envoke/crypto/ed25519"
	ld "github.com/zbo14/envoke/linked_data"
	"github.com/zbo14/envoke/spec"
)

type Api struct {
	partyId string
	logger  Logger
	priv    crypto.PrivateKey
	pub     crypto.PublicKey
}

func NewApi() *Api {
	return &Api{
		logger: NewLogger("api"),
	}
}

func (api *Api) AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/login_handler", api.LoginHandler)
	mux.HandleFunc("/register_handler", api.RegisterHandler)
	mux.HandleFunc("/compose_handler", api.ComposeHandler)
	mux.HandleFunc("/record_handler", api.RecordHandler)
	mux.HandleFunc("/right_handler", api.RightHandler)
	mux.HandleFunc("/publish_handler", api.PublishHandler)
	mux.HandleFunc("/release_handler", api.ReleaseHandler)
	mux.HandleFunc("/license_handler", api.LicenseHandler)
	mux.HandleFunc("/transfer_handler", api.TransferHandler)
	mux.HandleFunc("/search_handler", api.SearchHandler)
	mux.HandleFunc("/prove_handler", api.ProveHandler)
	mux.HandleFunc("/verify_handler", api.VerifyHandler)
}

func (api *Api) LoginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	form, err := MultipartForm(req)
	if err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	credentials, err := form.File["credentials"][0].Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	v := &struct {
		PartyId    string `json:"partyId"`
		PrivateKey string `json:"privateKey"`
	}{}
	p, err := ReadAll(credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = UnmarshalJSON(p, v); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := api.Login(v.PartyId, v.PrivateKey); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api *Api) RegisterHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	email := values.Get("email")
	ipi := values.Get("ipi")
	isni := values.Get("isni")
	memberIds := SplitStr(values.Get("memberIds"), ",")
	name := values.Get("name")
	password := values.Get("password")
	path := values.Get("path")
	proId := values.Get("proId")
	sameAs := values.Get("sameAs")
	_type := values.Get("type")
	if _, err = api.Register(email, ipi, isni, memberIds, name, password, path, proId, sameAs, _type); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("Registration successful!"))
}

func (api *Api) RightHandler(w http.ResponseWriter, req *http.Request) {
	if !api.LoggedIn() {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var right Data
	recipientId := values.Get("recipientId")
	recipientShares := MustAtoi(values.Get("recipientShares"))
	territory := SplitStr(values.Get("territory"), ",")
	_type := values.Get("type")
	validFrom := values.Get("validFrom")
	validThrough := values.Get("validThrough")
	if _type == "composition_right" {
		right, err = api.CompositionRight(recipientId, recipientShares, territory, validFrom, validThrough)
	} else if _type == "recording_right" {
		right, err = api.RecordingRight(recipientId, recipientShares, territory, validFrom, validThrough)
	} else {
		http.Error(w, ErrorAppend(ErrInvalidType, _type).Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, right)
}

func (api *Api) ComposeHandler(w http.ResponseWriter, req *http.Request) {
	if !api.LoggedIn() {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hfa := values.Get("hfa")
	iswc := values.Get("iswc")
	lang := values.Get("lang")
	lyrics := values.Get("lyrics")
	sameAs := values.Get("sameAs")
	title := values.Get("title")
	composition, err := api.Compose(hfa, iswc, lang, lyrics, sameAs, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, composition)
}

func (api *Api) RecordHandler(w http.ResponseWriter, req *http.Request) {
	if !api.LoggedIn() {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	form, err := MultipartForm(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	compositionId := form.Value["compositionId"][0]
	compositionRightId := form.Value["compositionRightId"][0]
	duration := form.Value["duration"][0]
	file, err := form.File["recording"][0].Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	isrc := form.Value["isrc"][0]
	mechanicalLicenseId := form.Value["mechanicalLicenseId"][0]
	performerId := form.Value["performerId"][0]
	publicationId := form.Value["publicationId"][0]
	recording, err := api.Record(compositionId, compositionRightId, duration, file, isrc, mechanicalLicenseId, performerId, publicationId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, recording)
}

func (api *Api) PublishHandler(w http.ResponseWriter, req *http.Request) {
	if !api.LoggedIn() {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	compositionsId := SplitStr(values.Get("compositionId"), ",")
	compositionRightIds := SplitStr(values.Get("compositionRightIds"), ",")
	title := values.Get("title")
	composition, err := api.Publish(compositionsId, compositionRightIds, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, composition)
}

func (api *Api) ReleaseHandler(w http.ResponseWriter, req *http.Request) {
	if !api.LoggedIn() {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	recordingIds := SplitStr(values.Get("recordingId"), ",")
	recordingRightIds := SplitStr(values.Get("recordingRightIds"), ",")
	title := values.Get("title")
	release, err := api.Release(recordingIds, recordingRightIds, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, release)
}

func (api *Api) LicenseHandler(w http.ResponseWriter, req *http.Request) {
	if !api.LoggedIn() {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var license Data

	recipientId := values.Get("recipientId")
	rightId := values.Get("rightId")
	territory := SplitStr(values.Get("territory"), ",")
	transferId := values.Get("transferId")
	_type := values.Get("type")
	usage := SplitStr(values.Get("usage"), ",")
	validFrom := values.Get("validFrom")
	validThrough := values.Get("validThrough")
	if _type == "mechanical_license" {
		compositionIds := SplitStr(values.Get("compositionIds"), ",")
		publicationId := values.Get("publicationId")
		license, err = api.MechanicalLicense(compositionIds, rightId, transferId, publicationId, recipientId, territory, usage, validFrom, validThrough)
	} else if _type == "master_license" {
		recordingIds := SplitStr(values.Get("recordingIds"), ",")
		releaseId := values.Get("releaseId")
		license, err = api.MasterLicense(recipientId, recordingIds, rightId, transferId, releaseId, territory, usage, validFrom, validThrough)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, license)
}

func (api *Api) ProveHandler(w http.ResponseWriter, req *http.Request) {
	if !api.LoggedIn() {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var sig crypto.Signature
	challenge := values.Get("challenge")
	_type := values.Get("type")
	switch _type {
	case "composition":
		compositionId := values.Get("compositionId")
		sig, err = ld.ProveComposer(challenge, compositionId, api.priv)
	case "composition_right":
		rightId := values.Get("rightId")
		publicationId := values.Get("publicationReleaseId")
		sig, err = ld.ProveCompositionRightHolder(challenge, rightId, api.priv, publicationId)
	case "composition_right_transfer":
		transferId := values.Get("transferId")
		publicationId := values.Get("publicationReleaseId")
		sig, err = ld.ProveCompositionRightTransferHolder(challenge, transferId, api.partyId, api.priv, publicationId)
	case "master_license":
		licenseId := values.Get("licenseId")
		sig, err = ld.ProveMasterLicenseHolder(challenge, licenseId, api.priv)
	case "mechanical_license":
		licenseId := values.Get("licenseId")
		sig, err = ld.ProveMechanicalLicenseHolder(challenge, licenseId, api.priv)
	case "publication":
		publicationId := values.Get("publicationId")
		sig, err = ld.ProvePublisher(challenge, api.priv, publicationId)
	case "recording":
		recordingId := values.Get("recordingId")
		sig, err = ld.ProvePerformer(challenge, api.priv, recordingId)
	case "recording_right":
		rightId := values.Get("rightId")
		releaseId := values.Get("publicationReleaseId")
		sig, err = ld.ProveRecordingRightHolder(challenge, api.priv, rightId, releaseId)
	case "recording_right_transfer":
		transferId := values.Get("transferId")
		releaseId := values.Get("publicationReleaseId")
		sig, err = ld.ProveRecordingRightTransferHolder(challenge, api.partyId, api.priv, transferId, releaseId)
	case "release":
		releaseId := values.Get("releaseId")
		sig, err = ld.ProveRecordLabel(challenge, api.priv, releaseId)
	default:
		http.Error(w, ErrorAppend(ErrInvalidType, _type).Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, sig)
}

func (api *Api) VerifyHandler(w http.ResponseWriter, req *http.Request) {
	if !api.LoggedIn() {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	challenge := values.Get("challenge")
	sig := new(ed25519.Signature)
	signature := values.Get("signature")
	if err := sig.FromString(signature); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_type := values.Get("type")
	switch _type {
	case "composition":
		compositionId := values.Get("compositionId")
		err = ld.VerifyComposer(challenge, compositionId, sig)
	case "composition_right":
		rightId := values.Get("rightId")
		publicationId := values.Get("publicationReleaseId")
		err = ld.VerifyCompositionRightHolder(challenge, rightId, publicationId, sig)
	case "composition_right_transfer":
		transferId := values.Get("transferId")
		publicationId := values.Get("publicationReleaseId")
		err = ld.VerifyCompositionRightTransferHolder(challenge, transferId, api.partyId, publicationId, sig)
	case "master_license":
		licenseId := values.Get("licenseId")
		err = ld.VerifyMasterLicenseHolder(challenge, licenseId, sig)
	case "mechanical_license":
		licenseId := values.Get("licenseId")
		err = ld.VerifyMechanicalLicenseHolder(challenge, licenseId, sig)
	case "publication":
		publicationId := values.Get("publicationId")
		err = ld.VerifyPublisher(challenge, publicationId, sig)
	case "recording":
		recordingId := values.Get("recordingId")
		err = ld.VerifyPerformer(challenge, recordingId, sig)
	case "recording_right":
		rightId := values.Get("rightId")
		releaseId := values.Get("publicationReleaseId")
		err = ld.VerifyRecordingRightHolder(challenge, rightId, releaseId, sig)
	case "recording_right_transfer":
		transferId := values.Get("transferId")
		releaseId := values.Get("publicationReleaseId")
		err = ld.VerifyRecordingRightTransferHolder(challenge, api.partyId, transferId, releaseId, sig)
	case "release":
		releaseId := values.Get("releaseId")
		err = ld.VerifyRecordLabel(challenge, releaseId, sig)
	default:
		http.Error(w, ErrorAppend(ErrInvalidType, _type).Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, "Verified signature!")
}

func (api *Api) SearchHandler(w http.ResponseWriter, req *http.Request) {
	if !api.LoggedIn() {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var model interface{}
	field := values.Get("field")
	_type := values.Get("type")
	switch _type {
	case "composition":
		compositionId := values.Get("compositionId")
		model, err = ld.QueryCompositionField(field, compositionId)
	case "master_license":
		licenseId := values.Get("licenseId")
		model, err = ld.QueryMasterLicenseField(field, licenseId)
	case "mechanical_license":
		licenseId := values.Get("licenseId")
		model, err = ld.QueryMechanicalLicenseField(field, licenseId)
	case "publication":
		publicationId := values.Get("publicationId")
		model, err = ld.QueryPublicationField(field, publicationId)
	case "recording":
		recordingId := values.Get("recordingId")
		model, err = ld.QueryRecordingField(field, recordingId)
	case "release":
		releaseId := values.Get("releaseId")
		model, err = ld.QueryReleaseField(field, releaseId)
	default:
		http.Error(w, "Expected publicationId or releaseId", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, model)
}

func (api *Api) TransferHandler(w http.ResponseWriter, req *http.Request) {
	if !api.LoggedIn() {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var transfer Data
	recipientId := values.Get("recipientId")
	recipientShares := MustAtoi(values.Get("recipientShares"))
	rightId := values.Get("rightId")
	transferId := values.Get("transferId")
	_type := values.Get("type")
	switch _type {
	case "composition_right_transfer":
		publicationId := values.Get("publicationReleaseId")
		transfer, err = api.TransferCompositionRight(rightId, transferId, publicationId, recipientId, recipientShares)
	case "recording_right_transfer":
		releaseId := values.Get("publicationReleaseId")
		transfer, err = api.TransferRecordingRight(rightId, transferId, recipientId, recipientShares, releaseId)
	default:
		http.Error(w, ErrorAppend(ErrInvalidType, _type).Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, transfer)
}

func (api *Api) LoggedIn() bool {
	switch {
	case api.partyId == "":
		api.logger.Warn("Party ID is not set")
	case api.priv == nil:
		api.logger.Warn("Private-key is not set")
	case api.pub == nil:
		api.logger.Warn("Public-key is not set")
	default:
		return true
	}
	api.logger.Error("LOGIN FAILED")
	return false
}

func (api *Api) Login(partyId, privstr string) error {
	priv := new(ed25519.PrivateKey)
	if err := priv.FromString(privstr); err != nil {
		return err
	}
	tx, err := ld.QueryAndValidateModel(partyId)
	if err != nil {
		return err
	}
	party := bigchain.GetTxData(tx)
	pub := bigchain.DefaultGetTxSender(tx)
	if !pub.Equals(priv.Public()) {
		return ErrInvalidKey
	}
	api.partyId = partyId
	api.priv = priv
	api.pub = pub
	partyName := spec.GetName(party)
	api.logger.Info(Sprintf("SUCCESS %s is logged in", partyName))
	return nil
}

func (api *Api) Register(email, ipi, isni string, memberIds []string, name, password, path, proId, sameAs, _type string) (Data, error) {
	priv, pub := ed25519.GenerateKeypairFromPassword(password)
	party := spec.NewParty(email, ipi, isni, memberIds, name, proId, sameAs, _type)
	tx := bigchain.DefaultIndividualCreateTx(party, pub)
	bigchain.FulfillTx(tx, priv)
	partyId, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS registered new party: " + name)
	file, err := CreateFile(path + "/credentials.json")
	if err != nil {
		return nil, err
	}
	v := Data{
		"partyId":    partyId,
		"privateKey": priv.String(),
	}
	WriteJSON(file, &v)
	return v, nil
}

func (api *Api) Compose(hfa, iswc, lang, lyrics, sameAs, title string) (Data, error) {
	composition := spec.NewComposition(api.partyId, hfa, iswc, lang, lyrics, title, sameAs)
	tx := bigchain.DefaultIndividualCreateTx(composition, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with composition")
	return Data{
		"composition": composition,
		"id":          id,
	}, nil
}

func (api *Api) Record(compositionId, compositionRightId, duration string, file io.Reader, isrc, mechanicalLicenseId, performerId, publicationId string) (Data, error) {
	// rs := MustReadSeeker(file)
	// meta, err := tag.ReadFrom(rs)
	// if err != nil {
	//	return nil, err
	// }
	// metadata := meta.Raw()
	recording := spec.NewRecording(compositionId, compositionRightId, duration, isrc, mechanicalLicenseId, performerId, publicationId)
	tx := bigchain.DefaultIndividualCreateTx(recording, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with recording")
	return Data{
		"id":        id,
		"recording": recording,
	}, nil
}

func (api *Api) Publish(compositionIds, compositionRightIds []string, title string) (Data, error) {
	publication := spec.NewPublication(compositionIds, compositionRightIds, title, api.partyId)
	tx := bigchain.DefaultIndividualCreateTx(publication, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with publication")
	return Data{
		"id":          id,
		"publication": publication,
	}, nil
}

func (api *Api) Release(recordingIds, recordingRightIds []string, title string) (Data, error) {
	release := spec.NewRelease(title, recordingIds, recordingRightIds, api.partyId)
	tx := bigchain.DefaultIndividualCreateTx(release, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with release")
	return Data{
		"id":      id,
		"release": release,
	}, nil
}

func (api *Api) CompositionRight(recipientId string, recipientShares int, territory []string, validFrom, validThrough string) (Data, error) {
	tx, err := bigchain.GetTx(recipientId)
	if err != nil {
		return nil, err
	}
	recipientPub := bigchain.DefaultGetTxSender(tx)
	compositionRight := spec.NewCompositionRight(recipientId, api.partyId, territory, validFrom, validThrough)
	tx = bigchain.IndividualCreateTx(recipientShares, compositionRight, recipientPub, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with composition right")
	return Data{
		"compositionRight": compositionRight,
		"id":               id,
	}, nil
}

func (api *Api) RecordingRight(recipientId string, recipientShares int, territory []string, validFrom, validThrough string) (Data, error) {
	tx, err := bigchain.GetTx(recipientId)
	if err != nil {
		return nil, err
	}
	recipientPub := bigchain.DefaultGetTxSender(tx)
	recordingRight := spec.NewRecordingRight(recipientId, api.partyId, territory, validFrom, validThrough)
	tx = bigchain.IndividualCreateTx(recipientShares, recordingRight, recipientPub, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with recording right")
	return Data{
		"recordingRight": recordingRight,
		"id":             id,
	}, nil
}

func (api *Api) MechanicalLicense(compositionIds []string, compositionRightId, compositionRightTransferId, publicationId, recipientId string, territory, usage []string, validFrom, validThrough string) (Data, error) {
	mechanicalLicense := spec.NewMechanicalLicense(compositionIds, compositionRightId, compositionRightTransferId, publicationId, recipientId, api.partyId, territory, usage, validFrom, validThrough)
	tx := bigchain.DefaultIndividualCreateTx(mechanicalLicense, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with mechanical license")
	return Data{
		"id":                id,
		"mechanicalLicense": mechanicalLicense,
	}, nil
}

func (api *Api) MasterLicense(recipientId string, recordingIds []string, recordingRightId, recordingRightTransferId, releaseId string, territory, usage []string, validFrom, validThrough string) (Data, error) {
	masterLicense := spec.NewMasterLicense(recipientId, recordingIds, recordingRightId, recordingRightTransferId, releaseId, api.partyId, territory, usage, validFrom, validThrough)
	tx := bigchain.DefaultIndividualCreateTx(masterLicense, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with master license")
	return Data{
		"id":            id,
		"masterLicense": masterLicense,
	}, nil
}

func (api *Api) TransferCompositionRight(compositionRightId, compositionRightTransferId, publicationId, recipientId string, recipientShares int) (Data, error) {
	var output, totalShares int
	var txId string
	if !EmptyStr(compositionRightTransferId) {
		compositionRightTransfer, err := ld.ValidateCompositionRightTransfer(compositionRightTransferId)
		if err != nil {
			return nil, err
		}
		if api.partyId == spec.GetRecipientId(compositionRightTransfer) {
			totalShares = spec.GetRecipientShares(compositionRightTransfer)
		} else if api.partyId == spec.GetSenderId(compositionRightTransfer) {
			totalShares = spec.GetSenderShares(compositionRightTransfer)
			output = 1
		} else {
			return nil, ErrorAppend(ErrCriteriaNotMet, "partyId does not match recipientId or senderId of TRANSFER tx")
		}
		compositionRightId = spec.GetCompositionRightId(compositionRightTransfer)
		txId = spec.GetTxId(compositionRightTransfer)
	} else {
		tx, err := bigchain.GetTx(compositionRightId)
		if err != nil {
			return nil, err
		}
		compositionRight := bigchain.GetTxData(tx)
		totalShares = spec.GetRecipientShares(compositionRight)
		txId = compositionRightId
	}
	tx, err := bigchain.GetTx(recipientId)
	if err != nil {
		return nil, err
	}
	recipientPub := bigchain.DefaultGetTxSender(tx)
	senderShares := totalShares - recipientShares
	if senderShares < 0 {
		return nil, ErrorAppend(ErrCriteriaNotMet, "cannot transfer this many shares")
	}
	if senderShares == 0 {
		tx = bigchain.IndividualTransferTx(recipientShares, compositionRightId, txId, output, recipientPub, api.pub)
	} else {
		tx = bigchain.DivisibleTransferTx([]int{recipientShares, senderShares}, compositionRightId, txId, output, []crypto.PublicKey{recipientPub, api.pub}, api.pub)
	}
	bigchain.FulfillTx(tx, api.priv)
	txId, err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	compositionRightTransfer := spec.NewCompositionRightTransfer(compositionRightId, publicationId, recipientId, api.partyId, txId)
	tx = bigchain.DefaultIndividualCreateTx(compositionRightTransfer, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with composition right transfer")
	return Data{
		"compositionRightTransfer": compositionRightTransfer,
		"id": id,
	}, nil
}

func (api *Api) TransferRecordingRight(recordingRightId, recordingRightTransferId, recipientId string, recipientShares int, releaseId string) (Data, error) {
	var output, totalShares int
	var txId string
	if !EmptyStr(recordingRightTransferId) {
		recordingRightTransfer, err := ld.ValidateRecordingRightTransfer(recordingRightTransferId)
		if err != nil {
			return nil, err
		}
		if api.partyId == spec.GetRecipientId(recordingRightTransfer) {
			totalShares = spec.GetRecipientShares(recordingRightTransfer)
		} else if api.partyId == spec.GetSenderId(recordingRightTransfer) {
			totalShares = spec.GetSenderShares(recordingRightTransfer)
			output = 1
		} else {
			return nil, ErrorAppend(ErrCriteriaNotMet, "partyId does not match recipientId or senderId of TRANSFER tx")
		}
		recordingRightId = spec.GetRecordingRightId(recordingRightTransfer)
		txId = spec.GetTxId(recordingRightTransfer)
	} else {
		tx, err := bigchain.GetTx(recordingRightId)
		if err != nil {
			return nil, err
		}
		recordingRight := bigchain.GetTxData(tx)
		totalShares = spec.GetRecipientShares(recordingRight)
		txId = recordingRightId
	}
	tx, err := bigchain.GetTx(recipientId)
	if err != nil {
		return nil, err
	}
	recipientPub := bigchain.DefaultGetTxSender(tx)
	senderShares := totalShares - recipientShares
	if senderShares < 0 {
		return nil, ErrorAppend(ErrCriteriaNotMet, "cannot transfer this many shares")
	}
	if senderShares == 0 {
		tx = bigchain.IndividualTransferTx(recipientShares, recordingRightId, txId, output, recipientPub, api.pub)
	} else {
		tx = bigchain.DivisibleTransferTx([]int{recipientShares, senderShares}, recordingRightId, txId, output, []crypto.PublicKey{recipientPub, api.pub}, api.pub)
	}
	bigchain.FulfillTx(tx, api.priv)
	txId, err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	recordingRightTransfer := spec.NewRecordingRightTransfer(recipientId, recordingRightId, releaseId, api.partyId, txId)
	tx = bigchain.DefaultIndividualCreateTx(recordingRightTransfer, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with recording right transfer")
	return Data{
		"id": id,
		"recordingRightTransfer": recordingRightTransfer,
	}, nil
}
