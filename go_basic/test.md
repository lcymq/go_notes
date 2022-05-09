# test

### 单元测试
- 白盒测试

- 测试命令：`go test`
    - `go test`命令是一个按照一定约定和组织的测试代码的驱动程序。
    - 在包目录里，所有以`_test.go`为后缀名的源代码文件都是`go test`测试的一部分，不会被`go build`编译到最终可执行文件中。
    
- 在`*_test.go`文件中有三种类型的函数：
    - 单元测试函数
    - 基准测试函数
    - 示例函数
    
    | 类型     | 格式                    | 作用                           |
    | -------- | ----------------------- | ------------------------------ |
    | 测试函数 | 函数名前缀为`Test`      | 测试程序的一些逻辑行为是否正确 |
    | 基准函数 | 函数名前缀为`Benchmark` | 测试函数的性能                 |
    | 示例函数 | 函数名前缀为`Example`   | 为文档提供示例文档             |
    
    `go test`命令会遍历所有的`*_test.go`文件中符合上述命名规则的函数，然后生成一个临时的main包用于调用相应的测试函数，然后构建并运行、报告测试结果，最后清理测试中生成的临时文件。

### 测试函数
##### 测试函数的格式
每个测试函数必须导入`testing`包，测试函数的基本格式（签名）如下：
```go
func TestName(t *testing.T) {
    // ...
}
```
测试函数的名字必须以`Test`开头，可选的后缀名必须以大写字母开头
```go
func TestSum(t *Testing.T) {}
```
参数`t`用于报告测试失败和附加的日志信息。
`testing.T`拥有的方法：
```go
func (c *T) Error(args ...interface{})
func (c *T) Fail()
func (c *T) FailNow()
func (c *T) Failed() bool
func (c *T) Fatal(args ...interface{})
func (c *T) Log(args ...interface{})
func (c *T) Logf(format string, args ...interface{})
func (c *T) Name() string
func (c *T) Parallel()
func (c *T) Run(name string, f func(t *T)) bool
func (c *T) Skip(args ...interface{})
func (c *T) SkipNow()
func (c *T) Skipf(format string, args ...interface{})
func (c *T) Skipped() bool
```
例子：
```go
func Split(str string, sep string) []string {
	res := []string{}
	temp := ""
	for i := range str {
		if string(str[i]) == sep {
            if temp != "" {
                res = append(res, temp)
			    temp = ""
            }
		} else {
			temp += string(str[i])
		}
	}
	if temp != "" {
		res = append(res, temp)
	}
	return res
}

func TestSplit(t *testing.T) {
	got := Split("a:b:", ":")
	want := []string{"a", "b"}
	if !reflect.DeepEqual(want, got) { // 切片是引用类型，不能直接比较
		t.Errorf("expected:%v, got:%v", want, got)
	} else {
		t.Log(got)
	}
}

// 测试组
func TestSplit(t *testing.T) {
    type test struct {
        str string
        sep string
        want []string
    }
    tests := map[string]test {
        "case_1": {str: "a:b:c", sep: ":", want: []string{"a", "b", "c"}},
        "case_2": {str: "a:b:", sep: ":", want: []string{"a", "b"}},
        "case_3": {str: "a::c", sep: ":", want: []string{"a", "c"}},
    }
    for name, tt := range tests {
        t.Run(name, func(t *testing.T) {
            got := Split(tt.str, tt.sep)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("expected:%v, got:%v", want, got)
            }
        }) 
    }
}
```
##### 相关命令：
- 跑所有测试：
    ```
    go test -v
    ```
- 跑某个测试用例：
    ```
    go test -run=TestSplit/case_3
    ```

##### 测试覆盖率
Go提供内置功能来检查测试覆盖率：
    ```
    go test -cover
    ```
`-coverprofile`参数，用来将覆盖率相关的记录信息输出到一个文件c.out：
    ```
    go test -cover -coverprofile=c.out
    ```
上面命令会将覆盖率相关信息输出到当前文件夹下面的`c.out`文件中，然后执行`go tool cover -html=c.out`，使用cover工具来处理生成的几率信息，该命令会打开本地的浏览器窗口生成一个HTML报告。
