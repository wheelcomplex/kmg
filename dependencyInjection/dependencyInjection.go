package dependencyInjection

const (
	ScopeSingleton = "singleton"
	ScopePrototype = "prototype"
	ScopeRequest   = "request"
)

type ContainerAwareInterface interface {
	SetContainer(container *Container)
}
