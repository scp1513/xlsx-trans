package trans

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func md5file(fn string) string {
	file, err := os.Open(fn)
	if err != nil {
		return err.Error()
	}
	hash := md5.New()
	var data [4096]byte
	for {
		n, err := file.Read(data[:])
		if n == 0 && err != nil {
			break
		}
		hash.Write(data[:n])
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func md5str(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func isXlsxFile(filename string) bool {
	ok, _ := filepath.Match("[^~$]*.xlsx", filename)
	return ok
}

func printTime(num int, str string, val int) {
	str = "\t" + str + ":"
	l := len(str)
	for i := l; i < num; i++ {
		str += " "
	}
	fmt.Printf("%s%.2f(ms)\n", str, float64(val)/1000000)
}

func sortKeys(m map[int]map[string]interface{}) []int {
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

// Sprintf 格式化字符串
func Sprintf(file, field string, row int, format string, msg ...interface{}) string {
	prefix := fmt.Sprintf("[%d:%s]: ", row, field)
	return fmt.Sprintf(prefix+format, msg...)
}

// Errorf 格式化字符串
func Errorf(file, field string, row int, format string, msg ...interface{}) error {
	prefix := fmt.Sprintf("[%d:%s]: ", row, field)
	return fmt.Errorf(prefix+format, msg...)
}
