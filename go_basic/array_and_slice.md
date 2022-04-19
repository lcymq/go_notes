# Array & Slice

### slice 初始化

```go 
// 声明slice是一个切片，并且初始化，默认值是1,2,3.长度len是3
slice1 := make([]int, 3) // 开辟3个空间，默认值是0
// 
slice1[0] = 100
fmt.Println(slice1) 
```

```
输出：
[100 0 0]
```



### slice开辟内存空间

```go
// 声明slice是一个切片，同时给slice分配空间，3个空间，初始化值是0
var slice1 := make
slice1 := make([]int, 3)
var slice2 []int = make([]int, 3)
fmt.Println(slice1)
```

```
输出：
[0 0 0]
```

### 判断slice是否为空

```go
var slice1 []int
if slice1 == nil {
	fmt.Println("this is nil")
} else {
	fmt.Println("this is not nil")
}
```

### 切片的追加

- slice的length和capacity
	```go
	// 创建一个numbers切片，length=3，capacity=5
	var numbers = make([]int, 3, 5)
	```
- 追加时，若length>capacity, 则向后追加与之前同样大小的capacity
	```go
	var numbers = make([]int, 3, 5)
	for i := 0; i < 10; i++ {
		numbers = append(numbers, 6)
		fmt.Printf("len=%d, cap=%d, slice=%v\n", len(numbers), cap(numbers), numbers)
	}
	```
	输出：
	```
	len=4, cap=5, slice=[0 0 0 6]
	len=5, cap=5, slice=[0 0 0 6 6]
	len=6, cap=10, slice=[0 0 0 6 6 6]
	len=7, cap=10, slice=[0 0 0 6 6 6 6]
	len=8, cap=10, slice=[0 0 0 6 6 6 6 6]
	len=9, cap=10, slice=[0 0 0 6 6 6 6 6 6]
	len=10, cap=10, slice=[0 0 0 6 6 6 6 6 6 6]
	len=11, cap=20, slice=[0 0 0 6 6 6 6 6 6 6 6]
	len=12, cap=20, slice=[0 0 0 6 6 6 6 6 6 6 6 6]
	len=13, cap=20, slice=[0 0 0 6 6 6 6 6 6 6 6 6 6]
	```

- 若不指定capacity，则默认capacity=length
	```go
	var numbers1 = make([]int, 3)
	fmt.Printf("len=%d, cap=%d, slice=%v\n", len(numbers1), cap(numbers1), numbers1)
	```

### 切片的截取

- 切片的截取，包含第一个元素，不包含第二个元素
	```go
	s := []int{1, 2, 3}
	s1 := s[0:2]
	fmt.Println(s1) // [1 2]
	s1 = s[:]
	fmt.Println(s1) // [1 2 3]
	s1 = s[:3]
	fmt.Println(s1) // [1 2 3]
	s1 = s[2:]
	fmt.Println(s1) // [3]
	```
	
- 截取的切片和被截取的切片，指向同一块内存区域
- 所有从同一个slice上截取的切片，capacity都相同

  ```go
	s := []int{1, 2, 3}
  
	s1 := s[0:1]
	s2 := s[0:2]
  
	fmt.Printf("%p\n", &s[0])
	fmt.Printf("%p\n", &s1[0])
	fmt.Printf("%p\n", &s2[0])
  
	s1[0] = 10
	fmt.Println("len(s) =", len(s), "cap(s) =", cap(s), s)
	fmt.Println("len(s1) =", len(s1), "cap(s1) =", cap(s1), s1)
	fmt.Println("len(s2) =", len(s2), "cap(s2) =", cap(s2), s2)
  
	s2 = append(s2, 100)
	fmt.Printf("%p\n", &s[0])
	fmt.Printf("%p\n", &s1[0])
	fmt.Printf("%p\n", &s2[0])
  
	s1[0] = 10
	fmt.Println("len(s) =", len(s), "cap(s) =", cap(s), s)
	fmt.Println("len(s1) =", len(s1), "cap(s1) =", cap(s1), s1)
	fmt.Println("len(s2) =", len(s2), "cap(s2) =", cap(s2), s2)
  
	s = append(s, 10)
	s = append(s, 10)
	s = append(s, 10)
	s = append(s, 10)
	fmt.Println("len(s) =", len(s), "cap(s) =", cap(s), s)
	fmt.Println("len(s1) =", len(s1), "cap(s1) =", cap(s1), s1)
	fmt.Println("len(s2) =", len(s2), "cap(s2) =", cap(s2), s2)
	fmt.Printf("%p\n", &s[0])
	fmt.Printf("%p\n", &s1[0])
	fmt.Printf("%p\n", &s2[0])
  ```

  输出：
  
  ```
  0xc0000101c8
  0xc0000101c8
  0xc0000101c8
  len(s) = 3 cap(s) = 3 [10 2 3]
  len(s1) = 1 cap(s1) = 3 [10]
  len(s2) = 2 cap(s2) = 3 [10 2]
  0xc0000101c8
  0xc0000101c8
  0xc000004078
  len(s) = 3 cap(s) = 3 [10 2 100]
  len(s1) = 1 cap(s1) = 3 [10]
  len(s2) = 3 cap(s2) = 3 [10 2 100]
  len(s) = 7 cap(s) = 12 [10 2 100 10 10 10 10]
  len(s1) = 1 cap(s1) = 3 [10]
  len(s2) = 3 cap(s2) = 3 [10 2 100]
  0xc00001a0c0
  0xc0000101c8
  0xc0000101c8
  ```
  
  - 被截取的切片capacity的改变不会影响切片的大小

### 切片的复制

`copy` 是深度复制

```go
s := []int{1, 2, 3}
s1 := s[0:2]
s2 := make([]int, 3)
copy(s2, s)
fmt.Printf("%p\t%v\n", &s[0], s)
fmt.Printf("%p\t%v\n", &s1[0], s1)
fmt.Printf("%p\t%v\n", &s2[0], s2)
```
输出：
```
0xc0000101c8    [1 2 3]
0xc0000101c8    [1 2]
0xc0000101e0    [1 2 3]
```

可见s和s1指向同一块内存地址，而s2是不同的内存地址