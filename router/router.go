package router

import (
	"gin_log/controller"
	"github.com/gin-gonic/gin"
)

func APIRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/http", controller.ReceiveLog)
	return r
}
