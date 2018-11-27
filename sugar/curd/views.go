//
// __author__ = "Miller"
// Date: 2018/11/15
//

package curd

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 列表
func Index(c *gin.Context) {
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
