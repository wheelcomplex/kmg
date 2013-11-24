package dependencyInjection

import (
	"errors"
	"fmt"
	"github.com/bronze1man/kmg/kmgReflect"
	"reflect"
)

type DefinitionInitType int

const (
	DefinitionFromInst DefinitionInitType = iota
	DefinitionFromFactory
	DefinitionFromType
)

/*
 definition of a service
 if you set Inst,not set Type or TypeReflect, TypeReflect and TypeName will come from Inst.
 if you set Type, TypeReflect and TypeName will come from Type, TypeReflect will be ignore.
 if you set TypeReflect, TypeName will come from TypeReflect.
 if Inst not assignable to TypeReflect, Init() will return err.
 Do not set TypeName, it does not work, golang can not get TypeReflect from TypeName.
*/
type Definition struct {
	Id          string
	Scope       string
	Inst        interface{}
	Factory     func(c *Container) (interface{}, error)
	Tags        []string
	InitType    DefinitionInitType
	Type        interface{}
	TypeReflect reflect.Type
	TypeName    string
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

//complete some default values.
func (definition *Definition) Init() (err error) {
	//guess Scope
	if definition.Scope == "" {
		definition.Scope = ScopeSingleton
	}
	isReflectFromInst := false
	//guess TypeName
	if definition.Inst != nil &&
		definition.Type == nil &&
		definition.TypeReflect == nil {
		definition.TypeReflect = reflect.TypeOf(definition.Inst)
		isReflectFromInst = true
	}
	if definition.Type != nil {
		definition.TypeReflect = reflect.TypeOf(definition.Type)
	}
	if definition.TypeReflect != nil {
		definition.TypeName, _ = kmgReflect.GetTypeFullName(definition.TypeReflect)
	}
	//guess Id
	if definition.Id == "" {
		definition.Id = definition.TypeName
	}
	if definition.Id == "" {
		return errors.New("definition not has id,and can not guess it")
	}

	//not assignable
	if !isReflectFromInst &&
		definition.Inst != nil &&
		!reflect.TypeOf(definition.Inst).AssignableTo(definition.TypeReflect) {
		return fmt.Errorf("Inst is not assignable to TypeReflect serviceId:%s", definition.Id)
	}
	//guess InitType
	if definition.InitType == 0 {
		switch {
		case definition.Inst != nil:
			definition.InitType = DefinitionFromInst
		case definition.Factory != nil:
			definition.InitType = DefinitionFromFactory
		case definition.TypeReflect != nil:
			definition.InitType = DefinitionFromType
		default:
			return fmt.Errorf("definition not has init way, id: %s", definition.Id)
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
	case DefinitionFromType:
		return definition.getInstFromType(c)
	default:
		return nil, fmt.Errorf("invalid definition init type,id: %s", definition.Id)
	}
	return nil, fmt.Errorf("invalid definition init type,id: %s", definition.Id)
}

func (definition *Definition) getInstFromType(c *Container) (interface{}, error) {
	t := definition.TypeReflect
	isPtr := false
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		isPtr = true
	}
	var valueP reflect.Value
	switch t.Kind() {
	case reflect.Ptr:
		return nil, fmt.Errorf("can not get inst from type %T", reflect.Zero(definition.TypeReflect).Interface())
	case reflect.Struct:
		valueP = reflect.New(t)
		value := valueP.Elem()
		for i := 0; i < t.NumField(); i++ {
			ft := t.Field(i)
			//guess from name
			if c.IsActiveService(ft.Name) && c.definition_map[ft.Name].TypeReflect == ft.Type {
				finst, err := c.Get(ft.Name)
				if err != nil {
					return nil, err
				}
				value.Field(i).Set(reflect.ValueOf(finst))
				continue
			}
			//guess from type
			fullName, ok := kmgReflect.GetTypeFullName(ft.Type)
			if ok {
				if c.IsActiveService(fullName) && c.definition_map[fullName].TypeReflect == ft.Type {
					finst, err := c.Get(fullName)
					if err != nil {
						return nil, err
					}
					value.Field(i).Set(reflect.ValueOf(finst))
					continue
				}
			}
			//zero value not need to set
		}
	default:
		valueP = reflect.New(t)
	}
	if isPtr {
		return valueP.Interface(), nil
	}
	return valueP.Elem().Interface(), nil
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
