/*
# __author__ = "Mr.chai"
# Date: 2018/12/6
*/
package utils

import (
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()

	// 禁止 logrus 的输出
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err!= nil{
		fmt.Println("err", err)
	}
	Log.Out = src
	Log.SetLevel(logrus.DebugLevel)
	apiLogPath := "api.log"
	logWriter, err := rotatelogs.New(
		strings.Split(apiLogPath,".")[0]+".%Y-%m-%d.log",
		rotatelogs.WithLinkName(apiLogPath), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(30*24*time.Hour), // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Second), // 日志切割时间间隔
	)
	fmt.Println("x")
	writeMap := lfshook.WriterMap{
		logrus.DebugLevel: logWriter,
		logrus.InfoLevel:  logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.FatalLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{})
	Log.AddHook(lfHook)
}