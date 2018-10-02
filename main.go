package main

import (
	_ "miller-blogs/public"
	_ "miller-blogs/middleware"
	_ "github.com/Go-SQL-Driver/MySQL"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "miller-blogs/models"
	_ "miller-blogs/routers"
)


func main() {
	orm.Debug = true
	beego.Run()

}


