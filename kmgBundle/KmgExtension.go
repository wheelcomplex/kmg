package kmgBundle

import (
	"database/sql"
	"github.com/bronze1man/kmg/ajkApi"
	"github.com/bronze1man/kmg/console/command"
	"github.com/bronze1man/kmg/dependencyInjection"
	"github.com/bronze1man/kmg/kmgSql"
	"github.com/bronze1man/kmg/sessionStore"
	_ "github.com/go-sql-driver/mysql"
)

type KmgExtension struct {
}

func (extension *KmgExtension) LoadDependencyInjection(
	c *dependencyInjection.ContainerBuilder) error {
	//ajkapi
	c.MustSetFactory("ajkApi.ApiManager", func(c *dependencyInjection.Container) (interface{}, error) {
		return ajkApi.NewApiManagerFromContainer(c), nil
	}, "")
	c.MustSetFactory("ajkApi.JsonHttpHandler", func(c *dependencyInjection.Container) (interface{}, error) {
		return &ajkApi.JsonHttpHandler{
			ApiManager:          c.MustGet("ajkApi.ApiManager").(ajkApi.ApiManagerInterface),
			SessionStoreManager: c.MustGet("sessionStore.Manager").(*sessionStore.Manager),
		}, nil
	}, "")

	//sessionStore
	c.MustSet("sessionStore.Provider", sessionStore.NewMemoryProvider(), "")
	c.MustSetFactory("sessionStore.Manager", func(c *dependencyInjection.Container) (interface{}, error) {
		return &sessionStore.Manager{
			c.MustGet("sessionStore.Provider").(sessionStore.Provider),
		}, nil
	}, "")

	databaseDsn := c.Parameters["databaseDsn"]
	databaseType := c.Parameters["databaseType"]
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
		Inst: &command.GoFmt{},
	}).AddTag("command")

	c.MustSetDefinition(&dependencyInjection.Definition{
		Inst: &command.GoRun{},
	}).AddTag("command")

	c.MustSetDefinition(&dependencyInjection.Definition{
		Inst: &command.WatchCmd{},
	}).AddTag("command")

	c.MustSetDefinition(&dependencyInjection.Definition{
		Inst: &command.GoWatch{},
	}).AddTag("command")
	return nil
}
