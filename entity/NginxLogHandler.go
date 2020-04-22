package entity

// 业务日志处理接口
type NginxLogInterface interface {
	CountLog() bool
}

type NginxLogHandler struct {
	Inputs *LogContent
}

func (nlh NginxLogHandler) CountLog() bool {
	return true
}
