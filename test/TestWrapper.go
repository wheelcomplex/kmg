package test

import (
	"reflect"
	"strings"
)

func TestWrapper(T Fatalfer, testObject FatalferAware) {
	testObject.SetFatalfer(T)
	tov := reflect.Indirect(reflect.ValueOf(testObject))
	tot := tov.Type()
	for i := 0; i < tov.NumMethod(); i++ {
		tm := tot.Method(i)
		if !strings.HasPrefix(tm.Name, "Test") {
			continue
		}
		tmt := tm.Type
		//no argument
		if tmt.NumIn() != 1 {
			T.Fatalf("Testfunction:%s should not have input argument", tm.Name)
			return
		}
		if tmt.NumOut() != 0 {
			T.Fatalf("Testfunction:%s should not have output argument", tm.Name)
			return
		}
		tov.Method(i).Call([]reflect.Value{})
	}
}
