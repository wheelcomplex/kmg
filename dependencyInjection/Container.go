package dependencyInjection

import (
	"errors"
	"fmt"
)

var ServiceIdNotExistError = errors.New("service id not exist")
var ScopeNotExistError = errors.New("scope not exist")
var ScopeNotActiveError = errors.New("scope not active")
var CanNotSetScopePrototypeByObjError = errors.New("can not set service in scope prototype by obj")
var CanNotSetNotActiveScopeByObjError = errors.New("can not set service in not active scope by obj")
var ParentScopeNotExistError = errors.New("container can not leave current scope,becase there is no a parent scope")

type Container struct {
	definition_map map[string]*Definition //share inst between scope
	service_map    map[string]interface{} //only current scope
	scope_map      map[string]*Container  //contain all scope
	scope          string
	parent         *Container
}
type Definition struct {
	Id      string
	Scope   string
	Factory func(c *Container) (interface{}, error)
}

//a simple container implement,which can only set object in runtime.
func NewContainer() *Container {
	container := &Container{
		service_map:    make(map[string]interface{}),
		definition_map: make(map[string]*Definition),
		scope_map:      make(map[string]*Container),
		scope:          ScopeSingleton,
	}
	container.scope_map[ScopeSingleton] = container
	container.scope_map[ScopePrototype] = container
	return container
}

func (c *Container) Get(id string) (service interface{}, err error) {
	definition, ok := c.definition_map[id]
	if !ok {
		return nil, ServiceIdNotExistError
	}
	if !c.IsScopeActive(definition.Scope) {
		return nil, ScopeNotActiveError
	}
	container := c.scope_map[definition.Scope]
	//TODO Concurrent safe
	service, ok = container.service_map[id]
	if ok {
		return service, nil
	}
	service, err = definition.Factory(c)
	if err != nil {
		return nil, err
	}
	container.service_map[id] = service
	return
}
func (c *Container) MustGet(id string) interface{} {
	service, err := c.Get(id)
	if err != nil {
		panic(err)
	}
	return service
}

//TODO Concurrent safe
func (c *Container) Set(id string, obj interface{}, scope string) error {
	if scope == "" {
		scope = ScopeSingleton
	}
	if scope == ScopePrototype {
		return CanNotSetScopePrototypeByObjError
	}
	if !c.IsScopeActive(scope) {
		return CanNotSetNotActiveScopeByObjError
	}
	definition, ok := c.definition_map[id]
	if !ok {
		definition = &Definition{}
		c.definition_map[id] = definition
	}
	definition.Scope = scope
	definition.Id = id

	c.scope_map[scope].service_map[id] = obj
	return nil
}
func (c *Container) MustSet(id string, obj interface{}, scope string) {
	err := c.Set(id, obj, scope)
	if err != nil {
		panic(err)
	}
}
func (c *Container) SetFactory(id string, factory func(c *Container) (interface{}, error), scope string) error {
	if scope == "" {
		scope = ScopeSingleton
	}
	//remove old inst
	c.removeInst(id)
	definition, ok := c.definition_map[id]
	if !ok {
		definition = &Definition{}
		c.definition_map[id] = definition
	}
	definition.Scope = scope
	definition.Id = id
	definition.Factory = factory
	return nil
}
func (c *Container) MustSetFactory(id string, factory func(c *Container) (interface{}, error), scope string) {
	err := c.SetFactory(id, factory, scope)
	if err != nil {
		panic(err)
	}
}

func (c *Container) Has(id string) bool {
	_, ok := c.definition_map[id]
	return ok
}
func (c *Container) EnterScope(scope string) (*Container, error) {
	if c.IsScopeActive(scope) {
		return nil, errors.New(fmt.Sprintf("you enter %s scope twice", scope))
	}
	scope_map := make(map[string]*Container)
	for k, v := range c.scope_map {
		scope_map[k] = v
	}
	container := &Container{
		service_map:    make(map[string]interface{}),
		definition_map: c.definition_map,
		scope:          scope,
		parent:         c,
	}
	scope_map[scope] = container
	container.scope_map = scope_map
	return container, nil
}
func (c *Container) LeaveScope() (*Container, error) {
	if c.parent != nil {
		return c.parent, nil
	}
	return nil, ParentScopeNotExistError
}
func (c *Container) IsScopeActive(scope string) bool {
	_, ok := c.scope_map[scope]
	return ok
}

//try to remove saved inst
func (c *Container) removeInst(id string) {
	definition, ok := c.definition_map[id]
	if !ok {
		return
	}
	container, ok := c.scope_map[definition.Scope]
	if !ok {
		return
	}
	delete(container.service_map, id)
}

/*

func (c *Container) Get(id string) (interface{}, error) {
	thisContainer := c
	for thisContainer != nil {
		obj, ok := c.service_map[id]
		if ok == true {
			return obj, nil
		}
		thisContainer = c.parent
	}
	return nil, ServiceIdNotExistError

}
func (c *Container) Set(id string, obj interface{}, scope string) error {
	if scope == "" {
		scope = ScopeContainer
	}
	thisContainer := c
	for thisContainer != nil {
		if thisContainer.scope == scope {
			thisContainer.service_map[id] = obj
			return nil
		}
		thisContainer = c.parent
	}
	return ScopeNotExistError
}
func (c *Container) Has(id string) bool {
	thisContainer := c
	for thisContainer != nil {
		_, ok := c.service_map[id]
		if ok == true {
			return true
		}
		thisContainer = c.parent
	}
	return false
}
func (c *Container) EnterScope(name string) (*Container, error) {
	thisContainer := c
	for thisContainer != nil {
		if c.scope == name {
			return nil, errors.New(fmt.Sprintf("you enter %s scope twice", name))
		}
		thisContainer = c.parent
	}
	return &Container{
		service_map: make(map[string]interface{}),
		scope:       name,
		parent:      c,
	}, nil
}

func (c *Container) LeaveScope() (*Container, error) {
	if c.parent == nil {
		return nil, ContainerLeaveScopeError
	}
	return c.parent, nil
}

func (c *Container) IsScopeActive(name string) bool {
	thisContainer := c
	for thisContainer != nil {
		if c.scope == name {
			return true
		}
		thisContainer = c.parent
	}
	return false
}

func (c *Container) MustSet(id string, obj interface{}, scope string) {
	err := c.Set(id, obj, scope)
	if err != nil {
		panic(err)
	}
}
*/
