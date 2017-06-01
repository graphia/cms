package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/rs/xid"
)

const (
	// HoursUntilExpiry is number of hours until token expires
	HoursUntilExpiry = 1
)

// ValidateTokenMiddleware validates the JWT token
func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var usernameFromToken string

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := FailureResponse{Message: "Unauthorized access, invalid JWT"}
		json, err := json.Marshal(response)
		if err != nil {
			panic(err)
		}
		w.Write(json)
		return
	}

	// some checks to make sure the request is from a valid source
	usernameFromToken = token.Claims.(jwt.MapClaims)["sub"].(string)
	Debug.Println("Retrieved username from token", usernameFromToken)

	// check that the user still exists, if this fails it's probably been deleted
	user, err := getUserByUsername(usernameFromToken)
	if err != nil {
		response := FailureResponse{Message: "Cannot match user with token"}
		json, err := json.Marshal(response)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(json)
		return
	}

	// check that the user is active in the system, if they're not disallow
	if !user.Active {
		response := FailureResponse{Message: "User has been deactivated"}
		json, err := json.Marshal(response)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(json)
		return
	}

	// check that the token matches the found user's stored token, if it doesn't
	// it's likely that the user has logged in again and is using the *old* token
	if user.TokenString != token.Raw {
		response := FailureResponse{Message: "Token is out of date"}
		json, err := json.Marshal(response)
		if err != nil {
			Debug.Println("WTFBBQ")

			panic(err)
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(json)
		return
	}

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

func newToken(user User) (jwt.Token, error) {

	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)

	// JWT claims as per https://tools.ietf.org/html/rfc7519#section-4.1

	// token expires at
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(HoursUntilExpiry)).Unix()

	// token issued at
	claims["iat"] = time.Now().Unix()

	// username
	claims["sub"] = user.Username

	// unique ID for this token
	claims["jti"] = xid.New()

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