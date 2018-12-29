/*
# __author__ = "Mr.chai"
# Date: 2018/12/6
*/
package main

import (
	"fmt"
	"regexp"
)

func main() {
	u := `^/www/\d+/aaab$`
fmt.Println(	regexp.Match(u, []byte("/www/4/aaab")))
	//req, err := regexp.Compile(u)
	//fmt.Println(req, err)
	//fmt.Println(req.MatchString("/www/4/aaab"))
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