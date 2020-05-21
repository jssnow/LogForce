package entity

import (
	"gin_log/common"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 业务日志处理接口
type NginxLogInterface interface {
	Count() bool
}

type NginxLogHandler struct {
	Inputs      *LogContent
	CountResult *NginxAnalysisMap
}

type NginxAccessLogInfo struct {
	TimeLocal                    time.Time
	BytesSent                    int
	Path, Method, Scheme, Status string
	UpstreamTime, RequestTime    float64
}

//统计分析结果暂存内存中
type NginxAnalysis struct {
	//访问次数
	UrlCount map[string]int
	//平均请求时间
	UrlTimeAverage map[string]float64
	//最大请求时间
	UrlTimeMax map[string]float64
	//最小请求时间
	UrlTimeMin map[string]float64
}

//nginx分析结果带锁
type NginxAnalysisMap struct {
	NginxAnalysisResult map[string]map[string]NginxAnalysis
	sync.RWMutex
}

func (nlh *NginxLogHandler) Count() bool {
	var accessLog NginxAccessLogInfo
	/****************** 处理日志内容 ********************/
	//请求时间
	accessLog.RequestTime, _ = strconv.ParseFloat(strings.Trim(nlh.Inputs.RequestTime, "\""), 64)

	//处理url
	reqSli := strings.Split(nlh.Inputs.RequestUrl, " ")
	if len(reqSli) != 3 {
		common.Log.Info("nginx日志中接口不符合规则:", nlh.Inputs.RequestUrl)
		return true
	}
	//请求的类型
	accessLog.Method = reqSli[0]

	u, err := url.Parse(reqSli[1])
	if err != nil {
		common.Log.Info("url解析失败:", err)
		return true
	}
	//路径
	accessLog.Path = u.Path

	/****************** 统计分析结果 ********************/
	nlh.CountResult.Lock()
	if _, ok := nlh.CountResult.NginxAnalysisResult[nlh.Inputs.Project]; !ok {
		nlh.CountResult.NginxAnalysisResult[nlh.Inputs.Project] = make(map[string]NginxAnalysis)

	}

	if _, ok := nlh.CountResult.NginxAnalysisResult[nlh.Inputs.Project][nlh.Inputs.ProjectEnv]; !ok {
		nlh.CountResult.NginxAnalysisResult[nlh.Inputs.Project][nlh.Inputs.ProjectEnv] = NginxAnalysis{
			UrlCount:       make(map[string]int),
			UrlTimeAverage: make(map[string]float64),
			UrlTimeMax:     make(map[string]float64),
			UrlTimeMin:     make(map[string]float64),
		}
	}

	//记录接口请求的数量
	analysisResult := nlh.CountResult.NginxAnalysisResult[nlh.Inputs.Project][nlh.Inputs.ProjectEnv]
	analysisResult.UrlCount[accessLog.Path] += 1
	//记录接口的所求请求时间,持久化时用来计算平均响应时间
	analysisResult.UrlTimeAverage[accessLog.Path] += accessLog.RequestTime

	//记录接口的最大和最小的请求时间
	if accessLog.RequestTime > analysisResult.UrlTimeMax[accessLog.Path] {
		analysisResult.UrlTimeMax[accessLog.Path] = accessLog.RequestTime
	}
	if analysisResult.UrlTimeMin[accessLog.Path] == 0 {
		analysisResult.UrlTimeMin[accessLog.Path] = accessLog.RequestTime
	} else {
		if accessLog.RequestTime < analysisResult.UrlTimeMin[accessLog.Path] {
			analysisResult.UrlTimeMin[accessLog.Path] = accessLog.RequestTime
		}
	}

	//统计完成
	nlh.CountResult.NginxAnalysisResult[nlh.Inputs.Project][nlh.Inputs.ProjectEnv] = analysisResult
	common.Log.Info(nlh.CountResult.NginxAnalysisResult)
	nlh.CountResult.Unlock()
	return true
}
