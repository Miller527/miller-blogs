/*
# __author__ = "Mr.chai"
# Date: 2018/9/9
*/
package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
)

type IndexController struct {
	beego.Controller
}

func (this *IndexController) Get () {
	fmt.Println(this.Data)
	this.TplName = "index.html"
	beego.Run()
}