# go routine

### 开启go routine
- 程序启动之后会创建一个主goroutine去执行main函数
```go
go myFunc() //开启一个单独的go routine去执行myFunc()函数
```

### go routine调用匿名函数
```go
func main() {
	for i := 0; i < 1000; i++ {
		go func() {
			fmt.Println(i)
		}()
	}
	fmt.Println("main")
	time.Sleep(time.Second)
}
```
输出结果会出现有很多重复数字的这种情况：
```
547
959
817
962
962
818
963
745
963
822
963
963
963
822
963
964
964
964
964
964
964
964
964
965
965
965
965
965
965
965
965
745
965
965
```
解决方法：把i当做参数传进去
```go
func main() {
	for i := 0; i < 1000; i++ {
		go func(i int) {
			fmt.Println(i) //用的是函数参数的那个i，不是外面的i
		}(i)
	}
	fmt.Println("main")
	time.Sleep(time.Second)
}
```

### `goroutine`什么时候结束
- `goroutine`对应的函数结束了，`goroutine`就结束了
- `main`函数执行完了，由`main`函数创建的那些`goroutine`也都结束了

那么，如何让`main`函数知道它创建的`goroutine`都结束了？
### sync.WaitGroup

```go
var wg sync.WaitGroup

func f(i int) {
	defer wg.Done() //做完一次f()计数器减1
	fmt.Println(i)
}

func main() {
	for i:=0; i<10; i++ {
		wg.Add(1) //每次进循环则计数器加1
		go f(i)
	}
	wg.Wait() //等待wg的计数器减为0，不会实时检查，只有运行到这一步才检查计算器是否为0
	fmt.Println("main")
}
```

### goroutine和线程的区别
##### 可增长的栈
- 操作系统线程一般都有固定的栈内存(通常是2MB)
- 一个`goroutine`的栈在其声明周期开始时只有很小的栈(典型情况下2KB)
- `goroutine`的栈大小是可变的，可以按需增大或缩小
- `goroutine`的栈大小限制可以达到1GB，但极少会用到这么大空间
- 因为一般`goroutine`栈的初始大小很小，所以一次创建十万左右的`goroutine`也是可以的

### goroutine调度 - GMP模型
`GMP`是Go语言运行时(runtime)层面的实现，是Go语言自己实现的一套调度系统，和操作系统调度的OS线程不一样。
- `G`(goroutine): 本`goroutine`信息 + 与`P`绑定的信息
- `M`(machine): 是Go运行时`runtime`对操作系统内核线程的虚拟，`M`与内核线程是一一对应的关系。每个`goroutine`最终是要放到`M`上执行的
- `P`(Processor): 管理一组`goroutine`队列，P里面会存储当前`goroutine`运行的上下文环境(函数指针，堆栈地址，地址便捷)

- `P`与`M`是一一对应的：
    - `P`管理着一组`G`挂在`M`上运行
    - 当一个`G`长久阻塞在一个`M`上时，runtime会新建一个`M`，阻塞`G`所在的`P`会把其他的`G`挂载在新建的`M`上。
    - 当旧的`G`阻塞完成或者认为其已经死掉时，回收旧的`M`。

- `P`的个数是通过`runtime.GOMAXPROCS`设定(最大256)，Go1.5版本之后默认为物理线程数。在并发最大的时候会增加一些`P`和`M`，但不会太多，切换太频繁的话得不偿失。

- 单从线程调度讲
    - 其他语言的OS线程是由OS内核来调度的，`goroutine`则是由Go运行时runtime自己的调度器调度的，
    - Go的调度器使用一个称为$m:n$的调度技术(复用/调度m个goroutine到n个OS线程)
    - `goroutine`的调度是在用户态下完成的，不涉及内核态和用户态之间的频繁切换，包括内存的分配与释放，都是在用户态维护者一大块内存池，不直接调用系统的`malloc`函数(除非内存池需要改变)，成本比调度OS线程低很多。
    - 充分利用了多核的硬件资源，近似的把若干`goroutine`均分在物理线程上，再加上本身`goroutine`的超轻量，保证了Go调度方面的性能。
    - (Java是OS完成线程的切换的，而Go是runtime切换的)

```go
func a(wg *sync.WaitGroup) {
	defer wg.Done()
	for i:=0; i<100; i++ {
		fmt.Printf("A:%d\n", i)
	}
}

func b(wg *sync.WaitGroup) {
	defer wg.Done()
	for i:=0; i<100; i++ {
		fmt.Printf("B:%d\n", i)
	}
}

func main() {
	// runtime.GOMAXPROCS(1) // 1 = 只用一个P干活
	// runtime.GOMAXPROCS(16) // 16 = 任务被分到16个线程去干活，本机有16个CPU，16表示跑满了整个CPU（默认是CPU的逻辑核心数，默认跑满整个CPU
	runtime.GOMAXPROCS(0) // 参数<1，GOMAXpROCS = CPU逻辑核心数
	runtime.GOMAXPROCS(runtime.NumCPU()) // 获取CPU个数，再当做参数传进GOMAXPROCS
	var wg sync.WaitGroup
	wg.Add(2)
	go a(&wg)
	go b(&wg)
	wg.Wait()
}
```