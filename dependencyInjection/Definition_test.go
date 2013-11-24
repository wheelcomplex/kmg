package dependencyInjection

import (
	"github.com/bronze1man/kmg/test"
	"reflect"
	"testing"
)

type DefinitionTestT1 struct {
	A int
	B *DefinitionTestT2
	C int
	D DefinitionTestI1
}
type DefinitionTestT2 struct {
	A int
}
type DefinitionTestI1 interface {
	GetNum() int
}
type DefinitionTestT3 struct {
	A int
}

func (t3 *DefinitionTestT3) GetNum() int {
	return t3.A
}
func TestDefinitionInitInst(ot *testing.T) {
	t := test.NewTestTools(ot)
	c := NewContainer()
	//init from Inst
	inst := &DefinitionTestT1{}
	def := &Definition{
		Inst: inst,
	}
	err := def.Init()
	t.Equal(err, nil)
	t.Equal(def.Id, "github.com/bronze1man/kmg/dependencyInjection.DefinitionTestT1")
	t.Equal(def.InitType, DefinitionFromInst)
	t.Equal(def.Scope, ScopeSingleton)
	t.Equal(def.TypeReflect, reflect.TypeOf((*DefinitionTestT1)(nil)))
	t.Equal(def.TypeName, "github.com/bronze1man/kmg/dependencyInjection.DefinitionTestT1")

	ret, err := def.GetInst(c)
	t.Equal(err, nil)
	t.Equal(ret, inst)
}

func TestDefinitionInitInterface(ot *testing.T) {
	t := test.NewTestTools(ot)
	c := NewContainer()
	//init from Inst
	inst := &DefinitionTestT3{}
	def := &Definition{
		Inst:        inst,
		TypeReflect: reflect.TypeOf((*DefinitionTestI1)(nil)).Elem(),
	}
	err := def.Init()
	t.Equal(err, nil)
	t.Equal(def.Id, "github.com/bronze1man/kmg/dependencyInjection.DefinitionTestI1")
	t.Equal(def.InitType, DefinitionFromInst)
	t.Equal(def.Scope, ScopeSingleton)
	t.Equal(def.TypeReflect, reflect.TypeOf((*DefinitionTestI1)(nil)).Elem())
	t.Equal(def.TypeName, "github.com/bronze1man/kmg/dependencyInjection.DefinitionTestI1")

	ret, err := def.GetInst(c)
	t.Equal(err, nil)
	t.Equal(ret, inst)
}

func TestDefinitionInitScope(ot *testing.T) {
	t := test.NewTestTools(ot)
	//scope
	def := &Definition{
		Inst:  &DefinitionTestT1{},
		Scope: ScopePrototype,
	}
	err := def.Init()
	t.Equal(err, nil)
	t.Equal(def.Scope, ScopePrototype)

}

func TestDefinitionInitFactory(ot *testing.T) {
	t := test.NewTestTools(ot)
	c := NewContainer()
	// init from factory
	def := &Definition{
		Id: "123",
		Factory: func(c *Container) (interface{}, error) {
			return &DefinitionTestT1{A: 5}, nil
		},
	}
	err := def.Init()
	t.Equal(err, nil)
	t.Equal(def.Id, "123")
	t.Equal(def.InitType, DefinitionFromFactory)
	t.Equal(def.Scope, ScopeSingleton)

	ret, err := def.GetInst(c)
	t.Equal(err, nil)
	t.Equal(ret.(*DefinitionTestT1).A, 5)
}
func TestDefinitionInitType(ot *testing.T) {
	t := test.NewTestTools(ot)
	c := NewContainer()
	//init from type
	def := &Definition{
		Type: (*DefinitionTestT1)(nil),
	}
	err := def.Init()
	t.Equal(err, nil)
	t.Equal(def.Id, "github.com/bronze1man/kmg/dependencyInjection.DefinitionTestT1")
	t.Equal(def.InitType, DefinitionFromType)
	t.Equal(def.Scope, ScopeSingleton)
	t.Equal(def.TypeReflect, reflect.TypeOf((*DefinitionTestT1)(nil)))
	t.Equal(def.TypeName, "github.com/bronze1man/kmg/dependencyInjection.DefinitionTestT1")

	cb := NewContainerBuilder()
	cb.MustSet("A", 9, "")
	cb.MustSet("github.com/bronze1man/kmg/dependencyInjection.DefinitionTestT2",
		&DefinitionTestT2{A: 16}, "")
	cb.MustSetDefinition(&Definition{
		Id: "C",
		Factory: func(c *Container) (interface{}, error) {
			return 15, nil
		},
		Type:  16,
		Scope: ScopeRequest,
	})
	cb.MustSetDefinition(&Definition{
		Factory: func(c *Container) (interface{}, error) {
			return &DefinitionTestT3{A: 17}, nil
		},
		TypeReflect: reflect.TypeOf((*DefinitionTestI1)(nil)).Elem(),
	})
	c = cb.MustCompile()
	ret, err := def.GetInst(c)
	t.Equal(err, nil)
	t.Equal(ret.(*DefinitionTestT1).A, 9)
	t.Equal(ret.(*DefinitionTestT1).B.A, 16)
	t.Equal(ret.(*DefinitionTestT1).C, 0) // not active scope field will not set
	t.Equal(ret.(*DefinitionTestT1).D.GetNum(), 17)
}
