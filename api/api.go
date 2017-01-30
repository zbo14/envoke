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
)

var signature = ""

type Api struct {
	cli    *minio.Client
	logger types.Logger
	priv   *ed25519.PrivateKey
	pub    *ed25519.PublicKey
	user   *types.User
	userId string
}

func NewApi() *Api {
	logger := types.NewLogger("api")
	return &Api{
		logger: logger,
	}
}

func (api *Api) AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/create_project", api.CreateProject)
	mux.HandleFunc("/stream_track", api.StreamTrack)
	mux.HandleFunc("/user_login", api.UserLogin)
	mux.HandleFunc("/user_register", api.UserRegister)
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

func UserFromValues(values url.Values) *types.User {
	email := values.Get("email")
	region := values.Get("region")
	password := values.Get("password")
	_type := values.Get("type")
	username := values.Get("username")
	return types.NewUser(email, region, password, _type, username)
}

func (api *Api) GenerateId(key string) string {
	hex := api.priv.Sign([]byte(key)).ToHex()
	return hex[:ID_SIZE]
}

func (api *Api) NewProject(projectTitle string) map[string]interface{} {
	artistName := api.user.Username
	data := coala.NewWork(projectTitle, artistName)
	return data
}

func (api *Api) NewTrack(projectId, trackTitle, location, trackURL string) map[string]interface{} {
	example := "" //what should example be?
	isManifestation := true
	date := DateString()
	data := coala.NewDigitalManifestation(trackTitle, example, isManifestation, projectId, date, location, trackURL)
	return data
}

func (api *Api) UserRegister(w http.ResponseWriter, req *http.Request) {
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
	// New user
	user := UserFromValues(values)
	// Generate keypair from password
	priv, pub := ed25519.GenerateKeypair(user.Password)
	data := make(map[string]interface{})
	// Sign user bytes
	json := MustMarshalJSON(user)
	data["user_signature"] = priv.Sign(json).ToB58()
	// send request to IPDB
	t := bigchain.GenerateTransaction(data, nil, pub)
	t.Fulfill(priv, pub)
	/*
		id, err := bigchain.PostTransaction(t)
		if err != nil {
			api.logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	*/
	id := t.Id
	signature = data["user_signature"].(string)
	api.logger.Info("User: " + string(MustMarshalIndentJSON(user)))
	api.logger.Info("Signature: " + signature)
	userInfo := NewUserInfo(id, priv, pub)
	WriteJSON(w, userInfo)
}

func (api *Api) UserLogin(w http.ResponseWriter, req *http.Request) {
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
	// PrivId
	priv := new(ed25519.PrivateKey)
	priv58 := values.Get("private_key")
	if err = priv.FromB58(priv58); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// User Id
	userId := values.Get("user_id")
	/*
		// Check that transaction with id exists
		t, err := bigchain.GetTransaction(userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		// Get user signature
		data := t.GetData()
		sig := new(ed25519.Signature)
		if err := sig.FromB58(data["user_signature"].(string)); err != nil {
			http.Error(w, Error("Failed to verify user signature").Error(), http.StatusUnauthorized)
			return
		}
	*/
	// User
	user := UserFromValues(values)
	json := MustMarshalJSON(user)
	pub := priv.Public()
	sig := new(ed25519.Signature)
	err = sig.FromB58(signature)
	Check(err)
	Println(pub, string(json), sig)
	if !pub.Verify(json, sig) {
		http.Error(w, Error("Failed to verify user signature").Error(), http.StatusUnauthorized)
		return
	}
	api.cli, err = NewClient(MINIO_ENDPOINT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	api.priv = priv
	api.pub = pub
	api.user = user
	api.userId = userId
	WriteJSON(w, NewLogin(api.user.Type))
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
	if api.user == nil {
		http.Error(w, "User-profile is not set", http.StatusUnauthorized)
		return
	}
	if api.userId == "" {
		http.Error(w, "User-id is not set", http.StatusUnauthorized)
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

func (api *Api) CreateProject(w http.ResponseWriter, req *http.Request) {
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
	if api.user == nil {
		http.Error(w, "User-profile is not set", http.StatusUnauthorized)
		return
	}
	if api.userId == "" {
		http.Error(w, "User-id is not set", http.StatusUnauthorized)
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
	err = api.cli.MakeBucket(projectTitle, api.user.Region)
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
