package main

import (
	"LogForce/common"
	"LogForce/routers"
	"LogForce/services"
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


	initNotice()
	initCount()

	// 获取监听端口
	port := common.Config.GetString("port")
	r.Run("0.0.0.0:" + port)
}

// 初始化错误通知
func initNotice() {
	if common.OpenNotice {
		// 启动发送邮件和钉钉的协程
		services.SendWarnService()
	}
}

// 初始化日志统计
func initCount() {
	if common.OpenBusinessCount || common.OpenNginxCount {
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
		// 启动定时刷新统计结果到mysql协程
		go services.Cron()
	}
}
