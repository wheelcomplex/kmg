package kmgQiniu

import (
	"fmt"
	"os"
	"path/filepath"
)

//单线程上传目录
//检查同文件名和同内容的文件是否存在,如果内容相同便不上传
func UploadDir(ctx *Context, localRoot string, remoteRoot string) (err error) {
	return filepath.Walk(localRoot, func(path string, info os.FileInfo, inErr error) (err error) {
		if inErr != nil {
			return inErr
		}
		if info.IsDir() {
			return
		}
		relPath, err := filepath.Rel(localRoot, path)
		if err != nil {
			return fmt.Errorf("[qiniuUploadDir] filepath.Rel err:%s", err.Error())
		}
		remotePath := filepath.Join(remoteRoot, relPath)
		err = UploadFileCheckExist(ctx, path, remotePath)
		if err != nil {
			return fmt.Errorf("[qiniuUploadDir] qiniuUploadFile err:%s", err.Error())
		}
		return
	})
}
