package kmgQiniu

import (
	"strings"
)

//正规化传给七牛的远程路径
//解决windows目录分隔符和开头的"/"的问题
func NormalizeRemotePath(path string) string {
	return strings.TrimLeft(strings.Replace(path, "\\", "/", -1), "/")
}
