package main

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"

	"github.com/scp1513/xlsx-trans"
)

func encLua(fileName string, rows []trans.XlsxRow, head []trans.Head) ([]byte, error) {
	var buff bytes.Buffer
	buff.WriteString("local data = {\n")
	var count = 0
	for i := 0; i < len(rows); i++ {
		row := rows[i]
		writeTab(&buff, 1)
		writeNameAndEqual(&buff, count+1)
		buff.WriteString("{\n")
		//buff.WriteString("[" + strconv.FormatFloat(v[i]["id"].(float64), 'f', -1, 64) + "] = {\n")
		for j := 0; j < len(head); j++ {
			field, ok := row[head[j].NameEN]
			if !ok {
				continue
			}
			encVal(fileName, i+1, &buff, head[j].NameEN, field, 2)
		}
		writeTab(&buff, 1)
		buff.WriteString("},\n")
		count++
	}
	buff.WriteString("}\n\nreturn data")
	return buff.Bytes(), nil
}

// func encVal(fileName string, lineno int, buff *bytes.Buffer, name string, val interface{}, tab int) {
// 	writeTab(buff, tab)
// 	switch v := val.(type) {
// 	case []interface{}:
// 		writeNameAndEqual(buff, name)
// 		buff.WriteString("{\n")
// 		for i := 0; i < len(v); i++ {
// 			if vvv, ok := v[i].(map[string]interface{}); ok {
// 				//encVal(fileName, lineno, buff, name, vvv, tab+1)
// 				writeTab(buff, tab+1)
// 				writeNameAndEqual(buff, strconv.FormatInt(int64(i+1), 10))
// 				buff.WriteString("{\n")
// 				var keys []string
// 				for k := range vvv {
// 					keys = append(keys, k)
// 				}
// 				sort.Strings(keys)
// 				for j := 0; j < len(keys); j++ {
// 					encVal(fileName, lineno, buff, keys[j], vvv[keys[j]], tab+2)
// 				}
// 				writeTab(buff, tab+1)
// 				buff.WriteString("},\n")
// 			} else if vvv, ok := v[i].([]interface{}); ok {
// 				encVal(fileName, lineno, buff, strconv.FormatInt(int64(i+1), 10), vvv, tab+1)
// 			} else {
// 				writeTab(buff, tab+1)
// 				writeVal(fileName, lineno, buff, name, v[i])
// 				buff.WriteString(",\n")
// 			}
// 		}
// 		writeTab(buff, tab)
// 		buff.WriteString("},\n")
// 	case string, int, uint, int64, uint64, float32, float64, nil:
// 		writeNameAndEqual(buff, name)
// 		writeVal(fileName, lineno, buff, name, v)
// 		buff.WriteString(",\n")
// 	case []string:
// 		writeNameAndEqual(buff, name)
// 		buff.WriteString("{\n")
// 		for i := 0; i < len(v); i++ {
// 			writeTab(buff, tab+1)
// 			buff.WriteString("\"" + v[i] + "\",\n")
// 		}
// 		writeTab(buff, tab)
// 		buff.WriteString("},\n")
// 	case []int64:
// 		writeNameAndEqual(buff, name)
// 		buff.WriteString("{\n")
// 		for i := 0; i < len(v); i++ {
// 			writeTab(buff, tab+1)
// 			writeVal(fileName, lineno, buff, name, v[i])
// 			buff.WriteString(",\n")
// 		}
// 		writeTab(buff, tab)
// 		buff.WriteString("},\n")
// 	case []float64:
// 		writeNameAndEqual(buff, name)
// 		buff.WriteString("{\n")
// 		for i := 0; i < len(v); i++ {
// 			writeTab(buff, tab+1)
// 			writeVal(fileName, lineno, buff, name, v[i])
// 			buff.WriteString(",\n")
// 		}
// 		writeTab(buff, tab)
// 		buff.WriteString("},\n")
// 	case map[string]interface{}:
// 		writeNameAndEqual(buff, name)
// 		buff.WriteString("{\n")
// 		var keys []string
// 		for k := range v {
// 			keys = append(keys, k)
// 		}
// 		sort.Strings(keys)
// 		for i := 0; i < len(keys); i++ {
// 			encVal(fileName, lineno, buff, keys[i], v[keys[i]], tab+1)
// 		}
// 		writeTab(buff, tab)
// 		buff.WriteString("},\n")
// 	case []map[string]interface{}:
// 		writeNameAndEqual(buff, name)
// 		buff.WriteString("{\n")
// 		for i := 0; i < len(v); i++ {
// 			writeTab(buff, tab+1)
// 			buff.WriteString("{\n")
// 			var keys []string
// 			for k := range v[i] {
// 				keys = append(keys, k)
// 			}
// 			sort.Strings(keys)
// 			for j := 0; j < len(keys); j++ {
// 				encVal(fileName, lineno, buff, keys[j], v[i][keys[j]], tab+2)
// 			}
// 			writeTab(buff, tab+1)
// 			buff.WriteString("},\n")
// 		}
// 		writeTab(buff, tab)
// 		buff.WriteString("},\n")
// 	default:
// 		panic(fmt.Errorf("%s 第%d行 %s字段 无效的值类型1 %#v", fileName, lineno, name, val))
// 	}
// }

func encVal(fileName string, lineno int, buff *bytes.Buffer, name interface{}, val interface{}, tab int) {
	writeTab(buff, tab)
	switch v := val.(type) {
	case []interface{}:
		writeNameAndEqual(buff, name)
		buff.WriteString("{\n")
		for i := 0; i < len(v); i++ {
			if vvv, ok := v[i].(map[string]interface{}); ok {
				//encVal(fileName, lineno, buff, name, vvv, tab+1)
				writeTab(buff, tab+1)
				writeNameAndEqual(buff, i+1)
				buff.WriteString("{\n")
				var keys []string
				for k := range vvv {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for j := 0; j < len(keys); j++ {
					encVal(fileName, lineno, buff, keys[j], vvv[keys[j]], tab+2)
				}
				writeTab(buff, tab+1)
				buff.WriteString("},\n")
			} else if vvv, ok := v[i].([]interface{}); ok {
				encVal(fileName, lineno, buff, strconv.FormatInt(int64(i+1), 10), vvv, tab+1)
			} else {
				writeTab(buff, tab+1)
				writeVal(fileName, lineno, buff, name, v[i])
				buff.WriteString(",\n")
			}
		}
		writeTab(buff, tab)
		buff.WriteString("},\n")
	case string, int, uint, int64, uint64, float32, float64, nil:
		writeNameAndEqual(buff, name)
		writeVal(fileName, lineno, buff, name, v)
		buff.WriteString(",\n")
	case []string:
		writeNameAndEqual(buff, name)
		buff.WriteString("{\n")
		for i := 0; i < len(v); i++ {
			writeTab(buff, tab+1)
			buff.WriteString("\"" + v[i] + "\",\n")
		}
		writeTab(buff, tab)
		buff.WriteString("},\n")
	case []int64:
		writeNameAndEqual(buff, name)
		buff.WriteString("{\n")
		for i := 0; i < len(v); i++ {
			writeTab(buff, tab+1)
			writeVal(fileName, lineno, buff, name, v[i])
			buff.WriteString(",\n")
		}
		writeTab(buff, tab)
		buff.WriteString("},\n")
	case []float64:
		writeNameAndEqual(buff, name)
		buff.WriteString("{\n")
		for i := 0; i < len(v); i++ {
			writeTab(buff, tab+1)
			writeVal(fileName, lineno, buff, name, v[i])
			buff.WriteString(",\n")
		}
		writeTab(buff, tab)
		buff.WriteString("},\n")
	case map[string]interface{}:
		writeNameAndEqual(buff, name)
		buff.WriteString("{\n")
		var keys []string
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for i := 0; i < len(keys); i++ {
			encVal(fileName, lineno, buff, keys[i], v[keys[i]], tab+1)
		}
		writeTab(buff, tab)
		buff.WriteString("},\n")
	case []map[string]interface{}:
		writeNameAndEqual(buff, name)
		buff.WriteString("{\n")
		for i := 0; i < len(v); i++ {
			writeTab(buff, tab+1)
			buff.WriteString("{\n")
			var keys []string
			for k := range v[i] {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for j := 0; j < len(keys); j++ {
				encVal(fileName, lineno, buff, keys[j], v[i][keys[j]], tab+2)
			}
			writeTab(buff, tab+1)
			buff.WriteString("},\n")
		}
		writeTab(buff, tab)
		buff.WriteString("},\n")
	default:
		panic(fmt.Errorf("%s 第%d行 %s字段 无效的值类型1 %#v", fileName, lineno, name, val))
	}
}

func writeTab(buff *bytes.Buffer, n int) {
	for i := 0; i < n; i++ {
		buff.WriteString("\t")
	}
}

func writeNameAndEqual(buff *bytes.Buffer, name interface{}) {
	switch v := name.(type) {
	case string:
		_, err := strconv.ParseInt(v, 0, 64)
		if err != nil {
			buff.WriteString(v + " = ")
		} else {
			buff.WriteString("[" + v + "] = ")
		}
	case int:
		buff.WriteString("[" + strconv.FormatInt(int64(v), 10) + "] = ")
	case uint:
		buff.WriteString("[" + strconv.FormatUint(uint64(v), 10) + "] = ")
	case int32:
		buff.WriteString("[" + strconv.FormatInt(int64(v), 10) + "] = ")
	case uint32:
		buff.WriteString("[" + strconv.FormatUint(uint64(v), 10) + "] = ")
	case int64:
		buff.WriteString("[" + strconv.FormatInt(v, 10) + "] = ")
	case uint64:
		buff.WriteString("[" + strconv.FormatUint(v, 10) + "] = ")
	default:
		panic(fmt.Errorf("无效的fieldName %#v", name))
	}
}

func writeVal(fileName string, lineno int, buff *bytes.Buffer, name interface{}, v interface{}) {
	switch val := v.(type) {
	case string:
		buff.WriteString(`"` + val + `"`)
	case int:
		buff.WriteString(strconv.FormatInt(int64(val), 10))
	case uint:
		buff.WriteString(strconv.FormatUint(uint64(val), 10))
	case int64:
		buff.WriteString(strconv.FormatInt(val, 10))
	case uint64:
		buff.WriteString(strconv.FormatUint(val, 10))
	case float32:
		buff.WriteString(strconv.FormatFloat(float64(val), 'f', -1, 64))
	case float64:
		buff.WriteString(strconv.FormatFloat(val, 'f', -1, 64))
	case nil:
		buff.WriteString("nil")
	default:
		panic(fmt.Errorf("%s 第%d行 %v字段 无效的值类型2 %#v", fileName, lineno, name, val))
	}
}

// lua注释
func genLuaNotes(head []trans.Head) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("--[=[\n")
	for i := 0; i < len(head); i++ {
		if head[i].Flag == trans.FlagSrv {
			buf.WriteString("[server] ")
		} else if head[i].Flag == trans.FlagCli {
			buf.WriteString("[client] ")
		} else if head[i].Flag == trans.FlagDouble {
			buf.WriteString("[double] ")
		} else {
			buf.WriteString("[read]   ")
		}
		buf.WriteString(head[i].NameEN)
		buf.WriteString(": ")
		buf.WriteString(head[i].NameCN)
		buf.WriteString("\n")
	}
	buf.WriteString("]=]\n\n")
	return buf.Bytes(), nil
}
