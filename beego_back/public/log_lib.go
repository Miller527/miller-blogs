//
// __author__ = "Miller"
// Date: 2018/9/30
//
package public

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"path"
)

var loggers = make(map[string]*logs.BeeLogger)

var logPath string

func logRegister() {
	loggersName := beego.AppConfig.Strings("loggers")
	level, _ := beego.AppConfig.Int("log_level")
	logConfig := map[string]interface{}{
		"level":level,
		"maxlines":0,
		"maxsize":0,
	}
	for _, name := range loggersName {
		logConfig["filename"] = path.Join(logPath, name + ".log")
		byteLogConfig, err := json.Marshal(logConfig)
		newLogger := logs.NewLogger()
		if err == nil {
			newLogger.SetLogger(logs.AdapterFile, string(byteLogConfig))
			newLogger.SetLogFuncCallDepth(3)
			newLogger.EnableFuncCallDepth(true)
			loggers[name] = newLogger
		}
	}
	if len(loggers) == 0 {
		newLogger := logs.NewLogger()
		newLogger.SetLogger(logs.AdapterConsole, `{"level":1,"color":true}`)
		fmt.Println("Lor register error, use default")
	}
}

func GetLogger(name string) *logs.BeeLogger {
	fmt.Println(loggers)
	return loggers[name]
}

func init() {
	dirPath, err := DirVerify("logs")
	if err == nil {
		logPath = dirPath
		logRegister()
	} else {
		panic(err)
	}
}
