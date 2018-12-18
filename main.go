//
// __author__ = "Miller"
// Date: 2018/11/15
//

package main

import (
	"miller-blogs/sugar"
)

func main() {
	sugar.Config("settings/config.json")
	sugar.Register("E:/GoProject/miller-blogs/src/miller-blogs/models", "json",nil)
	sugar.App.DBAlias("miller_blogs","blogs")
	sugar.App.TBAlias("miller_blogs","role","userrole")
	//sugar.SetAdmin(sugar.Config{Prefix:"sugar"})
	//sugar.App.Start(false)

	//urls.AdApp.LoadHTMLGlob("sugar/rbac/*")
	//
	//sugar.AppInit(urls.AdApp,"","")
	//go urls.AdApp.Run("0.0.0.0:9090")
}
