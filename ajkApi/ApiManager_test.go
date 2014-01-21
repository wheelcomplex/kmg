package ajkApi

import (
	"errors"
	//"github.com/bronze1man/kmg/dependencyInjection"
	//"github.com/bronze1man/kmg/kmgTest"
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
	return
	/*
		var err error
		t := &test.TestTools{T: ot}
		container := dependencyInjection.NewContainer()
		apiManager := NewApiManagerFromContainer(container)

		testService := &TestService{}
		err = container.Set("a/b.TestService", testService, "")
		t.Equal(err, nil)

		err = apiManager.RpcCall(
			nil, "a/b.TestService.TestFunc1", func(meta *ApiFuncMeta) error {
				t.Equal(meta.MethodName, "TestFunc1")
				t.Equal(meta.AttachObject.Interface(), testService)
				return nil
			})
		t.Equal(err, nil)
	*/
}
