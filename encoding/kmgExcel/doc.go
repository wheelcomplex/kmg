/*
excel read and write package
Grid: []map[string]string  same thing like mysql data struct
Array: [][]string     excel data struct, Data should be in this format to read from excel or write to excel
TitleArray: excel data struct with first row as title(map key) of following data.

Read from xlsx:
    XlsxFile2Array
    XlsxFileSheetIndex2Array
    XlsxFileFirstSheet2ArrayTrim
    UnmarshalFile
	UnmarshalFileSkipLine

Write to xlsx:
    Array2XlsxFile
    Array2XlsxIo

Tools:
   TitleArrayToGrid
   Trim2DArray

TODO find and remove all unnecessary function.
*/

package kmgExcel
