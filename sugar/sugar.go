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
	"path/filepath"
	"runtime"
	"strings"
)

type AdminConf struct {
	AccessControl string
	Address       string // 0.0.0.0:9090
	Prefix        string
	Extend        string //

	BackupSuffix      string // 扩展名
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

func (conf *AdminConf) AddWhite(urlSlice ...string) {
	for _, urlTmp := range urlSlice {
		conf.whiteUrls = append(conf.whiteUrls, urlTmp)
	}

}

func (conf *AdminConf) AddBlack(urlSlice ...string) {
	for _, urlTmp := range urlSlice {
		conf.blackUrls = append(conf.blackUrls, urlTmp)
	}
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
	conf   AdminConf
	Sugar  *gin.Engine
	tables map[string]map[string]*TableConf // database table list
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


func checkConfFileList(confType string, fileList []string) {

	tmpTables := map[string][]string{}

	for _, file := range fileList {
		fileNameFields := strings.Split(filepath.Base(file), ".")

		lenField := len(fileNameFields)
		if lenField == 4 {
			// 以backup后缀的为备份文件，不需要解析
			if fileNameFields[len(fileNameFields)-1] != App.conf.BackupSuffix {
				panic(TableConfFileNameError)
			}
			continue
		} else if lenField != 3 {
			panic(TableConfFileNameError)
		}
		// 过滤配置文件设置
		if strings.ToLower(fileNameFields[2]) != confType {
			continue
		}

		dbName, tbName, extName := fileNameFields[0], fileNameFields[1], fileNameFields[2]

		// 更新一个表的pool
		Dbm.UpdateDBPool(dbName)
		if _, ok := tmpTables[dbName]; ! ok {
			tmpTables[dbName] = Dbm.showTables(dbName)
		}
		tbList := tmpTables[dbName]

		fmt.Println(fileNameFields, utils.CamelString(filepath.Base(file)))
		fmt.Println(dbName, tbName, extName, tbList)
		if ! utils.InStringSlice(tbName, tbList) {
			errStr := fmt.Sprintf("RegisterTableNotFound: Register table name '%s' not found.", tbName)
			utils.PanicCheck(errors.New(errStr))
		}
		if _, ok := App.tables[dbName]; ! ok {
			App.tables[dbName] = make(map[string]*TableConf)
		}
		App.tables[dbName][tbName] = &TableConf{}

		fmt.Println(tbList, tmpTables)
		fmt.Println(tbList,  App.tables)
	}

}

// 注册表配置, 遍历一个目录, 数据表文件命名规则: database.table.json/yml/xml
func Register(confType string, confPath string,  analy analyzer) {
	confType = strings.ToLower(confType)
	if ! utils.InStringSlice(confType, confTypeList) {
		panic(TableConfTypeError)
	}
	changeAnalyzer(analy, confType)

	// 库名和表文件名, 用于支持跨库的表注册操作, 就要有多个数据库连接池
	var fileList []string
	fileList = readDir(confPath, fileList)
	fmt.Println("x", confType, fileList)

	checkConfFileList(confType, fileList)
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

func init() {
	c := AdminConf{}
	c.CheckParams()
	App = appAdmin{conf: c, tables: map[string]map[string]*TableConf{}}

}

var App appAdmin

func SetAdmin(conf AdminConf) {
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
