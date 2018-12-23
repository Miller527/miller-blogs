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
	sugar.SetAdmin(&sugar.AdminConf{Prefix: "manager",AccessControl:"rbac"})

	sugar.Settings("settings/config.json")
	sugar.Register("E:/GoProject/miller-blogs/src/miller-blogs/models", "json", nil)
	sugar.App.DBAlias("miller_blogs", "blogs")
	sugar.App.TBAlias("miller_blogs", "role", "userrole")
	sugar.SetAuthenticate(rbac.Register, true)
	sugar.App.Start(false)

	//urls.AdApp.LoadHTMLGlob("sugar/rbac/*")
	//sugar.AppInit(urls.AdApp,"","")
	//go urls.AdApp.Run("0.0.0.0:9090")
}




//fmt.Println(sessionStore)
//session := sessions.Default(c)
//v := session.Get("count")
//if v == nil {
//count = 0
//} else {
//count = v.(int)
//count++
//}
//session.Set("count", count)
//session.Save()