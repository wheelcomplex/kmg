package dependencyInjection

// complie passed after all extension loaded,load some tag as ...
type CompilePassInterface interface {
	CompilePass(c *ContainerBuilder) error
}
