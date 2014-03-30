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
	manager.MustAdd(&command.GoCrossCompileInit{})
	manager.MustAdd(&command.GoCrossCompile{})
	manager.MustAdd(&command.GoTest{})
	manager.MustAdd(&command.GoAllImport{})

	manager.MustAdd(&command.CurrentDir{})
	manager.MustAdd(&command.CurrentPlatform{})

	manager.MustAdd(&command.WatchCmd{})
	manager.MustAdd(&command.Yaml2Json{})
	manager.MustAdd(&command.Json2Yaml{})
	manager.MustAdd(&command.Xlsx2Yaml{})
	manager.MustAdd(&command.ParpareReflect{})
	manager.ExecuteGlobal()
}
