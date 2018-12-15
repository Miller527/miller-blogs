/*
# __author__ = "Mr.chai"
# Date: 2018/12/14
*/
package sugar

import (
	"fmt"
	"github.com/pkg/errors"
	"miller-blogs/sugar/utils"
	"reflect"
	"strings"
)

var confTypeList = []string{
	"yml",
	"json",
	"xml",
}


type descConf struct {
	Name    string
	Display string
	Field   []string
	Title   []string
	Filter  []string
	Desc    map[string]string
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
	dumps()
	dump()
	loads()
	load()

}

var defaultAnalyzer analyzer

type jsonAnalyzer struct {
	data []byte
}

func (ana jsonAnalyzer) dumps() {

}
func (ana jsonAnalyzer) dump() {

}

func (ana jsonAnalyzer) load() {

}

func (ana jsonAnalyzer) loads() {

}

type yamlAnalyzer struct {
	data []byte

}

func (ana yamlAnalyzer) dumps() {

}
func (ana yamlAnalyzer) dump() {

}

func (ana yamlAnalyzer) load() {

}

func (ana yamlAnalyzer) loads() {

}

type xmlAnalyzer struct {
	data []byte
}

func (ana xmlAnalyzer) dumps() {

}
func (ana xmlAnalyzer) dump() {

}

func (ana xmlAnalyzer) load() {

}

func (ana xmlAnalyzer) loads() {

}




func analyzerInit(confType string) analyzer{
	switch confType {
	case "json":
		return jsonAnalyzer{}
	case "yaml":
		return yamlAnalyzer{}
	case "xml":
		return xmlAnalyzer{}
	default:
		errStr := fmt.Sprintf("AnalyzerInitError: init analyzer type '%s' error.",confType)
		panic(errors.New(errStr))
	}
}

// 注册表分析器接口修改
func changeAnalyzer(analy analyzer, confType string){
	if analy == nil{
		defaultAnalyzer = analyzerInit(confType)
	}else {
		defaultAnalyzer = analy
	}
}