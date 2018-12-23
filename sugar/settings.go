//
// __author__ = "Miller"
// Date: 2018/11/15
//

package sugar

import (
	"encoding/json"
	"errors"
	"fmt"
	"miller-blogs/sugar/utils"
	"os"
	"time"
)

type configuration struct {
	Mode           string
	Version        string
	AllKafkaHosts  map[string][]string `json:"kafkahosts"`
	Kafkatimeout   time.Duration
	Kafkaconnsleep time.Duration
	DBConfig      DBConfig
	Session       SessionConfig
}

func (conf *configuration) GetKafkaHosts() ([]string, error) {
	val, ok := conf.AllKafkaHosts[conf.Mode]
	fmt.Println(val, conf.Mode, conf.AllKafkaHosts)
	if ok {
		return val, nil
	}
	return nil, errors.New("SettingsError: Not found kafka hosts")
}

var settings configuration

func Settings(confPath string) {

	if confPath == ""{
		confPath = "./settings/config.json"
	}
	file, err := os.Open(confPath)
	defer file.Close()

	utils.PanicCheck(err)

	decoder := json.NewDecoder(file)

	settings = configuration{}

	err = decoder.Decode(&settings)
	utils.PanicCheck(err)
}


