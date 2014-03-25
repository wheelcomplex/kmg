package kmgQiniu

import (
	"fmt"
	qiniuIo "github.com/qiniu/api/io"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/rpc"
)

//上传文件,检查同名和同内容文件
//先找cdn上是不是已经有一样的文件了,以便分文件断点续传,再上传
func UploadFileCheckExist(ctx *Context, localPath string, remotePath string) (err error) {
	remotePath = NormalizeRemotePath(remotePath)
	entry, err := ctx.client.Stat(nil, ctx.bucket, remotePath)
	if err != nil {
		if !(err.(*rpc.ErrorInfo) != nil && err.(*rpc.ErrorInfo).Err == "no such file or directory") {
			return err
		}
	}
	expectHash, err := ComputeHashFromFile(localPath)
	if err != nil {
		return
	}
	//already have a file with same context and same key,do nothing
	if entry.Hash == expectHash {
		return
	}
	return UploadFileWithHash(ctx, localPath, remotePath, expectHash)
}

//上传文件,检查返回的hash和需要的hash是否一致
func UploadFileWithHash(ctx *Context, localPath string, remotePath string, expectHash string) (err error) {
	var ret qiniuIo.PutRet
	var extra = &qiniuIo.PutExtra{
		CheckCrc: 1,
	}
	putPolicy := rs.PutPolicy{
		Scope: ctx.bucket + ":" + remotePath,
	}
	uptoken := putPolicy.Token(nil)
	err = qiniuIo.PutFile(nil, &ret, uptoken, remotePath, localPath, extra)
	if err != nil {
		return
	}
	if ret.Hash != expectHash {
		return fmt.Errorf("[UploadFileWithHash][remotePath:%s] ret.Hash:[%s]!=expectHash[%s] ", remotePath, ret.Hash, expectHash)
	}

	return
}
