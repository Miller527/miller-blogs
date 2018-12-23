/*
# __author__ = "Mr.chai"
# Date: 2018/12/22
*/
package utils

import (
	"math/rand"
	"time"
)

func RandomInt(end int) int{
	rand.Seed(time.Now().UnixNano() )
	return rand.Intn(end)

}
func RandomRangeInt(start,end int) int{
	for  {
		res := RandomInt(end)
		if res >= start{
			return res
		}
	}
}

func RandomFloat64() float64{
	rand.Seed(time.Now().UnixNano() )
	return rand.Float64()

}