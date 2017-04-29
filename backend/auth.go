package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/dgrijalva/jwt-go/request"
)

// ValidateTokenMiddleware validates the JWT token
func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})

	if err == nil {
		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			response := FailureResponse{Message: "Invalid credentials"}
			json, err := json.Marshal(response)
			if err != nil {
				panic(err)
			}
			w.Write(json)
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized access to this resource")
	}

}

// JSONResponse is a helper function to jsonify and send a response
func JSONResponse(response interface{}, w http.ResponseWriter) {

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func newToken() (jwt.Token, error) {

	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()

	// TODO add user-specific claims

	token.Claims = claims

	return *token, nil

}

func newTokenString(token jwt.Token) (string, error) {
	ts, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}
	return ts, err
}
