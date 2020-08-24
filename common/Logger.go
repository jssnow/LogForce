package common

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

// 全局日志变量
//var Log = log.New()

// 初始化配置读取
func InitLog() {
	//logInstance := log.New()
	// 获取日志输出方式
	t := Config.GetString("log.out_type")
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     false,
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	switch t {
	case "file":
		logFilePath := "./runtime/log"
		logFileName := Config.GetString("log.file_name")

		err := os.MkdirAll(logFilePath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}

		// 日志文件
		fileName := path.Join(logFilePath, logFileName)

		// 写入文件
		src, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE, os.ModePerm)
		if err != nil {
			fmt.Println("err", err)
		}
		//logInstance.Out = src
		log.SetOutput(src)
		break
	case "console":
		// 输出控制台
		log.SetOutput(os.Stdout)
		//logInstance.Out = os.Stdout
		break
	default:
		//logInstance.Out = os.Stdout
		log.SetOutput(os.Stdout)
		break
	}
}
