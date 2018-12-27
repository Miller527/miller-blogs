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
	"reflect"
	"strconv"
)

func handleVerifyLogin(c *gin.Context) {
	fmt.Println("handleVerifyLogin")
	username := c.PostForm("username")
	password := c.PostForm("password")
	coordX := c.PostForm("coordX")
	coordY := c.PostForm("coordY")
	x, errx := strconv.Atoi(coordX)
	y, erry := strconv.Atoi(coordY)
	fmt.Println("x", x, errx)
	fmt.Println("y", y, erry)
	// 验证码请求错误
	if errx != nil && erry == nil {
		c.JSON(http.StatusForbidden, ResMsg(403, "验证码错误."))
		return
	}
	session := sessions.Default(c)
	sessionCoordX := session.Get("coordX")
	sessionCoordY := session.Get("coordY")
	fmt.Println("session", sessionCoordX, reflect.TypeOf(sessionCoordX), sessionCoordY)

	// todo 用更安全的方式取判断session的正确性
	if sessionCoordY.(int) < 0 && len(sessionCoordX.([]int)) != 2 {
		c.JSON(http.StatusInternalServerError, ResMsg(500, "服务端验session验证码错误."))
		return
	}
	if x < sessionCoordX.([]int)[0] && x > sessionCoordX.([]int)[1] || y != sessionCoordY {
		c.JSON(http.StatusForbidden, ResMsg(403, "验证码验证失败."))
		return
	}
	if username == "" || password == "" {
		c.JSON(http.StatusForbidden, ResMsg(403, "用户名或密码输入为空."))
		return
	}
	fmt.Println("xxxxxxx", sugar.App.DB)
	sqlCmd := `SELECT id FROM userinfo WHERE uid=? and password=? AND status=1`

	stmt, err := sugar.App.DB.DefaultDB.Prepare(sqlCmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResMsg(500, "用户查询失败."))
		return
	}
	result, err := sugar.App.DB.SelectValues(stmt, username, password)
	fmt.Println(result)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ResMsg(500, "用户查询失败."))
		return
	}
	if len(result) != 1 {
		c.JSON(http.StatusForbidden, ResMsg(403, "用户或密码输入错误, 请重新输入."))
		return
	}
	fmt.Println(result)
	sqlCmd = `SELECT * FROM permission WHERE id in (SELECT permission_id FROM role_permission WHERE role_id
 in (SELECT role_id FROM userinfo_role WHERE userinfo_id=?)) and status=1`
	stmt, err = sugar.App.DB.DefaultDB.Prepare(sqlCmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResMsg(500, "用户查询失败."))
		return
	}
	permissions, err := sugar.App.DB.SelectDict(stmt, result[0])
	MenuList(permissions)
	fmt.Println("xxx", permissions)
	fmt.Println(username, password, coordX, reflect.TypeOf(coordY))
	c.JSON(http.StatusOK, ResMsg(200, "登录成功."))

}
func disposeMenus(menus SortedMenu, child *Menu, lev bool) SortedMenu {

	for i, m := range menus {
		fmt.Println(child.Id ,child.ParentId ,m.Id)
		if child.ParentId == m.Id {
			fmt.Println(child.ParentId ,m.Id,i)
			m.Children = SortedInsert(m.Children, child)
			fmt.Println(m.Children)
			return menus
		}else if len(m.Children) != 0 {
			m.Children = disposeMenus(m.Children, child, false)
			for _, v := range m.Children{
				if v.Id == child.Id {
					return menus
				}
			}
		} else {
			if i == len(menus)-1 && lev{
				menus = append(menus, child)
				return menus
			}
		}
	}
	return menus
}

func MenuList(permiss []map[string]interface{}) SortedMenu {
	var menus SortedMenu

	for _, line := range permiss {
		fmt.Println("line-------------", line)
		v, e := json.Marshal(line)
		menu := &Menu{}
		if e != nil {
			fmt.Println("e----------------", e)
			continue
		}
		err := json.Unmarshal(v, menu)
		if err != nil {
			fmt.Println("err----------------", err)
			continue

		}
		if len(menus) == 0 {
			menus = append(menus, menu)
			continue
		}
		menus = disposeMenus(menus, menu, true)
	}

	fmt.Println("xxxxxxxxxxxxx", menus)
	jsonBytes, _ := json.Marshal(menus)
	fmt.Println(string(jsonBytes))
	return menus
}



func SortedInsert(menus SortedMenu, menu *Menu) SortedMenu {
	if len(menus) == 0 {
		return append(menus, menu)
	}

	for i, v := range menus {
		if menu.Sort < v.Sort {
			s := append(SortedMenu{}, menus[i:]...)
			return append(append(menus[:i], menu), s...)
		}
	}
	return append(menus, menu)
}
func ResMsg(status int, msg string) map[string]interface{} {
	return map[string]interface{}{"status": status, "msg": msg}
}
func handleLogin(c *gin.Context) {
	c.String(http.StatusOK, "handleLogin")

}

type SortedMenu [] *Menu

// todo 这里的类型是否能够改成正常的类型
type Menu struct {
	Id       int
	Title    string
	Url      string
	Icon     string
	Children SortedMenu
	ParentId int `json:"parent_id"`
	Sort     int
	IsMenu   int `json:"is_menu"`
	IsRegex  int `json:"is_regex"`
}

func xxx() {
	//m := Menu{Children:SortedMenu{}}
	//{{}}
}
