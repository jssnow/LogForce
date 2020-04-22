package services

import (
	"fmt"
	"gin_log/entity"
	"strings"
)

type JsonInput struct {
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
	Project string `json:"project"`
	// 日志所属项目的环境
	ProjectEnv string `json:"project_env"`
}

// 处理日志逻辑
func DealLogs(logs []JsonInput) {

	var isBusiness bool
	var isNginx bool
	for _, v := range logs {
		logHandle := entity.LogHandler{&v}
		res := logHandle.AnalysisTag()
		if !res {
			break
		}

		isBusiness = false
		isNginx = false
		// 根据日志的tag字段分析该条日志的类型和所属项目的信息
		if len(v.Tags) > 0 {
			for _, value := range v.Tags {
				if strings.Contains(value, "project_") {
					project := strings.Split(value, "_")
					if len(project) < 3 {
						break
					}
					v.Project = project[1]
					v.ProjectEnv = project[2]
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

			if isBusiness {
				// 分模块、分系统、分等级，记录异常报警数量，每十分钟（可配置）记录一次
				//filter.CountErrorsFilter(&v)
				// 发送邮件报警
				//warn.SendNotice(&v)
			} else if isNginx {
				// 分析nginx访问日志
				//filter.CountNginxFilter(&v)
			}
		}
		fmt.Println(v)
	}
}
