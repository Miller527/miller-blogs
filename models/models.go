/*
# __author__ = "Mr.chai"
# Date: 2018/11/27
*/
package models

import (
	"fmt"
	"miller-blogs/sugar/curd"
	"reflect"
)

type Role struct {
	Id   int
	Rid  string
	Name string
}



func init() {
	dbc := curd.DBConfig{
		Driver:   "mysql",
		Username: "root",
		Password: "woaichenni",
		Protocol: "tcp",
		Address:  "127.0.0.1:3306",
		DBName:   "miller_blogs",
	}
	curd.DbmInit(dbc)

	{
		r := Role{}
		reflect.TypeOf(r).String()
		fmt.Println()
	}

	x := &curd.TableConf{
		Slice:[]interface{}{"id","rid","name"},
		Title:[]string{"ID","角色ID","角色名称"} ,
		Desc: Role{},
	}

	fmt.Println(x.Name())
	curd.AccessControl("rbac")
	curd.Register(x)

}
