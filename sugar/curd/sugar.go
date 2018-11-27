//
// __author__ = "Miller"
// Date: 2018/11/24
//

package curd

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/gin-gonic/gin"
	"strings"
)

var TableConfig = make(map[string]*TableConf)

var methods = []string{
	"Index",
	"Get",
	"Add",
	"Update",
	"Delete",
	"MulitDelete",
	"MulitAdd",
	"MulitUpdate",
}

var AccessControlType = "rbac"

// 配置表接口
type SugarTable interface {

	verifyName() (string, bool)
}
// todo struct 定义表结构、考虑怎么动态取值的方式


type TableConf struct {
	Name    string
	Slice   string
	Title   map[string]string
	Methods []string
	Desc  interface{}

}

func (tc *TableConf) PrefixName(pre string) {
	tc.Name = pre + tc.Name
	orm.RegisterModel()
}

func (tc *TableConf) verifyName() (string, bool) {
	if ! InSlice(tc.Name, tables) {
		return tc.Name, false
	}
	return tc.Name, true
}

func (tc *TableConf) verifyTitle() bool {
	sqlCmd := "select COLUMN_NAME,DATA_TYPE from information_schema.COLUMNS where table_schema=? AND table_name=?"
	stmt, err := Dbm.Db.Prepare(sqlCmd)

	res, err := Dbm.SelectSlice(stmt,Dbm.Conf.DBName,tc.Name)
	if err != nil{
		fmt.Println(err)
	}
	for res.Next(){
		var name string

		err := res.Scan(&name)

		if err == nil{
			tables = append(tables, name)
		}
	}
	return false
}

func NewTable(){

}

// 注册表配置
func Register(st ...*TableConf) {
	for _, tb := range st {
		name, ok := tb.verifyName()
		if ! ok {
			panic(errors.New("SugarTable: database not found [" + name + "] table"))
		}
		if _, ok := TableConfig[name]; ok {
			panic(errors.New("SugarTable: table [" + name + "] has already registered"))
		}
		TableConfig[name] = tb
	}
}

// 根据注册表增加路由数据
func AppInit(app *gin.Engine, prefix string, middle gin.HandlerFunc) {
	if prefix == "" {
		prefix = "/sugar"
	}
	if ! strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	if ! strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}
	sugarGroup := app.Group(prefix)
	// todo 是否可以增加其他类型中间件
	if middle != nil {
		sugarGroup.Use(middle)
	}
	initGroup(sugarGroup)
}


func AccessControl(acType string) {
	if acType == "rbac" || acType == "static" {
		AccessControlType = acType
	} else {
		panic(errors.New("SugarControl: Access control type must 'rbac' or 'static'"))
	}
}
