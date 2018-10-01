/*
# __author__ = "Mr.chai"
# Date: 2018/9/13
*/
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
	RbacController
	Permissions []orm.ParamsList
}

// 获取orm对象和日志
func (loginMC *LoginController) Prepare() {
	loginMC.GetOrmer()
	loginMC.ResponseData = make(map[string]interface{})
}
// 登录页面
func (loginMC *LoginController) Get() {
	fmt.Println(loginMC.Data["RouterPattern"])
	fmt.Println(loginMC.CruSession)
	fmt.Println("1111===",loginMC.GetSession("permissions"))
	fmt.Println("2222===",loginMC.GetSession(loginMC.Ctx.GetCookie(beego.AppConfig.String("sessionname"))))
	fmt.Println("3333====",loginMC.CruSession)

	fmt.Println(loginMC.Ctx.GetCookie(beego.AppConfig.String("sessionname")))
	loginMC.TplName = "manager/login.html"
}

// 登录请求, 将用户权限写入Session
func (loginMC *LoginController) Post() {



	// 登录认证
	userName := loginMC.GetString("username")
	userPwd := loginMC.GetString("userpwd")
	permissions := loginMC.GetSession("permissions")

	var user models.UserInfo

	err := loginMC.OrmObj.QueryTable("user_info").Filter("uid", userName).
		Filter("password", userPwd).One(&user)

	if err == orm.ErrMultiRows {
		loginMC.UpdateResponseMsg(40100,"后台异常, 请稍后重试", nil)
	} else if err == orm.ErrNoRows {
		loginMC.UpdateResponseMsg(30100,"用户或密码错误", nil)
	} else {
		loginMC.UpdateResponseMsg(20100,"登录成功", nil)
		if permissions == nil {
			loginMC.PermissionsSession(userName)
		}
	}
	fmt.Println(loginMC.ResponseData)
	loginMC.Data["json"] = &loginMC.ResponseData
	loginMC.ServeJSON()
	return

}

// 获取用户的权限信息
func (loginMC *LoginController) QueryPermissions(userName string) bool {
	_, err := loginMC.OrmObj.QueryTable("permission").
		Filter("Roles__Role__Users__UserInfo__Uid", userName).
		ValuesList(&loginMC.Permissions, "name", "url", "type")
	fmt.Println(err, loginMC.Permissions)

	if err == nil{
		return true
	}
	return false
}

// 将用户的权限信息写入session, 放到缓存中
func (loginMC *LoginController) PermissionsSession(userName string) {
	if status := loginMC.QueryPermissions(userName) ; status{
		fmt.Println(loginMC.Permissions)
		if  byteStr ,err := json.Marshal(loginMC.Permissions); err == nil{
			go loginMC.SetSession("permissions",string(byteStr))
			go loginMC.SetSession("uid",userName)
		}else {
			fmt.Println("json error")
		}
	}else{
		fmt.Println("wori")
	}
}






















//
//
//
//
//// 站点管理控制器
//type SiteManagerController struct {
//	CurdBaseController
//}
//
//func (siteMC *SiteManagerController) Get() {
//
//}
//
//func (siteMC *SiteManagerController) Post() {
//
//}
//
//func (siteMC *SiteManagerController) Put() {
//
//}
//
//func (siteMC *SiteManagerController) Del() {
//
//}
//
//
//
type HeaderData struct {
	Name string
}
//
// 管理后台首页，获取数据渲染首页
type IndexManagerController struct {
	CurdController
}

// 带着cookie才可以（RBAC权限）
func (indexMC *IndexManagerController) Get() {

	fmt.Println("1111---",indexMC.GetSession("permissions"))
	fmt.Println("2222---",indexMC.GetSession(indexMC.Ctx.GetCookie(beego.AppConfig.String("sessionname"))))
	fmt.Println("3333---",indexMC.CruSession)
	indexMC.Layout = "manager/base.html"

	indexMC.TplName = "manager/Hui-admin/index.html"

	indexMC.LayoutSections = make(map[string]string)
	indexMC.LayoutSections["HeadMeta"] = "manager/Hui-admin/headmeta.html"
	indexMC.LayoutSections["Header"] = "manager/Hui-admin/header.html"
	indexMC.LayoutSections["LeftMenu"] = "manager/Hui-admin/leftmenu.html"
	indexMC.Data["headerData"] = &HeaderData{"Miller"}
}
//
//// TODO 用户管理
//type UserManagerController struct {
//	BlogBaseController
//	UserList []orm.ParamsList
//}
//
//// 基础, 用户登录认证, 初始化表字段等
//func (userMC *UserManagerController) Prepare() {
//	userMC.Login()
//	userMC.TitleInit()
//}
//
//// 初始化表字段, 所有CURD相关的表都得重写该字段
//func (userMC *UserManagerController) TitleInit() {
//	userMC.DisplayTitle = []string{"ID", "帐号", "昵称", "邮箱", "手机号", "头像", "类型", "角色id",
//		"创建时间", "更新时间"}
//	userMC.FieldTitle = []string{"id", "uid", "nick_name", "email", "phone", "mugshot", "type",
//		"role__name", "created_time", "updated_time"}
//	fmt.Println(userMC.DisplayTitle)
//	fmt.Println(userMC.FieldTitle)
//}
//
//// 用户管理Get接口, 获取用户列表
//func (userMC *UserManagerController) Get() {
//	userMC.GetUserList()
//	userMC.Layout = "manager/base.html"
//	userMC.TplName = "manager/Hui-admin/admin_list.html"
//	userMC.LayoutSections = make(map[string]string)
//	userMC.LayoutSections["HeadMeta"] = "manager/Hui-admin/headmeta.html"
//	userMC.LayoutSections["Header"] = "manager/Hui-admin/header.html"
//	userMC.LayoutSections["LeftMenu"] = "manager/Hui-admin/leftmenu.html"
//	userMC.Data["headerData"] = &HeaderData{"Miller"}
//	userMC.Data["tableHeader"] = &userMC.DisplayTitle
//	userMC.Data["tableField"] = &userMC.FieldTitle
//
//	userMC.Data["tableData"] = &userMC.UserList
//}
//
//// 获取用户列表
//func (userMC *UserManagerController) GetUserList() {
//	ormObj := orm.NewOrm()
//	qs, err := ormObj.QueryTable("blog_user").ValuesList(&userMC.UserList, userMC.FieldTitle...)
//	fmt.Println(qs, err)
//	fmt.Println()
//}
//
//// 用户管理Post接口, 添加用户
//func (userMC *UserManagerController) Post() {
//
//}
//
//// 文章管理
//type ArticleManagerController struct {
//	BlogBaseController
//	ArticleList []orm.ParamsList
//}
//
////基础, 用户登录认证, 初始化表字段等
//func (articleMC *ArticleManagerController) Prepare() {
//	articleMC.Login()
//	articleMC.TitleInit()
//}
//
////初始化表字段, 所有CURD相关的表都得重写该字段
//func (articleMC *ArticleManagerController) TitleInit() {
//	articleMC.DisplayTitle = []string{"ID", "标题", "简介", "内容", "分类", "类型", "发布状态", "评论状态",
//		"权重", "访问量", "公共数据", "创建时间", "更新时间"}
//	articleMC.FieldTitle = []string{"id", "title", "intro", "body", "status", "article_type_id__name", "category__name",
//		"comment_status", "order_weight", "pv", "pub_date", "created_time", "updated_time"}
//	fmt.Println(articleMC.DisplayTitle)
//	fmt.Println(articleMC.FieldTitle)
//}
//func (articleMC *ArticleManagerController) GetUserList() {
//	ormObj := orm.NewOrm()
//	qs, err := ormObj.QueryTable("article").ValuesList(&articleMC.ArticleList, articleMC.FieldTitle...)
//	fmt.Println(articleMC.ArticleList)
//	fmt.Println(qs, err)
//	fmt.Println()
//}
//
//func (articleMC *ArticleManagerController) Get() {
//	articleMC.GetUserList()
//	articleMC.Layout = "manager/base.html"
//	articleMC.TplName = "manager/Hui-admin/article_list.html"
//
//	articleMC.LayoutSections = make(map[string]string)
//	articleMC.LayoutSections["HeadMeta"] = "manager/Hui-admin/headmeta.html"
//	articleMC.LayoutSections["Header"] = "manager/Hui-admin/header.html"
//	articleMC.LayoutSections["LeftMenu"] = "manager/Hui-admin/leftmenu.html"
//	articleMC.Data["headerData"] = &HeaderData{"Miller"}
//	articleMC.Data["tableHeader"] = &articleMC.DisplayTitle
//	articleMC.Data["tableField"] = &articleMC.FieldTitle
//	articleMC.Data["tableData"] = &articleMC.ArticleList
//}
//
//func (articleMC *ArticleManagerController) Post() {
//
//}
