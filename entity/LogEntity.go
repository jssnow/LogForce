package entity

import "strings"

type LogContent struct {
	// nginx 日志
	BytesSent   string `json:"body_bytes_sent"`
	RequestUrl  string `json:"request"`
	Status      string `json:"status"`
	RequestTime string `json:"x_forwarded_for"`

	// 业务日志
	Content    string   `json:"content"`
	Level      string   `json:"level"`
	HostName   string   `json:"host_name"`
	ModuleName string   `json:"module_name"`
	Tags       []string `json:"tags"`

	// 日志所属项目
	Project string
	// 日志所属项目的环境
	ProjectEnv string
	// 日志类型
	Type string
}

// 公共处理接口
type LogInterface interface {
	AnalysisTag() bool
}

// 日志公共处理
type LogHandler struct {
	Log *LogContent
}

// 处理日志的tag
func (lh *LogHandler) AnalysisTag() bool {
	if len(lh.Log.Tags) < 1 {
		return false
	}

	for _, value := range lh.Log.Tags {
		if strings.Contains(value, "project_") {
			project := strings.Split(value, "_")
			if len(project) < 3 {
				return false
			}
			lh.Log.Project = project[1]
			lh.Log.ProjectEnv = project[2]
		}

		// 日志类型
		if value == "business" || value == "nginx" {
			lh.Log.Type = value
		}
	}

	if lh.Log.Type == "" {
		return false
	}
	return true
}
