# Channel

单纯地将函数并发执行是没有意义的。函数与函数间需要交换数据才能体现并发执行函数的意义。

虽然可以使用共享内存进行数据交换，但是共享内存在不同的goroutine中容易发生竞态问题。为了保证数据交换的正确性，必须使用互斥量堆内存进行枷锁，这种做法势必造成性能问题。


Go语言的并发模型是`CSP(Communicating Sequential Procsses)`，提倡**通过通信共享内存**而不是**通过共享内存而实现通信**。

如果说`goroutine`是Go程序并发的执行体，`channel`就是它们之间的连接。`channel`是可以让一个`goroutine`发送特定值到另一个`goroutine`的通信机制。

Go语言中的通道`channel`是一种特殊类型。通道像一个传送带或者队列，总是遵循先入先出(First In First Out)的规则，保证收发数据的顺序，每一个通道都是一个具体类型的导管，也就是声明channel的时候需要为其指定元素类型。

### channel类型
声明通道类型的格式如下：
```go
var 变量 chan 元素类型
```
```go
var ch1 chan int //声明一个传递整型的通道
var ch2 chan bool //声明一个传递布尔型的通道
var ch3 chan []int //声明一个传递int切片的通道
```

### chan是指针类型，初始化后才能使用
用`make`关键字初始化`chan`
```go
var a chan int
var b chan int

func main() {
	fmt.Println(b)        // nil
	a = make(chan int)    // 不带缓冲区的通道初始化
	b = make(chan int, 16) // 带缓冲区的通道初始化
    fmt.Println(b)
}
```
不带缓冲区的通道不能接收数据：
```go
func main() {
	b := make(chan int)
	b <- 10
	b = make(chan int, 16)
	fmt.Println(b)
}
```
报错
```
fatal error: all goroutines are asleep - deadlock!
goroutine 1 [chan send]:
main.main()
	d:/workspace/test/src/go/test/main.go:45 +0x65
```
将`b <- 10`另起一个goroutine去接收值，可以解决这个错误：
```go
func main() {
	b := make(chan int)
	go func(){
		b <- 10
	}()
	// b = make(chan int, 16)
	fmt.Println(b)
}
```
此时输出：
```
0xc000018240
```
或者启一个goroutine去接收传进来的值：
```go
func main() {
	b := make(chan int)
	go func(){
		<- b
	}()
	b <- 10
	// b = make(chan int, 16)
	fmt.Println(b)
}
```
此时输出：
```
0xc0001020c0
```
检查通道b是否真的接收到传入的值：
```go
func main() {
	b := make(chan int)
	go func(){
		x := <- b
		fmt.Println("后台goroutine从通道b中取出了:", x)
	}()
	b <- 10
	fmt.Println("10发送到通道b中了...")
	// b = make(chan int, 16)
	fmt.Println(b)
}
```
输出：
```
10发送到通道b中了...
后台goroutine从通道b中取出了: 10
0xc000018240
```

###### 带有一个缓冲区的通道
一次只能存一个数据。
以下代码妄图存两个`int`进chan，做梦！
```go
func main() {
	b := make(chan int, 1)
	b <- 10
	fmt.Println("10发送到通道b中了...")
	b <- 20
	fmt.Println("20发送到通道b中了...")
	x := <- b
	fmt.Println("从通道b中取出了:", x)
	fmt.Println(b)
	close(b)
}
```
报错信息：
```
10发送到通道b中了...
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
	d:/workspace/test/src/go/test/main.go:48 +0x185
```

想要正常运行，缓冲区申请大一点:
```go
func main() {
	b := make(chan int, 10)
	b <- 10
	fmt.Println("10发送到通道b中了...")
	b <- 20
	fmt.Println("20发送到通道b中了...")
	x := <- b
	fmt.Println("从通道b中取出了:", x)
	fmt.Println(b)
	close(b)
}
```
输出：
```
10发送到通道b中了...
20发送到通道b中了...
从通道b中取出了: 10
0xc00013e000
```
### 发送
将一个值发送到通道中
```go
ch <- 10
```

### 接收
从一个通道中接收值
```go
x := <-ch // 从ch中接收值并赋值给变量x
<-ch      //从ch中接收值，忽略结果
```

### 关闭
通过内置的`close`函数来关闭通道
```go
close(ch)
```
注意：
- 只有在通知接收方goroutine所有的数据都发送完毕的时候才需要关闭通道。通道是可以被垃圾回收机制回收的。
- 它和关闭文件不一样，在结束操作之后关闭文件是必须要做的，但关闭通道不是必须的。
- 通道不关闭可能会死锁！`fatal error: all goroutines are asleep - deadlock!`

例：
```go
package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup
var once sync.Once

func f1(ch1 chan int) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		ch1 <- i
	}
	close(ch1)
}

func f2(ch1, ch2 chan int) {
	defer wg.Done()
	for {
		x, ok := <-ch1
		if !ok {
			break
		}
		ch2 <- x * x
	}
	once.Do(func() { close(ch2) }) //确保某个操作只执行一次
}

func main() {
	a := make(chan int, 100)
	b := make(chan int, 100)
	wg.Add(3)
	go f1(a)
	go f2(a, b)
	go f2(a, b)
	wg.Wait()
	for ret := range b {
		fmt.Println(ret)
	}
}
```

### chan的长度
```go
len(ch)
```

### 从chan取值，当ok为false时...
```go
x, ok := <-ch1
```
当`ok`为`false`时，`x`为相应数据类型的默认值。
### 单向通道
多用在函数的参数里

- 只能接收值的通道
```go
func f1(ch1 chan<- int) {
    defer wg.Done()
    for i:=0;i<100; i++ {
        ch1 <- i
    }
    // <-ch1 // 会报错
    close(ch1)
}
```
- 只能取值的通道
```go
func f2(ch1 <-chan int, ch2 chan<- int) {
	defer wg.Done()
	for {
		x, ok := <-ch1
		if !ok {
			break
		}
		ch2 <- x * x
	}
	once.Do(func() { close(ch2) }) //确保某个操作只执行一次
}
```

### 通道总结
`channel`常见的异常总结

| channel | nil   | 非空                         | 空的               | 满了                         | 没满                         |
| ------- | ----- | ---------------------------- | ------------------ | ---------------------------- | ---------------------------- |
| 接收    | 阻塞  | 接收值                       | 阻塞               | 接收值                       | 接收值                       |
| 发送    | 阻塞  | 发送值                       | 发送值             | 阻塞                         | 发送值                       |
| 关闭    | panic | 关闭成功，读完数据后返回零值 | 关闭成功，返回零值 | 关闭成功，读完数据后返回零值 | 关闭成功，读完数据后返回零值 |

关闭已经关闭的`channel`也会引发`panic`。



### worker pool (goroutine 池)

编写代码实现一个计算随机数的每个位置数字之和的程序，要求使用`goroutine`和`channel`构建生产者和消费者模式，可以指定启动的goroutine数量 - `worker pool`模式。

在工作中我们通常会使用`workerpool`模式，控制`goroutine`的数量，防止`goroutine`泄露和暴涨。

```go
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
使用goroutine和channel实现一个计算int64随机数个位数和的程序
1. 开启一个goroutine循环生成int64类型的随机数，发送到jobChan
2. 开启24个goroutine从jobChan中取出随机数计算各位数的和，将结果发送到resultChan
3. 主goroutine从resultChan取出结果并打印到终端输出
*/

var wg sync.WaitGroup

type Job struct {
	Value int64
}

type Result struct {
	Job *Job
	Sum int64
}

var jobChan = make(chan *Job, 100)
var resultChan = make(chan *Result, 100)

func producer(job chan<- *Job) {
	// 循环生成int64类型的随机数，发送到jobChan
	rand.Seed(time.Now().Unix())
	for {
		x := rand.Int63()
		newJob := &Job{
			Value: x,
		}
		job <- newJob
		time.Sleep(time.Millisecond * 500)
	}
}

func consumer(job <-chan *Job, result chan<- *Result) {
	// 从jobChan中取出随机数计算各位数的和，将结果发送到resultChan
	for {
		j := <-job
		var sum int64
		n := j.Value
		for n > 0 {
			sum += n % 10
			n = n / 10
		}
		newResult := &Result{
			Job: j,
			Sum: sum,
		}
		resultChan <- newResult
	}
}

func main() {
	wg.Add(1)
	go producer(jobChan)
	for i := 0; i < 24; i++ {
		go consumer(jobChan, resultChan)
	}
	for result := range resultChan {
		fmt.Printf("value:%d sum:%d\n", result.Job.Value, result.Sum)
	}
	wg.Wait()
}
```

### 什么是竞态条件？
- 当多个线程尝试访问和修改相同的数据(内存地址)时，就会出现竞态条件。
- e.g.当一个线程试图增加一个整数而另一个线程试图读取它。
- 如果变量是只读的，就不会有竞态条件。


### worker pool (goroutine池)
- 需求：编写代码实现一个计算随机数的每个位置数字之和的程序，要求使用`goroutine`和`channel`构建生产者和消费者模式，可以指定启动的goroutine数量 - `worker pool`模式。
- 在工作中，通常需要使用`worker pool`模式，控制`goroutine`的数量，防止`goroutine`泄露和暴涨。
```go
func worker(id int, jobs <-chan int, result chan<- int) {
	for job := range jobs {
		fmt.Printf("worker %d start job %d\n", id, job)
		time.Sleep(time.Second)
		fmt.Printf("worker %d end job %d\n", id, job)
		result <- job * 2
	}
}

func main() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	// 开启3个goroutine
	for w:=1; w<=3; w++ {
		go worker(w, jobs, results)
	}
	// 12个任务
	for j:=1; j<=12; j++ {
		jobs <- j
	}
	close(jobs)
	// 输出结果
	for result:=1; result<=12; result++ {
		<-results
	}
	close(results)
}
```


- 使用goroutine和channel实现一个计算int64随机数个位数和的程序
	1. 开启一个goroutine循环生成int64类型的随机数，发送到jobChan
	2. 开启24个goroutine从jobChan中取出随机数计算各位数的和，将结果发送到resultChan
	3. 主goroutine从resultChan取出结果并打印到终端输出
```go
// 获取各位上的数
func splitNum(n int64) []int64 {
	res := []int64{}
	for n > 0 {
		res = append(res, n % 10)
		n = n / 10
	}
	return res
}

func getSum(n int64) int64 {
	s := splitNum(n)
	sum := int64(0)
	for _, v := range s {
		sum += v
	}
	return sum
}

type job struct {
	value int64
}

type result struct {
	job *job
	sum int64
}

var wg sync.WaitGroup

func producer(jobChan chan<- *job) {
	defer wg.Done()
	// 循环生成int64类型的随机数，发送到jobChan
	for {
		x := rand.Int63()
		newJob := &job {
			value: x,
		}
		jobChan <- newJob
		time.Sleep(time.Millisecond * 500)
	}
	
	
}

func consumer(jobChan <-chan *job, resultChan chan<- *result) {
	defer wg.Done()
	// defer close(resultChan)
	// 从jobChan中取出随机数计算各位数的和，将结果发送到resultChan
	for {
		jobTemp := <-jobChan
		sum := getSum(jobTemp.value)
		newResult := &result{
			job:jobTemp, 
			sum: sum,
		}
		resultChan <- newResult
	}
	
}

func main() {
	var jobChan = make(chan *job, 100)
	var resultChan = make(chan *result, 100)
	wg.Add(1)
	go producer(jobChan)
	// 开启24个goroutine执行consumer
	wg.Add(24)
	for i:=0; i<24; i++ {
		go consumer(jobChan, resultChan)
	}
	for result := range resultChan {
		fmt.Println(result.job.value, result.sum)
	}
	wg.Wait()
	close(jobChan)
	close(resultChan)
}
```
