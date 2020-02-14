package main

import (
	"api"
	"net/http"
	"fmt"

	. "config"
	"github.com/julienschmidt/httprouter"
)

var GlobalConfig Config = Config{}

func main() {
	mux := httprouter.New()
	mux.GET("/topic/:lang", api.Api_topic)
	mux.GET("/pickedupusers/:lang", api.Api_Pickedupusers)

	GlobalConfig.Init()
	bindIp := GlobalConfig.BindIp()
	bindPort := GlobalConfig.BindPort()

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", bindIp, bindPort),
		Handler: mux,
	}
	server.ListenAndServe()
}
