package trans

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func splitFileName(filename string) (string, string, string) {
	s := strings.TrimSuffix(path.Base(filename), path.Ext(filename))
	ss := strings.Split(s, " - ")
	if len(ss) >= 2 {
		return strings.TrimSpace(s), strings.TrimSpace(ss[1]), strings.TrimSpace(ss[0])
	} else {
		return s, s, s
	}
}

func checkFileMd5() (diff, all []*_Mapping, err error) {
	var list []_MD5Pair
	md5Data, _ := ioutil.ReadFile(path.Join(conf.XlsxDir, "md5.json"))
	json.Unmarshal(md5Data, &list)

	infos, err := ioutil.ReadDir(conf.XlsxDir)
	if err != nil {
		return
	}

	for _, info := range infos {
		filename := info.Name()
		if info.IsDir() || !isXlsxFile(filename) {
			continue
		}
		basename, nameEN, nameCN := splitFileName(filename)

		mapping := _Mapping{CNFile: nameCN, XlsxFile: basename, OutFile: nameEN}
		all = append(all, &mapping)

		found := false
		for i := 0; i < len(list); i++ {
			if list[i].FileName == basename {
				found = true
				if md5file(path.Join(conf.XlsxDir, filename)) != list[i].Md5 {
					diff = append(diff, &mapping)
				}
				break
			}
		}
		if !found {
			diff = append(diff, &mapping)
		}
	}

	return
}

func writeFileMd5(all []*_Mapping) {
	var list []_MD5Pair
	for _, v := range all {
		md5 := md5file(path.Join(conf.XlsxDir, v.XlsxFile+".xlsx"))
		list = append(list, _MD5Pair{FileName: v.XlsxFile, Md5: md5})
	}
	data, _ := json.MarshalIndent(&list, "", "\t")
	file, _ := os.OpenFile(path.Join(conf.XlsxDir, "md5.json"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	file.Write(data)
	file.Close()
}
