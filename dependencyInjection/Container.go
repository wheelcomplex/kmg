package dependencyInjection

import (
	"errors"
	"fmt"
)

var ServiceIdNotExistError = errors.New("service id not exist")
var ScopeNotExistError = errors.New("scope not exist")
var ContainerLeaveScopeError = errors.New("container can not leave current scope,becase there is no a parent scope")

type container struct {
	service_map map[string]interface{}
	scope       string
	parent      *container
}

//a simple container implement,which can only set object in runtime.
func NewContainer() Container {
	return &container{
		service_map: make(map[string]interface{}),
		scope:       ScopeContainer,
	}
}
func (c *container) Get(id string) (interface{}, error) {
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
func (c *container) Set(id string, obj interface{}, scope string) error {
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
func (c *container) Has(id string) bool {
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
func (c *container) EnterScope(name string) (Container, error) {
	thisContainer := c
	for thisContainer != nil {
		if c.scope == name {
			return nil, errors.New(fmt.Sprintf("you enter %s scope twice", name))
		}
		thisContainer = c.parent
	}
	return &container{
		service_map: make(map[string]interface{}),
		scope:       name,
		parent:      c,
	}, nil
}

func (c *container) LeaveScope() (Container, error) {
	if c.parent == nil {
		return nil, ContainerLeaveScopeError
	}
	return c.parent, nil
}

func (c *container) IsScopeActive(name string) bool {
	thisContainer := c
	for thisContainer != nil {
		if c.scope == name {
			return true
		}
		thisContainer = c.parent
	}
	return false
}

func (c *container) MustSet(id string,obj interface {},scope string){
	err:=c.Set(id,obj,scope)
	if err!=nil{
		panic(err)
	}
}
