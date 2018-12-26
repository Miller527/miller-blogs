/*
# __author__ = "Mr.chai"
# Date: 2018/12/6
*/
package main

import (
	"fmt"
)


func main() {
	arr := []int{6, 5, 4, 3, 2, 1, 0}

	a2 := []int{}

	for _ ,v :=range arr{
		a2 = l(a2, v)
		fmt.Println("Sorted arr: ", a2)

	}
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
	return x
}