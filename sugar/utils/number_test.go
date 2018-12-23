/*
# __author__ = "Mr.chai"
# Date: 2018/12/22
*/
package utils

import (
	"fmt"
	"testing"
)

func TestRandomRangeInt(t *testing.T) {
	fmt.Println(RandomRangeInt(10,40))
}
func TestRandomFloat64(t *testing.T) {
	fmt.Println(RandomFloat64())
}