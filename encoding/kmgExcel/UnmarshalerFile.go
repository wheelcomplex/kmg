package kmgExcel

import (
	"fmt"
	"github.com/bronze1man/kmg/typeTransform"
)

/*
 Unmarshal excel into a array of a struct

 * you must put data in first sheet
 * accept []struct{Key1 value;Key2 value} or []map[string]string
 * first excel row as struct key,remain excel row as struct value
 * type can not convert will fail
 * lack key will ignore
*/
func UnmarshalFile(filePath string, obj interface{}) (err error) {
	rawArray, err := XlsxFile2Array(filePath)
	if err != nil {
		return
	}
	gridArray, err := TitleArrayToGrid(rawArray[0])
	if err != nil {
		return
	}
	err = typeTransform.Transform(gridArray, obj)
	if err != nil {
		return
	}
	return
}

//Unmarshal excel into a array of a struct and skip some line from start
func UnmarshalFileSkipLine(filePath string, obj interface{}, skipline int) (err error) {
	rawArray, err := XlsxFile2Array(filePath)
	if err != nil {
		return
	}
	if len(rawArray[0]) < skipline+1 {
		return fmt.Errorf("[kmgExcel.UnmarshalFileSkipLine]filePath:%s len(rawArray[0])<skipline+1 len(rawArray[0]):%d skipline:%d",
			filePath, len(rawArray[0]), skipline+1)
	}
	gridArray, err := TitleArrayToGrid(rawArray[0][skipline:])
	if err != nil {
		return
	}
	err = typeTransform.Transform(gridArray, obj)
	if err != nil {
		return
	}
	return
}
