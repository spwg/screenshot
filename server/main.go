// Binary main runs the screenshot server.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spwg/screenshot/internal/datastore"
	"github.com/spwg/screenshot/internal/home"
	"github.com/spwg/screenshot/internal/image"
	"github.com/spwg/screenshot/internal/upload"
)

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run() error {
	log.Default().SetFlags(log.LstdFlags | log.Lshortfile)
	assetsDirectory := flag.String("assets_dir", "", "The directory from which to read html/css/js/ts files.")
	saveDirectory := flag.String("upload_dir", "", "The directory to save uploaded files in.")
	flag.Parse()
	directory, err := datastore.NewDirectory(*saveDirectory)
	if err != nil {
		return err
	}
	he := home.NewEndpoint(*assetsDirectory)
	ue := upload.NewEndpoint(directory)
	ie := image.NewEndpoint(directory)
	http.HandleFunc("/", he.Handler)
	http.HandleFunc("/upload", ue.Handler)
	http.HandleFunc("/image/", ie.Handler)
	addr := "localhost:10987"
	fmt.Printf("Listening on http://%v\n", addr)
	http.ListenAndServe(addr, nil)
	return nil
}
