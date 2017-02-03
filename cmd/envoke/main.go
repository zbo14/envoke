package main

import (
	"github.com/zballs/envoke/api"
	. "github.com/zballs/envoke/util"
	"net/http"
)

func main() {

	CreatePages(
		"envoke_vocab",
		"listen",
		"login",
		"upload",
	)

	RegisterTemplates(
		"envoke_vocab.json",
		"listen.html",
		"login.html",
		"upload.html",
	)

	// Create request multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/listen", TemplateHandler("listen.html"))
	mux.HandleFunc("/login", TemplateHandler("login.html"))
	mux.HandleFunc("/spec", TemplateHandler("envoke_vocab.json"))
	mux.HandleFunc("/upload", TemplateHandler("upload.html"))
	fs := http.Dir("static/")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(fs)))

	// Create api
	api := api.NewApi()

	// Add routes to multiplexer
	api.AddRoutes(mux)

	// Start HTTP server with multiplexer
	http.ListenAndServe(":8888", mux)
}
