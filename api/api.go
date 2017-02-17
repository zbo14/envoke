package api

import (
	"io"
	"net/http"
	"time"

	"github.com/dhowden/tag"
	"github.com/zbo14/envoke/bigchain"
	// "github.com/zbo14/envoke/chroma"
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
	mux.HandleFunc("/composition", api.CompositionHandler)
	mux.HandleFunc("/recording", api.RecordingHandler)
	mux.HandleFunc("/license", api.LicenseHandler)
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
	agentId := values.Get("agent_id")
	privstr := values.Get("private_key")
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
	numRights, err := Atoi(values.Get("numRights"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rights := make([]Data, numRights)
	for i := 0; i < numRights; i++ {
		n := Itoa(i)
		context := SplitStr(values.Get("context"+n), ",")
		exclusive, err := ParseBool(values.Get("exclusive" + n))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		percentageShares := values.Get("percentageShares" + n)
		rightHolderId := values.Get("rightHolderId" + n)
		usage := SplitStr(values.Get("usage"+n), ",")
		validFrom, err := ParseDateStr(values.Get("validFrom" + n))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		validTo, err := ParseDateStr(values.Get("validTo" + n))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rights[i], err = api.Right(context, exclusive, percentageShares, rightHolderId, usage, validFrom, validTo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	composerId := values.Get("composerId")
	publisherId := values.Get("publisherId")
	title := values.Get("title")
	composition, err := api.Composition(composerId, publisherId, rights, title)
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
	form, err := MultipartForm(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	numRights, err := Atoi(form.Value["numRights"][0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rights := make([]Data, numRights)
	for i := 0; i < numRights; i++ {
		n := Itoa(i)
		context := SplitStr(form.Value["context"+n][0], ",")
		exclusive, err := ParseBool(form.Value["exclusive"+n][0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		percentageShares := form.Value["percentageShares"+n][0]
		rightHolderId := form.Value["rightHolderId"+n][0]
		usage := SplitStr(form.Value["usage"+n][0], ",")
		validFrom, err := ParseDateStr(form.Value["validFrom"+n][0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		validTo, err := ParseDateStr(form.Value["validTo"+n][0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rights[i], err = api.Right(context, exclusive, percentageShares, rightHolderId, usage, validFrom, validTo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
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
	recording, err := api.Recording(compositionId, file, labelId, performerId, producerId, publishingLicenseId, rights)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, recording)
}

func (api *Api) LicenseHandler(w http.ResponseWriter, req *http.Request) {
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
	recordingId := values.Get("recordingId")
	validFrom, err := ParseDateStr(values.Get("validFrom"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	validTo, err := ParseDateStr(values.Get("validTo"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var license Data
	switch licenseType {
	case
		spec.LICENSE_TYPE_MECHANICAL,
		spec.LICENSE_TYPE_SYNCHRONIZATION:
		license, err = api.PublishingLicense(compositionId, licenseeId, licenseType, validFrom, validTo)
	case spec.LICENSE_TYPE_MASTER:
		license, err = api.RecordingLicense(licenseeId, licenseType, recordingId, validFrom, validTo)
	default:
		http.Error(w, ErrorAppend(ErrInvalidType, licenseType).Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, license)
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

// should we do login or just registration via partner org?
// having artist identity on envoke might ease attribution
// e.g. album/track contains uri to artist profile in db
// but artist must be verified by partner org before they
// create profile..

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
	composition = spec.NewComposition(composerId, publisherId, rights, title)
	if _, err = ld.ValidateComposition(composition, api.pub); err != nil {
		return nil, err
	}
	tx := bigchain.GenerateTx(composition, nil, bigchain.CREATE, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	composition["id"], err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with composition")
	return composition, nil
}

func (api *Api) Recording(compositionId string, file io.Reader, labelId, performerId, producerId, publishingLicenseId string, rights []Data) (Data, error) {
	s, _ := MustTeeSeeker(file) //r
	meta, err := tag.ReadFrom(s)
	if err != nil {
		return nil, err
	}
	metadata := meta.Raw()
	// fingerprint, err := chroma.NewFingerprint(r)
	fingerprint := "V0VHa09XR0xXb2VnbGt3ZW93ZWZ3ZUZ3ZWZ3Z3dlZ2VnZ2VyZ2U"
	if err != nil {
		return nil, err
	}
	recording := spec.NewRecording(compositionId, fingerprint, labelId, performerId, producerId, publishingLicenseId, rights)
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

func (api *Api) Right(context []string, exclusive bool, percentageShares, rightHolderId string, usage []string, validFrom, validTo time.Time) (Data, error) {
	right := spec.NewRight(context, exclusive, percentageShares, rightHolderId, usage, validFrom, validTo)
	if err := spec.ValidRight(right); err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS created right")
	return right, nil
}

func (api *Api) PublishingLicense(compositionId, licenseeId, licenseType string, validFrom, validTo time.Time) (license Data, err error) {
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

func (api *Api) RecordingLicense(licenseeId, licenseType, recordingId string, validFrom, validTo time.Time) (license Data, err error) {
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
