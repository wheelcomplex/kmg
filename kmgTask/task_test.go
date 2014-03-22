package kmgTask

import "testing"

func TestOne(t *testing.T) {
	tm := NewLimitThreadTaskManager(1)
	result_chan := make(chan int)
	var result_var int
	task := TaskFunc(func() {
		result_chan <- 1
		result_var = 2
	})
	tm.AddTask(task)
	tm.AddTaskNewThread(task)
	if <-result_chan != 1 {
		t.Fatalf("not run result_chan not match")
	}
	if <-result_chan != 1 {
		t.Fatalf("not run result_chan not match")
	}
	tm.Close()
	if result_var != 2 {
		t.Fatalf("not run result_var not match")
	}
}

func TestOneThread(t *testing.T) {
	tm := NewLimitThreadTaskManager(1)
	var result_var int
	result_chan := make(chan int)
	task := TaskFunc(func() {
		result_var = 2
		result_chan <- 1
	})
	tm.AddTask(task)
	tm.AddTask(task)
	for i := 0; i < 2; i++ {
		if <-result_chan != 1 {
			t.Fatalf("result_chan not match")
		}
	}
	tm.Close()
	if result_var != 2 {
		t.Fatalf("result_var not match")
	}
}

func BenchmarkMulitThread(b *testing.B) {
	tm := NewLimitThreadTaskManager(10)
	tm.AddTask(TaskFunc(func() {
		task := TaskFunc(func() {
		})
		for i := 0; i < b.N; i++ {
			tm.AddTask(task)
		}
	}))

	tm.Close()
}
