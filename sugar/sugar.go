//
// __author__ = "Miller"
// Date: 2018/11/24
//

package sugar

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"miller-blogs/sugar/utils"
	"path"
	"reflect"
	"runtime"
	"strings"
)

type Config struct {
	AccessControl     string
	Address           string
	Prefix            string
	Extend            string
	Relative          string
	relativeKey       string
	buttons       []string
	globalMiddlewares []gin.HandlerFunc
	groupMiddlewares  []gin.HandlerFunc
	loginFunc         gin.HandlerFunc
	//Static        string
	whiteUrls []string
	blackUrls []string
}

func (conf *Config) AddWhite(urlSlice ...string) {
	for _, urlTmp := range urlSlice {
		conf.whiteUrls = append(conf.whiteUrls, urlTmp)
	}

}
func (conf *Config) AddBlack(urlSlice ...string) {
	for _, urlTmp := range urlSlice {
		conf.blackUrls = append(conf.blackUrls, urlTmp)
	}
}

func (conf *Config) CheckParams() {
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
func (conf *Config) checkPrefix() {
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

func (conf *Config) checkRelative() {
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

func (conf *Config) checkExtend() {
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
	conf  Config
	Sugar *gin.Engine
	groupRouter
	registry map[string]*TableConf
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



//// 注册表配置
func Register(desc interface{}, method []int, table Table) {
	if table == nil{
		tabl := &TableConf{Desc:desc, Methods:method}
	}
	for _, tc := range tcList {

		if len(tc.Field) != len(tc.Title) {
			panic(errors.New("SugarTable: Table field length unequal to title length"))
		}
		name := tc.Name()
		if ! verifyName(name) {
			panic(errors.New("SugarTable: database not found [" + name + "] table"))
		}

		if ! verifyField(tc) {
			panic(errors.New("SugarTable: Table [" + name + "] Field error"))
		}

		if _, ok := App.registry[name]; ok {
			panic(errors.New("SugarTable: table [" + name + "] has already registered"))
		}

		App.registry[name] = tc
	}
}

func init() {
	c := Config{}
	x := make(map[*TableConf]interface{})
	App = appAdmin{conf: c,registry:}

}

var App appAdmin

func SetAdmin(conf Config) {
	App.conf = conf
	App.InitApp(gin.Logger(), gin.Recovery())
	rg := App.InitGroup()
	App.groupRouter = groupRouter{group: rg, conf: App.conf}
	App.groupRouter.init()
}

// 全局的中间件
func AddGlobalMiddles(middles ...gin.HandlerFunc) {
	App.conf.globalMiddlewares = append(App.conf.globalMiddlewares, middles...)
}

// 单纯的Group中间件
func AddGroupMiddles(middles ...gin.HandlerFunc) {
	App.conf.groupMiddlewares = append(App.conf.groupMiddlewares, middles...)
}

// 配置表接口
type Table interface {
	Name() string
	DisplayName() string
}
//
type TableConf struct {
	Display     string
	DisplayJoin bool
	Field       []string
	Title       []string
	Methods     []int
	Desc interface{}

}
//
func (tc *TableConf) Name() string {
	tmpSlice := strings.Split(reflect.TypeOf(tc.Desc).String(), ".")
	return utils.SnakeString(tmpSlice[len(tmpSlice)-1])
}

func (tc *TableConf) DisplayName() string {
	if tc.Display == "" {
		return tc.Name()
	}
	return tc.Display + " ( " + tc.Name() + " )"
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
