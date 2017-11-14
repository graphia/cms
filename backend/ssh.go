package main

import (
	"fmt"
	"io"

	"github.com/gliderlabs/ssh"
	//gossh "golang.org/x/crypto/ssh"
)

func setupSSH() {

	ssh.Handle(func(s ssh.Session) {

		// authorizedKey := gossh.MarshalAuthorizedKey(s.PublicKey())

		// Debug.Println("authorizedKey", authorizedKey)
		//io.WriteString(s, fmt.Sprintf("public key used by %s:\n", s.User()))

		io.WriteString(s, "Graphia: Connection successful")
	})

	publicKeyOption := ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {

		return true
	})

	go ssh.ListenAndServe(fmt.Sprintf(":%s", config.SSHListenPort), nil, publicKeyOption)

}
