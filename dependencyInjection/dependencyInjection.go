package dependencyInjection


type ContainerBuilder struct{

}

const (
	ScopeContainer = "container"
	ScopePrototype = "prototype"
	ScopeRequest = "request"
)

// For end user
type ContainerInterface interface{
	Get(id string)(interface {},error)
	//set a new service into container
	//pass "" to scope to use container scope
	Set(id string,obj interface{},scope string) error
	Has(id string) bool
	//return a new container with that scope
	EnterScope(name string) (ContainerInterface,error)
	//leave current scope
	LeaveScope() (ContainerInterface,error)
	IsScopeActive(name string) bool
}
type ContainerAwareInterface interface{
	SetContainer(container ContainerInterface)
}



