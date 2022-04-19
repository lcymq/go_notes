# select 多路复用
### select
- 可处理一个或多个channel的发送/接收操作
- 如果多个case同时满足，select会随机选择一个
- 对于没有case的`select{}`会一直等，可用于阻塞`main`函数

在某些场景下我们需要同时从多个通道接收数据。通道在接收数据时，如果没有数据可以接收将发生阻塞。
如下方法可以解决：
```go
for {
    // 尝试从ch1接收值
    data, ok := <- ch1
    // 尝试从ch2接收值
    data, ok := <- ch2
    // ...
}
```
这种方式虽然可以实现从多个通道接收值的需求，但是运行性能会差很多。为了应对这种场景，Go内置了`select`关键字，可以同时响应多个通道的操作。select的使用类似于switch语句，它有一些列case分支和一个默认的分支。每个case会对应一个通道的通信过程（接收或发送）。`select`会一直等待，知道某个`case`的通信操作完成时，就会执行`case`分支对应的语句。如下：
```go
select { // 从以下case中随机选一个
    case <- ch1: // 从ch1中取出值并废弃掉
        // ...
    case data := <- ch2: // 从ch2中取出值并放入data中
        // ...
    case ch3 <- data:  // 将data的值放进ch3里
        // ...
    default:
        // 默认操作
}
```
