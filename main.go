package main

import (
	"gin_log/common"
	"gin_log/routers"
	"gin_log/services"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)



func main() {
	// 初始化配置
	common.InitConfig()

	// 设置程序运行模式
	gin.SetMode(common.Config.GetString("run_mode"))

	// 初始化日志
	common.InitLog()

	// 注册所有路由
	r := routers.APIRouter()

	// 启动发送邮件和钉钉的协程
	services.SendWarnService()

	// 连接数据库
	var err error
	// 添加前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "ts_" + defaultTableName;
	}
	common.Db, err = gorm.Open("mysql", common.GetDbConfigString())
	if err != nil {
		common.Log.Error(err)
	}
	defer common.Db.Close()
	// 禁止表名加s
	common.Db.SingularTable(true)

	// 获取监听端口
	port := common.Config.GetString("port")
	r.Run("127.0.0.1:" + port)
}