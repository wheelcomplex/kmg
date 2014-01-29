package kmgBundle

import (
	"fmt"
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/dependencyInjection"
)

type CommandCompilePass struct {
}

func (extension *CommandCompilePass) CompilePass(
	c *dependencyInjection.ContainerBuilder) error {
	ids := []string{}
	for _, def := range c.GetTaggedDefinition("command") {
		ids = append(ids, def.Id)
	}
	//ajkapi
	c.MustSetDefinition(&dependencyInjection.Definition{
		Type: (*console.Manager)(nil),
		Factory: func(c *dependencyInjection.Container) (interface{}, error) {
			manager := console.NewManager()
			for _, id := range ids {
				obj, err := c.Get(id)
				if err != nil {
					return nil, err
				}
				command, ok := obj.(console.Command)
				if !ok {
					return nil, fmt.Errorf("service %s register as command but not implement command interface", id)
				}
				err = manager.Add(command)
				if err != nil {
					return nil, err
				}
			}
			return manager, nil
		},
	})
	return nil
}
