//
// __author__ = "Miller"
// Date: 2018/11/15
//

package main

import (
	_ "miller-blogs/settings"

	_ "miller-blogs/models"
	"miller-blogs/sugar/curd"
	"miller-blogs/urls"
	"fmt"

)


func main() {
	curd.AppInit(urls.AdApp,"",nil)
	err := urls.AdApp.Run("0.0.0.0:9090")
	fmt.Println(err)
}
