package home

import (
	"log"
	"net/http"
	"os"
	"path"
)

type Endpoint struct {
	assetsDirectory string
}

func NewEndpoint(assetsDirectory string) *Endpoint {
	return &Endpoint{assetsDirectory}
}

func (e *Endpoint) Handle(resp http.ResponseWriter, req *http.Request) {
	name := path.Join(e.assetsDirectory, "home.html")
	b, err := os.ReadFile(name)
	if err != nil {
		log.Printf("Failed to read %q: %v", name, err)
		http.Error(resp, "Failed to read html file.", http.StatusInternalServerError)
		return
	}
	resp.Write(b)
}
