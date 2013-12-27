package excel

import (
	"github.com/bronze1man/kmg/errors"
	"strings"
)

func TitleArrayToGrid(titleArray [][]string) (output []map[string]string, err error) {
	titleArray = Trim2DArray(titleArray)
	lenTitleArray := len(titleArray)
	if lenTitleArray <= 1 {
		return []map[string]string{}, nil
	}
	output = make([]map[string]string, lenTitleArray-1)
	titles := titleArray[0]
	//titles = trimRightRowString(titles)
	lenTitles := len(titles)
	for rowIndex := 1; rowIndex < lenTitleArray; rowIndex++ {
		row := titleArray[rowIndex]
		r := make(map[string]string)
		lenRow := len(row)
		//there can be some blank row of left of row
		if lenRow > lenTitles {
			return nil, errors.Sprintf(
				"titleArrayToGrid not all row length less or equal than first row length,"+
					"rowIndex: %d thisRowLen:%d firstRowLen:%d", rowIndex, lenRow, lenTitles)
		}
		for cellIndex := 0; cellIndex < lenRow; cellIndex++ {
			cell := row[cellIndex]
			r[titles[cellIndex]] = cell
		}
		output[rowIndex-1] = r
	}
	return output, nil
}

func trimRightRowString(row []string) []string {
	for i := len(row) - 1; i >= 0; i-- {
		if strings.Trim(row[i], " ") != "" {
			return row[:i+1]
		}
	}
	return []string{}
}
