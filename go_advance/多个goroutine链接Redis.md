# gorountine连接Redis
当goroutine连接Redis时，可以把调用这个goroutine的goroutine里的conn传过去吗？
答：不可以！

先上测试代码

```go
func f(conn redis.Conn) {
	// fmt.Println("golang连接redis")

	// conn, err := redis.Dial("tcp", "localhost:6379", redis.DialPassword(""))
	// if err != nil {
	// 	fmt.Println("连接redis失败",err)
	// 	return
	// }
	// defer conn.Close()
	_, err := conn.Do("Set", "golangtest", "yes")
	if err != nil {
		fmt.Println(err)
	}
	reply, err := redis.String(conn.Do("Get", "golangtest"))
	fmt.Println("键golangtest的值为: ", reply)
	fmt.Println("键golangtest的值为err: ", err)
}

func f2() {
	conn, err := redis.Dial("tcp", "localhost:6379", redis.DialPassword(""))
	if err != nil {
		fmt.Println("连接redis失败",err)
		return
	}
	defer conn.Close()
	fmt.Println("连接 redis success ...")
	reply, err := redis.String(conn.Do("Get", "golangtest"))
	fmt.Println("键golangtest的值为: ", reply)
	fmt.Println("键golangtest的值为err: ", err)
	for i:=0; i<3; i++ {
		go f(conn)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	f2()
	wg.Wait()
}
```



全程关注`conn`

- 链接上Redis之后，`conn`是有值的：

  ![image-20220325184744163](D:\workspace\test\src\go\go_advance\img\image-20220325184744163.png)

- 但真正执行`f()`的时候，conn显示已经关了：

  ![image-20220325184926606](D:\workspace\test\src\go\go_advance\img\image-20220325184926606.png)

得出结论：不能将Redis的`conn`传到子`goroutine`里。



### 尝试在子`goroutine`里重新链接Redis

```go
package main

import (
	"fmt"
	"sync"

	"github.com/gomodule/redigo/redis"
)

func f() {
	fmt.Println("golang连接redis")

	conn, err := redis.Dial("tcp", "localhost:6379", redis.DialPassword(""))
	if err != nil {
		fmt.Println("连接redis失败",err)
		return
	}
	defer conn.Close()
	_, err = conn.Do("Set", "golangtest", "yes")
	if err != nil {
		fmt.Println(err)
	}
	reply, err := redis.String(conn.Do("Get", "golangtest"))
	fmt.Println("键golangtest的值为: ", reply)
	fmt.Println("键golangtest的值为err: ", err)
}

func f2() {
	conn, err := redis.Dial("tcp", "localhost:6379", redis.DialPassword(""))
	if err != nil {
		fmt.Println("连接redis失败",err)
		return
	}
	defer conn.Close()
	fmt.Println("连接 redis success ...")
	reply, err := redis.String(conn.Do("Get", "golangtest"))
	fmt.Println("键golangtest的值为: ", reply)
	fmt.Println("键golangtest的值为err: ", err)
	for i:=0; i<3; i++ {
		go f()
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	f2()
	wg.Wait()
}
```

然后成功跑通啦！！

![image-20220325185220110](D:\workspace\test\src\go\go_advance\img\image-20220325185220110.png)



### 备注

有小伙伴问为什么要在`main`里加上`sync.WaitGroup`，这里统一回答一下：

是为了不让本进程在所有`goroutine`运行完之前停止。进程停了，由它创建的所有`goroutine`都被销毁了。
