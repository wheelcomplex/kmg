package kmgTask

import "sync"

//implement TaskManager
//一个可以限制同时运行线程数目的taskmanager
type LimitThreadTaskManager struct {
	task_chan  chan Task      //任务通道
	num_thread int            //线程数目
	wg         sync.WaitGroup //等待任务完成
}

const limit_thread_task_buffer_size = 10000

var debug_num int32

//新建管理器
func NewLimitThreadTaskManager(num_thread int) *LimitThreadTaskManager {
	tm := &LimitThreadTaskManager{
		task_chan:  make(chan Task, limit_thread_task_buffer_size),
		num_thread: num_thread,
	}
	tm.run()
	return tm
}

// 开始运行管理器(异步)
func (t *LimitThreadTaskManager) run() {
	for i := 0; i < t.num_thread; i++ {
		go func() {
			for {
				task, ok := <-t.task_chan
				if ok == false {
					return
				}
				task.Run()
				t.wg.Done()
			}
		}()
	}
}

// 添加一个任务
func (t *LimitThreadTaskManager) AddTask(task Task) {
	t.wg.Add(1)
	t.task_chan <- task
}

//等待所有任务完成
func (t *LimitThreadTaskManager) Wait() {
	t.wg.Wait()
}

//关闭管理器
//需要等待所有任务完成后,返回
func (t *LimitThreadTaskManager) Close() {
	defer close(t.task_chan)
	t.Wait()
}

// 添加一个任务,运行在新线程上(不在线程限制里面)
func (t *LimitThreadTaskManager) AddTaskNewThread(task Task) {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		task.Run()

	}()
}
