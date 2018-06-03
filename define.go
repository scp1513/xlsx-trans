package trans

import (
	"io"
)

// XlsxRow xlsx行数据
type XlsxRow = map[string]interface{}

// XlsxRows xlsx所有行数据
type XlsxRows = map[int]XlsxRow

// Output 输出格式
type Output interface {
	Type() string
	Dir() string
	Flag() FlagType
	BeforeAll()
	AfterAll()
	Before(heads []Head, file string, w io.Writer)
	After(heads []Head, file string, w io.Writer)
	Row(heads []Head, data XlsxRow, row int, w io.Writer)
}

type _Config struct {
	XlsxDir string
}

// xlsx文件名与输出文件名映射关系
type _Mapping struct {
	XlsxFile string
	CNFile   string
	OutFile  string
}

type _TransResult struct {
	err         error
	mapping     *_Mapping
	loadTime    int
	analyseTime int
	marshalTime map[string]int
}

// Head xlsx头信息
type Head struct {
	col    int
	NameCN string
	NameEN string
	Flag   FlagType
	Type   string
}

type _MD5Pair struct {
	FileName string
	Md5      string
}

// FlagType 标记类型
type FlagType = string

const (
	// FlagRead 只读
	FlagRead = "r"
	// FlagSrv 服务器
	FlagSrv = "s"
	// FlagCli 客户端
	FlagCli = "c"
	// FlagDouble 双端
	FlagDouble = "d"
)

func isFlagValid(v FlagType) bool {
	return v == FlagSrv || v == FlagCli || v == FlagDouble
}
