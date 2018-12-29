//
// __author__ = "Miller"
// Date: 2018/11/15
//

package sugar

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"math"
	"miller-blogs/sugar/utils"
	"net/http"
	"reflect"
	"strings"
)



// 登录页面
func HandlerLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"path": App.Config.Static,
		"urlprefix": App.Config.Prefix,
		"site": "bootstrap-cerulean",
	})
}

// 登录验证
func HandlerVerifyLogin(c *gin.Context) {

	c.Redirect(http.StatusMovedPermanently, App.Config.Prefix+"index.html")
}

// 登出页面
func HandlerLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	c.Redirect(http.StatusFound, App.Config.Prefix+"login")
}

// 首页
func HandlerIndex(c *gin.Context) {
	session := sessions.Default(c)
	fmt.Println(session.Get("permission"))
	fmt.Println(session.Get("menu"))
	c.HTML(http.StatusOK, "index.html", gin.H{
		"path": App.Config.Static,

		"urlprefix": App.Config.Prefix,
		"site": "bootstrap-cerulean",
	})
}

//列表
func HandlerList(c *gin.Context) {
	dbName := c.Param(App.Config.ExtendKey)
	tbInfo, ok := App.Registry[dbName]
	fmt.Println(App.Config.relativeKey)
	fmt.Println(tbInfo)
	fmt.Println(c.Params)
	tb, ok := tbInfo[c.Param(App.Config.relativeKey)]
	if !ok {
		c.HTML(http.StatusNotFound, "error.html", gin.H{})
		return
	}
	queryCmd := fmt.Sprintf("SELECT %s FROM %s", strings.Join(tb.Field, ","), tb.Name)
	Db, ok := Dbm.DBPool[dbName]
	stmt, err := Db.Prepare(queryCmd)
	result, err := Dbm.SelectSlice(stmt, tb)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{})
		return
	}

	// todo 是不是可以在取数据那里处理, 根据权限生成这个点击字符串，同样前端的多操作按钮也要根据权限去判断生成
	var newResult [][]string
	for _, line := range result {
		v := `<i class="glyphicon glyphicon-zoom-in icon-white"></i>&nbsp;
             <i class="glyphicon glyphicon-edit icon-white"></i>&nbsp;
             <i class="glyphicon glyphicon-trash icon-white"></i>`

		newResult = append(newResult, append(line, v))
		//fmt.Println(line)
	}
	fmt.Println(result)
	res, err := json.Marshal(map[string]interface{}{
		"data": newResult,
	})
	if err == nil {
		c.String(http.StatusOK, string(res))
		return
	}
	c.String(http.StatusInternalServerError, "")
}

func HandlerCurd(c *gin.Context) {
	//var line []*TableConf
	//var tables [][]*TableConf
	//count := 0
	//for _, val := range App.Registry {
	//	if len(val.Field) >= 5 {
	//		tables = append(tables, []*TableConf{val})
	//	} else {
	//		line = append(line, val)
	//		if len(line) == 2 {
	//			tables = append(tables, line)
	//			line = []*TableConf{}
	//			continue
	//		}
	//	}
	//	if count == len(App.registry)-1 {
	//		tables = append(tables, line)
	//	}
	//}
	v := `<i id="detailBtn" class="glyphicon glyphicon-zoom-in icon-white"></i>&nbsp;
             <i id="updateBtn" class="glyphicon glyphicon-edit icon-white"></i>&nbsp;
             <i id="deleteBtn"class="glyphicon glyphicon-trash icon-white"></i>`
	c.HTML(http.StatusOK, "table.html", gin.H{
		//"tables": tables,
		"config": v,
		"site":   "static/css/bootstrap-cerulean.min.css",
	})
}

// 详情的一条
func HandlerGet(c *gin.Context) {
	c.String(http.StatusOK, "Get")
}

// 添加
func HandlerAdd(c *gin.Context) {
	fmt.Println("Params", c.Params)

	fmt.Println("c.Request.PostForm", c.Request.PostForm)
	fmt.Println("c.Request.ParseForm()", c.Request.ParseForm())
	fmt.Println("c.Request.PostForm", c.Request.PostForm)
	fmt.Println("c.Request.Body", c.Request.Body)
	fmt.Println(c.GetQuery("name"))
	fmt.Println("c.Query", c.Query("name"))
	x, _ := c.GetQueryArray("name")
	fmt.Println("type", reflect.TypeOf(x), x)

	fmt.Println(c.Query("name"))
	fmt.Println("drc", c.Query("sex"))
	c.String(http.StatusOK, "Add")
}

// 更新一条
func HandlerUpdate(c *gin.Context) {
	c.String(http.StatusOK, "Update")
}

// 删除
func HandlerDelete(c *gin.Context) {
	c.String(http.StatusOK, "Delete")
}

//// 多删
//func HandlerMulitDelete(c *gin.Context) {
//	c.String(http.StatusOK, "MulitDelete")
//}
//
//// 多加
//func HandlerMulitAdd(c *gin.Context) {
//	c.String(http.StatusOK, "MulitDelete")
//
//}
//
//// 多更
//func HandlerMulitUpdate(c *gin.Context) {
//	c.String(http.StatusOK, "MulitDelete")
//
//}

// 获取验证码
func SlideCode(c *gin.Context) {
	s := App.Config.Static

	width := 360
	height := 176
	img := []string{
		s + "slide_code/images/ver-0.png",
		s + "slide_code/images/ver-1.png",
		s + "slide_code/images/ver-2.png",
		s + "slide_code/images/ver-3.png",
	}

	img_src := img[utils.RandomInt(len(img))]

	plSize := 48
	Padding := 20
	minX := plSize + Padding
	maxX := width - Padding - plSize - plSize/6
	minY := height - Padding - plSize - plSize/6
	maxY := Padding
	deviation := 4 // 滑动偏移量

	X := RandomCoord(minX, maxX)
	Y := RandomCoord(minY, maxY)
	result := map[string]interface{}{"width": width, "height": height,
		"img_src": img_src, "pl_size": plSize,
		"padding": Padding, "x": X , "y": Y,
		"deviation": deviation}
	session := sessions.Default(c)

	session.Set("coordX",[]int{X - 10 - deviation, X-10+deviation})
	session.Set("coordY",Y)

	err  := session.Save()
	fmt.Println("save", err)
	c.JSON(http.StatusOK,result)

}

func RandomCoord(minn, maxn int) int {
	rangFloat := float64(maxn - minn)
	randFloat := utils.RandomFloat64()

	if math.Round(randFloat*rangFloat) == 0 {
		return minn + 1

	} else if int(math.Round(randFloat*float64(maxn))) == maxn {
		return maxn - 1
	} else {
		return minn + int(math.Round(randFloat*rangFloat)) - 1
	}
}
