////
//// __author__ = "Miller"
//// Date: 2018/11/24
////
//
package sugar
//
//import (
//	"errors"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"miller-blogs/sugar/utils"
//	"net/http"
//	"strings"
//)
//
//
//const (
//	List = iota
//	Get
//	Add
//	Update
//	Delete
//	//MulitDelete
//	//MulitAdd
//	//MulitUpdate
//)
//
//var methods = []int{
//	List,
//	Get,
//	Add,
//	Update,
//	Delete,
//	//MulitDelete,
//	//MulitAdd,
//	//MulitUpdate,
//}
//
//type groupRouter struct {
//	conf Config
//	group *gin.RouterGroup
//	login gin.HandlerFunc
//}
//
//// if gr.conf.AccessControl == "rbac" {
////			if ! utils.InStringSlice(c.Param(gr.conf.Prefix), tables) && strings.IndexAny(urlTmp, gr.conf.Extend)==-1{
////				// todo 404页面报错
////				fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
////				c.String(http.StatusNotFound, "404")
////				c.Abort()
////			}
////			}
//
//func (gr *groupRouter) staticMiddle() gin.HandlerFunc {
//	return 	func(c *gin.Context){
//			urlTmp := c.Request.URL.String()
//			fmt.Println(urlTmp)
//			fmt.Println(gr.conf.Prefix + gr.conf.Extend)
//			prifix := gr.conf.Prefix + gr.conf.Extend
//			if strings.HasPrefix(urlTmp, prifix) {
//				urlTmp := strings.Split(urlTmp, prifix)[1]
//				urlTmp = strings.Split(urlTmp, "/")[0]
//				fmt.Println(urlTmp,"xxxxxxxxx")
//
//				c.Params = append(c.Params, gin.Param{gr.conf.relativeKey, urlTmp})
//			}
//			if utils.InStringSlice(urlTmp, gr.conf.whiteUrls) {
//				c.Next()
//				return
//			}
//			if utils.InStringSlice(urlTmp, gr.conf.blackUrls) {
//				// todo 404页面报错
//				c.String(http.StatusNotFound, "404")
//				c.Abort()
//		}
//	}
//}
//func (gr *groupRouter) initMiddle(){
//	if gr.conf.AccessControl == "static" {
//		gr.conf.groupMiddlewares = append(gr.conf.groupMiddlewares, gr.staticMiddle())
//	}
//	gr.group.Use(gr.conf.groupMiddlewares...)
//
//}
//
//func (gr *groupRouter) init() {
//	gr.initMiddle()
//	gr.staticRouter()
//	gr.dynamicRouter()
//}
//
//func (gr *groupRouter) staticRouter() {
//
//	if gr.conf.loginFunc != nil{
//		gr.group.POST("/login", gr.conf.loginFunc )
//
//	}else {
//		gr.group.POST("/login", HandlerVerifyLogin)
//	}
//	gr.group.GET("/login", HandlerLogin)
//	gr.group.GET("/tables", HandlerCurd)
//	gr.group.GET("/index", HandlerIndex)
//	gr.group.GET("/index.html", HandlerIndex)
//}
//
//
//
//// 通过配置限制限制访问权限, 不分用户
//func (gr *groupRouter) groupRouter(relative string , methodList []int) {
//	extend := gr.conf.Extend
//	for _, method := range methodList {
//		switch method {
//		case List:
//			// 单表查询
//			gr.group.GET(extend+relative+"/list", HandlerList)
//		case Get:
//			// 单行查询详情
//			gr.group.GET(extend+relative+"/get/:id", HandlerGet)
//
//		case Add:
//			// 添加一行或多行
//			gr.group.POST(extend+relative+"/add", HandlerAdd)
//		case Update:
//			//更新一行或多行
//			gr.group.PUT(extend+relative+"/update", HandlerUpdate)
//		case Delete:
//			//删除一行或多行
//			gr.group.DELETE(extend+relative+"/delete", HandlerDelete)
//		//
//		//case MulitAdd:
//		//	gr.group.DELETE(extend+relative+"/delete/:id", HandlerMulitDelete)
//		//case MulitUpdate:
//		//	gr.group.PUT(extend+relative+"/delete", HandlerMulitUpdate)
//		//case MulitDelete:
//		//	gr.group.PUT(extend+relative+"/mulit-update", HandlerMulitUpdate)
//		default:
//			panic(errors.New("SugarTable: table [" + relative + "] method error"))
//		}
//	}
//}
//func (gr *groupRouter) dynamicRouter() {
//	// 通过rbac控制访问权限
//	fmt.Println("conf.AccessControl",gr.conf.AccessControl)
//
//	if gr.conf.AccessControl == "rbac" {
//		gr.groupRouter(gr.conf.Relative,  methods)
//		return
//	}
//	// 通过配置控制访问权限, 路由信息会很多
//	for k, tc := range App.registry {
//		if gr.conf.AccessControl == "static" || tc.Methods == nil {
//			gr.groupRouter(k, methods)
//			continue
//		}
//		for _, v := range tc.Methods {
//			if ! utils.InIntSlice(v, methods) {
//				panic(errors.New("SugarTable: table [" + tc.Name(nil) + "] method error"))
//			}
//		}
//		gr.groupRouter(k, tc.Methods)
//	}
//}
