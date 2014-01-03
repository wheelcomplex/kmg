package excel

import (
	"github.com/bronze1man/kmg/typeTransform"
)

/*
 unmarshal some stuff from excel File
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
