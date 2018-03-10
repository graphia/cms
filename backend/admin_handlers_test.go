package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
)

func setupMiddlewareAdminTestServer() *httptest.Server {

	// The Admin MW depends on the user being set by the Token MW,
	// so both are chained
	n := negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.HandlerFunc(ValidateAdminMiddleware),
		negroni.Wrap(adminRouter()),
	)

	server := httptest.NewServer(n)
	return server
}

func TestAdminMiddleware(t *testing.T) {

	repoPath := "../tests/tmp/repositories/auth_handlers"
	setupSmallTestRepo(repoPath)

	server := setupMiddlewareAdminTestServer()

	target := fmt.Sprintf("%s/%s", server.URL, "api/admin/users/cookie.kwan")

	tests := []struct {
		name    string
		user    User
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Admin user",
			user:    ck,
			wantErr: false,
		},
		{
			name:    "Non-Admin user",
			user:    mh,
			wantErr: true,
			errMsg:  "User does not have sufficient privileges",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db.Drop("User")
			setupTestKeys()

			_ = createUser(tt.user)
			u, _ := getUserByUsername(tt.user.Username)

			client := &http.Client{}

			payload, _ := json.Marshal(tt.user)
			buff := bytes.NewBuffer(payload)

			req, _ := http.NewRequest("PATCH", target, buff)

			req = authorizeRequest(u, req)

			resp, _ := client.Do(req)

			if tt.wantErr {
				var fr FailureResponse
				json.NewDecoder(resp.Body).Decode(&fr)

				assert.Equal(t, tt.errMsg, fr.Message)
				return
			}

			var sr SuccessResponse
			json.NewDecoder(resp.Body).Decode(&sr)

			assert.Equal(t, http.StatusCreated, resp.StatusCode)

		})
	}

}

func authorizeRequest(user User, req *http.Request) *http.Request {

	// create a token and set it in the db
	token, _ := newToken(user)
	tokenString, _ := newTokenString(token)
	user.setToken(tokenString)

	// and add the auth header
	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Bearer %s", tokenString),
	)

	return req
}

func TestAPICreateUser(t *testing.T) {

	// initial setup
	server := setupMiddlewareAdminTestServer()
	setupTestKeys()

	// clear db
	db.Drop("User")

	// create an admin
	createUser(ck)
	admin, _ := getUserByUsername(ck.Username)

	target := fmt.Sprintf("%s/%s", server.URL, "api/admin/users")

	payload, _ := json.Marshal(mh)
	buff := bytes.NewBuffer(payload)

	req, _ := http.NewRequest("POST", target, buff)

	req = authorizeRequest(admin, req)

	client := &http.Client{}
	Debug.Println(client)
	resp, _ := client.Do(req)

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

	// initial setup
	server := setupMiddlewareAdminTestServer()
	setupTestKeys()

	// clear db
	db.Drop("User")

	// create an admin
	createUser(ck)
	admin, _ := getUserByUsername(ck.Username)

	var target = fmt.Sprintf("%s/%s", server.URL, "api/admin/users")

	invalid := User{
		//Username: "misshoover",
		Email:    "e.hoover@springfield.k12.us",
		Name:     "Elizabeth Hoover",
		Password: "SuperSecret123",
	}

	payload, _ := json.Marshal(invalid)

	buff := bytes.NewBuffer(payload)

	client := &http.Client{}

	req, _ := http.NewRequest("POST", target, buff)
	req = authorizeRequest(admin, req)

	resp, _ := client.Do(req)

	var receiver map[string]string

	json.NewDecoder(resp.Body).Decode(&receiver)

	// API should return the right response
	assert.Contains(t, receiver["Username"], "is a required field")

	// Should still equal 1 as our second attempt failed
	au, _ := allUsers()
	assert.Equal(t, 1, len(au))

}
