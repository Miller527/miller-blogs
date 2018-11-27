//
// __author__ = "Miller"
// Date: 2018/11/15
//

package urls

import (
	"github.com/gin-gonic/gin"
)


var AdApp *gin.Engine

func url(app *gin.Engine) {

	//behavior := app.Group("/ad-behavior/"  )

	//behavior.Use(middleware.AppKeyFilter())
	//{
	//	behavior.GET("/:appKey", apps.AdBehavior)
	//	behavior.POST("/:appKey", apps.AdBehavior)
	//}

}


func init()  {
	AdApp = gin.Default()
	url(AdApp)
}