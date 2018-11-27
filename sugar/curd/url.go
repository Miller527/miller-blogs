//
// __author__ = "Miller"
// Date: 2018/11/24
//

package curd

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func initGroup(rg *gin.RouterGroup) {
	// 通过rbac控制访问权限
	if AccessControlType == "rbac" {
		initRouter("tablename", rg, methods)
		return
	}
	// 通过配置控制访问权限
	for k, tc := range TableConfig {
		if AccessControlType == "rbac" || tc.Methods == nil {
			initRouter(k, rg, methods)
			continue
		}
		for _, v := range tc.Methods {
			if ! InSlice(v, methods) {
				panic(errors.New("SugarTable: table [" + tc.Name + "] method [" + v + "] is error"))
			}
		}
		initRouter(k, rg, tc.Methods)
	}
}

// 通过配置限制限制访问权限, 不分用户
func initRouter(name string, rg *gin.RouterGroup, methodList []string) {
	for _, method := range methodList {
		switch method {
		case "Index":
			rg.GET("/:"+name, Index)
		case "Get":
			rg.GET("/:"+name+"/:id", Get)
		case "Add":
			rg.POST("/:"+name+"/add", Add)
		case "Update":
			rg.POST("/:"+name+"/mulit-add", MulitAdd)
		case "Delete":
			rg.PUT("/:"+name+"/update", Update)
		case "MulitDelete":
			rg.PUT("/:"+name+"/mulit-update", MulitUpdate)
		case "MulitAdd":
			rg.DELETE("/:"+name+"/delete/:id", Delete)
		case "MulitUpdate":
			rg.DELETE("/:"+name+"/mulit-delete", MulitDelete)
		default:
			panic(errors.New("SugarTable: table [" + name + "] method [" + method + "] is error"))
		}
	}
}