# runtime包

runtime指的是运行时。
- 记录了程序运行时的堆栈信息
- gc进行垃圾回收


```go

func main() {
	getInfo(1)
}

// pc - 函数调用的相关信息
// file - 文件名
// line - 行数
// ok - 如果能够取到相关信息，则ok=true, 否则ok=false
func getInfo(n int) {
	pc, file, line, ok := runtime.Caller(0) // n = 层数，调用的层数，main函数调用是0，getInfo调用是1, ...
	if !ok {
		fmt.Printf("runtime.Caller() failed\n")
		return
	}
	fmt.Println(pc)   // 4409405
	fmt.Println(file) // e:/workspace/go_works/go/worklog/test.go
	fmt.Println(line) // 9

    funcName := runtime.FuncForPC(pc).Name() // main.main

    fileName := path.Base(file)
	fmt.Println(fileName) // test.go

    funcName = strings.Split(funcName, ".")[1]
	fmt.Println(funcName) // main
}
```