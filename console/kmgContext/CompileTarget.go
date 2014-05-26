package kmgContext

import "strings"

// see http://golang.org/doc/install/source to get all possiable GOOS and GOARCH
// should be something like "windows_amd64","darwin_386",etc..
type CompileTarget string

func (target CompileTarget) GetGOOS() string {
	part := strings.Split(string(target), "_")
	return part[0]
}
func (target CompileTarget) GetGOARCH() string {
	part := strings.Split(string(target), "_")
	return part[1]
}
