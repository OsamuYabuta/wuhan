package main

import (
	"api"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	mux := httprouter.New()
	mux.GET("/topic/:lang", api.Api_topic)
	mux.GET("/pickedupusers/:lang", api.Api_Pickedupusers)

	server := http.Server{
		Addr:    "localhost:8000",
		Handler: mux,
	}
	server.ListenAndServe()
}
