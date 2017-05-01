package main

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

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
		Password: []byte("SuperSecret123"),
	}
	ck = User{
		Username: "cookie.kwan",
		Email:    "cookie@springfield-realtors.com",
		Name:     "Cookie Kwan",
		Password: []byte("P@ssword!"),
	}
}

func TestCreateUser(t *testing.T) {

	_ = createUser(mh)

	retrievedUser := User{}
	_ = db.One("Username", "misshoover", &retrievedUser)

	assert.Equal(t, mh.Name, retrievedUser.Name)
	assert.Equal(t, mh.Username, retrievedUser.Username)
	assert.Equal(t, mh.Email, retrievedUser.Email)

}

func TestPasswordIsEncrypted(t *testing.T) {

	initialPassword := mh.Password

	retrievedUser := User{}
	_ = db.One("Username", "misshoover", &retrievedUser)

	// ensure that the value written to the db has been modified
	assert.NotEqual(t, initialPassword, retrievedUser.Password)

	// ensure comparing our generated hash with the initial password doesn't error (i.e. it matches)
	assert.Nil(t, bcrypt.CompareHashAndPassword(retrievedUser.Password, initialPassword))
}

func TestCreateDuplicateUser(t *testing.T) {

	var err error

	_ = createUser(mh)
	err = createUser(mh)

	// make sure error message is correct
	assert.Contains(t, err.Error(), "already exists")

	// and check that only one misshoover was created
	criteria := User{Username: "misshoover"}
	matches, _ := db.Count(&criteria)
	assert.Equal(t, matches, 1)
}

func TestFindUser(t *testing.T) {

	_ = createUser(ck)

	cookieKwan, _ := getUser("cookie.kwan")
	assert.NotZero(t, cookieKwan.ID)

	_, err := getUser("not.miss.hoover.ok")
	assert.Contains(t, err.Error(), "not found")

}

func TestAllUsers(t *testing.T) {
	users, _ := allUsers()

	var names []string

	for u := range users {
		names = append(names, users[u].Name)
	}
	assert.Equal(t, 2, len(users))

	assert.Contains(t, names, ck.Name)
	assert.Contains(t, names, mh.Name)
}
