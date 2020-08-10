package controllers

import (
	"LogForce/entity"
	"LogForce/services"
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

var jsonInput entity.LogContent

// 接收日志控制器
func TcpReceiveLog(logInfo []byte) {
	// logstash 推送tcp codec选项设置为json_lines则每条数据后都会加上一个\n，因此这里使用该标识来切分日志
	// output {
	//	tcp {
	//		host  => "192.168.0.2"
	//		port  => 8888
	//		codec => json_lines
	//	}
	// }
	logInfo = bytes.TrimSuffix(logInfo, []byte("\n"))
	// 判断是否无数据
	logInfoLen := len(logInfo)
	if logInfoLen == 0 {
		return
	}
	if bytes.Contains(logInfo, []byte("\n")) {
		// 处理多条一起发送过来的情况
		logs := bytes.Split(logInfo, []byte("\n"))
		for _, v := range logs {
			logInfoLen = len(v)
			if logInfoLen == 0 {
				continue
			}
			unmarshalLog(v, logInfoLen)
		}
	} else {
		unmarshalLog(logInfo, logInfoLen)
	}

	return
}

// 反解析日志，统计监控数据
func unmarshalLog(logInfo []byte, len int) {
	err := json.Unmarshal(logInfo, &jsonInput)
	if err != nil {
		log.Errorf("err:%v,log:%s", err, logInfo)
	}
	jsonInput.ContentLength = len
	services.DoDealLog(jsonInput)
}
