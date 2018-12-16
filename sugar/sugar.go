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
	"runtime"
	"strings"
)

// url路径规则: /Prifix/Extend/tbName/...
type AdminConf struct {
	AccessControl string
	Address       string            // 0.0.0.0:9090
	Prefix        string            // 标注
	Extend        string            // 匹配数据库名字 默认 :dbname/
	extendKey     string            // 默认 dbname
	DatabaseAlias map[string]string //数据库别名, 用于修改url路径

	Relative    string // 默认 :tablename/
	relativeKey string // 默认 tablename

	BackupSuffix string // 注册表配置文件的扩展名

	//buttons           []string
	globalMiddlewares []gin.HandlerFunc
	groupMiddlewares  []gin.HandlerFunc
	loginFunc         gin.HandlerFunc
	//Static        string
	whiteUrls []string
	blackUrls []string
}

func (conf *AdminConf) AddWhite(urlSlice ...string) {
	for _, urlTmp := range urlSlice {
		if utils.InStringSlice(urlTmp, conf.whiteUrls) {
			fmt.Printf("Warning: URL '%s' already in whiteUrls.", urlTmp)
			continue
		}
		conf.whiteUrls = append(conf.whiteUrls, urlTmp)
	}
}

func (conf *AdminConf) DelWhite(urlSlice ...string) {
	conf.whiteUrls = utils.DelStringSliceEles(conf.whiteUrls, urlSlice...)
}

func (conf *AdminConf) AddBlack(urlSlice ...string) {
	for _, urlTmp := range urlSlice {
		if utils.InStringSlice(urlTmp, conf.blackUrls) {
			fmt.Printf("Warning: URL '%s' already in blackUrls.", urlTmp)
			continue
		}
		conf.blackUrls = append(conf.blackUrls, urlTmp)
	}
}

func (conf *AdminConf) DelBlack(urlSlice ...string) {
	conf.blackUrls = utils.DelStringSliceEles(conf.blackUrls, urlSlice...)
}

func (conf *AdminConf) CheckParams() {
	if conf.AccessControl == "" {
		conf.AccessControl = "static"
	} else if conf.AccessControl != "rbac" && conf.AccessControl != "static" {
		panic(errors.New("SugarAdminError: Access control type must 'rbac' or 'static'"))
	}

	if conf.Address == "" {
		conf.Address = "0.0.0.0:9090"
	}

	conf.checkPrefix()
	conf.checkBackupSuffix()
	conf.checkRelative()
	conf.checkExtend()
	//conf.checkStatic()
}

func (conf *AdminConf) checkBackupSuffix() {
	if conf.BackupSuffix == "" {
		conf.BackupSuffix = "backup"
		return
	}
	for _, b := range conf.BackupSuffix {
		if b < 'a' || b > 'Z' {
			panic(errors.New("SugarAdminError: BackupSuffix only be case letters"))
		}
	}
}

func (conf *AdminConf) checkPrefix() {
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

func (conf *AdminConf) checkRelative() {
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

func (conf *AdminConf) checkExtend() {
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
	Config AdminConf
	Sugar  *gin.Engine
	tables map[string]map[string]*descConf // database table list
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
	app.Sugar.Static(app.Config.Prefix+"static", tplPath)
	//app.Sugar.Static(app.Prefix+app.Static, tplPath)

}
func (app *appAdmin) InitApp(middleware ...gin.HandlerFunc) {
	app.Config.CheckParams()
	app.new(middleware...)
	app.htmlGlob()
	app.static()

}

func (app *appAdmin) InitGroup(middles ...gin.HandlerFunc) *gin.RouterGroup {

	sugarGroup := app.Sugar.Group(app.Config.Prefix)
	sugarGroup.Use(middles...)
	return sugarGroup

}

func (app *appAdmin) Start(back bool) {

	if back {
		go app.Sugar.Run(app.Config.Address)
	} else {
		app.Sugar.Run(app.Config.Address)
	}
}

var defaultTableHandle = &defaultDescAnalyzer{}

type tableDesc interface {
	DisplayName() string
}

var TableConfDirError = errors.New("TableConfDirError: Register dir is none.")
var TableConfTypeError = errors.New("TableConfTypeError: Register type is not supported.")
var TableConfPathError = errors.New("TableConfPathError: Register configuration file path error.")
var TableConfFileNameError = errors.New("TableConfFileNameError: Register configuration file name error.")

// 遍历目录获取所有的配置文件
func readDir(dirPath string, fileList []string) []string {
	flist, e := ioutil.ReadDir(dirPath)
	if e != nil {
		panic(TableConfPathError)
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

func checkConfFileList(fileList []string) {
	tmpTables := map[string][]string{}
	for _, file := range fileList {
		dbName, tbName, err := defaultAnalyzer.verifyPath(file)
		utils.PanicCheck(err)
		// 更新一个表的pool
		Dbm.UpdateDBPool(dbName)
		if _, ok := tmpTables[dbName]; ! ok {
			tmpTables[dbName] = Dbm.showTables(dbName)
		}
		tbList := tmpTables[dbName]

		if ! utils.InStringSlice(tbName, tbList) {
			errStr := fmt.Sprintf("RegisterTableNotFound: Register table name '%s' not found.", tbName)
			utils.PanicCheck(errors.New(errStr))
		}
		if _, ok := App.tables[dbName]; ! ok {
			App.tables[dbName] = make(map[string]*descConf)
		}
		dc, err := defaultAnalyzer.dump()
		checkDescConf(dc, tbName)
		fmt.Println("xxxxxxxxxxxxxx", dc)
		utils.PanicCheck(err)
		App.tables[dbName][tbName] = dc
	}
}

func checkDescConf(dc *descConf, tbName string) {
	if dc.Name == "" {
		dc.Name = tbName
	}
	if dc.Methods == nil {
		dc.Methods = append(dc.Methods, methods...)
	}

}

// 注册表配置, 遍历一个目录, 数据表文件命名规则: database.table.json/yml/xml
func Register(confPath string, confType string, analy analyzer) {
	confType = strings.ToLower(confType)

	changeAnalyzer(confType, analy)

	// 库名和表文件名, 用于支持跨库的表注册操作, 就要有多个数据库连接池
	var fileList []string
	fileList = readDir(confPath, fileList)

	checkConfFileList(fileList)

	if App.tables == nil {
		fmt.Println("Warning: No form to register was found.")
		return
	}

	for dbName, tbInfo := range App.tables {
		fmt.Println(dbName, tbInfo)
	}
	//if fileDic == nil {
	//	panic(TableConfDirError)
	//}
	//b, _ := json.Marshal(fileDic)

	//fmt.Println(string(b))

	//fmt.Println(fileList)
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

var App appAdmin

func init() {
	// 基本的配置文件
	Config("")
	c := AdminConf{DatabaseAlias: map[string]string{}}
	c.CheckParams()
	App = appAdmin{Config: c, tables: map[string]map[string]*descConf{}}

}


func SetAdmin(conf AdminConf) {
	App.Config = conf
	App.InitApp(gin.Logger(), gin.Recovery())
	rg := App.InitGroup()
	print(rg)
	//App.groupRouter = groupRouter{group: rg, conf: App.conf}
	//App.groupRouter.init()
}

// 全局的中间件
func AddGlobalMiddles(middles ...gin.HandlerFunc) {
	App.Config.globalMiddlewares = append(App.Config.globalMiddlewares, middles...)
}

// 单纯的Group中间件
func AddGroupMiddles(middles ...gin.HandlerFunc) {
	App.Config.groupMiddlewares = append(App.Config.groupMiddlewares, middles...)
}
