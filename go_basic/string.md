# string
golang中的基本数据类型之一


### 生成xx开头的字符串
```go
keys := make([]string, 0, 100)
for i:=0; i<100; i++ {
	key := fmt.Sprintf("xx%02d", i)
	keys = append(keys, key)
}
fmt.Println(keys)
```


### strconv包
- Atoi, Itoa
```go
n, err := strconv.Atoi(i)
s := strconv.Itoa(ret)
```

```go
func main() {
	i := "97"
	ret, err := strconv.Atoi(i)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v", ret) // 97
	n := strconv.Itoa(ret)
	fmt.Printf("%#v", n) // "97"
}
```

- parse
```go

```

```go
func main() {
	i := "97"
	ret, err := strconv.ParseInt(i, 10, 64) // 10进制，64位(int64)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v", ret) // 97
}
```

### unsafe.Sizeof()
```go
str1 := "abc"
println("string1:", unsafe.Sizeof(str1)) // 16
str2 := "abcdef"
println("string2:", unsafe.Sizeof(str2)) // 16
```

为什么字符串类型的 unsafe.Sizeof() 一直是16呢？
实际上字符串类型对应一个结构体，该结构体有两个域，第一个域是指向该字符串的指针，第二个域是字符串的长度，每个域占8个字节，但是并不包含指针指向的字符串的内容，这也就是为什么sizeof始终返回的是16。
组成可以理解成此结构体
```go
typedef struct {
    char* buffer;
    size_t len;
} string;
```
