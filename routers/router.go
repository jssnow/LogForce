package routers

import (
	"gin_log/controllers"
	"github.com/gin-gonic/gin"
)

func APIRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/http", controllers.ReceiveLog)
	return r
}
