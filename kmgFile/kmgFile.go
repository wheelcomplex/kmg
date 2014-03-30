package kmgFile

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func IsDotFile(path string) bool {
	if path == "./" {
		return false
	}
	base := filepath.Base(path)
	if strings.HasPrefix(base, ".") {
		return true
	}
	return false
}

func GetFileBaseWithoutExt(p string) string {
	return filepath.Base(p[:len(p)-len(filepath.Ext(p))])
}

func WriteFile(path string, content []byte) (err error) {
	return ioutil.WriteFile(path, content, os.FileMode(0777))
}

func FileExist(path string) (exist bool, err error) {
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, err
}
