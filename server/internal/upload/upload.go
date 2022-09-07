// Package upload provides an http endpoint for uploading images.
package upload

import (
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/spwg/screenshot/internal/datastore"
)

// Endpoint is the implementation of the upload endpoint. It accepts images and
// associated metadata and saves them.
type Endpoint struct {
	db datastore.Database
}

// NewEndpoint creates an *Endpoint backed by the db.
func NewEndpoint(db datastore.Database) *Endpoint {
	return &Endpoint{db}
}

// Handler is an HTTP handler that accepts uploaded files.
func (e *Endpoint) Handler(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if code, err := validateRequest(req); err != nil {
		http.Error(resp, err.Error(), code)
		return
	}
	if err := req.ParseMultipartForm(32 << 10); err != nil {
		log.Print(err)
		http.Error(resp, "Failed to parse request form.", http.StatusBadRequest)
		return
	}
	file, header, err := req.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			http.Error(resp, "No file was uploaded.", http.StatusBadRequest)
			return
		}
		log.Print(err)
		http.Error(resp, "Failed to parse uploaded file.", http.StatusInternalServerError)
		return
	}
	if contentType := header.Header.Get("Content-Type"); contentType != "image/png" {
		http.Error(resp, fmt.Sprintf("%q is not %q, it's %q.", escape(header.Filename), "image/png", contentType), http.StatusBadRequest)
		return
	}
	log.Printf("%q %v bytes", escape(header.Filename), header.Size)
	b, err := io.ReadAll(file)
	if err != nil {
		log.Print(err)
		http.Error(resp, "Failed to read uploaded file.", http.StatusBadRequest)
		return
	}
	url := req.Form.Get("url")
	id, err := e.db.Save(ctx, b, url)
	if err != nil {
		log.Print(err)
		http.Error(resp, "Failed to save the upload.", http.StatusInternalServerError)
		return
	}
	http.Redirect(resp, req, "/image/"+id.String(), http.StatusTemporaryRedirect)
}

// validateRequest makes sure that the request is POST.
func validateRequest(req *http.Request) (int, error) {
	if req.Method != http.MethodPost {
		return http.StatusMethodNotAllowed, fmt.Errorf("http method %q not allowed", req.Method)
	}
	return http.StatusOK, nil
}

func escape(s string) string {
	s = html.EscapeString(s)
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}
