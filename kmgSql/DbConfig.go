package kmgSql

import (
	"fmt"
)

type DbConfig struct {
	Username string
	Password string
	Host     string
	DbName   string
}

func (config *DbConfig) GetDsn() string {
	return fmt.Sprintf("%s:%s@%s/%s?charset=utf8&timeout=5s",
		config.Username,
		config.Password,
		config.Host,
		config.DbName)
}
