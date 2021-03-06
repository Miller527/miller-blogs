/*
# __author__ = "Mr.chai"
# Date: 2018/12/14
*/
package sugar
//
//type RawFilter interface {
//	match(val string, re string)bool
//	Filter(rawData []string, reList []string)[]int
//}
//
//
//type regexFilter struct {
//	role map[string]string
//}
//
//// 正则表达式
//func(rf regexFilter) match(val string, re string)bool{
//	return true
//}
//
//
//func (rf regexFilter) filter(rawData []string, fieldType []string)[]int{
//	return []int{}
//}
//
//var filterRole = map[string]string{	"x":"xx","u":"uu",	"p":"pp",	"q":"qq"}

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


var defaultFilter iFieldFilter



type iFieldFilter interface {
	// 过滤器, 返回msg和状态
	Filter(str string) (string, bool)
	// todo 处理所有
	FilterAll(fields FilterItem) (map[string]string, bool)
}

type FilterItem map[string]FilterInfo

type Serializer struct {
	iFieldFilter
}

type FilterInfo struct {
	Rule string
	Length  int
}

// 过滤器, 类型为正则表达式, 还没想好其他的过滤方式
type RegexFilter struct {
	Rules map[string]string
}

func (rgf Serializer) Filter(str string) (string, bool){

	return "",true
}
func (rgf Serializer) FilterAll(items FilterItem) (map[string]string, bool){

	return nil,true
}

func GetFilter()  {
}

func SetFilter(filter iFieldFilter){
	defaultFilter = defaultFilter
}

func init() {
	defaultFilter = Serializer{}
}