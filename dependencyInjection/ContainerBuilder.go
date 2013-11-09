package dependencyInjection

type ContainerBuilder struct {
}

func (builder *ContainerBuilder) SetFactoryFunc(id string, new func() interface{}) {

}
