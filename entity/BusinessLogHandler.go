package entity

import (
	"gin_log/services"
	"strings"
)

// 公共处理接口
type LogInterface interface {
	AnalysisTag() bool
}
type LogHandler struct {
	Inputs *services.JsonInput
}

// 处理日志的tag
func (lh LogHandler) AnalysisTag() bool {
	if len(lh.Inputs.Tags) < 1 {
		return false
	}

	for _, value := range lh.Inputs.Tags {
		if strings.Contains(value, "project_") {
			project := strings.Split(value, "_")
			if len(project) < 3 {
				break
			}
			lh.Inputs.Project = project[1]
			lh.Inputs.ProjectEnv = project[2]
		}

		// 业务日志
		if value == "business" {
			isBusiness = true
		}

		// nginx 日志
		if value == "nginx" {
			isNginx = true
		}
	}

	return true
}

// 业务日志处理接口
type BusinessLogInterface interface {
	Counter() bool
	SendNotice() bool
}

type BusinessLogHandler struct {
	Inputs *services.JsonInput
}

func (blh BusinessLogHandler) Counter() bool {
	// 获取配置中需要统计的日志级别

	//countLevel := common.Config.GetString("business_log_count.")
	//
	////如果没有配置则统计所有级别的错误
	//if len(CountErrorLevel) > 0 {
	//	//只统计配置了的错误级别的错误
	//	if !common.IsStringExistsInSlice(JsonInput.Level, CountErrorLevel) {
	//		return true
	//	}
	//}
	return true
}
func (blh BusinessLogHandler) SendNotice() bool {
	// 获取配置中需要统计的日志级别

	//countLevel := common.Config.GetString("business_log_count.")
	//
	////如果没有配置则统计所有级别的错误
	//if len(CountErrorLevel) > 0 {
	//	//只统计配置了的错误级别的错误
	//	if !common.IsStringExistsInSlice(JsonInput.Level, CountErrorLevel) {
	//		return true
	//	}
	//}
	return true
}
