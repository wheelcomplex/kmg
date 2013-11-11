package kmgBundle

import (
	"github.com/bronze1man/kmg/ajkApi"
	"github.com/bronze1man/kmg/dependencyInjection"
	"github.com/bronze1man/kmg/sessionStore"
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
	return nil
}
