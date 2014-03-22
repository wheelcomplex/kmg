package kmgSsh

import (
	"code.google.com/p/go.crypto/ssh"
	"fmt"
	"io"
)

type keychain struct {
	key ssh.Signer
}

func (k *keychain) Key(i int) (ssh.PublicKey, error) {
	if i != 0 {
		return nil, nil
	}
	return k.key.PublicKey(), nil
}

func (k *keychain) Sign(i int, rand io.Reader, data []byte) (sig []byte, err error) {
	if i != 0 {
		return nil, fmt.Errorf("[keychain.Sign]invalid key index:%d", i)
	}
	return k.key.Sign(rand, data)
}

//DialWithCertFile
//addr like 127.0.0.1:22
//username is the username on that remote machine
//clientKey is the content of ~/.ssh/id_rsa
func DialWithCertFile(addr string, username string, clientKeyBytes []byte) (client *ssh.ClientConn, err error) {
	signer, err := ssh.ParsePrivateKey(clientKeyBytes)
	if err != nil {
		return
	}
	clientKey := &keychain{signer}
	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.ClientAuth{
			ssh.ClientAuthKeyring(clientKey),
		},
	}
	client, err = ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("[DialWithCertFile] Failed to dial: %s", err.Error())
	}
	return
}
