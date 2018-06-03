package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"github.com/scp1513/xlsx-trans"
)

type JSONOutput struct {
	rows []trans.XlsxRow
	flag trans.FlagType
	dir  string
}

func (j *JSONOutput) Type() string { return "json" }

func (j *JSONOutput) Flag() trans.FlagType { return j.flag }

func (j *JSONOutput) Dir() string { return j.dir }

func (j *JSONOutput) BeforeAll() {

}

func (j *JSONOutput) AfterAll() {

}

func (j *JSONOutput) Before(heads []trans.Head, file string, w io.Writer) {

}

func (j *JSONOutput) After(heads []trans.Head, file string, w io.Writer) {
	if len(j.rows) == 0 {
		return
	}
	m := map[string]interface{}{"data": j.rows}
	b, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		panic(fmt.Errorf("[%s] json serialize error: %s", file, err.Error()))
	}
	w.Write(b)
	j.rows = nil
}

func (j *JSONOutput) Row(heads []trans.Head, data trans.XlsxRow, row int, w io.Writer) {
	j.rows = append(j.rows, data)
}

type LUAOutput struct {
	rows []trans.XlsxRow
	flag trans.FlagType
	dir  string
}

func (l *LUAOutput) Type() string { return "lua" }

func (l *LUAOutput) Flag() trans.FlagType { return l.flag }

func (l *LUAOutput) Dir() string { return l.dir }

func (l *LUAOutput) BeforeAll() {

}

func (l *LUAOutput) AfterAll() {

}

func (l *LUAOutput) Before(heads []trans.Head, file string, w io.Writer) {

}

func (l *LUAOutput) After(heads []trans.Head, file string, w io.Writer) {
	if len(l.rows) == 0 {
		return
	}
	notes, err := genLuaNotes(heads)
	b, err := encLua(file, l.rows, heads)
	if err != nil {
		panic(fmt.Errorf("[%s] lua serialize error: %s", file, err.Error()))
	}
	w.Write(notes)
	w.Write(b)
	l.rows = nil
}

func (l *LUAOutput) Row(heads []trans.Head, data trans.XlsxRow, row int, w io.Writer) {
	l.rows = append(l.rows, data)
}

func main() {
	var xlsxDir, jsonDir, luaDir string
	flag.StringVar(&xlsxDir, "i", "./xlsx", "")
	flag.StringVar(&jsonDir, "j", "./json", "")
	flag.StringVar(&luaDir, "l", "./lua", "")
	flag.Parse()

	trans.RegOutput(&JSONOutput{flag: trans.FlagSrv, dir: jsonDir})
	trans.RegOutput(&LUAOutput{flag: trans.FlagCli, dir: luaDir})

	trans.Xlsx(xlsxDir)
}
