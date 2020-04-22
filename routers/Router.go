package routers

import (
	"gin_log/common"
	"gin_log/controllers"
	"gin_log/middleware"
	"github.com/gin-gonic/gin"
)

func APIRouter() *gin.Engine {
	r := gin.New()

	// 设置输出

	r.Use(middleware.Logger(common.Log))
	r.POST("/http", controllers.ReceiveLog)
	return r
}
