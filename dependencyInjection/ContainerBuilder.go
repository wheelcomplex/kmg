package dependencyInjection

type ContainerBuilder struct {
	definition_map map[string]*Definition
	extensions     []ExtensionInterface
	compliePasses  []CompilePassInterface
	bootes         []BootInterface
	Parameters     map[string]string
}

func NewContainerBuilder() *ContainerBuilder {
	return &ContainerBuilder{
		definition_map: make(map[string]*Definition),
		Parameters:     make(map[string]string),
	}
}
func (builder *ContainerBuilder) AddExtension(extension ExtensionInterface) {
	builder.extensions = append(builder.extensions, extension)
}
func (builder *ContainerBuilder) AddCompilePass(compliePass CompilePassInterface) {
	builder.compliePasses = append(builder.compliePasses, compliePass)
}
func (builder *ContainerBuilder) AddBoot(boot BootInterface) {
	builder.bootes = append(builder.bootes, boot)
}
func (builder *ContainerBuilder) GetDefinition(id string) (definition *Definition, exist bool) {
	definition, exist = builder.definition_map[id]
	return
}

func (builder *ContainerBuilder) SetDefinition(definition *Definition) error {
	err := definition.Init()
	if err != nil {
		return err
	}
	builder.definition_map[definition.Id] = definition
	return nil
}

func (builder *ContainerBuilder) MustSetDefinition(definition *Definition) *Definition {
	err := builder.SetDefinition(definition)
	if err != nil {
		panic(err)
	}
	return definition
}

// will return a empty slice if tag not exist
func (builder *ContainerBuilder) GetTaggedDefinition(tag string) []*Definition {
	//TODO can use cache improve performance
	definitions := []*Definition{}
	for _, v := range builder.definition_map {
		if v.HasTag(tag) {
			definitions = append(definitions, v)
		}
	}
	return definitions
}
func (builder *ContainerBuilder) Set(id string, obj interface{}, scope string) error {
	if scope == "" {
		scope = ScopeSingleton
	}
	if scope == ScopePrototype {
		return CanNotSetScopePrototypeByObjError
	}
	definition := NewDefinitionFromInst(obj)
	definition.Scope = scope
	definition.Id = id

	builder.definition_map[id] = definition
	return nil
}
func (builder *ContainerBuilder) MustSet(id string, obj interface{}, scope string) {
	err := builder.Set(id, obj, scope)
	if err != nil {
		panic(err)
	}
}
func (builder *ContainerBuilder) SetFactory(id string, factory func(c *Container) (interface{}, error), scope string) error {
	if scope == "" {
		scope = ScopeSingleton
	}
	definition := NewDefinitionFromFactory(factory)
	definition.Scope = scope
	definition.Id = id
	builder.definition_map[id] = definition
	return nil
}
func (builder *ContainerBuilder) MustSetFactory(id string, factory func(c *Container) (interface{}, error), scope string) {
	err := builder.SetFactory(id, factory, scope)
	if err != nil {
		panic(err)
	}
}

func (builder *ContainerBuilder) Compile() (c *Container, err error) {
	for k, v := range builder.Parameters {
		err = builder.Set("parameter."+k, v, "")
		if err != nil {
			return
		}
	}
	for _, v := range builder.extensions {
		err = v.LoadDependencyInjection(builder)
		if err != nil {
			return
		}
	}

	for _, v := range builder.compliePasses {
		err = v.CompilePass(builder)
		if err != nil {
			return
		}
	}
	c = &Container{
		definition_map: builder.definition_map,
	}
	c.init()
	for _, v := range builder.bootes {
		err = v.Boot(c)
		if err != nil {
			return
		}
	}
	return
}
func (builder *ContainerBuilder) MustCompile() (c *Container) {
	c, err := builder.Compile()
	if err != nil {
		panic(err)
	}
	return c
}
