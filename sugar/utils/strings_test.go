/*
# __author__ = "Mr.chai"
# Date: 2018/11/27
*/
package utils

import (
	"fmt"
	"testing"
)

func TestSnakeString(t *testing.T) {
	fmt.Println(SnakeString("AABc"))
	fmt.Println(SnakeString("AaBc"))
	fmt.Println(SnakeString("Aabc"))
	fmt.Println(SnakeString("AabB"))
	fmt.Println(SnakeString("AaBB"))
}

func TestCamelString(t *testing.T) {
	fmt.Println(CamelString("a_a_bc"))
	fmt.Println(CamelString("aa_bc"))
	fmt.Println(CamelString("aabc"))
	fmt.Println(CamelString("aab_b"))
	fmt.Println(CamelString("aa_b_b"))
}