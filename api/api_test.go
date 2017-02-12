package api

import (
	"net/url"
	"testing"

	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/spec/core"
)

func TestApi(t *testing.T) {
	api := NewApi()
	// Register
	values := make(url.Values)
	values.Set("email", "artist@email.com")
	values.Set("name", "artist_name")
	values.Set("password", "canyouguess")
	values.Set("type", core.ARTIST)
	registerMessage, err := api.Register(values)
	if err != nil {
		t.Error(err.Error())
	}
	PrintJSON(registerMessage)
	// Login
	agentId := registerMessage.AgentId
	privstr := registerMessage.PrivKey
	_type := values.Get("type")
	loginMessage, err := api.Login(agentId, privstr, _type)
	if err != nil {
		t.Error(err.Error())
	}
	PrintJSON(loginMessage)
	// TODO: Album, Track, Sign, Verify
}
