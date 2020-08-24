package routers

import (
	"LogForce/controllers"
	"LogForce/middleware"
	"github.com/gin-gonic/gin"
)

// http receive log api
func APIRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logger())
	r.POST("/http", controllers.ReceiveLog)
	return r
}

// system monitoring api
func MonitorRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logger())
	r.GET("/monitor", controllers.Monitor)
	return r
}
