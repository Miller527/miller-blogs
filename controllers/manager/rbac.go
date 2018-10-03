//
// __author__ = "Miller"
// Date: 2018/10/3
//
package manager

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	. "miller-blogs/controllers/base"
	"miller-blogs/models"
)


// 登录管理后台控制器（管理端）
type LoginController struct {
	Permissions []orm.ParamsList
	SessionPermissions []string
	BlogController
}

// 登录页面
func (loginMC *LoginController) Get() {
	ip := loginMC.Ctx.Input.Context.Request.RemoteAddr
	fmt.Println(ip )
	fmt.Println(logger)
	logger.Info("User login request IP [%s].", ip )
	fmt.Println(	 loginMC.Ctx.Input.Cookie(beego.AppConfig.String("sessionname")))
	if loginMC.GetSession(beego.AppConfig.String("permission_key")) != nil {
		loginMC.Ctx.Redirect(302, beego.AppConfig.String("manager_router_prefix") + "/index")
		return
	}
	loginMC.TplName = loginMC.GetManagerPagePath("login.html")
}

// 登录请求, 将用户权限写入Session
func (loginMC *LoginController) Post() {
	// 登录认证
	userId := loginMC.GetString("username")
	userPwd := loginMC.GetString("userpwd")
	permissions := loginMC.GetSession(beego.AppConfig.String("permission_key"))

	var user models.UserInfo

	err := loginMC.OrmObj.QueryTable("user_info").Filter("uid", userId).
		Filter("password", userPwd).One(&user)

	if err == orm.ErrMultiRows {
		loginMC.UpdateResponseMsg(40100,"后台异常, 请稍后重试", nil)
	} else if err == orm.ErrNoRows {
		loginMC.UpdateResponseMsg(30100,"用户或密码错误", nil)
	} else {
		loginMC.UpdateResponseMsg(20100,"登录成功", nil)
		if permissions == nil {
			loginMC.writePermissionsSession(userId)
		}
	}
	fmt.Println(loginMC.ResponseData)
	loginMC.Data["json"] = &loginMC.ResponseData
	loginMC.ServeJSON()
	return
}

// 注销请求, 将用户存放session的数据删除
func (loginMC *LoginController) Del() {

}

// 获取用户的权限信息
func (loginMC *LoginController) queryPermissions(userName string) bool {
	_, err := loginMC.OrmObj.QueryTable("permission").
		Filter("Roles__Role__Users__UserInfo__Uid", userName).
		ValuesList(&loginMC.Permissions, "title", "url", "is_menu")
	fmt.Println(loginMC.Permissions)
	if err == nil{
		return true
	}
	return false
}

// 将用户的权限信息写入session, 放到缓存中
func (loginMC *LoginController) writePermissionsSession(uid string) {
	if status := loginMC.queryPermissions(uid) ; status{
		fmt.Println(loginMC.Permissions)
		if  byteStr ,err := json.Marshal(loginMC.Permissions); err == nil{
			go loginMC.SetSession(beego.AppConfig.String("permission_key"),string(byteStr))
			go loginMC.SetSession("uid",uid)
		}else {
			fmt.Println("json error")
		}
	}else{
		fmt.Println("wori")
	}
}
