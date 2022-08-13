package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/spwg/screenshot/internal/home"
	"github.com/spwg/screenshot/internal/upload"
)

func main() {
	assetsDirectory := flag.String("assets_dir", "", "The directory from which to read html/css/js/ts files.")
	flag.Parse()
	he := home.NewEndpoint(*assetsDirectory)
	http.HandleFunc("/", he.Handle)
	http.HandleFunc("/upload", upload.HandleUpload)
	addr := "localhost:10987"
	fmt.Printf("Listening on http://%v\n", addr)
	http.ListenAndServe(addr, nil)
}
