package api

import (
	"github.com/dhowden/tag"
	"github.com/minio/minio-go"
	"github.com/zballs/envoke/bigchain"
	"github.com/zballs/envoke/coala"
	"github.com/zballs/envoke/crypto/ed25519"
	"github.com/zballs/envoke/types"
	. "github.com/zballs/envoke/util"
	"net/http"
	"net/url"
	"time"
)

const (
	EXPIRY_TIME = 1000 * time.Second
	ID_SIZE     = 63

	MINIO_ENDPOINT   = "http://127.0.0.1:9000"
	MINIO_ACCESS_KEY = "N3R2IT5XGCOMVIAUI25K"
	MINIO_SECRET_KEY = "I9zaxZWzbdvpbQO0hT6+bBaEJyHJF78RA2wAFNvJ"
	MINIO_REGION     = "us-east"
)

var signature = ""

type Api struct {
	cli       *minio.Client
	logger    types.Logger
	partner   coala.Data
	partnerId string
	priv      *ed25519.PrivateKey
	pub       *ed25519.PublicKey
}

func NewApi() *Api {
	logger := types.NewLogger("api")
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

func PartnerFromValues(values url.Values) coala.Data {
	email := values.Get("email")
	id := values.Get("id")
	name := values.Get("name")
	_type := values.Get("type")
	return coala.NewOrganization(coala.JSON, id, email, name, _type)
	// return coala.NewPartner(coala.JSON,id,email,name,_type)
}

func (api *Api) GenerateId(key string) string {
	hex := api.priv.Sign([]byte(key)).ToHex()
	return hex[:ID_SIZE]
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
		id, err := bigchain.PostTransaction(t)
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
	// PrivKey
	priv := new(ed25519.PrivateKey)
	priv58 := values.Get("private_key")
	if err = priv.FromB58(priv58); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// Sign partner data
	sig := priv.Sign(json)
	// Query tx with id
	tx, err := bigchain.GetTransaction(partner["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Verify signature
	data := tx.GetData()
	json = MustMarshalJSON(data)
	pub := priv.Public()
	if !pub.Verify(data, sig) {
		http.Error(w, Error("Failed to verify signature").Error(), http.StatusUnauthorized)
		return
	}
	api.cli, err = NewClient(MINIO_ENDPOINT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	api.priv = priv
	api.pub = pub
	api.partner = partner
	api.partnerId = partnerId
	WriteJSON(w, NewLogin(api.partner.Type))
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
	if api.priv == nil {
		http.Error(w, "Private-key is not set", http.StatusUnauthorized)
		return
	}
	if api.pub == nil {
		http.Error(w, "Public-key is not set", http.StatusUnauthorized)
		return
	}
	if api.partner == nil {
		http.Error(w, "Partner-profile is not set", http.StatusUnauthorized)
		return
	}
	if api.partnerId == "" {
		http.Error(w, "Partner-id is not set", http.StatusUnauthorized)
		return
	}
	// Get request data
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	projectTitle := values.Get("project_title")
	trackTitle := values.Get("track_title")
	// projectId := api.GenerateId(projectTitle)
	// trackId := api.GenerateId(trackTitle)
	/*
		trackId := values.Get("track_id")
		t, err := bigchain.GetTransaction(trackId)
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
	presignedURL, err := api.cli.PresignedGetObject(projectTitle, trackTitle+".mp3", EXPIRY_TIME, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, NewStream("", projectTitle, trackTitle, presignedURL.String()))
}

func (api *Api) ReleaseAlbum(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
		return
	}
	if api.cli == nil {
		http.Error(w, "Minio-client is not set", http.StatusUnauthorized)
		return
	}
	// Make sure we're logged in
	if api.priv == nil {
		http.Error(w, "Private-key is not set", http.StatusUnauthorized)
		return
	}
	if api.pub == nil {
		http.Error(w, "Public-key is not set", http.StatusUnauthorized)
		return
	}
	if api.partner == nil {
		http.Error(w, "Partner-profile is not set", http.StatusUnauthorized)
		return
	}
	if api.partnerId == "" {
		http.Error(w, "Partner-id is not set", http.StatusUnauthorized)
		return
	}
	// Get request data
	form, err := MultipartForm(req)
	if err != nil {
		http.Error(w, "Failed to read request data", http.StatusBadRequest)
		return
	}
	// Project title
	projectTitle := form.Value["project_title"][0]
	projectId := api.GenerateId(projectTitle)
	/*
		// Create new project
		data := api.NewProject(projectTitle)
		// Generate and send transaction to IPDB
		t := bigchain.GenerateTransaction(data, nil, api.pub)
		t.Fulfill(api.priv, api.pub)
		projectId, err := bigchain.PostTransaction(t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	*/
	// Location
	// location := form.Value["location"][0]
	/*
		if exists, _ := cli.BucketExists(projectTitle); exists {
			http.Error(w, "You already have project with title="+projectTitle, http.StatusBadRequest)
			return
		}
	*/
	err = api.cli.MakeBucket(projectTitle, api.partner.Region)
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
		trackId := api.GenerateId(trackTitle)
		// trackURL := ""
		// Upload track to minio
		_, err = api.cli.PutObject(projectTitle, track.Filename, r, "audio/mp3")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		file.Close()
		trackIds[i] = trackId
		/*
			data = api.NewTrack(projectId, trackTitle, location, trackURL)
			// Generate and send transaction to IPDB
			t := bigchain.GenerateTransaction(data, metadata, api.pub)
			trackIds[i], err = bigchain.PostTransaction(t)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		*/
	}
	projectInfo := NewProjectInfo(projectId, trackIds)
	WriteJSON(w, projectInfo)
}
