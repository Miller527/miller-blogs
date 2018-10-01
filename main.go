package main

import (
	"encoding/gob"
	_ "github.com/Go-SQL-Driver/MySQL"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "miller-blogs/models"
	_ "miller-blogs/routers"
)


func init()  {
	gob.Register(&[]orm.ParamsList{})//必须在encoding/gob编码解码前进行注册
}
func main() {
	//ormObj := orm.NewOrm()
	//qs := ormObj.QueryTable("user")
	//
	//
	//a :=qs.Filter("name", "Miller")
	//fmt.Println(a.Count())
	//user := models.User{}
	//fmt.Println(a.One(&user,"profile__age"))
	//fmt.Println(user)
	//fmt.Println(qs.Exclude("profile__isnull", true).Filter("name","miller"))
	//var maps []orm.Params
	//ormObj.QueryTable("user").Values(&maps, "id", "name", "profile", "profile__age")
	//fmt.Println(maps)
	orm.Debug = true
	beego.Run()


	//fmt.Println([]byte("Hello World"))
}


//func main() {
	//beego.SetStaticPath("/img", "static")
	//beego.SetStaticPath("/down2", "download2")
	//beego.Run()

//}

