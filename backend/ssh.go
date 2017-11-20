package main

import (
	"fmt"
	"io"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

// PublicKey holds a User's Public Key
type PublicKey struct {
	ID          int `storm:"id,increment"`
	UserID      int
	Name        string
	Raw         []byte
	Fingerprint string `storm:"unique"`
}

// User returns the Public Key's assoicated User
func (pk PublicKey) User() (user User, err error) {
	err = db.One("ID", pk.UserID, &user)
	if err != nil && pk.UserID != 0 {
		return user, fmt.Errorf("Could not find user %d", pk.UserID)
	}
	return user, err
}

// File returns the keyfile's contents
func (pk PublicKey) File() (str string, err error) {

	key, err := gossh.ParsePublicKey(pk.Raw)
	if err != nil {
		return str, fmt.Errorf("Key could not be parsed %s", pk.Fingerprint)
	}

	return string(gossh.MarshalAuthorizedKey(key)), err
}

func setupSSH() {

	ssh.Handle(func(s ssh.Session) {

		// authorizedKey := gossh.MarshalAuthorizedKey(s.PublicKey())

		// Debug.Println("authorizedKey", authorizedKey)
		//io.WriteString(s, fmt.Sprintf("public key used by %s:\n", s.User()))

		io.WriteString(s, "Graphia: Connection successful")
	})

	publicKeyOption := ssh.PublicKeyAuth(func(ctx ssh.Context, providedKey ssh.PublicKey) bool {
		// use ssh.KeysEqual here...
		var retrievedKeyRecord PublicKey
		var err error

		fp := gossh.FingerprintSHA256(providedKey)
		Debug.Println("authenticating SSH key", fp)

		err = db.One("Fingerprint", fp, &retrievedKeyRecord)
		if err != nil {
			// do something, key probably not found
			Warning.Println("No key found", fp)
			return false
		}

		user, err := retrievedKeyRecord.User()
		if err != nil {
			Error.Println("Key not associated with user", fp)
			return false
		}
		Debug.Println("Found matching fingerprint for user", user.Username)

		retrievedKey, err := gossh.ParsePublicKey(retrievedKeyRecord.Raw)
		if err != nil {
			Warning.Println("Key could not be parsed", retrievedKeyRecord)
			return false
		}

		if ssh.KeysEqual(providedKey, retrievedKey) {
			Debug.Println("Keys match, auth enabled for", user.Username)
			return true
		}

		Warning.Println("Keys don't match", fp, user.Username)
		return false

	})

	go ssh.ListenAndServe(fmt.Sprintf(":%s", config.SSHListenPort), nil, publicKeyOption)

}
