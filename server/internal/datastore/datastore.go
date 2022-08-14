// Package datastore provides functionality to save and fetch uploads.
package datastore

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/google/uuid"
)

const (
	// ImageFileName is the filename that all images will have.
	ImageFileName = "image.png"
	// ImageURLFileName is the filename for the URL of the image.
	ImageURLFileName = "url.text"
	// filePermissions are the permissions that each saved file will have.
	filePermissions = 0644
)

// Database is the interface for database implementations.
type Database interface {
	// Save will persist the image and url. The returned uuid can be used to
	// reference the saved upload.
	Save(ctx context.Context, image []byte, url string) (uuid.UUID, error)
	// Fetch retrieves a saved upload using a uuid.
	Fetch(ctx context.Context, id uuid.UUID) (image []byte, url string, err error)
}

// Directory is an implementation of a data storage system for uploaded files
// based on a root directory. Images and metadata will be stored in
// subdirectories with UUIDs as names.
type Directory struct {
	root string
}

// NewDirectory creates a *Directory based on root.
func NewDirectory(root string) (*Directory, error) {
	if err := os.MkdirAll(root, 0777); err != nil {
		return nil, fmt.Errorf("MkdirAll %q failed: %v", root, err)
	}
	return &Directory{root}, nil
}

// Save writes the image and url into the directory and returns a UUID for it.
func (d *Directory) Save(_ context.Context, image []byte, url string) (uuid.UUID, error) {
	id := uuid.New()
	p := path.Join(d.root, id.String())
	if err := os.Mkdir(p, 0777); err != nil {
		return uuid.UUID{}, fmt.Errorf("mkdir %q failed: %v", p, err)
	}
	imagePath := path.Join(p, ImageFileName)
	if err := os.WriteFile(imagePath, image, filePermissions); err != nil {
		return uuid.UUID{}, fmt.Errorf("write file %q failed: %v", imagePath, err)
	}
	urlPath := path.Join(p, ImageURLFileName)
	if err := os.WriteFile(urlPath, []byte(url), filePermissions); err != nil {
		return uuid.UUID{}, fmt.Errorf("write file %q failed: %v", urlPath, err)
	}
	return id, nil
}

// Fetch retrieves an upload.
func (d *Directory) Fetch(_ context.Context, id uuid.UUID) ([]byte, string, error) {
	imagePath := path.Join(d.root, id.String(), ImageFileName)
	urlPath := path.Join(d.root, id.String(), ImageURLFileName)
	image, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, "", fmt.Errorf("read file %q failed: %v", imagePath, err)
	}
	urlBytes, err := os.ReadFile(urlPath)
	if err != nil {
		return nil, "", fmt.Errorf("read file %q failed: %v", urlPath, err)
	}
	url := string(urlBytes)
	return image, url, nil
}
