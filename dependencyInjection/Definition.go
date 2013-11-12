package dependencyInjection

import (
	"github.com/bronze1man/kmg/errors"
	"reflect"
	//"fmt"
	"github.com/bronze1man/kmg/kmgReflect"
)

type DefinitionInitType int

const (
	DefinitionFromInst DefinitionInitType = iota
	DefinitionFromFactory
)

type Definition struct {
	Id       string
	Scope    string
	Inst     interface{}
	Factory  func(c *Container) (interface{}, error)
	Tags     []string
	InitType DefinitionInitType
}

func NewDefinitionFromInst(inst interface{}) *Definition {
	return &Definition{Inst: inst, InitType: DefinitionFromInst}
}
func NewDefinitionFromFactory(factory func(c *Container) (interface{}, error)) *Definition {
	return &Definition{Factory: factory, InitType: DefinitionFromFactory}
}
func (definition *Definition) HasTag(tag string) bool {
	for _, v := range definition.Tags {
		if v == tag {
			return true
		}
	}
	return false
}

func (definition *Definition) AddTag(tag string) *Definition {
	definition.Tags = append(definition.Tags, tag)
	return definition
}
func (definition *Definition) Init() error {
	if definition.Scope == "" {
		definition.Scope = ScopeSingleton
	}
	definition.guessId()
	if definition.Id == "" {
		return errors.Sprintf("definition not has id", definition.Id)
	}
	if definition.InitType == 0 {
		switch {
		case definition.Inst != nil:
			definition.InitType = DefinitionFromInst
		case definition.Factory != nil:
			definition.InitType = DefinitionFromFactory
		default:
			return errors.Sprintf("definition not has init way, id: %s", definition.Id)
		}
	}
	return nil
}
func (definition *Definition) GetInst(c *Container) (interface{}, error) {
	switch definition.InitType {
	case DefinitionFromInst:
		return definition.Inst, nil
	case DefinitionFromFactory:
		return definition.Factory(c)
	default:
		return nil, errors.Sprintf("invalid definition init type,id: %s", definition.Id)
	}
	return nil, errors.Sprintf("invalid definition init type,id: %s", definition.Id)
}

func (definition *Definition) guessId() {
	if definition.Id != "" {
		return
	}
	if definition.Inst == nil {
		return
	}
	name, ok := kmgReflect.GetTypeFullName(reflect.TypeOf(definition.Inst))
	if ok {
		definition.Id = name
	}
}
