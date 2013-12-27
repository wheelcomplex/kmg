package excel

import (
	"github.com/tealeg/xlsx"
)

//output index mean=> sheet ,row ,cell ,value
func XlsxFile2Array(path string) ([][][]string, error) {
	file, err := xlsx.OpenFile(path)
	if err != nil {
		return nil, err
	}
	output := [][][]string{}
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
		//ignore empty Sheet
		if len(s) == 0 {
			continue
		}
		output = append(output, s)
	}
	return output, nil
}

//output index mean=> row ,cell ,value
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
