//
// __author__ = "Miller"
// Date: 2018/11/15
//

package main

import (
	//_ "miller-blogs/models"
	//_ "miller-blogs/settings"
	"miller-blogs/sugar"
)


func main() {
	 sugarApp := sugar.SugarAdmin{Relative:""}
	 sugarApp.Init()
	 sugarApp.Sugar.Run("0.0.0.0:9090")
	//sugar.AccessControl("static")

	//urls.AdApp.LoadHTMLGlob("sugar/rbac/*")

	//sugar.AppInit(urls.AdApp,"","")
	//go urls.AdApp.Run("0.0.0.0:9090")
	//time.Sleep(10000*time.Second)
}
