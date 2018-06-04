package trans

import (
	"strings"

	"github.com/tealeg/xlsx"
)

var validType = []string{
	"lua",
	"json",
	"num", "number", "int", "uint",
	"str", "string",
}

func getXlsxData(name string) (heads []Head, datas XlsxRows, err error) {
	f, err := xlsx.OpenFile(name + ".xlsx")
	if err != nil {
		return
	}
	sheet := f.Sheets[0]
	heads = getHeads(sheet)
	datas = make(XlsxRows)
	for row := 4; row < sheet.MaxRow; row++ {
		data := make(XlsxRow)
		emptyLine := false
		for i := 0; i < len(heads); i++ {
			col := heads[i].col
			val := sheet.Cell(row, col).Value
			trim := strings.TrimSpace(val)
			if i == 0 && (trim == "" || trim[0] == '#') {
				emptyLine = true
				break
			}
			data[heads[i].NameEN] = val
		}
		if !emptyLine {
			datas[row+1] = data
		}
	}
	return
}

func getHeads(sheet *xlsx.Sheet) (ret []Head) {
	for col := 0; col < 9999; col++ {
		val := getCellVal(sheet, 1, col, "")
		if val == "" {
			continue
		}
		flag := getCellVal(sheet, 3, col, FlagRead)
		if !isFlagValid(flag) {
			continue
		}
		typ := getCellVal(sheet, 2, col, "")
		typeValid := false
		for _, v := range validType {
			if v == typ {
				typeValid = true
				break
			}
		}
		if !typeValid {
			continue
		}
		ret = append(ret, Head{
			col:    col,
			NameCN: getCellVal(sheet, 0, col, ""),
			NameEN: val,
			Type:   typ,
			Flag:   flag,
		})
	}
	return
}

func getCellVal(sheet *xlsx.Sheet, row, col int, def string) string {
	str := strings.TrimSpace(sheet.Cell(row, col).Value)
	if str == "" {
		return def
	}
	return str
}
