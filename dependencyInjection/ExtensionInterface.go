package dependencyInjection

type ExtensionInterface interface {
	LoadDependencyInjection(c *Container) error
}
