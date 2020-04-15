package controller

import "github.com/gin-gonic/gin"

type JsonInput struct {
	//nginx日志
	BytesSent   string `json:"body_bytes_sent"`
	RequestUrl  string `json:"request"`
	Status      string `json:"status"`
	RequestTime string `json:"x_forwarded_for"`

	//业务日志
	Content    string   `json:"content"`
	Level      string   `json:"level"`
	HostName   string   `json:"host_name"`
	ModuleName string   `json:"module_name"`
	Tags       []string `json:"tags"`
	//日志所属项目:用户中心,大商城,
	Project string `json:"project"`
}

func ReceiveLog(c *gin.Context) {
	var jsonInputs []JsonInput
	if err := c.ShouldBindJSON(&jsonInputs); err != nil {
		c.JSON(500, gin.H{
			"message": "参数解析错误!",
		})
	}
	c.JSON(200, gin.H{
		"message": "pongsds",
	})
}
