package kmgTest

import (
	"fmt"
	"reflect"
	"strings"
)

//use T for error report,automatic call every function start with Test
//function start with Test must not have input arguments,output arguments.
func TestWarpper(T TestingTB, testObject TestingTBAware) {
	testObject.SetTestingTB(T)
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
			fmt.Printf("[kmgTest.TestWarpper] Testfunction:%s should not have input argument\n", tm.Name)
			T.FailNow()
			return
		}
		if tmt.NumOut() != 0 {
			fmt.Printf("[kmgTest.TestWarpper] Testfunction:%s should not have output argument\n", tm.Name)
			T.FailNow()
			return
		}
		tov.Method(i).Call([]reflect.Value{})
	}
}
