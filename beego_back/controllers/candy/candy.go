package base

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"miller-blogs/beego_back/models"
	"path"
	"sort"
)



type Config interface {
	GetHeader()
}

type baseConfig struct {
	Model interface{}
	Header []string




	Response map[string]interface{}
}

func (conf *baseConfig) GetHeader(){

}


func NewConfig(model interface{})Config{
	c := new(baseConfig)
	c.Model = new(model)
	return c
}
// 选择框
func choiceBox (in string)(out string){
	out = in + "world"
	return
}

// 删除按钮
func deleteBtn(in string)(out string){
	out = in + "world"
	return
}
// 更新按钮
func updateBtn(in string)(out string){
	out = in + "world"
	return
}

// 发布按钮（隐藏或不隐藏）
func publishBtn(in string)(out string){
	out = in + "world"
	return
}

func RegisterFunc(){
	beego.AddFuncMap("choiceBox",choiceBox)
	beego.AddFuncMap("deleteBtn",deleteBtn)
	beego.AddFuncMap("updateBtn",updateBtn)
	beego.AddFuncMap("publishBtn",publishBtn)
}

func init(){
//	注册函数（前面的按钮）
RegisterFunc()
}










// 统一的返回值格式
type ResponseMsg struct {
	Status int                    `json:"status"`
	Msg    string                 `json:"msg"`
	Data   map[string]interface{} `json:"data"`
}

// Public base struct
type baseController struct {
	Conf *Config
	beego.Controller
}

// 按照需求重写该字段
func (base *baseController) Prepare() {
	base.DefaultInit()
}

func (base *baseController) DefaultInit() {
	base.SiteManager()
	base.Initialization()
}

// 获取orm对象
func (base *baseController) Initialization() {
	base.OrmObj = orm.NewOrm()
	base.ResponseData = ResponseMsg{}
}

// 站点管理
func (base *baseController) SiteManager() {
	base.managerSite = path.Join("manager", beego.AppConfig.String("manager_file_path"))
	base.BlogSite = "xxx"
}

// 获取返回页面文件路径
func (base *baseController) GetManagerPagePath(filename string) string {
	return path.Join(base.managerSite, filename)
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
	PermissionsData []*models.Permission
	permissionTitle []string
	PathGuidance    []map[string]string
	baseController
}

// 按照需求重写Prepare和ExtendFieldInit
func (rbac *rbacBaseController) Prepare() {
	rbac.DefaultInit()
	rbac.ExtendFieldInit()
}

// 初始化permissionTitle字段
func (rbac *rbacBaseController) ExtendFieldInit() {
	rbac.permissionTitle = []string{"id", "title", "url", "icon", "is_menu", "parent", "ButtonPid"}


}

func (rbac *rbacBaseController) DefaultInit() {
	rbac.SiteManager()
	rbac.Initialization()
	rbac.InitPathGuidance()
}


func (rbac *rbacBaseController) InitPathGuidance() {
	//rbac.PathGuidance = []map[string]string{}

	indexUrl := beego.AppConfig.String("manager_router_prefix") + "/index"
	rbac.PathGuidance = append(rbac.PathGuidance, map[string]string{"url": indexUrl, "title": "首页"})

	fmt.Println("===============================================", rbac.PathGuidance)
}


// 获取用户的权限信息
func (rbac *rbacBaseController) queryPermissions(userName string) bool {
	_, err := rbac.OrmObj.QueryTable("permission").
		Filter("Roles__Role__Users__UserInfo__Uid", userName).
		All(&rbac.PermissionsData, rbac.permissionTitle...)
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
	sideMenu := menuData{}
	json.Unmarshal([]byte(leftMenuStr), &sideMenu)
	return rbac.orderMenu(sideMenu)

}

// 菜单排序
func (rbac *rbacBaseController) orderMenu(menu menuData) (sideMenu []menuData) {
	var sortKeys []string
	for k := range menu {

		sortKeys = append(sortKeys, k)
	}
	sort.Strings(sortKeys)

	for _, k := range sortKeys {
		menuChildren := menu[k].(map[string]interface{})["children"]
		if len(menuChildren.(map[string]interface{})) == 0 {
			delete(menu, k)
			continue
		}

		rbac.updateMenuStatus(menu[k].(map[string]interface{}))
		sideMenu = append(sideMenu, menu[k].(map[string]interface{}))
	}
	return
}

// 菜单状态修改
func (rbac *rbacBaseController) updateMenuStatus(menu menuData) () {
	if len(rbac.Ctx.Request.Header["menus-id"]) == 1 {
		subMenuIdStr := rbac.Ctx.Request.Header["menus-id"][0]


		val, ok := menu["children"].(map[string]interface{})[subMenuIdStr]
		if ok {
			menuUrl := func() string{
				if 	val, ok := menu["url"]; ok && val != nil {
					fmt.Println(val,ok)
					return val.(string)
				}else {
					return ""
				}

			}()

			rbac.PathGuidance = append(rbac.PathGuidance,
				map[string]string{"url": menuUrl, "title": menu["title"].(string)})


			rbac.PathGuidance = append(rbac.PathGuidance,
				map[string]string{"url": val.(map[string]interface{})["url"].(string),
					"title": val.(map[string]interface{})["title"].(string)})
			fmt.Println("------------------------------", rbac.PathGuidance)
			val.(map[string]interface{})["class"] = "current"
			menu["class"] = "selected"
			menu["style"] = "display: block"
		}

		//for _, val := range menu["children"].(map[string] interface{}) {
		//	fmt.Println(rbac.Ctx.Request.RequestURI)
		//	fmt.Println(val.(map[string]interface{})["url"])
		//	if val.(map[string]interface{})["url"] == rbac.Ctx.Request.RequestURI {
		//		val.(map[string]interface{})["class"] =
		//		menu["class"] = "selected"
		//		menu["style"] = "display: block"
		//		break
		//	}
		//	}else {
		//		menu["style"] = "display: none"
		//	}
		//
		//}

	}
}

// 初始化session字典
func (rbac *rbacBaseController) sessionInit(uid string) {
	go rbac.SetSession(beego.AppConfig.String("session_uid_key"), uid)
	var permissionsList []map[string]interface{}

	menu := menuDict{}

	for _, obj := range rbac.PermissionsData {

		permission := map[string]interface{}{
			"id":    obj.Id,
			"title": obj.Title,
			"url":   obj.Url,
			"parent": func() interface{} {
				if obj.Parent != nil {
					return obj.Parent.Id
				}
				return nil

			}(),
			"button_pid": func() interface{} {
				if obj.ButtonPid != nil {
					return obj.ButtonPid.Id
				}
				return nil
			}(),
		}
		if obj.Url != "" {
			// 过滤一级菜单（必须有url的才是基本权限）
			permissionsList = append(permissionsList, permission)
		}

		if obj.Url == "" && obj.IsMenu {

			// 一级菜单的判断
			if _, ok := menu[obj.Id]; !ok {
				menu[obj.Id] = menuData{}
				menu[obj.Id]["children"] = menuDict{}
			}
			menu[obj.Id]["title"] = obj.Title
			menu[obj.Id]["icon"] = obj.Icon

		} else if obj.Url != "" && obj.Parent != nil {
			if _, ok := menu[obj.Parent.Id]; !ok {
				menu[obj.Parent.Id]["children"] = menuDict{}
			}
			if obj.ButtonPid == nil {
				// 二级菜单的判断
				if _, ok := menu[obj.Parent.Id]["children"].(menuDict)[obj.Id]; !ok {
					menu[obj.Parent.Id]["children"].(menuDict)[obj.Id] = menuData{}
					menu[obj.Parent.Id]["children"].(menuDict)[obj.Id]["button"] = []string{}
				}
				menu[obj.Parent.Id]["children"].(menuDict)[obj.Id]["url"] = obj.Url
				menu[obj.Parent.Id]["children"].(menuDict)[obj.Id]["title"] = obj.Title
			} else {
				// button的判断
				if _, ok := menu[obj.Parent.Id]["children"].(menuDict)[obj.ButtonPid.Id]; !ok {
					menu[obj.Parent.Id]["children"].(menuDict)[obj.ButtonPid.Id] = menuData{}
					menu[obj.Parent.Id]["children"].(menuDict)[obj.ButtonPid.Id]["button"] = []string{}

				}
				tmp := append(menu[obj.Parent.Id]["children"].(menuDict)[obj.ButtonPid.Id]["button"].([]string),
					obj.Url)
				menu[obj.Parent.Id]["children"].(menuDict)[obj.ButtonPid.Id]["button"] = tmp
			}

		}
	}

	if byteStr, err := json.Marshal(menu); err == nil {
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
	DisplayTitle []string //前端显示的表字段名字
	FieldTitle   []string //数据库的字段名，和上面的显示字段从前一一对应, 长度可以不一样，后边缺少的字段是前端自定义的字段
	HeaderData   *HeaderData
	rbacBaseController
}

// 按照需求重写Prepare和ExtendFieldInit
func (curd *curdBaseController) Prepare() {
	curd.ExtendFieldInit()
	curd.DefaultInit()

}

func (curd *curdBaseController) ExtendFieldInit() {}


func (curd *curdBaseController) DefaultInit() {
	curd.SiteManager()
	curd.Initialization()
	curd.InitPathGuidance()
	curd.initUid()
}


func (curd *curdBaseController) initUid() {
	uid := curd.GetSession(beego.AppConfig.String("session_uid_key"))
	curd.HeaderData = &HeaderData{uid.(string)}
}

// 获取model的自定义的默认字段
func (curd *curdBaseController) DefaultFiledTitles(model struct{}) {

}

// 生成模范的返回数据
func (curd *curdBaseController) ResponseTemplate(htmlName string) {
	curd.Layout = curd.GetManagerPagePath("head_start.html")
	curd.TplName = curd.GetManagerPagePath(htmlName)
	curd.LayoutSections = make(map[string]string)
	curd.LayoutSections["HeadMeta"] = curd.GetManagerPagePath("headmeta.html")
	curd.LayoutSections["Header"] = curd.GetManagerPagePath("head.html")
	curd.LayoutSections["LeftMenu"] = curd.GetManagerPagePath("menu.html")
	curd.LayoutSections["Navigation"] = curd.GetManagerPagePath("topbar.html")
	curd.getTemplateData()

}

// 获取模板的头部数据
func (curd *curdBaseController) getTemplateData() {
	sideMenu := curd.GetLeftMenu()
	curd.Data["sideMenu"] = &sideMenu
	curd.Data["pathGuidance"] = &curd.PathGuidance
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
