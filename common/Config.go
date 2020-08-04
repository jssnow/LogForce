package common

import (
	"fmt"
	"github.com/spf13/viper"
)

// 全局配置变量
var Config = viper.New()

// 是否开启业务日志统计
var OpenBusinessCount = true
// 是否开启nginx日志统计
var OpenNginxCount = true

// 是否开启错误通知
var OpenNotice = true

// 初始化配置读取
func InitConfig() {
	//添加读取的配置文件路径
	Config.AddConfigPath("./config/")
	Config.SetConfigName("config")
	Config.SetConfigType("yaml")
	if err := Config.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("配置读取失败 %s", err.Error()))
	}

	// 业务日志统计
	OpenBusinessCount = Config.GetBool("open_business_count")

	// nginx日志统计
	OpenNginxCount = Config.GetBool("open_nginx_count")

	// 通知
	OpenNotice = Config.GetBool("open_notice")
}