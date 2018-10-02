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

	// 中间件
	// 登录认证, 通过session管理
	beego.InsertFilter("/manager/*", beego.BeforeRouter, middleware.ManagerLoginFilter)
	beego.InsertFilter("/*", beego.BeforeRouter, middleware.PermissionVerify)

}
