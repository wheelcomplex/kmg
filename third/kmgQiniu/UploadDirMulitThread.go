package kmgQiniu

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgTask"
	"github.com/qiniu/api/rs"
	"os"
	"path/filepath"
)

type uploadFileRequest struct {
	localPath  string
	remotePath string
	expectHash string
}

//多线程上传目录
//1.某个文件仅在一个线程中上传,
//2.检查同名和同内容的文件是否已经存在了,如果存在,且hash相同便不上传(断点续传)
func UploadDirMulitThread(ctx *Context, localRoot string, remoteRoot string) (err error) {
	tm := kmgTask.NewLimitThreadErrorHandleTaskManager(5, 3)
	defer tm.Close()
	requestList := []uploadFileRequest{}
	//dispatch task 分配任务
	err = filepath.Walk(localRoot, func(path string, info os.FileInfo, inErr error) (err error) {
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
		expectHash, err := ComputeHashFromFile(path)
		if err != nil {
			return
		}
		requestList = append(requestList, uploadFileRequest{
			localPath:  path,
			remotePath: NormalizeRemotePath(remotePath),
			expectHash: expectHash,
		})
		return
	})
	//群发状态询问消息减少网络连接数量,加快速度
	entryPathList := make([]rs.EntryPath, len(requestList))
	for i, req := range requestList {
		entryPathList[i].Bucket = ctx.bucket
		entryPathList[i].Key = req.remotePath
	}
	batchRet, err := ctx.client.BatchStat(nil, entryPathList)
	if err != nil {
		return
	}
	for i, ret := range batchRet {
		//验证hash
		if ret.Error != "" && ret.Error != "no such file or directory" {
			return fmt.Errorf("[UploadDirMulitThread] [remotePath:%s]ctx.client.BatchStat err[%s]",
				requestList[i].remotePath, ret.Error)
		}
		if ret.Data.Hash == requestList[i].expectHash {
			continue
		}
		tm.AddTask(func() (err error) {
			return UploadFileWithHash(ctx, requestList[i].localPath, requestList[i].remotePath, requestList[i].expectHash)
		})
	}
	tm.Wait()
	if err != nil {
		return err
	}
	err = tm.GetError()
	return
}
