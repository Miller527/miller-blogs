/*
# __author__ = "Mr.chai"
# Date: 2018/12/21
*/
package rbac

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"miller-blogs/sugar/utils"
	"net/http"
)

// session中存储生成的所有的前端左侧菜单代码，不依赖于原来的内容
// 存储路由规则、静态和动态





func RbacLoginMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		log.Println("RbacLoginMiddle start")
		fmt.Println(c.Request.URL.String())
		fmt.Println(c.Request.URL.User)
		fmt.Println(c.Request.Host)
		fmt.Println(c.Request.RequestURI)
		fmt.Println(c.Request)
		// Process request

		url := c.Request.RequestURI
		loginUrl := ParamsRbac.loginUrl

		if utils.InStringSlice(url,ParamsRbac.blackList){
			// todo 无权限页面
			c.JSON(http.StatusForbidden, ResMsg(403,"无访问权限"))
			c.Abort()
			return
		}
		perStr := session.Get("permission")
		per := Permissions{}


		if ! utils.InStringSlice(url, ParamsRbac.whiteList){

			if perStr == nil {
				c.Redirect(http.StatusFound,loginUrl)
				c.Abort()
				return
			}
			err := json.Unmarshal([]byte(perStr.(string)), &per)

			if  err != nil {
				c.Redirect(http.StatusFound,loginUrl)
				c.Abort()
				return
			}
			// 静态路由校验
			if !utils.InStringSlice(url, per.Static) && ! regexUrlVerify(url, per.Regex){
				c.JSON(http.StatusForbidden, ResMsg(403,"无访问权限"))
				c.Abort()
				return
			}
			fmt.Println("sssssssssssssssssssssss", per.Static)
			fmt.Println("rrrrrrrrrrrrrrrrrrrrrrr", per.Regex)
		}else {
			if 	perStr != nil{
				err := json.Unmarshal([]byte(perStr.(string)), &per)

				if  err != nil {
					c.Redirect(http.StatusFound,loginUrl)
					c.Abort()
					return
				}

			}

			if url == ParamsRbac.loginUrl && (len(per.Static) > 1 || len(per.Regex) >0) {
				fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxx")
				fmt.Println(per.Static)
				fmt.Println(per.Regex)
				fmt.Println(len(per.Static) > 1 || len(per.Regex) >0)
				c.Redirect(http.StatusFound,ParamsRbac.indexUrl)
				c.Abort()
				return
			}
		}


		c.Next()
		fmt.Println(session.Get("permission"))

		log.Println("RbacLoginMiddle end")
	}
}

func regexUrlVerify(url string, regs []string) bool{
	return true
}


func BehaviorLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("BehaviorLog start ")

		c.Next()

		log.Println("BehaviorLog end ")

	}
}