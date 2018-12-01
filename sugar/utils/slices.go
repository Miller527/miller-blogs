/*
# __author__ = "Mr.chai"
# Date: 2018/11/30
*/
package utils

func InStringSlice(v string, sl []string) bool {

	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}
func InIntSlice(v int, sl []int) bool {

	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}