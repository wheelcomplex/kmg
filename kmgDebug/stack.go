package kmgDebug

import (
	"fmt"
	"runtime"
)

type StackFuncCall struct {
	File     string
	Line     int
	Pc       uintptr
	FuncName string
}

type Stack []StackFuncCall

func (s *Stack) ToString() (output string) {
	output = ""
	for _, call := range *s {
		output += fmt.Sprintf("%s\n\t%s:%d:%x\n", call.FuncName, call.File, call.Line, call.Pc)
	}
	return
}

func GetCurrentStack(skip int) (stack *Stack) {
	pcs := make([]uintptr, 32)
	thisLen := runtime.Callers(skip+2, pcs)
	s := make(Stack, thisLen)
	stack = &s
	for i := 0; i < thisLen; i++ {
		f := runtime.FuncForPC(pcs[i])
		file, line := f.FileLine(pcs[i])
		(*stack)[i] = StackFuncCall{
			Pc:       pcs[i],
			FuncName: f.Name(),
			File:     file,
			Line:     line,
		}
	}
	return
}
