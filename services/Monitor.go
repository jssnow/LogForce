package services

import (
	"fmt"
	"runtime"
	"sync"
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

type MonitorOut struct {
	RecordTime    time.Time
	Num           int
	ContentLength string
	Goroutines    int
}

type MonitorOutMap struct {
	MonitorOutResult []MonitorOut
	sync.RWMutex
}

var MonitorOutMapWithLock = MonitorOutMap{
	MonitorOutResult: []MonitorOut{},
}

// Initialize monitoring timing task
func InitMonitorCron() {
	ticker := time.NewTicker(time.Minute * 5)
	for {
		<-ticker.C

		MonitorInfo.Lock()
		out := MonitorOut{time.Now(), MonitorInfo.Num, FormatFileSize(MonitorInfo.ContentLength), runtime.NumGoroutine()}
		MonitorInfo.ContentLength = 0
		MonitorInfo.Num = 0
		MonitorInfo.Unlock()

		MonitorOutMapWithLock.Lock()
		MonitorOutMapWithLock.MonitorOutResult = append(MonitorOutMapWithLock.MonitorOutResult, out)
		MonitorOutMapWithLock.Unlock()
	}
}

func GetMonitorData() []MonitorOut {
	MonitorOutMapWithLock.Lock()
	result := MonitorOutMapWithLock.MonitorOutResult
	MonitorOutMapWithLock.Unlock()
	return result
}

// two decimal places
func FormatFileSize(fileSize int) (size string) {
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
