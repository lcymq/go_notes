# defer
- defer是先进后出的结构

- return 比 defer 先执行

- 由于`defer`语句延迟调用的特性，`defer`语句能非常方便的处理资源释放问题。比如：资源清理、文件关闭、解锁及记录时间等。

  ![image-20220119222706377](C:\Users\cloris.l\AppData\Roaming\Typora\typora-user-images\image-20220119222706377.png)

```go
package main

import (
	"fmt"
)

func changeValue(p *int) {
	*p = 10
}

func swap(a *int, b *int) {
	*a ^= *b
	*b ^= *a
	*a ^= *b
}

func func1() int {
	fmt.Println("func1")
	return 1
}

func func2() int {
	fmt.Println("func2")
	return 2
}

func func3() int {
	fmt.Println("func3")
	return 3
}

func defer_call() (int, int, int) {
	defer func1()
	defer func2()
	defer func3()

	return func1(), func2(), func3()
}

func main() {
	defer_call()
}

```

输出内容为：
```
func1
func2
func3
func3
func2
func1
```

### defer 题目

###### 执行顺序
```go
func main() {
    fmt.Println("start")
    defer fmt.Println(1)
    defer fmt.Println(2)
    defer fmt.Println(3)
    fmt.Println("end")
}
```
解答：先顺序执行defer以外的语句，再倒序执行defer语句
```
start
end
3
2
1
```

```go
func f1() int {
    x := 5
    defer func() {
        x++
    }()
    return x // 1.返回值赋值 2. defer 3. 真正的RET指令
}

func f2() (x int) {
    defer func() {
        x++
    }()
    return 5 // 返回值=x
}

func f3() (y int) {
    x := 5
    defer func() {
        x++ // 修改的是x
    }()
    return x // 返回值 = y = x = 5, defer改的是x不是y
}

func f4() (x int) {
    defer func(x int) {
        x++ // 改变的是函数中x的副本
    }(x)
    return 5 //返回值 = x = 5
}

func f5() (x int) {
    defer func(x int) int {
        x++
        return x // return了但没人接收
    }(x)
    return 5
}

func f6() (x int) {
    defer func(x *int) { // 传x的指针到匿名函数
        (*x)++
    }(&x)
    return 5
}

func main() {
    fmt.Println(f1()) // 5
    fmt.Println(f2()) // 6
    fmt.Println(f3()) // 5
    fmt.Println(f4()) // 5
    fmt.Println(f5()) // 5
    fmt.Println(f6()) // 6
}
```

```go
func calc(index string, a, b int) int {
    ret := a + b
    fmt.Println(index, a, b, ret)
    return ret
}

func main() {
    a := 1
    b := 2
    defer calc("1", a, calc("10", a, b))
    a = 0
    defer calc("2", a, calc("20", a, b))
    b = 1
}
```
输出：
```go
10 1 2 3
20 0 2 2
2 0 2 2
1 1 3 4
```
解:
`defer`中的函数会执行到只剩最外一层，然后放进stack
0. stack上压入`a:=1` `b:=2`
1. `calc("1", a, calc("10", a, b))`
2. 算出`calc("10", a, b)`, 输出`"10" 1 2 3`
3. stack上存了`calc("1", a, 3)`
4. `a = 0`
5. `calc("2", a, calc("20", a, b))`
6. 算出`calc("20", a, b)`, 输出`"20", 0, 2, 2`
7. stack上存了`defer calc("2", a, 2)`
8. `b = 1`
9. `calc("2", a, 2)`, 输出`"2", 0, 2, 2`
10. `calc("1", a, 3)`, 输出`"1", 1, 3, 4`而不是`"1", 0, 3, 3`, 因为在stack中存的参数`a=1`