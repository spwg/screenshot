// Package image provides an http endpoint for fetching images.
package image

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/google/uuid"
	"github.com/spwg/screenshot/internal/datastore"
)

var page = template.Must(template.New("").Parse(`
<html>
<body>
<a href="{{.ImageSRC}}"><img src="{{.ImageSRC}}"/></a>
<p><a href="//{{.ImageURL}}" target="_blank">{{.ImageURL}}</a>
</body>
</html>
`))

// Endpoint is an http endpoint for displaying images.
type Endpoint struct {
	db datastore.Database
}

// NewEndpoint returns an *Endpoint backed by the db.
func NewEndpoint(db datastore.Database) *Endpoint {
	return &Endpoint{db}
}

func (e *Endpoint) Handler(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	split := strings.Split(strings.TrimPrefix(req.URL.Path, "/"), "/")
	found := len(split) == 2
	if !found {
		http.Error(resp, fmt.Sprintf("Invalid path %q.", req.URL.Path), http.StatusBadRequest)
		return
	}
	s := strings.TrimSuffix(split[1], ".png")
	id, err := uuid.Parse(s)
	if err != nil {
		log.Print(err)
		http.Error(resp, fmt.Sprintf("%q is not a valid ID.", s), http.StatusBadRequest)
		return
	}
	image, url, err := e.db.Fetch(ctx, id)
	if err != nil {
		log.Print(err)
		http.NotFound(resp, req)
		return
	}
	if strings.HasSuffix(split[1], ".png") {
		if _, err := resp.Write(image); err != nil {
			log.Print(err)
		}
		return
	}
	if err := page.Execute(resp, struct {
		ImageSRC string
		ImageURL string
	}{
		"/image/" + id.String() + ".png",
		url,
	}); err != nil {
		log.Print(err)
		http.Error(resp, "Failed to render html page.", http.StatusInternalServerError)
		return
	}
}
