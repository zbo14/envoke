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
	mux.HandleFunc("/right", api.RightHandler)
	mux.HandleFunc("/composition_info", api.CompositionInfoHandler)
	mux.HandleFunc("/recording_info", api.RecordingInfoHandler)
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
	infoId := values.Get("infoId")
	percentageShares := values.Get("percentageShares")
	validFrom := values.Get("validFrom")
	validTo := values.Get("validTo")
	right, err := api.Right(infoId, percentageShares, validFrom, validTo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, right)
}

func (api *Api) CompositionInfoHandler(w http.ResponseWriter, req *http.Request) {
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
	publisherId := values.Get("publisherId")
	title := values.Get("title")
	info, err := api.CompositionInfo(composerId, publisherId, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, info)
}

func (api *Api) RecordingInfoHandler(w http.ResponseWriter, req *http.Request) {
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
	file, err := form.File["recording"][0].Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	labelId := form.Value["labelId"][0]
	performerId := form.Value["performerId"][0]
	producerId := form.Value["producerId"][0]
	publishingLicenseId := form.Value["publishingLicenseId"][0]
	info, err := api.RecordingInfo(compositionId, file, labelId, performerId, producerId, publishingLicenseId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, info)
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
	infoId := values.Get("infoId")
	rightIds := SplitStr(values.Get("rightIds"), ",")
	composition, err := api.Composition(infoId, rightIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, composition)
}

func (api *Api) RecordingHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	infoId := values.Get("infoId")
	rightIds := SplitStr(values.Get("rightIds"), ",")
	recording, err := api.Recording(infoId, rightIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
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
	tx := bigchain.IndividualCreateTx(agent, pub)
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

func (api *Api) CompositionInfo(composerId, publisherId, title string) (info Data, err error) {
	info = spec.NewCompositionInfo(composerId, publisherId, title)
	if _, err = ld.ValidateCompositionInfo(info, api.pub); err != nil {
		return nil, err
	}
	tx := bigchain.IndividualCreateTx(info, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	info["id"], err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with composition info")
	return info, nil
}

func (api *Api) RecordingInfo(compositionId string, file io.Reader, labelId, performerId, producerId, publishingLicenseId string) (info Data, err error) {
	// rs := MustReadSeeker(file)
	// meta, err := tag.ReadFrom(rs)
	// if err != nil {
	// 	return nil, err
	// }
	// metadata := meta.Raw()
	info = spec.NewRecordingInfo(compositionId, labelId, performerId, producerId, publishingLicenseId)
	if _, err = ld.ValidateRecordingInfo(info, api.pub); err != nil {
		return nil, err
	}
	tx := bigchain.IndividualCreateTx(info, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	info["id"], err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with recording info")
	return info, nil
}

func (api *Api) Composition(infoId string, rightIds []string) (composition Data, err error) {
	composition = spec.NewComposition(infoId, rightIds)
	if _, err = ld.ValidateComposition(composition, api.pub); err != nil {
		return nil, err
	}
	tx := bigchain.IndividualCreateTx(composition, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	composition["id"], err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with composition")
	return composition, nil
}

func (api *Api) Recording(infoId string, rightIds []string) (recording Data, err error) {
	recording = spec.NewRecording(infoId, rightIds)
	if _, err = ld.ValidateRecording(recording, api.pub); err != nil {
		return nil, err
	}
	tx := bigchain.IndividualCreateTx(recording, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	recording["id"], err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with recording")
	return recording, nil
}

func (api *Api) Right(infoId, percentageShares, validFrom, validTo string) (right Data, err error) {
	right = spec.NewRight(infoId, percentageShares, validFrom, validTo)
	if err = spec.ValidRight(right); err != nil {
		return nil, err
	}
	tx := bigchain.IndividualCreateTx(right, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	right["id"], err = bigchain.PostTx(tx)
	if err != nil {
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
	tx := bigchain.IndividualCreateTx(license, api.pub)
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
	tx := bigchain.IndividualCreateTx(license, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	license["id"], err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with recording license")
	return license, nil
}
