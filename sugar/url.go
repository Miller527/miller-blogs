//
// __author__ = "Miller"
// Date: 2018/11/24
//

package sugar

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"miller-blogs/sugar/curd"
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

type SugarRouter struct {
	AccessControl string
	Prefix        string
	Extend        string
	Relative      string
	WhiteUrls     []string
	BlackUrls     []string
}

func (sr *SugarRouter) initMiddle() gin.HandlerFunc {

	return func(c *gin.Context) {
		urlTmp := c.Request.URL.String()

		if sr.AccessControl == "static" {

			if strings.HasPrefix(urlTmp, sr.Prefix) {
				urlTmp := strings.Split(urlTmp, sr.Prefix)[1]
				urlTmp = strings.Split(urlTmp, "/")[0]
				c.Params = append(c.Params, gin.Param{sr.Prefix, urlTmp})
			}
			if utils.InStringSlice(urlTmp, sr.WhiteUrls) {
				fmt.Println("xxxxxxxxxxxxx")
				c.Next()
				return
			}
			if utils.InStringSlice(urlTmp, sr.BlackUrls) {
				// todo 404页面报错
				c.String(http.StatusNotFound, "404")
				c.Abort()

			}
		}

		if sr.AccessControl == "rbac" && ! utils.InStringSlice(c.Param(sr.Prefix), tables) {
			// todo 404页面报错
			c.String(http.StatusNotFound, "404")
			c.Abort()
		}

	}
}

func (sr *SugarRouter) Router(rg *gin.RouterGroup) {
	sr.addRouter(rg)
	sr.initUrl(rg)
}

func (sr *SugarRouter) addRouter(rg *gin.RouterGroup) {
	rg.GET("/login", curd.Login)
	rg.GET("/curd", curd.Curd)
	rg.POST("/login", curd.VerifyLogin)
	rg.GET("/index", curd.Index)
	rg.GET("/index.html", curd.Index)
}

// 通过配置限制限制访问权限, 不分用户
func (sr *SugarRouter) groupRouter(relative string, rg *gin.RouterGroup, methodList []int) {
	for _, method := range methodList {
		switch method {
		case List:
			rg.GET(sr.Extend+relative+"/list", curd.List)
		case Get:
			rg.GET(sr.Extend+relative+"/get/:id", curd.Get)

		case Add:
			rg.POST(sr.Extend+relative+"/add", curd.Add)
		case Update:
			rg.PUT(sr.Extend+relative+"/update", curd.Update)
		case Delete:
			rg.DELETE(sr.Extend+relative+"/update", curd.Delete)

		case MulitAdd:
			rg.DELETE(sr.Extend+relative+"/delete/:id", curd.MulitDelete)
		case MulitUpdate:
			rg.PUT(sr.Extend+relative+"/mulit-delete", curd.MulitUpdate)
		case MulitDelete:
			rg.PUT(sr.Extend+relative+"/mulit-update", curd.MulitUpdate)
		default:
			panic(errors.New("SugarTable: table [" + relative + "] method error"))
		}
	}
}
func (sr *SugarRouter) initUrl(rg *gin.RouterGroup) {
	// 通过rbac控制访问权限
	if sr.AccessControl == "rbac" {
		sr.groupRouter(sr.Relative, rg, methods)
		return
	}
	// 通过配置控制访问权限, 路由信息会很多
	for k, tc := range Registry {
		if sr.AccessControl == "rbac" || tc.Methods == nil {
			sr.groupRouter(k, rg, methods)
			continue
		}
		for _, v := range tc.Methods {
			if ! utils.InIntSlice(v, methods) {
				panic(errors.New("SugarTable: table [" + tc.Name() + "] method error"))
			}
		}
		sr.groupRouter(k, rg, tc.Methods)
	}
}
