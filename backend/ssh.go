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
	Raw         []byte
	Fingerprint string
}

// User returns the Public Key's assoicated User
func (pk PublicKey) User() (user User, err error) {
	err = db.One("ID", pk.UserID, &user)
	if err != nil && pk.UserID != 0 {
		return user, fmt.Errorf("Could not find user %d", pk.UserID)
	}
	return user, err
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

		err = db.One("Fingerprint", fp, &retrievedKeyRecord)
		if err != nil {
			// do something, key probably not found
			Warning.Println("No key found", fp)
			return false
		}

		retrievedKey, err := gossh.ParsePublicKey(retrievedKeyRecord.Raw)
		if err != nil {
			Warning.Println("Key could not be parsed", retrievedKeyRecord)
			return false
		}

		if !ssh.KeysEqual(providedKey, retrievedKey) {
			Warning.Println("Keys don't match")
			return false
		}

		return true
	})

	go ssh.ListenAndServe(fmt.Sprintf(":%s", config.SSHListenPort), nil, publicKeyOption)

}
