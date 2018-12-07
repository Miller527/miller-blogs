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
	"net/http"
	"strings"
)


const (
	List = iota
	Get
	Add
	Update
	Delete
	MulitDelete
	MulitAdd
	MulitUpdate
)

var methods = []int{
	List,
	Get,
	Add,
	Update,
	Delete,
	MulitDelete,
	MulitAdd,
	MulitUpdate,
}

type groupRouter struct {
	conf Config
	group *gin.RouterGroup
	login gin.HandlerFunc
}

func (gr *groupRouter) initMiddle() gin.HandlerFunc {

	return func(c *gin.Context) {
		urlTmp := c.Request.URL.String()

		if gr.conf.AccessControl == "static" {

			if strings.HasPrefix(urlTmp, gr.conf.Prefix) {
				urlTmp := strings.Split(urlTmp, gr.conf.Prefix)[1]
				urlTmp = strings.Split(urlTmp, "/")[0]
				c.Params = append(c.Params, gin.Param{gr.conf.Prefix, urlTmp})
			}
			if utils.InStringSlice(urlTmp, gr.conf.whiteUrls) {
				fmt.Println("xxxxxxxxxxxxx")
				c.Next()
				return
			}
			if utils.InStringSlice(urlTmp, gr.conf.blackUrls) {
				// todo 404页面报错
				c.String(http.StatusNotFound, "404")
				c.Abort()

			}
		}

		if gr.conf.AccessControl == "rbac" && ! utils.InStringSlice(c.Param(gr.conf.Prefix), tables) {
			// todo 404页面报错
			c.String(http.StatusNotFound, "404")
			c.Abort()
		}

	}
}

func (gr *groupRouter) init() {
	gr.addRouter()
	gr.initUrl()
}

func (gr *groupRouter) addRouter() {
	gr.group.Use(gr.conf.groupMiddlewares...)
	if gr.conf.loginFunc != nil{
		gr.group.POST("/login", gr.conf.loginFunc )

	}else {
		gr.group.POST("/login", HandlerVerifyLogin)

	}
	gr.group.GET("/login", HandlerLogin)
	gr.group.GET("/tables", HandlerCurd)
	gr.group.GET("/index", HandlerIndex)
	gr.group.GET("/index.html", HandlerIndex)

}


// 通过配置限制限制访问权限, 不分用户
func (gr *groupRouter) groupRouter(relative string , methodList []int) {
	extend := gr.conf.Extend
	for _, method := range methodList {
		switch method {
		case List:
			// 单表查询
			gr.group.GET(extend+relative+"/list", HandlerList)
		case Get:
			// 单行查询详情
			gr.group.GET(extend+relative+"/get/:id", HandlerGet)

		case Add:
			// 添加一行或多行
			gr.group.POST(extend+relative+"/add", HandlerAdd)
		case Update:
			//更新一行或多行
			gr.group.PUT(extend+relative+"/update", HandlerUpdate)
		case Delete:
			//删除一行或多行
			gr.group.DELETE(extend+relative+"/delete", HandlerDelete)

		case MulitAdd:
			gr.group.DELETE(extend+relative+"/delete/:id", HandlerMulitDelete)
		case MulitUpdate:
			gr.group.PUT(extend+relative+"/delete", HandlerMulitUpdate)
		case MulitDelete:
			gr.group.PUT(extend+relative+"/mulit-update", HandlerMulitUpdate)
		default:
			panic(errors.New("SugarTable: table [" + relative + "] method error"))
		}
	}
}
func (gr *groupRouter) initUrl() {
	// 通过rbac控制访问权限
	fmt.Println("conf.AccessControl",gr.conf.AccessControl)

	if gr.conf.AccessControl == "rbac" {
		gr.groupRouter(gr.conf.Relative,  methods)
		return
	}
	// 通过配置控制访问权限, 路由信息会很多
	for k, tc := range App.registry {
		if gr.conf.AccessControl == "static" || tc.Methods == nil {
			gr.groupRouter(k, methods)
			continue
		}
		for _, v := range tc.Methods {
			if ! utils.InIntSlice(v, methods) {
				panic(errors.New("SugarTable: table [" + tc.Name() + "] method error"))
			}
		}
		gr.groupRouter(k, tc.Methods)
	}
}
