package controllers

import (
	"gin_log/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var Log = log.New()

// 接收日志控制器
func ReceiveLog(c *gin.Context) {
	var jsonInputs []services.JsonInput
	if err := c.ShouldBindJSON(&jsonInputs); err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"message": "日志数据解析错误!",
		})
		return
	}

	// 判断是否无数据
	if len(jsonInputs) == 0 {
		c.JSON(200, gin.H{
			"message": "无日志!",
		})
		return
	}
	services.DealLogs(jsonInputs)
	Log.Info("qeqweqe")

	c.JSON(200, gin.H{
		"message": "成功!",
	})
	return
}
