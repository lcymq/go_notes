# 基本数据类型 和 string 之间的转换

### 基本数据类型 转 string类型
###### 方式1：fmt.Sprintf("%参数",表达式)
- 参数需要和表达式的数据类型相匹配
- `fmt.Sprintf()`会返回转换后的字符串
    ```go
    var str string
	var num1 int = 99
	var num2 float64 = 123.45
	var b bool = true
	var myChar byte = 'a'

	str = fmt.Sprintf("%d", num1)
	fmt.Printf("type = %T, str = %q\n", str, str)
    // type = string, str = "99"

	str = fmt.Sprintf("%f", num2)
	fmt.Printf("type = %T, str = %q\n", str, str)
    // type = string, str = "123.450000"

	str = fmt.Sprintf("%t", b)
	fmt.Printf("type = %T, str = %q\n", str, str)
    // type = string, str = "true"

	str = fmt.Sprintf("%c", myChar)
	fmt.Printf("type = %T, str = %q\n", str, str)
    // type = string, str = "a"
    ```

###### 方式2：使用strconv包的函数
包中函数包括：
```go
    func FormatBool(b bool) string
    func FormatInt(i int64, base int) string
    func FormatUint(i uint64, base int) string
    func FormatFloat(f float64, fmt byte, prec, bitSize int) string
    func Itoa(i int) string
```
e.g.
```go
    var num3 int = 99
	var num4 float64 = 23.456
	var b2 bool = true

	str = strconv.FormatInt(int64(num3), 10)
	fmt.Printf("type = %T, str = %q\n", str, str)
    // type = string, str = "99"

	// 'f'是格式，10表示小数位保留10位，64表示这个小数是float64
	str = strconv.FormatFloat(num4, 'f', 10, 64)
	fmt.Printf("type = %T, str = %q\n", str, str)
    //type = string, str = "23.4560000000"

	str = strconv.FormatBool(b2)
	fmt.Printf("type = %T, str = %q\n", str, str)
    // type = string, str = "true"

    var num5 = 4567
	str = strconv.Itoa(num5)
	fmt.Printf("type = %T, str = %q\n", str, str)
    // type = string, str = "4567"
```

### string类型 转 基本数据类型
###### 使用strconv包的函数
```go
func ParseBool(str string) (value bool, err error)
func ParseFloat(s string, bitSize int) (f float64, err error)
func ParseInt(s string, base int, bitSize int) (i int64, err error)
func ParseUnit (s string, b int, bitSize int) (n uint 64, err error)
```
- string -> bool
    ```go
    var str string = "true"
	var b bool
	b, _ = strconv.ParseBool(str)
	fmt.Printf("type = %T, b = %t\n", b, b) 
    // type = bool, b = true
    ```
- string -> int
    ```go
    var str2 string = "123456789"
	var n1 int64
	n1, _ = strconv.ParseInt(str2, 10, 64)
	fmt.Printf("type = %T, b = %d\n", n1, n1)
    // type = int64, b = 123456789
    ```

- string -> float
    ```go
    var str3 string = "123.456"
	var f1 float64
	f1, _ = strconv.ParseFloat(str3, 10)
	fmt.Printf("type = %T, b = %f\n", f1, f1)
    // type = float64, b = 123.456000
    ```
- 注意：因为返回的是int64或float64，若希望得到int32，float32等，如下处理
    ```go
    var num5 int32
    num5 = int32(num)

### 注意事项
- 在将String类型转成基本数据类型时，要确保String类型转成有效的数据，比如我们可以把“123”转成一个整数，但不能把“hello”转成一个整数，如果这样做，Golang直接将其转换为0

    ```go
        var str2 string = "hello"
	    var n1 int64
	    n1, _ = strconv.ParseInt(str2, 10, 64) 
	    fmt.Printf("type = %T, b = %d\n", n1, n1)
        // type = int64, b = 0
    ```
- 
    ```go
        var str string = "hello"
	    var b bool
	    b, _ = strconv.ParseBool(str)
	    fmt.Printf("type = %T, b = %t\n", b, b)
        // type = bool, b = false
    ```