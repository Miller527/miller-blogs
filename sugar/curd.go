//
// __author__ = "Miller"
// Date: 2018/11/15
//

package sugar

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// 首页

func HandlerIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
	})
}

//列表
func HandlerList(c *gin.Context) {
	tb, ok := Registry[	c.Param(RelativePath)]
	fmt.Println(tb)
	if !ok{
		c.HTML(http.StatusNotFound, "error.html", gin.H{})
		return
	}
	queryCmd := fmt.Sprintf("SELECT %s FROM %s",strings.Join(tb.Field, ","), tb.Name())

	stmt, err := Dbm.Db.Prepare(queryCmd)
	fmt.Println(stmt,err)
	fmt.Println(strings.Join(tb.Field,","))
	result,err := Dbm.SelectSlice(stmt, tb)
	fmt.Println(result,err)
	if err != nil  {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{})
		return
	}
	fmt.Println(result)
	res,err := json.Marshal(map[string][][]string{"data":result})
	if err == nil{
		c.String(http.StatusOK, string(res))
		return
	}
	c.String(http.StatusInternalServerError, "")
}

func HandlerCurd(c *gin.Context) {
fmt.Println(	Registry)

	var tn []string
	for name, _ := range Registry{
		tn = append(tn, name)
		fmt.Println(tn)
	}
	c.HTML(http.StatusOK, "table.html", gin.H{
		"tables":Registry,
	})
}

// 详情的一条
func HandlerGet(c *gin.Context) {
	c.String(http.StatusOK, "Get")
}

// 添加
func HandlerAdd(c *gin.Context) {
	c.String(http.StatusOK, "Add")
}

// 更新一条
func HandlerUpdate(c *gin.Context) {
	c.String(http.StatusOK, "Update")
}

// 删除
func HandlerDelete(c *gin.Context) {
	c.String(http.StatusOK, "Delete")
}

// 多删
func HandlerMulitDelete(c *gin.Context){
	c.String(http.StatusOK, "MulitDelete")
}

// 多加
func HandlerMulitAdd(c *gin.Context){
	c.String(http.StatusOK, "MulitDelete")

}

// 多更
func HandlerMulitUpdate(c *gin.Context){
	c.String(http.StatusOK, "MulitDelete")

}

// 登录页面
func HandlerLogin(c *gin.Context) {
	//c.String(http.StatusOK, "Login")

	c.HTML(http.StatusOK, "login.html", gin.H{
		"path": "/sugar/static",

	})

}

// 登录验证
func HandlerVerifyLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "verify-login.html", gin.H{
		"path": "/sugar/static",
	})
}