# struct

Go 语言没有“类”的概念，也不支持“类”的集成等面向对象的概念。Go中通过结构体的内嵌再配合接口比面向对象具有更高的扩展性和灵活性。

### 类型别名 和 自定义类型

##### 自定义类型
- 自定义类型用`type`关键字
- 自定义类型定义了一个全新的类型：可以基于内置的基本类型定义，也可以通过`struct`定义。
```go
type MyInt int
```
注意此时`MyInt`是一个全新的类型，它具有`int`的特性，但不是`int`！

验证：
```go
type MyInt int

func main() {
    var n MyInt
    n = 100
    fmt.Printf("%T %d\n", n, n) // main.MyInt 100
}
```

##### 类型别名
```go
type YourInt = int
```
此时`YourInt`还是`int`类型，只是起了一个别的名字

验证：
```go
type YourInt = int

func main() {
    var n YourInt
    n = 100
    fmt.Printf("%T %d\n", n, n) // int 100
}
```

### 结构体
```go
type 类型名 struct {
    字段名 字段类型
    字段名 字段类型
}
```

### 匿名结构体
多用于临时场景
```go
func main() {
    var s struct {
        x string
        y int
    }
    s.x = "..."
    s.y = 100
    fmt.Printf("%T %v\n", s, s) //struct { x string; y int } {... 100}
}
```

注意：
- 结构体是值类型
```go
type A struct {
    x string
    y int
}

func f(x A) {
    x.x = "yyy"
}

func main() {
    var a A
    a.x = "xxx"
    a.y = 111
    f(a)
    fmt.Println(a) //{xxx 111}, 修改的是副本
}
```
要想改变结构体内部值，需要传指针
```go
func f2(x *A) {
    x.x = "yyy"
    // (*x).x = "yyy" 不需要写成这样，go里有语法糖，自动根据指针找到对应的变量
}

func main() {
    var a A
    a.x = "xxx"
    a.y = 111
    f2(&a)
    fmt.Println(a) //{yyy 111}
}
```

##### 创建指针类型结构体
1. 通过`new`关键字对结构体进行实例化，得到的是结构体的地址。
```go
a := new(A)
fmt.Println(a) //&{ 0}
```

### 结构体初始化
1. 方法一：声明一个该结构体类型的变量
```go
var p person
p.name = "123"
p.age = 18
```
2. 方法二：键值对初始化 = 声明的同时初始化
```go
var p2 = person {
    name: "123",
    age: 18
}
```
3. 方法三：值列表初始化
```go
var p3 = person {
    "123",
    18,
}
```

### 方法和接收者
```go
// d代表了是那只狗调用wang这个方法
// d是接收者
// 接收者表示的是调用方法的具体类型变量，多用类型名首字母小写表示
func (d *dog) wang() {
    fmt.Printf("%s:汪汪汪~", d.name)
}
func main() {
    d1 := newDog("dahuang")
}
```
注意：只能给自己定义的类型添加方法
```go
func (i int) funcName() {
    // 这个方法不行，因为int是内置类型
}
```

### 结构体中的匿名字段
结构体允许其成员字段在声明时没有字段名而只有类型，这种没有名字的字段就称为匿名字段。
```go
type Person struct {
	string
	int
}
```
想要拿出来：
```go
p1 := &Person{
    "123",
    1,
}
fmt.Println(p1.string)
fmt.Println(p1.int)

```
匿名字段适用于：字段比较少也比较简单的场景
注意：不常用

### 嵌套结构体
结构体嵌套可以实现继承
用重命名类型来创建新的类型

### 结构体模拟实现其他语言的继承
```go
type animal struct {
    name string
}

// animal的方法
func (a *animal)move() {
    fmt.Printf("%s会动!", a.name)
}

// 狗类
type dog struct {
    feet uint8
    animal // animal拥有的方法，dog此时也拥有了
}

// 给dog实现一个汪汪汪的方法
func (d dog)wang() {
    fmt.Printf("%s is barking: wang~", d.name)
}

func main() {
    d1 := dog {
        animal: animal{name: "ddd"},
        feet: 4,
    }
    fmt.Println(d1) // {4 {ddd}}
    d1.wang() // ddd is barking: wang~
    d1.move() // ddd会动!
}
```

### 结构体 与 json

Go语言中，结构体与json能互相转换

##### 结构体 -> json : 序列化
```go
json.Marshal()
```
```go
func main() {
    p1 := Person{
        Name : "dog",
        Age : 2,
    }

    p, err := json.Marshal(p1)
    if err != nil {
        fmt.Printf("marshal failed, err: %v\n", err)
        return
    }
    fmt.Println(string(p)) //{"Name":"dog","Age":2}
}
```
当用json去解析的时候
```go
type AntiqueRepairLiquidUp struct {
	ID                int `json:"id"`
	ResourceID        int `json:"resourceId"`
	Amount            int `json:"amount"`
	AimResourceID     int `json:"aimResourceId"`
	DoubleProbability int `json:"doubleProbability"`
}
```
还可以写成多个：
```go
type Person struct {
    Name string `json:"name" db:"name" ini:"name"`
}
```

##### json -> 结构体 : 序列化
```go
str := `{"name":"墨祈","age":16}`
var p1 person
json.Unmarshal([]byte(str), &p1) //函数只能改p1的副本，不能改其本身，所以传引用进去
```

# 结构体

### type声明一种已存在的数据类型的别名

```go
type myint int
func main() {
	var a myint = 10
	fmt.Printf("%T\n", a) //main.myint
}
```

### type 定义一个结构体

```go
var book1 Book
```

### 结构体传递参数

```go
func changeBook1(book Book) {
	// 传递一个book的副本
	book.auth = "李四"
	fmt.Printf("%T, %p\n", book, &book)
	fmt.Println("book1 =", book)

}

func main() {
	var book1 Book
	book1.title = "张三奇遇记"
	book1.auth = "张三"
	fmt.Printf("%T, %p\n", book1, &book1)
	fmt.Println("book1 =", book1)
	changeBook1(book1)

	fmt.Printf("%T, %p\n", book1, &book1)
	fmt.Println("book1 =", book1)
}
```

输出：

```
main.Book, 0xc0000463c0
book1 = {张三奇遇记 张三}
main.Book, 0xc000046440
book1 = {张三奇遇记 李四}
main.Book, 0xc0000463c0
book1 = {张三奇遇记 张三}
```

想要真正改变book中的值，需要使用指针传递：

```go
type Book struct {
	title string
	auth  string
}

func changeBook1(book Book) {
	// 传递一个book的副本
	book.auth = "李四"
	fmt.Printf("%T, %p\n", book, &book)
	fmt.Println("book1 =", book)
}

func changeBook2(book *Book) {
	// 传递一个指针
	book.auth = "李四"
}

func main() {
	var book1 Book
	book1.title = "张三奇遇记"
	book1.auth = "张三"
	fmt.Printf("%T, %p\n", book1, &book1)
	fmt.Println("book1 =", book1)

	changeBook1(book1)

	fmt.Printf("%T, %p\n", book1, &book1)
	fmt.Println("book1 =", book1)

	changeBook2(&book1)
	fmt.Printf("%T, %p\n", book1, &book1)
	fmt.Println("book1 =", book1)

}
```

输出：

```
book1 = {张三奇遇记 张三}
main.Book, 0xc000046440
book1 = {张三奇遇记 李四}
main.Book, 0xc0000463c0
book1 = {张三奇遇记 张三}
main.Book, 0xc0000463c0
book1 = {张三奇遇记 李四}
```