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
	"github.com/astaxie/beego/utils"
	"miller-blogs/controllers/base"
	"strings"
)

// app management background white list
var urlWhiteList []string

// app blog black list
var urlBlackList []string

var urlPrefix = beego.AppConfig.String("manager_router_prefix")

// login session filter
func ManagerLoginFilter(ctx *context.Context) {
	url := ctx.Request.RequestURI
	_, ok := ctx.Input.Session("uid").(string)
	if ok &&url == urlPrefix + "/login" {
		ctx.Redirect(302, urlPrefix + "/index")
	} else if !ok && url != urlPrefix + "/login" {
		ctx.Redirect(302, urlPrefix + "/login")
	}
}

// permission verify filter
func PermissionVerify(ctx *context.Context) {
	url := ctx.Request.RequestURI

	if blackOk := utils.InSlice(url, urlBlackList); blackOk {
		responseData := base.ResponseMsg{}
		responseData.Status = 30200
		responseData.Msg = "无权访问"
		ctx.Output.JSON(responseData, true, false)
		return
	}

	if whiteOk := utils.InSlice(url, urlWhiteList); ! whiteOk && strings.HasPrefix(url, urlPrefix) {
		permissionsStr, ok := ctx.Input.Session(beego.AppConfig.String("permission_key")).(string)

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
			fmt.Println(val[1], url)
			if url == val[1] {
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
	urlWhiteList = beego.AppConfig.Strings("url_white_list")
	urlBlackList = beego.AppConfig.Strings("url_black_list")
}

func init() {
	getUrlRequestList()
}
