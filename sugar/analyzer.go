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
	"reflect"
	"strings"
)

var confTypeList = []string{
	"yml",
	"yaml",
	"json",
	"xml",
}

type descConf struct {
	Name    string
	Display string
	Field   []string
	Title   []string
	Filter  map[string]string	// todo 字段和对应的类型
	Desc    map[string]string	// todo 查询数据库生成
	Left    bool
	Right   bool
	Methods []int
}

// 配置表接口
type TableHandle interface {
	//Name(desc tableDesc) string
	//DisplayName(desc tableDesc) string
	ParseDesc(desc tableDesc) *descConf
}

//
type TableConf struct {
	Left    bool
	Right   bool
	Methods []int
}
type defaultDescAnalyzer struct {
	Display     string
	DisplayJoin bool
	//Desc interface{}
}

func (da *defaultDescAnalyzer) ParseDesc(desc tableDesc) *descConf {
	//field, title, primary = da.getField()
	return &descConf{
		Name:    da.getName(desc),
		Display: da.getDisplay(desc),
		Field:   da.getField(desc),
		Title:   da.getTitle(desc),
	}
}

func (da *defaultDescAnalyzer) getName(desc tableDesc) string {
	tmpSlice := strings.Split(reflect.TypeOf(desc).String(), ".")
	return utils.SnakeString(tmpSlice[len(tmpSlice)-1])
}

func (da *defaultDescAnalyzer) getDisplay(desc tableDesc) string {
	return desc.DisplayName() + "(" + da.getName(desc) + ")"
}
func (da *defaultDescAnalyzer) getField(desc tableDesc) []string {
	//var fields []string
	//var titles []string
	//var primary string
	value := reflect.ValueOf(desc)
	//fmt.Println(value.CanSet())
	//te := value.Type()
	//n := value.Type().NumField()
	//if value.Kind() != reflect.Ptr {
	//	fmt.Println("xxxxxxxxxxxxxx")
	//}
	fmt.Println(value.Kind())
	//elem := value.Elem()
	//for i:=0;i< 3 ;i++{
	//	elemField := elem.Field(i)
	//	switch elemField.Kind() {
	//	case 	reflect.Struct:
	//		fmt.Println("xxxxxxxxxxxxxxxxxxxxxx")
	//	}
	//	//fmt.Println(te.Field(i), )
	//	//
	//	//fields = append(fields, utils.SnakeString(te.Field(i).Name))
	//	//
	//	//tit := te.Field(i).Tag.Get("title")
	//	//if tit != ""{
	//	//	titles = append(titles, tit)
	//	//}
	//	//
	//	//if primary == ""{
	//	//	primary = te.Field(i).Tag.Get("primary")
	//	//
	//	//}
	//	//fmt.Println(tit,te.Field(i).Tag.Get("title"),primary)
	//
	//}
	//fmt.Println(fields, titles, primary)
	//for t.Elem().
	//field := t.Elem().Field(0)
	//fmt.Println(field.Tag)
	return nil
}
func (da *defaultDescAnalyzer) getTitle(desc tableDesc) []string {
	var titles []string
	return titles
}

//
func (da *defaultDescAnalyzer) Name(desc interface{}) string {
	return ""
}

func (da *defaultDescAnalyzer) DisplayName(desc interface{}) string {
	return ""
}

//// 验证表名字
//func verifyName(name string) bool {
//	if ! utils.InStringSlice(name, tables) {
//		return false
//	}
//	return true
//}
//func verifyField(tc *TableConf) bool {
//	sqlCmd := `select COLUMN_NAME as name,DATA_TYPE as dataType
//			   from information_schema.COLUMNS
//			   where table_schema=? AND table_name=?`
//	stmt, err := Dbm.Db.Prepare(sqlCmd)
//	type desc struct {
//		name     string
//		dataType string
//	}
//	column := &TableConf{
//		Field: []string{"name", "dataType"},
//		Desc:  &desc{},
//	}
//	result, err := Dbm.SelectValues(stmt, column, Dbm.Conf.DBName, tc.Name())
//	if err != nil {
//		fmt.Println("verifyField", result, err)
//		return false
//	}
//	for i, f := range tc.Field {
//		f = utils.SnakeString(f)
//		tc.Field[i] = f
//		if ! utils.InStringSlice(f, result) {
//			return false
//		}
//	}
//	tc.Title = append(tc.Title, "操作")
//	return true
//}

type analyzer interface {
	dump() (*descConf, error)
	dumps(r io.Reader) (*descConf, error)
	load(desc *descConf) error
	loads(r io.Writer, desc *descConf) error
	verifyPath(fp string) (string, string, error)
}

var defaultAnalyzer analyzer

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
		}
	} else if lenField != 3 || strings.ToLower(nameFields[2]) != jana.FileSuffix { // 过滤配置文件设置
		return "", "", TableConfFileNameError
	}
	jana.confPath = fp
	fmt.Println("confPath1", jana.confPath)

	return nameFields[0], nameFields[1], nil
}

func (jana *jsonAnalyzer) dump() (*descConf, error) {
	fmt.Println("confPath2", jana.confPath)
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
func analyzerInit(confType string) analyzer {
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
func changeAnalyzer(confType string, analy analyzer) {
	if analy == nil && utils.InStringSlice(confType, confTypeList) {
		defaultAnalyzer = analyzerInit(confType)
	} else if analy != nil && !utils.InStringSlice(confType, confTypeList) {
		defaultAnalyzer = analy
	} else {
		panic(TableConfTypeError)
	}
}
