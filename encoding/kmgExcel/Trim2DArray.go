package kmgExcel

//remove empty cell right(every line will have same length)
//remove all empty line (include some empty line in between)
func Trim2DArray(input [][]string) (output [][]string) {
	if len(input) == 0 {
		return [][]string{}
	}
	//find max line cell number
	MaxLineCellNumber := len(TrimRightRowString(input[0]))
	LineNumber := 0
	for _, row := range input {
		thisLineCellNumber := len(TrimRightRowString(row))
		if thisLineCellNumber != 0 {
			LineNumber++
		}
		if MaxLineCellNumber < thisLineCellNumber {
			MaxLineCellNumber = thisLineCellNumber
		}
	}
	output = make([][]string, LineNumber)
	i := 0
	for _, row := range input {
		thisLineCellNumber := len(TrimRightRowString(row))
		if thisLineCellNumber == 0 {
			continue
		}
		if MaxLineCellNumber <= len(row) {
			output[i] = row[:MaxLineCellNumber]
		} else {
			retRow := make([]string, 0, MaxLineCellNumber)
			retRow = append(retRow, row...)
			remainCellNumber := MaxLineCellNumber - len(row)
			for i := 0; i < remainCellNumber; i++ {
				retRow = append(retRow, "")
			}
			output[i] = retRow
		}

		i++
	}
	return output
}
