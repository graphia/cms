package main

import (
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/bcrypt"
	gossh "golang.org/x/crypto/ssh"
)

var (
	// ErrUserNotExists returned when the user can't be found
	ErrUserNotExists = errors.New("user does not exist")
)

// UserCredentials is the subset of User required for auth
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LimitedUser is a 'safe' subset of user data that we can
// send out via the API. Password is omitted
type LimitedUser struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Admin           bool   `json:"admin"`
	Active          bool   `json:"active"`
	ConfirmationKey string `json:"confirmation_key"`
}

// User holds all information specific to a user
type User struct {
	ID              int    `json:"id" storm:"id,increment"`
	Name            string `json:"name" validate:"required,min=3,max=64"`
	Username        string `json:"username" storm:"index,unique" validate:"required,min=3,max=32"`
	Password        string `json:"password" validate:"required,min=6"`
	Email           string `json:"email" storm:"unique" validate:"required,email"`
	Active          bool   `json:"active"`
	Admin           bool   `json:"admin"`
	TokenString     string `json:"token_string" storm:"unique"`
	ConfirmationKey string `json:"confirmation_key" storm:"unique"`
}

// PasswordUpdate used by users to modify their password
type PasswordUpdate struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}

func (u User) addPublicKey(name, raw string) error {

	// make sure the key is valid and populate fingerprint
	parsed, _, _, _, err := ssh.ParseAuthorizedKey([]byte(raw))
	if err != nil {
		return fmt.Errorf("invalid key")
	}

	fp := gossh.FingerprintSHA256(parsed)

	pk := PublicKey{
		UserID:      u.ID,
		Name:        name,
		Raw:         parsed.Marshal(),
		Fingerprint: fp,
	}

	return db.Save(&pk)

}

func (u User) keys() (pks []PublicKey, err error) {

	err = db.Find("UserID", u.ID, &pks)

	// if no matching records are found, that's ok, just
	// return the empty slice and ignore the error
	if err != nil && err.Error() == "not found" {
		Warning.Println("No SSH keys found for", u.Username)
		return pks, nil
	}

	if err != nil {
		Warning.Println("User key query failed for", u.Username)
		return pks, err
	}

	return pks, err
}

func (u User) limitedUser() LimitedUser {
	return LimitedUser{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Name:     u.Name,
		Admin:    u.Admin,
		Active:   u.Active,
		// ConfirmationKey: u.ConfirmationKey, (risky to leak this, let's omit)
	}
}

func (u User) reload() (User, error) {
	user, err := getUserByID(u.ID)
	if err != nil {
		return u, err
	}
	return user, err
}

func (u User) setToken(tokenString string) error {
	return db.UpdateField(&u, "TokenString", tokenString)
}

func (u User) unsetToken() error {
	return db.UpdateField(&u, "TokenString", "")
}

func (u User) setRandomConfirmationKey() error {
	nk := generateRandomConfirmationKey()
	return db.UpdateField(&u, "ConfirmationKey", nk)
}

func (u User) delete() error {
	return db.DeleteStruct(&u)
}

func (u User) setPassword(pw string) error {

	Debug.Println("password valid, saving")

	bcryptedPassword, err := generateBcryptedPassword(pw)
	if err != nil {
		return err
	}

	return db.UpdateField(&u, "Password", bcryptedPassword)
}

func generateBcryptedPassword(pw string) (cpw string, err error) {
	var bcpw []byte
	if !(len(pw) >= 6) {
		return cpw, fmt.Errorf("password must be at least 6 characters")
	}

	bcpw, err = bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	cpw = string(bcpw)
	return cpw, err

}

func (u User) checkPassword(pw string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw))

	if err != nil {
		return fmt.Errorf("passwords don't match")
	}
	return nil
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
		return user, ErrUserNotExists
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

func getLimitedUserByConfirmationKey(ck string) (limitedUser LimitedUser, err error) {
	var user User
	err = db.One("ConfirmationKey", ck, &user)

	if user.ID == 0 {
		Warning.Println("Cannot find user with ConfirmationKey", ck)
		return limitedUser, ErrUserNotExists
	}

	return user.limitedUser(), err
}

func getUserByConfirmationKey(ck string) (user User, err error) {

	err = db.One("ConfirmationKey", ck, &user)

	if user.ID == 0 {
		Warning.Println("Cannot find user with ConfirmationKey", ck)
		return user, ErrUserNotExists
	}

	return user, err
}

func getUserByFingerprint(pk gossh.PublicKey) (user User, err error) {

	var pub PublicKey

	fp := gossh.FingerprintSHA256(pk)
	Debug.Println("searching with fingerprint", fp)

	err = db.One("Fingerprint", fp, &pub)
	if err != nil {
		return user, fmt.Errorf("not found, fingerprint: %s", fp)
	}

	user, err = pub.User()

	if err != nil {
		return user, fmt.Errorf("no user found for public key %d", pub.ID)
	}
	return user, err
}

func createUser(user User) (err error) {

	bcryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(bcryptedPassword)
	user.ConfirmationKey = generateRandomConfirmationKey()

	err = validate.Struct(user)
	if err != nil {
		return err
	}

	Info.Println("Validation passed, saving", user)

	go sendEmailConfirmation(user)

	err = db.Save(&user)
	return err
}

func updateUser(user User, updates LimitedUser) (err error) {

	// note username isn't updateable yet
	user.Name = updates.Name
	user.Email = updates.Email
	user.Active = updates.Active
	user.Admin = updates.Admin

	if user.Username != updates.Username {
		Warning.Printf("attempt to change username from %s to %s denied", user.Username, updates.Username)
	}

	// make sure everything's above board
	err = validate.Struct(user)
	if err != nil {
		Warning.Println("validation failed for user", user, updates)
		return err
	}

	err = db.Save(&user)
	return err
}

func activateUser(user User, password string) (err error) {

	if user.Active {
		return fmt.Errorf("%s already active", user.Username)
	}

	user.Password, err = generateBcryptedPassword(password)
	if err != nil {
		return
	}

	user.Active = true

	return db.Save(&user)

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

func generateRandomConfirmationKey() string {

	const length = 32
	var chars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	clen := len(chars)

	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			panic("uniuri: error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}

}
