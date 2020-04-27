package entity

import (
	"gin_log/common"
	"strings"
	"sync"
)

// 业务日志处理接口
type BusinessLogInterface interface {
	Counter() bool
	SendNotice() bool
}

type BusinessLogCountResult struct {
	sync.RWMutex
	Data map[string]map[string]map[string]map[string]int
}

type BusinessLogHandler struct {
	Inputs      *LogContent
	CountResult *BusinessLogCountResult
}

// 统计业务日志数量
func (blh *BusinessLogHandler) Counter() bool {
	// 获取配置中需要统计的日志级别
	CountErrorLevel := common.Config.GetStringSlice(blh.Inputs.ProjectEnv + ".business_log_count_level")
	// 如果没有配置则统计所有级别的错误
	if len(CountErrorLevel) > 0 {
		// 只统计配置了的错误级别的错误
		if !common.IsStringExistsInSlice(blh.Inputs.Level, CountErrorLevel) {
			return true
		}
	}
	//按照项目->模块->级别统计日志数量
	if blh.Inputs.Project != "" && blh.Inputs.ModuleName != "" && blh.Inputs.Level != "" {
		//加锁
		blh.CountResult.Lock()
		if _, ok := blh.CountResult.Data[blh.Inputs.Project]; !ok {
			moduleMap := make(map[string]map[string]map[string]int)
			blh.CountResult.Data[blh.Inputs.Project] = moduleMap
		}
		if _, ok := blh.CountResult.Data[blh.Inputs.Project][blh.Inputs.ModuleName]; !ok {
			levelMap := make(map[string]map[string]int)
			blh.CountResult.Data[blh.Inputs.Project][blh.Inputs.ModuleName] = levelMap
		}
		if _, ok := blh.CountResult.Data[blh.Inputs.Project][blh.Inputs.ModuleName][blh.Inputs.ProjectEnv]; !ok {
			envMap := make(map[string]int)
			blh.CountResult.Data[blh.Inputs.Project][blh.Inputs.ModuleName][blh.Inputs.ProjectEnv] = envMap
		}

		blh.CountResult.Data[blh.Inputs.Project][blh.Inputs.ModuleName][blh.Inputs.ProjectEnv][blh.Inputs.Level]++
		blh.CountResult.Unlock()
	}
	return true
}

// 发送业务日志报警
func (blh *BusinessLogHandler) SendNotice() bool {
	// 获取配置中需要统计的日志级别
	WarnErrorLevel := common.Config.GetStringSlice(blh.Inputs.ProjectEnv + ".business_log_warn_level")
	ErrorNoticeSend := common.Config.GetStringSlice(blh.Inputs.ProjectEnv + strings.ToLower(blh.Inputs.Level) + "_level")
	//如果没有配置则统计所有级别的错误
	if len(WarnErrorLevel) > 0 {
		//对不同的错误级别发送不同的警报
		if common.IsStringExistsInSlice(blh.Inputs.Level, WarnErrorLevel) {
			if sendNotices, ok := ErrorNoticeSend[blh.Inputs.Level][Data.Project]; ok {
				// 取出对应模块的配置

				for k, v := range sendNotices {
					if k == "mail" {
						//发送邮件
						mail := common.MailContent{
							MailAddress: v,
							Content:     Data.Content,
						}
						common.MailAddress <- mail
					}

					if k == "dingding" {
						if limit.Allow() {
							dingContent := format.DingDingFormat(Data)
							ding := common.DingDingContent{
								DingDingToken: v,
								Content:       dingContent,
							}
							//发送钉钉
							common.DingDingSend <- ding
						} else {
							beego.Error("钉钉发送频率过快,已限制错误通知速度")
						}

					}
				}
			}
		}
	}
	return true
}
