package common

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

func InitLog() {
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

		fileName := path.Join(logFilePath, logFileName)

		// 写入文件
		src, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeAppend|os.ModePerm)
		if err != nil {
			fmt.Println("err", err)
		}
		log.SetOutput(src)
		break
	case "console":
		// 输出控制台
		log.SetOutput(os.Stdout)
		break
	default:
		log.SetOutput(os.Stdout)
		break
	}
}
