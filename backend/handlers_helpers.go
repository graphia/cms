package main

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
)

// JSONResponse is a helper function to jsonify and send a response
func JSONResponse(response interface{}, status int, w http.ResponseWriter) {

	json, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func hasImageExt(uri string) bool {

	for _, extension := range config.FileCategories["images"] {
		if strings.HasSuffix(uri, extension) {
			return true
		}
	}

	return false
}

func extractImagePath(uri string) string {

	var imagePath = []string{config.Repository}

	// append the parts of the path required to serve the file
	// from any point in the hierarchy
	for _, part := range strings.Split(uri, "/") {

		// relative image paths should definitely not include the
		// `/cms/` segment, so skip it
		if part == "cms" {
			continue
		}

		// if the filename is present in the URI with the
		// markdown extension, trim it to get the corresponding
		// directory
		if strings.HasSuffix(part, ".md") {
			part = strings.TrimSuffix(part, ".md")
		}

		imagePath = append(imagePath, part)

	}

	return filepath.Join(imagePath...)
}
