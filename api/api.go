package api

import (
	"github.com/dhowden/tag"
	"github.com/minio/minio-go"
	"github.com/zballs/envoke/bigchain"
	"github.com/zballs/envoke/crypto/ed25519"
	"github.com/zballs/envoke/spec"
	. "github.com/zballs/envoke/util"
	"net/http"
	"net/url"
	"time"
)

const (
	EXPIRY_TIME = 1000 * time.Second
	ID_SIZE     = 63
	LABEL       = "label"
	PUBLISHER   = "publisher"

	MINIO_ENDPOINT   = "http://127.0.0.1:9000"
	MINIO_ACCESS_KEY = "N3R2IT5XGCOMVIAUI25K"
	MINIO_SECRET_KEY = "I9zaxZWzbdvpbQO0hT6+bBaEJyHJF78RA2wAFNvJ"
	MINIO_REGION     = "us-east-1"
)

type Api struct {
	cli     *minio.Client
	logger  Logger
	partner spec.Data
	priv    *ed25519.PrivateKey
	pub     *ed25519.PublicKey
}

func NewApi() *Api {
	logger := NewLogger("api")
	return &Api{
		logger: logger,
	}
}

func (api *Api) AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/partner_login", api.PartnerLogin)
	mux.HandleFunc("/partner_register", api.PartnerRegister)
	mux.HandleFunc("/release_album", api.ReleaseAlbum)
	mux.HandleFunc("/stream_track", api.StreamTrack)
	mux.HandleFunc("/user_login", api.UserLogin)
}

func GenerateId(key string) string {
	hash := Shake256([]byte(key), ID_SIZE)
	return BytesToHex(hash)
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

// Partner

func PartnerFromValues(values url.Values) spec.Data {
	email := values.Get("email")
	id := values.Get("id")
	name := values.Get("name")
	_type := values.Get("type")
	Println(email, id, name, _type)
	switch _type {
	case LABEL:
		return spec.Data{} //TODO: determine which data model to use
	case PUBLISHER:
		return spec.Data{}
	default:
		panic("Unexpected partner type: " + _type)
	}
}

func (api *Api) PartnerRegister(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
		return
	}
	// Get request data
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, "Failed to read request data", http.StatusBadRequest)
		return
	}
	// New partner
	partner := PartnerFromValues(values)
	// Generate keypair from password
	password := values.Get("password")
	priv, pub := ed25519.GenerateKeypair(password)
	// send request to IPDB
	tx := bigchain.GenerateTx(partner, nil, pub)
	tx.Fulfill(priv)
	/*
		id, err := bigchain.PostTx(t)
		if err != nil {
			api.logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	*/
	id := tx.Id
	api.logger.Info("Partner: " + string(MustMarshalIndentJSON(partner)))
	partnerInfo := NewPartnerInfo(id, priv, pub)
	WriteJSON(w, partnerInfo)
}

func (api *Api) PartnerLogin(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
		return
	}
	// Get request data
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, "Failed to read request data", http.StatusBadRequest)
		return
	}
	// Partner
	partner := PartnerFromValues(values)
	json := MustMarshalJSON(partner)
	Println(json)
	// PrivKey
	priv := new(ed25519.PrivateKey)
	priv58 := values.Get("private_key")
	if err = priv.FromB58(priv58); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// Sign partner data
	sig := priv.Sign(json)
	Println(sig)
	/*
		// Query tx with id
		tx, err := bigchain.GetTx(partner["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		// Verify signature
		data := tx.GetData()
		json = MustMarshalJSON(data)
		pub := priv.Public()
		if !pub.Verify(data, sig) {
			http.Error(w, "Failed to verify signature", http.StatusUnauthorized)
			return
		}
	*/
	api.cli, err = NewClient(MINIO_ENDPOINT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	api.priv = priv
	api.pub = priv.Public()
	api.partner = partner
	w.Write([]byte("Login successful!"))
}

// User

func (api *Api) UserLogin(w http.ResponseWriter, req *http.Request) {
	// Should be post request
	if req.Method != http.MethodPost {
		http.Error(w, Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
		return
	}
	// Get request data
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, "Failed to read request data", http.StatusBadRequest)
		return
	}
	// Get partner id
	id := values.Get("id")
	// Query IPDB
	tx, err := bigchain.GetTx(id)
	if err != nil {
		api.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	partner := tx.GetData()
	loginURL := partner["login"]
	Println(loginURL)
	// TODO: send POST request with user info to login endpoint
	// If login is successful:
	api.cli, err = NewClient(MINIO_ENDPOINT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	api.partner = partner
	w.Write([]byte("Login successful!"))
}

func (api *Api) StreamTrack(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
		return
	}
	// Make sure we're logged in
	if api.cli == nil {
		http.Error(w, "Minio-client is not set", http.StatusUnauthorized)
		return
	}
	if api.partner == nil {
		http.Error(w, "Partner-profile is not set", http.StatusUnauthorized)
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
	/*
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
	*/
	// Get track url
	presignedURL, err := api.cli.PresignedGetObject(albumTitle, trackTitle+".mp3", EXPIRY_TIME, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, NewStream("", albumTitle, trackTitle, presignedURL.String()))
}

func (api *Api) ReleaseAlbum(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
		return
	}
	// Make sure we're logged in properly
	if api.cli == nil {
		http.Error(w, "Minio-client is not set", http.StatusUnauthorized)
		return
	}
	if api.partner == nil {
		http.Error(w, "Could not identify partner", http.StatusUnauthorized)
		return
	}
	// Get request data
	form, err := MultipartForm(req)
	if err != nil {
		http.Error(w, "Failed to read request data", http.StatusBadRequest)
		return
	}
	artistName := form.Value["artist_name"][0]
	albumTitle := form.Value["album_title"][0]
	albumId := GenerateId(artistName + albumTitle)
	/*
		if exists, _ := cli.BucketExists(albumTitle); exists {
			http.Error(w, "You already have album with title="+albumTitle, http.StatusBadRequest)
			return
		}
		// Generate one-off keypair for now
		priv, pub := ed25519.GenerateKeypair(password)
		// TODO: create new album
		// Generate and send transaction to IPDB
		tx := bigchain.GenerateTx(data, nil, pub)
		t.Fulfill(priv)
		albumId, err := bigchain.PostTx(tx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	*/
	err = api.cli.MakeBucket(albumTitle, MINIO_REGION)
	Check(err)
	// Tracks
	tracks := form.File["tracks"]
	trackIds := make([]string, len(tracks))
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
		// Track info
		trackTitle := meta.Title()
		trackId := GenerateId(artistName + albumTitle + trackTitle)
		// trackURL := ""
		// Upload track to minio
		_, err = api.cli.PutObject(albumTitle, track.Filename, r, "audio/mp3")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		file.Close()
		trackIds[i] = trackId
		/*
			// TODO: create new track
			// Generate and send transaction to IPDB
			tx := bigchain.GenerateTx(data, metadata, pub)
			trackIds[i], err = bigchain.PostTx(tx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		*/
	}
	albumInfo := NewAlbumInfo(albumId, trackIds)
	WriteJSON(w, albumInfo)
}
