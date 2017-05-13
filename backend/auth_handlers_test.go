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
	"github.com/urfave/negroni"
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
		Password: "SuperSecret123",
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
		Password: "SuperSecret123",
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
		Password: "SuperSecret123",
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

// Initial User Tests

func TestAuthCreateInitialUser(t *testing.T) {

	// setup
	db.Drop("User")
	server = httptest.NewServer(unprotectedRouter())
	var target = fmt.Sprintf("%s/%s", server.URL, "auth/create_initial_user")

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

func TestAuthCreateInitialUserWhenOneExists(t *testing.T) {

	// setup
	db.Drop("User")
	server = httptest.NewServer(unprotectedRouter())
	var target = fmt.Sprintf("%s/%s", server.URL, "auth/create_initial_user")

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

func TestAuthCreateInitialUserWithErrors(t *testing.T) {

	// get rid of users first
	db.Drop("User")

	server = httptest.NewServer(unprotectedRouter())

	var target = fmt.Sprintf("%s/%s", server.URL, "auth/create_initial_user")

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

// Auth middleware tests

func setupMiddlewareProtectedTestServer() *httptest.Server {
	n := negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(protectedRouter()),
	)

	server := httptest.NewServer(n)
	return server
}

func TestProtectedMiddlewareWithToken(t *testing.T) {

	server := setupMiddlewareProtectedTestServer()

	db.Drop("User")
	setupTestKeys()

	repoPath := "../tests/tmp/repositories/auth_handlers"
	setupSmallTestRepo(repoPath)

	_ = createUser(ck)
	cookieKwan, _ := getUserByUsername("cookie.kwan")

	token, _ := newToken(cookieKwan)
	tokenString, _ := newTokenString(token)
	setToken(cookieKwan, tokenString)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files/document_1.md")

	//	Debug.Println(token)
	client := &http.Client{}

	req, _ := http.NewRequest("GET", target, nil)
	//resp, _ := http.Get(target)
	Debug.Println(token.Raw)

	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Bearer %s", tokenString),
	)

	resp, _ := client.Do(req)
	var file File
	json.NewDecoder(resp.Body).Decode(&file)

	// ensure the file 'looks' correct
	assert.Contains(t, file.Filename, "document_1.md")
	assert.Contains(t, file.Path, "documents")
}

func TestProtectedMiddlewareNoToken(t *testing.T) {

	db.Drop("User")
	setupTestKeys()

	repoPath := "../tests/tmp/repositories/auth_handlers"
	setupSmallTestRepo(repoPath)

	server := setupMiddlewareProtectedTestServer()

	_ = createUser(ck)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files/document_1.md")

	//	Debug.Println(token)

	resp, _ := http.Get(target)
	var fr FailureResponse

	json.NewDecoder(resp.Body).Decode(&fr)

	assert.Contains(t, fr.Message, "Unauthorized")
}

func TestProtectedMiddlewareOutdatedToken(t *testing.T) {

	server := setupMiddlewareProtectedTestServer()

	db.Drop("User")
	setupTestKeys()

	repoPath := "../tests/tmp/repositories/auth_handlers"
	setupSmallTestRepo(repoPath)

	_ = createUser(ck)
	cookieKwan, _ := getUserByUsername("cookie.kwan")

	// create the first token and assign it to the user
	tokenOne, _ := newToken(cookieKwan)
	tokenOneString, _ := newTokenString(tokenOne)
	setToken(cookieKwan, tokenOneString)

	// now create a second one which we'll attempt to connect with
	tokenTwo, _ := newToken(cookieKwan)
	tokenTwoString, _ := newTokenString(tokenTwo)
	setToken(cookieKwan, tokenTwoString)

	cookieKwan, _ = getUserByUsername("cookie.kwan")

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files/document_1.md")

	//	Debug.Println(token)
	client := &http.Client{}

	req, _ := http.NewRequest("GET", target, nil)

	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Bearer %s", tokenOneString),
	)

	resp, _ := client.Do(req)

	var fr FailureResponse
	json.NewDecoder(resp.Body).Decode(&fr)
	Debug.Println(resp.Body)

	// ensure response claims token is out of date

	assert.Equal(t, resp.StatusCode, 401)
	assert.Contains(t, fr.Message, "Token is out of date")
}

func TestProtectedMiddlewareDeletedUser(t *testing.T) {

	server := setupMiddlewareProtectedTestServer()

	db.Drop("User")
	setupTestKeys()

	repoPath := "../tests/tmp/repositories/auth_handlers"
	setupSmallTestRepo(repoPath)

	_ = createUser(ck)
	cookieKwan, _ := getUserByUsername("cookie.kwan")

	// create the first token and assign it to the user
	token, _ := newToken(cookieKwan)
	tokenString, _ := newTokenString(token)
	setToken(cookieKwan, tokenString)

	_ = deleteUser(cookieKwan)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files/document_1.md")

	//	Debug.Println(token)
	client := &http.Client{}

	req, _ := http.NewRequest("GET", target, nil)

	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Bearer %s", tokenString),
	)

	resp, _ := client.Do(req)

	var fr FailureResponse
	json.NewDecoder(resp.Body).Decode(&fr)
	Debug.Println(resp.Body)

	// ensure response claims token is out of date

	assert.Equal(t, resp.StatusCode, 401)
	assert.Contains(t, fr.Message, "Cannot match user with token")
}

func TestProtectedMiddlewareDeactivatedUser(t *testing.T) {

	server := setupMiddlewareProtectedTestServer()

	db.Drop("User")
	setupTestKeys()

	repoPath := "../tests/tmp/repositories/auth_handlers"
	setupSmallTestRepo(repoPath)

	_ = createUser(ck)
	cookieKwan, _ := getUserByUsername("cookie.kwan")

	// create the first token and assign it to the user
	token, _ := newToken(cookieKwan)
	tokenString, _ := newTokenString(token)
	setToken(cookieKwan, tokenString)

	_ = deactivateUser(cookieKwan)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files/document_1.md")

	//	Debug.Println(token)
	client := &http.Client{}

	req, _ := http.NewRequest("GET", target, nil)

	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Bearer %s", tokenString),
	)

	resp, _ := client.Do(req)

	var fr FailureResponse
	json.NewDecoder(resp.Body).Decode(&fr)
	Debug.Println(resp.Body)

	// ensure response claims token is out of date

	assert.Equal(t, resp.StatusCode, 401)
	assert.Contains(t, fr.Message, "User has been deactivated")
}
