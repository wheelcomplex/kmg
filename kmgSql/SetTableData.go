package kmgSql

import (
	"database/sql"
	"fmt"
	"github.com/bronze1man/kmg/errors"
	"launchpad.net/goyaml"
	"strings"
)

func (db *Db) MustSetTablesDataYaml(yaml string) {
	err := db.SetTablesDataYaml(yaml)
	if err != nil {
		panic(err)
	}
}
func (db *Db) SetTablesDataYaml(yaml string) (err error) {
	data := make(map[string][]map[string]string)
	err = goyaml.Unmarshal([]byte(yaml), &data)
	if err != nil {
		return err
	}
	return db.SetTablesData(data)
}

// Set some tables data in this database.
// It will delete all data in that table,and insert new data.
// mostly for test
func (db *Db) SetTablesData(data map[string][]map[string]string) (err error) {
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	err = setTablesDataTransaction(data, tx)
	if err != nil {
		errRoll := tx.Rollback()
		if errRoll != nil {
			return errors.Sprintf("error [transaction] %s,[rollback] %s", err, errRoll)
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
func setTablesDataTransaction(data map[string][]map[string]string, tx *sql.Tx) error {
	for tableName, tableData := range data {
		_, err := tx.Exec(fmt.Sprintf("truncate `%s`", tableName))
		if err != nil {
			return err
		}
		for _, row := range tableData {
			colNameList := []string{}
			placeHolderNum := len(row)
			valueList := []interface{}{}
			for name, value := range row {
				colNameList = append(colNameList, name)
				valueList = append(valueList, value)
			}
			sqlColNamePart := "`" + strings.Join(colNameList, "`, `") + "`"
			sqlValuePart := strings.Repeat("?, ", placeHolderNum-1) + "?"
			sql := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", tableName, sqlColNamePart, sqlValuePart)
			_, err := tx.Exec(sql, valueList...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
