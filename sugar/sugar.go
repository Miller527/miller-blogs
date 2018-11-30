//
// __author__ = "Miller"
// Date: 2018/11/24
//

package sugar

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"runtime"
	"strings"
)


type SugarAdmin struct {
	AccessControl string
	Address string
	Prefix string
	Relative string
	Sugar *gin.Engine
}

func (sa *SugarAdmin)checkParams() {
	if sa.AccessControl == ""{
		sa.AccessControl = "rbac"
	}else if  sa.AccessControl != "rbac" && sa.AccessControl != "static"{
		panic(errors.New("SugarAdminError: Access control type must 'rbac' or 'static'"))
	}

	if  sa.Address == "" {
		sa.Address = "0.0.0.0:9090"
	}

	sa.checkPrefix()
	sa.checkRelative()
}

//func relativePathVerify(relativeKey string)string{

//	return relativeKey
//}

func (sa *SugarAdmin)checkRelative() {
		if sa.Relative == "" {
			if sa.AccessControl == "rbac"{
				sa.Relative = ":tablename/"
			}else {
				sa.Relative = "tablename/"

			}
		}else {
			for _, b := range sa.Relative{
				if b < 'a' || b > 'Z' || b != '_'{
					panic(errors.New("SugarAdminError: Relative only be case letters"))
				}
			}
			sa.Relative+="/"
		}
}
func (sa *SugarAdmin)checkPrefix() {
	if sa.Prefix == "" {
		sa.Prefix =  "/sugar/"
	}
	if ! strings.HasPrefix(sa.Prefix, "/") {
		sa.Prefix = "/" + sa.Prefix
	}
	if ! strings.HasSuffix(sa.Prefix, "/") {
		sa.Prefix += "/"
	}
}

func (sa *SugarAdmin)new(middleware ...gin.HandlerFunc){
	sa.Sugar = gin.New()
	sa.Sugar.Use(middleware...)
}

func (sa *SugarAdmin)htmlGlob(){
	_, file, _, ok :=runtime.Caller(0)
	if ! ok {
		fmt.Println()
		panic(errors.New("SugarAdminError: get template path error"))
	}
	tplPath := path.Join(path.Dir(file), "template","*")
	sa.Sugar.LoadHTMLGlob(tplPath)

}
func (sa *SugarAdmin)Init(middleware ...gin.HandlerFunc){
	sa.checkParams()
	sa.new(middleware...)
	sa.htmlGlob()


	 sa.Sugar.GET("/verify-login", verifyLogin)

}


func verifyLogin(c *gin.Context) {

	c.HTML(http.StatusOK, "login.html", gin.H{
	})
}






//
//var TableRegister = make(map[string]*TableConf)
//
//
//
//var accessControlType = "rbac"
//
//// 配置表接口
//type SugarTable interface {
//	verifyName() (string, bool)
//}
//
//// todo struct 定义表结构、考虑怎么动态取值的方式
//
//type TableConf struct {
//	Field   []string
//	Title   []string
//	Methods []string
//	Desc    interface{}
//}
//
//func (tc *TableConf) Name() string {
//	tmpSlice := strings.Split(reflect.TypeOf(tc.Desc).String(), ".")
//	return utils.SnakeString(tmpSlice[len(tmpSlice)-1])
//}
//
//func (tc *TableConf) PrefixName(pre string) {
//	//tc.Name = pre + tc.Name
//	//orm.RegisterModel()
//}
//
//func verifyName(tc *TableConf) bool {
//	if ! InSlice(tc.Name(), tables) {
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
//
//	type desc struct {
//		name     string
//		dataType string
//	}
//	column := &TableConf{
//		Field: []string{"name", "dataType"},
//		Desc:  &desc{},
//	}
//
//	result, err := Dbm.SelectSlice(stmt, column, Dbm.Conf.DBName, tc.Name())
//	if err != nil {
//		fmt.Println("verifyField", result, err)
//		return false
//	}
//
//	for _, line := range result {
//		if ! InSlice(line[0].(string), tc.Field) {
//			return false
//		}
//	}
//	return true
//}
//
//func NewTable() {
//
//}
//
//// 注册表配置
//func Register(tcList ...*TableConf) {
//	for _, tc := range tcList {
//		name := tc.Name()
//		if ! verifyField(tc) {
//			panic(errors.New("SugarTable: Table [" + name + "] Field error"))
//		}
//		if ! verifyName(tc) {
//			panic(errors.New("SugarTable: database not found [" + name + "] table"))
//		}
//		if _, ok := TableRegister[name]; ok {
//			panic(errors.New("SugarTable: table [" + name + "] has already registered"))
//		}
//
//		TableRegister[name] = tc
//	}
//}
//var pathKey string
//var prefixPath string
//
//var whiteList []string
//
//var blackList []string
//
//func appendWhite(prefixPath string){
//	//whiteList = append(whiteList,prefixPath+"/login")
//}
//
//func appendblack(prefixPath string){
//
//}
//
//

//

//
//
// 根据注册表增加路由数据
//func AppInit(app *gin.Engine, prefix string, relativeKey string, middlewares ...gin.HandlerFunc) {
//
//	pathKey = relativePathVerify(relativeKey)
//	sugarGroup := app.Group(prefixPath)
//
//	for _, middle := range middlewares{
//		sugarGroup.Use(middle)
//	}
//	curd.initGroup(sugarGroup,pathKey)
//}
