package common

import (
	"fmt"
	"github.com/spf13/viper"
)

// 全局配置变量
var Config = viper.New()

// 初始化配置读取
func InitConfig() {
	//添加读取的配置文件路径
	Config.AddConfigPath("./config/")
	Config.SetConfigName("config")
	Config.SetConfigType("yaml")
	if err := Config.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("配置读取失败 %s", err.Error()))
	}
}
