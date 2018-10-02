//
// __author__ = "Miller"
// Date: 2018/10/1
//
package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	"miller-blogs/controllers/base"
	"miller-blogs/public"
	"strings"
)

// app management background white list
var urlWhiteList []string

// app blog black list
var urlBlackList []string

// login session filter
func ManagerLoginFilter(ctx *context.Context) {
	_, ok := ctx.Input.Session("uid").(string)
	if ok && ctx.Request.RequestURI == "/manager/login" {
		ctx.Redirect(302, "/manager/index")
	} else if !ok && ctx.Request.RequestURI != "/manager/login" {
		ctx.Redirect(302, "/manager/login")
	}
}

// permission verify filter
func PermissionVerify(ctx *context.Context) {
	url := ctx.Request.RequestURI
	if _, blackOk := public.ElementInList(urlBlackList, ctx.Request.RequestURI); blackOk {
		responseData := base.ResponseMsg{}
		responseData.Status = 30200
		responseData.Msg = "无权访问"
		ctx.Output.JSON(responseData, true, false)
		return
	}
	_, whiteOk := public.ElementInList(urlWhiteList, url)
	if ! whiteOk && strings.HasPrefix(url, "/manager") {
		permissionsStr, ok := ctx.Input.Session("permissions").(string)

		fmt.Println(permissionsStr, ok)
		var permissionsData []orm.ParamsList
		err := json.Unmarshal([]byte(permissionsStr), &permissionsData)

		responseData := base.ResponseMsg{}
		if ! ok || err != nil {
			responseData.Status = 40200
			responseData.Msg = "访问权限异常"
			ctx.Output.JSON(responseData, true, false)
			return
		}
		status := false
		for _, val := range permissionsData {
			fmt.Println(val[1], ctx.Request.RequestURI)
			if ctx.Request.RequestURI == val[1] {
				status = true
				break
			}
		}
		if ! status {
			responseData.Status = 30200
			responseData.Msg = "无权访问"
			ctx.Output.JSON(responseData, true, false)
			return
		}
	}

}

// initialization
func getUrlRequestList() {
	urlWhiteList = beego.AppConfig.Strings("urlwhitelist")
	urlBlackList = beego.AppConfig.Strings("urlblacklist")
}

func init() {
	getUrlRequestList()
}
