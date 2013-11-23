package kmgFile

import (
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
