package kmgTask

import "sync"

//implement TaskManager
//简单版taskmanager
//不做任何限制
type SimpleTaskManager struct {
	wg sync.WaitGroup //等待任务完成
}

// 添加一个任务
func (t *SimpleTaskManager) AddTask(task Task) {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		task.Run()
	}()
}

//等待所有任务完成
func (t *SimpleTaskManager) Wait() {
	t.wg.Wait()
}

//关闭管理器
//需要等待所有任务完成后,返回
func (t *SimpleTaskManager) Close() {
	t.Wait()
}
