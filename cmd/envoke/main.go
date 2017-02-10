package main

import (
	"github.com/zbo14/envoke/api"
	. "github.com/zbo14/envoke/common"
	"net/http"
)

func main() {

	CreatePages(
		"login",
		"music",
		"sign",
		"verify",
	)

	RegisterTemplates(
		"login.html",
		"music.html",
		"sign.html",
		"verify.html",
	)

	// Create request multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/login", TemplateHandler("login.html"))
	mux.HandleFunc("/music", TemplateHandler("music.html"))
	mux.HandleFunc("/sign", TemplateHandler("sign.html"))
	mux.HandleFunc("/verify", TemplateHandler("verify.html"))
	fs := http.Dir("static/")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(fs)))

	// Create api
	api := api.NewApi()

	// Add routes to multiplexer
	api.AddRoutes(mux)

	// Start HTTP server with multiplexer
	http.ListenAndServe(":8888", mux)
}
