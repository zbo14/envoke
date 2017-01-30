package main

import (
	"github.com/zballs/envoke/api"
	. "github.com/zballs/envoke/util"
	"net/http"
)

func main() {

	CreatePages(
		"artist",
		"listener",
		"login",
	)

	RegisterTemplates(
		"artist.html",
		"listener.html",
		"login.html",
	)

	// Create request multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/artist", TemplateHandler("artist.html"))
	mux.HandleFunc("/listener", TemplateHandler("listener.html"))
	mux.HandleFunc("/login", TemplateHandler("login.html"))
	fs := http.Dir("static/")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(fs)))

	// Create api
	api := api.NewApi()

	// Add routes to multiplexer
	api.AddRoutes(mux)

	// Start HTTP server with multiplexer
	http.ListenAndServe(":8888", mux)
}
