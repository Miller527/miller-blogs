/*
# __author__ = "Mr.chai"
# Date: 2018/11/30
*/
package utils

func InStringSlice(val string, sl []string) bool {

	for _, v := range sl {
		if v == val {
			return true
		}
	}
	return false
}
func InIntSlice(val int, sl []int) bool {

	for _, v := range sl {
		if v == val {
			return true
		}
	}
	return false
}
func InInterfaceSlice(val interface{}, sl []interface{}) bool {

	for _, v := range sl {
		if v == val {
			return true
		}
	}
	return false
}

func DelStringSliceEle(sl []string, val string) []string {
	for i := 0; i < len(sl); i++ {
		if sl[i] == val {
			sl = append(sl[:i], sl[i+1:]...)
			break
		}
	}
	return sl
}
func DelStringSliceEles(sl []string, vals... string) []string {
	for _, v := range vals{
		for i := 0; i < len(sl); i++ {
			if sl[i] == v {
				sl = append(sl[:i], sl[i+1:]...)
				i--
			}
		}
	}
	return sl
}

func EleIndexStringSlice(v string, sl []string)int{
	for i,val :=range sl{
		if val == v {
			return i
		}
	}
	return -1
}