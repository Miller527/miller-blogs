/*
# __author__ = "Mr.chai"
# Date: 2018/9/9
*/
package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type TimeModel struct {
	UpdatedTime time.Time `orm:"auto_now;type(datetime)" explain:"更新时间"`
	CreatedTime time.Time `orm:"auto_now_add;type(datetime)" explain:"创建时间"`
}

// TODO Rbac基于角色的权限管理
// 登陆管理后台使用
type UserInfo struct {
	Id       int     `explain:"id"`
	Uid  string  `orm:"size(32)" explain:"帐号"`
	Password string  `orm:"size(128)" explain:"密码"`
	Nickname string  `orm:"size(64)" explain:"昵称"`
	Email    string  `orm:"size(64)" explain:"邮箱"`
	Phone    string  `orm:"size(11)" explain:"手机"`
	Roles    []*Role `orm:"rel(m2m)" explain:"角色"`
	//Site     *BlogSite `orm:"rel(fk)" explain:"站点"`
	Mugshot string `orm:"size(11)" explain:"头像"` // 头像地址, 通过配置的路径拼接
	Type    int8   `orm:"size(11)" explain:"类型"` // 账号类型, 0管理员, 1普通用户, 往后排是第三方用户
	// 普通用户和第三方用户不能后台管理
	TimeModel
}

//// 用户登录历史
//type LoginHistory struct {
//
//}

// 用户角色, 访客角色, 只能评论
type Role struct {
	Id          int           `explain:"id"`
	Rid         string        `orm:"unique" explain:"角色"`
	Name        string        `orm:"size(64)" explain:"名称"`
	Users       []*UserInfo   `orm:"reverse(many)" explain:"关联用户"`
	Permissions []*Permission `orm:"rel(m2m)" explain:"权限列表"`
}

//manager开头, 用户权限, 即菜单, 一级菜单主要是些详细报表, 二级菜单主要是基本功能菜单, 权限不是菜单
type Permission struct {
	Id     int         `orm:"on_delete(cascade)" explain:"id"`
	Title  string      `orm:"size(128)"  explain:"权限"`
	Url    string      `orm:"size(256)" explain:"链接地址"`
	Icon   string      `orm:"null;size(128)" explain:"菜单图标"`
	IsMenu bool        `explain:"菜单状态"`                                     // 该权限是否是菜单权限
	Parent *Permission `orm:"null;rel(fk)" explain:"父菜单"` // 必须是子菜单才能有父菜单
	Roles  []*Role     `orm:"reverse(many);on_delete(cascade)" explain:"所属角色"`
}

// TODO 公共参数
//// 站点模版
//type BlogSite struct {
//	Id                 int
//	SiteName           string `orm:"size(60)"`
//	SitePath           string `orm:"size(60)"`
//	SiteDescription    string //网站描述
//	SiteSeoDescription string //网站SEO描述
//	SiteKeyword        string //网站关键字
//	Pv                 int    //站点浏览量(一天一个IP记录一个有效浏览量)
//	AdStatus           bool   //是否开启广告
//	CommentStatus      bool   //是否开启评论
//	CommentAnonymous   bool   //是否开启匿名评论（开启匿名评论首先要开启评论，匿名用户通过session自动创建用户）
//	Template           string //模板类型
//	RecordCode         string //网站备案号
//	PoliceRecordStatus string //是否显示公安备案号
//	PoliceRecordCode   string //公安备案号
//	Users              []*UserInfo `orm:"reverse(many)"`
//	TagType            *TagType    `orm:"rel(fk)"`	//标签云类型, 根据模板名字匹配
//	TimeModel
//}
//
//

//// 管理后台登录历史
//type ManagerLoginHistory struct {
//	Id   int
//	Ip   string
//	User *BlogUser `orm:"rel(fk)"`
//	TimeModel
//}
//
//
//
//
//
//// 访客, 靠角色维护
////type BlogVisitor struct {
////	Id       int
////	Uid      string `orm:"size(32)"`
////	PassWord string `orm:"size(128)"`
////	NickName string `orm:"size(64)"`
////	Email    string `orm:"size(64)"`
////	Mugshot  string
////	TimeModel
////}
//
//// 广告
//type Ad struct {
//	Id     int
//	Name   string `orm:"size(64)"`
//	Url    string
//	Type   string //image、text、video
//	Status bool   // 该广告是否显示
//	TimeModel
//}
//
//// 友链
//type BlogRoll struct {
//	Id     int
//	Name   string `orm:"size(64)"`
//	Url    string
//	Level  string // 友链等级
//	Status bool   // 该友链是否显示
//	TimeModel
//}
//
//// 文章类型
//type ArticleType struct {
//	Id     int
//	Name   string
//	Status bool //是否禁用
//	Article []*Article `orm:"reverse(many)"`
//	TimeModel
//}
//
//// 文章分类
//type Category struct {
//	Id     int
//	Name   string
//	Status bool //是否禁用
//	Article []*Article `orm:"reverse(many)"`
//	TimeModel
//}
//
//// 标签云类型, 根据模板匹配
//type TagType struct {
//	Id   int
//	Name string
//	Site []*BlogSite `orm:"reverse(many)"`
//	TimeModel
//}
//
//// 标签, 标签云直接groupby计算
//type Tag struct {
//	Id       int
//	Name     string
//	Articles []*Article `orm:"reverse(many)"`
//	TimeModel
//}
//
//// 文章
//type Article struct {
//	Id            int
//	Title         string
//	Intro          string       `orm:"type(text)"`
//	Body          string       `orm:"type(text)"`
//	ArticleType   *ArticleType `orm:"rel(fk)"` // 文章分类(article、album、video等)
//	Category      *Category    `orm:"rel(fk)"` // 文章类型(Python、Golang、区块链等)
//	Status        bool                         // 文章发布状态(草稿、发布)
//	CommentStatus bool                         // 评论状态，是否允许评论
//	OrderWeight   int                          // 文章排序权重(针对于首页)
//	Pv            int                          // 文章浏览量
//	Tags          []*Tag `orm:"rel(m2m)"`
//	PubDate       string
//	TimeModel
//}
//
//// 评论
//type Comment struct {
//	Id         int
//	Body       string   `orm:"type(text)"`
//	Article    *Article `orm:"rel(fk)"`
//	SubComment *Comment `orm:"rel(one)"` //子评论
//	TimeModel
//}

//type Profile struct {
//	Id      int
//	Age     int16
//	//User    *User    `orm:"reverse(one)"`  // 设置一对一反向关系(可选)
//	Profile *Profile `orm:"rel(one)"`      // OneToOne relation
//	Post    []*Post  `orm:"reverse(many)"` // 设置一对多的反向关系
//}
//
//type Post struct {
//	Id    int
//	Title string
//	//User  *User  `orm:"rel(fk)"` //设置一对多关系
//	Tags  []*Tag `orm:"rel(m2m)"`
//}

type twoLevelMenu struct {
	title    string
	url      string
	children [] map[string]interface{}
}
type leftMenuList map[string]interface{}

// 初始化数据表
func Initialization() {

	z := []*Permission{}
	obj := orm.NewOrm()
	a, b := obj.QueryTable("permission").
		Filter("Roles__Role__Users__UserInfo__Uid", "Miller").

		//ValuesList(&z, "id", "title", "url", "is_menu", "parent")
		All(&z, "id", "title", "url", "is_menu", "parent")
	fmt.Println(a, b)

	permissionsList := make([]string, 10)
	menuDict := make(map[int]interface{}, 10)

	for _, obj := range z {
		fmt.Println(obj)
		permissionsList = append(permissionsList, obj.Url)
		permission := leftMenuList{
			"title": obj.Title,
			"url":   obj.Url,
		}
		if obj.Id != 0 && obj.IsMenu {
			if _, ok := menuDict[obj.Id]; ! ok && obj.Parent == nil{
				menuDict[obj.Id] = leftMenuList{
					"title":    obj.Title,
					"url":      obj.Url,
					"children": [] leftMenuList{},
				}
			} else {
				menuDict[obj.Parent.Id].(leftMenuList)["children"] = append(menuDict[obj.Parent.Id].(leftMenuList)["children"].([] leftMenuList), permission)
			}
		}
	}
	fmt.Println(menuDict)
	fmt.Println(menuDict[1])
	fmt.Println(menuDict[1].(leftMenuList)["children"])
	c,err := json.Marshal(menuDict)
	fmt.Println(err,string(c))
}
func init() {
	dataSource := beego.AppConfig.String("data_source")
	driverName := beego.AppConfig.String("data_driver")
	fmt.Println(dataSource, driverName)

	orm.RegisterDriver(driverName, orm.DRMySQL)
	orm.RegisterDataBase("default", driverName, dataSource, 30)

	orm.RegisterModel(new(UserInfo), new(Role), new(Permission))
	orm.RunSyncdb("default", false, true)

	Initialization()
}
