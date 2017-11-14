package main

import (
	"fmt"

	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/bcrypt"
	gossh "golang.org/x/crypto/ssh"
)

// PublicKey holds a User's Public Key
type PublicKey struct {
	ID          int `storm:"id,increment"`
	UserID      int
	Raw         []byte
	Fingerprint string
}

// User returns the Public Key's assoicated User
func (pk PublicKey) User() (user User, err error) {
	err = db.One("ID", pk.UserID, &user)
	if err != nil {
		return user, fmt.Errorf("Could not find user %d", pk.UserID)
	}
	return user, err
}

// UserCredentials is the subset of User required for auth
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LimitedUser is a 'safe' subset of user data that we can
// send out via the API. Password is omitted
type LimitedUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// User holds all information specific to a user
type User struct {
	ID          int    `json:"id" storm:"id,increment"`
	Name        string `json:"name" validate:"required,min=3,max=64"`
	Username    string `json:"username" storm:"unique" validate:"required,min=3,max=32"`
	Password    string `json:"password" validate:"required,min=6"`
	Email       string `json:"email" storm:"unique" validate:"email,required"`
	Active      bool   `json:"active"`
	TokenString string `json:"token_string" storm:"unique"`
}

func (u User) addPublicKey(raw string) error {
	//fields := strings.Fields(raw)

	// make sure the key is valid and populate fingerprint
	parsed, _, _, _, err := ssh.ParseAuthorizedKey([]byte(raw))
	if err != nil {
		return fmt.Errorf("invalid key")
	}

	fp := gossh.FingerprintSHA256(parsed)

	pk := PublicKey{
		UserID:      u.ID,
		Raw:         parsed.Marshal(),
		Fingerprint: fp,
	}

	return db.Save(&pk)

}

func (u User) limitedUser() LimitedUser {
	return LimitedUser{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Name:     u.Name,
	}
}

func getUserByID(id int) (user User, err error) {
	err = db.One("ID", id, &user)

	if user.ID == 0 {
		Warning.Println("Cannot find user with ID", id)
		return user, fmt.Errorf("User not found: %d", id)
	}

	return user, err
}

func getUserByUsername(username string) (user User, err error) {
	err = db.One("Username", username, &user)

	if user.ID == 0 {
		Warning.Println("Cannot find user with Username", username)
		return user, fmt.Errorf("not found username: %s", username)
	}

	Debug.Println("Found user", username)

	return user, err
}

func getUserByEmail(email string) (user User, err error) {
	err = db.One("Username", email, &user)

	if user.ID == 0 {
		Warning.Println("Cannot find user with email address", email)
		return user, fmt.Errorf("not found email: %s", email)
	}

	Debug.Println("Found user", email)

	return user, err
}

func getLimitedUserByUsername(username string) (limitedUser LimitedUser, err error) {
	var user User
	err = db.One("Username", username, &user)

	if user.ID == 0 {
		Warning.Println("Cannot find user with Username", username)
		return limitedUser, fmt.Errorf("not found: %s", username)
	}

	return user.limitedUser(), err
}

func getUserByFingerprint(pk gossh.PublicKey) (user User, err error) {

	fp := gossh.FingerprintSHA256(pk)
	Debug.Println("searching with fingerprint", fp)
	err = db.One("Fingerprint", fp, &user)

	if user.ID == 0 {
		Warning.Println("Cannot find user with ssh public key fingerprint", fp)
		return user, fmt.Errorf("not found ssh public key fingerprint: %s", fp)
	}

	Debug.Println("Found user by fingerprint", fp)

	return user, err
}

func createUser(user User) (err error) {

	bcryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(bcryptedPassword)

	err = validate.Struct(user)
	if err != nil {
		return err
	}

	Info.Println("Validation passed, saving", user)

	err = db.Save(&user)
	if err != nil {
		return fmt.Errorf("User cannot be created, %v", err)
	}
	return nil
}

func allUsers() (limitedUsers []LimitedUser, err error) {

	var users []User

	err = db.All(&users)

	for _, user := range users {
		limitedUsers = append(limitedUsers, user.limitedUser())
	}

	if err != nil {
		return nil, err
	}
	return limitedUsers, nil
}

func countUsers() (qty int, err error) {
	qty, err = db.Count(&User{})
	if err != nil {
		return -1, err
	}
	return qty, err
}

// Actually do a 'hard' delete
func deleteUser(user User) bool {
	err := db.DeleteStruct(&user)
	if err != nil {
		// handle
	}
	return true
}

// Soft delete, prevent user from signing in and
// accessing the API
func deactivateUser(user User) error {
	err := db.UpdateField(&user, "Active", false)
	if err != nil {
		return err
	}
	return nil
}

func reactivateUser(user User) error {
	return db.UpdateField(&user, "Active", true)

}

func setToken(user User, tokenString string) error {
	return db.UpdateField(&user, "TokenString", tokenString)
}
