package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/gliderlabs/ssh"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	gossh "golang.org/x/crypto/ssh"
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

	sb = User{
		ID:       4,
		Username: "selma.bouvier",
		Email:    "selma.bouvier@macguyver-fans.org",
		Name:     "Selma Bouvier",
		Password: "ilubjubjub",
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

func TestUser_addPublicKey(t *testing.T) {

	db.Drop("User")
	db.Drop("PublicKey")

	_ = createUser(ds)
	certsPath := "../tests/backend/certificates"

	validPub, _ := ioutil.ReadFile(filepath.Join(certsPath, "valid.pub"))
	validPubParsed, _, _, _, _ := ssh.ParseAuthorizedKey(validPub)
	tooShortPub, _ := ioutil.ReadFile(filepath.Join(certsPath, "too_short.pub"))
	invalidPub, _ := ioutil.ReadFile(filepath.Join(certsPath, "invalid.pub"))

	type args struct {
		raw string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		errMsg  string
		want    PublicKey
	}{
		{
			name: "Set valid public key",
			args: args{
				raw: string(validPub),
			},
			want: PublicKey{
				Raw:         validPubParsed.Marshal(),
				Fingerprint: gossh.FingerprintSHA256(validPubParsed),
			},
		},
		{
			name: "Too short public key",
			args: args{
				raw: string(tooShortPub),
			},
			wantErr: true,
			errMsg:  "invalid key",
		},
		{
			name: "Invalid key",
			args: args{
				raw: string(invalidPub),
			},
			wantErr: true,
			errMsg:  "invalid key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			user, _ := getUserByUsername(ds.Username)
			err := user.addPublicKey(tt.args.raw)

			if tt.wantErr && (err == nil) {
				t.Fatal("Error expected, none found")
			}

			if tt.wantErr {
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			// should be no errors
			assert.Nil(t, err)

			// one public key should have been created for this user
			var matchingKeys []PublicKey
			db.Find("UserID", user.ID, &matchingKeys)
			assert.Equal(t, 1, len(matchingKeys))

			// and the public key should have the right attributes
			var pk PublicKey
			db.One("UserID", user.ID, &pk)
			assert.Equal(t, tt.want.Fingerprint, pk.Fingerprint)
			assert.Equal(t, tt.want.Raw, pk.Raw)

		})
	}
}

func Test_getUserByFingerprint(t *testing.T) {

	db.Drop("User")
	db.Drop("PublicKey")

	certsPath := "../tests/backend/certificates"

	validPub, _ := ioutil.ReadFile(filepath.Join(certsPath, "valid.pub"))
	parsedValidPub, _, _, _, _ := gossh.ParseAuthorizedKey(validPub)

	type args struct {
		raw         string
		user        User
		fingerprint string
		key         gossh.PublicKey
	}
	tests := []struct {
		name         string
		args         args
		wantUsername string
		wantErr      bool
		errMsg       string
		generateCert bool
		generateUser bool
	}{
		{
			name: "Successful match",
			args: args{
				raw:         string(validPub),
				user:        mh,
				fingerprint: gossh.FingerprintSHA256(parsedValidPub),
				key:         parsedValidPub,
			},
			wantUsername: mh.Username,
			wantErr:      false,
			generateCert: true,
			generateUser: true,
		},
		{
			name: "Key not found by fingerprint",
			args: args{
				raw:         string(validPub),
				user:        mh,
				fingerprint: gossh.FingerprintSHA256(parsedValidPub),
				key:         parsedValidPub,
			},
			wantUsername: ds.Username,
			wantErr:      true,
			errMsg: fmt.Sprintf(
				"not found, fingerprint: %s",
				gossh.FingerprintSHA256(parsedValidPub),
			),
			generateCert: false,
			generateUser: true,
		},
		{
			name: "Key found but no matching user",
			args: args{
				raw:         string(validPub),
				user:        mh,
				fingerprint: gossh.FingerprintSHA256(parsedValidPub),
				key:         parsedValidPub,
			},
			wantUsername: ds.Username,
			wantErr:      true,
			errMsg:       "no user found for public key 1", // it *should* be the only one
			generateCert: true,
			generateUser: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db.Drop("User")
			db.Drop("PublicKey")

			_ = createUser(mh)

			var user User
			if tt.generateUser {
				user, _ = getUserByUsername(tt.wantUsername)
			}

			if tt.generateCert {
				pk := PublicKey{
					UserID:      user.ID,
					Raw:         tt.args.key.Marshal(),
					Fingerprint: tt.args.fingerprint,
				}

				db.Save(&pk)
			}

			gotUser, err := getUserByFingerprint(tt.args.key)

			if tt.wantErr {
				assert.Equal(t, err.Error(), tt.errMsg)
				return
			}

			assert.Equal(t, tt.wantUsername, gotUser.Username)

		})
	}
}
