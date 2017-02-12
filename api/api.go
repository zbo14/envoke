package api

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/url"

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
	publisherId := form.Value["publisher_id"][0]
	tracks := form.File["tracks"]
	// Extract metadata from tracks, send album to Bigchain/IPDB
	albumMessage, err := api.Album(publisherId, title, tracks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, albumMessage)
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
	publisherId := form.Value["publisher_id"][0]
	file, err := form.File["tracks"][0].Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Extract track metadata and send to BigchainDB/IPDB
	trackMessage, err := api.Track("", file, 0, publisherId)
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
	musicId := values.Get("music_id")
	// Linked-data signature
	txInfo, err := api.Sign(musicId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, txInfo)
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
	verifyMessage := api.Verify(signatureId)
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
	tx := bigchain.GenerateTx(agent, nil, bigchain.CREATE, pub)
	bigchain.FulfillTx(tx, priv)
	// send request to IPDB
	agentId, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	_type := spec.GetEntityType(agent)
	api.logger.Info("SUCCESS registered new " + _type)
	return NewRegisterMessage(agentId, priv.String(), pub.String()), nil
}

func (api *Api) Album(publisherId, title string, tracks []*multipart.FileHeader) (*AlbumMessage, error) {
	// Generate and send tx with album
	album := spec.NewAlbum(api.agentId, publisherId, title)
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
		trackMessage, err := api.Track(albumId, file, i+1, "")
		if err != nil {
			return nil, err
		}
		trackIds[i] = trackMessage.TrackId
	}
	return NewAlbumMessage(albumId, trackIds), nil
}

func (api *Api) Track(albumId string, file multipart.File, number int, publisherId string) (*TrackMessage, error) {
	// Extract metadata
	meta, err := tag.ReadFrom(file)
	if err != nil {
		return nil, err
	}
	metadata := meta.Raw()
	trackTitle := meta.Title()
	// Create new track
	track := spec.NewTrack(albumId, api.agentId, number, publisherId, trackTitle)
	// Generate and send tx with track
	tx := bigchain.GenerateTx(track, metadata, bigchain.CREATE, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	trackId, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with track")
	return NewTrackMessage(trackId), nil
}

func (api *Api) Sign(musicId string) (*SignMessage, error) {
	// Validate linked music
	music, err := ld.ValidateLdMusicId(musicId)
	if err != nil {
		return nil, err
	}
	// Check that user agentId == music publisher_id
	if api.agentId != spec.GetMusicPublisher(music) {
		return nil, ErrInvalidId
	}
	json := MustMarshalJSON(music)
	sig := api.priv.Sign(json)
	signature := spec.NewSignature(musicId, api.agentId, sig)
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

func (api *Api) Verify(signatureId string) *VerifyMessage {
	signature, err := ld.ValidateLdSignatureId(signatureId)
	if err != nil {
		api.logger.Info("FAILURE could not verify signature")
		return NewVerifyMessage(err.Error(), nil, false)
	}
	api.logger.Info("SUCCESS verified signature")
	return NewVerifyMessage("", signature, true)
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
