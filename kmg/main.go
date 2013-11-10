package main

import (
	"github.com/bronze1man/kmg/buildCommand"
	"github.com/bronze1man/kmg/console"
)

func main() {
	manager := console.NewManager()
	manager.MustAdd(&buildCommand.FmtCommand{})
	manager.MustAdd(&buildCommand.RunCommand{})
	manager.MustAdd(&buildCommand.WatchCommand{})
	manager.ExecuteGlobal()

}
