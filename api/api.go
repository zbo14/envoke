package api

import (
	"github.com/dhowden/tag"
	"github.com/minio/minio-go"
	"github.com/zbo14/envoke/bigchain"
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/crypto"
	"github.com/zbo14/envoke/crypto/ed25519"
	"github.com/zbo14/envoke/crypto/rsa"
	"github.com/zbo14/envoke/spec"
	// "github.com/zbo14/envoke/spec/coala"
	mo "github.com/zbo14/envoke/spec/music_ontology"
	"net/http"
	"net/url"
	"time"
)

const (
	EXPIRY_TIME = 1000 * time.Second
	ID_SIZE     = 47
	IMPL        = spec.JSON

	MINIO_ENDPOINT   = "http://127.0.0.1:9000"
	MINIO_ACCESS_KEY = "N3R2IT5XGCOMVIAUI25K"
	MINIO_SECRET_KEY = "I9zaxZWzbdvpbQO0hT6+bBaEJyHJF78RA2wAFNvJ"
	MINIO_REGION     = "us-east-1"
)

type Api struct {
	artist  spec.Data
	cli     *minio.Client
	logger  Logger
	priv    crypto.PrivateKey
	partner spec.Data
	pub     crypto.PublicKey
}

func NewApi() *Api {
	logger := NewLogger("api")
	return &Api{
		logger: logger,
	}
}

func (api *Api) AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/artist_login", api.ArtistLogin)
	// mux.HandleFunc("/artist_register", api.ArtistRegister)
	mux.HandleFunc("/listen_track", api.ListenTrack)
	mux.HandleFunc("/partner_login", api.PartnerLogin)
	mux.HandleFunc("/partner_register", api.PartnerRegister)
	mux.HandleFunc("/upload_album", api.UploadAlbum)
}

func GenerateId(key string) string {
	hash := Shake256([]byte(key), ID_SIZE)
	return Base64UrlEncode(hash)
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
	// email := values.Get("email")
	name := values.Get("name")
	openId := values.Get("open_id")
	partnerId := values.Get("partner_id")
	return mo.NewArtist(IMPL, name, openId, partnerId)
	// return coala.NewArtist(spec.JSON, email, name, openId, partnerId)
}

func PartnerFromValues(values url.Values) spec.Data {
	_type := values.Get("type")
	// email := values.Get("email")
	login := values.Get("login")
	name := values.Get("name")
	openId := values.Get("open_id")
	switch _type {
	case mo.LABEL:
		lc := values.Get("label_code")
		return mo.NewLabel(IMPL, lc, login, name, openId)
		// return coala.NewLabel(IMPL, id, email, login, name)
	case mo.PUBLISHER:
		return mo.NewPublisher(IMPL, login, name, openId)
		// return coala.NewPublisher(IMPL, id, email, login, name)
		// TODO: add more partner types?
	}
	panic("Unexpected partner type: " + _type)
}

// Partner

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
	// Generate keypair from password
	// password := values.Get("password")
	// priv, pub := ed25519.GenerateKeypair(password)
	privRSA, pubRSA := rsa.GenerateKeypair()
	pub := mo.NewPublicKey(IMPL, pubRSA.EncodePEM())
	pubTx := bigchain.GenerateTx(pub, nil, pubRSA)
	pubTx.Fulfill(privRSA)
	pubId := pubTx.Id
	// New partner
	partner := PartnerFromValues(values)
	mo.AddPublicKey(IMPL, partner, pubId)
	partnerTx := bigchain.GenerateTx(partner, nil, pubRSA)
	partnerTx.Fulfill(privRSA)
	partnerId := partnerTx.Id
	mo.AddOwner(IMPL, partnerId, pub)
	pubTx.SetData(pub) //update tx data
	/*
		// send requests to IPDB
		id, err := bigchain.PostTx(pubTx)
		Check(err)
		if id != pubId {
			Panicf("Expected id=%s; got id=%s\n", pubId, id)
		}
		id, err = bigchain.PostTx(partnerTx)
		Check(err)
		if id != partnerId {
			Panicf("Expected id=%s; got id=%s\n", partnerId, id)
		}
	*/
	api.logger.Info("Partner: " + string(MustMarshalIndentJSON(partner)))
	api.logger.Info("Public_Key: " + string(MustMarshalIndentJSON(pub)))
	partnerInfo := NewPartnerInfo(partnerId, privRSA, pubRSA)
	WriteJSON(w, partnerInfo)
}

// TODO: change to RSA keys
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
	if err = priv.FromString(priv58); err != nil {
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
	api.partner = partner
	api.priv = priv
	api.pub = priv.Public()
	w.Write([]byte("Login successful!"))
}

// Artist

// should we do login or just registration via partner org?
// having artist identity on envoke might ease attribution
// e.g. album/track contains uri to artist profile in db
// but artist must be verified by partner org before they
// create profile..

func (api *Api) ArtistLogin(w http.ResponseWriter, req *http.Request) {
	//Should be post request
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
	partnerId := values.Get("partner_id")
	// Query IPDB
	tx, err := bigchain.GetTx(partnerId)
	if err != nil {
		api.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	partnerLogin := tx.GetValue("login")
	Println(partnerLogin)
	// TODO: send POST request with artist credentials to login url
	// get artist info in response?
	// If login via partner is successful:
	artist := ArtistFromValues(values)
	password := values.Get("password")
	/*
		// Generate keypair from password
		password := values.Get("password")
		priv, pub := ed25519.GenerateKeypair(password)
		// send request to IPDB
		tx = bigchain.GenerateTx(partner, nil, pub)
		tx.Fulfill(priv)
		id, err := bigchain.PostTx(tx)
		if err != nil {
			api.logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	*/
	// artistId := tx.Id
	api.logger.Info("Artist: " + string(MustMarshalIndentJSON(artist)))
	api.artist = artist
	// Generate one-off keypair for now
	api.priv, api.pub = ed25519.GenerateKeypair(password)
	w.Write([]byte("Login successful!"))
}

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

func (api *Api) UploadAlbum(w http.ResponseWriter, req *http.Request) {
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
	form, err := MultipartForm(req)
	if err != nil {
		http.Error(w, "Failed to read request data", http.StatusBadRequest)
		return
	}
	artistName := api.artist["name"].(string)
	albumTitle := form.Value["album_title"][0]
	// albumLocation := form.Value["album_location"][0]
	albumId := GenerateId(artistName + albumTitle)
	// datePublished := DateString()
	/*
		if exists, _ := cli.BucketExists(albumTitle); exists {
			http.Error(w, "You already have album with title="+albumTitle, http.StatusBadRequest)
			return
		}
		album := mo.NewRecord(spec.JSON, "", albumTitle, 0, partnerId) //rename to record
		// album := coala.NewAlbum(spec.JSON, "", albumTitle, api.artist)
		// Generate album tx
		albumTx := bigchain.GenerateTx(album, nil, api.pub)
		albumTx.Fulfill(api.priv)
	*/
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
			track := mo.NewTrack(spec.JSON, "", trackTitle, api.artist, i, albumId)
			// track := coala.NewTrack(spec.JSON, "", trackTitle, nil, albumId, api.artist, datePublished, albumLocation, trackURL)
			// Generate and send track tx
			trackTx := bigchain.GenerateTx(track, metadata, api.pub)
			trackTx.Fulfill(api.priv)
			trackIds[i], err = bigchain.PostTx(trackTx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		*/
	}
	/*
		mo.AddTracks(spec.JSON, album, trackIds)
		albumTx.SetData(album)
		albumId, err := bigchain.PostTx(albumTx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	*/
	albumInfo := NewAlbumInfo(albumId, trackIds)
	WriteJSON(w, albumInfo)
}
