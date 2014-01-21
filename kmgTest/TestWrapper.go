package kmgTest

import (
	"reflect"
	"strings"
)

//use T for error report,automatic call every function start with Test
//function start with Test must not have input arguments,output arguments.
func TestWarpper(T Fatalfer, testObject FatalferAware) {
	testObject.SetFatalfer(T)
	tov := reflect.ValueOf(testObject)
	//tov := reflect.Indirect(reflect.ValueOf(testObject))
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
