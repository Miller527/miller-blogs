//
// __author__ = "Miller"
// Date: 2018/11/24
//

package sugar

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"miller-blogs/sugar/utils"
	"path"
	"reflect"
	"runtime"
	"strings"
)

type AdminConfig struct {
	AccessControl     string
	Address           string
	Prefix            string
	Extend            string
	Relative          string
	relativeKey       string
	buttons           []string
	globalMiddlewares []gin.HandlerFunc
	groupMiddlewares  []gin.HandlerFunc
	loginFunc         gin.HandlerFunc
	//Static        string
	whiteUrls []string
	blackUrls []string
}

func (conf *AdminConfig) AddWhite(urlSlice ...string) {
	for _, urlTmp := range urlSlice {
		conf.whiteUrls = append(conf.whiteUrls, urlTmp)
	}

}

func (conf *AdminConfig) AddBlack(urlSlice ...string) {
	for _, urlTmp := range urlSlice {
		conf.blackUrls = append(conf.blackUrls, urlTmp)
	}
}

func (conf *AdminConfig) CheckParams() {
	if conf.AccessControl == "" {
		conf.AccessControl = "static"
	} else if conf.AccessControl != "rbac" && conf.AccessControl != "static" {
		panic(errors.New("SugarAdminError: Access control type must 'rbac' or 'static'"))
	}

	if conf.Address == "" {
		conf.Address = "0.0.0.0:9090"
	}

	conf.checkPrefix()
	conf.checkRelative()
	conf.checkExtend()
	//conf.checkStatic()
}

func (conf *AdminConfig) checkPrefix() {
	if conf.Prefix == "" {
		conf.Prefix = "/"
		return
	}
	if ! strings.HasPrefix(conf.Prefix, "/") {
		conf.Prefix = "/" + conf.Prefix
	}
	if ! strings.HasSuffix(conf.Prefix, "/") {
		conf.Prefix += "/"
	}
}

func (conf *AdminConfig) checkRelative() {
	if conf.Relative == "" {
		if conf.AccessControl == "rbac" {
			conf.Relative = ":tablename/"
		} else {
			conf.Relative = "tablename/"
		}
		conf.relativeKey = "tablename"

	} else {
		for _, b := range conf.Relative {
			if b < 'a' || b > 'Z' || b != '_' {
				panic(errors.New("SugarAdminError: Relative only be case letters"))
			}
		}
		conf.relativeKey = conf.Relative

		conf.Relative += "/"
		if conf.AccessControl == "rbac" {
			conf.Relative = ":" + conf.Relative
		}
	}
}

func (conf *AdminConfig) checkExtend() {
	if conf.Extend == "" {
		conf.Extend = "curd/"
		return
	}
	if strings.HasPrefix(conf.Extend, "/") {
		conf.Extend = conf.Extend[1:]
	}
	if ! strings.HasSuffix(conf.Extend, "/") {
		conf.Extend += "/"
	}
}

// todo 前端代码里怎么改动（预加载选择主题时候的路径问题）
//func (conf *Config) checkStatic() {
//	if conf.Static == "" {
//		conf.Static = "static"
//		return
//	}
//	if strings.HasPrefix(conf.Static, "/") {
//		conf.Static = conf.Static[1:]
//	}
//	if ! strings.HasSuffix(conf.Static, "/") {
//		conf.Static = conf.Static[:len(conf.Static)-1]
//	}
//}

type appAdmin struct {
	conf  AdminConfig
	Sugar *gin.Engine
	//groupRouter
	//registry map[string]*TableConf
}

func (app *appAdmin) new(middleware ...gin.HandlerFunc) {
	app.Sugar = gin.New()
	app.Sugar.Use(middleware...)
}

func (app *appAdmin) htmlGlob() {
	_, file, _, ok := runtime.Caller(0)
	if ! ok {
		panic(errors.New("SugarAdminError: get template path error"))
	}
	tplPath := path.Join(path.Dir(file), "template", "**", "*")
	app.Sugar.LoadHTMLGlob(tplPath)

}

func (app *appAdmin) static() {
	_, file, _, ok := runtime.Caller(0)
	if ! ok {
		panic(errors.New("SugarAdminError: get template path error"))
	}
	tplPath := path.Join(path.Dir(file), "static")
	app.Sugar.Static(app.conf.Prefix+"static", tplPath)
	//app.Sugar.Static(app.Prefix+app.Static, tplPath)

}
func (app *appAdmin) InitApp(middleware ...gin.HandlerFunc) {
	app.conf.CheckParams()
	app.new(middleware...)
	app.htmlGlob()
	app.static()

}

func (app *appAdmin) InitGroup(middles ...gin.HandlerFunc) *gin.RouterGroup {

	sugarGroup := app.Sugar.Group(app.conf.Prefix)
	sugarGroup.Use(middles...)
	return sugarGroup

}

func (app *appAdmin) Start(back bool) {

	if back {
		go app.Sugar.Run(app.conf.Address)
	} else {
		app.Sugar.Run(app.conf.Address)
	}
}

var defaultTableHandle = &defaultDescAnalyzer{}

type tableDesc interface {
	DisplayName() string
}

// 遍历目录获取所有的配置文件
func readDir(dirPath string, fileList []string) []string {
	flist, e := ioutil.ReadDir(dirPath)
	if e != nil {
		return nil
	}
	for _, f := range flist {
		if f.IsDir() {
			fileList = readDir(dirPath+"/"+f.Name(), fileList)
		} else {
			fileList = append(fileList, dirPath+"/"+f.Name())
		}

	}
	return fileList
}

var TableConfDirError = errors.New("TableConfDir: Register dir is none.")

// 注册表配置, 遍历一个目录
func Register(confType string, confPath string) {
	if ! utils.InStringSlice(confType, confTypeList) {

	}

	var fileList []string
	fileList = readDir(confPath, fileList)
	if fileList == nil {
		panic(TableConfDirError)
	}
	fmt.Println(fileList)
	//handle := tc
	//if handle == nil {
	//	handle = defaultTableHandle
	//}
	//tableConfig := handle.ParseDesc(desc)
	//fmt.Println(tableConfig.Display, tableConfig.Name, tableConfig.Methods, tableConfig.Field)
	//
	//if len(tc.Field) != len(tc.Title) {
	//	panic(errors.New("SugarTable: Table field length unequal to title length"))
	//}
	//name := tc.Name()
	//if ! verifyName(name) {
	//	panic(errors.New("SugarTable: database not found [" + name + "] table"))
	//}
	//
	//if ! verifyField(tc) {
	//	panic(errors.New("SugarTable: Table [" + name + "] Field error"))
	//}
	//
	//if _, ok := App.registry[name]; ok {
	//	panic(errors.New("SugarTable: table [" + name + "] has already registered"))
	//}
	//
	//App.registry[name] = tc
}

func init() {
	//c := Config{}
	//App = appAdmin{conf: c, registry: map[string]*TableConf{}}

}

var App appAdmin

func SetAdmin(conf AdminConfig) {
	App.conf = conf
	App.InitApp(gin.Logger(), gin.Recovery())
	rg := App.InitGroup()
	print(rg)
	//App.groupRouter = groupRouter{group: rg, conf: App.conf}
	//App.groupRouter.init()
}

// 全局的中间件
func AddGlobalMiddles(middles ...gin.HandlerFunc) {
	App.conf.globalMiddlewares = append(App.conf.globalMiddlewares, middles...)
}

// 单纯的Group中间件
func AddGroupMiddles(middles ...gin.HandlerFunc) {
	App.conf.groupMiddlewares = append(App.conf.groupMiddlewares, middles...)
}

type descConf struct {
	Name    string
	Display string
	Field   []string
	Title   []string
	Filter 	[]string
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
func (tc *defaultDescAnalyzer) Name(desc interface{}) string {
	return ""
}

func (tc *defaultDescAnalyzer) DisplayName(desc interface{}) string {
	return ""

}

//
//// 验证表名字
//func verifyName(name string) bool {
//	if ! utils.InStringSlice(name, tables) {
//		return false
//	}
//	return true
//}
//
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
//
//	for i, f := range tc.Field {
//		f = utils.SnakeString(f)
//		tc.Field[i] = f
//		if ! utils.InStringSlice(f, result) {
//			return false
//		}
//	}
//	tc.Title = append(tc.Title, "操作")
//
//	return true
//}
