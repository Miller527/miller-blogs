//
// __author__ = "Miller"
// Date: 2018/11/24
//

package sugar

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/ioutil"
	"log"
	"miller-blogs/sugar/utils"
	"path"
	"runtime"
	"strings"
	"unicode"
)

var TableConfDirError = errors.New("TableConfDirError: Register dir is none.")
var DBAliasSetError = errors.New("DBAliasSetError: database alias is None.")
var TBAliasSetError = errors.New("TBAliasSetError: table alias is None.")
var TableConfTypeError = errors.New("TableConfTypeError: Register type is not supported.")
var TableConfPathError = errors.New("TableConfPathError: Register configuration file path error.")
var TableConfFileNameError = errors.New("TableConfFileNameError: Register configuration file name error.")
var TableConfBackupWarning = errors.New("TableConfBackupWarning: Register configuration file is buckup.")

// url路径规则: /Prifix/Extend/tbName/...
type AdminConf struct {
	AccessControl string
	Address       string // 0.0.0.0:9090
	Static        string
	StaticPrefix  bool   // 是否使用前缀
	Prefix        string // 前缀
	Relative      string // 默认 :tablename/
	RelativeKey   string // 默认 tablename
	ExtendKey     string // 中间件获取数据库名
	BackupSuffix  string // 注册表配置文件的扩展名

	//buttons           []string
	globalMiddlewares []gin.HandlerFunc
	groupMiddlewares  []gin.HandlerFunc
	LoginFunc         gin.HandlerFunc
	VerifyLoginFunc   gin.HandlerFunc
	whiteUrls         []string
	blackUrls         []string
}
func (conf *AdminConf)  WhiteUrls() []string {
	return conf.whiteUrls

}
func(conf *AdminConf)  BlackUrls() []string {
	return conf.blackUrls

}
// 增加全局中间件
func (conf *AdminConf) AddGlobalMiddle(middles ...gin.HandlerFunc) {
	conf.globalMiddlewares = append(conf.globalMiddlewares, middles...)
}

// 增加组中间件
func (conf *AdminConf) AddGroupMiddle(middles ...gin.HandlerFunc) {
	conf.groupMiddlewares = append(conf.groupMiddlewares, middles...)

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
	conf.checkStatic()
	conf.checkBackupSuffix()
	conf.checkRelative()
	conf.checkExtend()
	conf.AddWhite(conf.Prefix + "login")
	conf.AddWhite(conf.Prefix + "verify-login")
}

func (conf *AdminConf) checkBackupSuffix() {
	if conf.BackupSuffix == "" {
		conf.BackupSuffix = "backup"
		return
	}
	for _, b := range conf.BackupSuffix {
		if ! unicode.IsLetter(b) {
			panic(errors.New("SugarAdminError: BackupSuffix only be case letters"))
		}
	}
}

// 判断url上的前缀字段
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
func (conf *AdminConf) checkStatic() {
	if conf.Static == "" {
		conf.Static = "static/"
	} else {
		if strings.HasPrefix(conf.Static, "/") {
			conf.Static = conf.Static[1:]
		}
		if ! strings.HasSuffix(conf.Static, "/") {
			conf.Static = conf.Static +"/"
		}
	}
	if conf.StaticPrefix {
		conf.Static = conf.Prefix + conf.Static
	}else {
		conf.Static =  "/" + conf.Static

	}
}

// 判断url上的表名字段
func (conf *AdminConf) checkRelative() {

	if conf.Relative == "" {
		if conf.AccessControl == "rbac" {
			conf.Relative = ":tablename/"
		}
		conf.RelativeKey = "tablename"

	} else {
		for _, b := range conf.Relative {
			if ! unicode.IsLetter(b) && b != '_' && b != ':' && b != '/' {

				panic(errors.New("SugarAdminError: Relative only be case letters"))
			}
		}
		conf.RelativeKey = conf.Relative

		if conf.AccessControl == "rbac" && !strings.HasPrefix(conf.Relative, ":") {
			conf.Relative = ":" + conf.Relative
		}
		if !strings.HasSuffix(conf.Relative, "/") {
			conf.Relative += "/"

		}
	}
}

func (conf *AdminConf) checkExtend() {
	for _, b := range conf.ExtendKey {
		if b < 'a' || b > 'Z' || b != '_' {
			panic(errors.New("SugarAdminError: Extend only be case letters"))
		}
	}

}

var App appAdmin

type appAdmin struct {
	Config      *AdminConf
	Sugar       *gin.Engine
	Registry    map[string]map[string]*descConf // database table list
	GroupRouter groupRouter
	//registry map[string]*TableConf
	// 先从别名中找, 然后从原名中找
	databaseAlias map[string]string            //数据库名对应别名, 用于修改url路径
	aliasDatabase map[string]string            //数据库别名对应的数据库名
	tableAlias    map[string]map[string]string // 数据库下的表名和别名
	aliasTable    map[string]map[string]string // 数据库下的别名和表名

	DB *DBManager
}

func (app *appAdmin) WhiteUrls() []string {
	return app.Config.whiteUrls

}
func (app *appAdmin) BlackUrls() []string {
	return app.Config.blackUrls

}

// 数据库别名设置
func (app *appAdmin) DBAlias(dbName, dbAlias string) {
	if dbAlias == "" {
		panic(DBAliasSetError)
	}
	if _, ok := app.Registry[dbName]; ok {
		app.databaseAlias[dbName] = dbAlias
		app.aliasDatabase[dbAlias] = dbName
	} else {
		fmt.Println("DBAliasSetWarning: Not found database name in registry.")
	}
}

// 表别名设置, 需要指定库
func (app *appAdmin) TBAlias(dbName, tbName, tbAlias string) {
	if tbAlias == "" {
		panic(TBAliasSetError)
	}
	tbInfo, ok := app.Registry[dbName]
	if !ok {
		fmt.Println("DBAliasSetWarning: Not found database name in registry.")
		return
	}
	_, ok = tbInfo[tbName]
	if !ok {
		fmt.Println("TBAliasSetWarning: Not found table name in registry.")
		return
	}
	if app.tableAlias[dbName] == nil {
		app.tableAlias[dbName] = map[string]string{}
	}
	app.tableAlias[dbName][tbName] = tbAlias
	if app.aliasTable[dbName] == nil {
		app.aliasTable[dbName] = map[string]string{}
	}

	app.aliasTable[dbName][tbAlias] = tbName
}

// 全局的中间件
func (app *appAdmin) AddGlobalMiddle(middles ...gin.HandlerFunc) {
	app.Config.AddGlobalMiddle(middles...)
}

// 单纯的Group中间件
func (app *appAdmin) AddGroupMiddle(middles ...gin.HandlerFunc) {
	app.Config.AddGroupMiddle(middles...)
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

func (app *appAdmin) staticFile() {
	_, file, _, ok := runtime.Caller(0)
	if ! ok {
		panic(errors.New("SugarAdminError: get template path error"))
	}
	tplPath := path.Join(path.Dir(file), "static")

	app.Sugar.Static(app.Config.Static, tplPath)

}
func (app *appAdmin) InitApp(middles ...gin.HandlerFunc) {
	app.Config.CheckParams()
	app.new(middles...)
	app.htmlGlob()
	app.staticFile()

}

// 生效session
func (app *appAdmin) UseSession(name string, store sessions.Store) {
	if store != nil {
		sessionStore = store
	}
	app.Sugar.Use(sessions.Sessions(name, store))
}

func (app *appAdmin) InitGroup(middles ...gin.HandlerFunc) *gin.RouterGroup {

	sugarGroup := app.Sugar.Group(app.Config.Prefix)
	sugarGroup.Use(middles...)
	return sugarGroup

}

func (app *appAdmin) globalMiddle() {
	app.Sugar.Use(sessions.Sessions(settings.Session.Name, sessionStore))
	if app.Config.AccessControl == "static" {
		app.Config.AddGlobalMiddle(staticGlobalMiddle())
	}
	app.Sugar.Use(app.Config.globalMiddlewares...)
}
func (app *appAdmin) groupMiddle() {
	app.Sugar.Use(sessions.Sessions(settings.Session.Name, sessionStore))
	if app.Config.AccessControl == "static" {
		app.Config.AddGroupMiddle(staticGroupMiddle())

	}
	app.Sugar.Use(app.Config.groupMiddlewares...)
}
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Set example variable
		c.Set("example", "12345")

		// before request
		log.Print("before request")

		c.Next()

		// after request
		log.Print("after request")

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}

func (app *appAdmin) Start(back bool) {
	App.DB = &Dbm
	App.globalMiddle()

	rg := App.InitGroup()
	App.groupMiddle()
	App.GroupRouter = groupRouter{group: rg, conf: App.Config}
	App.GroupRouter.init()
	rg.Use(App.Config.groupMiddlewares...)

	if back {
		go app.Sugar.Run(app.Config.Address)
	} else {
		app.Sugar.Run(app.Config.Address)
	}
}

// 注册表配置, 遍历一个目录, 数据表文件命名规则: database.table.json/yml/xml
func Register(confPath string, confType string, analy IAnalyzer) {
	confType = strings.ToLower(confType)

	changeAnalyzer(confType, analy)

	// 库名和表文件名, 用于支持跨库的表注册操作, 就要有多个数据库连接池
	var fileList []string
	fileList = readDir(confPath, fileList)

	checkDescFileList(fileList)

	if App.Registry == nil {
		fmt.Println("Warning: No form to register was found.")
		return
	}

	for dbName, tbInfo := range App.Registry {
		for _, tb := range tbInfo {
			fmt.Println(dbName, tb)
		}
	}
}

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

// 对注册表json做解析
func checkDescFileList(fileList []string) {
	dbList := Dbm.showDatabase()
	tmpTables := map[string][]string{}
	for _, file := range fileList {
		dbName, tbName, err := defaultAnalyzer.verifyPath(file)
		if err == TableConfBackupWarning {
			continue
		} else {
			utils.PanicCheck(err)
		}
		// 验证数据库
		if !utils.InStringSlice(dbName, dbList) {
			errStr := fmt.Sprintf("RegisterDatabaseNotFound: Register database name '%s' not found.", dbName)
			utils.PanicCheck(errors.New(errStr))
		}
		// 更新一个表的pool
		Dbm.UpdateDBPool(dbName)
		if _, ok := tmpTables[dbName]; ! ok {
			tmpTables[dbName] = Dbm.showTables(dbName)
		}

		// 验证数据表
		tbList := tmpTables[dbName]
		if ! utils.InStringSlice(tbName, tbList) {
			errStr := fmt.Sprintf("RegisterTableNotFound: Register table name '%s' not found.", tbName)
			panic(errors.New(errStr))
		}

		if _, ok := App.Registry[dbName]; ! ok {
			App.Registry[dbName] = make(map[string]*descConf)
		}

		dc, err := defaultAnalyzer.dump()
		utils.PanicCheck(err)
		checkDesc(tbName, dbName, dc)

		App.Registry[dbName][tbName] = dc
	}
}

// 对解析后的注册表做基本数据校验修改
func checkDesc(tbName, dbName string, dc *descConf) {
	// 以文件名的表明为准
	dc.Name = tbName
	if len(dc.Methods) == 0 && App.Config.AccessControl != "rbac" {
		dc.Methods = append(dc.Methods, methods...)
	}

	if dc.Filter == nil {
		// todo 这里根据字段类型生成过滤器
		dc.Filter = map[string]FilterInfo{}

	}
	if dc.DescType == nil {
		dc.DescType = map[string]string{}
	}
	if dc.Foreign == nil {
		dc.Foreign = map[string]string{}

	}
	updateDesc(dbName, dc)

}

// 检验字段配置是否为空, 为空的话进行全填充
func verifyDescField(dc *descConf) bool {
	if len(dc.Field) != len(dc.Title) {
		errStr := fmt.Sprintf("RegisterTableFieldError: Register table name '%s' field lendth error.", dc.Name)
		panic(errors.New(errStr))
	}
	// todo 这里的长度问题和表结构的顺序问题
	if len(dc.Field) == 0 {
		return true
	}
	return false
}

// 更新相关字段, 主键、和表结构
func updateDesc(dbName string, dc *descConf) {

	fieldStatus := verifyDescField(dc)

	sqlCmd := `select COLUMN_NAME,DATA_TYPE,CHARACTER_MAXIMUM_LENGTH, 
			   COLUMN_KEY from information_schema.COLUMNS
			   where table_schema=? AND table_name=?`
	stmt, err := Dbm.DefaultDB.Prepare(sqlCmd)
	utils.PanicCheck(err)

	result, err := Dbm.SelectSlice(stmt, dbName, dc.Name)
	utils.PanicCheck(err)

	for _, line := range result {
		if line[3] == "PRI" && dc.Primary == "" {
			dc.Primary = line[0]
		}
		dc.DescField = append(dc.DescField, line[0])
		dc.DescType[line[0]] = line[1]
		if fieldStatus {
			dc.Field = append(dc.Field, line[0])
			dc.Title = append(dc.Title, line[0])
		}
	}
	updateDescLeft(dc)
	updateDescRight(dc)
}
func updateDescLeft( dc *descConf){
	if dc.Left {
		// todo 验证左侧html代码合法性
		if dc.LeftHtml == ""{
			dc.LeftHtml = template.HTML(`<label><input class="checkall" type="checkbox" value="all">选择</label>`)
		}

	}

}
func updateDescRight( dc *descConf){
	if dc.Right {
		// todo 验证左侧html代码合法性
		if dc.RightHtml == ""{
			dc.RightHtml = template.HTML(`操作`)
		}

	}
}

func SetAdmin(conf *AdminConf) {
	App.Config = conf
	App.InitApp(App.Config.globalMiddlewares...)
	App.AddGlobalMiddle(gin.Logger(), gin.Recovery(), GetMenu())

}

func SetAuthenticate(handle digestHandler) {
	if handle != nil {
		handle(App.Config)
	}
}

type digestHandler func(ac *AdminConf)

func init() {
	// 基本的配置文件
	Settings("")
	pluginInit()
	App = appAdmin{
		Registry:      map[string]map[string]*descConf{},
		databaseAlias: map[string]string{},
		aliasDatabase: map[string]string{},
		tableAlias:    map[string]map[string]string{},
		aliasTable:    map[string]map[string]string{},
	}
}

func pluginInit() {
	// 数据库连接池初始化
	DBMInit(settings.DBConfig)
	fmt.Println(settings.Session)
	InitSession(settings.Session)
}
