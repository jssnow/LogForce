package services

import (
	"gin_log/common"
	"gin_log/entity"
)

var LogCountWithLock = entity.BusinessLogCountResult{
	Data: make(map[string]map[string]map[string]map[string]int),
}

var NginxCountWithLock = entity.NginxAnalysisMap{
	NginxAnalysisResult: make(map[string]map[string]entity.NginxAnalysis),
}

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
			businessHandle := entity.BusinessLogHandler{&v, &LogCountWithLock}
			businessHandle.Counter()
			businessHandle.SendNotice()
			common.Log.Info(businessHandle.CountResult.Data)
			break
		case "nginx":
			nginxHandle := entity.NginxLogHandler{&v, &NginxCountWithLock}
			nginxHandle.Count()
			break
		}
	}

}
