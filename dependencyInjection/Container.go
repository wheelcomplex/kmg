package dependencyInjection

import (
	"errors"
	"fmt"
	"github.com/bronze1man/kmg/kmgReflect"
	"reflect"
	"sync"
)

var ServiceIdNotExistError = errors.New("service id not exist")
var ServiceTypeNotExistError = errors.New("service type not exist")
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
	serivce_lock   sync.RWMutex
}

//a simple container implement,which can only set object in runtime.
func NewContainer() *Container {
	container := &Container{
		definition_map: make(map[string]*Definition),
	}
	container.init()
	return container
}

func (c *Container) init() {
	c.scope = ScopeSingleton
	c.service_map = make(map[string]interface{})
	c.scope_map = make(map[string]*Container)
	c.scope_map[ScopeSingleton] = c
	c.scope_map[ScopePrototype] = c
}

//init all singleton object, to avoid data race,eat memory stuff
func (c *Container) build() (err error) {
	for name, def := range c.definition_map {
		if def.Scope != ScopeSingleton {
			continue
		}
		c.serivce_lock.RLock()
		_, ok := c.service_map[name]
		c.serivce_lock.RUnlock()
		if ok {
			return
		}
		service, err := def.GetInst(c)
		if err != nil {
			return err
		}
		c.serivce_lock.Lock()
		c.service_map[name] = service
		c.serivce_lock.Unlock()
	}
	return nil
}
func (c *Container) Get(id string) (service interface{}, err error) {
	definition, ok := c.definition_map[id]
	if !ok {
		return nil, fmt.Errorf("service id:%s not exist", id)
	}
	if !c.IsScopeActive(definition.Scope) {
		return nil, fmt.Errorf("service id:%s scope:%s not active", id, definition.Scope)
	}
	container := c.scope_map[definition.Scope]
	//TODO Concurrent safe not work...
	c.serivce_lock.RLock()
	service, ok = container.service_map[id]
	c.serivce_lock.RUnlock()
	if ok {
		return service, nil
	}
	service, err = definition.GetInst(c)
	if err != nil {
		return nil, err
	}
	if definition.Scope != ScopePrototype {
		c.serivce_lock.Lock()
		container.service_map[id] = service
		c.serivce_lock.Unlock()
	}
	return
}
func (c *Container) MustGet(id string) interface{} {
	service, err := c.Get(id)
	if err != nil {
		panic(err)
	}
	return service
}
func (c *Container) MustGetString(id string) string {
	service, err := c.Get(id)
	if err != nil {
		panic(err)
	}
	return service.(string)
}
func (c *Container) GetByType(ti interface{}) (service interface{}, err error) {
	var t reflect.Type
	switch ti.(type) {
	case reflect.Type:
		t = ti.(reflect.Type)
	default:
		t = reflect.TypeOf(ti)
	}
	typeName := kmgReflect.GetTypeFullName(t)
	if typeName == "" {
		return nil, ServiceTypeNotExistError
	}
	return c.Get(typeName)
}
func (c *Container) MustGetByType(ti interface{}) (service interface{}) {
	service, err := c.GetByType(ti)
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
	definition := &Definition{}
	definition.Scope = scope
	definition.Id = id
	definition.Inst = obj
	err := definition.Init()
	if err != nil {
		return err
	}
	c.definition_map[id] = definition
	c.serivce_lock.Lock()
	defer c.serivce_lock.Unlock()
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
	definition.InitType = DefinitionFromFactory
	return nil
}
func (c *Container) MustSetFactory(id string, factory func(c *Container) (interface{}, error), scope string) {
	err := c.SetFactory(id, factory, scope)
	if err != nil {
		panic(err)
	}
}

//has this service, maybe not active
func (c *Container) Has(id string) bool {
	_, ok := c.definition_map[id]
	return ok
}

// has this service, and active
func (c *Container) IsActiveService(id string) bool {
	def, ok := c.definition_map[id]
	if !ok {
		return false
	}
	if !c.IsScopeActive(def.Scope) {
		return false
	}
	return true
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
	c.serivce_lock.Lock()
	delete(container.service_map, id)
	c.serivce_lock.Unlock()
}
