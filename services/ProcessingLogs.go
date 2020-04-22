package services

import (
	"fmt"
	"gin_log/entity"
)

// 处理日志逻辑
func DealLogs(logs []entity.LogContent) {
	for _, v := range logs {
		logHandle := entity.LogHandler{&v}
		res := logHandle.AnalysisTag()
		if !res {
			break
		}

		switch v.Type {
		case "business":
			businessHandle := entity.BusinessLogHandler{&v}
			businessHandle.Counter()
			businessHandle.SendNotice()
			break
		case "nginx":
			nginxHandle := entity.NginxLogHandler{&v}
			nginxHandle.CountLog()
		}
		fmt.Println(v)
	}
}
