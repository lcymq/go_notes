# func

### 语法
```go
func 函数名(参数)(返回值){
    函数体
}
```

### 参数的类型简写
当参数中连续多个参数的类型一致时，我们可以将非最后一个参数的类型省略
```go
func f(x, y, z int, m, n string) int {
    return x + y
}
```

### 可变长参数
```go
func f(x string, y ...int) { // y的类型是切片 []int

}
```

### 没有默认参数
Go语言中没有默认参数的概念


### defer
Go语言中的`defer`语句会将其后面跟随的语句进行延迟处理。在`defer`归属的函数即将返回时，将延迟处理的语句按`defer`定义的逆序进行执行，也就是说，先被`defer`的语句最后被执行，最后被`defer`的语句最先被执行。

- `defer`语句执行的时机：
    ```go
        return x // 返回值=x; 运行defer; RET指令
    ```
- Go语言中函数的return不是原子操作，在底层是分为两步来执行的：
    1. 返回值赋值
    2. 真正的RET返回


### 函数作为参数类型
```go
func f(x func() int) {
    ret := x()
    fmt.Println(ret)
}
```

### 函数作为返回值
```go
func f0() {

}
func f(x func() int) func() {
    return f0
}
func main() {
	fmt.Println(f) //0x516420 一个地址
}
```

### 匿名函数
- 没有名字的函数
- 函数内部没有办法声明带名字的函数，可以用匿名函数
    ```go
    func main() {
    	f := func(x, y int) {
    		fmt.Println(x + y)
    	}
    	f(10, 20)
    }
    ```
- 如果只是调用一次的函数，还可以写成**立即调用函数**
    ```go
    func(x, y int) {
		fmt.Println(x + y)
	}(10, 20)
    ```
### 闭包
闭包 = 函数 + 环境
```go
func x(m, n int) {
	fmt.Println(m+n)
}

// 定义一个函数对f2进行包装
func f(x func(int, int), m, n int) {
	temp := func() {
		x(m, n)
	}
	temp()
}

func main() {
	f(x, 10, 20) // 30
}
```

```go
func x(m, n int) {
	fmt.Println(m+n)
}

// 定义一个函数对f2进行包装
func f(x func(int, int), m, n int) func() {
	temp := func() {
		x(m, n)
	}
	return temp
}

func main() {
	ret := f(x, 10, 20)
	ret() // 实现了调用ret即为调用x
}
```

```go
func adder() func(int) int {
	var x = 100
	return func(y int) int {
		x += y
		return x
	}
}

func main() {
	ret := adder()
	fmt.Println(ret(200))
}
```

```go
func calc(base int) (func(int) int, func(int) int) {
	add := func(i int) int {
		base += i
		return base
	}

	sub := func(i int) int {
		base -= i
		return base
	}
	return add, sub
}

func main() {
	f1, f2 := calc(10)
	fmt.Println(f1(1), f2(2)) // 11 9
	fmt.Println(f1(3), f2(4)) // 12 8
	fmt.Println(f1(5), f2(6)) // 13 7
    // 这三行共用同一个base(10)
}
```

### 闭包实现多态?
```go
func makeSuffixFunc(suffix string) func(string) string {
    return func(name string) string {
        if !strings.HasSuffix(name, suffix) {
            return name + suffix
        }
        return name
    }
}

func main() {
    jpgFunc := makeSuffixFunc(".jpg")
    txtFunc := makeSuffixFunc(".txt")

    fmt.Println(jpgFunc("test"))
    fmt.Println(jpgFunc("jpg"))
    fmt.Println(txtFunc("test"))
    fmt.Println(txtFunc("txt"))
    
}
```

### 闭包的延迟绑定
```go
func foo0() func() {
    x := 1
    f := func() {
        fmt.Printf("foo0 val = %d\n", x)
    }
    x = 10
    return f
}

func main() {
	foo0()() // 11
}
```

```go
func foo7(x int) []func() {
    var fs []func()
    values := []int{1, 2, 3, 5}
    for _, val := range values {
        fs = append(fs, func() {
            fmt.Printf("foo7 val = %d\n", x+val)
        })
    }
    return fs
}

func main() {
	f7s := foo7(11)
	for _, f7 := range f7s {
		f7()
	}
}
```