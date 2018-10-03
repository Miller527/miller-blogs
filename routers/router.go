package routers

import (
	"github.com/astaxie/beego"
	"miller-blogs/controllers/manager"
	"miller-blogs/middleware"
)

var urlPrefix = beego.AppConfig.String("manager_router_prefix")

// curd路由生成
func getLinkNamespace (apiUrl string, control beego.ControllerInterface) beego.LinkNamespace{
	return beego.NSNamespace(apiUrl,
			beego.NSRouter("/list", control, "get:Get"),
			beego.NSRouter("/add", control,"post:Post"),
			beego.NSRouter("/update", control,"put:Put"),
			beego.NSRouter("/del", control,"delete:Delete"),
		)
}


func init() {

	managerNameSpace := beego.NewNamespace(urlPrefix,
			// 登录和登出
			beego.NSRouter("/login", &manager.LoginController{}),
			// 管理后台的首页
			beego.NSRouter("/index", &manager.ManagerIndexController{}),
			// 权限管理
			getLinkNamespace("/permission",&manager.PermissionManagerController{}),

			// 中间件, 通过session管理登录和首页的跳转
			beego.NSBefore(middleware.ManagerLoginFilter),
		)

	// 注册后台管理的路由空间
	beego.AddNamespace(managerNameSpace)

	//// 中间件, 登录认证, 通过session里的权限管理访问
	beego.InsertFilter("/*", beego.BeforeRouter, middleware.PermissionVerify)


	//注册 namespace
    //beego.Router("/manager/user", &manager.UserManagerController{})
    //beego.Router("/manager/article-list", &manager.ArticleManagerController{})
    //beego.Router("/manager/admin-list", &manager.UserManagerController{})
	//beego.AutoRouter(&manager.UserManagerController{})
	//beego.AutoPrefix(urlPrefix,&manager.UserManagerController{})

	//fmt.Println(beego.URLFor("UserManagerController"))
	//fmt.Println(beego.URLFor("UserManagerController.Get"))
	//fmt.Println(beego.URLFor("UserManagerController.Initialization"))
	//fmt.Println(beego.URLFor("UserManagerController.getUserList"))
	//fmt.Println(beego.URLFor("UserManagerController.Prepare"))

	// Blog router
	//beego.Router("/index", &controllers.IndexController{})


}
