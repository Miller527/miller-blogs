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


// 定义基本数据类型, 通过数据类型匹配判断, 如果是选择和操作字段, 那么就不该做操作, 
const (
	SELECTED="SELECTED" // 选择字段

	INT="int"
	STR="str"
	SLICE="slice"
	IMG="img"
	TXT="txt"

	OPERATE="OPERATE"	// 操作字段

)