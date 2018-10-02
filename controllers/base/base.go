package base

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
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
	SitePath     string //站点文件路径
	OrmObj       orm.Ormer
	ResponseData ResponseMsg
	BlogLogger   *logs.BeeLogger
}

// 按照需求重写该字段
func (base *baseController) Prepare() {
	base.Initialization()
}

// 获取orm对象
func (base *baseController) Initialization() {
	base.OrmObj = orm.NewOrm()
	base.ResponseData = ResponseMsg{}

}

// 站点管理
func (base *baseController) SiteManager() {

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

// CURD base struct
type curdBaseController struct {
	DisplayTitle []string //前端显示的表字段名字
	FieldTitle   []string //数据库的字段名，和上面的显示字段从前一一对应, 长度可以不一样，后边缺少的字段是前端自定义的字段
}

// RBAC base struct
type rbacBaseController struct {
}

// 用户登录认证
func (rbac *rbacBaseController) Login() {

}

// Blog struct
type BlogController struct {
	baseController
	curdBaseController
	rbacBaseController
}

// Curd struct
type CurdController struct {
	baseController
	curdBaseController
}

// Rbac struct
type RbacController struct {
	baseController
	rbacBaseController
}
