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
		Type:  (*ajkApi.SessionIdContainer)(nil),
		Scope: dependencyInjection.ScopeRequest,
	})
	c.MustSetDefinition(&dependencyInjection.Definition{
		Type: (*sessionStore.Session)(nil),
		Factory: func(c *dependencyInjection.Container) (interface{}, error) {
			session := c.MustGetByType((*ajkApi.SessionIdContainer)(nil)).(*ajkApi.SessionIdContainer)
			provider := c.MustGetByType((*sessionStore.Manager)(nil)).(*sessionStore.Manager)
			store, err := provider.Load(session.SessionId)
			session.SessionId = store.Id
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

	c.MustSetDefinition(&dependencyInjection.Definition{
		Factory: func(c *dependencyInjection.Container) (interface{}, error) {
			DbConfig := &kmgSql.DbConfig{
				Username: c.MustGetString("Parameter.DatabaseUsername"),
				Password: c.MustGetString("Parameter.DatabasePassword"),
				Host:     c.MustGetString("Parameter.DatabaseHost"),
				DbName:   c.MustGetString("Parameter.DatabaseDbName"),
			}

			if c.MustGetString("Parameter.Env") == "test" {
				DbConfig.DbName = c.MustGetString("Parameter.DatabaseTestDbName")
			}
			return DbConfig, nil
		},
		Type: (*kmgSql.DbConfig)(nil),
	})
	c.MustSetDefinition(&dependencyInjection.Definition{
		Id: "kmgSql.DbConfig.Dsn",
		Factory: func(c *dependencyInjection.Container) (interface{}, error) {
			return c.MustGetByType((*kmgSql.DbConfig)(nil)).(*kmgSql.DbConfig).GetDsn(), nil
		},
	})
	//kmgSql
	c.MustSetDefinition(&dependencyInjection.Definition{
		Id: "kmgSql.godb",
		Factory: func(c *dependencyInjection.Container) (interface{}, error) {
			return sql.Open("mysql", c.MustGetString("kmgSql.DbConfig.Dsn"))
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
