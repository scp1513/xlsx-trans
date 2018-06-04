package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"github.com/scp1513/xlsx-trans"
)

type XlsxFiles map[string][]trans.XlsxRow

type JSONOutput struct {
	files XlsxFiles
	flag  trans.FlagType
	dir   string
}

func (j *JSONOutput) Type() string { return "json" }

func (j *JSONOutput) Flag() trans.FlagType { return j.flag }

func (j *JSONOutput) Dir() string { return j.dir }

func (j *JSONOutput) BeforeAll() {

}

func (j *JSONOutput) AfterAll() {

}

func (j *JSONOutput) Before(file string, heads []trans.Head, w io.Writer) {

}

func (j *JSONOutput) After(file string, heads []trans.Head, w io.Writer) {
	datas := j.files[file]
	if len(datas) == 0 {
		return
	}
	m := map[string]interface{}{"data": datas}
	b, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		panic(fmt.Errorf("[%s] json serialize error: %s", file, err.Error()))
	}
	w.Write(b)
	j.files[file] = nil
}

func (j *JSONOutput) Row(file string, heads []trans.Head, data trans.XlsxRow, row int, w io.Writer) {
	datas := j.files[file]
	datas = append(datas, data)
	j.files[file] = datas
}

type LUAOutput struct {
	files XlsxFiles
	flag  trans.FlagType
	dir   string
}

func (l *LUAOutput) Type() string { return "lua" }

func (l *LUAOutput) Flag() trans.FlagType { return l.flag }

func (l *LUAOutput) Dir() string { return l.dir }

func (l *LUAOutput) BeforeAll() {

}

func (l *LUAOutput) AfterAll() {

}

func (l *LUAOutput) Before(file string, heads []trans.Head, w io.Writer) {

}

func (l *LUAOutput) After(file string, heads []trans.Head, w io.Writer) {
	datas := l.files[file]
	if len(datas) == 0 {
		return
	}
	notes, err := genLuaNotes(heads)
	b, err := encLua(file, datas, heads)
	if err != nil {
		panic(fmt.Errorf("[%s] lua serialize error: %s", file, err.Error()))
	}
	w.Write(notes)
	w.Write(b)
	l.files[file] = nil
}

func (l *LUAOutput) Row(file string, heads []trans.Head, data trans.XlsxRow, row int, w io.Writer) {
	datas := l.files[file]
	datas = append(datas, data)
	l.files[file] = datas
}

func main() {
	var xlsxDir, jsonDir, luaDir string
	flag.StringVar(&xlsxDir, "i", "./xlsx", "")
	flag.StringVar(&jsonDir, "j", "./json", "")
	flag.StringVar(&luaDir, "l", "./lua", "")
	flag.Parse()

	trans.RegOutput(&JSONOutput{flag: trans.FlagSrv, dir: jsonDir, files: make(XlsxFiles)})
	trans.RegOutput(&LUAOutput{flag: trans.FlagCli, dir: luaDir, files: make(XlsxFiles)})

	trans.Xlsx(xlsxDir)
}
