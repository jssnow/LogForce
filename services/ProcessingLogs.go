package services

import (
	"LogForce/common"
	"LogForce/entity"
)

var LogCountWithLock = entity.BusinessLogCountResult{
	Data: make(map[string]map[string]map[string]map[string]int),
}

var NginxCountWithLock = entity.NginxAnalysisMap{
	NginxAnalysisResult: make(map[string]map[string]entity.NginxAnalysis),
}

var MonitorInfo = entity.Monitor{Num: 0, ContentLength: 0}

// 处理日志逻辑
func DealLogs(logs []entity.LogContent) {
	for _, v := range logs {
		DoDealLog(v)
	}
}

// 执行处理日志
func DoDealLog(singleLog entity.LogContent) {
	logHandle := entity.LogHandler{&singleLog, &MonitorInfo}
	// 监控系统吞吐量
	logHandle.DoMonitor()
	res := logHandle.AnalysisTag()
	if !res {
		return
	}

	switch singleLog.Type {
	case "business":
		businessHandle := entity.BusinessLogHandler{&singleLog, &LogCountWithLock}
		if common.OpenBusinessCount {
			businessHandle.Count()
		}
		if common.OpenNotice {
			businessHandle.SendNotice()
		}
		break
	case "nginx":
		if common.OpenNginxCount {
			nginxHandle := entity.NginxLogHandler{&singleLog, &NginxCountWithLock}
			nginxHandle.Count()
		}
		break
	}
}
