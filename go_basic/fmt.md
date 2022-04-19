# fmt
### 键盘输入语句
- `fmt.Scanln()`
- `fmt.Scanf()`
```go
    var name string
	var age byte
	var sal float32
	var isPass bool
	fmt.Scanln(&name, &age, &sal, &isPass)
	fmt.Println(name, age, sal, isPass)
```
输入时候用空格隔开
```shell
name 10 10000000000.1279369842 true
```

### 格式化输入输出

- `%p`：输出指针
###### 通用占位符
- `%v`: 值的默认格式表示
- `%+v`: 类似`%v`, 但输出结构体时会添加字段名
- `%#v`: 输出详细信息，值的Go语法表示
- `%T`: 打印值得类型
