/*
# __author__ = "Mr.chai"
# Date: 2018/12/21
*/
package rbac

import (
	"github.com/gin-gonic/gin"
	"miller-blogs/sugar"
	"miller-blogs/sugar/utils"
)

// session中存储生成的所有的前端左侧菜单代码，不依赖于原来的内容
// 存储路由规则、静态和动态

func rbacLoginMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {


		urlTmp := c.Request.URL.String()
		if utils.InStringSlice(urlTmp, sugar.App.WhiteUrls()) {
			c.Next()
			return
		}
		//if utils.InStringSlice(urlTmp, sugar.App.BlackUrls()) {
		//	// todo 404页面报错
		//	c.String(http.StatusNotFound, "404")
		//	c.Abort()
		//}
		//
		//prefix := App.Config.Prefix
		//if ! strings.HasPrefix(urlTmp, prefix) {
		//	c.Next()
		//}
		//urlTmp = urlTmp[len(prefix):]
		//
		//urlFields := strings.Split(urlTmp, "/")
		//fmt.Println(urlFields)
		//fmt.Println(len(urlFields))
		//if len(urlFields) > 1{
		//	c.Params = append(c.Params,
		//		gin.Param{gr.conf.relativeKey, urlFields[0]},
		//		gin.Param{gr.conf.ExtendKey, urlFields[1]},
		//	)
		//
		//
		//}

	}
}