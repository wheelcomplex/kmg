package kmgSsh

import (
	"code.google.com/p/go.crypto/ssh"
	"fmt"
)

type typePassword string

func (pass typePassword) Password(user string) (password string, err error) {
	return string(pass), nil
}
func DialWithPassword(addr string, username string, password string) (client *ssh.ClientConn, err error) {
	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.ClientAuth{
			ssh.ClientAuthPassword(typePassword(password)),
		},
	}
	client, err = ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("[DialWithPassword] Failed to dial: %s", err.Error())
	}
	return
}
