package controllers

import (
	"LogForce/common"
	"LogForce/entity"
	"LogForce/services"
	"github.com/gin-gonic/gin"
	"runtime"
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

	common.Log.Infof("成功收到并处理处理%d条日志", len(jsonInputs))
	//监控是否有goroutine泄露
	common.Log.Infof("the number of goroutines: %d", runtime.NumGoroutine())

	c.JSON(200, gin.H{
		"message": "成功!",
	})
	return
}
