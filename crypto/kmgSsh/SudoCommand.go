package kmgSsh

import (
	"bytes"
	"code.google.com/p/go.crypto/ssh"
	"fmt"
	"github.com/howeyc/gopass"
	"os"
)

//run sudo a command in remote and ask password in current console,
//Stdout will goto current console too.
//Stderr will treat as an error from remote command
func SudoCommand(client *ssh.Client, cmd string) (err error) {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("[SudoCommand] Failed to create session: %s", err.Error())
	}
	defer session.Close()
	w, err := session.StdinPipe()
	if err != nil {
		return
	}
	session.Stdout = os.Stdout
	errBuf := &bytes.Buffer{}
	session.Stderr = errBuf
	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return fmt.Errorf("[SudoCommand] request for pseudo terminal failed: %s", err)
	}
	err = session.Start(cmd)
	if err != nil {
		return fmt.Errorf("[SudoCommand] Failed to Start: %s", err.Error())
	}
	pass := gopass.GetPasswd()
	w.Write(append(pass, byte('\n')))
	err = session.Wait()
	if err != nil {
		return fmt.Errorf("[SudoCommand] session.Wait err: %s", err.Error())
	}
	if errBuf.Len() != 0 {
		return fmt.Errorf("[SudoCommand] remote err: %s", err.Error())
	}
	return nil
}
