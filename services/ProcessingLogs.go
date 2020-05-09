package services

import (
	"gin_log/common"
	"gin_log/entity"
)

var ErrorCountWithLock = entity.BusinessLogCountResult{
	Data: make(map[string]map[string]map[string]map[string]int),
}

var AnalysisResults = entity.InterfaceAnalysisMap{
	InterfaceAnalysisResult: make(map[string]entity.InterfaceAnalysis),
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
			businessHandle := entity.BusinessLogHandler{&v, &ErrorCountWithLock}
			businessHandle.Counter()
			businessHandle.SendNotice()
			common.Log.Info(businessHandle.CountResult.Data)
			break
		case "nginx":
			nginxHandle := entity.NginxLogHandler{&v, &AnalysisResults}
			nginxHandle.Count()
		}
	}


}
