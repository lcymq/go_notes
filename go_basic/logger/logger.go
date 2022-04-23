package logger

import "strings"

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
