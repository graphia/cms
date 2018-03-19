package main

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/rs/xid"
)

type userKey int

const (
	// HoursUntilExpiry is number of hours until token expires
	HoursUntilExpiry         = 1
	currentUser      userKey = 0
)

// ValidateTokenMiddleware validates the JWT token
func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var usernameFromToken string
	var fr FailureResponse

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})

	if err != nil {
		fr = FailureResponse{Message: "Unauthorized access, invalid JWT"}
		JSONResponse(fr, http.StatusUnauthorized, w)
		return
	}

	// some checks to make sure the request is from a valid source
	usernameFromToken = token.Claims.(jwt.MapClaims)["sub"].(string)
	Debug.Println("Retrieved username from token", usernameFromToken)

	// check that the user still exists, if this fails it's probably been deleted
	user, err := getUserByUsername(usernameFromToken)

	if err != nil {
		fr = FailureResponse{Message: "Cannot match user with token"}
		JSONResponse(fr, http.StatusUnauthorized, w)
		return
	}

	// check that the user is active in the system, if they're not disallow
	if !user.Active {
		fr = FailureResponse{Message: "User has been deactivated"}
		JSONResponse(fr, http.StatusUnauthorized, w)
		return
	}

	// check that the token matches the found user's stored token, if it doesn't
	// it's likely that the user has logged in again and is using the *old* token
	if user.TokenString != token.Raw {
		fr = FailureResponse{Message: "Token is out of date"}
		JSONResponse(fr, http.StatusUnauthorized, w)
		return
	}

	if !token.Valid {
		fr = FailureResponse{Message: "Invalid credentials"}
		JSONResponse(fr, http.StatusUnauthorized, w)
		return
	}

	ctx := newContextWithCurrentUser(r.Context(), r, user)
	next(w, r.WithContext(ctx))

}

// ValidateAdminMiddleware ensures the user is an administrator
func ValidateAdminMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var user User
	var fr FailureResponse

	user = getCurrentUser(r.Context())

	if !user.Admin {
		Warning.Printf("Regular user %s trying to access admin area", user.Username)
		fr = FailureResponse{Message: "User does not have sufficient privileges"}
		JSONResponse(fr, http.StatusUnauthorized, w)
		return
	}

	Debug.Println("User is admin, continuing")

	next(w, r)

}

// Add the constant with key currentUser to the context
func newContextWithCurrentUser(ctx context.Context, req *http.Request, user User) context.Context {
	return context.WithValue(ctx, currentUser, user)
}

// Retrive the current user from the request's context
func getCurrentUser(ctx context.Context) User {
	return ctx.Value(currentUser).(User)
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
