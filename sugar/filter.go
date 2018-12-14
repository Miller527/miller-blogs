/*
# __author__ = "Mr.chai"
# Date: 2018/12/14
*/
package sugar


type RawFilter interface {
	match(val string, re string)bool
	Filter(rawData []string, reList []string)[]int
}


type regexFilter struct {
	role map[string]string
}

// 正则表达式
func(rf regexFilter) match(val string, re string)bool{
	return true
}


func (rf regexFilter) filter(rawData []string, fieldType []string)[]int{
	return []int{}
}

var filterRole = map[string]string{	"x":"xx","u":"uu",	"p":"pp",	"q":"qq"}
var defaultFilter = regexFilter{filterRole}

const (
	INT="int"
	STR="str"
	SLICE="slice"
	INT="int"
	INT="int"
	INT="int"
	INT="int"
	INT="int"

)