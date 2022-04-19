package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

type MysqlConfig struct {
	Address  string `ini:"address"`
	Port     int    `ini:"port"`
	Username string `ini:"username"`
	Password string `ini:"password"`
}

type RedisConfig struct {
	Host string `ini:"host"`
	Port string `ini:"port"`
}

func loadIni(fileName string, data interface{}) (err error) {
	// 0. 参数的校验
	// 0.1 传进来的data参数必须是指针类型(因为需要在函数中对其进行赋值)
	t := reflect.TypeOf(data) // *main.MysqlConfig
	if t.Kind() != reflect.Ptr {
		err = errors.New("data param should be a pointer") // 格式化输出之后返回一个error类型
		return
	}
	// 0.2 传进来的data参数必须是结构体类型指针(因为配置文件中各种键值对需要赋值给结构体的字段)
	if t.Elem().Kind() != reflect.Struct {
		err = errors.New("data param should be a struct pointer")
	}
	// 1. 读文件得到字节类型数据
	fileObj, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}
	lineSlice := strings.Split(string(fileObj), "\r\n")
	fmt.Println(lineSlice)
	// 2. 一行一行的读取数据
	for idx, line := range lineSlice {
		// 2.0 去掉字符串首尾空格
		line = strings.TrimSpace(line)
		// 2.1 如果是注释就跳过
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		// 2.2 如果是[开头的就表示是节(section)
		if line[0] != '[' || line[len(line)-1] != ']' {
			err = fmt.Errorf("line:%d syntax error", idx+1)
			return
		}
		// 2.3 如果不是[开头就是=分割的键值对
	}

	return
}

func main() {
	var mc MysqlConfig
	err := loadIni("./conf.ini", &mc)
	if err != nil {
		fmt.Printf("load ini failed, err:%v\n", err)
		return
	}
	// fmt.Println(mc)
}
