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
	agent   Data
	agentId string
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
	mux.HandleFunc("/right", api.RightHandler)
	mux.HandleFunc("/publish", api.PublishHandler)
	mux.HandleFunc("/release", api.ReleaseHandler)
	mux.HandleFunc("/license", api.LicenseHandler)
	mux.HandleFunc("/transfer", api.TransferHandler)
	mux.HandleFunc("/search", api.SearchHandler)
}

func (api *Api) LoginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	agentId := values.Get("agentId")
	privstr := values.Get("privateKey")
	if err := api.Login(agentId, privstr); err != nil {
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
	name := values.Get("name")
	password := values.Get("password")
	socialMedia := values.Get("socialMedia")
	msg, err := api.Register(email, name, password, socialMedia)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, msg)
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
	compositionId := values.Get("compositionId")
	percentageShares := MustAtoi(values.Get("percentageShares"))
	recordingId := values.Get("recordingId")
	rightHolderId := values.Get("rightHolderId")
	territory := SplitStr(values.Get("territory"), ",")
	validFrom := values.Get("validFrom")
	validTo := values.Get("validTo")
	switch {
	case !EmptyStr(compositionId):
		right, err = api.CompositionRight(compositionId, percentageShares, rightHolderId, territory, validFrom, validTo)
	case !EmptyStr(recordingId):
		right, err = api.RecordingRight(percentageShares, recordingId, rightHolderId, territory, validFrom, validTo)
	default:
		http.Error(w, "Expected compositionId or recordingId", http.StatusBadRequest)
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
	composerId := values.Get("composerId")
	hfa := values.Get("hfa")
	iswc := values.Get("iswc")
	publisherId := values.Get("publisherId")
	title := values.Get("title")
	info, err := api.Compose(composerId, hfa, iswc, publisherId, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, info)
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
	compositionRightId := form.Value["compositionRightId"][0]
	file, err := form.File["recording"][0].Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	isrc := form.Value["isrc"][0]
	labelId := form.Value["labelId"][0]
	performerId := form.Value["performerId"][0]
	producerId := form.Value["producerId"][0]
	publicationId := form.Value["publicationId"][0]
	recording, err := api.Record(compositionRightId, file, isrc, labelId, performerId, producerId, publicationId)
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
	compositionId := values.Get("compositionId")
	rightIds := SplitStr(values.Get("rightIds"), ",")
	composition, err := api.Publish(compositionId, rightIds)
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
	mechanicalLicenseId := values.Get("mechanicalLicenseId")
	recordingId := values.Get("recordingId")
	rightIds := SplitStr(values.Get("rightIds"), ",")
	release, err := api.Release(mechanicalLicenseId, recordingId, rightIds)
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
	licenseeId := values.Get("licenseeId")
	publicationId := values.Get("publicationId")
	releaseId := values.Get("releaseId")
	rightId := values.Get("rightId")
	territory := SplitStr(values.Get("territory"), ",")
	validFrom := values.Get("validFrom")
	validTo := values.Get("validTo")
	switch {
	case !EmptyStr(publicationId):
		license, err = api.MechanicalLicense(licenseeId, publicationId, rightId, territory, validFrom, validTo)
	case !EmptyStr(releaseId):
		license, err = api.MasterLicense(licenseeId, releaseId, rightId, territory, validFrom, validTo)
	default:
		http.Error(w, "Expected publicationId or releaseId", http.StatusBadRequest)
		return
	}
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
		tx, err := bigchain.GetTx(publicationId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		publication := bigchain.GetTxData(tx)
		pub := bigchain.GetTxSigner(tx)
		model, err = ld.QueryPublicationField(field, publication, pub)
	case !EmptyStr(releaseId):
		tx, err := bigchain.GetTx(releaseId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		release := bigchain.GetTxData(tx)
		pub := bigchain.GetTxSigner(tx)
		model, err = ld.QueryReleaseField(field, release, pub)
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
	output := MustAtoi(values.Get("output"))
	percentageShares := MustAtoi(values.Get("percentageShares"))
	publicationId := values.Get("publicationId")
	recipientId := values.Get("recipientId")
	releaseId := values.Get("releaseId")
	rightId := values.Get("rightId")
	switch {
	case !EmptyStr(publicationId):
		transfer, err = api.TransferCompositionRight(output, percentageShares, publicationId, recipientId, rightId)
	case !EmptyStr(releaseId):
		transfer, err = api.TransferRecordingRight(output, percentageShares, recipientId, releaseId, rightId)
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
	case api.agent == nil:
		api.logger.Warn("Agent profile is not set")
	case api.agentId == "":
		api.logger.Warn("Agent ID is not set")
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

func (api *Api) Login(agentId, privstr string) error {
	priv := new(ed25519.PrivateKey)
	if err := priv.FromString(privstr); err != nil {
		return err
	}
	tx, err := bigchain.GetTx(agentId)
	if err != nil {
		return err
	}
	agent := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(agent); err != nil {
		return err
	}
	pub := bigchain.GetTxSigner(tx)
	if !pub.Equals(priv.Public()) {
		return ErrInvalidKey
	}
	api.agent = agent
	api.agentId = agentId
	api.priv = priv
	api.pub = pub
	agentName := spec.GetAgentName(agent)
	api.logger.Info(Sprintf("SUCCESS %s is logged in", agentName))
	return nil
}

type RegisterMessage struct {
	AgentId string `json:"agent_id"`
	PrivKey string `json:"private_key"`
	PubKey  string `json:"public_key"`
}

func (api *Api) Register(email, name, password, socialMedia string) (*RegisterMessage, error) {
	priv, pub := ed25519.GenerateKeypairFromPassword(password)
	agent := spec.NewAgent(email, name, socialMedia)
	if err := spec.ValidAgent(agent); err != nil {
		return nil, err
	}
	tx := bigchain.DefaultIndividualCreateTx(agent, pub)
	bigchain.FulfillTx(tx, priv)
	agentId, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS registered new agent: " + name)
	return &RegisterMessage{
		AgentId: agentId,
		PrivKey: priv.String(),
		PubKey:  pub.String(),
	}, nil
}

func (api *Api) Compose(composerId, hfa, iswc, publisherId, title string) (composition Data, err error) {
	composition = spec.NewComposition(composerId, hfa, iswc, publisherId, title)
	if err = ld.ValidateComposition(composition, api.pub); err != nil {
		return nil, err
	}
	tx := bigchain.DefaultIndividualCreateTx(composition, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	composition["id"], err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with composition")
	return composition, nil
}

func (api *Api) Record(compositionRightId string, file io.Reader, isrc, labelId, performerId, producerId, publicationId string) (recording Data, err error) {
	// rs := MustReadSeeker(file)
	// meta, err := tag.ReadFrom(rs)
	// if err != nil {
	// 	return nil, err
	// }
	// metadata := meta.Raw()
	recording = spec.NewRecording(compositionRightId, isrc, labelId, performerId, producerId, publicationId)
	if err = ld.ValidateRecording(recording, api.pub); err != nil {
		return nil, err
	}
	tx := bigchain.DefaultIndividualCreateTx(recording, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	recording["id"], err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with recording")
	return recording, nil
}

func (api *Api) Publish(compositionId string, rightIds []string) (publication Data, err error) {
	publication = spec.NewPublication(compositionId, rightIds)
	if err = ld.ValidatePublication(publication, api.pub); err != nil {
		panic(err)
		return nil, err
	}
	tx := bigchain.DefaultIndividualCreateTx(publication, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	publication["id"], err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with publication")
	return publication, nil
}

func (api *Api) Release(mechanicalLicenseId, recordingId string, rightIds []string) (release Data, err error) {
	release = spec.NewRelease(mechanicalLicenseId, recordingId, rightIds)
	if err = ld.ValidateRelease(release, api.pub); err != nil {
		return nil, err
	}
	tx := bigchain.DefaultIndividualCreateTx(release, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	release["id"], err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with release")
	return release, nil
}

func (api *Api) CompositionRight(compositionId string, percentageShares int, rightHolderId string, territory []string, validFrom, validTo string) (right Data, err error) {
	tx, err := bigchain.GetTx(rightHolderId)
	if err != nil {
		return nil, err
	}
	pub := bigchain.GetTxSigner(tx)
	right = spec.NewCompositionRight(compositionId, territory, validFrom, validTo)
	tx = bigchain.IndividualCreateTx(percentageShares, right, pub, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	right["id"], err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS created composition right")
	return right, nil
}

func (api *Api) RecordingRight(percentageShares int, recordingId, rightHolderId string, territory []string, validFrom, validTo string) (right Data, err error) {
	tx, err := bigchain.GetTx(rightHolderId)
	if err != nil {
		return nil, err
	}
	pub := bigchain.GetTxSigner(tx)
	right = spec.NewRecordingRight(recordingId, territory, validFrom, validTo)
	if err = spec.ValidRecordingRight(right); err != nil {
		return nil, err
	}
	tx = bigchain.IndividualCreateTx(percentageShares, right, pub, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	right["id"], err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS created recording right")
	return right, nil
}

func (api *Api) MechanicalLicense(licenseeId, publicationId, rightId string, territory []string, validFrom, validTo string) (license Data, err error) {
	license = spec.NewLicense(licenseeId, api.agentId, publicationId, "", rightId, territory, spec.LICENSE_MECHANICAL, validFrom, validTo)
	if err = ld.ValidateMechanicalLicense(license, api.pub); err != nil {
		return nil, err
	}
	tx := bigchain.DefaultIndividualCreateTx(license, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	license["id"], err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with mechanical license")
	return license, nil
}

func (api *Api) MasterLicense(licenseeId, releaseId, rightId string, territory []string, validFrom, validTo string) (license Data, err error) {
	license = spec.NewLicense(licenseeId, api.agentId, "", releaseId, rightId, territory, spec.LICENSE_MASTER, validFrom, validTo)
	if err = ld.ValidateMasterLicense(license, api.pub); err != nil {
		return nil, err
	}
	tx := bigchain.DefaultIndividualCreateTx(license, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	license["id"], err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with master license")
	return license, nil
}

// At the moment, only supports transfers between one party and another
// Output should be 0 if it's an original right or right was self-assigned; 1 if we received right

func (api *Api) TransferCompositionRight(output, percentageShares int, publicationId, recipientId, rightId string) (Data, error) {
	rightTx, originalRightId, err := ld.ValidateCompositionRightHolder(publicationId, api.agentId, rightId)
	if err != nil {
		return nil, err
	}
	tx, err := bigchain.GetTx(recipientId)
	if err != nil {
		return nil, err
	}
	pub := bigchain.GetTxSigner(tx)
	totalShares := bigchain.GetOutputAmount(bigchain.GetTxOutput(rightTx, output))
	keepShares := totalShares - percentageShares
	if keepShares < 0 {
		return nil, ErrorAppend(ErrCriteriaNotMet, "agent cannot transfer this many shares")
	}
	if keepShares == 0 {
		tx = bigchain.IndividualTransferTx(percentageShares, originalRightId, rightId, output, pub, api.pub)
	} else {
		tx = bigchain.DivisibleTransferTx([]int{keepShares, percentageShares}, originalRightId, rightId, output, []crypto.PublicKey{api.pub, pub}, api.pub)
	}
	bigchain.FulfillTx(tx, api.priv)
	if _, err = bigchain.PostTx(tx); err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent composition right transfer tx")
	return tx, nil
}

func (api *Api) TransferRecordingRight(output, percentageShares int, recipientId, releaseId, rightId string) (Data, error) {
	rightTx, originalRightId, err := ld.ValidateRecordingRightHolder(releaseId, api.agentId, rightId)
	if err != nil {
		return nil, err
	}
	tx, err := bigchain.GetTx(recipientId)
	if err != nil {
		return nil, err
	}
	pub := bigchain.GetTxSigner(tx)
	totalShares := bigchain.GetOutputAmount(bigchain.GetTxOutput(rightTx, output))
	keepShares := totalShares - percentageShares
	if keepShares < 0 {
		return nil, ErrorAppend(ErrCriteriaNotMet, "agent cannot transfer this many shares")
	}
	if keepShares == 0 {
		tx = bigchain.IndividualTransferTx(percentageShares, originalRightId, rightId, output, pub, api.pub)
	} else {
		tx = bigchain.DivisibleTransferTx([]int{percentageShares, keepShares}, originalRightId, rightId, output, []crypto.PublicKey{pub, api.pub}, api.pub)
	}
	bigchain.FulfillTx(tx, api.priv)
	if _, err = bigchain.PostTx(tx); err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent recording right transfer tx")
	return tx, nil
}
