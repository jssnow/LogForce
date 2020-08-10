package entity

import (
	"strings"
	"sync"
)

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
	// 该条日志的长度
	ContentLength int
}

// 公共处理接口
type LogInterface interface {
	AnalysisTag() bool
}

// 日志公共处理
type LogHandler struct {
	Log     *LogContent
	Monitor *Monitor
}

type Monitor struct {
	sync.RWMutex
	Num           int
	ContentLength int
}

// 处理日志的tag
func (lh *LogHandler) AnalysisTag() bool {
	if len(lh.Log.Tags) < 1 {
		return false
	}

	for _, value := range lh.Log.Tags {
		if strings.Contains(value, "project_") {
			project := strings.Split(value, "_")
			len := len(project)
			if len < 2 {
				return false
			}
			lh.Log.Project = project[1]
			if len == 3 {
				lh.Log.ProjectEnv = project[2]
			} else {
				lh.Log.ProjectEnv = "default"
			}
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

// 监控
func (lh *LogHandler) DoMonitor() bool {
	// 监控数据
	lh.Monitor.Lock()
	lh.Monitor.Num += 1
	lh.Monitor.ContentLength += lh.Log.ContentLength
	lh.Monitor.Unlock()
	return true
}
