package main

import (
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/console/command"
)

func main() {
	manager := console.NewManager()
	manager.MustAdd(&command.Fmt{})
	manager.MustAdd(&command.Run{})
	manager.MustAdd(&command.Watch{})
	manager.MustAdd(&command.WatchCmd{})
	manager.MustAdd(&command.CurrentDir{})
	manager.MustAdd(&command.Yaml2Json{})
	manager.MustAdd(&command.Xlsx2Json{})
	manager.ExecuteGlobal()

}
