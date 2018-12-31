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