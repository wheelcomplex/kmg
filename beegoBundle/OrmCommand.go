package beegoBundle

import (
	"github.com/astaxie/beego/orm"
	"github.com/bronze1man/kmg/console"
	//"fmt"
	"os"
)

type OrmCommand struct {
}

func (command *OrmCommand) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "beego.orm",
		Short: "beego orm command"}
}
func (command *OrmCommand) Execute(context *console.Context) error {
	os.Args[1] = "orm"
	orm.RunCommand()
	return nil
}
