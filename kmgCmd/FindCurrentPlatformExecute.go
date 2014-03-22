package kmgCmd

import (
	"errors"
	"os/exec"
	"runtime"
)

var ErrNotFound = errors.New("executable file not found.")

func FindCurrentPlatformExecute(prefix string) (path string, err error) {
	tryPath := []string{
		prefix,
	}
	//arch 386 may be compatible with amd64
	//TODO find real system arch
	var archList []string
	switch runtime.GOARCH {
	case `386`:
		archList = []string{`386`, `amd64`}
	case `amd64`:
		archList = []string{`amd64`, `386`}
	default:
		archList = []string{runtime.GOARCH}
	}
	switch runtime.GOOS {
	case `windows`:
		for _, arch := range archList {
			tryPath = append(tryPath, prefix+"_"+runtime.GOOS+"_"+arch+".exe")
		}
	default:
		for _, arch := range archList {
			tryPath = append(tryPath, prefix+"_"+runtime.GOOS+"_"+arch)
		}
	}
	for _, p := range tryPath {
		path, err = exec.LookPath(p)
		if err != nil {
			if err == exec.ErrNotFound {
				continue
			} else {
				return "", err
			}
		}
		return path, err
	}
	return "", ErrNotFound
}
