package routers

import (
	"LogForce/common"
	"LogForce/controllers"
	"LogForce/middleware"
	"github.com/gin-gonic/gin"
)

func APIRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logger(common.Log))
	r.POST("/http", controllers.ReceiveLog)
	return r
}
