package trans

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	luaJson "github.com/layeh/gopher-json"
	"github.com/yuin/gopher-lua"
)

var (
	conf    _Config
	outputs []Output

	statePool chan *lua.LState
)

func init() {
	size := runtime.NumCPU()
	statePool = make(chan *lua.LState, size)
	for i := 0; i < size; i++ {
		state := lua.NewState()
		luaJson.Preload(state)
		state.DoString("json = require('json')")
		statePool <- state
	}
}

func getLuaState() *lua.LState {
	return <-statePool
}

func releaseLuaState(s *lua.LState) {
	statePool <- s
}

// RegOutput 注册输出格式
func RegOutput(output Output) {
	outputs = append(outputs, output)
}

// Xlsx do
func Xlsx(xlsxDir string) {
	conf.XlsxDir = xlsxDir

	os.MkdirAll(conf.XlsxDir, 0644)
	for _, v := range outputs {
		os.MkdirAll(v.Dir(), 0644)
	}

	diff, _, err := checkFileMd5()
	if err != nil {
		return
	}
	if len(diff) == 0 {
		return
	}

	chans := make(chan _TransResult, len(diff))
	for _, mapping := range diff {
		go transOne(mapping, chans)
	}
	fmt.Println()
	var results []_TransResult
	for i := 0; i < len(diff); i++ {
		result := <-chans
		results = append(results, result)
		if result.err != nil {
			fmt.Printf("\n%s\n", result.mapping.CNFile)
			fmt.Printf("\t%s\n", result.err.Error())
		} else {
			fmt.Printf("\n%s\n", result.mapping.CNFile)
			printTime(20, "load", result.loadTime)
			printTime(20, "analyse", result.analyseTime)
			for _, output := range outputs {
				t, ok := result.marshalTime[output.Type()]
				if ok {
					printTime(20, "["+output.Type()+"]", t)
				}
			}
		}
	}
	fmt.Println()
	//writeFileMd5(all)
}

func transOne(mapping *_Mapping, chans chan<- _TransResult) {
	fmt.Printf("%s -> %s\n", mapping.CNFile, mapping.OutFile)
	var err error
	result := _TransResult{mapping: mapping, marshalTime: make(map[string]int)}
	defer func() {
		result.err = err
		chans <- result
	}()

	begin := time.Now().UnixNano()
	heads, rawDatas, err := getXlsxData(path.Join(conf.XlsxDir, mapping.XlsxFile))
	if err != nil {
		return
	}
	result.loadTime = int(time.Now().UnixNano() - begin)

	begin = time.Now().UnixNano()
	rows := sortKeys(rawDatas)
	srv, cli, err := analyseDatas(mapping.OutFile, rawDatas, heads)
	if err != nil {
		return
	}
	result.analyseTime = int(time.Now().UnixNano() - begin)

	for _, output := range outputs {
		begin = time.Now().UnixNano()
		outFlag := output.Flag()
		var datas XlsxRows
		if outFlag == FlagSrv {
			datas = srv
		} else if outFlag == FlagCli {
			datas = cli
		} else {
			datas = rawDatas
		}
		if datas == nil {
			continue
		}
		err = toFile(output, mapping, heads, datas, rows)
		result.marshalTime[output.Type()] = int(time.Now().UnixNano() - begin)
	}
}

func analyseDatas(file string, datas XlsxRows, heads []Head) (srv, cli XlsxRows, err error) {
	if heads[0].Flag == FlagDouble || heads[0].Flag == FlagSrv {
		srv = make(XlsxRows)
	}
	if heads[0].Flag == FlagDouble || heads[0].Flag == FlagCli {
		cli = make(XlsxRows)
	}
	for row, rawRow := range datas {
		srvRow := make(XlsxRow)
		cliRow := make(XlsxRow)
		for i := 0; i < len(heads); i++ {
			head := heads[i]
			v := rawRow[head.NameEN]
			switch head.Type {
			case "lua":
				if strings.TrimSpace(v.(string)) == "" {
					v = nil
				} else {
					val, e := parseLua(v)
					if e != nil {
						err = Errorf(file, head.NameEN, row, e.Error())
						return
					}
					v = val
				}
			case "json":
				if strings.TrimSpace(v.(string)) == "" {
					v = nil
				} else {
					var val interface{}
					err = json.Unmarshal(([]byte)(v.(string)), &val)
					if err != nil {
						err = Errorf(file, head.NameEN, row, err.Error())
						return
					}
					v = val
				}
			case "num", "number":
				if strings.TrimSpace(v.(string)) == "" {
					v = 0
				} else {
					val, e := strconv.ParseFloat(strings.TrimSpace(v.(string)), 64)
					if e != nil {
						err = Errorf(file, head.NameEN, row, e.Error())
						return
					}
					v = val
				}
			case "int":
				if strings.TrimSpace(v.(string)) == "" {
					v = 0
				} else {
					val, e := strconv.ParseInt(strings.TrimSpace(v.(string)), 0, 64)
					if e != nil {
						err = Errorf(file, head.NameEN, row, e.Error())
						return
					}
					v = val
				}
			case "uint":
				if strings.TrimSpace(v.(string)) == "" {
					v = 0
				} else {
					val, e := strconv.ParseUint(strings.TrimSpace(v.(string)), 0, 64)
					if e != nil {
						err = Errorf(file, head.NameEN, row, e.Error())
						return
					}
					v = val
				}
			}
			if srv != nil && (head.Flag == FlagDouble || head.Flag == FlagSrv) {
				asignPath(head.NameEN, srvRow, v)
			}
			if cli != nil && (head.Flag == FlagDouble || head.Flag == FlagCli) {
				asignPath(head.NameEN, cliRow, v)
			}
		}
		if srv != nil {
			srv[row] = srvRow
		}
		if cli != nil {
			cli[row] = cliRow
		}
	}
	return
}

func parseLua(v interface{}) (interface{}, error) {
	luaState := getLuaState()
	defer releaseLuaState(luaState)

	// luaState.SetGlobal("obj", lua.LString(v.(string)))
	if err := luaState.DoString("obj = " + v.(string)); err != nil {
		return nil, err
	}
	if err := luaState.DoString("return json.encode(obj)"); err != nil {
		return nil, err
	}
	lval := luaState.Get(luaState.GetTop())
	var jval interface{}
	err := json.Unmarshal(([]byte)(lval.String()), &jval)
	return jval, err
}

func toFile(output Output, mapping *_Mapping, heads []Head, datas XlsxRows, rows []int) (err error) {
	file, err := os.OpenFile(path.Join(output.Dir(), mapping.OutFile+"."+output.Type()), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	defer func() {
		if e := recover(); e != nil {
			switch v := e.(type) {
			case error:
				err = v
			default:
				err = fmt.Errorf("%v", v)
			}
		}
	}()
	output.Before(heads, mapping.XlsxFile, file)
	for _, row := range rows {
		output.Row(heads, datas[row], row, file)
	}
	output.After(heads, mapping.XlsxFile, file)
	return
}

func asignPath(fieldPath string, elem XlsxRow, v interface{}) {
	nameSl := strings.Split(fieldPath, ".")
	for k := 0; k < len(nameSl)-1; k++ {
		this := elem[nameSl[k]]
		if this == nil {
			elem[nameSl[k]] = make(map[string]interface{})
			this = elem[nameSl[k]]
		}
		elem = this.(map[string]interface{})
	}
	elem[nameSl[len(nameSl)-1]] = v
}

func buildLuaStr(str string) string {
	return "return json.encode(obj)"
}
