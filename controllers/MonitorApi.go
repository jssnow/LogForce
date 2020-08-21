package controllers

import (
	"LogForce/services"
	"github.com/gin-gonic/gin"
)

// 接收日志控制器
func Monitor(c *gin.Context) {
	num, qps, size := services.GetMonitorData()
	c.JSON(200, gin.H{
		"最近一分钟日志数":   num,
		"吞吐量:":       qps,
		"最近一分钟日志大小:": size,
	})
	return
}
