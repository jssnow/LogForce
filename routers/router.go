package routers

import (
	"fmt"
	"gin_log/controllers"
	"gin_log/middleware"
	"github.com/gin-gonic/gin"
	"os"
	"path"
)

func init() {
	// 日志
	logFilePath := "./"
	logFileName := "log.log"

	// 日志文件
	fileName := path.Join(logFilePath, logFileName)

	// 写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	controllers.Log.Out = src
}

func APIRouter() *gin.Engine {
	r := gin.New()

	// 设置输出

	r.Use(middleware.Logger(controllers.Log))
	r.POST("/http", controllers.ReceiveLog)
	return r
}
