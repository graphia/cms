package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	mh = User{
		//ID:       1,
		Username: "misshoover",
		Email:    "e.hoover@springfield.k12.us",
		Name:     "Elizabeth Hoover",
		Password: "SuperSecret123",
		Active:   true,
	}

	ck = User{
		//ID:       2,
		Username: "cookie.kwan",
		Email:    "cookie@springfield-realtors.com",
		Name:     "Cookie Kwan",
		Password: "P@ssword!",
		Active:   true,
	}

	ds = User{
		ID:       3,
		Username: "dolph.starbeam",
		Email:    "dolph.starbeam@springfield.k12.us",
		Name:     "Dolph Starbeam",
		Password: "mightypig",
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

	// initialPassword as a byte slice
	ipwBytes := []byte(mh.Password)

	_ = createUser(mh)

	retrievedUser := User{}
	_ = db.One("Username", "misshoover", &retrievedUser)

	// retrieved password as a byteslice
	rpwBytes := []byte(retrievedUser.Password)

	// ensure that the value written to the db has been modified
	assert.NotEqual(t, rpwBytes, ipwBytes)

	// ensure comparing our generated hash with the initial password doesn't error (i.e. it matches)
	assert.Nil(t, bcrypt.CompareHashAndPassword(rpwBytes, ipwBytes))
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

	lu := ds.limitedUser()

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

func Test_setPublicKey(t *testing.T) {

	_ = createUser(ds)
	certsPath := "../tests/backend/certificates"

	validPub, _ := ioutil.ReadFile(filepath.Join(certsPath, "valid.pub"))
	validPubFields := strings.Fields(string(validPub))

	tooShortPub, _ := ioutil.ReadFile(filepath.Join(certsPath, "too_short.pub"))
	invalidPub, _ := ioutil.ReadFile(filepath.Join(certsPath, "invalid.pub"))

	type args struct {
		user User
		key  string
	}
	tests := []struct {
		name    string
		args    args
		want    publicKey
		wantErr bool
		errMsg  string
	}{
		{
			name: "Set valid public key",
			args: args{
				user: ds,
				key:  string(validPub),
			},
			want: publicKey{
				Type:    validPubFields[0],
				Raw:     []byte(validPubFields[1]),
				Comment: validPubFields[2],
			},
		},
		{
			name: "Too short public key",
			args: args{
				user: ds,
				key:  string(tooShortPub),
			},
			wantErr: true,
			errMsg:  "invalid key",
		},
		{
			name: "Invalid key",
			args: args{
				user: ds,
				key:  string(invalidPub),
			},
			wantErr: true,
			errMsg:  "invalid key",
		},
	}
	for _, tt := range tests {
		err := setPublicKey(tt.args.user, tt.args.key)

		if tt.wantErr && (err == nil) {
			t.Fatal("Error expected, none found")
		}

		if tt.wantErr {
			assert.Contains(t, err.Error(), tt.errMsg)
			continue
		}

		// should be no errors
		assert.Nil(t, err)

		// Make sure record saved correctly
		user, err := getUserByUsername("dolph.starbeam")
		assert.NotNil(t, user.PublicKey)
		assert.Equal(t, tt.want, user.PublicKey)
	}
}
