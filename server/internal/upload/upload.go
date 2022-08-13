package upload

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// HandleUpload is an HTTP handler that accepts uploaded files.
//
// TODO: also read the URL of the screenshot from the form.
// TODO: save the file and redirect to the endpoint to display it.
func HandleUpload(resp http.ResponseWriter, req *http.Request) {
	if code, err := validateRequest(req); err != nil {
		http.Error(resp, err.Error(), code)
		return
	}
	if err := req.ParseMultipartForm(32 << 10); err != nil {
		log.Printf("%v", err)
		http.Error(resp, "Failed to parse request form.", http.StatusBadRequest)
		return
	}
	file, header, err := req.FormFile("file1")
	if err != nil {
		log.Printf("%v", err)
		http.Error(resp, "Failed to read form file.", http.StatusBadRequest)
		return
	}

	if _, err := io.ReadAll(file); err != nil {
		log.Printf("%v", err)
		http.Error(resp, "Failed to read uploaded file.", http.StatusBadRequest)
		return
	}
	io.WriteString(resp, fmt.Sprintf("Read %v bytes.", header.Size))
}

// validateRequest makes sure that the request is POST.
func validateRequest(req *http.Request) (int, error) {
	if req.Method != http.MethodPost {
		return http.StatusMethodNotAllowed, fmt.Errorf("http method %q not allowed", req.Method)
	}
	return http.StatusOK, nil
}
