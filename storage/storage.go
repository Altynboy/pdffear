package storage

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// Storage defines the interface for file storage.
type Storage interface {
	Save(file multipart.File, header *multipart.FileHeader, destDir string) (string, error)
}

// LocalStorage implements Storage for local filesystem.
type LocalStorage struct{}

// NewLocalStorage creates a new LocalStorage.
func NewLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

// Save saves a multipart file to the destination directory.
func (s *LocalStorage) Save(file multipart.File, header *multipart.FileHeader, destDir string) (string, error) {
	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", err
	}

	destPath := filepath.Join(destDir, header.Filename)
	dst, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return destPath, nil
}
