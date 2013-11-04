package ajkApi

import "kmg/dependencyInjection"

type containerAwareApiManager struct{
	c dependencyInjection.ContainerInterface
}
func NewApiManagerFromContainer(c dependencyInjection.ContainerInterface)ApiManagerInterface{
	return &containerAwareApiManager{c:c}
}
func (manager *containerAwareApiManager) RpcCall(
		session *Session,
		name string,
		input interface{},
		output interface{},
	)error{

}
