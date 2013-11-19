package fsnotify

import "time"

//watch and run a function
// you should handle error in Watcher and add some files to watch (at least)
type Runner struct {
	*Watcher
	//wait WaitTime to run your function (default 0.2s)
	WaitTime time.Duration
	//whether init run this function (default true)
	InitRun bool
}

func NewRunner(bufferSize int) (*Runner, error) {
	watcher, err := NewWatcher(bufferSize)
	if err != nil {
		return nil, err
	}
	runner := &Runner{Watcher: watcher,
		WaitTime: time.Duration(0.2 * float64(time.Second)),
		InitRun:  true,
	}
	return runner, err
}

//this function will block forever
func (runner *Runner) Run(work func()) {
	lastHappendTime := time.Now()
	//start app when command start
	work()
	for {
		event := <-runner.Watcher.Event
		if event.Time.Before(lastHappendTime) {
			continue
		}
		//wait 200ms to prevent multiple restart in short time
		time.Sleep(runner.WaitTime)
		lastHappendTime = time.Now()
		work()
	}
}
