/*
# __author__ = "Mr.chai"
# Date: 2018/12/6
*/
package main

import (
	"fmt"
	"strings"
)

func main() {
	u := "name"
	fmt.Println(strings.Split(u,"-"))
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