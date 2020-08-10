package services

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"runtime"
	"time"
)

const (
	B   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
	GiB // 1073741824
	TiB // 1099511627776 (exceeds 1 << 32)
	PiB // 1125899906842624
)

// 程序内部监控
func InitMonitor() {
	// 定时输出(1分钟)
	ticker := time.NewTicker(time.Second * 60)
	for {
		<-ticker.C
		// 携程数量
		log.Infof("the number of goroutines: %d", runtime.NumGoroutine())
		MonitorInfo.Lock()
		// 日志数量
		log.Infof("最近一分钟日志数: %d", MonitorInfo.Num)
		log.Infof("吞吐量: %d", MonitorInfo.Num/60)
		// 日志大小
		log.Infof("最近一分钟日志大小: %s", formatFileSize(MonitorInfo.ContentLength))
		MonitorInfo.Num = 0
		MonitorInfo.ContentLength = 0
		MonitorInfo.Unlock()
	}
}

// 字节的单位转换 保留两位小数
func formatFileSize(fileSize int) (size string) {
	if fileSize < KiB {
		return fmt.Sprintf("%dB", fileSize)
	} else if fileSize < MiB {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/KiB)
	} else if fileSize < GiB {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/MiB)
	} else if fileSize < TiB {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/GiB)
	} else if fileSize < PiB {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/TiB)
	} else {
		return fmt.Sprintf("%d B", fileSize)
	}
}
