package main

import (
	"github.com/zbo14/envoke/api"
	. "github.com/zbo14/envoke/common"
	"net/http"
)

func main() {

	CreatePages(
		"composition_recording",
		"license",
		"login_register",
	)

	RegisterTemplates(
		"composition_recording.html",
		"license.html",
		"login_register.html",
	)

	// Create request multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/composition_recording", TemplateHandler("composition_recording.html"))
	mux.HandleFunc("/license", TemplateHandler("license.html"))
	mux.HandleFunc("/login_register", TemplateHandler("login_register.html"))
	fs := http.Dir("static/")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(fs)))

	// Create api
	api := api.NewApi()

	// Add routes to multiplexer
	api.AddRoutes(mux)

	// Start HTTP server with multiplexer
	http.ListenAndServe(":8888", mux)
}
