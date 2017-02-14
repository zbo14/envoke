package api

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/dhowden/tag"
	"github.com/zbo14/envoke/bigchain"
	"github.com/zbo14/envoke/chroma"
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
	mux.HandleFunc("/album", api.AlbumHandler)
	mux.HandleFunc("/track", api.TrackHandler)
	mux.HandleFunc("/right", api.RightHandler)
	mux.HandleFunc("/sign", api.SignHandler)
	mux.HandleFunc("/verify", api.VerifyHandler)
}

func (api *Api) LoginHandler(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	// Get request values
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	agentId := values.Get("agent_id")
	privstr := values.Get("private_key")
	_type := values.Get("type")
	// Login
	loginMessage, err := api.Login(agentId, privstr, _type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(loginMessage))
}

func (api *Api) RegisterHandler(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	// Get request values
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Register new user
	registerMessage, err := api.Register(values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, registerMessage)
}

func (api *Api) AlbumHandler(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	// Check that we're logged in
	if !api.LoggedIn() {
		http.Error(w, ErrInvalidLogin.Error(), http.StatusUnauthorized)
		return
	}
	// Multipart form
	form, err := MultipartForm(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	title := form.Value["album_title"][0]
	labelId := form.Value["label_id"][0]
	publisherId := form.Value["publisher_id"][0]
	tracks := form.File["tracks"]
	// Extract metadata from tracks, send album to Bigchain/IPDB
	albumMessage, err := api.Album(labelId, publisherId, title, tracks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, albumMessage)
}

func (api *Api) RightHandler(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	// Check that we're logged in
	if !api.LoggedIn() {
		http.Error(w, ErrInvalidLogin.Error(), http.StatusUnauthorized)
		return
	}
	// Get request values
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	values.Set("issuer_id", api.agentId)
	// Music Rights
	rightMessage, err := api.Right(values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, rightMessage)
}

func (api *Api) TrackHandler(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	// Check that we're logged in
	if !api.LoggedIn() {
		http.Error(w, ErrInvalidLogin.Error(), http.StatusUnauthorized)
		return
	}
	// Multipart form
	form, err := MultipartForm(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	labelId := form.Value["label_id"][0]
	publisherId := form.Value["publisher_id"][0]
	file, err := form.File["tracks"][0].Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Extract track metadata and send to BigchainDB/IPDB
	trackMessage, err := api.Track("", file, labelId, 0, publisherId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, trackMessage)
}

func (api *Api) SignHandler(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	// Check that we're logged in
	if !api.LoggedIn() {
		http.Error(w, ErrInvalidLogin.Error(), http.StatusUnauthorized)
		return
	}
	// Get request values
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	rightId := values.Get("right_id") //signing rights for now
	signMessage, err := api.Sign(rightId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, signMessage)
}

func (api *Api) VerifyHandler(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	// Get request values
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	signatureId := values.Get("signature_id")
	// Verify linked-data signature
	verifyMessage, err := api.Verify(signatureId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, verifyMessage)
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

func (api *Api) Login(agentId, privstr, _type string) (string, error) {
	// PrivKey
	priv := new(ed25519.PrivateKey)
	if err := priv.FromString(privstr); err != nil {
		return "", err
	}
	// Query tx with agent id
	tx, err := bigchain.GetTx(agentId)
	if err != nil {
		return "", err
	}
	// Validate agent
	agent := bigchain.GetTxData(tx)
	if !spec.ValidAgentWithType(agent, _type) {
		return "", ErrorAppend(ErrInvalidModel, _type)
	}
	// Check that privkey matches agent pubkey
	pub := spec.GetAgentPublicKey(agent)
	if !bytes.Equal(priv.Public().Bytes(), pub.Bytes()) {
		return "", ErrInvalidKey
	}
	api.agent = agent
	api.agentId = agentId
	api.priv = priv
	api.pub = pub
	agentName := spec.GetAgentName(agent)
	api.logger.Info("SUCCESS you are logged in")
	return NewLoginMessage(agentName), nil
}

func (api *Api) Register(values url.Values) (*RegisterMessage, error) {
	// Generate keypair from password
	password := values.Get("password")
	priv, pub := ed25519.GenerateKeypairFromPassword(password)
	values.Set("public_key", pub.String())
	// New agent
	agent, err := AgentFromValues(values)
	if err != nil {
		return nil, err
	}
	// Check that we created a valid agent
	if !spec.ValidAgent(agent) {
		return nil, ErrorAppend(ErrInvalidModel, spec.GetType(agent))
	}
	tx := bigchain.GenerateTx(agent, nil, bigchain.CREATE, pub)
	bigchain.FulfillTx(tx, priv)
	// send request to IPDB
	agentId, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	_type := spec.GetType(agent)
	api.logger.Info("SUCCESS registered new " + _type)
	return NewRegisterMessage(agentId, priv.String(), pub.String()), nil
}

func (api *Api) Album(labelId, publisherId, title string, tracks []*multipart.FileHeader) (*AlbumMessage, error) {
	// New album
	album := spec.NewAlbum(api.agentId, labelId, publisherId, title)
	// Check that we generated a valid linked-data album
	if err := ld.ValidateAlbum(album); err != nil {
		return nil, err
	}
	// Send tx with album
	tx := bigchain.GenerateTx(album, nil, bigchain.CREATE, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	albumId, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with album")
	// It would be great if we could batch write tracks
	trackIds := make([]string, len(tracks))
	for i, track := range tracks {
		file, err := track.Open()
		if err != nil {
			return nil, err
		}
		trackMessage, err := api.Track(albumId, file, "", i+1, "")
		if err != nil {
			return nil, err
		}
		trackIds[i] = trackMessage.TrackId
	}
	return NewAlbumMessage(albumId, trackIds), nil
}

func (api *Api) Track(albumId string, file multipart.File, labelId string, number int, publisherId string) (*TrackMessage, error) {
	s, r := MustTeeSeeker(file)
	// Extract metadata
	meta, err := tag.ReadFrom(s)
	if err != nil {
		return nil, err
	}
	metadata := meta.Raw()
	trackTitle := meta.Title()
	// Get acoustic fingerprint
	fingerprint, err := chroma.NewFingerprint(r)
	if err != nil {
		return nil, err
	}
	// New track
	track := spec.NewTrack(albumId, api.agentId, fingerprint, labelId, number, publisherId, trackTitle)
	// Check that we generated a valid linked-data track
	if err := ld.ValidateTrack(track); err != nil {
		return nil, err
	}
	// Send tx with track
	tx := bigchain.GenerateTx(track, metadata, bigchain.CREATE, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	trackId, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with track")
	return NewTrackMessage(trackId), nil
}

func (api *Api) Right(values url.Values) (*RightMessage, error) {
	musicId := values.Get("music_id")
	tx, err := bigchain.GetTx(musicId)
	if err != nil {
		return nil, err
	}
	music := bigchain.GetTxData(tx)
	sig := api.priv.Sign(MustMarshalJSON(music))
	context := Split(values.Get("context"), ",")
	issuerId := values.Get("issuer_id")
	issuerType := values.Get("issuer_type")
	percentageShares := values.Get("percentage_shares")
	recipientId := values.Get("recipient_id")
	usage := Split(values.Get("usage"), ",")
	validFrom, err := ParseDateStr(values.Get("valid_from"))
	if err != nil {
		return nil, err
	}
	validTo, err := ParseDateStr(values.Get("valid_to"))
	if err != nil {
		return nil, err
	}
	// New right
	right := spec.NewRight(context, issuerId, issuerType, musicId, percentageShares, recipientId, sig, usage, validFrom, validTo)
	// Check that we generated a valid linked-data right
	if err = ld.ValidateRight(right); err != nil {
		return nil, err
	}
	tx = bigchain.GenerateTx(right, nil, bigchain.CREATE, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	rightId, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with right")
	return NewRightMessage(rightId), nil
}

func (api *Api) Sign(modelId string) (*SignMessage, error) {
	// Check that model id matches regex
	if !MatchString(spec.ID_REGEX, modelId) {
		return nil, ErrInvalidId
	}
	// Validate linked-data model
	model, err := ld.ValidateModelId(modelId)
	if err != nil {
		return nil, err
	}
	// Check that agent, model meet criteria
	if !MeetsCriteria(api.agentId, model) {
		return nil, ErrCriteriaNotMet
	}
	sig := api.priv.Sign(MustMarshalJSON(model))
	signature := spec.NewSignature(modelId, api.agentId, sig)
	// Check that we generated a valid linked-data signature
	if err := ld.ValidateSignature(signature); err != nil {
		return nil, err
	}
	// Send tx with signature to IPDB
	tx := bigchain.GenerateTx(signature, nil, bigchain.CREATE, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	signatureId, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with signature")
	return NewSignMessage(signatureId), nil
}

func (api *Api) Verify(signatureId string) (*VerifyMessage, error) {
	// Validate linked-data signature
	signature, err := ld.ValidateSignatureId(signatureId)
	if err != nil {
		api.logger.Info("FAILURE could not verify signature")
		return NewVerifyMessage(err.Error(), nil, false), nil
	}
	agentId := spec.GetSignatureSigner(signature)
	modelId := spec.GetSignatureModel(signature)
	tx, err := bigchain.GetTx(modelId)
	if err != nil {
		return nil, err
	}
	model := bigchain.GetTxData(tx)
	// Check that agent, model meet criteria
	if !MeetsCriteria(agentId, model) {
		api.logger.Info("FAILURE could not verify signature")
		return NewVerifyMessage(ErrCriteriaNotMet.Error(), nil, false), nil
	}
	api.logger.Info("SUCCESS verified signature")
	return NewVerifyMessage("", signature, true), nil
}

func AgentFromValues(values url.Values) (Data, error) {
	email := values.Get("email")
	name := values.Get("name")
	pub := new(ed25519.PublicKey)
	pubstr := values.Get("public_key")
	if err := pub.FromString(pubstr); err != nil {
		return nil, err
	}
	switch values.Get("type") {
	case spec.ARTIST:
		return spec.NewArtist(email, name, pub), nil
	case spec.LABEL:
		return spec.NewLabel(email, name, pub), nil
	case spec.ORGANIZATION:
		return spec.NewOrganization(email, name, pub), nil
	case spec.PUBLISHER:
		return spec.NewPublisher(email, name, pub), nil
		// TODO: add more partner types?
	}
	return nil, ErrInvalidType
}

func MeetsCriteria(agentId string, model Data) bool {
	if spec.IsRight(model) {
		recipientId := spec.GetRightRecipient(model)
		return agentId == recipientId
	}
	return false
}
