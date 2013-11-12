package beegoBundle

import (
	"github.com/astaxie/beego/orm"
	"github.com/bronze1man/kmg/dependencyInjection"
)

type BeegoExtension struct {
}

func (extension *BeegoExtension) LoadDependencyInjection(
	c *dependencyInjection.ContainerBuilder) error {

	orm.RegisterDataBase("default", c.Parameters["databaseType"],
		c.Parameters["databaseDsn"])

	c.MustSetDefinition(&dependencyInjection.Definition{
		Id:   "beego.command.orm",
		Inst: &OrmCommand{},
	}).AddTag("command")

	c.MustSetDefinition(&dependencyInjection.Definition{
		Id: "beego.orm",
		Factory: func(c *dependencyInjection.Container) (interface{}, error) {
			return orm.NewOrm(), nil
		},
		Scope: dependencyInjection.ScopePrototype,
	})
	return nil
}
