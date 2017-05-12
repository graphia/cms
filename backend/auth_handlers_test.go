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
	"gopkg.in/go-playground/validator.v9"
)

const (
	testPrivKeyPath = "../tests/backend/keys/app.test.rsa"
	testPubKeyPath  = "../tests/backend/keys/app.test.rsa.pub"
)

type TokenAttributes struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

func init() {
	validate = validator.New()
}

// Auth Tests

func TestAuthLoginHandler(t *testing.T) {
	// setup
	db.Drop("User")
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

	// Ensure token has the right attributes
	assert.Equal(t, "RS256", ta.Algorithm)
	assert.Equal(t, "JWT", ta.Type)

	// Make sure that we've set the user's TokenString to equal
	// the returned Token
	user, _ := getUserByUsername("misshoover")
	assert.Equal(t, receiver.Token, user.TokenString)

}

func TestAuthInvalidLoginHandler(t *testing.T) {

	// setup
	db.Drop("User")
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

	// Make sure that we've not set the user's TokenString
	user, _ := getUserByUsername("misshoover")
	assert.Empty(t, user.TokenString)

}

func TestAuthNonExistantLoginHandler(t *testing.T) {

	// setup
	db.Drop("User")
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

func TestAuthCreateInitialUser(t *testing.T) {

	// setup
	db.Drop("User")
	server = httptest.NewServer(unprotectedRouter())
	var target = fmt.Sprintf("%s/%s", server.URL, "auth/create_initial_user")

	mh := User{
		Username: "misshoover",
		Email:    "e.hoover@springfield.k12.us",
		Name:     "Elizabeth Hoover",
		Password: []byte("SuperSecret123"),
	}

	payload, _ := json.Marshal(mh)

	b := bytes.NewBuffer(payload)

	client := &http.Client{}

	req, _ := http.NewRequest("POST", target, b)

	resp, _ := client.Do(req)

	var receiver SuccessResponse

	json.NewDecoder(resp.Body).Decode(&receiver)

	// API should return the right response
	assert.Contains(t, "User created", receiver.Message)

	// Only one user should have been created
	au, _ := allUsers()
	assert.Equal(t, 1, len(au))

	u, _ := getUserByUsername(mh.Username)

	// The user's credentials should match the input
	assert.Equal(t, mh.Username, u.Username)
	assert.Equal(t, mh.Name, u.Name)
	assert.Equal(t, mh.Email, u.Email)
	assert.True(t, u.Active)

	// Except the passwors should now be encrypted
	assert.NotEqual(t, mh.Password, u.Password)
}

func TestAuthCreateInitialUserWhenOneExists(t *testing.T) {

	// setup
	db.Drop("User")
	server = httptest.NewServer(unprotectedRouter())
	var target = fmt.Sprintf("%s/%s", server.URL, "auth/create_initial_user")

	mh := User{
		Username: "misshoover",
		Email:    "e.hoover@springfield.k12.us",
		Name:     "Elizabeth Hoover",
		Password: []byte("SuperSecret123"),
	}

	// First ensure a user exists
	_ = createUser(mh)

	payload, _ := json.Marshal(mh)

	b := bytes.NewBuffer(payload)

	client := &http.Client{}

	req, _ := http.NewRequest("POST", target, b)

	resp, _ := client.Do(req)

	var receiver FailureResponse

	json.NewDecoder(resp.Body).Decode(&receiver)

	// API should return the right response
	assert.Contains(t, "Users already exist. Cannot create initial user", receiver.Message)

	// Should still equal 1 as our second attempt failed
	au, _ := allUsers()
	assert.Equal(t, 1, len(au))

}

func TestAuthCreateInitialUserWithErrors(t *testing.T) {

	// get rid of users first
	db.Drop("User")

	server = httptest.NewServer(unprotectedRouter())

	var target = fmt.Sprintf("%s/%s", server.URL, "auth/create_initial_user")

	mh := User{
		//Username: "misshoover",
		Email:    "e.hoover@springfield.k12.us",
		Name:     "Elizabeth Hoover",
		Password: []byte("SuperSecret123"),
	}

	payload, _ := json.Marshal(mh)

	b := bytes.NewBuffer(payload)

	client := &http.Client{}

	req, _ := http.NewRequest("POST", target, b)

	resp, _ := client.Do(req)

	var receiver map[string]string

	json.NewDecoder(resp.Body).Decode(&receiver)

	// API should return the right response
	assert.Contains(t, receiver["Username"], "is a required field")

	// Should still equal 1 as our second attempt failed
	au, _ := allUsers()
	assert.Equal(t, 0, len(au))

}
