package kmgExcel

import (
	"github.com/tealeg/xlsx"
	//"fmt"
)

// output index mean=> sheet ,row ,cell ,value
// not remove any cells
func XlsxFile2Array(path string) ([][][]string, error) {
	file, err := xlsx.OpenFile(path)
	if err != nil {
		return nil, err
	}
	output := [][][]string{}
	//fmt.Printf("%#v\n",file.Sheet["工作表1"])
	for _, sheet := range file.Sheets {
		s := [][]string{}
		for _, row := range sheet.Rows {
			if row == nil {
				continue
			}
			r := []string{}
			for _, cell := range row.Cells {
				r = append(r, cell.String())
			}
			s = append(s, r)
		}
		output = append(output, s)
	}
	return output, nil
}

//output index mean=> row ,cell ,value
// not remove any cells
func XlsxFileSheetIndex2Array(path string, index int) ([][]string, error) {
	file, err := xlsx.OpenFile(path)
	if err != nil {
		return nil, err
	}
	output := [][]string{}
	sheet := file.Sheets[index]
	for _, row := range sheet.Rows {
		if row == nil {
			continue
		}
		r := []string{}
		for _, cell := range row.Cells {
			r = append(r, cell.String())
		}
		output = append(output, r)
	}
	return output, nil
}

//output index mean=> row ,cell ,value
// remove all right and bottom black lines
func XlsxFileFirstSheet2ArrayTrim(path string) (output [][]string, err error) {
	output, err = XlsxFileSheetIndex2Array(path, 0)
	if err != nil {
		return
	}
	output = Trim2DArray(output)
	return
}
