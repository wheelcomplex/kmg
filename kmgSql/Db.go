package kmgSql

import (
	"database/sql"
)

//a wrap of database/sql.Db
type Db struct {
	*sql.DB
}
