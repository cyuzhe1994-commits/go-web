package go_web

import (
	"fmt"
	"log"
	"time"
)

type IFrameWorkLog interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
}

// 定义日志级别常量
const (
	DefaultFrameWorkLogLevelDebug = iota
	DefaultFrameWorkLogLevelInfo
	DefaultFrameWorkLogLevelWarn
	DefaultFrameWorkLogLevelError
)

type DefaultFrameWorkLog struct {
	*log.Logger
	level    int
	location *time.Location
}

func NewDefaultFrameWorkLog(level int, location *time.Location) IFrameWorkLog {
	if level < DefaultFrameWorkLogLevelDebug || level > DefaultFrameWorkLogLevelError {
		level = DefaultFrameWorkLogLevelInfo // 默认级别为 Info
	}
	if location == nil {
		location = time.UTC // 默认使用本地时区
	}
	return &DefaultFrameWorkLog{
		Logger:   log.Default(),
		level:    level,
		location: location,
	}
}

func (l *DefaultFrameWorkLog) Debug(format string, v ...interface{}) {
	l.log(DefaultFrameWorkLogLevelDebug, format, v...)
}

func (l *DefaultFrameWorkLog) Info(format string, v ...interface{}) {
	l.log(DefaultFrameWorkLogLevelInfo, format, v...)
}

func (l *DefaultFrameWorkLog) Warn(format string, v ...interface{}) {
	l.log(DefaultFrameWorkLogLevelWarn, format, v...)
}

func (l *DefaultFrameWorkLog) Error(format string, v ...interface{}) {
	l.log(DefaultFrameWorkLogLevelError, format, v...)
}

func (l *DefaultFrameWorkLog) log(level int, format string, v ...interface{}) {
	// 只有当前级别大于等于设置级别时才打印
	if level < l.level {
		return
	}

	prefix := "[INFO]"
	if level == DefaultFrameWorkLogLevelError {
		prefix = "[ERROR]"
	} else if level == DefaultFrameWorkLogLevelDebug {
		prefix = "[DEBUG]"
	} else if level == DefaultFrameWorkLogLevelWarn {
		prefix = "[WARN]"
	}

	fmt.Printf("%s %s - %s\n", time.Now().In(l.location).Format(time.RFC3339), prefix, fmt.Sprintf(format, v...))
}
