package main

import (
	"E:/workspace/go_works/go/go_basic/logger"
)

func main() {
	log := logger.NewLog()
	for {
		id := 10010
		name := "墨香"
		log.Debug("这是一条Debug日志, id:%d, name:%s, err:%v", id, name, err)
		log.Info("这是一条Info日志")
		log.Warn("这是一条Warning日志")
		log.Error("这是一条Error日志, id:%d, name:%s, err:%v", err)
		time.Sleep(2 * time.Second)
	}
}
