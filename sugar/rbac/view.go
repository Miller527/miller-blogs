/*
# __author__ = "Mr.chai"
# Date: 2018/12/21
*/
package rbac

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"miller-blogs/sugar"
	"net/http"
	"strconv"
)

func handlerVerifyLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	coordX := c.PostForm("coordX")
	coordY := c.PostForm("coordY")
	x, errx := strconv.Atoi(coordX)
	y, erry := strconv.Atoi(coordY)

	// 验证码请求错误
	if errx != nil && erry == nil {
		c.JSON(http.StatusForbidden, ResMsg(403, "验证码错误."))
		return
	}
	// 获取登录验证码
	session := sessions.Default(c)
	sessionCoordX := session.Get("coordX")
	sessionCoordY := session.Get("coordY")

	// todo 用更安全的方式取判断session的正确性
	if sessionCoordY.(int) < 0 && len(sessionCoordX.([]int)) != 2 {
		c.JSON(http.StatusInternalServerError, ResMsg(500, "服务端session验证码生成错误."))
		return
	}
	if x < sessionCoordX.([]int)[0] && x > sessionCoordX.([]int)[1] || y != sessionCoordY {
		c.JSON(http.StatusForbidden, ResMsg(403, "验证码验证失败."))
		return
	}

	sqlCmd := `SELECT id FROM userinfo WHERE uid=? and password=? AND status=1`

	stmt, err := sugar.App.DB.DefaultDB.Prepare(sqlCmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResMsg(500, "用户查询失败."))
		return
	}
	result, err := sugar.App.DB.SelectValues(stmt, username, password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ResMsg(500, "用户查询失败."))
		return
	}
	if len(result) != 1 {
		c.JSON(http.StatusForbidden, ResMsg(403, "用户或密码输入错误, 请重新输入."))
		return
	}
	sqlCmd = `SELECT * FROM permission WHERE id in (SELECT permission_id FROM role_permission WHERE role_id
			  in (SELECT role_id FROM userinfo_role WHERE userinfo_id=?)) and status=1`
	stmt, err = sugar.App.DB.DefaultDB.Prepare(sqlCmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResMsg(500, "用户查询失败."))
		return
	}
	reult, err := sugar.App.DB.SelectDict(stmt, result[0])
	menuJson, permiss := MenuList(reult)
	if menuJson == "" && permiss == "" {
		c.JSON(http.StatusOK, ResMsg(403, "没有任何访问权限."))
		return
	}
	session.Set("menu", menuJson)
	session.Set("permission", permiss)
	session.Set("username", username)
	err = session.Save()
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusOK, ResMsg(500, "权限初始化失败, 请稍后再试."))
		return
	}
	c.JSON(http.StatusOK, ResMsg(200, "登录成功."))

}

// 处理菜单列表
func disposeMenus(menus sugar.SortedMenu, child *sugar.Menu, pid int) sugar.SortedMenu {
	for i, m := range menus {
		if child.ParentId == m.Id {
			m.Children = SortedInsert(m.Children, child)
			return menus
		} else if len(m.Children) != 0 {
			m.Children = disposeMenus(m.Children, child, m.Id)
			for _, v := range m.Children {
				if v.Id == child.Id {
					return menus
				}
			}
		}
		if i+1 == len(menus) && pid == child.ParentId {
			menus = append(menus, child)
		}
	}
	return menus
}

// 生成菜单列表和权限表
func MenuList(pList []map[string]interface{}) (string, string) {
	fmt.Println(pList)
	var menus sugar.SortedMenu
	var permiss = &Permissions{}
	for _, line := range pList {

		isRegex, rok := line["is_regex"]
		isMenu, mok := line["is_menu"]
		url, uok := line["url"]

		if rok && mok && uok {
			if isRegex.(int) == 1 {
				permiss.Regex = append(permiss.Regex, url.(string))
			} else {
				permiss.Static = append(permiss.Static, url.(string))

			}
		}
		if isMenu.(int) != 1 {
			continue
		}
		v, e := json.Marshal(line)
		if e != nil {
			return "", ""
		}

		menu := &sugar.Menu{}
		err := json.Unmarshal(v, menu)
		if err != nil {
			return "", ""

		}
		if len(menus) == 0 {
			menus = append(menus, menu)
			continue
		}
		menus = disposeMenus(menus, menu, 0)
	}
	menuByte, err := json.Marshal(menus)
	if err != nil {
		return "", ""
	}
	permissByte, err := json.Marshal(permiss)
	if err != nil {
		return "", ""
	}
	return string(menuByte), string(permissByte)
}

// 顺序插入
func SortedInsert(menus sugar.SortedMenu, menu *sugar.Menu) sugar.SortedMenu {
	if len(menus) == 0 {
		return append(menus, menu)
	}

	for i, v := range menus {
		if menu.Sort < v.Sort {
			s := append(sugar.SortedMenu{}, menus[i:]...)
			return append(append(menus[:i], menu), s...)
		}
	}
	return append(menus, menu)
}

// 获取状态信息
func ResMsg(status int, msg string) map[string]interface{} {
	return map[string]interface{}{"status": status, "msg": msg}
}

func handlerLogin(c *gin.Context) {
	session := sessions.Default(c)
	permiss := Permissions{Static:[]string{ParamsRbac.urlPrefix + "slidecode"}}
	permissByte, _ := json.Marshal(permiss)

	session.Set("permission", string(permissByte))
	err := session.Save()
	if err != nil {
		c.JSON(http.StatusOK, ResMsg(500, "权限初始化失败, 请稍后再试."))
		return
	}
	c.HTML(http.StatusOK, "login.html", gin.H{
		"path": ParamsRbac.staticPath,
		"urlprefix": ParamsRbac.urlPrefix,
		"site": "bootstrap-cerulean",
	})
}


type Permissions struct {
	Static []string
	Regex  []string
}

//func init(){
//gob.Register(Permissions)
//}
