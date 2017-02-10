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
	"github.com/zbo14/envoke/spec"
	mo "github.com/zbo14/envoke/spec/music_ontology"
)

const (
	IMPL = spec.JSON

	ALBUM     = "album"
	SIGNATURE = "signature"
	TRACK     = "track"
)

type Api struct {
	logger Logger
	priv   crypto.PrivateKey
	pub    crypto.PublicKey
	user   spec.Data
	userId string
}

func NewApi() *Api {
	return &Api{
		logger: NewLogger("api"),
	}
}

func (api *Api) AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/artist_login", api.ArtistLogin)
	mux.HandleFunc("/artist_register", api.ArtistRegister)
	mux.HandleFunc("/partner_login", api.PartnerLogin)
	mux.HandleFunc("/partner_register", api.PartnerRegister)
	mux.HandleFunc("/album", api.Album)
	mux.HandleFunc("/track", api.Track)
	mux.HandleFunc("/sign", api.Sign)
	mux.HandleFunc("/verify", api.Verify)
}

func (api *Api) ArtistLogin(w http.ResponseWriter, req *http.Request) {
	api.Login(w, req)
}

func (api *Api) ArtistRegister(w http.ResponseWriter, req *http.Request) {
	Register(w, req, ArtistFromValues)
}

func (api *Api) PartnerLogin(w http.ResponseWriter, req *http.Request) {
	api.Login(w, req)
}

func (api *Api) PartnerRegister(w http.ResponseWriter, req *http.Request) {
	Register(w, req, PartnerFromValues)
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
	if api.pub == nil {
		http.Error(w, "Pubkey is not set", http.StatusUnauthorized)
		return
	}
	if api.user == nil {
		http.Error(w, "User profile is not set", http.StatusUnauthorized)
		return
	}
	if api.userId == "" {
		http.Error(w, "User ID is not set", http.StatusUnauthorized)
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
	pub := priv.Public()
	// Query tx with user id
	userId := values.Get("user_id")
	tx, err := bigchain.GetTx(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	Println(tx)
	user := bigchain.GetTxData(tx)
	pubstr := mo.GetPublicKey(user)
	if pub.String() != pubstr {
		http.Error(w, ErrInvalidKey.Error(), http.StatusUnauthorized)
		return
	}
	api.priv = priv
	api.pub = pub
	api.user = user
	api.userId = userId
	w.Write([]byte("Login successful!"))
}

func Register(w http.ResponseWriter, req *http.Request, userFromValues func(url.Values) spec.Data) {
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
	// New user
	user := userFromValues(values)
	tx := bigchain.GenerateTx(user, nil, bigchain.CREATE, pub)
	bigchain.FulfillTx(tx, priv)
	// send request to IPDB
	id, err := bigchain.PostTx(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, NewUserInfo(id, priv.String(), pub.String()))
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
	music := bigchain.GetTxData(tx)
	publisherId := mo.GetPublisher(music)
	if api.userId != publisherId {
		http.Error(w, ErrInvalidId.Error(), http.StatusUnauthorized)
		return
	}
	json := MustMarshalJSON(music)
	sigstr := api.priv.Sign(json).String()
	sig := mo.NewSignature(IMPL, api.userId, sigstr)
	// Send tx with signature to IPDB
	tx = bigchain.GenerateTx(sig, nil, bigchain.CREATE, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, NewActionInfo(id, SIGNATURE))
}

func (api *Api) SendTrack(file multipart.File, number int, publisherId, recordId interface{}) (string, error) {
	// Extract metadata
	meta, err := tag.ReadFrom(file)
	if err != nil {
		return "", err
	}
	metadata := meta.Raw()
	trackTitle := meta.Title()
	// Create new track
	track := mo.NewTrack(IMPL, api.userId, number, publisherId, recordId, trackTitle)
	// Generate and send tx with track
	tx := bigchain.GenerateTx(track, metadata, bigchain.CREATE, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	trackId, err := bigchain.PostTx(tx)
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
	trackId, err := api.SendTrack(file, 0, publisherId, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, NewActionInfo(trackId, TRACK))
}

func (api *Api) AlbumFromRequest(w http.ResponseWriter, req *http.Request) {
	form, err := MultipartForm(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	albumTitle := form.Value["album_title"][0]
	publisherId := form.Value["publisher_id"][0]
	album := mo.NewRecord(IMPL, api.userId, 0, publisherId, albumTitle)
	// Generate and send tx with album
	tx := bigchain.GenerateTx(album, nil, bigchain.CREATE, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	albumId, err := bigchain.PostTx(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	albumInfo := NewActionInfo(albumId, ALBUM)
	WriteJSON(w, albumInfo)
	tracks := form.File["tracks"]
	// It would be great if we could batch write tracks
	for i, track := range tracks {
		file, err := track.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		trackId, err := api.SendTrack(file, i+1, nil, albumId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		WriteJSON(w, NewActionInfo(trackId, TRACK))
	}
}

func (api *Api) Verify(w http.ResponseWriter, req *http.Request) {
	// Should be GET request
	if req.Method != http.MethodGet {
		http.Error(w, Sprintf("Expected GET request; got %s request", req.Method), http.StatusBadRequest)
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
	music := bigchain.GetTxData(tx)
	publisherId := mo.GetPublisher(music)
	if publisherId != "" {
		api.VerifyMusic(music, signatureId, w)
		return
	}
	albumId := mo.GetRecord(music)
	if albumId == "" {
		WriteJSON(w, NewQueryResult(ErrInvalidId.Error(), false))
		return
	}
	tx, err = bigchain.GetTx(albumId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	album := bigchain.GetTxData(tx)
	api.VerifyMusic(album, signatureId, w)
}

func (api *Api) VerifyMusic(music spec.Data, signatureId string, w http.ResponseWriter) {
	publisherId := mo.GetPublisher(music)
	tx, err := bigchain.GetTx(signatureId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	signature := bigchain.GetTxData(tx)
	signerId := mo.GetSigner(signature)
	if publisherId != signerId {
		WriteJSON(w, NewQueryResult(ErrInvalidId.Error(), false))
		return
	}
	tx, err = bigchain.GetTx(publisherId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	publisher := bigchain.GetTxData(tx)
	keystr := mo.GetPublicKey(publisher)
	key := new(ed25519.PublicKey)
	err = key.FromString(keystr)
	Check(err)
	sigstr := mo.GetSig(signature)
	sig := new(ed25519.Signature)
	err = sig.FromString(sigstr)
	Check(err)
	if !key.Verify(MustMarshalJSON(music), sig) {
		WriteJSON(w, NewQueryResult(ErrInvalidSignature.Error(), false))
		return
	}
	WriteJSON(w, NewQueryResult("", true))
}

func ArtistFromValues(values url.Values) spec.Data {
	name := values.Get("name")
	openId := values.Get("open_id")
	pub := values.Get("public_key")
	return mo.NewArtist(IMPL, name, openId, pub)
}

func PartnerFromValues(values url.Values) spec.Data {
	_type := values.Get("type")
	name := values.Get("name")
	openId := values.Get("open_id")
	pub := values.Get("public_key")
	switch _type {
	case mo.LABEL:
		lc := values.Get("label_code")
		return mo.NewLabel(IMPL, lc, name, openId, pub)
	case mo.PUBLISHER:
		return mo.NewPublisher(IMPL, name, openId, pub)
		// TODO: add more partner types?
	}
	panic(ErrInvalidType.Error())
}
