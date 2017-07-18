package main

import (
	"encoding/json"
	"net/http"
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
