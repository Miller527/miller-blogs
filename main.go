//
// __author__ = "Miller"
// Date: 2018/11/15
//

package main

import (
	"miller-blogs/sugar"
	"miller-blogs/sugar/rbac"
)

func main() {

	//conf := &sugar.AdminConf{Prefix: ,AccessControl:""}
	sugar.OriginalAdminConf.Prefix = "sugar"
	sugar.OriginalAdminConf.AccessControl = "rbac"
	sugar.SetAdmin()

	sugar.Settings("settings/config.json")
	sugar.Register("models", "json", nil)
	sugar.App.DBAlias("miller_blogs", "blogs")
	sugar.App.TBAlias("miller_blogs", "role", "role")
	//rbac.BlackList("/manager/index","/favicon.ico")
	rbac.ParamsRbac.SetAdmin(sugar.OriginalAdminConf)
	rbac.ParamsRbac.WhiteList("/sugar/login")
	sugar.SetAuthenticate(rbac.Register)
	sugar.App.Start(false)

	//urls.AdApp.LoadHTMLGlob("sugar/rbac/*")
	//sugar.AppInit(urls.AdApp,"","")
	//go urls.AdApp.Run("0.0.0.0:9090")
}
