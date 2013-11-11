package dependencyInjection

//modular register some dependency
type ExtensionInterface interface {
	LoadDependencyInjection(c *ContainerBuilder) error
}
