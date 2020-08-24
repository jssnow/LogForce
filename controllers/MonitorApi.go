package controllers

import (
	"LogForce/services"
	"github.com/gin-gonic/gin"
)

func Monitor(c *gin.Context) {
	res := services.GetMonitorData()
	c.IndentedJSON(200, res)
	return
}
