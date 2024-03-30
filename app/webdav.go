package main

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type WebDAVClient struct {
	Backend *S3Client
}

func NewWebDAVClient() *WebDAVClient {
	return &WebDAVClient{
		Backend: NewS3Client(),
	}
}

func (h *WebDAVClient) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	switch r.Method {
		case "GET":
			h.Get_Profind(w, r)
		case "PROPFIND":
			h.Get_Profind(w, r)
		case "PUT":
			h.Put(w, r)
		case "DELETE":
			h.Delete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *WebDAVClient) Get_Profind(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "" && strings.HasSuffix(r.URL.Path, "/") {
		h.Propfind(w, r)
	} else {
		h.Get(w, r)
	}
}

func (h *WebDAVClient) Get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path
	result, err := h.Backend.GetObject(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer result.Body.Close()

	for k, v := range result.Metadata {
		w.Header().Set(k, *v)
	}
	w.Header().Set("Content-Type", "application/octet-stream")

	_, err = io.Copy(w, result.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebDAVClient) Propfind(w http.ResponseWriter, r *http.Request) {
	keyPrefix := r.URL.Path[1:]
	if keyPrefix != "" && !strings.HasSuffix(keyPrefix, "/") {
        keyPrefix += "/"
    }
	result, err := h.Backend.ListObjects(keyPrefix)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	html := `
	<html>
	<head>
		<title>Index of ` + r.URL.Path + `</title>
		<style>
			th, td { text-align: left; padding: 0.5em; }
			th { border-bottom: 1px solid #eee; }
			body {
				background-color: black;
				color: white;
				font-family: sans-serif;
			}
			a {
				color: white;
			}
		</style>
	</head>
	<body>
		<h1>Index of ` + r.URL.Path + `</h1>
		<table>
			<tr><th>Name</th><th>Last Modified</th><th>Size</th></tr>`
	
	parentpath := strings.Join(strings.Split(r.URL.Path, "/")[0:len(strings.Split(r.URL.Path, "/"))-2], "/")
	if keyPrefix == "/"{
		parentpath = ""
	}
	html += `<tr><td><a href="` + parentpath + "/" + `">../</a></td><td>-</td><td>-</td></tr>`

	for _, prefix := range result.CommonPrefixes {
		html += `<tr><td><a href="` + "/" + *prefix.Prefix + `">` + *prefix.Prefix + `</a></td><td>-</td><td>-</td></tr>`
	}
	for _, obj := range result.Contents {
		href := path.Join("/", *obj.Key)
		modified := obj.LastModified.String()
		modified = strings.Split(modified, ".")[0]
		size := formatByte(*obj.Size)
		html += `<tr><td><a href="` + href + `">` + *obj.Key + `</a></td><td>` + modified + `</td><td>` + size + `</td></tr>`
	}

	html += `</table></body></html>`
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(html))
}

func formatByte (size int64) string {
	if size < 1024 {
		return strconv.FormatInt(size, 10) + " Bytes"
	}
	size = size / 1024
	if size < 1024 {
		return strconv.FormatInt(size, 10) + " KB"
	}
	size = size / 1024
	if size < 1024 {
		return strconv.FormatInt(size, 10) + " MB"
	}
	size = size / 1024
	return strconv.FormatInt(size, 10) + " GB"
}

func (h *WebDAVClient) Put(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path
	fmt.Println(key, r.Body)
	_, err := h.Backend.PutObject(key, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *WebDAVClient) Delete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path
	output, err := h.Backend.DeleteObject(key)
	fmt.Println(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
