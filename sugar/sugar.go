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

type SugarAdmin struct {
	AccessControl string
	Address       string
	Prefix        string
	Extend        string
	Relative      string
	//Static        string
	Sugar         *gin.Engine
	whiteUrls     []string
	blackUrls     []string
}

func (sa *SugarAdmin) AddWhite(urlSlice ...string) {
	for _, urlTmp := range urlSlice {
		sa.whiteUrls = append(sa.whiteUrls, urlTmp)
	}

}
func (sa *SugarAdmin) AddBlack(urlSlice ...string) {
	for _, urlTmp := range urlSlice {
		sa.blackUrls = append(sa.blackUrls, urlTmp)
	}
}
func (sa *SugarAdmin) checkParams() {
	if sa.AccessControl == "" {
		sa.AccessControl = "rbac"
	} else if sa.AccessControl != "rbac" && sa.AccessControl != "static" {
		panic(errors.New("SugarAdminError: Access control type must 'rbac' or 'static'"))
	}

	if sa.Address == "" {
		sa.Address = "0.0.0.0:9090"
	}

	sa.checkPrefix()
	sa.checkRelative()
	sa.checkExtend()
	//sa.checkStatic()
}

func (sa *SugarAdmin) checkRelative() {
	if sa.Relative == "" {
		if sa.AccessControl == "rbac" {
			sa.Relative = ":tablename/"
		} else {
			sa.Relative = "tablename/"

		}
	} else {
		for _, b := range sa.Relative {
			if b < 'a' || b > 'Z' || b != '_' {
				panic(errors.New("SugarAdminError: Relative only be case letters"))
			}
		}
		sa.Relative += "/"
	}
}
func (sa *SugarAdmin) checkPrefix() {
	if sa.Prefix == "" {
		sa.Prefix = "/"
		return
	}
	if ! strings.HasPrefix(sa.Prefix, "/") {
		sa.Prefix = "/" + sa.Prefix
	}
	if ! strings.HasSuffix(sa.Prefix, "/") {
		sa.Prefix += "/"
	}
}
func (sa *SugarAdmin) checkExtend() {
	if sa.Extend == "" {
		sa.Extend = "curd/"
		return
	}
	if strings.HasPrefix(sa.Extend, "/") {
		sa.Extend = sa.Extend[1:]
	}
	if ! strings.HasSuffix(sa.Extend, "/") {
		sa.Extend += "/"
	}
}

// todo 前端代码里怎么改动（预加载选择主题时候的路径问题）
//func (sa *SugarAdmin) checkStatic() {
//	if sa.Static == "" {
//		sa.Static = "static"
//		return
//	}
//	if strings.HasPrefix(sa.Static, "/") {
//		sa.Static = sa.Static[1:]
//	}
//	if ! strings.HasSuffix(sa.Static, "/") {
//		sa.Static = sa.Static[:len(sa.Static)-1]
//	}
//}

func (sa *SugarAdmin) new(middleware ...gin.HandlerFunc) {
	sa.Sugar = gin.New()
	sa.Sugar.Use(middleware...)
}

func (sa *SugarAdmin) htmlGlob() {
	_, file, _, ok := runtime.Caller(0)
	if ! ok {
		panic(errors.New("SugarAdminError: get template path error"))
	}
	tplPath := path.Join(path.Dir(file), "template", "**", "*")
	sa.Sugar.LoadHTMLGlob(tplPath)

}

func (sa *SugarAdmin) static() {
	_, file, _, ok := runtime.Caller(0)
	if ! ok {
		panic(errors.New("SugarAdminError: get template path error"))
	}
	tplPath := path.Join(path.Dir(file), "static")
	sa.Sugar.Static(sa.Prefix+"static", tplPath)
	//sa.Sugar.Static(sa.Prefix+sa.Static, tplPath)

}
func (sa *SugarAdmin) InitApp(middleware ...gin.HandlerFunc) {
	sa.checkParams()
	sa.new(middleware...)
	sa.htmlGlob()
	sa.static()

}

func (sa *SugarAdmin) InitGroup(middlewares ...gin.HandlerFunc) *gin.RouterGroup {

	sugarGroup := sa.Sugar.Group(sa.Prefix)

	for _, middle := range middlewares {
		sugarGroup.Use(middle)
	}

	return sugarGroup

}

func (sa *SugarAdmin) InitUrl(rg *gin.RouterGroup) {
	sr := SugarRouter{
		AccessControl: sa.AccessControl,
		Prefix:        sa.Prefix,
		Extend:        sa.Extend,
		Relative:      sa.Relative,
		WhiteUrls:     sa.whiteUrls,
		BlackUrls:     sa.blackUrls,
	}

	sr.Router(rg)
}

func (sa *SugarAdmin) Start(back bool) {
	sa.InitApp(gin.Logger(), gin.Recovery())
	rg := sa.InitGroup()
	sa.InitUrl(rg)
	if back {
		go sa.Sugar.Run(sa.Address)
	} else {
		sa.Sugar.Run(sa.Address)
	}
}

//// 注册表配置
func Register(tcList ...*TableConf) {
	for _, tc := range tcList {

		if len(tc.Field) != len(tc.Title) {
			panic(errors.New("SugarTable: Table field length unequal to title length"))
		}
		name := tc.Name()
		if ! verifyField(tc) {
			panic(errors.New("SugarTable: Table [" + name + "] Field error"))
		}
		if ! verifyName(name) {
			panic(errors.New("SugarTable: database not found [" + name + "] table"))
		}
		if _, ok := Registry[name]; ok {
			panic(errors.New("SugarTable: table [" + name + "] has already registered"))
		}

		Registry[name] = tc
	}
}

var Registry = make(map[string]*TableConf)

// 配置表接口
type SugarTable interface {
	Name() string
	DisplayName() string
}

type TableConf struct {
	Display     string
	DisplayJoin bool
	Field       []string
	Title       []string
	Methods     []int
	Desc        interface{}
}

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

func verifyName(name string) bool {
	if ! utils.InStringSlice(name, tables) {
		return false
	}
	return true
}

func verifyField(tc *TableConf) bool {
	sqlCmd := `select COLUMN_NAME as name,DATA_TYPE as dataType
			   from information_schema.COLUMNS
			   where table_schema=? AND table_name=?`
	stmt, err := Dbm.Db.Prepare(sqlCmd)

	type desc struct {
		name     string
		dataType string
	}
	column := &TableConf{
		Field: []string{"name", "dataType"},
		Desc:  &desc{},
	}

	result, err := Dbm.SelectSlice(stmt, column, Dbm.Conf.DBName, tc.Name())
	if err != nil {
		fmt.Println("verifyField", result, err)
		return false
	}

	for _, line := range result {
		if ! utils.InStringSlice(line[0].(string), tc.Field) {
			return false
		}
	}
	return true
}

//
//func NewTable() {
//
//}
//
// 根据注册表增加路由数据
//func AppInit(app *gin.Engine, prefix string, relativeKey string, middlewares ...gin.HandlerFunc) {
//

//}

//func (tc *TableConf) PrefixName(pre string) {
//	//tc.Name = pre + tc.Name
//	//orm.RegisterModel()
//}
