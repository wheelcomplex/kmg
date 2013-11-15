package dependencyInjection

type BootInterface interface {
	Boot(c *Container) error
}
