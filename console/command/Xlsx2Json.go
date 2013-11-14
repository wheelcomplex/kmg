package command

import (
	"encoding/json"
	"flag"
	"github.com/bronze1man/kmg/console"
	"github.com/tealeg/xlsx"
	"github.com/bronze1man/kmg/errors"
	"strings"
)

type Xlsx2Json struct {
	filePath *string
	format *string
}

func (command *Xlsx2Json) GetNameConfig() *console.NameConfig {
	return &console.NameConfig{Name: "Xlsx2Json",
		Short: "convert from xlsx(Microsoft excel 2007) to json",
	}
}
func (command *Xlsx2Json) ConfigFlagSet(f *flag.FlagSet) {
	command.filePath = f.String("input", "", "input file path")
	command.format = f.String("format","grid","output json format(grid,raw)")
}
func (command *Xlsx2Json) Execute(context *console.Context) error {
	if *command.filePath == ""{
		return errors.New("need input file")
	}
	file, err := xlsx.OpenFile(*command.filePath)
	if err != nil {
		return err
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
		if len(s)==0{
			continue
		}
		output = append(output, s)
	}
	var outByte []byte
	switch *command.format{
	case "raw":
		outByte, err = json.Marshal(output)
		if err != nil {
			return err
		}
	case "grid":
		o:=[][]map[string]string{}
		for _,s:=range output{
			o1,err:=titleArrayToGrid(s)
			if err!=nil{
				return err
			}
			o = append(o,o1)
		}
		outByte, err = json.Marshal(o)
		if err != nil {
			return err
		}
	default:
		return errors.Sprintf("not support output format: %s",command.format)
	}
	_, err = context.Stdout.Write(outByte)
	if err != nil {
		return err
	}
	return nil
}
func titleArrayToGrid(titleArray [][]string)(output []map[string]string,err error){
	lenTitleArray := len(titleArray)
	if lenTitleArray <=1{
		return []map[string]string{},nil
	}
	output = make([]map[string]string,lenTitleArray-1)
	titles := titleArray[0]
	titles = trimLeftRowString(titles)
	lenTitles := len(titles)
	for rowIndex:=1;rowIndex<lenTitleArray;rowIndex++{
		row:=titleArray[rowIndex]
		//ignore blank cell left of row
		row = trimLeftRowString(row)
		//ignore blank line
		if len(row)==0{
			continue
		}
		r:=make(map[string]string)
		lenRow := len(row)
		//there can be some blank row of left of row
		if lenRow>lenTitles{
			return nil,errors.Sprintf(
				"titleArrayToGrid not all row length less or equal than first row length,"+
				"rowIndex: %d thisRowLen:%d firstRowLen:%d",rowIndex,lenRow,lenTitles)
		}
		for cellIndex:=0;cellIndex<lenRow;cellIndex++{
			cell:=row[cellIndex]
			r[titles[cellIndex]] = cell
		}
		output[rowIndex-1] = r
	}
	return output,nil
}

func trimLeftRowString(row []string)([]string){
	for i:=len(row)-1;i>=0;i--{
		if strings.Trim(row[i]," ")!=""{
			return row[:i]
		}
	}
	return []string{}
}
