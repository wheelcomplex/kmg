package errors

import "fmt"

import "runtime"

type MsgWarpError struct {
	Err error
	Msg string
}

func (err *MsgWarpError) Error() string {
	return fmt.Sprintf("[%s] %v", err.Msg, err.Err)
}

//add more message to an error
func AddMessage(err error, msg string) error {
	return &MsgWarpError{Err: err, Msg: msg}
}

type CallerWarpError struct {
	Err  error
	File string
	Line int
}

func (err *CallerWarpError) Error() string {
	return fmt.Sprintf("[%s:%d] %v", err.File, err.Line, err.Err)
}
func AddCaller(err error) error {
	errOut := &CallerWarpError{Err: err}
	_, errOut.File, errOut.Line, _ = runtime.Caller(1)
	return errOut
}

type SprintfWrapError struct {
	args   []interface{}
	format string
}

func (err *SprintfWrapError) Error() string {
	return fmt.Sprintf(err.format, err.args...)
}
func Sprintf(format string, args ...interface{}) error {
	return &SprintfWrapError{format: format, args: args}
}
