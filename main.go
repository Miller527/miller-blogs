package main

import (
	_ "github.com/Go-SQL-Driver/MySQL"
	_ "github.com/astaxie/beego/session/redis"

	_ "miller-blogs/middleware"
	_ "miller-blogs/models"
	_ "miller-blogs/public" // 注册日志
	_ "miller-blogs/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func main() {
	orm.Debug = true
	beego.Run()

}
