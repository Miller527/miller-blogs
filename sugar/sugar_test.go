/*
# __author__ = "Mr.chai"
# Date: 2018/12/13
*/
package sugar

import (
	"fmt"
	"testing"
)


type TestRole struct {
	Id   []uint `title:"ID" primary:"true" filter:"这里是个正则表达式用来做输入过滤用"`
	Rid  string `title:"角色ID"`
	Name string
	Table *TableConf
}

func (role TestRole) DisplayName ()string{
	return "权限表"
}
var t = &TableConf{true,false,[]int{11}}
var x = TestRole{Table:t}
func TestDefaultDescAnalyzer(t *testing.T) {

	ana := defaultDescAnalyzer{}
	tc := ana.ParseDesc(x)
	fmt.Println("Name",tc.Name)
	fmt.Println("Display",tc.Display)
	fmt.Println("Field",tc.Field)
	fmt.Println("Title",tc.Title)
	fmt.Println("Methods",tc.Methods)
}