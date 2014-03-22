package kmgSsh

import (
	"bytes"
	"code.google.com/p/go.crypto/ssh"
	"fmt"
	"github.com/bronze1man/kmg/archive/kmgTar"
	"github.com/bronze1man/kmg/kmgCmd"
	"io"
	"io/ioutil"
	"os"
	"path"
)

//transfer a directory to remote host,with single tcp and tar.gz format
//Stdout will be ignore.
//Stderr will treat as an error from remote command
func CopyDirectory(client *ssh.ClientConn, selfPath string, remotePath string) (err error) {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("[CopyDirectory]Failed to create session: %s", err.Error())
	}
	defer session.Close()
	w, err := session.StdinPipe()
	if err != nil {
		return
	}
	defer w.Close()
	errBuf := &bytes.Buffer{}
	session.Stderr = errBuf
	escapedRemotePath := kmgCmd.BashEscape(remotePath)
	err = session.Start(fmt.Sprintf("mkdir -p %s;tar -xz -C %s", escapedRemotePath, escapedRemotePath))
	if err != nil {
		return fmt.Errorf("[CopyDirectory] Failed to Run: %s", err.Error())
	}
	err = kmgTar.PackageDirectoryTarGz(selfPath, w)
	if err != nil {
		return
	}
	err = w.Close()
	if err != nil {
		return
	}
	err = session.Wait()
	if err != nil {
		return fmt.Errorf("[CopyDirectory] session.Wait() err:%s", err.Error())
	}
	if errBuf.Len() != 0 {
		return fmt.Errorf("[CopyDirectory] remote: %s", string(errBuf.Bytes()))
	}
	return nil
}

//transfer a file to remote host,with single tcp and gz format
//Stdout will be ignore.
//Stderr will treat as an error from remote command
//TODO test
func CopyFile(client *ssh.ClientConn, selfPath string, remotePath string) (err error) {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("[CopyFile]Failed to create session: %s", err.Error())
	}
	defer session.Close()
	w, err := session.StdinPipe()
	if err != nil {
		return
	}
	defer w.Close()
	errReader, err := session.StderrPipe()
	if err != nil {
		return
	}
	escapedRemotePath := kmgCmd.BashEscape(remotePath)
	escapedRemoteDir := kmgCmd.BashEscape(path.Dir(remotePath))
	err = session.Start(fmt.Sprintf("mkdir -p %s;cat > %s", escapedRemoteDir, escapedRemotePath))
	if err != nil {
		return fmt.Errorf("[CopyFile] Failed to Start: %s", err.Error())
	}
	f, err := os.Open(selfPath)
	if err != nil {
		return
	}
	defer f.Close()
	//gzipW := gzip.NewWriter(w)
	//defer gzipW.Close()
	_, err = io.Copy(w, f)
	if err != nil {
		return fmt.Errorf("[CopyFile] io.Copy: %s", err)
	}
	//err = gzipW.Close()
	//if err!=nil{
	//	return
	//}
	err = w.Close()
	if err != nil {
		return
	}
	errBytes, err := ioutil.ReadAll(errReader)
	if err != nil {
		return
	}
	if len(errBytes) != 0 {
		return fmt.Errorf("[CopyFile] remote: %s", string(errBytes))
	}
	err = session.Wait()
	if err != nil {
		return fmt.Errorf("[CopyFile] session.Wait() err:%s", err.Error())
	}
	return nil
}
