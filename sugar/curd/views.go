//
// __author__ = "Miller"
// Date: 2018/11/15
//

package curd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 列表
func Index(c *gin.Context) {
	fmt.Println(c.Params)

	c.String(http.StatusOK, "Index")
}

// 详情的一条
func Get(c *gin.Context) {
	c.String(http.StatusOK, "Get")
}

// 添加
func Add(c *gin.Context) {
	c.String(http.StatusOK, "Add")
}

// 更新一条
func Update(c *gin.Context) {
	c.String(http.StatusOK, "Update")
}

// 删除
func Delete(c *gin.Context) {
	c.String(http.StatusOK, "Delete")
}

// 多删
func MulitDelete(c *gin.Context){
	c.String(http.StatusOK, "MulitDelete")
}

// 多加
func MulitAdd(c *gin.Context){
	c.String(http.StatusOK, "MulitDelete")

}

// 多更
func MulitUpdate(c *gin.Context){
	c.String(http.StatusOK, "MulitDelete")

}

// 登录页面
func Login(c *gin.Context) {
	//c.String(http.StatusOK, "Login")

	c.HTML(http.StatusOK, "login.html", gin.H{
	})

}

// 登录验证
func verifyLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "verify-login.html", gin.H{
	})
}