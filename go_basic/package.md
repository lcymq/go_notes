# package 包

注意：只有main包可以编译成可执行文件

### 包的导入
- `import`导入语句通常放在文件开头包声明语句的下面。
- 导入的包名需要使用双引号包裹起来。
- 包名是从$GOPATH/src/$后开始计算的，使用`/`进行路径分隔。
- Go中禁止循环导入包：比如A包需要引入B包，B包又需要引入A包。
```go
import (
    "fmt"
)
```

#### 包的别名
```go
import (
    api "server/apigrpc"
)
```

### 匿名导入包
如果只希望导入包，而不是使用包内部的数据时，可以使用匿名导入包。具体格式如下：
```go
import _ "包的路径"
```
使用场景：比如需要其中的`init()`函数等。

### init() 函数
- `import`语句会自动触发(调用)包内部`init()`函数。
- `init()`没有参数也没有返回值
- `init()`只能被自动调用，不能在代码中主动调用。

### init()函数执行时机
1. 全局声明
2. `init()`
3. `main()`

![image-20220201235524677](C:\Users\cloris.l\AppData\Roaming\Typora\typora-user-images\image-20220201235524677.png)