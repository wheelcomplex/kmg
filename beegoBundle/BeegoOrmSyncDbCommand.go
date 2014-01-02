package beegoBundle

import (
	"flag"
	"github.com/astaxie/beego/orm"
	"github.com/bronze1man/kmg/console"
	"github.com/bronze1man/kmg/dependencyInjection"
	"os"
	"time"
)

type BeegoOrmSyncDbCommand struct {
	C   *dependencyInjection.Container
	env string
}

func (command *BeegoOrmSyncDbCommand) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{
		Name:  "BeegoOrmSyncDb",
		Short: "beego orm command",
	}
}

func (command *BeegoOrmSyncDbCommand) ConfigFlagSet(flag *flag.FlagSet) {
	flag.StringVar(&command.env, "env", "dev", "database env(dev,test)")
}
func (command *BeegoOrmSyncDbCommand) Execute(context *console.Context) error {
	os.Args = []string{
		os.Args[0], "orm", "syncdb",
	}
	command.C.MustSet("Parameter.Env", command.env, "")

	orm.RegisterDataBase("default", "mysql", command.C.MustGetString("kmgSql.DbConfig.Dsn"))
	orm.SetDataBaseTZ("default", time.UTC)

	orm.RunCommand()
	return nil
}
