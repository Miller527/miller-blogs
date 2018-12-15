/*
# __author__ = "Mr.chai"
# Date: 2018/12/15
*/
package utils

func PanicCheck(err error){
	if err != nil{
		panic(err)
	}
}