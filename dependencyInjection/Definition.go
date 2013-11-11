package dependencyInjection

import (
	"errors"
)

type DefinitionInitType int

const (
	DefinitionFromInst DefinitionInitType = iota
	DefinitionFromFactory
)

var InvalidDefinitionInitTypeError = errors.New("invalid definition init type")

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

func (definition *Definition) GetInst(c *Container) (interface{}, error) {
	switch definition.InitType {
	case DefinitionFromInst:
		return definition.Inst, nil
	case DefinitionFromFactory:
		return definition.Factory(c)
	default:
		return nil, InvalidDefinitionInitTypeError
	}
	return nil, InvalidDefinitionInitTypeError
}
