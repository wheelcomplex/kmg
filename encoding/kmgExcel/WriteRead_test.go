package kmgExcel

import (
	"archive/zip"
	"bytes"
	"github.com/bronze1man/kmg/kmgTest"
	"github.com/tealeg/xlsx"
	"testing"
	//"fmt"
)

func TestWriteRead(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	buf := &bytes.Buffer{}
	inData := [][]string{
		[]string{"中文"},
		[]string{"1"},
	}
	err := Array2XlsxIo(inData, buf)
	t.Equal(err, nil)
	r := bytes.NewReader(buf.Bytes())
	zr, err := zip.NewReader(r, int64(buf.Len()))
	t.Equal(err, nil)
	xlsxFile, err := xlsx.ReadZipReader(zr)
	t.Equal(err, nil)
	outData, err := xlsx2ArrayXlsxFile(xlsxFile)
	t.Equal(err, nil)
	t.Equal(outData[0], inData)
}
