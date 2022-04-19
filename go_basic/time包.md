# time包

### 常用方法
```go
now = time.Now() // 当前时间:2022-02-03 17:51:26.2889372 -0800 PST m=+0.005177001
now.Year() //2022
now.Month() //February
now.Date() //2022 February 3
now.Hour() //17
now.Minute() //51
now.Second() //26
```

### 时间戳
```go
now = time.Now()
timestamp1 := now.Unix()      //秒级时间戳      1643940358
timestamp2 := now.UnixMilli() //毫秒级时间戳    1643940358828
timestamp3 := now.UnixMicro() //微秒级时间戳    1643940358828319
timestamp4 := now.UnixNano()  //纳秒级时间戳    1643940358828319100
```

### time包中的常量
```go
const (
	Nanosecond  Duration = 1
	Microsecond          = 1000 * Nanosecond
	Millisecond          = 1000 * Microsecond
	Second               = 1000 * Millisecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
)
```

### 时间的比较和运算
#### Add
`func (t Time) Add(d Duration) Time`
```go
func main() {
    now := time.Now()
    later := now.Add(time.Hour)
    fmt.Println(later)
}
```

#### Sub
`func (t Time) Sub(u Time) Duration`
- 返回一个时间段`t-u`
- 如果结果超出了Duration可以表示的最大值/最小值，将返回最大值/最小值。要获取时间点`t-d`（d为Duration），可以使用`t.Add(-d)`

#### Equal
`func (t Time) Equal(u Time) bool`
判断两个时间是否相同，会考虑时区的影响，因此不同时区标准的时间也可以正确比较。本方法和`t==u`不同，这种方法还会比较地点和时区信息。

#### Before
`func (t Time) Before(u Time) bool`
如果t代表的时间点在u之前，返回真；否则返回假。

#### After
`func (t Time) After(u Time) bool`
如果t代表的时间点在u之后，返回真；否则返回假。

#### UTC
`func (t Time) UTC() Time`
返回采用UTC和零时区，但指向同一时间点的Time

### Parse
```go
func f2() {
	now := time.Now() //本地的时间
	fmt.Println(now)
	//按照指定格式解析一个字符串格式的时间
	time.Parse("2006-01-02 15:04:05", "2022-02-05 12:36:00")

	//按照东八区的时区和格式解析一个字符串格式的时间
	//根据字符串加载时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("load loc failed, err:%v\n", err)
		return
	}
	//按照指定时区解析时间
	timeObj, err := time.ParseInLocation("2006-01-02 15:04:05", "2022-02-05 12:36:00", loc)
	if err != nil {
		fmt.Printf("parse time failed, err: %v\n", err)
	}
	fmt.Println(timeObj)
	//时间对象相减
	td := timeObj.Sub(now)
	fmt.Println(td)
}
```

### 定时器
使用`time.Tick(时间间隔)`来设置定时器，定时器本质山是一个通道(channel)。
```go
timer := time.Tick(time.Second)
for t := range timer {
	fmt.Println(t) // 一秒钟执行一次
}
```

### 时间格式化
时间类型有一个自带的方法`Format`进行格式化，需要注意的是Go语言中格式化时间模板不是常见的`Y-m-d H:M:S`而是使用Go的诞生时间2006年1月2日15点04分05秒（记忆口诀：2006 1 2 3 4 5）。
```go
fmt.Println(now.Format("2006-01-02"))
fmt.Println(now.Format("2006/01/02 15:04:05"))
fmt.Println(now.Format("2006/01/02 03:04:05 PM"))
fmt.Println(now.Format("2006/01/02 03:04:05.000 PM")) //带毫秒
```

### int64和timestamppb.Timestamp的转换
```go
func Int64ToTime(t int64) *timestamppb.Timestamp {
    timeS, _ := time.Parse(TIMEFORMAT, strconv.Itoa(int(t)))
    time := &timestamp.Timestamp{Seconds: timeS.Unix()}
    return time
}
```