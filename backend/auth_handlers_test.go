package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

const (
	testPrivKeyPath = "../tests/backend/keys/app.test.rsa"
	testPubKeyPath  = "../tests/backend/keys/app.test.rsa.pub"
)

type TokenAttributes struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

// Auth Tests

func TestAuthLoginHandler(t *testing.T) {

	setupTestKeys()

	server = httptest.NewServer(unprotectedRouter())

	var target = fmt.Sprintf("%s/%s", server.URL, "auth/login")

	mh := User{
		Username: "misshoover",
		Email:    "e.hoover@springfield.k12.us",
		Name:     "Elizabeth Hoover",
		Password: []byte("SuperSecret123"),
	}

	_ = createUser(mh)

	uc := &UserCredentials{
		Username: "misshoover",
		Password: "SuperSecret123",
	}

	payload, _ := json.Marshal(uc)

	b := bytes.NewBuffer(payload)

	client := &http.Client{}

	req, _ := http.NewRequest("POST", target, b)

	resp, _ := client.Do(req)

	var receiver Token

	json.NewDecoder(resp.Body).Decode(&receiver)

	assert.NotEmpty(t, receiver)

	decoded, _ := jwt.DecodeSegment(receiver.Token)

	var ta TokenAttributes

	if err := json.Unmarshal(decoded, &ta); err != nil {
		panic(err)
	}

	assert.Equal(t, "RS256", ta.Algorithm)
	assert.Equal(t, "JWT", ta.Type)

}

func TestAuthInvalidLoginHandler(t *testing.T) {

	setupTestKeys()

	server = httptest.NewServer(unprotectedRouter())

	var target = fmt.Sprintf("%s/%s", server.URL, "auth/login")

	mh := User{
		Username: "misshoover",
		Email:    "e.hoover@springfield.k12.us",
		Name:     "Elizabeth Hoover",
		Password: []byte("SuperSecret123"),
	}

	_ = createUser(mh)

	uc := &UserCredentials{
		Username: "misshoover",
		Password: "atotallyIncoRRecTPassw0rd",
	}

	payload, _ := json.Marshal(uc)

	b := bytes.NewBuffer(payload)

	client := &http.Client{}

	req, _ := http.NewRequest("POST", target, b)

	resp, _ := client.Do(req)

	var receiver FailureResponse

	json.NewDecoder(resp.Body).Decode(&receiver)

	assert.NotEmpty(t, receiver)
	assert.Contains(t, "Invalid credentials", receiver.Message)

}

func TestAuthNonExistantLoginHandler(t *testing.T) {

	setupTestKeys()

	server = httptest.NewServer(unprotectedRouter())

	var target = fmt.Sprintf("%s/%s", server.URL, "auth/login")

	mh := User{
		Username: "misshoover",
		Email:    "e.hoover@springfield.k12.us",
		Name:     "Elizabeth Hoover",
		Password: []byte("SuperSecret123"),
	}

	_ = createUser(mh)

	uc := &UserCredentials{
		Username: "mskrabappel",
		Password: "SuperSecret123",
	}

	payload, _ := json.Marshal(uc)

	b := bytes.NewBuffer(payload)

	client := &http.Client{}

	req, _ := http.NewRequest("POST", target, b)

	resp, _ := client.Do(req)

	var receiver FailureResponse

	json.NewDecoder(resp.Body).Decode(&receiver)

	assert.NotEmpty(t, receiver)
	assert.Contains(t, "User not found: mskrabappel", receiver.Message)

}
