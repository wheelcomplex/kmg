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
