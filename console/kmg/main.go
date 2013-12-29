package main

import (
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/console/command"
)

func main() {
	manager := console.NewManager()
	manager.MustAdd(&command.GoFmt{})
	manager.MustAdd(&command.GoRun{})
	manager.MustAdd(&command.GoBuild{})
	manager.MustAdd(&command.GoWatch{})
	manager.MustAdd(&command.GoCrossComplieInit{})
	manager.MustAdd(&command.GoCrossComplie{})
	manager.MustAdd(&command.WatchCmd{})
	manager.MustAdd(&command.CurrentDir{})
	manager.MustAdd(&command.Yaml2Json{})
	manager.MustAdd(&command.Xlsx2Yaml{})
	manager.MustAdd(&command.ParpareReflect{})
	manager.MustAdd(&command.GoTest{})
	manager.ExecuteGlobal()
}
