package api

import (
	"io"
	"net/http"

	"github.com/dhowden/tag"
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
	rights  []Data
}

func NewApi() *Api {
	return &Api{
		logger: NewLogger("api"),
	}
}

func (api *Api) AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/login", api.LoginHandler)
	mux.HandleFunc("/register", api.RegisterHandler)
	mux.HandleFunc("/right", api.RightHandler)
	mux.HandleFunc("/composition", api.CompositionHandler)
	mux.HandleFunc("/recording", api.RecordingHandler)
	mux.HandleFunc("/publishing_license", api.PublishingLicenseHandler)
	mux.HandleFunc("/recording_license", api.RecordingLicenseHandler)
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
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	percentageShares := values.Get("percentageShares")
	// rightHolderId := values.Get("rightHolderId")
	validFrom := values.Get("validFrom")
	validTo := values.Get("validTo")
	right, err := api.Right(percentageShares, validFrom, validTo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	api.rights = append(api.rights, right)
	WriteJSON(w, right)
}

func (api *Api) CompositionHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if api.rights == nil {
		http.Error(w, ErrorAppend(ErrCriteriaNotMet, "no composition rights").Error(), http.StatusBadRequest)
		return
	}
	composerId := values.Get("composerId")
	publisherId := values.Get("publisherId")
	title := values.Get("title")
	composition, err := api.Composition(composerId, publisherId, api.rights, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	api.rights = nil
	WriteJSON(w, composition)
}

func (api *Api) RecordingHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	form, err := MultipartForm(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if api.rights == nil {
		http.Error(w, ErrorAppend(ErrCriteriaNotMet, "no recording rights").Error(), http.StatusBadRequest)
		return
	}
	compositionId := form.Value["compositionId"][0]
	file, err := form.File["recording"][0].Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	labelId := form.Value["labelId"][0]
	performerId := form.Value["performerId"][0]
	producerId := form.Value["producerId"][0]
	publishingLicenseId := form.Value["publishingLicenseId"][0]
	recording, err := api.Recording(compositionId, file, labelId, performerId, producerId, publishingLicenseId, api.rights)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	api.rights = nil
	WriteJSON(w, recording)
}

func (api *Api) PublishingLicenseHandler(w http.ResponseWriter, req *http.Request) {
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
	licenseeId := values.Get("licenseeId")
	licenseType := values.Get("licenseType")
	validFrom := values.Get("validFrom")
	validTo := values.Get("validTo")
	license, err := api.PublishingLicense(compositionId, licenseeId, licenseType, validFrom, validTo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, license)
}

func (api *Api) RecordingLicenseHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	licenseeId := values.Get("licenseeId")
	licenseType := values.Get("licenseType")
	recordingId := values.Get("recordingId")
	validFrom := values.Get("validFrom")
	validTo := values.Get("validTo")
	license, err := api.RecordingLicense(licenseeId, licenseType, recordingId, validFrom, validTo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, license)
}

func (api *Api) SearchHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	field := values.Get("field")
	modelId := values.Get("modelId")
	model, err := ld.QueryModelIdField(field, modelId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, model)
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
	pub := bigchain.GetTxPublicKey(tx)
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
	tx := bigchain.GenerateTx(agent, nil, bigchain.CREATE, pub)
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

func (api *Api) Composition(composerId, publisherId string, rights []Data, title string) (composition Data, err error) {
	rightIds := make([]string, len(rights))
	for i, right := range rights {
		tx := bigchain.CreateTx(right, api.pub)
		bigchain.FulfillTx(tx, priv)
		rightIds[i], err = bigchain.PostTx(tx)
		if err != nil {
			return nil, err
		}
	}
	composition = spec.NewComposition(composerId, publisherId, rightIds, title)
	if _, err = ld.ValidateComposition(composition, api.pub); err != nil {
		return nil, err
	}
	tx := bigchain.CreateTx(composition, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	composition["id"], err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with composition")
	return composition, nil
}

func (api *Api) Recording(compositionId string, file io.Reader, labelId, performerId, producerId, publishingLicenseId string, rights []Data) (Data, error) {
	rs := MustReadSeeker(file)
	meta, err := tag.ReadFrom(rs)
	if err != nil {
		return nil, err
	}
	metadata := meta.Raw()
	recording := spec.NewRecording(compositionId, labelId, performerId, producerId, publishingLicenseId, rights)
	if _, err = ld.ValidateRecording(recording, api.pub); err != nil {
		return nil, err
	}
	tx := bigchain.GenerateTx(recording, metadata, bigchain.CREATE, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	recording["id"], err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with recording")
	return recording, nil
}

func (api *Api) Right(percentageShares, validFrom, validTo string) (Data, error) {
	right := spec.NewRight(percentageShares, validFrom, validTo)
	if err := spec.ValidRight(right); err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS created right")
	return right, nil
}

func (api *Api) PublishingLicense(compositionId, licenseeId, licenseType, validFrom, validTo string) (license Data, err error) {
	license = spec.NewPublishingLicense(compositionId, licenseeId, api.agentId, licenseType, validFrom, validTo)
	if _, err = ld.ValidatePublishingLicense(license, api.pub); err != nil {
		return nil, err
	}
	tx := bigchain.GenerateTx(license, nil, bigchain.CREATE, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	license["id"], err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with publishing license")
	return license, nil
}

func (api *Api) RecordingLicense(licenseeId, licenseType, recordingId, validFrom, validTo string) (license Data, err error) {
	license = spec.NewRecordingLicense(licenseeId, api.agentId, licenseType, recordingId, validFrom, validTo)
	if _, err = ld.ValidateRecordingLicense(license, api.pub); err != nil {
		return nil, err
	}
	tx := bigchain.GenerateTx(license, nil, bigchain.CREATE, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	license["id"], err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with recording license")
	return license, nil
}
