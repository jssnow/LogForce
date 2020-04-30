package services

import (
	"gin_log/common"
	"gin_log/entity"
	"gopkg.in/gomail.v2"
	"strconv"
)

//发送邮件方法
func SendMail() {
	common.Log.Info("发送邮件协程启动")
	//获取发送邮件的配置
	mConfig := common.Config.GetStringMapString("email")
	if mConfig == nil {
		common.Log.Error("缺少邮件发送配置")
	}

	for {
		select {
		case mailString := <-entity.MailSendChan:
			if len(mailString.MailAddress) > 0 {
				m := gomail.NewMessage()
				m.SetHeader("From", mConfig["from"])
				for _, v := range mailString.MailAddress {
					m.SetHeader("To", v)
				}
				m.SetHeader("Subject", mConfig["subject"])
				m.SetBody("text/html", mailString.Content)
				mPort, _ := strconv.Atoi(mConfig["port"])
				d := gomail.NewDialer(mConfig["host"], mPort, mConfig["username"], mConfig["password"])

				// 发送
				if err := d.DialAndSend(m); err != nil {
					common.Log.Error(err)
				}
				common.Log.Info("报警邮件发送成功")
			}

		}
	}

	return
}
