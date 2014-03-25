package kmgTask

import (
	"fmt"
	"sync"
)

type ErrorTaskFn func() (err error)

//提供线程数量限制和错误处理(重试+批量报错)的任务管理器
//适用于并发批量下载或上传
//TODO 测试出错情况,测试关闭情况
type LimitThreadErrorHandleTaskManager struct {
	ErrorArr   []error
	errorMutex sync.Mutex
	threadNum  int              //线程数目
	retryNum   int              //重试次数
	taskChan   chan ErrorTaskFn //任务通道
	wg         sync.WaitGroup   //等待任务完成
}

func NewLimitThreadErrorHandleTaskManager(threadNum int, retryNum int) *LimitThreadErrorHandleTaskManager {
	tm := &LimitThreadErrorHandleTaskManager{
		taskChan:  make(chan ErrorTaskFn),
		threadNum: threadNum,
		retryNum:  retryNum,
	}
	tm.run()
	return tm
}

func (m *LimitThreadErrorHandleTaskManager) run() {
	for i := 0; i < m.threadNum; i++ {
		go func() {
			for task := range m.taskChan {
				m.runOneTask(task)
			}
		}()
	}
}
func (m *LimitThreadErrorHandleTaskManager) runOneTask(task ErrorTaskFn) {
	defer m.wg.Done()
	retryNum := m.retryNum
Retry:
	err := task()
	if err != nil {
		fmt.Printf("%T %v\n", err, err)
		if retryNum >= 1 {
			retryNum--
			goto Retry
		}
		m.errorMutex.Lock()
		m.ErrorArr = append(m.ErrorArr, err)
		m.errorMutex.Unlock()
	}
}
func (m *LimitThreadErrorHandleTaskManager) AddTask(task ErrorTaskFn) {
	m.wg.Add(1)
	m.taskChan <- task
}

func (m *LimitThreadErrorHandleTaskManager) Wait() {
	m.wg.Wait()
}

func (m *LimitThreadErrorHandleTaskManager) GetError() error {
	totalErrStr := ""
	for _, err := range m.ErrorArr {
		totalErrStr += err.Error() + "\n"
	}
	if totalErrStr != "" {
		return fmt.Errorf("%s", totalErrStr)
	}
	return nil
}
func (m *LimitThreadErrorHandleTaskManager) Close() {
	close(m.taskChan)
}
