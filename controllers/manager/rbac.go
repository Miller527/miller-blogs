//
// __author__ = "Miller"
// Date: 2018/10/3
//
package manager

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	. "miller-blogs/controllers/base"
	"miller-blogs/models"
)


// 登录管理后台控制器（管理端）
type LoginController struct {
	RbacController
}

// 登录页面
func (loginMC *LoginController) Get() {
	ip := loginMC.Ctx.Input.Context.Request.RemoteAddr
	fmt.Println(ip )
	fmt.Println(logger)
	logger.Info("User login request IP [%s].", ip )
	if loginMC.GetSession(beego.AppConfig.String("session_permission_key")) != nil {
		indexUrl := loginMC.PathGuidance[0]["url"]
		loginMC.Ctx.Redirect(302, indexUrl)
		return
	}
	loginMC.TplName = loginMC.GetManagerPagePath("login.html")
}

// 登录请求, 将用户权限写入Session
func (loginMC *LoginController) Post() {
	// 登录认证
	userId := loginMC.GetString("username")
	userPwd := loginMC.GetString("userpwd")
	permissions := loginMC.GetSession(beego.AppConfig.String("session_permission_key"))

	var user models.UserInfo
	// TODO 做用户管理时候在这加上密码加密
	err := loginMC.OrmObj.QueryTable("user_info").Filter("uid", userId).
		Filter("password", userPwd).One(&user)

	if err == orm.ErrMultiRows {
		loginMC.UpdateResponseMsg(40100,"后台异常, 请稍后重试", nil)
	} else if err == orm.ErrNoRows {
		loginMC.UpdateResponseMsg(30100,"用户或密码错误", nil)
	} else {
		indexUrl := loginMC.PathGuidance[0]["url"]
		loginMC.UpdateResponseMsg(20100,"登录成功", map[string]interface{}{"url":indexUrl})
		if permissions == nil {
			// todo 维护登录时间, 更新session的过期时间
			loginMC.WriteSession(userId)
		}
	}
	loginMC.Data["json"] = &loginMC.ResponseData
	loginMC.ServeJSON()
	return
}

// 注销请求, 将用户存放session的数据删除
func (loginMC *LoginController) Del() {

}


