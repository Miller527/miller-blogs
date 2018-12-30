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
	"miller-blogs/sugar"
	"miller-blogs/sugar/utils"
	"net/http"
	"regexp"
)

// session中存储生成的所有的前端左侧菜单代码，不依赖于原来的内容
// 存储路由规则、静态和动态





func RbacLoginMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("RbacLoginMiddle start")

		session := sessions.Default(c)

		url := c.Request.RequestURI
		loginUrl := ParamsRbac.loginUrl

		if utils.InStringSlice(url,ParamsRbac.blackList){
			// todo 无权限页面
			c.JSON(http.StatusForbidden, ResMsg(403,"无访问权限"))
			c.Abort()
			return
		}
		perStr := session.Get("permission")
		menuStr := session.Get("menu")
		per := Permissions{}
		me := sugar.SortedMenu{}

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
			if menuStr != nil{
				err = json.Unmarshal([]byte(menuStr.(string)), &me)
				x, e := json.Marshal(me)
				fmt.Println("xxxxxxxxxxxxxxxxxxx", e, string(x),me)

			}

			// 静态、动态路由校验
			if !utils.InStringSlice(url, per.Static) && ! regexUrlVerify(url, per.Regex){
				c.JSON(http.StatusForbidden, ResMsg(403,"无访问权限"))
				c.Abort()
				return
			}
			fmt.Println("cccccccccccccccccccccccccc")
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
				fmt.Println(len(per.Static) > 1 || len(per.Regex) >0)
					c.Redirect(http.StatusFound,ParamsRbac.indexUrl)
					c.Abort()
				return
			}
		}


		c.Next()


		log.Println("RbacLoginMiddle end")
	}
}

func regexUrlVerify(url string, regs []string) bool{
	byteUrl := []byte(url)
	for _, r := range regs {
		status, err := regexp.Match(r, byteUrl)
		if err != nil{
			// todo 记录日志输出检查错误
			fmt.Println(err)
		}
		if status{
			return true
		}
	}
	return false
}


func BehaviorLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("BehaviorLog start ")

		c.Next()

		log.Println("BehaviorLog end ")

	}
}