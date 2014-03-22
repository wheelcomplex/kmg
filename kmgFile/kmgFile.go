package kmgFile

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func IsDotFile(path string) bool {
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
