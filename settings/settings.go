//
// __author__ = "Miller"
// Date: 2018/11/15
//

package settings



import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type configuration struct {
	Mode           string
	Version        string
	AllKafkaHosts  map[string][]string `json:"kafkahosts"`
	Kafkatimeout   time.Duration
	Kafkaconnsleep time.Duration
}

func (conf *configuration) GetKafkaHosts() ([]string, error) {
	val, ok := conf.AllKafkaHosts[conf.Mode]
	fmt.Println(val, conf.Mode, conf.AllKafkaHosts)
	if ok {
		return val, nil
	}
	return nil, errors.New("SettingsError: Not found kafka hosts")
}

var Settings configuration

func init() {
	//path.Join(dir, "settings",config.json)
	file, err := os.Open("settings/config.json")
	if err != nil {
		panic(err)
	}
	//关闭文件
	defer file.Close()

	//NewDecoder创建一个从file读取并解码json对象的*Decoder，解码器有自己的缓冲，并可能超前读取部分json数据。
	decoder := json.NewDecoder(file)

	Settings = configuration{}
	//Decode从输入流读取下一个json编码值并保存在v指向的值里
	err = decoder.Decode(&Settings)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println(Settings)
}