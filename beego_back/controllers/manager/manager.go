/*
# __author__ = "Mr.chai"
# Date: 2018/9/13
*/
package manager

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	. "miller-blogs/beego_back/controllerscontrollers/base"
	"miller-blogs/beego_back/publicback/public"
)

var logger = public.GetLogger("blog_manager")

// 管理后台首页，获取数据渲染首页
type IndexManagerController struct {
	CurdController
}

// 管理后台首页
func (indexMC *IndexManagerController) Get() {
	indexMC.ResponseTemplate("index.html")

}



// 管理后台首页，获取数据渲染首页
type PermissionManagerController struct {
	CurdController
}

// 按照需求重写Prepare和ExtendFieldInit
func (perMC *PermissionManagerController) Prepare() {
	perMC.ExtendFieldInit()
	perMC.DefaultInit()

}

func (perMC *PermissionManagerController) ExtendFieldInit() {

	perMC.DisplayTitle = []string{"ID", "帐号", "昵称", "邮箱", "手机号", "头像", "类型", "角色id",
		"创建时间", "更新时间"}
	perMC.FieldTitle = []string{"id", "uid", "nick_name", "email", "phone", "mugshot", "type",
		"role__name", "created_time", "updated_time"}
}

// 带着cookie才可以（RBAC权限）
func (perMC *PermissionManagerController) Get() {

	perMC.ResponseTemplate("permission_list.html")

}


func (perMC *PermissionManagerController) Post() {

	//fmt.Println("1111---",indexMC.GetSession("permissions"))

	perMC.Layout = "manager/base.html"

	perMC.TplName = "manager/Hui-admin/index.html"

	perMC.LayoutSections = make(map[string]string)
	perMC.LayoutSections["HeadMeta"] = "manager/Hui-admin/headmeta.html"
	perMC.LayoutSections["Header"] = "manager/Hui-admin/header.html"
	perMC.LayoutSections["LeftMenu"] = "manager/Hui-admin/menu.html"
	perMC.Data["headerData"] = &HeaderData{"Miller"}
}

func (perMC *PermissionManagerController) Put() {

	//fmt.Println("1111---",indexMC.GetSession("permissions"))

	perMC.Layout = "manager/base.html"

	perMC.TplName = "manager/Hui-admin/index.html"

	perMC.LayoutSections = make(map[string]string)
	perMC.LayoutSections["HeadMeta"] = "manager/Hui-admin/headmeta.html"
	perMC.LayoutSections["Header"] = "manager/Hui-admin/header.html"
	perMC.LayoutSections["LeftMenu"] = "manager/Hui-admin/menu.html"
	perMC.Data["headerData"] = &HeaderData{"Miller"}
}

func (perMC *PermissionManagerController) Del() {

	//fmt.Println("1111---",indexMC.GetSession("permissions"))
	perMC.Ctx.ResponseWriter.Write()
	perMC.Layout = "manager/base.html"

	perMC.TplName = "manager/Hui-admin/index.html"

	perMC.LayoutSections = make(map[string]string)
	perMC.LayoutSections["HeadMeta"] = "manager/Hui-admin/headmeta.html"
	perMC.LayoutSections["Header"] = "manager/Hui-admin/header.html"
	perMC.LayoutSections["LeftMenu"] = "manager/Hui-admin/menu.html"
	perMC.Data["headerData"] = &HeaderData{"Miller"}
}













//








//
//// TODO 用户管理
type UserManagerController struct {
	CurdController
	UserList []orm.ParamsList
}

// 基础, 用户登录认证, 初始化表字段等
func (userMC *UserManagerController) Prepare() {
	userMC.TitleInit()
}

// 初始化表字段, 所有CURD相关的表都得重写该字段
func (userMC *UserManagerController) TitleInit() {
	userMC.DisplayTitle = []string{"ID", "帐号", "昵称", "邮箱", "手机号", "头像", "类型", "角色id",
		"创建时间", "更新时间"}
	userMC.FieldTitle = []string{"id", "uid", "nick_name", "email", "phone", "mugshot", "type",
		"role__name", "created_time", "updated_time"}
	fmt.Println(userMC.DisplayTitle)
	fmt.Println(userMC.FieldTitle)
}

// 用户管理Get接口, 获取用户列表
//func (userMC *UserManagerController) Get() {
//
//	userMC.GetUserList()
//	userMC.Layout = "manager/base.html"
//	userMC.TplName = "manager/Hui-admin/admin_list.html"
//	userMC.LayoutSections = make(map[string]string)
//	userMC.LayoutSections["HeadMeta"] = "manager/Hui-admin/headmeta.html"
//	userMC.LayoutSections["Header"] = "manager/Hui-admin/header.html"
//	userMC.LayoutSections["LeftMenu"] = "manager/Hui-admin/menu.html"
//	userMC.Data["headerData"] = &HeaderData{"Miller"}
//	userMC.Data["tableHeader"] = &userMC.DisplayTitle
//	userMC.Data["tableField"] = &userMC.FieldTitle
//
//	userMC.Data["tableData"] = &userMC.UserList
//}

// 获取用户列表
func (userMC *UserManagerController) getUserList() {
	ormObj := orm.NewOrm()
	qs, err := ormObj.QueryTable("blog_user").ValuesList(&userMC.UserList, userMC.FieldTitle...)
	fmt.Println(qs, err)
	fmt.Println()
}

// 用户管理Post接口, 添加用户
func (userMC *UserManagerController) Post() {

}
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
//	articleMC.LayoutSections["LeftMenu"] = "manager/Hui-admin/menu.html"
//	articleMC.Data["headerData"] = &HeaderData{"Miller"}
//	articleMC.Data["tableHeader"] = &articleMC.DisplayTitle
//	articleMC.Data["tableField"] = &articleMC.FieldTitle
//	articleMC.Data["tableData"] = &articleMC.ArticleList
//}
//
//func (articleMC *ArticleManagerController) Post() {
//
//}
