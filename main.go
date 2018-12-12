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


	sugar.SetAdmin(sugar.Config{Prefix:"sugar"})
	sugar.App.Start(false)

	//urls.AdApp.LoadHTMLGlob("sugar/rbac/*")
	//
	//sugar.AppInit(urls.AdApp,"","")
	//go urls.AdApp.Run("0.0.0.0:9090")
}
