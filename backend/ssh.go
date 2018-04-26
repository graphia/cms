package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

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

// append the git executable path
func checkGitExecutable(name string) (ap string, err error) {
	if name != "git-upload-pack" && name != "git-receive-pack" {
		return ap, fmt.Errorf("Error: Only Git operations are permitted")
	}
	ap = filepath.Join(config.GitBinPath, name)
	return ap, err
}

func resolvePath(name string) (ap string, err error) {

	switch name {
	case "theme":
		return config.HugoThemePath, err
	case "content":
		return config.Repository, err
	default:
		return ap, fmt.Errorf("invalid repo")
	}

}

func setupSSH() {

	ssh.Handle(func(s ssh.Session) {

		Debug.Println("Incoming SSH connection")
		var msg string

		if s.User() != config.GitUser {
			io.WriteString(s, "Access denied\n")
			return
		}

		if len(s.Command()) == 0 {
			io.WriteString(s, "Graphia: Connection successful\n")
			return
		}

		// Only git-upload-pack and git-receive-pack are valid
		// get the operation (either git-upload-pack or git-receive-pack)
		// and the repo. the actual command is in the format:
		// git-upload-pack repo-name
		operation, err := checkGitExecutable(s.Command()[0])
		if err != nil {
			Error.Println("Invalid operation", err)
			io.WriteString(s, err.Error())
			return
		}

		rp, err := resolvePath(s.Command()[1])
		if err != nil {
			msg = "Invalid command"
			Error.Println(msg, s.Command()[1])
			io.WriteString(s, msg)
			return
		}

		// Make sure target repo exists at the given path
		Debug.Println("checking repo exists", rp)
		_, err = os.Stat(rp)
		if os.IsNotExist(err) {
			Error.Println("repo does not exist, abort", rp)
			msg = "repo does not exist at specified path"
			io.WriteString(s, msg)
			return
		}

		Debug.Println("executing", operation, rp)
		cmd := exec.Command(operation, rp)
		cmd.Env = append(os.Environ(), "SSH_ORIGINAL_COMMAND="+operation)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			Info.Printf("SSH: StdoutPipe: %v", err)
			return
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			Error.Printf("SSH: StderrPipe: %v", err)
			return
		}

		input, err := cmd.StdinPipe()
		if err != nil {
			Info.Printf("SSH: StdinPipe: %v", err)
			return
		}

		err = cmd.Start()
		if err != nil {
			fmt.Printf("SSH: Start: %v", err)
			return
		}

		go io.Copy(input, s)
		io.Copy(s, stdout)
		io.Copy(s.Stderr(), stderr)

		err = cmd.Wait()
		if err != nil {
			Error.Printf("SSH: Wait: %v", err)
			return
		}

		s.SendRequest("exit-status", false, []byte{0, 0, 0, 0})
		return

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
