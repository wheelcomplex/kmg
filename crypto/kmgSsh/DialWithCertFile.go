package kmgSsh

import (
	"code.google.com/p/go.crypto/ssh"
	"fmt"
)

//DialWithCertFile
//addr like 127.0.0.1:22
//username is the username on that remote machine
//clientKey is the content of ~/.ssh/id_rsa
func DialWithCertFile(addr string, username string, clientKeyBytes []byte) (client *ssh.Client, err error) {
	signer, err := ssh.ParsePrivateKey(clientKeyBytes)
	if err != nil {
		return
	}
	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}
	client, err = ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("[DialWithCertFile] Failed to dial: %s", err.Error())
	}
	return
}
