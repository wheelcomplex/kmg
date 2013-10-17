package dependencyInjection

type ContainerBuilder struct{

}

type ContainerInterface interface{

}
type ContainerAwareInterface interface{
	SetContainer(container ContainerInterface)
}
