package main

import (
	"net/http"
)

var Cfg *Config

func main() {
	var err error
	Cfg, err = LoadConfig()
	if err != nil {
		Logoutput("Unable to load config", "error")
		return
	}

	Logoutput("Webdav server started", "info")
	Logoutput("Log level: "+Cfg.Loglevel, "info")
	webdav := NewWebDAVClient()
	Logoutput("Starting server on port "+Cfg.Port, "info")
	Logoutput("Base URL: "+Cfg.BaseURL, "info")
	http.Handle("/", webdav)
	http.ListenAndServe(":"+Cfg.Port, nil)
}