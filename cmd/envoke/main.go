package main

import (
	"github.com/zbo14/envoke/api"
	. "github.com/zbo14/envoke/common"
	"net/http"
)

func main() {

	CreatePages(
		"compose_publish",
		"login_register",
		"right_license",
		"record_release",
	)

	RegisterTemplates(
		"compose_publish.html",
		"login_register.html",
		"right_license.html",
		"record_release.html",
	)

	// Create request multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/compose_publish", TemplateHandler("compose_publish.html"))
	mux.HandleFunc("/login_register", TemplateHandler("login_register.html"))
	mux.HandleFunc("/right_license", TemplateHandler("right_license.html"))
	mux.HandleFunc("/record_release", TemplateHandler("record_release.html"))
	fs := http.Dir("static/")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(fs)))

	// Create api
	api := api.NewApi()

	// Add routes to multiplexer
	api.AddRoutes(mux)

	// Start HTTP server with multiplexer
	http.ListenAndServe(":8888", mux)
}
