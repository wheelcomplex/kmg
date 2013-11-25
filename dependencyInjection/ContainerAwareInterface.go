package dependencyInjection

type ContainerAwareInterface interface {
	SetContainer(container *Container)
}
