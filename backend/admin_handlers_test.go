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

			assert.Equal(t, http.StatusOK, resp.StatusCode)

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

func Test_apiUpdateUserHandler(t *testing.T) {

	repoPath := "../tests/tmp/repositories/auth_handlers"
	setupSmallTestRepo(repoPath)

	server := setupMiddlewareAdminTestServer()

	tests := []struct {
		name              string
		updates           LimitedUser
		targetSuffix      string
		wantErr           bool
		errMsg            string
		code              int
		wantValidationErr bool
		validationErrors  map[string]string
	}{
		{
			name:         "Successful Update",
			targetSuffix: "api/admin/users/selma.bouvier",
			updates: LimitedUser{
				Username: "selma.mcclure",
				Name:     "Selma McClure",
				Email:    "selma.mcclure@macguyver-fans.org",
				Active:   true,
				Admin:    true,
			},
			wantErr: false,
			code:    http.StatusOK,
		},
		{
			name:         "User not found",
			targetSuffix: "api/admin/users/edna.krabappel",
			updates:      LimitedUser{},
			wantErr:      true,
			errMsg:       "No user edna.krabappel",
			code:         http.StatusNotFound,
		},
		{
			name:         "Validation failures",
			targetSuffix: "api/admin/users/selma.bouvier",
			updates: LimitedUser{
				Name:  "Se",       // too short
				Email: "abc@.com", // not an email address
			},
			code:              http.StatusBadRequest,
			wantValidationErr: true,
			validationErrors: map[string]string{
				"Name":  "must be at least 3 characters",
				"Email": "is not a valid email address",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db.Drop("User")
			setupTestKeys()

			// add the admin who will be performing the update
			_ = createUser(ck)
			updater, _ := getUserByUsername("cookie.kwan")

			// create the user who will be updated and
			// stash the original record so we can check the
			// updates worked later
			_ = createUser(sb)
			prior, _ := getUserByUsername("selma.bouvier")

			Debug.Println("diff", tt.updates.Email, prior.Email)

			client := &http.Client{}

			payload, _ := json.Marshal(tt.updates)
			buff := bytes.NewBuffer(payload)

			target := fmt.Sprintf("%s/%s", server.URL, tt.targetSuffix)

			req, _ := http.NewRequest("PATCH", target, buff)

			req = authorizeRequest(updater, req)

			resp, _ := client.Do(req)

			assert.Equal(t, tt.code, resp.StatusCode)

			if tt.wantValidationErr {

				actualErrors := make(map[string]string)
				json.NewDecoder(resp.Body).Decode(&actualErrors)

				assert.Equal(t, tt.validationErrors, actualErrors)
				return
			}

			if tt.wantErr {
				var fr FailureResponse
				json.NewDecoder(resp.Body).Decode(&fr)

				assert.Equal(t, tt.errMsg, fr.Message)
				assert.Equal(t, tt.code, resp.StatusCode)
				return
			}

			// make sure all allowed fields are being updated in
			// this request
			assert.NotEqual(t, tt.updates.Email, prior.Email)
			assert.NotEqual(t, tt.updates.Name, prior.Name)
			assert.NotEqual(t, tt.updates.Admin, prior.Admin)
			assert.NotEqual(t, tt.updates.Active, prior.Active)

			var sr SuccessResponse

			json.NewDecoder(resp.Body).Decode(&sr)

			// make sure attributes have changed
			subsequent, _ := getUserByID(prior.ID)

			assert.Equal(t, tt.updates.Email, subsequent.Email)
			assert.Equal(t, tt.updates.Name, subsequent.Name)
			assert.Equal(t, tt.updates.Admin, subsequent.Admin)
			assert.Equal(t, tt.updates.Active, subsequent.Active)

			// make sure username hasn't changed if it is supplied
			assert.Equal(t, "selma.bouvier", prior.Username)

		})
	}

}
