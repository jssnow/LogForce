package entity

// 业务日志处理接口
type BusinessLogInterface interface {
	Counter() bool
	SendNotice() bool
}

type BusinessLogHandler struct {
	Inputs *LogContent
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
