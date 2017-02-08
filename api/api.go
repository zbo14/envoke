package api

import (
	// "bytes"
	"github.com/dhowden/tag"
	"github.com/minio/minio-go"
	"github.com/zbo14/envoke/bigchain"
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/ed25519"
	"github.com/zbo14/envoke/spec"
	mo "github.com/zbo14/envoke/spec/music_ontology"
	"net/http"
	"net/url"
	"time"
)

const (
	EXPIRY_TIME = 1000 * time.Second
	IMPL        = spec.JSON

	MINIO_ENDPOINT   = "http://127.0.0.1:9000"
	MINIO_ACCESS_KEY = "N3R2IT5XGCOMVIAUI25K"
	MINIO_SECRET_KEY = "I9zaxZWzbdvpbQO0hT6+bBaEJyHJF78RA2wAFNvJ"
	MINIO_REGION     = "us-east-1"

	ALBUM     = "album"
	SIGNATURE = "signature"
)

type Api struct {
	cli    *minio.Client
	user   spec.Data
	logger Logger
	userId string
	priv   *ed25519.PrivateKey
}

func NewApi() *Api {
	logger := NewLogger("api")
	return &Api{
		logger: logger,
	}
}

func (api *Api) AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/artist_login", api.ArtistLogin)
	mux.HandleFunc("/artist_register", api.ArtistRegister)
	mux.HandleFunc("/partner_login", api.PartnerLogin)
	mux.HandleFunc("/partner_register", api.PartnerRegister)
	mux.HandleFunc("/sign_album", api.SignAlbum)
	mux.HandleFunc("/upload_album", api.UploadAlbum)
}

// Minio client

func NewClient(endpoint, accessId, secretId string) (*minio.Client, error) {
	host, secure := HostSecure(endpoint)
	cli, err := minio.New(host, accessId, secretId, secure)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func HostSecure(rawurl string) (string, bool) {
	url := MustParseUrl(rawurl)
	return url.Host, url.Scheme == "https"
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

// should trackIds be trackURLs?
// should we persist an asset for each track on the album,
// or should we just persist an asset for the album?
func (api *Api) AlbumFromRequest(w http.ResponseWriter, req *http.Request) {
	form, err := MultipartForm(req)
	if err != nil {
		http.Error(w, ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	albumTitle := form.Value["album_title"][0]
	publisherId := form.Value["publisher_id"][0]
	Println(publisherId)
	/*
		if exists, _ := cli.BucketExists(albumTitle); exists {
			panic("You already have album with title=" + albumTitle)
		}
		album := mo.NewRecord(IMPL, api.userId, 0, publisherId, albumTitle)
		// Generate album tx
		albumTx := bigchain.GenerateTx(album, nil, api.pub)
		albumId := albumTx.Id
	*/
	albumId := ""
	err = api.cli.MakeBucket(albumTitle, MINIO_REGION)
	Check(err)
	// Tracks
	tracks := form.File["tracks"]
	trackIds := make([]string, len(tracks))
	// It would be great if we could batch write tracks
	for i, track := range tracks {
		file, err := track.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s, r := MustTeeSeeker(file)
		// Extract metadata
		meta, err := tag.ReadFrom(s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// metadata := meta.Raw()
		trackTitle := meta.Title()
		Println(trackTitle)
		trackURL := ""
		// Upload track to minio
		_, err = api.cli.PutObject(albumTitle, track.Filename, r, "audio/mp3")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		file.Close()
		/*
			track := mo.NewTrack(IMPL, api.userId, i, albumId, trackTitle)
			// Generate track tx
			trackTx := bigchain.GenerateTx(track, metadata, api.priv.Public())
			trackIds[i] = trackTx.Id
		*/
		trackIds[i] = trackURL
	}
	/*
		mo.AddTracks(IMPL, album, trackIds)
		albumTx.SetData(album) //update tx data
		_, err := bigchain.PostTx(albumTx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	*/
	albumInfo := NewTxInfo(trackIds, albumId, ALBUM)
	WriteJSON(w, albumInfo)
}

func (api *Api) SignatureFromRequest(w http.ResponseWriter, req *http.Request) {
	// Get request data
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	albumId := values.Get("album_id")
	// Query IPDB
	albumTx, err := bigchain.GetTx(albumId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	album := bigchain.GetTxData(albumTx)
	json := MustMarshalJSON(album)
	sigstr := api.priv.Sign(json).String()
	sig := mo.NewSignature(IMPL, api.userId, sigstr)
	// Send tx with signature to IPDB
	sigTx := bigchain.GenerateTx(sig, nil, bigchain.CREATE, api.priv.Public())
	sigId, err := bigchain.PostTx(sigTx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sigInfo := NewTxInfo(nil, sigId, SIGNATURE)
	WriteJSON(w, sigInfo)
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
	userInfo := NewUserInfo(id, priv.String(), pub.String())
	WriteJSON(w, userInfo)
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
	/*
		api.cli, err = NewClient(MINIO_ENDPOINT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	*/
	api.priv = priv
	api.user = user
	api.userId = userId
	w.Write([]byte("Login successful!"))
}

func (api *Api) HandleAction(w http.ResponseWriter, req *http.Request, handler http.HandlerFunc) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
		return
	}
	// Make sure we're logged in properly
	if api.cli == nil {
		http.Error(w, "Minio client is not set", http.StatusUnauthorized)
		return
	}
	if api.priv == nil {
		http.Error(w, "Privkey is not set", http.StatusUnauthorized)
		return
	}
	if api.user == nil {
		http.Error(w, "User profile is not set", http.StatusUnauthorized)
		return
	}
	if api.userId == "" {
		http.Error(w, "User id is not set", http.StatusUnauthorized)
		return
	}
	handler(w, req)
}

func (api *Api) ArtistRegister(w http.ResponseWriter, req *http.Request) {
	Register(w, req, ArtistFromValues)
}

func (api *Api) PartnerRegister(w http.ResponseWriter, req *http.Request) {
	Register(w, req, PartnerFromValues)
}

func (api *Api) ArtistLogin(w http.ResponseWriter, req *http.Request) {
	api.Login(w, req)
}

func (api *Api) PartnerLogin(w http.ResponseWriter, req *http.Request) {
	api.Login(w, req)
}

func (api *Api) SignAlbum(w http.ResponseWriter, req *http.Request) {
	api.HandleAction(w, req, api.AlbumFromRequest)
}

func (api *Api) UploadAlbum(w http.ResponseWriter, req *http.Request) {
	api.HandleAction(w, req, api.SignatureFromRequest)
}

/*
func (api *Api) ListenTrack(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
		return
	}
	// Make sure we're logged in properly
	if api.artist == nil {
		http.Error(w, "Could not identify artist", http.StatusUnauthorized)
		return
	}
	if api.cli == nil {
		http.Error(w, "Minio-client is not set", http.StatusUnauthorized)
		return
	}
	// Get request data
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	albumTitle := values.Get("album_title")
	trackTitle := values.Get("track_title")
	trackId := values.Get("track_id")
	t, err := bigchain.GetTx(trackId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	streamAddr := t.GetValue("url").(string)
	if streamAddr == "" {
		http.Error(w, "Could not find track url", http.StatusNotFound)
		return
	}
	// Get track url
	presignedURL, err := api.cli.PresignedGetObject(albumTitle, trackTitle+".mp3", EXPIRY_TIME, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, NewStream("", albumTitle, trackTitle, presignedURL.String()))
}
*/
