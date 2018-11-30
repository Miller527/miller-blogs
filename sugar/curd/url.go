//
// __author__ = "Miller"
// Date: 2018/11/24
//

package curd

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)
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

func initMiddle() gin.HandlerFunc {

	return func(c *gin.Context) {
		urlTmp := c.Request.URL.String()
		if InSlice(urlTmp, whiteList){
			fmt.Println("xxxxxxxxxxxxx")
			c.Next()
		}
		if accessControlType == "rbac" && ! InSlice(c.Param(pathKey), tables) {
			// todo 404页面报错
			c.String(http.StatusNotFound, "404")
			c.Abort()
		}

		if accessControlType == "static" {

			if strings.HasPrefix(urlTmp, prefixPath) {
				urlTmp := strings.Split(urlTmp, prefixPath)[1]
				urlTmp = strings.Split(urlTmp, "/")[0]
				c.Params = append(c.Params, gin.Param{pathKey, urlTmp})
			}
		}

	}
}

func addRouter(rg *gin.RouterGroup){
	rg.GET("/login", Login)
	rg.GET("/verify-login", verifyLogin)
}

func initGroup(rg *gin.RouterGroup, pathKey string) {

	rg.Use(initMiddle())
	addRouter(rg)
	// 通过rbac控制访问权限
	if accessControlType == "rbac" {
		groupRouter("/:"+pathKey, rg, methods)
		return
	}
	// 通过配置控制访问权限, 路由信息会很多
	for k, tc := range TableRegister {
		if accessControlType == "rbac" || tc.Methods == nil {
			groupRouter(k, rg, methods)
			continue
		}
		for _, v := range tc.Methods {
			if ! InSlice(v, methods) {
				panic(errors.New("SugarTable: table [" + tc.Name() + "] method [" + v + "] is error"))
			}
		}
		groupRouter(k, rg, tc.Methods)
	}
}

// 通过配置限制限制访问权限, 不分用户
func groupRouter(pathKey string, rg *gin.RouterGroup, methodList []string) {
	for _, method := range methodList {
		switch method {
		case "Index":
			rg.GET(pathKey + "/list", Index)
		case "Get":
			rg.GET(pathKey+"/get/:id", Get)
		case "Add":
			rg.POST(pathKey+"/add", Add)
		case "Update":
		case "Delete":
			rg.PUT(pathKey+"/update", Update)
		case "MulitDelete":
			rg.PUT(pathKey+"/mulit-update", MulitUpdate)
		case "MulitAdd":
			rg.DELETE(pathKey+"/delete/:id", Delete)
		case "MulitUpdate":
			rg.DELETE(pathKey+"/mulit-delete", MulitDelete)
		default:
			panic(errors.New("SugarTable: table [" + pathKey + "] method [" + method + "] is error"))
		}
	}
}
