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

			// create a token and set it in the db
			token, _ := newToken(u)
			tokenString, _ := newTokenString(token)
			u.setToken(tokenString)

			payload, _ := json.Marshal(tt.user)
			buff := bytes.NewBuffer(payload)

			req, _ := http.NewRequest("PATCH", target, buff)

			req.Header.Add(
				"Authorization",
				fmt.Sprintf("Bearer %s", tokenString),
			)

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
