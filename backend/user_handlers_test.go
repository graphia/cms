package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"io/ioutil"

	"github.com/stretchr/testify/assert"
)

func TestAPIListUsers(t *testing.T) {
	server = httptest.NewServer(protectedRouter())
	db.Drop("User")

	_ = createUser(mh)
	_ = createUser(ck)

	target := fmt.Sprintf("%s/%s", server.URL, "api/users")

	client := &http.Client{}

	req, _ := http.NewRequest("GET", target, nil)

	resp, _ := client.Do(req)

	var receiver []User

	json.NewDecoder(resp.Body).Decode(&receiver)

	assert.Equal(t, 2, len(receiver))

}

func TestAPIGetUser(t *testing.T) {
	server = httptest.NewServer(protectedRouter())
	db.Drop("User")

	_ = createUser(mh)

	target := fmt.Sprintf("%s/%s", server.URL, "api/users/misshoover")
	client := &http.Client{}
	req, _ := http.NewRequest("GET", target, nil)
	resp, _ := client.Do(req)

	buf := new(bytes.Buffer)

	// use TeeReader so we can read the buffer more than once
	tee := io.TeeReader(resp.Body, buf)

	var receiver User
	json.NewDecoder(tee).Decode(&receiver)
	Debug.Println("receiver", receiver)

	// ensure the result contains the correct information
	assert.Equal(t, mh.Username, receiver.Username)

	// ensure password isn't included in the raw result
	raw, _ := ioutil.ReadAll(buf)
	assert.NotContains(t, string(raw), "password")
}

func TestAPICreateUser(t *testing.T) {
	server = httptest.NewServer(protectedRouter())
	db.Drop("User")

	target := fmt.Sprintf("%s/%s", server.URL, "api/users")

	payload, _ := json.Marshal(mh)
	buff := bytes.NewBuffer(payload)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", target, buff)
	resp, _ := client.Do(req)

	Debug.Println(resp)

	// should get a 201 back if everything's successful
	assert.Equal(t, 201, resp.StatusCode)

	// make sure new user has been properly created
	user, _ := getUserByUsername(mh.Username)
	assert.Equal(t, mh.Username, user.Username)
	assert.Equal(t, mh.Name, user.Name)
	assert.Equal(t, mh.Email, user.Email)
	assert.True(t, user.Active)

	// ensure passwords don't match (i.e. encrypted)
	assert.NotEqual(t, mh.Password, user.Password)

}

func TestAPICreateUserWithErrors(t *testing.T) {

	// get rid of users first
	db.Drop("User")

	server = httptest.NewServer(protectedRouter())

	var target = fmt.Sprintf("%s/%s", server.URL, "api/users")

	invalid := User{
		//Username: "misshoover",
		Email:    "e.hoover@springfield.k12.us",
		Name:     "Elizabeth Hoover",
		Password: []byte("SuperSecret123"),
	}

	payload, _ := json.Marshal(invalid)

	buff := bytes.NewBuffer(payload)

	client := &http.Client{}

	req, _ := http.NewRequest("POST", target, buff)

	resp, _ := client.Do(req)

	var receiver map[string]string

	json.NewDecoder(resp.Body).Decode(&receiver)

	// API should return the right response
	assert.Contains(t, receiver["Username"], "is a required field")

	// Should still equal 1 as our second attempt failed
	au, _ := allUsers()
	assert.Equal(t, 0, len(au))

}
