/*
# __author__ = "Mr.chai"
# Date: 2018/12/14
*/
package sugar

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"miller-blogs/sugar/utils"
	"os"
	"path/filepath"
	"strings"
)

var confTypeList = []string{
	"yml",
	"yaml",
	"json",
	"xml",
}

type tableDesc interface {
	DisplayName() string
}

type descConf struct {
	Name      string
	Display   string
	Primary   string            // todo 扩展, 手动配置和自动查询数据库配置, 目前自动匹配
	Foreign   map[string]string // todo 扩展, 手动配置和自动查询数据库配置, 目前自动匹配, 还没用到就没写呢
	Field     []string
	Title     []string
	Filter    map[string]string // todo 字段和对应的类型、根据DescType匹配生成
	DescField []string          // todo 查询数据库生成, 如果Field没有配置, 那么使用该字段
	DescType  map[string]string // todo 查询数据库生成
	Left      bool
	Right     bool
	Methods   []int
}

func (dc *descConf) DisplayName() string {
	if dc.Display ==""{
		return dc.Name
	}
	return dc.Display + "(" + dc.Name + ")"
}

type IAnalyzer interface {
	dump() (*descConf, error)
	dumps(r io.Reader) (*descConf, error)
	load(desc *descConf) error
	loads(r io.Writer, desc *descConf) error
	verifyPath(fp string) (string, string, error)
}

var defaultAnalyzer IAnalyzer

type jsonAnalyzer struct {
	FileSuffix string
	confPath   string
}

func (jana *jsonAnalyzer) verifyPath(fp string) (string, string, error) {
	nameFields := strings.Split(filepath.Base(fp), ".")
	lenField := len(nameFields)
	if lenField == 4 {
		// 以backup后缀的为备份文件，不需要解析
		if nameFields[len(nameFields)-1] != App.Config.BackupSuffix {
			return "", "", TableConfFileNameError
		} else {
			return "", "", TableConfBackupWarning
		}
	} else if lenField != 3 || strings.ToLower(nameFields[2]) != jana.FileSuffix { // 过滤配置文件设置
		return "", "", TableConfFileNameError
	}
	jana.confPath = fp

	return nameFields[0], nameFields[1], nil
}

func (jana *jsonAnalyzer) dump() (*descConf, error) {
	file, err := os.Open(jana.confPath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	return jana.dumps(file)
}

func (jana *jsonAnalyzer) dumps(r io.Reader) (*descConf, error) {
	decoder := json.NewDecoder(r)
	desc := &descConf{}

	err := decoder.Decode(desc)
	if err == nil {
		return desc, nil
	}
	return nil, err
}

func (jana *jsonAnalyzer) load(desc *descConf) error {
	file, err := os.Create(jana.confPath)
	defer file.Close()
	if err != nil {
		return err
	}
	return jana.loads(file, desc)
}

func (jana *jsonAnalyzer) loads(r io.Writer, desc *descConf) error {
	encoder := json.NewEncoder(r)
	return encoder.Encode(desc)

}

type yamlAnalyzer struct {
	FileSuffix string
	confPath   string
}

func (yana *yamlAnalyzer) verifyPath(fp string) (string, string, error) {
	return "", "", nil
}
func (yana *yamlAnalyzer) dumps(r io.Reader) (*descConf, error) {
	return nil, nil

}
func (yana *yamlAnalyzer) dump() (*descConf, error) {
	return nil, nil

}
func (yana *yamlAnalyzer) load(desc *descConf) error {
	return nil
}
func (yana *yamlAnalyzer) loads(r io.Writer, desc *descConf) error {
	return nil

}

type xmlAnalyzer struct {
	FileSuffix string
	confPath   string
}

func (xana *xmlAnalyzer) verifyPath(fp string) (string, string, error) {
	return "", "", nil

}
func (xana *xmlAnalyzer) dumps(r io.Reader) (*descConf, error) {
	return nil, nil

}
func (xana *xmlAnalyzer) dump() (*descConf, error) {
	return nil, nil

}

func (xana *xmlAnalyzer) load(desc *descConf) error {
	return nil

}

func (xana *xmlAnalyzer) loads(r io.Writer, desc *descConf) error {
	return nil
}
func analyzerInit(confType string) IAnalyzer {
	switch confType {
	case "json":
		return &jsonAnalyzer{FileSuffix: confType}
	case "yaml", "yml":
		return &yamlAnalyzer{FileSuffix: confType}
	case "xml":
		return &xmlAnalyzer{FileSuffix: confType}
	default:
		errStr := fmt.Sprintf("AnalyzerInitError: init analyzer type '%s' error.", confType)
		panic(errors.New(errStr))
	}
}

// 注册表分析器接口修改
func changeAnalyzer(confType string, analy IAnalyzer) {
	if analy == nil && utils.InStringSlice(confType, confTypeList) {
		defaultAnalyzer = analyzerInit(confType)
	} else if analy != nil && !utils.InStringSlice(confType, confTypeList) {
		defaultAnalyzer = analy
	} else {
		panic(TableConfTypeError)
	}
}
