package kmgSsh

import (
	"code.google.com/p/go.crypto/ssh"
	"fmt"
	"github.com/howeyc/gopass"
	"io/ioutil"
	"os"
)

type consoleAskPassword struct {
	user string
	addr string
}

func (p consoleAskPassword) Password() (password string, err error) {
	fmt.Printf("[ssh] password for %s@%s", p.user, p.addr)
	password = string(gopass.GetPasswd())
	return password, nil
}

//TODO 某种认证方法只有一个会被使用,需要多次猜测
func DialInConsole(addr string, username string) (client *ssh.Client, err error) {
	//find cert file
	pathList := certFilePathList()
	authList := []ssh.AuthMethod{}
	for _, path := range pathList {
		clientKeyBytes, err := ioutil.ReadFile(path)
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("[DialInConsole] ioutil.ReadFile() err:%s", err)
			}
		} else {
			signer, err := ssh.ParsePrivateKey(clientKeyBytes)
			if err != nil {
				return nil, fmt.Errorf("[DialInConsole] ssh.ParsePrivateKey err:%s", err)
			}
			//clientKey := &keychain{signer}
			authList = append(authList, ssh.PublicKeys(signer))
		}
	}
	authList = append(authList, ssh.PasswordCallback(func() (secret string, err error) {
		fmt.Printf("[ssh] password for %s@%s", username, addr)
		secret = string(gopass.GetPasswd())
		return
	}))
	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: authList,
	}
	client, err = ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("[DialInConsole] Failed to dial: %s", err.Error())
	}
	return
}
