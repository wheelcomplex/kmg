package kmgExcel

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"text/template"
)

//write raw data into xlsx file
//data key means: rowIndex,columnIndex,value
func Array2XlsxFile(data [][]string, path string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	err = Array2XlsxIo(data, f)
	if err != nil {
		return
	}
	return
}

//write raw data into a io.Writer
//data key means: rowIndex,columnIndex,value
//TODO work with reader...
func Array2XlsxIo(data [][]string, w io.Writer) (err error) {
	zw := zip.NewWriter(w)
	defer zw.Close()

	// sharedStrings
	err = array2XlsxWriteSharedStrings(zw, data)
	if err != nil {
		return
	}
	// sheel1
	err = array2XlsxWriteSheet1(zw, data)
	if err != nil {
		return
	}

	for filename, content := range fixedFileContent {
		thisW, err := zw.Create(filename)
		if err != nil {
			return err
		}
		_, err = thisW.Write(content)
		if err != nil {
			return err
		}
	}

	return
}
func array2XlsxWriteSharedStrings(zw *zip.Writer, data [][]string) (err error) {
	siList := []xlsxSharedStringSi{xlsxSharedStringSi{""}}
	for _, row := range data {
		for _, v1 := range row {
			if v1 == "" { //ignore blank cell can save space
				continue
			}
			siList = append(siList, xlsxSharedStringSi{v1})
		}

	}
	sst := xlsxSharedStringSst{
		Xmlns:  xmlNs,
		Count:  len(siList),
		SiList: siList,
	}
	thisW, err := zw.Create(sharedStringsFileName)
	_, err = thisW.Write([]byte(xml.Header))
	if err != nil {
		return
	}
	encoder := xml.NewEncoder(thisW)
	err = encoder.Encode(sst)
	if err != nil {
		return
	}
	return
}

func array2XlsxWriteSheet1(zw *zip.Writer, data [][]string) (err error) {
	rowList := make([]xlsxRow, len(data))
	totalIndex := 1
	MaxCellIndex := 0
	for rowIndex, row := range data {
		rowList[rowIndex].C = make([]xlsxC, len(row))
		if len(row) > MaxCellIndex {
			MaxCellIndex = len(row) - 1
		}
		for cellIndex, v1 := range row {
			index := totalIndex
			if v1 == "" { //ignore blank cell can save space
				index = 0
			}
			rowList[rowIndex].C[cellIndex] = xlsxC{
				R: CoordinateXy2Excel(cellIndex, rowIndex),
				T: "s",
				V: index,
			}
			totalIndex++
		}
	}
	sheetData := xlsxSheetData{
		Row: rowList,
	}
	thisW, err := zw.Create(sheel1FileName)
	xmlBytes, err := xml.Marshal(sheetData)
	if err != nil {
		return
	}
	err = sheelTpl.Execute(thisW, struct {
		MaxPosition string
		SheetData   string
	}{
		MaxPosition: CoordinateXy2Excel(MaxCellIndex, len(data)-1),
		SheetData:   string(xmlBytes),
	})
	if err != nil {
		return
	}
	return
}

//坐标系变换,从xy坐标系变化成excel的字符坐标系
//xy坐标从(0,0)开始,excel坐标从A1开始
func CoordinateXy2Excel(collomnIndex int, rowIndex int) (output string) {
	if collomnIndex < 0 || rowIndex < 0 {
		panic(fmt.Errorf("[CoordinateXy2Excel] collomnIndex[%d]<0||y[%d]<0", collomnIndex, rowIndex))
	}
	for {
		output = string(collomnIndex%26+int('A')) + output
		collomnIndex = collomnIndex/26 - 1
		if collomnIndex < 0 {
			break
		}
	}
	/*
		for reference,通过下面的代码推导出上面的循环代码
		if ((x/26-1)/26-1)>=0{
		    output+=string((x/26-1)/26-1+int('A'))+string((x/26-1)%26+int('A'))+string(x%26+int('A'))
		}else if (x/26-1)>=0{
			output+=string(x/26-1+int('A'))+string(x%26+int('A'))
		}else{
			output+=string(x+int('A'))
		}
	*/
	output += strconv.Itoa(rowIndex + 1)
	return output
}

//file content from wps xlsx file.
var fixedFileContent = map[string][]byte{
	"[Content_Types].xml": []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default ContentType="application/vnd.openxmlformats-package.relationships+xml" Extension="rels"/><Default ContentType="application/xml" Extension="xml"/><Override ContentType="application/vnd.openxmlformats-officedocument.extended-properties+xml" PartName="/docProps/app.xml"/><Override ContentType="application/vnd.openxmlformats-package.core-properties+xml" PartName="/docProps/core.xml"/><Override ContentType="application/vnd.openxmlformats-officedocument.custom-properties+xml" PartName="/docProps/custom.xml"/><Override ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sharedStrings+xml" PartName="/xl/sharedStrings.xml"/><Override ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml" PartName="/xl/styles.xml"/><Override ContentType="application/vnd.openxmlformats-officedocument.theme+xml" PartName="/xl/theme/theme1.xml"/><Override ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml" PartName="/xl/workbook.xml"/><Override ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml" PartName="/xl/worksheets/sheet1.xml"/></Types>`),
	"_rels/.rels": []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" Target="docProps/app.xml"/><Relationship Id="rId3" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/><Relationship Id="rId4" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/custom-properties" Target="docProps/custom.xml"/></Relationships>`),
	"docProps/app.xml": []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties" xmlns:vt="http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes"><Application>WPS Office 个人版</Application><HeadingPairs><vt:vector size="2" baseType="variant"><vt:variant><vt:lpstr>工作表</vt:lpstr></vt:variant><vt:variant><vt:i4>1</vt:i4></vt:variant></vt:vector></HeadingPairs><TitlesOfParts><vt:vector size="1" baseType="lpstr"><vt:lpstr>1</vt:lpstr></vt:vector></TitlesOfParts></Properties>`),
	"docProps/core.xml": []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:dcmitype="http://purl.org/dc/dcmitype/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><dc:creator>Administrator</dc:creator><dcterms:created xsi:type="dcterms:W3CDTF">2014-04-01T19:29:42Z</dcterms:created><dcterms:modified xsi:type="dcterms:W3CDTF">2014-04-01T19:29:50Z</dcterms:modified></cp:coreProperties>`),
	"docProps/custom.xml": []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/custom-properties" xmlns:vt="http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes"><property fmtid="{D5CDD505-2E9C-101B-9397-08002B2CF9AE}" pid="2" name="KSOProductBuildVer"><vt:lpwstr>2052-9.1.0.4468</vt:lpwstr></property></Properties>`),
	"xl/_rels/workbook.xml.rels": []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/><Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="theme/theme1.xml"/><Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/><Relationship Id="rId4" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/sharedStrings" Target="sharedStrings.xml"/></Relationships>`),
	"xl/theme/theme1.xml": []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<a:theme xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" name="Office 主题​​"><a:themeElements><a:clrScheme name="Office"><a:dk1><a:sysClr val="windowText" lastClr="000000"/></a:dk1><a:lt1><a:sysClr val="window" lastClr="FFFFFF"/></a:lt1><a:dk2><a:srgbClr val="1F497D"/></a:dk2><a:lt2><a:srgbClr val="EEECE1"/></a:lt2><a:accent1><a:srgbClr val="4F81BD"/></a:accent1><a:accent2><a:srgbClr val="C0504D"/></a:accent2><a:accent3><a:srgbClr val="9BBB59"/></a:accent3><a:accent4><a:srgbClr val="8064A2"/></a:accent4><a:accent5><a:srgbClr val="4BACC6"/></a:accent5><a:accent6><a:srgbClr val="F79646"/></a:accent6><a:hlink><a:srgbClr val="0000FF"/></a:hlink><a:folHlink><a:srgbClr val="800080"/></a:folHlink></a:clrScheme><a:fontScheme name="Office"><a:majorFont><a:latin typeface="Cambria"/><a:ea typeface=""/><a:cs typeface=""/><a:font script="Arab" typeface="Times New Roman"/><a:font script="Beng" typeface="Vrinda"/><a:font script="Cans" typeface="Euphemia"/><a:font script="Cher" typeface="Plantagenet Cherokee"/><a:font script="Deva" typeface="Mangal"/><a:font script="Ethi" typeface="Nyala"/><a:font script="Geor" typeface="Sylfaen"/><a:font script="Gujr" typeface="Shruti"/><a:font script="Guru" typeface="Raavi"/><a:font script="Hang" typeface="맑은 고딕"/><a:font script="Hans" typeface="宋体"/><a:font script="Hant" typeface="新細明體"/><a:font script="Hebr" typeface="Times New Roman"/><a:font script="Jpan" typeface="ＭＳ Ｐゴシック"/><a:font script="Khmr" typeface="MoolBoran"/><a:font script="Knda" typeface="Tunga"/><a:font script="Laoo" typeface="DokChampa"/><a:font script="Mlym" typeface="Kartika"/><a:font script="Mong" typeface="Mongolian Baiti"/><a:font script="Orya" typeface="Kalinga"/><a:font script="Sinh" typeface="Iskoola Pota"/><a:font script="Syrc" typeface="Estrangelo Edessa"/><a:font script="Taml" typeface="Latha"/><a:font script="Telu" typeface="Gautami"/><a:font script="Thaa" typeface="MV Boli"/><a:font script="Thai" typeface="Tahoma"/><a:font script="Tibt" typeface="Microsoft Himalaya"/><a:font script="Uigh" typeface="Microsoft Uighur"/><a:font script="Viet" typeface="Times New Roman"/><a:font script="Yiii" typeface="Microsoft Yi Baiti"/></a:majorFont><a:minorFont><a:latin typeface="Calibri"/><a:ea typeface=""/><a:cs typeface=""/><a:font script="Arab" typeface="Arial"/><a:font script="Beng" typeface="Vrinda"/><a:font script="Cans" typeface="Euphemia"/><a:font script="Cher" typeface="Plantagenet Cherokee"/><a:font script="Deva" typeface="Mangal"/><a:font script="Ethi" typeface="Nyala"/><a:font script="Geor" typeface="Sylfaen"/><a:font script="Gujr" typeface="Shruti"/><a:font script="Guru" typeface="Raavi"/><a:font script="Hang" typeface="맑은 고딕"/><a:font script="Hans" typeface="宋体"/><a:font script="Hant" typeface="新細明體"/><a:font script="Hebr" typeface="Arial"/><a:font script="Jpan" typeface="ＭＳ Ｐゴシック"/><a:font script="Khmr" typeface="DaunPenh"/><a:font script="Knda" typeface="Tunga"/><a:font script="Laoo" typeface="DokChampa"/><a:font script="Mlym" typeface="Kartika"/><a:font script="Mong" typeface="Mongolian Baiti"/><a:font script="Orya" typeface="Kalinga"/><a:font script="Sinh" typeface="Iskoola Pota"/><a:font script="Syrc" typeface="Estrangelo Edessa"/><a:font script="Taml" typeface="Latha"/><a:font script="Telu" typeface="Gautami"/><a:font script="Thaa" typeface="MV Boli"/><a:font script="Thai" typeface="Tahoma"/><a:font script="Tibt" typeface="Microsoft Himalaya"/><a:font script="Uigh" typeface="Microsoft Uighur"/><a:font script="Viet" typeface="Arial"/><a:font script="Yiii" typeface="Microsoft Yi Baiti"/></a:minorFont></a:fontScheme><a:fmtScheme name="Office"><a:fillStyleLst><a:solidFill><a:schemeClr val="phClr"/></a:solidFill><a:gradFill rotWithShape="1"><a:gsLst><a:gs pos="0"><a:schemeClr val="phClr"><a:tint val="50000"/><a:satMod val="300000"/></a:schemeClr></a:gs><a:gs pos="35000"><a:schemeClr val="phClr"><a:tint val="37000"/><a:satMod val="300000"/></a:schemeClr></a:gs><a:gs pos="100000"><a:schemeClr val="phClr"><a:tint val="15000"/><a:satMod val="350000"/></a:schemeClr></a:gs></a:gsLst><a:lin ang="16200000" scaled="1"/></a:gradFill><a:gradFill rotWithShape="1"><a:gsLst><a:gs pos="0"><a:schemeClr val="phClr"><a:shade val="51000"/><a:satMod val="130000"/></a:schemeClr></a:gs><a:gs pos="80000"><a:schemeClr val="phClr"><a:shade val="93000"/><a:satMod val="130000"/></a:schemeClr></a:gs><a:gs pos="100000"><a:schemeClr val="phClr"><a:shade val="94000"/><a:satMod val="135000"/></a:schemeClr></a:gs></a:gsLst><a:lin ang="16200000" scaled="0"/></a:gradFill></a:fillStyleLst><a:lnStyleLst><a:ln w="9525" cap="flat" cmpd="sng" algn="ctr"><a:solidFill><a:schemeClr val="phClr"><a:shade val="95000"/><a:satMod val="105000"/></a:schemeClr></a:solidFill><a:prstDash val="solid"/></a:ln><a:ln w="25400" cap="flat" cmpd="sng" algn="ctr"><a:solidFill><a:schemeClr val="phClr"/></a:solidFill><a:prstDash val="solid"/></a:ln><a:ln w="38100" cap="flat" cmpd="sng" algn="ctr"><a:solidFill><a:schemeClr val="phClr"/></a:solidFill><a:prstDash val="solid"/></a:ln></a:lnStyleLst><a:effectStyleLst><a:effectStyle><a:effectLst><a:outerShdw blurRad="40000" dist="20000" dir="5400000" rotWithShape="0"><a:srgbClr val="000000"><a:alpha val="38000"/></a:srgbClr></a:outerShdw></a:effectLst></a:effectStyle><a:effectStyle><a:effectLst><a:outerShdw blurRad="40000" dist="23000" dir="5400000" rotWithShape="0"><a:srgbClr val="000000"><a:alpha val="35000"/></a:srgbClr></a:outerShdw></a:effectLst></a:effectStyle><a:effectStyle><a:effectLst><a:outerShdw blurRad="40000" dist="23000" dir="5400000" rotWithShape="0"><a:srgbClr val="000000"><a:alpha val="35000"/></a:srgbClr></a:outerShdw></a:effectLst><a:scene3d><a:camera prst="orthographicFront"><a:rot lat="0" lon="0" rev="0"/></a:camera><a:lightRig rig="threePt" dir="t"><a:rot lat="0" lon="0" rev="1200000"/></a:lightRig></a:scene3d><a:sp3d><a:bevelT w="63500" h="25400"/></a:sp3d></a:effectStyle></a:effectStyleLst><a:bgFillStyleLst><a:solidFill><a:schemeClr val="phClr"/></a:solidFill><a:gradFill rotWithShape="1"><a:gsLst><a:gs pos="0"><a:schemeClr val="phClr"><a:tint val="40000"/><a:satMod val="350000"/></a:schemeClr></a:gs><a:gs pos="40000"><a:schemeClr val="phClr"><a:tint val="45000"/><a:satMod val="350000"/><a:shade val="99000"/></a:schemeClr></a:gs><a:gs pos="100000"><a:schemeClr val="phClr"><a:shade val="20000"/><a:satMod val="255000"/></a:schemeClr></a:gs></a:gsLst><a:path path="circle"><a:fillToRect l="50000" t="-80000" r="50000" b="180000"/></a:path></a:gradFill><a:gradFill rotWithShape="1"><a:gsLst><a:gs pos="0"><a:schemeClr val="phClr"><a:tint val="80000"/><a:satMod val="300000"/></a:schemeClr></a:gs><a:gs pos="100000"><a:schemeClr val="phClr"><a:shade val="30000"/><a:satMod val="200000"/></a:schemeClr></a:gs></a:gsLst><a:path path="circle"><a:fillToRect l="50000" t="50000" r="50000" b="50000"/></a:path></a:gradFill></a:bgFillStyleLst></a:fmtScheme></a:themeElements><a:objectDefaults><a:spDef><a:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="0" cy="0"/></a:xfrm><a:custGeom><a:avLst/><a:gdLst><a:gd name="_h" fmla="val 21600"/><a:gd name="_w" fmla="val 21600"/></a:gdLst><a:ahLst/><a:cxnLst/><a:pathLst><a:path w="21600" h="21600"/></a:pathLst></a:custGeom><a:gradFill rotWithShape="0"><a:gsLst><a:gs pos="100000"><a:srgbClr val="9CBEE0"/></a:gs><a:gs pos="0"><a:srgbClr val="BBD5F0"/></a:gs></a:gsLst><a:lin ang="5400000" scaled="0"/></a:gradFill><a:ln w="15875" cap="flat" cmpd="sng" algn="ctr"><a:solidFill><a:srgbClr val="739CC3"/></a:solidFill><a:prstDash val="solid"/><a:miter lim="200000"/></a:ln></a:spPr><a:bodyPr/><a:lstStyle/></a:spDef></a:objectDefaults><a:extraClrSchemeLst/></a:theme>`),
	"xl/styles.xml": []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><numFmts count="4"><numFmt numFmtId="43" formatCode="_ * #,##0.00_ ;_ * \-#,##0.00_ ;_ * &quot;-&quot;??_ ;_ @_ "/><numFmt numFmtId="44" formatCode="_ &quot;￥&quot;* #,##0.00_ ;_ &quot;￥&quot;* \-#,##0.00_ ;_ &quot;￥&quot;* &quot;-&quot;??_ ;_ @_ "/><numFmt numFmtId="41" formatCode="_ * #,##0_ ;_ * \-#,##0_ ;_ * &quot;-&quot;_ ;_ @_ "/><numFmt numFmtId="42" formatCode="_ &quot;￥&quot;* #,##0_ ;_ &quot;￥&quot;* \-#,##0_ ;_ &quot;￥&quot;* &quot;-&quot;_ ;_ @_ "/></numFmts><fonts count="1"><font><sz val="12"/><name val="宋体"/><charset val="134"/></font></fonts><fills count="2"><fill><patternFill patternType="none"/></fill><fill><patternFill patternType="gray125"/></fill></fills><borders count="1"><border><left/><right/><top/><bottom/><diagonal/></border></borders><cellStyleXfs count="6"><xf numFmtId="0" fontId="0" fillId="0" borderId="0"><alignment vertical="center"/></xf><xf numFmtId="43" fontId="0" fillId="0" borderId="0" applyFont="0" applyFill="0" applyBorder="0" applyAlignment="0" applyProtection="0"><alignment vertical="center"/></xf><xf numFmtId="44" fontId="0" fillId="0" borderId="0" applyFont="0" applyFill="0" applyBorder="0" applyAlignment="0" applyProtection="0"><alignment vertical="center"/></xf><xf numFmtId="41" fontId="0" fillId="0" borderId="0" applyFont="0" applyFill="0" applyBorder="0" applyAlignment="0" applyProtection="0"><alignment vertical="center"/></xf><xf numFmtId="9" fontId="0" fillId="0" borderId="0" applyFont="0" applyFill="0" applyBorder="0" applyAlignment="0" applyProtection="0"><alignment vertical="center"/></xf><xf numFmtId="42" fontId="0" fillId="0" borderId="0" applyFont="0" applyFill="0" applyBorder="0" applyAlignment="0" applyProtection="0"><alignment vertical="center"/></xf></cellStyleXfs><cellXfs count="1"><xf numFmtId="0" fontId="0" fillId="0" borderId="0" xfId="0"><alignment vertical="center"/></xf></cellXfs><cellStyles count="6"><cellStyle name="常规" xfId="0" builtinId="0"/><cellStyle name="千位分隔" xfId="1" builtinId="3"/><cellStyle name="货币" xfId="2" builtinId="4"/><cellStyle name="千位分隔[0]" xfId="3" builtinId="6"/><cellStyle name="百分比" xfId="4" builtinId="5"/><cellStyle name="货币[0]" xfId="5" builtinId="7"/></cellStyles></styleSheet>`),
	"xl/workbook.xml": []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><fileVersion appName="xl" lastEdited="3" lowestEdited="5" rupBuild="9302"/><workbookPr/><bookViews><workbookView windowWidth="13400" windowHeight="14160"/></bookViews><sheets><sheet name="1" sheetId="1" r:id="rId1"/></sheets><calcPr calcId="144525"/></workbook>`),
}

const sheel1FileName = "xl/worksheets/sheet1.xml"
const sharedStringsFileName = "xl/sharedStrings.xml"

var sheelTpl = template.Must(template.New("main").Parse(`<?xml version="1.0" encoding="UTF-8"?>
<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"
xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
xmlns:xdr="http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing"
xmlns:x14="http://schemas.microsoft.com/office/spreadsheetml/2009/9/main"
xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006">
<dimension ref="A1:{{.MaxPosition}}"/>
{{.SheetData}}
</worksheet>
`))

type xlsxSharedStringSst struct {
	XMLName xml.Name             `xml:"sst"`
	Xmlns   string               `xml:"xmlns,attr"`
	Count   int                  `xml:"count,attr"`
	SiList  []xlsxSharedStringSi `xml:"si"`
}
type xlsxSharedStringSi struct {
	T string `xml:"t"`
}
type xlsxSheetData struct {
	XMLName xml.Name  `xml:"sheetData"`
	Row     []xlsxRow `xml:"row"`
}
type xlsxRow struct {
	C []xlsxC `xml:"c"`
}
type xlsxC struct {
	R string `xml:"r,attr"`
	T string `xml:"t,attr"`
	V int    `xml:"v"`
}

const xmlNs = "http://schemas.openxmlformats.org/spreadsheetml/2006/main"

/*
for reference of contents of wps xlsx files which dynamically generated from this program
sharedStrings.xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<sst xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" count="2">
  <si><t>中文</t></si>
  <si><t>哈哈</t></si>
</sst>

sheet1.xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"
xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
xmlns:xdr="http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing"
xmlns:x14="http://schemas.microsoft.com/office/spreadsheetml/2009/9/main"
xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006">
<dimension ref="A1:C1492"/>
<sheetData>
<row>
  <c r="A1" t="s"><v>0</v></c>
  <c r="B1"><v>1</v></c>
</row>
<row r="2">
  <c r="A2"><v>2</v></c>
  <c r="B2"><v>3</v></c>
  <c r="C2" t="s"><v>1</v></c>
</row>
</sheetData>
</worksheet>
*/
