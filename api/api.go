package api

import (
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/dhowden/tag"
	"github.com/zbo14/envoke/bigchain"
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/crypto"
	"github.com/zbo14/envoke/crypto/ed25519"
	"github.com/zbo14/envoke/spec/core"
)

type Api struct {
	agent   *core.Agent
	agentId string
	logger  Logger
	priv    crypto.PrivateKey
}

func NewApi() *Api {
	return &Api{
		logger: NewLogger("api"),
	}
}

func (api *Api) AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/login", api.Login)
	mux.HandleFunc("/register", api.Register)
	mux.HandleFunc("/album", api.Album)
	mux.HandleFunc("/track", api.Track)
	mux.HandleFunc("/sign", api.Sign)
	mux.HandleFunc("/verify", api.Verify)
}

func (api *Api) Album(w http.ResponseWriter, req *http.Request) {
	api.HandleAction(w, req, api.AlbumFromRequest)
}

func (api *Api) Track(w http.ResponseWriter, req *http.Request) {
	api.HandleAction(w, req, api.TrackFromRequest)
}

func (api *Api) Sign(w http.ResponseWriter, req *http.Request) {
	api.HandleAction(w, req, api.SignatureFromRequest)
}

func (api *Api) HandleAction(w http.ResponseWriter, req *http.Request, handler http.HandlerFunc) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
		return
	}
	// Make sure we're logged in properly
	if api.priv == nil {
		http.Error(w, "Privkey is not set", http.StatusUnauthorized)
		return
	}
	if api.agent == nil {
		http.Error(w, "Agent profile is not set", http.StatusUnauthorized)
		return
	}
	if api.agentId == "" {
		http.Error(w, "Agent ID is not set", http.StatusUnauthorized)
		return
	}
	handler(w, req)
}

// should we do login or just registration via partner org?
// having artist identity on envoke might ease attribution
// e.g. album/track contains uri to artist profile in db
// but artist must be verified by partner org before they
// create profile..

func (api *Api) Login(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	// Get request data
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	// PrivKey
	priv := new(ed25519.PrivateKey)
	privstr := values.Get("private_key")
	if err := priv.FromString(privstr); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// Query tx with agent id
	agentId := values.Get("agent_id")
	tx, err := bigchain.GetTx(agentId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	agent, ok := tx.GetData().(*core.Agent)
	if !ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_type := values.Get("type")
	if agent.Type != _type {
		http.Error(w, ErrInvalidType.Error(), http.StatusBadRequest)
		return
	}
	api.agent = agent
	api.agentId = agentId
	api.priv = priv
	w.Write([]byte("Login successful!"))
}

func (api *Api) Register(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	// Get request data
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	// Generate keypair from password
	password := values.Get("password")
	priv, pub := ed25519.GenerateKeypairFromPassword(password)
	values.Set("public_key", pub.String())
	// New agent
	agent, err := AgentFromValues(values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tx := bigchain.GenerateTx(agent, nil, bigchain.CREATE, pub)
	tx.Fulfill(priv)
	// send request to IPDB
	id, err := tx.Post()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, NewAgentInfo(id, priv.String(), pub.String()))
}

func (api *Api) SignatureFromRequest(w http.ResponseWriter, req *http.Request) {
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	musicId := values.Get("music_id")
	// Query IPDB
	tx, err := bigchain.GetTx(musicId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	music := tx.GetData()
	publisherId, err := core.GetMusicPublisher(music)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if api.agentId != publisherId {
		http.Error(w, ErrInvalidId.Error(), http.StatusUnauthorized)
		return
	}
	json := MustMarshalJSON(music)
	sig := api.priv.Sign(json)
	signature := core.NewSignature(api.agentId, sig)
	// Send tx with signature to IPDB
	tx = bigchain.GenerateTx(signature, nil, bigchain.CREATE, api.agent.PubKey)
	tx.Fulfill(api.priv)
	id, err := tx.Post()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	api.logger.Info("SUCCESS sent tx with signature")
	WriteJSON(w, NewActionInfo(id, core.SIGNATURE))
}

func (api *Api) SendTrack(albumId string, file multipart.File, number int, publisherId string) (string, error) {
	// Extract metadata
	meta, err := tag.ReadFrom(file)
	if err != nil {
		return "", err
	}
	metadata := meta.Raw()
	trackTitle := meta.Title()
	// Create new track
	track := core.NewTrack(albumId, api.agentId, number, publisherId, trackTitle)
	// Generate and send tx with track
	tx := bigchain.GenerateTx(track, metadata, bigchain.CREATE, api.agent.PubKey)
	tx.Fulfill(api.priv)
	trackId, err := tx.Post()
	if err != nil {
		return "", err
	}
	return trackId, nil
}

func (api *Api) TrackFromRequest(w http.ResponseWriter, req *http.Request) {
	form, err := MultipartForm(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	publisherId := form.Value["publisher_id"][0]
	tracks := form.File["tracks"]
	file, err := tracks[0].Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	trackId, err := api.SendTrack("", file, 0, publisherId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	api.logger.Info("SUCCESS sent tx with track")
	WriteJSON(w, NewActionInfo(trackId, core.TRACK))
	tx, err := bigchain.GetTx(trackId)
	Check(err)
	Println(string(MustMarshalIndentJSON(tx)))
}

func (api *Api) AlbumFromRequest(w http.ResponseWriter, req *http.Request) {
	form, err := MultipartForm(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	albumTitle := form.Value["album_title"][0]
	publisherId := form.Value["publisher_id"][0]
	album := core.NewAlbum(api.agentId, publisherId, albumTitle)
	// Generate and send tx with album
	tx := bigchain.GenerateTx(album, nil, bigchain.CREATE, api.agent.PubKey)
	tx.Fulfill(api.priv)
	albumId, err := tx.Post()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	api.logger.Info("SUCCESS sent tx with album")
	albumInfo := NewActionInfo(albumId, core.ALBUM)
	WriteJSON(w, albumInfo)
	tracks := form.File["tracks"]
	// It would be great if we could batch write tracks
	for i, track := range tracks {
		file, err := track.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		trackId, err := api.SendTrack(albumId, file, i+1, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		WriteJSON(w, NewActionInfo(trackId, core.TRACK))
	}
}

func (api *Api) Verify(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	musicId := values.Get("music_id")
	signatureId := values.Get("signature_id")
	tx, err := bigchain.GetTx(musicId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	music := tx.GetData()
	publisherId, err := core.GetMusicPublisher(music)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if publisherId != "" {
		api.VerifyMusic(music, signatureId, w)
		return
	}
	albumId := music.(*core.Track).AlbumId
	if albumId == "" {
		WriteJSON(w, NewQueryResult(ErrInvalidId.Error(), false))
		return
	}
	tx, err = bigchain.GetTx(albumId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	album := tx.GetData()
	api.VerifyMusic(album, signatureId, w)
}

func (api *Api) VerifyMusic(music interface{}, signatureId string, w http.ResponseWriter) {
	publisherId, err := core.GetMusicPublisher(music)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tx, err := bigchain.GetTx(signatureId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	signature := tx.GetData().(*core.Signature)
	signerId := signature.SignerId
	if publisherId != signerId {
		WriteJSON(w, NewQueryResult(ErrInvalidId.Error(), false))
		return
	}
	tx, err = bigchain.GetTx(publisherId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	publisher := tx.GetData().(*core.Agent)
	pub := publisher.PubKey
	if !pub.Verify(MustMarshalJSON(music), signature.Value) {
		WriteJSON(w, NewQueryResult(ErrInvalidSignature.Error(), false))
		return
	}
	api.logger.Info("SUCCESS verified release")
	WriteJSON(w, NewQueryResult("", true))
}

func AgentFromValues(values url.Values) (*core.Agent, error) {
	email := values.Get("email")
	name := values.Get("name")
	pub := new(ed25519.PublicKey)
	pubstr := values.Get("public_key")
	if err := pub.FromString(pubstr); err != nil {
		return nil, err
	}
	_type := values.Get("type")
	switch _type {
	case core.ARTIST:
		return core.NewArtist(email, name, pub), nil
	case core.LABEL:
		return core.NewLabel(email, name, pub), nil
	case core.ORGANIZATION:
		return core.NewOrganization(email, name, pub), nil
	case core.PUBLISHER:
		return core.NewPublisher(email, name, pub), nil
		// TODO: add more partner types?
	}
	return nil, ErrInvalidType
}
