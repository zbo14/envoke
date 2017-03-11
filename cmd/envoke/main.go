package main

import (
	"net/http"

	"github.com/zbo14/envoke/api"
	. "github.com/zbo14/envoke/common"
)

func main() {

	CreatePages(
		"compose_publish",
		"login_register",
		"record_release",
		"right_license",
		"schema",
	)

	RegisterTemplates(
		"compose_publish.html",
		"login_register.html",
		"record_release.html",
		"right_license.html",
		"schema.html",
	)

	// Create request multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/compose_publish", TemplateHandler("compose_publish.html"))
	mux.HandleFunc("/login_register", TemplateHandler("login_register.html"))
	mux.HandleFunc("/right_license", TemplateHandler("right_license.html"))
	mux.HandleFunc("/record_release", TemplateHandler("record_release.html"))
	mux.HandleFunc("/schema", TemplateHandler("schema.html"))
	fs := http.Dir("static/")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(fs)))

	// Create api
	api := api.NewApi()

	// Add routes to multiplexer
	api.AddRoutes(mux)

	// Start HTTP server with multiplexer
	Println(http.ListenAndServe(":8888", mux))
}
