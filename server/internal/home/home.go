// Package home provides an http endpoint implementation for the home page.
package home

import (
	"log"
	"net/http"
	"os"
	"path"
)

// Endpoint is the home page endpoint.
type Endpoint struct {
	assetsDirectory string
}

// NewEndpoint creates an *Endpoint.
func NewEndpoint(assetsDirectory string) *Endpoint {
	return &Endpoint{assetsDirectory}
}

// Handler is an http handler implementation for the home page endpoint.
func (e *Endpoint) Handler(resp http.ResponseWriter, req *http.Request) {
	if !(req.URL.Path == "" || req.URL.Path == "/") {
		http.NotFound(resp, req)
		return
	}
	name := path.Join(e.assetsDirectory, "home.html")
	b, err := os.ReadFile(name)
	if err != nil {
		log.Printf("Failed to read %q: %v", name, err)
		http.Error(resp, "Failed to read html file.", http.StatusInternalServerError)
		return
	}
	if _, err := resp.Write(b); err != nil {
		log.Print(err)
		return
	}
}
