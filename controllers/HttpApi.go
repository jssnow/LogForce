package controllers

import (
	"gin_log/entity"
	"gin_log/services"
	"github.com/gin-gonic/gin"
)

// 接收日志控制器
func ReceiveLog(c *gin.Context) {
	var jsonInputs []entity.LogContent
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

	c.JSON(200, gin.H{
		"message": "成功!",
	})
	return
}
