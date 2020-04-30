package services

import "gin_log/common"

//启动发送协程
func SendWarnService() {
	//默认启动10个
	numInt := 10
	maiNum := common.Config.GetInt("send_mail_goroutines_num")
	dingDingNum := common.Config.GetInt("send_dingding_goroutines_num")

	if maiNum > 0 {
		numInt = maiNum
	}
	for i := 0; i < numInt; i++ {
		go SendMail()
	}

	numInt = 10
	if dingDingNum > 0 {
		numInt = dingDingNum
	}
	for i := 0; i < numInt; i++ {
		go SendDingDing()
	}
}
