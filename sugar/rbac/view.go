/*
# __author__ = "Mr.chai"
# Date: 2018/12/21
*/
package rbac

import "github.com/gin-gonic/gin"

func handleVerifyLogin(c *gin.Context) {


}

func handleLogin(c *gin.Context) {


}

type SortedMenu [] Menu


type Menu struct {
	Id int
	Name string
	Url string
	Icon string
	Children SortedMenu
	Sort int
	IsMenu bool
	IsRegex bool
}

func xxx()  {
	//m := Menu{Children:SortedMenu{}}
	//{{}}
}