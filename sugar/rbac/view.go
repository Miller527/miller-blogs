/*
# __author__ = "Mr.chai"
# Date: 2018/12/21
*/
package rbac

import (
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
	fmt.Println("x",x,errx)
	fmt.Println("y",y,erry)
	// 验证码请求错误
	if errx != nil && erry == nil {
		c.JSON(http.StatusForbidden,ResMsg(403,"验证码错误."))
		return
	}
	session := sessions.Default(c)
	sessionCoordX := session.Get("coordX")
	sessionCoordY := session.Get("coordY")
	fmt.Println("session",sessionCoordX, reflect.TypeOf(sessionCoordX), sessionCoordY)

	// todo 用更安全的方式取判断session的正确性
	if sessionCoordY.(int) < 0 && len(sessionCoordX.([]int)) != 2  {
		c.JSON(http.StatusInternalServerError,ResMsg(500,"服务端验session验证码错误."))
		return
	}
	if x < sessionCoordX.([]int)[0] && x > sessionCoordX.([]int)[1]  || y !=  sessionCoordY{
		c.JSON(http.StatusForbidden,ResMsg(403,"验证码验证失败."))
		return
	}
	if username =="" || password == ""{
		c.JSON(http.StatusForbidden,ResMsg(403,"用户名或密码输入为空."))
		return
	}
	fmt.Println("xxxxxxx",sugar.App.DB)
	sqlCmd := `SELECT id FROM userinfo WHERE uid=? and password=?`

	stmt, err := sugar.App.DB.DefaultDB.Prepare(sqlCmd)
	if err != nil{
		c.JSON(http.StatusInternalServerError,ResMsg(500,"用户查询失败."))
		return
	}
	result, err := sugar.App.DB.SelectSlice(stmt, username, password)
	if err != nil{
		c.JSON(http.StatusInternalServerError,ResMsg(500,"用户查询失败."))
		return
	}
	if len(result) != 1{
		c.JSON(http.StatusForbidden,ResMsg(403,"用户或密码输入错误, 请重新输入."))
		return
	}

	//sqlCmd = `SELECT * FROM permission WHERE id in (SELECT menu_id FROM role_menu WHERE role_id
 //in (SELECT role_id FROM userinfo_role WHERE userinfo_id=?)) and status=1`
 //
	//stmt, err := sugar.App.DB.DefaultDB.Prepare(sqlCmd)

	fmt.Println(username,password, coordX,reflect.TypeOf(coordY))
	c.JSON(http.StatusOK,ResMsg(200,"登录成功."))


}

func ResMsg(status int, msg string)map[string]interface{}{
	return map[string]interface{}{"status":status, "msg":msg}
}
func handleLogin(c *gin.Context) {
	c.String(http.StatusOK, "handleLogin")


}

type SortedMenu [] Menu


type Menu struct {
	Id int
	Name string
	Url string
	Icon string
	Children SortedMenu
	Sort int
	IsMenu bool
	IsRegex bool
}

func xxx()  {
	//m := Menu{Children:SortedMenu{}}
	//{{}}
}