package ajkApi

import (
	"fmt"
	"kmg/dependencyInjection"
	"reflect"
	"strings"
)

type containerAwareApiManager struct {
	c *dependencyInjection.Container
}
type ApiFuncArgumentError struct {
	Reason  string
	ApiName string
}

func (err *ApiFuncArgumentError) Error() string {
	return fmt.Sprintf("api argument error, reason:%s, name:%s", err.Reason, err.ApiName)
}

type ApiFuncNotFoundError struct {
	Reason  string
	ApiName string
}

func (err *ApiFuncNotFoundError) Error() string {
	return fmt.Sprintf("api function not found, reason:%s, name:%s", err.Reason, err.ApiName)
}

/*
 container service + method -> api
 the api name will be "serviceName.methodName"
*/
func NewApiManagerFromContainer(c *dependencyInjection.Container) ApiManagerInterface {
	return &containerAwareApiManager{c: c}
}
func (manager *containerAwareApiManager) RpcCall(
	session *Session,
	name string,
	caller func(*ApiFuncMeta) error,
) error {
	dotP := strings.LastIndex(name, ".")
	if dotP == -1 {
		return &ApiFuncNotFoundError{Reason: "name not cantain .", ApiName: name}
	}
	c, err := manager.c.EnterScope(dependencyInjection.ScopeRequest)
	if err != nil {
		return err
	}
	err = c.Set("session", session, dependencyInjection.ScopeRequest)
	if err != nil {
		return err
	}

	defer c.LeaveScope()
	serviceName := name[:dotP]
	if !c.Has(serviceName) {
		return &ApiFuncNotFoundError{Reason: "service not exist", ApiName: name}
	}
	service, err := c.Get(serviceName)
	if err != nil {
		return err
	}
	serviceType := reflect.TypeOf(service)
	methodName := name[dotP+1:]
	method, ok := serviceType.MethodByName(methodName)
	if ok == false {
		return &ApiFuncNotFoundError{Reason: "method not on service", ApiName: name}
	}
	return caller(&ApiFuncMeta{IsMethod: true, Func: method.Func, AttachObject: reflect.ValueOf(service)})
}
