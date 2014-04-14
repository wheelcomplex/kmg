package kmgSql

import (
	"fmt"
)

type TransactionableDb interface {
	Begin() error
	Commit() error
	Rollback() error
}

//transaction callback on beego.orm,but not depend on it
func TransactionCallback(db TransactionableDb, f func() error) error {
	hasFinish := false
	defer func() { //panic的时候处理
		if !hasFinish {
			db.Rollback()
		}
	}()
	err := db.Begin()
	if err != nil {
		return err
	}
	err = f()
	if err != nil {
		errR := db.Rollback()
		hasFinish = true
		if errR != nil {
			return fmt.Errorf("rollback fail:%s,origin fail:%s", errR.Error(), err.Error())
		}
		return err
	}
	err = db.Commit()
	if err != nil {
		return err
	}
	hasFinish = true
	return nil
}
