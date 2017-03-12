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
	mux.HandleFunc("/login", api.LoginHandler)
	mux.HandleFunc("/register", api.RegisterHandler)
	mux.HandleFunc("/compose", api.ComposeHandler)
	mux.HandleFunc("/record", api.RecordHandler)
	mux.HandleFunc("/composition_right", api.CompositionRightHandler)
	mux.HandleFunc("/recording_right", api.RecordingRightHandler)
	mux.HandleFunc("/publish", api.PublishHandler)
	mux.HandleFunc("/release", api.ReleaseHandler)
	mux.HandleFunc("/mechanical_license", api.MechanicalLicenseHandler)
	mux.HandleFunc("/master_license", api.MasterLicenseHandler)
	mux.HandleFunc("/transfer", api.TransferHandler)
	mux.HandleFunc("/search", api.SearchHandler)
	// mux.HandleFunc("/prove", api.ProveHandler)
	// mux.HandleFunc("/verify", api.VerifyHandler)
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
	v := struct {
		PartyId    string `json:"partyId"`
		PrivateKey string `json:"privateKey"`
	}{}
	if err = ReadJSON(credentials, &v); err != nil {
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
	if err = api.Register(email, ipi, isni, memberIds, name, password, path, proId, sameAs, _type); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("Registration successful!"))
}

func (api *Api) CompositionRightHandler(w http.ResponseWriter, req *http.Request) {
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
	recipientId := values.Get("recipientId")
	recipientShares := MustAtoi(values.Get("recipientShares"))
	territory := SplitStr(values.Get("territory"), ",")
	validFrom := values.Get("validFrom")
	validThrough := values.Get("validThrough")
	right, err := api.CompositionRight(recipientId, recipientShares, territory, validFrom, validThrough)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, right)
}

func (api *Api) RecordingRightHandler(w http.ResponseWriter, req *http.Request) {
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
	recipientId := values.Get("recipientId")
	recipientShares := MustAtoi(values.Get("recipientShares"))
	territory := SplitStr(values.Get("territory"), ",")
	validFrom := values.Get("validFrom")
	validThrough := values.Get("validThrough")
	right, err := api.RecordingRight(recipientId, recipientShares, territory, validFrom, validThrough)
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
	ipi := values.Get("ipi")
	iswc := values.Get("iswc")
	pro := values.Get("pro")
	title := values.Get("title")
	composition, err := api.Compose(hfa, ipi, iswc, pro, title)
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

func (api *Api) MechanicalLicenseHandler(w http.ResponseWriter, req *http.Request) {
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
	compositionIds := SplitStr(values.Get("compositionIds"), ",")
	publicationId := values.Get("publicationId")
	recipientId := values.Get("recipientId")
	rightId := values.Get("rightId")
	territory := SplitStr(values.Get("territory"), ",")
	transferId := values.Get("transferId")
	usage := SplitStr(values.Get("usage"), ",")
	validFrom := values.Get("validFrom")
	validThrough := values.Get("validThrough")
	license, err := api.MechanicalLicense(compositionIds, rightId, transferId, publicationId, recipientId, territory, usage, validFrom, validThrough)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, license)
}

func (api *Api) MasterLicenseHandler(w http.ResponseWriter, req *http.Request) {
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
	recipientId := values.Get("recipientId")
	recordingIds := SplitStr(values.Get("recordingIds"), ",")
	releaseId := values.Get("releaseId")
	rightId := values.Get("rightId")
	territory := SplitStr(values.Get("territory"), ",")
	transferId := values.Get("transferId")
	usage := SplitStr(values.Get("usage"), ",")
	validFrom := values.Get("validFrom")
	validThrough := values.Get("validThrough")
	license, err := api.MasterLicense(recipientId, recordingIds, rightId, transferId, releaseId, territory, usage, validFrom, validThrough)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, license)
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
	publicationId := values.Get("publicationId")
	releaseId := values.Get("releaseId")
	switch {
	case !EmptyStr(publicationId):
		model, err = ld.QueryPublicationField(field, publicationId)
	case !EmptyStr(releaseId):
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
	publicationId := values.Get("publicationId")
	recipientId := values.Get("recipientId")
	recipientShares := MustAtoi(values.Get("recipientShares"))
	releaseId := values.Get("releaseId")
	rightId := values.Get("rightId")
	rightTransferId := values.Get("rightTransferId")
	switch {
	case !EmptyStr(publicationId):
		transfer, err = api.TransferCompositionRight(rightId, rightTransferId, publicationId, recipientId, recipientShares)
	case !EmptyStr(releaseId):
		transfer, err = api.TransferRecordingRight(rightId, rightTransferId, recipientId, recipientShares, releaseId)
	default:
		http.Error(w, "Expected publicationId or releaseId", http.StatusBadRequest)
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

func (api *Api) Register(email, ipi, isni string, memberIds []string, name, password, path, proId, sameAs, _type string) error {
	priv, pub := ed25519.GenerateKeypairFromPassword(password)
	party := spec.NewParty(email, ipi, isni, memberIds, name, proId, sameAs, _type)
	tx := bigchain.DefaultIndividualCreateTx(party, pub)
	bigchain.FulfillTx(tx, priv)
	partyId, err := bigchain.PostTx(tx)
	if err != nil {
		return err
	}
	api.logger.Info("SUCCESS registered new party: " + name)
	file, err := CreateFile(path)
	if err != nil {
		return err
	}
	WriteJSON(file, &Data{
		"partyId":    partyId,
		"privateKey": priv.String(),
	})
	return nil
}

func (api *Api) Compose(hfa, ipi, iswc, pro, title string) (Data, error) {
	composition := spec.NewComposition(api.partyId, hfa, iswc, title)
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
