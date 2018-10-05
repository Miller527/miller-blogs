package base

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"miller-blogs/models"
	"path"
	"reflect"
	"sort"
)

// 统一的返回值格式
type ResponseMsg struct {
	Status int                    `json:"status"`
	Msg    string                 `json:"msg"`
	Data   map[string]interface{} `json:"data"`
}

// Public base struct
type baseController struct {
	beego.Controller
	// Todo 后续改成数据库或者加载到缓存
	ManagerSite  string //管理后台站点文件路径, 暂定配置文件中配置
	BlogSite     string //博客站点文件路径, 暂定配置文件中配置
	OrmObj       orm.Ormer
	ResponseData ResponseMsg
	BlogLogger   *logs.BeeLogger
}

// 按照需求重写该字段
func (base *baseController) Prepare() {
	base.SiteManager()
	base.Initialization()
	base.ExtendFieldInit()
}

// 获取orm对象
func (base *baseController) Initialization() {
	base.OrmObj = orm.NewOrm()
	base.ResponseData = ResponseMsg{}
}

// 扩展字段的初始化
func (base *baseController) ExtendFieldInit() {}

// 站点管理
func (base *baseController) SiteManager() {
	base.ManagerSite = path.Join("manager", beego.AppConfig.String("manager_file_path"))
	base.BlogSite = "xxx"
}

// 获取返回页面文件路径
func (base *baseController) GetManagerPagePath(filename string) string {
	return path.Join(base.ManagerSite, filename)
}

// 更新ResponseData
func (base *baseController) UpdateResponseMsg(status int, msg string, data map[string]interface{}) {
	fmt.Println("base.ResponseData============", base.ResponseData)
	//base.ResponseData["status"] = status
	//base.ResponseData["msg"] = msg
	//base.ResponseData["data"] = data
	base.ResponseData.Status = status
	base.ResponseData.Msg = msg
	base.ResponseData.Data = data
	fmt.Println("base.ResponseData============", base.ResponseData)

}


type menuData map[string]interface{}
type menuDict map[int]menuData

// Rbac struct
type rbacBaseController struct {
	MenuDict        map[string]interface{}
	PermissionsData []*models.Permission
	permissionTitle []string
	baseController
}

// 配置信息, 修改数据库字段的话改这里
func (rbac *rbacBaseController) configInit() {
	rbac.permissionTitle = []string{"id", "title", "url", "is_menu", "parent"}
}

// 初始化permissionTitle字段
func (rbac *rbacBaseController) ExtendFieldInit() {
	rbac.configInit()
	//rbac.getUrlIndex ()
}

//// 获取orm请求后的permission字段索引
//func (rbac *rbacBaseController) getUrlIndex (){
//	for i, v := range rbac.permissionTitle{
//		if v == "url" {
//			rbac.urlIndex = i
//		}
//	}
//}

// 获取用户的权限信息
func (rbac *rbacBaseController) queryPermissions(userName string) bool {
	fmt.Println(userName)
	_, err := rbac.OrmObj.QueryTable("permission").
		Filter("Roles__Role__Users__UserInfo__Uid", userName).
		All(&rbac.PermissionsData, rbac.permissionTitle...)
	fmt.Println("permissions", err, rbac.PermissionsData)
	if err == nil {
		return true
	}
	return false
}

// 将用户的权限信息写入session, 放到缓存中
func (rbac *rbacBaseController) WriteSession(uid string) {
	if status := rbac.queryPermissions(uid); status {
		rbac.sessionInit(uid)
	} else {
		fmt.Println("wori")
	}
}


// 获取用户的左侧菜单
func (rbac *rbacBaseController) GetLeftMenu() []menuData {
	leftMenuStr, _ := rbac.GetSession(beego.AppConfig.String("session_menu_key")).(string)
	sideMenu := menuDict{}
	json.Unmarshal([]byte(leftMenuStr), &sideMenu)
	return rbac.orderMenu(sideMenu)

}

// 菜单排序
func (rbac *rbacBaseController) orderMenu(menuDict menuDict) (sideMenu []menuData) {
	leftMenuStr, _ := rbac.GetSession(beego.AppConfig.String("session_menu_key")).(string)
	json.Unmarshal([]byte(leftMenuStr), &menuDict)
	var sortKeys  []int
	for k := range menuDict {
		sortKeys = append(sortKeys, k)
	}
	sort.Ints(sortKeys)

	for _, k := range sortKeys {
		menu := menuDict[k]
		if len(menu["children"].([] interface{})) == 0{
			delete(menuDict, k)
			continue
		}
		rbac.updateMenuStatus(menuDict[k])
		sideMenu = append(sideMenu, menuDict[k])
	}
	return
}

// 菜单状态修改
func (rbac *rbacBaseController) updateMenuStatus(menu menuData) () {

	fmt.Println(menu["children"], reflect.TypeOf(menu["children"]))
	for _, val := range menu["children"].([] interface{}){
		fmt.Println(rbac.Ctx.Request.RequestURI)
		fmt.Println(val.(map[string]interface{})["url"])
		if val.(map[string]interface{})["url"] == rbac.Ctx.Request.RequestURI {
			val.(map[string]interface{})["class"] = "current"
			menu["class"] = "selected"
			menu["style"] = "display: block"
			break
		}
			//}else {
		//	menu["style"] = "display: none"
		//}

	}
	fmt.Println(menu)

}


// 初始化session字典
func (rbac *rbacBaseController) sessionInit(uid string) {
	go rbac.SetSession(beego.AppConfig.String("session_uid_key"), uid)
	var permissionsList  []string
	menuDict := menuDict{}

	for _, obj := range rbac.PermissionsData {
		permissionsList = append(permissionsList, obj.Url)
		permission := menuData{
			"title": obj.Title,
			"url":   obj.Url,
		}
		if obj.Id != 0 && obj.IsMenu {
			if _, ok := menuDict[obj.Id]; ! ok && obj.Parent == nil {
				menuDict[obj.Id] = menuData{
					"title":    obj.Title,
					"url":      obj.Url,
					"children": [] menuData{},
				}
			} else {
				tmpList := append(menuDict[obj.Parent.Id]["children"].([] menuData), permission)
				menuDict[obj.Parent.Id]["children"] = tmpList
			}
		}
	}
	if byteStr, err := json.Marshal(menuDict); err == nil {
		go rbac.SetSession(beego.AppConfig.String("session_menu_key"), string(byteStr))
	}
	if byteStr, err := json.Marshal(permissionsList); err == nil {
		go rbac.SetSession(beego.AppConfig.String("session_permission_key"), string(byteStr))
	}
}

type HeaderData struct {
	Name string
}
// CURD base struct
type curdBaseController struct {
	rbacBaseController
	DisplayTitle []string //前端显示的表字段名字
	FieldTitle   []string //数据库的字段名，和上面的显示字段从前一一对应, 长度可以不一样，后边缺少的字段是前端自定义的字段
	HeaderData *HeaderData
}

// 获取model的自定义的默认字段
func (curd *curdBaseController) DefaultFiledTitles(model struct{}) {

}

// 生成模范的返回数据
func (curd *curdBaseController) ResponseTemplate(htmlName string) {
	curd.Layout =  curd.GetManagerPagePath("base.html")
	curd.TplName = curd.GetManagerPagePath(htmlName)
	curd.LayoutSections = make(map[string]string)
	curd.LayoutSections["HeadMeta"] = curd.GetManagerPagePath("headmeta.html")
	curd.LayoutSections["Header"] = curd.GetManagerPagePath("header.html")
	curd.LayoutSections["LeftMenu"] = curd.GetManagerPagePath("leftmenu.html")
	curd.getTemplateData()

}

// 获取模板的头部数据
func (curd *curdBaseController) getTemplateData( ){
	sideMenu := curd.GetLeftMenu()
	curd.Data["sideMenu"] = &sideMenu
	curd.HeaderData.Name = "Miller"
	curd.Data["headerData"] = curd.HeaderData
	curd.Data["tableHeader"] = &curd.DisplayTitle
}


// Curd struct
type CurdController struct {
	curdBaseController
}

// Rbac struct
type RbacController struct {
	rbacBaseController
}

// Blog struct
type BlogController struct {
	baseController
}
