package services

import (
	"fmt"
	"gin_log/common"
	"gin_log/models"
	"time"
)

func Cron() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		<-ticker.C
		countErrorService()
	}

}

// 定时统计方法
func countErrorService() {

	//持久化业务日志统计结果
	if len(ErrorCountWithLock.Data) > 0 {
		var logs []models.BusinessLogErrorCount
		ErrorCountWithLock.Lock()
		for k, v := range ErrorCountWithLock.Data {
			//组装数据
			var businessLogErrorCount models.BusinessLogErrorCount
			businessLogErrorCount.App = k
			for kk, vv := range v {
				businessLogErrorCount.ModuleName = kk
				for kkk, vvv := range vv {
					businessLogErrorCount.Level = kkk
					businessLogErrorCount.Count = vvv
					logs = append(logs, businessLogErrorCount)
				}
			}
		}
		// 清空已有的数据重新统计
		ErrorCountWithLock.Data = make(map[string]map[string]map[string]map[string]int)
		// 记录统计结果的map操作完释放锁,写入mysql放后面,减少由于锁等待对日志处理性能的影响
		ErrorCountWithLock.Unlock()
		if len(logs) > 0 {
			//持久化分析结果
			successNums, err :=  common.Db.Create(logs)
			if err != nil {
				common.Log.Error(err)
			}

			common.Log.Info(fmt.Sprintf("成功写入%d条错误统计日志:%v", successNums, logs))
		} else {
			common.Log.Info("无错误日志统计数据")
		}

	} else {
		common.Log.Info("无业务错误日志统计数据")
	}
}
//
//	//持久化nginx日志分析结果
//	//if len(AnalysisResults.InterfaceAnalysisResult) > 0 {
//	//	o := orm.NewOrm()
//	//	var logs []models.LogAppAccess
//	//	AnalysisResults.Lock()
//	//	for k, v := range AnalysisResults.InterfaceAnalysisResult {
//	//		//组装数据
//	//		var appAccess models.LogAppAccess
//	//		appAccess.App = k
//	//		for url, count := range v.UrlCount {
//	//			urlTime := v.UrlTimeAverage[url]
//	//			urlTimeMin := v.UrlTimeMin[url] * 1000
//	//			urlTimeMax := v.UrlTimeMax[url] * 1000
//	//			if urlTime > 0 {
//	//				//计算平均时间
//	//				avgTime := (urlTime / float64(count)) * 1000
//	//				appAccess.Url = url
//	//				appAccess.AccessCount = count
//	//				appAccess.AccessAvgTime = avgTime
//	//				appAccess.AccessMinTime = urlTimeMin
//	//				appAccess.AccessMaxTime = urlTimeMax
//	//				logs = append(logs, appAccess)
//	//			}
//	//		}
//	//	}
//	//	//清空已有的数据重新统计
//	//	AnalysisResults.InterfaceAnalysisResult = make(map[string]InterfaceAnalysis)
//	//	AnalysisResults.Unlock()
//	//	if len(logs) > 0 {
//	//		//持久化分析结果
//	//		successNums, err := o.InsertMulti(100, logs)
//	//		if err != nil {
//	//			beego.Error(err)
//	//		}
//	//
//	//		beego.Info(fmt.Sprintf("成功写入%d条nginx access log统计日志:%v", successNums, logs))
//	//	} else {
//	//		beego.Info("无nginx access log日志统计可写入数据")
//	//	}
//	//} else {
//	//	beego.Info("无nginx access log日志统计数据")
//	//}
//
//	return nil
//}
