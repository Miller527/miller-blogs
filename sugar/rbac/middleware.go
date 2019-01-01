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
	"regexp"
	"strings"
)

// session中存储生成的所有的前端左侧菜单代码，不依赖于原来的内容
// 存储路由规则、静态和动态





func RbacLoginMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("RbacLoginMiddle start")

		session := sessions.Default(c)

		url := strings.Split(c.Request.RequestURI, "?")[0]
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

			// 静态、动态路由校验
			if !utils.InStringSlice(url, per.Static) && ! regexUrlVerify(url, per.Regex){
				c.JSON(http.StatusForbidden, ResMsg(403,"无访问权限"))
				c.Abort()
				return
			}
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
		//log.Println("BehaviorLog start ")

		c.Next()

		//log.Println("BehaviorLog end ")

	}
}





func groupParamsMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("rbacGroupMiddle start")

		urlTmp := c.Request.URL.String()
		fmt.Println(ParamsRbac.urlPrefix)
		fmt.Println(urlTmp)
		fmt.Println(c.Request.RequestURI)
		fmt.Println(c.Request.Host)
		//if utils.InStringSlice(urlTmp, ParamsRbac.whiteList) {
		//	c.Next()
		//	return
		//}
		//if utils.InStringSlice(urlTmp, ParamsRbac.blackList) {
		//	// todo 404页面报错
		//	c.String(http.StatusNotFound, "404")
		//	c.Abort()
		//}
		prefix := ParamsRbac.urlPrefix
		if strings.HasPrefix(urlTmp, prefix) {
			urlTmp = urlTmp[len(prefix):]
			urlFields := strings.Split(urlTmp, "/")
			if len(urlFields) == 1{
				dbTmp := strings.Split(urlFields[0],"-")

				c.Params = append(c.Params,
					gin.Param{ParamsRbac.extendKey, dbTmp[0]},
				)
			}else if len(urlFields) > 1{
				c.Params = append(c.Params,
					gin.Param{ParamsRbac.extendKey, urlFields[0]},
					gin.Param{ParamsRbac.relativeKey, urlFields[1]},
				)
			}
			fmt.Println("urlFields",urlFields)

		}

		c.Next()
		//
		//
		//fmt.Println(urlFields)
		//fmt.Println(len(urlFields))



		log.Println("rbacGroupMiddle end")
	}
}

func globalParamsMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("rbacGlobalMiddle start")

		session := sessions.Default(c)
		menuJson := session.Get("menu")

		if menuJson != nil {
			c.Set("menu",menuJson.(string))
		}else{
			fmt.Println("GetMenu set error ")
		}
		c.Next()


		log.Println("rbacGlobalMiddle end")
	}
}