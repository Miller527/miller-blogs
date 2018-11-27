/*
# __author__ = "Mr.chai"
# Date: 2018/11/25
*/
package curd

func InSlice(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}