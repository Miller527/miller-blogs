//
// __author__ = "Miller"
// Date: 2018/11/15
//

package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)


func AppKeyFilter() gin.HandlerFunc {

	//
	//
	//if err := json.Unmarshal(mc.Ctx.Input.RequestBody, &ob); err == nil {
	//	mc.DataList = append(mc.DataList, ob)
	//} else {
	//	str := string(mc.Ctx.Input.RequestBody)
	//	jsonStrList := strings.Split(str, "\\n")
	//	fmt.Println("jsonStrList", jsonStrList)
	//	for _, data := range jsonStrList {
	//		if err := json.Unmarshal([]byte(data), &ob); err == nil {
	//			mc.DataList = append(mc.DataList, ob)
	//		}
	//	}
	//}
	return func(c *gin.Context) {


		fmt.Println(c.Params)
	}
}

