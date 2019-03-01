/*
# __author__ = "Mr.chai"
# Date: 2018/12/6
*/
package main

import (
	"flag"
	"fmt"
)


var recusive bool
var test string
var level int

func init(){
	flag.StringVar(&test, "t", "def","rec")
	flag.IntVar(&level, "l", 1,"rec")
	flag.BoolVar(&recusive, "r", false,"rec")

	flag.Parse()

}

func main() {

	fmt.Printf("%v\n", recusive)
	fmt.Printf("%v\n", test)
	fmt.Printf("%v\n", level)

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