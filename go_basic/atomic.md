# atomic

| Method                                                       | explanation    |
| ------------------------------------------------------------ | -------------- |
| func LoadInt32(addr *int32)(val int32)<br/>func LoadInt64(addr *int64)(val int64)<br/>func LoadUint32(addr *uint32)(val uint32)<br/>func LoadUint64(addr *uint64)(val uint64)<br/>func LoadUintptr(addr *uintptr)(val uintptr)<br/>func LoadPointer(addr *unsafe.Pointer)(val unsafe.Pointer) | 读取操作       |
| func StoreInt32(addr *int32, val int32)<br/>func StoreInt64(addr *int64, val int64)<br/>func StoreUint32(addr *uint32, val uint32)<br/>func StoreUint64(addr *uint64, val uint64)<br/>func StoreUintptr(addr *uintptr, val uintptr)<br/>func StorePointer(addr *unsafe.Pointer, val unsafe.Pointer) | 写入操作       |
| func AddInt32(addr *int32, delta int32)(new int32)<br/>func AddInt64(addr *int64, delta int64)(new int64)<br/>func AddUint(addr *uint32, delta uint32)(new uint32)<br/>func AddUint64(addr *int64, delta int64)(new uint64)<br/>func AddUintptr(addr *uintptr, delta uintptr)(new uintptr) | 修改操作       |
| func SwapInt32(addr *int32, new int32)(old int32)<br/>func SwapInt64(addr *int64, new int64)(old int64)<br/>func SwapUint(addr *uint32, new uint32)(old uint32)<br/>func SwapUint64(addr *int64, new int64)(old uint64)<br/>func SwapUintptr(addr *uintptr, new uintptr)(old uintptr)<br/>func SwapPointer(addr *unsafe.Pointer, new unsafe.Pointer)(old unsafe.Pointer) | 交换操作       |
|                                                              | 比较并交换操作 |



### 同时给一个变量做递增操作

- 用`WaitGroup`的方法

  ```go
  package main
  
  import (
  	"fmt"
  	"sync"
  )
  
  var x int64
  var wg sync.WaitGroup
  
  func add() {
  	x++
  	wg.Done()
  }
  
  func main() {
  	wg.Add(1000)
  	for i:=0; i<1000; i++ {
  		go add()
  	}
  	wg.Wait()
  	fmt.Println(x)
  }
  ```

  结果每次都不一样：

  ![image-20220419021532488](D:\workspace\test\src\go\go_basic\img\image-20220419021532488.png)

- 给x加锁

  ```go
  package main
  
  import (
  	"fmt"
  	"sync"
  )
  
  var x int64
  var wg sync.WaitGroup
  var lock sync.Mutex
  
  func add() {
  	lock.Lock()
  	x++
  	lock.Unlock()
  	wg.Done()
  }
  
  func main() {
  	wg.Add(1000)
  	for i:=0; i<1000; i++ {
  		go add()
  	}
  	wg.Wait()
  	fmt.Println(x)
  }
  ```

  这次结果正确了:

  ![image-20220419021818542](D:\workspace\test\src\go\go_basic\img\image-20220419021818542.png)

- `atomic`包中的函数在底层实现了加锁操作

  ```go
  package main
  
  import (
  	"fmt"
  	"sync"
  	"sync/atomic"
  )
  
  var x int64
  var wg sync.WaitGroup
  var lock sync.Mutex
  
  func add() {
  	// lock.Lock()
  	// x++
  	// lock.Unlock()
  	atomic.AddInt64(&x, 1)
  	wg.Done()
  }
  
  func main() {
  	wg.Add(1000)
  	for i:=0; i<1000; i++ {
  		go add()
  	}
  	wg.Wait()
  	fmt.Println(x)
  }
  ```

  结果也是正确的：

  ![image-20220419022000493](D:\workspace\test\src\go\go_basic\img\image-20220419022000493.png)

Go源码中的`atomic`：

```assembly
TEXT runtime∕internal∕atomic·Xadd64(SB), NOSPLIT, $0-24
	MOVQ	ptr+0(FP), BX // 第一个参数保存到BX
	MOVQ	delta+8(FP), AX // 第二个参数保存到AX
	MOVQ	AX, CX  // 将第二个参数临时存到CX寄存器中
	LOCK			// LOCK指令进行锁住操作，实现对共享内存独占访问
	XADDQ	AX, 0(BX) // xaddq指令，实现寄存器AX的值与BX指向的内存存的值互换，
	// 并将这两个值的和存在BX指向的内存中，此时AX寄存器存的是第一个参数指向的值
	ADDQ	CX, AX ; 此时AX寄存器的值是Add操作之后的值，和0(BX)值一样
	MOVQ	AX, ret+16(FP) ; 返回值
	RET
```


另一个例子：
```go
func main() {
	x := int64(200)
	ok := atomic.CompareAndSwapInt64(&x, 200, 100)
	fmt.Println(ok, x)
}
```
