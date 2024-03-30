package main

import (
	"net/http"
)

func main() {
	webdav := NewWebDAVClient()
	http.Handle("/", webdav)
	http.ListenAndServe(":8080", nil)
}