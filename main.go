//
// __author__ = "Miller"
// Date: 2018/11/15
//

package main

import (
	"fmt"
	_ "miller-blogs/models"
	"miller-blogs/sugar/curd"
	"miller-blogs/urls"
)


func main() {
	curd.AppInit(urls.AdApp,"",nil)
	err := urls.AdApp.Run("0.0.0.0:9090")
	fmt.Println(err)
}

//type user struct{
//	name string
//	age int
//	feature map[string]interface{}
//}
//func main() {
//	var u interface{}
//	u=&user{}
//	value:=reflect.ValueOf(u)
//	if value.Kind()==reflect.Ptr{
//		elem:=value.Elem()
//		name:=elem.FieldByName("name")
//		if name.Kind()==reflect.String{
//			*(*string)(unsafe.Pointer(name.Addr().Pointer())) = "fangwendong"
//		}
//
//		age:=elem.FieldByName("age")
//		if age.Kind()==reflect.Int{
//			*(*int)(unsafe.Pointer(age.Addr().Pointer())) =24
//		}
//
//		feature:=elem.FieldByName("feature")
//		if feature.Kind()==reflect.Map{
//			*(*map[string]interface{})(unsafe.Pointer(feature.Addr().Pointer())) =map[string]interface{}{
//				"爱好":"篮球",
//				"体重":60,
//				"视力":5.2,
//			}
//		}
//
//	}
//
//	fmt.Println(u)
//}