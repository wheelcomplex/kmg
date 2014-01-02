package beegoBundle

import (
	"github.com/astaxie/beego/orm"
	"github.com/bronze1man/kmg/console"
	"os"
)

type BeegoOrmSyncDbCommand struct {
}

func (command *BeegoOrmSyncDbCommand) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{
		Name:  "BeegoOrmSyncDb",
		Short: "beego orm command",
	}
}
func (command *BeegoOrmSyncDbCommand) Execute(context *console.Context) error {
	os.Args = []string{
		os.Args[0], "orm", "syncdb",
	}
	orm.RunCommand()
	return nil
}
