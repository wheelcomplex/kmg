package kmgSql

import "fmt"

// get all data of a table
// mostly for test
func (db *Db) GetTableData(tableName string) (output []map[string]string, err error) {
	return db.QueryGrid(fmt.Sprintf("SELECT * FROM `%s`", tableName))
}
func (db *Db) MustGetTableData(tableName string) (output []map[string]string) {
	output, err := db.GetTableData(tableName)
	if err != nil {
		panic(err)
	}
	return
}
