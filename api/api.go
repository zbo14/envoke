package api

import (
	"fmt"
	// "github.com/gorilla/websocket"
	"github.com/tendermint/go-wire" //TODO: move to util
	tndr "github.com/tendermint/tendermint/types"
	"github.com/zballs/go_resonate/state"
	"github.com/zballs/go_resonate/types"
	. "github.com/zballs/go_resonate/util"
	"net/http"
)

const CHAIN = "res()nate"

type Api struct {
	blocks       chan *tndr.Block
	latestHeight int
	logger       types.Logger
	proxy        *types.Proxy
	user         *types.PrivateAccount
}

func NewApi(remote string) *Api {
	return &Api{
		blocks: make(chan *tndr.Block),
		logger: types.NewLogger("api"),
		proxy:  types.NewProxy(remote, "/websocket"),
	}
}

func (api *Api) AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/create_account", api.CreateAccount)
	mux.HandleFunc("/remove_account", api.RemoveAccount)
	mux.HandleFunc("/login", api.Login)
}

func Respond(w http.ResponseWriter, response interface{}) {
	json := JSON(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (api *Api) CreateAccount(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
		return
	}
	// Get request data
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, "Failed to read request data", http.StatusBadRequest)
	}
	password := values.Get("password")
	username := values.Get("username")
	// Create action
	action := types.NewAction([]byte(username), types.CREATE_ACCOUNT)
	// Generate new keypair
	priv, pub := GenerateKeypair(password)
	// Prepare and sign action
	action.Prepare(pub, 1) // pass sequence=1
	action.Sign(priv, CHAIN)
	// Broadcast tx
	result, err := api.proxy.BroadcastTx("sync", action.Tx())
	if err != nil {
		Respond(w, MessageCreateAccount(nil, err))
		return
	}
	if err = ResultToError(result); err != nil {
		Respond(w, MessageCreateAccount(nil, err))
		return
	}
	keypair := NewKeypairB58(pub, priv)
	Respond(w, MessageCreateAccount(keypair, nil))
}

func (api *Api) RemoveAccount(w http.ResponseWriter, req *http.Request) {
	// Make sure we're logged in
	var user *types.PrivateAccount
	if user = api.user; user == nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	// Should be a POST request
	if req.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
		return
	}
	// Create action
	action := types.NewAction(nil, types.REMOVE_ACCOUNT)
	// Prepare and sign action
	action.Prepare(user.PubKey, user.Sequence)
	action.Sign(user.PrivKey, CHAIN)
	// Broadcast tx
	result, err := api.proxy.BroadcastTx("sync", action.Tx())
	if err == nil {
		err = ResultToError(result)
	}
	Respond(w, MessageRemoveAccount(err))
}

func (api *Api) Login(w http.ResponseWriter, req *http.Request) {
	// Should be POST request
	if req.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("Expected POST request; got %s request", req.Method), http.StatusBadRequest)
		return
	}
	// Get request data
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, "Failed to read request data", http.StatusBadRequest)
		return
	}
	// PubKey
	pub := new(PublicKey)
	keystr := values.Get("pub_key")
	if err = pub.FromB58(keystr); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// PrivKey
	priv := new(PrivateKey)
	keystr = values.Get("priv_key")
	if err = priv.FromB58(keystr); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// Verify keypair
	signBytes := []byte(CHAIN)
	sig := priv.Sign(signBytes)
	verified := pub.Verify(signBytes, sig)
	if !verified {
		http.Error(w, "Invalid keypair", http.StatusUnauthorized)
		return
	}
	// Query account
	acckey := state.AccountKey(pub.Address())
	query := KeyQuery(acckey)
	result, err := api.proxy.TMSPQuery(query)
	if err == nil {
		err = ResultToError(result)
		if err == nil {
			err = wire.ReadBinaryBytes(result.Result.Data, &api.user.Account)
		}
	}
	Respond(w, MessageLogin(err))
	// Start ws
}
