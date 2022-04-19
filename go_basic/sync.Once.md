# sync.Once
`sync.Once`表示某个操作只执行一次。
例如：
- 只加载一次配置文件
- 只关闭一次通道

`sync.Once`只有一个`Do`方法:
```go
func (o *Once) Do(f func()) {}
```
注意：如果要执行的函数`f`需要传递参数就需要搭配闭包来使用。
```go
func f(ch1 <-chan int, ch2 chan<- int) {
	defer wg.Done()
	for {
		x, ok := <-ch1
		if !ok {
			break
		}
		ch2 <- x * x
	}
	f1 := func() {
		close(ch2)
	}
	once.Do(f)
}
```

### 加载配置文件示例：
延迟一个开销很大的初始化操作到真正用到它的时候再执行是一个很好的实践。因为预先初始化一个变量（比如在ini函数中完成初始化）会增加程序的启动耗时，而且有可能实际执行过程中这个变量没有用上，那么这个初始化操作就不是必须要做的。
```go
var icons map[string]image.Image

func loadIcons() {
	icons = map[string]image.Image{
		"left":  loadIcon("left.png"),
		"up":    loadIcon("up.png"),
		"right": loadIcon("right.png"),
		"down":  loadIcon("down.png"),
	}
}

func loadIcon() image.Image {
	return nil
}

// Icon被多个goroutine调用时不是并发安全的
func Icon(name string) image.Image {
	if icons == nil {
		loadIcons()
	}
	return icons[name]
}
```
以上代码存在的问题：
- 多个goroutine并发调用`Icon`函数时不是并发安全的，现代的编译器和CPU可能会在保证每个goroutine都满足串行一致的基础上自由重排访问内存的顺序。`loadIcons`函数可能会被重排为以下结果：
```go
func loadIcons() {
	icons = make(map[string]image.Image) // 当某一个goroutine执行了这一句后，icons便不再是nil，此时下面语句还没有执行，但其他goroutine不再能执行loadIcons()，此时会发生icons[name]找不到的情况
	icons["left"] = loadIcon("left.png"),
	icons["up"] = loadIcon("up.png"),
	icons["right"] = loadIcon("right.png"),
	icons["down"] = loadIcon("down.png"),
	
}
```

如果使用互斥锁解决该问题，会引发性能问题。
但可以使用`sync.Once`：
```go
var icons map[string]image.Image
var loadIconOnce sync.Once

func loadIcons() {
	icons = map[string]image.Image{
		"left":  loadIcon("left.png"),
		"up":    loadIcon("up.png"),
		"right": loadIcon("right.png"),
		"down":  loadIcon("down.png"),
	}
}

func loadIcon() image.Image {
	return nil
}

// Icon被是并发安全的
func Icon(name string) image.Image {
	loadIconOnce.Do(loadIcons)
	return icons[name]
}
```

