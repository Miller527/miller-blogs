//
// __author__ = "Miller"
// Date: 2018/11/15
//

package main

import (
	_ "miller-blogs/models"
	"miller-blogs/sugar"
)


func main() {
	 sugarApp := sugar.SugarAdmin{AccessControl:"rbac", Address:"0.0.0.0:9090",Prefix:"sugar",Extend:""}
	 sugarApp.Start(false)

	//sugar.AccessControl("static")

	//urls.AdApp.LoadHTMLGlob("sugar/rbac/*")

	//sugar.AppInit(urls.AdApp,"","")
	//go urls.AdApp.Run("0.0.0.0:9090")
}
