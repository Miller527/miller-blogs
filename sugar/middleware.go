/*
# __author__ = "Mr.chai"
# Date: 2018/12/21
*/
package sugar

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"miller-blogs/sugar/utils"
	"net/http"
	"strings"
)

// session中存储生成的所有的前端左侧菜单代码，不依赖于原来的内容
// 存储路由规则、静态和动态





func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("GetSession start")

		session := sessions.Default(c)
		menuJson := session.Get("menu")

		if menuJson != nil {
			c.Set("menu",menuJson.(string))
		}else{
			fmt.Println("GetMenu set error ")
		}
		c.Next()


		log.Println("GetSession end")
	}
}


func staticGroupMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("staticGroupMiddle start")

		urlTmp := c.Request.URL.String()
		if utils.InStringSlice(urlTmp, App.Config.whiteUrls) {
			c.Next()
			return
		}
		if utils.InStringSlice(urlTmp, App.Config.blackUrls) {
			// todo 404页面报错
			c.String(http.StatusNotFound, "404")
			c.Abort()
		}
		prefix := App.Config.Prefix
		if ! strings.HasPrefix(urlTmp, prefix) {
			c.Next()
		}
		urlTmp = urlTmp[len(prefix):]


		urlFields := strings.Split(urlTmp, "/")
		fmt.Println(urlFields)
		fmt.Println(len(urlFields))
		if len(urlFields) > 1{
			c.Params = append(c.Params,
				gin.Param{App.Config.RelativeKey, urlFields[0]},
				gin.Param{App.Config.ExtendKey, urlFields[1]},
			)


		}



		log.Println("staticGroupMiddle end")
	}
}

func staticGlobalMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("staticGlobalMiddle start")

		session := sessions.Default(c)
		menuJson := session.Get("menu")

		if menuJson != nil {
			c.Set("menu",menuJson.(string))
		}else{
			fmt.Println("GetMenu set error ")
		}
		c.Next()


		log.Println("staticGlobalMiddle end")
	}
}


