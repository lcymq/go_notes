package logger

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// 在终端写日志相关内容
type LogLv uint16

const (
	LOG_ERR LogLv = iota
	DEBUG
	TRACE
	INFO
	WARN
	ERROR
	FATAL
)

const (
	TIMEFORMAT = "2006-01-02 15:04:05"
)

func parseLogLevel(s string) LogLv {
	s = strings.ToLower(s)
	switch s {
	case "debug":
		return DEBUG
	case "trace":
		return TRACE
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return LOG_ERR
	}
}

// the struct of ConsoleLogger
type ConsoleLogger struct {
	Lv LogLv // 上线后只输出这个level以上的日志
}

func (c ConsoleLogger) enable(lv LogLv) bool {
	return c.Lv >= lv
}

// Newlog constructor
func NewLog(lvStr string) ConsoleLogger {
	level := parseLogLevel(lvStr)
	if level == LOG_ERR {
		panic("incorrect levec.")
	}
	return ConsoleLogger{
		Lv: level,
	}
}

func getInfo(n int) (string, string, int) {
	pc, file, line, ok := runtime.Caller(0) // n = 层数，调用的层数，main函数调用是0，getInfo调用是1, ...
	if !ok {
		fmt.Printf("runtime.Caller() failed\n")
		return "", "", 0
	}

	funcName := runtime.FuncForPC(pc).Name() // main.main
	return funcName, file, line
}

func log(lv string, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	now := time.Now()
	funcName, fileName, lineNo := getInfo(3)
	fmt.Printf("[%s] [DEBUG] [%s:%s:%d] %s\n", now.Format(TIMEFORMAT), fileName, funcName, lineNo, msg)
}

func (c ConsoleLogger) Debug(format string, a ...interface{}) {
	if c.enable(DEBUG) {
		log("DEBUG", format, a...)
	}
}

func (c ConsoleLogger) Info(format string, a ...interface{}) {
	if c.enable(INFO) {
		log("INFO", format, a...)
	}
}

func (c ConsoleLogger) Warn(format string, a ...interface{}) {
	if c.enable(WARN) {
		log("WARN", format, a...)
	}
}

func (c ConsoleLogger) Error(format string, a ...interface{}) {
	if c.enable(ERROR) {
		log("ERROR", format, a...)
	}

}

func (c ConsoleLogger) Fatal(format string, a ...interface{}) {
	if c.enable(FATAL) {
		log("FATAL", format, a...)
	}
}
