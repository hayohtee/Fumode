package uploader

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"mime"
	"path/filepath"
	"strings"
)

// generateRandomString creates a random base64 encoded string of specified length.
func generateRandomString(length int) (string, error) {
	data := make([]byte, length)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(data), nil
}

// getFileExtension gets the file extension based on content-type or filename.
func getFileExtension(filename, contentType string) string {
	// If content type is available, derive extension from it
	if contentType != "" {
		exts, err := mime.ExtensionsByType(contentType)
		if err == nil && len(exts) > 0 {
			return exts[0]
		}
	}
	// Fallback to using the extension from the filename
	return strings.ToLower(filepath.Ext(filename))
}

// generateUniqueFilename generates a unique filename for an image with a random string and correct extension
func generateUniqueFilename(filename, contentType string) (string, error) {
	randomString, err := generateRandomString(16)
	if err != nil {
		return "", err
	}

	fileExtension := getFileExtension(filename, contentType)
	if fileExtension == "" {
		return "", errors.New("unable to determine file extension")
	}

	uniqueFilename := randomString + fileExtension
	return uniqueFilename, nil
}
