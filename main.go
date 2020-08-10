package main

import (
	"LogForce/common"
	"LogForce/controllers"
	"LogForce/routers"
	"LogForce/services"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/evio"
)

func main() {
	// 初始化配置
	common.InitConfig()

	// 初始化日志
	common.InitLog()

	initNotice()
	initCount()
	go services.InitMonitor()

	mode := common.Config.Get("data_receiving_mode")
	log.Infof("数据接收方式%s", mode)

	switch mode {
	case "tcp":
		var events evio.Events
		events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {

			controllers.TcpReceiveLog(in)
			return
		}
		tcpPort := common.Config.GetString("tcp.tcp_port")
		log.Infof("tcp port %s", tcpPort)
		if err := evio.Serve(events, "tcp://0.0.0.0:"+tcpPort); err != nil {
			panic(err.Error())
		}
		break
	case "http":
		// 设置程序运行模式
		gin.SetMode(common.Config.GetString("run_mode"))

		// 注册所有路由
		r := routers.APIRouter()

		// 获取监听端口
		port := common.Config.GetString("http_port")
		r.Run("0.0.0.0:" + port)
		break
	default:
		log.Error("数据接收方式必须为tcp或者http")
	}
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
			return "ts_" + defaultTableName
		}
		common.Db, err = gorm.Open("mysql", common.GetDbConfigString())
		if err != nil {
			log.Error(err)
		}
		defer common.Db.Close()
		// 禁止表名加s
		common.Db.SingularTable(true)
		// 启动定时刷新统计结果到mysql协程
		go services.Cron()
	}
}
