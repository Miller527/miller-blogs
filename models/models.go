/*
# __author__ = "Mr.chai"
# Date: 2018/11/27
*/
package models

import (
	"miller-blogs/sugar/curd"
)

type Role struct {
	id   int
	rid  string
	name string
}

func init() {

	curd.DbmInit()
	curd.AccessControl("rbac")

	role := &curd.TableConf{
		Field: []string{"id", "rid", "name"},
		Title: []string{"ID", "角色ID", "角色名称"},
		Desc:  &Role{},
	}

	curd.Register(role)

}
