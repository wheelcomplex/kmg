package dependencyInjection

type ContainerBuilder struct {
}

const (
	ScopeContainer = "container"
	ScopePrototype = "prototype"
	ScopeRequest   = "request"
)

// For end user
type Container interface {
	//it is concurrency safe to call this function
	Get(id string) (interface{}, error)
	//set a new service into container
	//pass "" to scope to use container scope
	//it is not concurrency safe to call this function
	Set(id string, obj interface{}, scope string) error
	// set a new service into container
	// if error happen,it panic
	MustSet(id string, obj interface{}, scope string)
	Has(id string) bool
	//return a new container with that scope
	EnterScope(name string) (Container, error)
	//leave current scope
	LeaveScope() (Container, error)
	IsScopeActive(name string) bool
}
type ContainerAwareInterface interface {
	SetContainer(container Container)
}
