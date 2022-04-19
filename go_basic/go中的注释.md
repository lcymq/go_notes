# go中的注释

### 行注释
```go
// this is a line comment
```
### 块注释
```go
/*
this is a block comment
*/
```
### 函数(方法)注释
```go
/**
 * @description: connect to the database
 * @param {*psqlInfo} host
 * @return {*sql.DB} database
 */

```
### 结构(接口)注释
```go
// User   用户对象，定义了用户的基础信息
type User struct{
    Username  string // 用户名
    Email     string // 邮箱
}
```

### 包注释
```go
// @Title  请填写文件名称（需要改）
// @Description  请填写文件描述（需要改）
// @Author  请填写自己的真是姓名（需要改）  ${DATE} ${TIME}
// @Update  请填写自己的真是姓名（需要改）  ${DATE} ${TIME}
package ${GO_PACKAGE_NAME}
```


## 正确的代码规范
- 官方推荐使用行注释而不是块注释
- 使用`gofmt main.go`命令可以对go代码进行自动规范化(工作中看其他人写的而定)
- 使用`gofmt -w main.go`命令可以对go代码进行规范化并将改动写入代码中
- 运算符、变量左右要加空格