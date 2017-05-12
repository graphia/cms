package main

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/assert"
)

var (
	mh = User{
		//ID:       1,
		Username: "misshoover",
		Email:    "e.hoover@springfield.k12.us",
		Name:     "Elizabeth Hoover",
		Password: []byte("SuperSecret123"),
		Active:   true,
	}

	ck = User{
		//ID:       2,
		Username: "cookie.kwan",
		Email:    "cookie@springfield-realtors.com",
		Name:     "Cookie Kwan",
		Password: []byte("P@ssword!"),
		Active:   true,
	}

	ds = User{
		ID:       3,
		Username: "dolph.starbeam",
		Email:    "dolph.starbeam@springfield.k12.us",
		Name:     "Dolph Starbeam",
		Password: []byte("mightypig"),
		Active:   false,
	}
)

func init() {
	Debug.Println("Flushing test database")
	db = flushDB(config.Database)
}

func TestCreateUser(t *testing.T) {
	db.Drop("User")

	_ = createUser(mh)

	retrievedUser := User{}
	_ = db.One("Username", "misshoover", &retrievedUser)

	assert.Equal(t, mh.Name, retrievedUser.Name)
	assert.Equal(t, mh.Username, retrievedUser.Username)
	assert.Equal(t, mh.Email, retrievedUser.Email)

}

func TestPasswordIsEncrypted(t *testing.T) {
	db.Drop("User")

	initialPassword := mh.Password
	_ = createUser(mh)

	retrievedUser := User{}
	_ = db.One("Username", "misshoover", &retrievedUser)

	// ensure that the value written to the db has been modified
	assert.NotEqual(t, initialPassword, retrievedUser.Password)

	// ensure comparing our generated hash with the initial password doesn't error (i.e. it matches)
	assert.Nil(t, bcrypt.CompareHashAndPassword(retrievedUser.Password, initialPassword))
}

func TestCreateDuplicateUsers(t *testing.T) {
	db.Drop("User")

	var err error

	_ = createUser(mh)
	err = createUser(mh)

	// make sure error message is correct
	assert.Contains(t, err.Error(), "already exists")

	// and check that only one misshoover was created
	criteria := User{Username: "miss.hoover"}
	matches, _ := db.Count(&criteria)
	assert.Equal(t, matches, 1)
}

func TestGetUserByID(t *testing.T) {
	db.Drop("User")

	_ = createUser(ck)

	cookieKwan, _ := getUserByID(1)
	assert.NotZero(t, cookieKwan.ID)

	_, err := getUserByID(99)
	assert.Contains(t, err.Error(), "not found")

}

func TestGetUserByUsername(t *testing.T) {
	db.Drop("User")

	_ = createUser(ck)

	cookieKwan, _ := getUserByUsername("cookie.kwan")
	assert.NotZero(t, cookieKwan.ID)

	_, err := getUserByUsername("not.miss.hoover.ok")
	assert.Contains(t, err.Error(), "not found")

}

func TestAllUsers(t *testing.T) {
	db.Drop("User")

	_ = createUser(ds)
	_ = createUser(ck)
	_ = createUser(mh)

	users, _ := allUsers()

	assert.Equal(t, 3, len(users))

	var names []string

	for u := range users {
		//fmt.Println(users[u].Name)
		names = append(names, users[u].Name)
	}

	assert.Contains(t, names, ck.Name)
	assert.Contains(t, names, mh.Name)
	assert.Contains(t, names, ds.Name)

}

func TestConvertToLimitedUser(t *testing.T) {
	var lu LimitedUser
	lu = convertToLimitedUser(ds)

	assert.IsType(t, LimitedUser{}, lu)
	assert.Equal(t, ds.ID, lu.ID)
	assert.Equal(t, ds.Username, lu.Username)
	assert.Equal(t, ds.Name, lu.Name)
	assert.Equal(t, ds.Email, lu.Email)
}

func TestDeleteUser(t *testing.T) {

	var users []LimitedUser

	db.Drop("User")

	// create two users
	_ = createUser(ck)
	_ = createUser(ds)

	cookieKwan, _ := getUserByUsername("cookie.kwan")

	// there should be two users
	users, _ = allUsers()
	assert.Equal(t, 2, len(users))

	// delete one of them
	_ = deleteUser(cookieKwan)

	// now there should be one left
	users, _ = allUsers()
	assert.Equal(t, 1, len(users))
}

func TestDeactivateUser(t *testing.T) {
	var user User
	var username = "cookie.kwan" // active

	_ = createUser(ck)
	user, _ = getUserByUsername(username)

	assert.True(t, user.Active)

	_ = deactivateUser(user)

	user, _ = getUserByUsername(username)

	assert.False(t, user.Active)

}

func TestReactivateUser(t *testing.T) {
	var user User
	var username = "dolph.starbeam" // inactive

	_ = createUser(ck)
	user, _ = getUserByUsername(username)

	assert.False(t, user.Active)

	_ = reactivateUser(user)

	user, _ = getUserByUsername(username)

	assert.True(t, user.Active)

}
