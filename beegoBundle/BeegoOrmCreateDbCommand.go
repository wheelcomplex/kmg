package beegoBundle

import (
	"database/sql"
	"github.com/bronze1man/kmg/console"
	_ "github.com/go-sql-driver/mysql"
)

type BeegoOrmCreateDbCommand struct {
}

func (command *BeegoOrmCreateDbCommand) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{
		Name:  "BeegoOrmCreateDb",
		Short: "beego orm create db",
	}
}

func (command *BeegoOrmCreateDbCommand) Execute(context *console.Context) (err error) {
	db, err := sql.Open("mysql", "root:cd32d5e86e@tcp(sxz10.lqv.ca:3306)/?charset=utf8&timeout=5s")
	if err != nil {
		return
	}
	_, err = db.Exec("create database `monster_test`")
	if err != nil {
		return
	}
	return
}
