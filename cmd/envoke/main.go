package main

import (
	"github.com/zbo14/envoke/api"
	. "github.com/zbo14/envoke/util"
	"net/http"
)

func main() {

	CreatePages(
		"listen",
		"login",
		"upload",
	)

	RegisterTemplates(
		"listen.html",
		"login.html",
		"upload.html",
	)

	// Create request multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/listen", TemplateHandler("listen.html"))
	mux.HandleFunc("/login", TemplateHandler("login.html"))
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
