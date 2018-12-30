/*
# __author__ = "Mr.chai"
# Date: 2018/12/21
*/
package sugar

import (
	"encoding/json"
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
		menuStr := session.Get("menu")

		if menuStr != nil {
			me := SortedMenu{}
			err := json.Unmarshal([]byte(menuStr.(string)), &me)
			if err == nil{
				c.Set("menu",me)

			}else {
				fmt.Println("GetMenu set error ")
			}
		}
		c.Next()


		log.Println("GetSession end")
	}
}