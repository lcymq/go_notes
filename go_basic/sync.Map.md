# sync.Map

Go中自带的map不是线程安全的。

例如：
```go
var m = make(map[string]int)

func get(key string) int {
	return m[key]
}

func set(key string, value int) {
	m[key] = value
}

func main() {
	wg := sync.WaitGroup{}
	for i:=0; i<10; i++ {
		wg.Add(1)
		go func(n int) {
			key := strconv.Itoa(n)
			set(key, n)
			fmt.Printf("k=%v, v=%v\n", key, get(key))
			wg.Done()
		}(i)
	}
	wg.Wait()
}
```
编译时没问题，运行时报错：
```
fatal error: concurrent map writes
k=0, v=0
```

# sync.Map

Go原生map不是线程安全的，对它进行读写操作时，需要加锁。
`sync.Map`是并发安全的，在Go1.9引入。

### sync.Map 
- `sync.Map`是线程安全的，读取，插入，删除都是常数级时间复杂度
- `sync.Map`的零值是有效的，并且零值是一个空的map
- `sync.Map`在第一次使用后，不允许被copy

### 没有sync.Map时，如何处理map的并发？
- 方法一：并发读写map的思路是加一把大锁 (缺点：锁粒度较大，影响效率)
- 方法二：把一个map分成若干个小map，对key进行哈希，只操作相应的小map (缺点：实现起来比较复杂，容易出错)
使用 `sync.map` 之后，对 map 的读写，不需要加锁。并且它通过空间换时间的方式，使用 read 和 dirty 两个 map 来进行读写分离，降低锁时间来提高效率。

### sync.Map使用方法
```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var m sync.Map

	// 1. 写入
	m.Store(1, "Go")
	m.Store(2, "C")

	// 2. 读取
	lang, _ := m.Load(1)
	fmt.Println(lang) // Go

	// 3. 遍历
	m.Range(func(key, value interface{}) bool {
		id := key
		lang := value
		fmt.Println(id, lang) // 1 Go 2 C
		return true
	})

	// 4. 删除
	m.Delete(2)
	lang, ok := m.Load(2)
	fmt.Println(lang, ok) // <nil> false

	// 5. 读取 / 写入
	m.LoadOrStore(3, "C++")
	lang, _ = m.Load(3)
	fmt.Println(lang) // C++
}
```
