package kmgExcel

import (
	"fmt"
)

//transform from title array to grid
//titleArray should call Trim2DArray already
func TitleArrayToGrid(titleArray [][]string) (output []map[string]string, err error) {
	titleArray = Trim2DArray(titleArray)
	lenTitleArray := len(titleArray)
	if lenTitleArray <= 1 {
		return []map[string]string{}, nil
	}
	output = make([]map[string]string, lenTitleArray-1)
	titles := titleArray[0]
	lenTitles := len(titles)
	for rowIndex := 1; rowIndex < lenTitleArray; rowIndex++ {
		row := titleArray[rowIndex]
		r := make(map[string]string)
		lenRow := len(row)
		//there can be some blank row of left of row
		if lenRow > lenTitles {
			return nil, fmt.Errorf(
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

//transform from grid to title array
//keys not in title will return an error
// key in title but not in element map ,value will be "" in titleArray
// len(grid)==0 will return a titleArray only have title in it.
func GridToTitleArrayWithTitle(grid []map[string]string, title []string) (titleArray [][]string, err error) {
	if len(grid) == 0 {
		return [][]string{title}, nil
	}
	titleArray = make([][]string, len(grid)+1)
	lenTitle := len(title)
	titleIndexMap := make(map[string]int, lenTitle)
	titleArray[0] = title
	for i, key := range title {
		titleIndexMap[key] = i
	}
	for i, row := range grid {
		titleArray[i+1] = make([]string, lenTitle)
		for key, value := range row {
			j, exist := titleIndexMap[key]
			if !exist {
				return nil, fmt.Errorf("[GridToTitleArray][i:%d] key[%s] in grid Row,but not in title", i, key)
			}
			titleArray[i+1][j] = value
		}
	}
	return
}
