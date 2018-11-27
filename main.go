//
// __author__ = "Miller"
// Date: 2018/11/15
//

package main

import (
	"fmt"
	"miller-blogs/sugar/curd"
	"miller-blogs/urls"
)

func main() {
	dbc := curd.DBConfig{
		Driver:   "mysql",
		Username: "root",
		Password: "woaichenni",
		Protocol: "tcp",
		Address:  "127.0.0.1:3306",
		DBName:   "supervision",
	}
	curd.DbmInit(dbc)
	//"root:woaichenni@tcp(127.0.0.1:3306)/supervision?charset=utf8"
	type sss struct {
		xx string
		pp string
	}
	x := &curd.TableConf{Name:"app_list",Desc:sss{}}

	curd.AccessControl("rbac")
	curd.Register(x)
	curd.AppInit(urls.AdApp,"",nil)
	err := urls.AdApp.Run("0.0.0.0:9090")
	fmt.Println(err)
}
