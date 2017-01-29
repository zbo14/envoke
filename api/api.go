package api

import (
	"github.com/dhowden/tag"
	"github.com/minio/minio-go"
	"github.com/zballs/go_resonate/bigchain"
	"github.com/zballs/go_resonate/coala"
	"github.com/zballs/go_resonate/crypto/ed25519"
	"github.com/zballs/go_resonate/types"
	. "github.com/zballs/go_resonate/util"
	"net/http"
	"net/url"
)

const (
	MINIO_ENDPOINT = "192.168.99.100:9000"
	MINIO_ID       = "AH6LQWO6JCXWEKUX17CM"
	MINIO_SECRET   = "JD++tfxKzZdbq9jPpU6j4pxPLs++BZ3ZCm2e6jzk"
	USE_SSL        = false
)

var signature = ""

type Api struct {
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
	mux.HandleFunc("/play_track", api.PlayTrack)
	mux.HandleFunc("/user_login", api.UserLogin)
	mux.HandleFunc("/user_register", api.UserRegister)
}

func UserFromValues(values url.Values) *types.User {
	email := values.Get("email")
	region := values.Get("region")
	password := values.Get("password")
	_type := values.Get("type")
	username := values.Get("username")
	return types.NewUser(email, region, password, _type, username)
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
	data["user_signature"] = priv.Sign(json).String()
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
	// PrivKey
	priv := new(ed25519.PrivateKey)
	keystr := values.Get("private_key")
	if err = priv.FromString(keystr); err != nil {
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
		if err := sig.FromString(data["user_signature"].(string)); err != nil {
			http.Error(w, Error("Failed to verify user signature").Error(), http.StatusUnauthorized)
			return
		}
	*/
	// User
	user := UserFromValues(values)
	json := MustMarshalJSON(user)
	pub := priv.Public()
	sig := new(ed25519.Signature)
	err = sig.FromString(signature)
	Check(err)
	Println(pub, string(json), sig)
	if !pub.Verify(json, sig) {
		http.Error(w, Error("Failed to verify user signature").Error(), http.StatusUnauthorized)
		return
	}
	api.priv = priv
	api.pub = pub
	api.user = user
	api.userId = userId
	WriteJSON(w, "Logged in!")
}

func (api *Api) PlayTrack(w http.ResponseWriter, req *http.Request) {
	// Should be GET request
	if req.Method != http.MethodGet {
		http.Error(w, Sprintf("Expected GET request; got %s request", req.Method), http.StatusBadRequest)
		return
	}
	// Make sure we're logged in
	if api.priv == nil {
		http.Error(w, "Privkey is not set", http.StatusUnauthorized)
		return
	}
	if api.pub == nil {
		http.Error(w, "Pubkey is not set", http.StatusUnauthorized)
		return
	}
	if api.user == nil {
		http.Error(w, "User is not set", http.StatusUnauthorized)
		return
	}
	if api.userId == "" {
		http.Error(w, "User Id is not set", http.StatusUnauthorized)
		return
	}
	// Get request data
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	/*
		trackId := values.Get("track_id")
		t, err := bigchain.GetTransaction(trackId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		playAddr := t.GetValue("url").(string)
		if playAddr == "" {
			http.Error(w, "Could not find track url", http.StatusNotFound)
			return
		}
	*/
	projectTitle := values.Get("project_title")
	trackTitle := values.Get("track_title")
	cli, err := minio.New(MINIO_ENDPOINT, MINIO_ID, MINIO_SECRET, USE_SSL)
	Check(err)
	object, err := cli.GetObject(projectTitle, trackTitle)
	Check(err)
	Copy(w, object)
}

func (api *Api) CreateProject(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
		return
	}
	// Make sure we're logged in
	if api.priv == nil {
		http.Error(w, "Privkey is not set", http.StatusUnauthorized)
		return
	}
	if api.pub == nil {
		http.Error(w, "Pubkey is not set", http.StatusUnauthorized)
		return
	}
	if api.user == nil {
		http.Error(w, "User is not set", http.StatusUnauthorized)
		return
	}
	if api.userId == "" {
		http.Error(w, "User Id is not set", http.StatusUnauthorized)
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
	projectId := api.priv.Sign([]byte(projectTitle)).String() //for now
	// Location
	// location := form.Value["location"][0]
	// Initialize minio client object
	cli, err := minio.New(MINIO_ENDPOINT, MINIO_ID, MINIO_SECRET, USE_SSL)
	Check(err)
	if exists, _ := cli.BucketExists(projectTitle); exists {
		http.Error(w, "You already have project with title="+projectTitle, http.StatusBadRequest)
		return
	}
	err = cli.MakeBucket(projectTitle, api.user.Region)
	Check(err)
	// Tracks
	tracks := form.File["tracks"]
	trackIds := make([]string, len(tracks))
	for i, track := range tracks {
		file, err := track.Open()
		Check(err)
		// Get metadata
		meta, err := tag.ReadFrom(file)
		Check(err)
		// metadata := meta.Raw()
		Println(meta)
		// Track info
		trackTitle := meta.Title()
		// trackFormat := "audio/" + ToLower(string(meta.Format()))
		// trackURL := ""
		// Upload track to minio
		_, err = cli.PutObject(projectTitle, trackTitle, file, "audio/mp3")
		Check(err)
		file.Close()
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
		trackIds[i] = api.priv.Sign([]byte(trackTitle)).String() //for now
	}
	projectInfo := NewProjectInfo(projectId, trackIds)
	WriteJSON(w, projectInfo)
}
