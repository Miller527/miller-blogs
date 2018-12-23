/*
# __author__ = "Mr.chai"
# Date: 2018/12/21
*/
package rbac

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// session中存储生成的所有的前端左侧菜单代码，不依赖于原来的内容
// 存储路由规则、静态和动态

func RbacLoginMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(c)

		// Process request
		c.Next()

		fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	}
}



func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}