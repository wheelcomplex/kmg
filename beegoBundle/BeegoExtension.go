package beegoBundle

import (
	"github.com/astaxie/beego/orm"
	"github.com/bronze1man/kmg/dependencyInjection"
	"reflect"
	//"fmt"
	"time"
)

type BeegoExtension struct {
}

var HasRegisterDb bool

func (extension *BeegoExtension) LoadDependencyInjection(
	c *dependencyInjection.ContainerBuilder) error {
	if !HasRegisterDb {
		orm.RegisterDataBase("default", c.MustGetString("Parameter.databaseType"),
			c.MustGetString("Parameter.databaseDsn"))
		orm.SetDataBaseTZ("default", time.UTC)
	}
	HasRegisterDb = true

	c.MustSetDefinition(&dependencyInjection.Definition{
		Inst: &BeegoOrmSyncDbCommand{},
	}).AddTag("command")

	c.MustSetDefinition(&dependencyInjection.Definition{
		Inst: &BeegoOrmCreateDbCommand{},
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
