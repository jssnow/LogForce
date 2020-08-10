package services

import (
	"LogForce/entity"
	"bytes"
	log "github.com/sirupsen/logrus"
	"net/http"
)

//发送钉钉
func SendDingDing() {
	log.Info("发送钉钉协程启动")
	for {
		select {
		case dingString := <-entity.DingDingSendChan:
			if len(dingString.DingDingToken) > 0 {
				for _, value := range dingString.DingDingToken {
					dingDingURL := "https://oapi.dingtalk.com/robot/send?access_token=" + value
					//发送消息到钉钉群使用webhook
					resp, err := http.Post(dingDingURL, "application/json", bytes.NewBuffer(dingString.Content))
					//body, err := ioutil.ReadAll(resp.Body)
					if err != nil || resp.StatusCode != 200 {
						log.Error("报警钉钉发送请求失败:")
					} else {
						log.Info("报警钉钉发送请求成功")
					}
					resp.Body.Close()
				}
			}
		}
	}

	return
}
