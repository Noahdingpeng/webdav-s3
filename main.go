package main

import (
	"net/http"
	"log"
)

var Cfg *Config

func main() {
	var err error
	Cfg, err = LoadConfig()
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}
	
	webdav := NewWebDAVClient()
	http.Handle("/", webdav)
	log.Println("Server started on port " + Cfg.Port)
	http.ListenAndServe(":"+Cfg.Port, nil)

}