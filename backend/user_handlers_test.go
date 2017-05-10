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
