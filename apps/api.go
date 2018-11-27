//
// __author__ = "Miller"
// Date: 2018/11/15
//

package apps

import (
	"github.com/gin-gonic/gin"
	"net/http"
)
//var dataChannel = cache.GetDataChannel()
func AdBehavior(c *gin.Context) {
	defer c.Request.Body.Close()

	//ip := c.GetHeader("x-forwarded-for")
	appKey := c.Param("appKey")

	//topic := strings.Join([]string{"ad-behavior",appKey}, ".")
	//
	//s, _ := ioutil.ReadAll(c.Request.Body)
	//byteList := bytes.Split(s,[]byte("\n"))
	//
	//for _, line := range byteList{
	//	//var ob map[string]interface{}
	//	t := &t1{}
	//	if err := json.Unmarshal(line, t); err == nil {
	//		fmt.Println("xxx",t.aaaa)
	//		//verifyField(ip, topic, ob)
	//		//mc.DataList = append(mc.DataList, ob)
	//	}
	//}
	c.String(http.StatusOK, appKey)
}

//func verifyField( ip,topic string,data map[string]interface{},){
//	data["server_time"] = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
//	data["ip"] = ip
//	data["topic"] = topic
//	dataChannel <- data
//}