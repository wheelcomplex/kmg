package kmgBundle

import (
	"database/sql"
	"github.com/bronze1man/kmg/ajkApi"
	"github.com/bronze1man/kmg/dependencyInjection"
	"github.com/bronze1man/kmg/kmgSql"
	"github.com/bronze1man/kmg/sessionStore"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
)

type KmgExtension struct {
}

func (extension *KmgExtension) LoadDependencyInjection(
	c *dependencyInjection.ContainerBuilder) error {
	//ajkapi
	c.MustSetDefinition(&dependencyInjection.Definition{
		TypeReflect: reflect.TypeOf((*ajkApi.ApiManagerInterface)(nil)).Elem(),
		Factory: func(c *dependencyInjection.Container) (interface{}, error) {
			return ajkApi.NewApiManagerFromContainer(c), nil
		},
	})
	c.MustSetDefinition(&dependencyInjection.Definition{
		Type:  (*ajkApi.Session)(nil),
		Scope: dependencyInjection.ScopeRequest,
	})
	c.MustSetDefinition(&dependencyInjection.Definition{
		Type: (*sessionStore.Store)(nil),
		Factory: func(c *dependencyInjection.Container) (interface{}, error) {
			session := c.MustGetByType((*ajkApi.Session)(nil)).(*ajkApi.Session)
			provider := c.MustGetByType((*sessionStore.Manager)(nil)).(*sessionStore.Manager)
			store, err := provider.LoadStoreOrNewIfNotExist(session.Guid)
			session.Guid = store.Guid()
			return store, err
		},
		Scope: dependencyInjection.ScopeRequest,
	})
	c.MustSetDefinition(&dependencyInjection.Definition{
		Type: (*ajkApi.JsonHttpHandler)(nil),
	})

	//sessionStore
	c.MustSetDefinition(&dependencyInjection.Definition{
		TypeReflect: reflect.TypeOf((*sessionStore.Provider)(nil)).Elem(),
		Inst:        sessionStore.NewMemoryProvider(),
	})
	c.MustSetDefinition(&dependencyInjection.Definition{
		Type: (*sessionStore.Manager)(nil),
	})

	databaseDsn := c.MustGetString("Parameter.databaseDsn")
	databaseType := c.MustGetString("Parameter.databaseType")
	//kmgSql
	c.MustSetDefinition(&dependencyInjection.Definition{
		Id: "kmgSql.godb",
		Factory: func(c *dependencyInjection.Container) (interface{}, error) {
			return sql.Open(databaseType, databaseDsn)
		},
	})

	c.MustSetDefinition(&dependencyInjection.Definition{
		Id: "kmgSql.db",
		Factory: func(c *dependencyInjection.Container) (interface{}, error) {
			return &kmgSql.Db{
				c.MustGet("kmgSql.godb").(*sql.DB),
			}, nil
		},
	})

	// build command
	c.MustSetDefinition(&dependencyInjection.Definition{
		Type: (*ajkApi.GoHttpApiServerCommand)(nil),
	}).AddTag("command")
	return nil
}
