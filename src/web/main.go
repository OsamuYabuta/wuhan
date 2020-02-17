package main

import (
	"api"
	"fmt"
	"net/http"

	. "config"
	"github.com/julienschmidt/httprouter"
)

var GlobalConfig Config = Config{}

func main() {
	fmt.Println("main start...")
	mux := httprouter.New()
	mux.GET("/topic/:lang", api.Api_topic)
	mux.GET("/pickedupusers/:lang", api.Api_Pickedupusers)

	GlobalConfig.Init()
	bindIp := GlobalConfig.BindIp()
	bindPort := GlobalConfig.BindPort()

	fmt.Printf("bindIp:%s\n", bindIp)
	fmt.Printf("bindPort:%s\n", bindPort)

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", bindIp, bindPort),
		Handler: mux,
	}
	err := server.ListenAndServe()

	if err != nil {
		panic(err.Error())
	}
}
