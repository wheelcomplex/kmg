package kmgBundle

import (
	"database/sql"
	"github.com/bronze1man/kmg/ajkApi"
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

	//kmgSql
	c.MustSetDefinition(&dependencyInjection.Definition{
		Id: "kmgSql.godb",
		Factory: func(c *dependencyInjection.Container) (interface{}, error) {
			return sql.Open("mysql", "root:cd32d5e86e@tcp(git2.lqv.ca:3306)/monster_dev?charset=utf8")
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
	return nil
}
