package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mh, ck User
)

func init() {
	db = setupDBForTests(config.Database)
	mh = User{
		Username: "misshoover",
		Email:    "e.hoover@springfield.k12.us",
		Name:     "Elizabeth Hoover",
	}
	ck = User{
		Username: "cookie.kwan",
		Email:    "cookie@springfield-realtors.com",
		Name:     "Cookie Kwan",
	}
}

func TestCreateUser(t *testing.T) {

	returnedUser, err := createUser(mh)

	if err != nil {
		fmt.Println(err)
	}

	retrievedUser := User{}
	_ = db.One("Username", "misshoover", &retrievedUser)

	assert.Equal(t, returnedUser.Name, retrievedUser.Name)
	assert.Equal(t, returnedUser.Username, retrievedUser.Username)
	assert.Equal(t, returnedUser.Email, retrievedUser.Email)

}
func TestCreateDuplicateUser(t *testing.T) {
	var err error

	_, err = createUser(mh)
	_, err = createUser(mh)

	// make sure error message is correct
	assert.Contains(t, err.Error(), "already exists")

	// and check that only one misshoover was created
	criteria := User{Username: "misshoover"}
	matches, _ := db.Count(&criteria)
	assert.Equal(t, matches, 1)
}

func TestFindUser(t *testing.T) {
	_, _ = createUser(ck)

	cookieKwan, _ := getUser("cookie.kwan")
	assert.NotZero(t, cookieKwan.ID)

	_, err := getUser("not.miss.hoover.ok")
	assert.Contains(t, err.Error(), "not found")

}
