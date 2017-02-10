package main

import (
	"github.com/zbo14/envoke/api"
	. "github.com/zbo14/envoke/common"
	"net/http"
)

func main() {

	CreatePages(
		"artist",
		"login",
		"partner",
		"verification",
	)

	RegisterTemplates(
		"artist.html",
		"login.html",
		"partner.html",
		"verification.html",
	)

	// Create request multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/artist", TemplateHandler("artist.html"))
	mux.HandleFunc("/login", TemplateHandler("login.html"))
	mux.HandleFunc("/partner", TemplateHandler("partner.html"))
	mux.HandleFunc("/verification", TemplateHandler("verification.html"))
	fs := http.Dir("static/")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(fs)))

	// Create api
	api := api.NewApi()

	// Add routes to multiplexer
	api.AddRoutes(mux)

	// Start HTTP server with multiplexer
	http.ListenAndServe(":8888", mux)
}
