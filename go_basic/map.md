# map
map是引用类型，必须初始化才能使用。

### map初始化
语法：
```go
m = make(map[T1]T2, size)
```
关于size：
- 若已知map中会有多少key-value：make的时候就写好，避免程序在运行期间再动态扩容
- 若size未知，则省略不写

以下用法是错误的：panic: assignment to entry in nil map，因为这个map还没有初始化
```go
var m1 map[int]int
m1[1]=1
fmt.Println(m1)
```
改正：
```go
var m1 map[int]int
m1 = make(map[int]int, 10)
m1[1]=1
fmt.Println(m1)
```
### map的三种声明方式
```go
// 第一种声明方式
var m1 map[int]string
fmt.Println("m1 =", m1)
m1 = make(map[int]string, 10)
m1[0] = "Go"
m1[1] = "C++"
m1[2] = "Java"
fmt.Println("m1 =", m1)
// 第二种声明方式
m2 := make(map[int]string)
m2[1] = "1"
m2[2] = "'2'"
fmt.Println("m2 =", m2)
// 第三种声明方式
m3 := map[int]string{
	1: "2",
	2: "3",
	3: "4",
}
fmt.Println("m3 =", m3)
```

### key不存在时
若不存在这个key，则拿到对应值类型的零值
```go
var m map[int]int
fmt.Println(m[1])
```


### map遍历
1. 同时遍历key和value
    ```go
    for k, v := range m1 {
        fmt.Println(k)
    }
    ```
2. 只遍历key
    ```go
    for k := range m1 {
        fmt.Println(k)
    }
    ```

### 删除map中的值
语法：
```go
delete(map, key)
```
1. 删除已存在的键值对 - 直接删除
    ```go
    delete(m, 1)
    ```
2. 删除不存在的键值对 - no op 无操作
    ```go
    delete(m, 2)
    ```

### 获取map中的所有key
较快方法：
```go
var keys := make([]int, 0, len(m))
for key := range m {
    keys = append(keys, key)
}
//若还需要排序
sort.Ints(keys)
//按照排序后的key遍历map
for _, key := range keys {
    fmt.Println(key, m[key])
}
```

### map的使用方式
- map用key,value形式存取数据
    ```go
    cityMap := make(map[string]string)
	cityMap["1"] = "1"
	cityMap["2"] = "2"
	cityMap["3"] = "3"

	for key, value := range cityMap {
		fmt.Println(key, value)
	}
    ```
- main传参是传的引用类型
    ```go
    func changeMap(cityMap map[string]string) {
	    cityMap["1"] = "2"
    }

    func main() {
        cityMap := make(map[string]string)
        cityMap["1"] = "1"
        cityMap["2"] = "2"
        cityMap["3"] = "3"      
        for key, value := range cityMap {
        	fmt.Println(key, value)
        }       
        changeMap(cityMap)
        for key, value := range cityMap {
        	fmt.Println(key, value)
        }       
    }

    ```
    输出:
    ```
    1 1
    2 2
    3 3
    1 2
    2 2
    3 3
    ```

### map的坑
map在遍历时可以增删！！
如果用`len(m)`来作为循环的结束条件，可能会陷入死循环。
比如下面这段程序直接让我的vscode崩溃了 QAQ
```go
func main() {
	m := make(map[int]int)
	for i:=0; i<10;i++ {
		m[i] = i
	}
	j := 0
	for len(m) != 0 {
		delete(m, j)
		delete(m, len(m)-1)
		m[j*10] = j*10
		fmt.Println(m)
	}
}
```