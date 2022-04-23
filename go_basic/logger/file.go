package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

// 日志写入文件中

type FileLogger struct {
    Level LogLevel
    FilePath string // 文件保存的路径
    FileName string // 文件保存的文件名
    maxFileSize int64
    fileObj *os.File
    errFileObj *os.File
}

// 
func NewFileLogger(levelStr, filePath, fileName string, maxSize uint64) (*FileLogger) {
    logLevel, err := parseLogLevel(levelStr)
    if err != nil {
        panic(err)
    }
    fileLog := &FileLogger{
        Level: logLevel,
        FilePath: filePath,
        FileName: fileName,
        maxFileSize: int64(maxSize),
    }
    fileLog.initFile() // 按照文件路径和文件名将文件打开
    return fileLog
}

func (f *FileLogger) initFile() {
    fullFileName := path.Join(f.FilePath, f.FileName)
    fileObj, err := os.OpenFile(fullFileName, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
    if err != nil {
        fmt.Printf("open log file failed, err %v\n", err)
    }
    errFileObj, err := os.OpenFile(fullFileName+".err", os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
    if err != nil {
        fmt.Printf("open log file failed, err %v\n", err)
    }
    // 日志文件都打开后
    f.fileObj = fileObj
    f.errFileObj = errFileObj
}

func (f *FileLogger) Close() {
    f.fileObj.Close()
    f.errFileObj.Close()
}

func (f *FileLogger) enable(lv LogLv) bool {
	return f.Level >= lv
}

// Newlog constructor
func NewLog(lvStr string) FileLogger {
	level := parseLogLevel(lvStr)
	if level == LOG_ERR {
		panic("incorrect level.")
	}
	return FileLogger{
		Lv: level,
	}
}

func getInfo(n int) (string, string, int) {
	pc, file, line, ok := runtime.Caller(0) // n = 层数，调用的层数，main函数调用是0，getInfo调用是1, ...
	if !ok {
		fmt.Fprintf("runtime.Caller() failed\n")
		return "", "", 0
	}

	funcName := runtime.FuncForPC(pc).Name() // main.main
	return funcName, file, line
}

func log(lv string, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	now := time.Now()
	funcName, fileName, lineNo := getInfo(3)
	fmt.Fprintf(f.fileObj, "[%s] [DEBUG] [%s:%s:%d] %s\n", now.Format(TIMEFORMAT), fileName, funcName, lineNo, msg)
    if lv > ERROR { // 如果要记录的日志大于等于ERROR级别，我还要在err日志文件中在记录一遍
        fmt.Fprintf(f.errFileObj, "[%s] [DEBUG] [%s:%s:%d] %s\n", now.Format(TIMEFORMAT), fileName, funcName, lineNo, msg)
    }
}

func (f *FileLogger) Debug(format string, a ...interface{}) {
	if f.enable(DEBUG) {
		log("DEBUG", format, a...)
	}
}

func (f *FileLogger) Info(format string, a ...interface{}) {
	if f.enable(INFO) {
		log("INFO", format, a...)
	}

}

func (f *FileLogger) Warn(format string, a ...interface{}) {
	if f.enable(WARN) {
		log("WARN", format, a...)
	}
}

func (f *FileLogger) Error(format string, a ...interface{}) {
	if f.enable(ERROR) {
		log("ERROR", format, a...)
	}

}

func (f *FileLogger) Fatal(format string, a ...interface{}) {
	if f.enable(FATAL) {
		log("FATAL", format, a...)
	}
}
