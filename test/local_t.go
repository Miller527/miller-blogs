/*
# __author__ = "Mr.chai"
# Date: 2018/12/6
*/
package main

import (
	"fmt"
	"reflect"
)


var recusive bool
var test string
var level int

func init(){
	//flag.StringVar(&test, "t", "def","rec")
	//flag.IntVar(&level, "l", 1,"rec")
	//flag.BoolVar(&recusive, "r", false,"rec")


}

func fieldIsNone(f interface{}){
	x := reflect.ValueOf(f)
	y := reflect.TypeOf(f)

	for i:=0;i<x.NumField();i++{

		fmt.Println(x.Field(i).Interface())
		fmt.Println(y.Field(i).Type)
	}

	for ii := 0;ii<x.NumMethod();ii++{
		fmt.Println(x.Method(ii).Interface())
	}
}


type UserInfo struct {
	Name string
}
func (ui UserInfo) Print(){
	fmt.Println(ui.Name)
}


func hasfunc(x interface{}, name string)bool{
	v := reflect.ValueOf(x)
	m := v.MethodByName(name)
	if m.IsValid(){
		return true
	}
	return false
}

func getfunc(x interface{}, name string)reflect.Value{
	v := reflect.ValueOf(x)
	m := v.MethodByName(name)
	return m
}

func main() {

	//x:=1.11
	//yy :=reflect.ValueOf(&x)
	//fmt.Println(yy.Kind())
	//qq := yy.Elem()
	//fmt.Println(qq.CanSet())
	//qq.SetFloat(2.22)
	//fmt.Println(x)
	//
	//u := reflect.ValueOf(UserInfo{"Miller"})
	//m:=u.Method(0)
	//m.Call([]reflect.Value{})
if hasfunc(UserInfo{"Miller"}, "Print"){
	f := getfunc(UserInfo{"Miller"}, "Print")
	args := []reflect.Value{}
	f.Call(args)
}
}


func l (x []int,y int)[]int{

	if len(x) == 0{
		return append(x, y)
	}
	for i,v := range x{
		if y < v {
			c :=  append([]int{}, x[i:]...)
			return append(append(x[:i],y), c...)
		}
	}
	return append(x,y)
}