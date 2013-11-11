package fsnotify

type Watcher struct {
	Error chan error
	Event chan *FileEvent
}
