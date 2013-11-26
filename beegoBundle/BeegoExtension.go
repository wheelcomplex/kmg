package beegoBundle

import (
	"github.com/astaxie/beego/orm"
	"github.com/bronze1man/kmg/dependencyInjection"
	"reflect"
)

type BeegoExtension struct {
}

var HasRegisterDb bool

func (extension *BeegoExtension) LoadDependencyInjection(
	c *dependencyInjection.ContainerBuilder) error {
	if !HasRegisterDb {
		orm.RegisterDataBase("default", c.Parameters["databaseType"],
			c.Parameters["databaseDsn"])
	}
	HasRegisterDb = true

	c.MustSetDefinition(&dependencyInjection.Definition{
		Id:   "beego.command.orm",
		Inst: &OrmCommand{},
	}).AddTag("command")

	c.MustSetDefinition(&dependencyInjection.Definition{
		TypeReflect: reflect.TypeOf((*orm.Ormer)(nil)).Elem(),
		Factory: func(c *dependencyInjection.Container) (interface{}, error) {
			return orm.NewOrm(), nil
		},
		Scope: dependencyInjection.ScopeRequest,
	})
	return nil
}
