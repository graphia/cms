package main

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestNewToken(t *testing.T) {

	setupTestKeys()
	token, _ := newToken(mh)

	claims, _ := token.Claims.(jwt.MapClaims)

	// make sure that the new token has the supplied user's name
	assert.Equal(t, "misshoover", claims["sub"])

	// and we're using RS256
	assert.Equal(t, "RS256", token.Method.Alg())

	// FIXME should it be valid here?
	//assert.True(t, token.Valid)

}
