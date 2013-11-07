package ajkApi

import (
	//"kmg/dependencyInjection"
	//"kmg/test"
	"errors"
	"testing"
)

type TestService struct {
	TestFunc1Num int
	TestFunc4A   int
}

func (this *TestService) TestFunc1() {
	this.TestFunc1Num = 10
}
func (this *TestService) TestFunc2() error {
	return errors.New("test error")
}
func (this *TestService) TestFunc3(apiInput *struct{ a int }, apiOutput *struct{ b int }) error {
	apiOutput.b = apiInput.a + 1
	return nil
}

func (this *TestService) TestFunc4(apiInput *struct{ a int }) error {
	this.TestFunc4A = apiInput.a + 1
	return nil
}
func TestContainerAwareApiManager(ot *testing.T) {
	/*
		var err error
		t := &test.TestTools{T: ot}
		container := dependencyInjection.NewContainer()
		apiManager := NewApiManagerFromContainer(container)
		output:=&struct{b int}{}

		testService :=  &TestService{}
		err = container.Set("TestService",testService , "")
		t.Equal(err, nil)

		err = apiManager.RpcCall(
			nil, "TestService.TestFunc1", struct{}{}, output)
		t.Equal(err, nil)
		t.Equal(testService.TestFunc1Num,10)

		err = apiManager.RpcCall(nil, "TestService1.TestFunc1", struct{}{}, output)
		t.Equal(err.(*ApiFuncNotFoundError).Reason, "service not exist")

		err = apiManager.RpcCall(nil, "TestService.TestFunc3", &struct{a int}{a:1}, output)
		t.Equal(err,nil)
		t.Equal(output.b,2)

		err = apiManager.RpcCall(nil, "TestService.TestFunc4", &struct{a int}{a:5}, output)
		t.Equal(err,nil)
		t.Equal(testService.TestFunc4A,6)
	*/
}
