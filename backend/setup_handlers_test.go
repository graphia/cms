package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Initial User Tests

func Test_SetupCreateInitialUser(t *testing.T) {

	// setup
	db.Drop("User")
	server = httptest.NewServer(unprotectedRouter())
	var target = fmt.Sprintf("%s/%s", server.URL, "setup/create_initial_user")

	mh := User{
		Username: "misshoover",
		Email:    "e.hoover@springfield.k12.us",
		Name:     "Elizabeth Hoover",
		Password: "SuperSecret123",
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

func Test_SetupCreateInitialUserWhenOneExists(t *testing.T) {

	// setup
	db.Drop("User")
	server = httptest.NewServer(unprotectedRouter())
	var target = fmt.Sprintf("%s/%s", server.URL, "setup/create_initial_user")

	mh := User{
		Username: "misshoover",
		Email:    "e.hoover@springfield.k12.us",
		Name:     "Elizabeth Hoover",
		Password: "SuperSecret123",
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

func Test_SetupCreateInitialUserWithErrors(t *testing.T) {

	// get rid of users first
	db.Drop("User")

	server = httptest.NewServer(unprotectedRouter())

	var target = fmt.Sprintf("%s/%s", server.URL, "setup/create_initial_user")

	mh := User{
		//Username: "misshoover",
		Email:    "e.hoover@springfield.k12.us",
		Name:     "Elizabeth Hoover",
		Password: "SuperSecret123",
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

func Test_SetupAllowCreateInitialUserNoUsers(t *testing.T) {
	// clear the db, 0 users exist
	db.Drop("User")

	server = httptest.NewServer(unprotectedRouter())

	target := fmt.Sprintf("%s/%s", server.URL, "setup/create_initial_user")

	resp, _ := http.Get(target)
	var is SetupOption

	json.NewDecoder(resp.Body).Decode(&is)

	// one user exists, initial setup should be allowed
	assert.Equal(t, is, SetupOption{Enabled: true})
}
