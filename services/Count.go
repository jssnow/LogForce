package services

import (
	"LogForce/common"
	"LogForce/entity"
	"LogForce/models"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

func Cron() {
	timeInterval := common.Config.GetDuration("count_result_time_interval")
	ticker := time.NewTicker(time.Second * timeInterval)
	for {
		<-ticker.C
		countErrorService()
	}
}

// 定时统计方法
func countErrorService() {
	//持久化业务日志统计结果
	if len(LogCountWithLock.Data) > 0 {
		var logs []models.BusinessLogErrorCount
		LogCountWithLock.Lock()
		for k, v := range LogCountWithLock.Data {
			//组装数据
			var businessLogErrorCount models.BusinessLogErrorCount
			businessLogErrorCount.App = k
			for kk, vv := range v {
				businessLogErrorCount.ModuleName = kk
				for kkk, vvv := range vv {
					businessLogErrorCount.Env = kkk
					for level, count := range vvv {
						businessLogErrorCount.Level = level
						businessLogErrorCount.Count = count
					}
					logs = append(logs, businessLogErrorCount)
				}
			}
		}
		// 清空已有的数据重新统计
		LogCountWithLock.Data = make(map[string]map[string]map[string]map[string]int)
		// 记录统计结果的map操作完释放锁,写入mysql放后面,减少由于锁等待对日志处理性能的影响
		LogCountWithLock.Unlock()
		if len(logs) > 0 {
			// TODO gorm包2.0版本会支持批量插入功能,此处修改为批量插入
			for _, v := range logs {
				//持久化分析结果
				err := common.Db.Create(&v).Error
				if err != nil {
					log.Error(err)
				}
			}
			log.Info(fmt.Sprintf("成功写入错误统计日志:%v", logs))
		} else {
			log.Info("无错误日志统计数据")
		}

	} else {
		log.Info("无业务错误日志统计数据")
	}

	//持久化nginx日志分析结果
	if len(NginxCountWithLock.NginxAnalysisResult) > 0 {
		var logs []models.LogAppAccess
		NginxCountWithLock.Lock()
		for k, v := range NginxCountWithLock.NginxAnalysisResult {
			//组装数据
			var appAccess models.LogAppAccess
			appAccess.App = k
			for kk, vv := range v {
				appAccess.Env = kk
				for url, count := range vv.UrlCount {
					urlTime := vv.UrlTimeAverage[url]
					urlTimeMin := vv.UrlTimeMin[url] * 1000
					urlTimeMax := vv.UrlTimeMax[url] * 1000
					if urlTime > 0 {
						//计算平均时间
						avgTime := (urlTime / float64(count)) * 1000
						appAccess.Url = url
						appAccess.AccessCount = count
						appAccess.AccessAvgTime = avgTime
						appAccess.AccessMinTime = urlTimeMin
						appAccess.AccessMaxTime = urlTimeMax
						logs = append(logs, appAccess)
					}
				}
			}
		}
		//清空已有的数据重新统计
		NginxCountWithLock.NginxAnalysisResult = make(map[string]map[string]entity.NginxAnalysis)
		NginxCountWithLock.Unlock()
		if len(logs) > 0 {
			for _, v := range logs {
				//持久化分析结果
				err := common.Db.Create(&v).Error
				if err != nil {
					log.Error(err)
				}
			}
			log.Info(fmt.Sprintf("成功写入nginx access log统计日志:%v", logs))
		} else {
			log.Info("无nginx access log日志统计可写入数据")
		}
	} else {
		log.Info("无nginx access log日志统计数据")
	}

	return

}
