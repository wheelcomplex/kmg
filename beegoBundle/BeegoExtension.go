package beegoBundle

import (
	"github.com/astaxie/beego/orm"
	"github.com/bronze1man/kmg/dependencyInjection"
	"reflect"
	"sync"
	"time"
)

var RegisterDb sync.Once

type BeegoExtension struct {
}

func (extension *BeegoExtension) LoadDependencyInjection(
	c *dependencyInjection.ContainerBuilder) error {
	c.MustSetDefinition(&dependencyInjection.Definition{
		Type: (*BeegoOrmSyncDbCommand)(nil),
	}).AddTag("command")

	c.MustSetDefinition(&dependencyInjection.Definition{
		Type: (*BeegoOrmCreateDbCommand)(nil),
	}).AddTag("command")

	c.MustSetDefinition(&dependencyInjection.Definition{
		TypeReflect: reflect.TypeOf((*orm.Ormer)(nil)).Elem(),
		Factory: func(c *dependencyInjection.Container) (interface{}, error) {
			RegisterDb.Do(func() {
				orm.RegisterDataBase("default", "mysql", c.MustGetString("kmgSql.DbConfig.Dsn"))
				orm.SetDataBaseTZ("default", time.UTC)
			})
			return orm.NewOrm(), nil
		},
		Scope: dependencyInjection.ScopeRequest,
	})
	return nil
}
