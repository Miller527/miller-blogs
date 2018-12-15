package beego_back

import (
	"fmt"
	"net/http"
)

//import (
//	_ "github.com/Go-SQL-Driver/MySQL"
//	_ "github.com/astaxie/beego/session/redis"
//
//	_ "miller_blogs/middleware"
//	_ "miller_blogs/models"
//	_ "miller_blogs/public" // 注册日志
//	_ "miller_blogs/routers"
//
//	"github.com/astaxie/beego"
//	"github.com/astaxie/beego/orm"
//)

//func main() {
//	//orm.Debug = true
//	//beego.Run()
//	var x []int
//	if x == nil{
//		fmt.Println("x")
//	}
//}

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello Wrold!</h1>"))
	fmt.Println("Hello Wrold!")
}
func main() {
	//i := debug.SetMaxThreads(65536)
	//i = debug.SetMaxThreads(65536)
	//fmt.Println(i)
	http.HandleFunc("ad-behavior/v1/p2n4n7osvqsgttuw", sayHello)
	http.ListenAndServe("0.0.0.0:8000", nil)
	x :=&A{"xxxxx","q"}
	fmt.Println(x)
	var y []*A
	y = append(y, x)
	fmt.Println(y)
}

type A struct {
	host string
	port string
}
func (a *A) String() string {
	fmt.Println("a")
	return a.host+"qqq"+a.port
}



func f(x string){
	fmt.Println("xxx")
}