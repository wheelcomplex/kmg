package upyun

import "github.com/bronze1man/kmg/kmgTask"
import "github.com/bronze1man/kmg/kmgLog"
import "os"
import "strconv"
import "time"
import "sync"
import "path/filepath"

//批量upyun操作
type BulkUpyun struct {
	UpYun  *UpYun
	Tm     *kmgTask.LimitThreadTaskManager
	Logger kmgLog.Logger
}

//批量上传接口
//upload a file
func (obj *BulkUpyun) UploadFile(upyun_path, local_path string) {
	obj.Tm.AddTask(kmgTask.TaskFunc(func() {

		obj.Logger.Info("upload file: " + upyun_path)
		file, err := os.Open(local_path)
		if err != nil {
			obj.Logger.LogError(err)
			return
		}
		defer file.Close()
		err = obj.UpYun.WriteFile(upyun_path, file, true)
		if err != nil {
			obj.Logger.LogError(err)
			return
		}
		return
	}))
}

//upload a dir
func (obj *BulkUpyun) UploadDir(upyun_path, local_path string) {
	obj.Tm.AddTask(kmgTask.TaskFunc(func() {
		obj.Logger.Info("upload dir: " + upyun_path)
		dir, err := os.Open(local_path)
		if err != nil {
			obj.Logger.LogError(err)
			return
		}
		file_list, err := dir.Readdir(0)
		if err != nil {
			obj.Logger.LogError(err)
			return
		}
		err = obj.UpYun.MkDir(upyun_path, true)
		if err != nil {
			obj.Logger.LogError(err)
			return
		}
		for _, file_info := range file_list {
			file_name := file_info.Name()
			this_local_path := local_path + "/" + file_name
			this_upyun_path := upyun_path + "/" + file_name
			if file_info.IsDir() {
				obj.UploadDir(this_upyun_path, this_local_path)
			} else {
				obj.UploadFile(this_upyun_path, this_local_path)
			}
		}
		return
	}))
}

//download a file
func (obj *BulkUpyun) DownloadFile(upyun_path, local_path string) {
	obj.Tm.AddTask(kmgTask.TaskFunc(func() {
		obj.Logger.Info("download file: " + upyun_path)
		err := os.MkdirAll(filepath.Dir(local_path), os.FileMode(0777))
		if err != nil {
			obj.Logger.LogError(err)
			return
		}

		file, err := os.Create(local_path)
		if err != nil {
			obj.Logger.LogError(err)
			return
		}
		defer file.Close()
		err = obj.UpYun.ReadFile(upyun_path, file)
		if err != nil {
			obj.Logger.LogError(err)
			return
		}
		return
	}))
}

//resursive download a dir
func (obj *BulkUpyun) DownloadDir(upyun_path string, file_path string) {
	obj.Tm.AddTask(kmgTask.TaskFunc(func() {
		obj.Logger.Info("download dir: " + upyun_path)
		file_list, err := obj.UpYun.ReadDir(upyun_path)
		file_mode := os.FileMode(0777)
		if err != nil {
			obj.Logger.LogError(err)
			return
		}
		for _, file_info := range file_list {
			file_type := file_info.Type
			file_name := file_info.Name
			this_local_path := file_path + "/" + file_name
			this_upyun_path := upyun_path + "/" + file_name
			if file_type == "folder" {
				err := os.MkdirAll(this_local_path, file_mode)
				if err != nil {
					obj.Logger.Error("os.MkdirAll fail!" + err.Error())
					return
				}
				obj.DownloadDir(this_upyun_path, this_local_path)
			} else if file_type == "file" {
				obj.DownloadFile(this_upyun_path, this_local_path)
			} else {
				obj.Logger.Critical("unknow file type2:" + file_type)
				return
			}
		}
		return
	}))
}

//delete a file
func (obj *BulkUpyun) DeleteFile(upyun_path string) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	obj.deleteFile(upyun_path, wg)
	obj.Tm.AddTaskNewThread(kmgTask.TaskFunc(func() {
		wg.Wait()
	}))
}

//resursive delete a dir
func (obj *BulkUpyun) DeleteDir(upyun_path string) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	obj.deleteDir(upyun_path, wg)
	obj.Tm.AddTaskNewThread(kmgTask.TaskFunc(func() {
		wg.Wait()
	}))
}

//delete a file
func (obj *BulkUpyun) deleteFile(upyun_path string, finish_wg *sync.WaitGroup) {
	obj.Tm.AddTask(kmgTask.TaskFunc(func() {
		defer finish_wg.Done()
		obj.Logger.Info("delete file: " + upyun_path)
		err := obj.UpYun.DeleteFile(upyun_path)
		if err != nil {
			obj.Logger.Error("delete file failed!:" + upyun_path + ":" + err.Error())
			return
		}
		return
	}))
}

//we need to know when is finish delete all file in it ,so we can delete the dir
func (obj *BulkUpyun) deleteDir(upyun_path string, finish_wg *sync.WaitGroup) {
	obj.Tm.AddTask(kmgTask.TaskFunc(func() {
		wg := &sync.WaitGroup{}
		defer obj.Tm.AddTaskNewThread(kmgTask.TaskFunc(func() {
			wg.Wait()
			wg.Add(1)
			obj.deleteFile(upyun_path, wg)
			wg.Wait()
			finish_wg.Done()
		}))
		obj.Logger.Info("delete dir: " + upyun_path)
		file_list, err := obj.UpYun.ReadDir(upyun_path)
		if err != nil {
			obj.Logger.LogError(err)
			return
		}
		for _, file_info := range file_list {
			file_type := file_info.Type
			file_name := file_info.Name
			this_upyun_path := upyun_path + "/" + file_name
			if file_type == "folder" {
				wg.Add(1)
				obj.deleteDir(this_upyun_path, wg)
			} else if file_type == "file" {
				wg.Add(1)
				obj.deleteFile(this_upyun_path, wg)
			} else {
				obj.Logger.Critical("unknow file type2:" + file_type)
				return
			}
		}
		return
	}))
}

func (obj *BulkUpyun) GetFileType(upyun_path string) (file_type string, err error) {
	info, err := obj.UpYun.GetFileInfo(upyun_path)
	if err != nil {
		return
	}
	file_type = info["type"]
	return
}

func (obj *BulkUpyun) GetFileSize(upyun_path string) (size uint64, err error) {
	info, err := obj.UpYun.GetFileInfo(upyun_path)
	if err != nil {
		return
	}
	size, err = strconv.ParseUint(info["size"], 10, 64)
	if err != nil {
		return
	}
	return
}

func (obj *BulkUpyun) GetFileDate(upyun_path string) (date time.Time, err error) {
	info, err := obj.UpYun.GetFileInfo(upyun_path)
	if err != nil {
		return
	}
	date, err = time.Parse(time.RFC1123, info["date"])
	if err != nil {
		return
	}
	return
}
