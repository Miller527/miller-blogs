package routers

import (
	"miller-blogs/controllers/manager"
	"miller-blogs/controllers"
	"github.com/astaxie/beego"
	"miller-blogs/middleware"
)






func init() {
    beego.Router("/index", &controllers.IndexController{})
    beego.Router("/manager/login", &manager.LoginController{})
    beego.Router("/manager/index", &manager.IndexManagerController{})
    //beego.Router("/manager/user", &manager.UserManagerController{})
    //beego.Router("/manager/article-list", &manager.ArticleManagerController{})
    //beego.Router("/manager/admin-list", &manager.UserManagerController{})

	// 执行中间件
	//登录时候
	beego.InsertFilter("/manager/*", beego.BeforeStatic, middleware.LoginSessionFilter)

}
