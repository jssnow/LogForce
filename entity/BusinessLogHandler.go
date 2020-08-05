package entity

import (
	"encoding/json"
	"fmt"
	"LogForce/common"
	"golang.org/x/time/rate"
	"strings"
	"sync"
)

// 业务日志处理接口
type BusinessLogInterface interface {
	Count() bool
	SendNotice() bool
	DingDingFormat([]string) DingDingContent
	MailFormat([]string) MailContent
}

type BusinessLogCountResult struct {
	sync.RWMutex
	Data map[string]map[string]map[string]map[string]int
}

type BusinessLogHandler struct {
	Inputs      *LogContent
	CountResult *BusinessLogCountResult
}

type MailContent struct {
	// 需要发送该邮件的所有地址
	MailAddress []string
	// 邮件内容
	Content string
}

type DingDingContent struct {
	// 所有需要通知的钉钉
	DingDingToken []string
	// 钉钉通知内容
	Content []byte
}

// 发送邮件时使用的channel
var MailSendChan = make(chan MailContent, 200)

var DingDingSendChan = make(chan DingDingContent, 200)

// 钉钉通知速度限制
// 一分钟18个 防止触发钉钉限制,同时防止发送钉钉协程配置不够channel阻塞影响处理请求的能力
var dingLimit = rate.NewLimiter(0.3, 1)

// 统计业务日志数量
func (blh *BusinessLogHandler) Count() bool {
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
			blh.CountResult.Data[blh.Inputs.Project] = make(map[string]map[string]map[string]int)
		}
		if _, ok := blh.CountResult.Data[blh.Inputs.Project][blh.Inputs.ModuleName]; !ok {
			blh.CountResult.Data[blh.Inputs.Project][blh.Inputs.ModuleName] = make(map[string]map[string]int)
		}
		if _, ok := blh.CountResult.Data[blh.Inputs.Project][blh.Inputs.ModuleName][blh.Inputs.ProjectEnv]; !ok {
			blh.CountResult.Data[blh.Inputs.Project][blh.Inputs.ModuleName][blh.Inputs.ProjectEnv] = make(map[string]int)
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
	ErrorNoticeSend := common.Config.GetStringMapStringSlice(blh.Inputs.ProjectEnv + "." + strings.ToLower(blh.Inputs.Level) + "_notice" + "." + blh.Inputs.Project)
	//如果没有配置则统计所有级别的错误
	if len(WarnErrorLevel) > 0 {
		//对不同的错误级别发送不同的警报
		if common.IsStringExistsInSlice(blh.Inputs.Level, WarnErrorLevel) {
			fmt.Println(ErrorNoticeSend)
			if ErrorNoticeSend != nil {
				// 取出对应模块的配置
				for k, v := range ErrorNoticeSend {
					if k == "mail" {
						//发送邮件
						mail := blh.MailFormat(v)
						MailSendChan <- mail
					}

					if k == "dingding" {
						if dingLimit.Allow() {
							ding := blh.DingDingFormat(v)
							//发送钉钉
							DingDingSendChan <- ding
						} else {
							common.Log.Error("钉钉发送频率过快,已限制错误通知速度")
						}
					}
				}
			}
		}
	}
	return true
}

type DingFormat struct {
	MsgType  string            `json:"msgtype"`
	Markdown map[string]string `json:"markdown"`
	At       map[string]bool   `json:"at"`
}

// 格式化钉钉消息
func (blh *BusinessLogHandler) DingDingFormat(dingTokens []string) DingDingContent {
	//拼接消息字符串,截取中文字符串需要先转为[]rune类型,截取之后转为string
	contents := strings.Split(blh.Inputs.Content, "|")[0]
	runeContent := []rune(contents)
	var limitContent string
	if len(runeContent) > 1000 {
		limitContent = string(runeContent[:1000]) + "....更多内容请在邮件或者日志中查看"
	} else {
		limitContent = contents
	}
	text := fmt.Sprintf("### **服务器**  \n %s\n ### **错误级别**  \n %s\n ### **模块**  \n %s\n ### **错误内容** \n %s", blh.Inputs.HostName, blh.Inputs.Level, blh.Inputs.ModuleName, limitContent)
	var dingFormat = DingFormat{
		MsgType: "markdown",
		Markdown: map[string]string{
			"title": "业务日志系统错误报警",
		},
		At: map[string]bool{
			"isAtAll": false,
		},
	}
	dingFormat.Markdown["text"] = text
	dingJson, err := json.Marshal(dingFormat)
	if err != nil {
		common.Log.Error(err)
	}

	return DingDingContent{
		DingDingToken: dingTokens,
		Content:       dingJson,
	}
}

// 格式化邮件消息
func (blh *BusinessLogHandler) MailFormat(mails []string) MailContent {
	// TODO 格式化邮件报警
	return MailContent{
		MailAddress: mails,
		Content:     blh.Inputs.Content,
	}
}
