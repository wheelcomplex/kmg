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

func (db *Db) GetTableDataMap(tableName string, pkName string) (output map[string]map[string]string, err error) {
	grid, err := db.GetTableData(tableName)
	output = make(map[string]map[string]string)
	for _, row := range grid {
		pk, ok := row[pkName]
		if !ok {
			return nil, fmt.Errorf("GetTableDataMap: pkName colum not exits, tableName: %s, pkName: %s", tableName, pkName)
		}
		output[pk] = row
	}
	return
}

func (db *Db) MustGetTableDataMap(tableName string, pkName string) (output map[string]map[string]string) {
	output, err := db.GetTableDataMap(tableName, pkName)
	if err != nil {
		panic(err)
	}
	return
}
