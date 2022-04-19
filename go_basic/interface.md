# interface 接口
- 接口是一种类型，规定了变量有哪些方法
- 只要一个类型实现了该接口的所有方法，就是这种接口的类型
- 如果一个变量实现了接口中规定的所有方法，那么这么变量就实现了这个接口，可以称为这个接口类型的变量。
### 定义
```go
type interfaceName interface {
	method1(arg1, arg2, ...)(returnValue1, returnValue2, ...)
}
```
注意：下面这种用法不算是实现了`animal`接口，因为`eat(string)`和`eat()`不一样

![image-20220131190903839](C:\Users\cloris.l\AppData\Roaming\Typora\typora-user-images\image-20220131190903839.png)

### 接口的保存方式

接口的保存分为两部分：

值的类型和值的本身

![image-20220131193006889](C:\Users\cloris.l\AppData\Roaming\Typora\typora-user-images\image-20220131193006889.png)

初始时：动态类型和动态值都是nil

后面指定为什么类型和值就是什么类型和值，可以动态变化

![image-20220131193520248](C:\Users\cloris.l\AppData\Roaming\Typora\typora-user-images\image-20220131193520248.png)



### 同一个结构体可以实现多个接口

以下`cat`接口实现了多个`interface`

```go
type mover interface {
    move()
}
type eater interface {
    eat()
}

type cat struct {
    name string
    feet int8
}

func (c *cat) move() {
    fmt.Println("cat move...")
}

func (c *cat) eat() {
    fmt.Println("cat eat...")
}
```

### 接口可以嵌套 

```go
type animal interface {
    mover
    eater
}
```



### 看实例~

```go
type speaker interface {
    speak()
}

type cat struct {}
type dog struct {}
type person struct{}

func (c cat) speak() {
    fmt.Println("meow")
}
func (d dog) speak() {
    fmt.Println("Woof")
}
func (p person) speak() {
    fmt.Println("ah~")
}

func hit(x speaker) {
    // 接收一个参数，传进来什么就打什么
    x.speak()
}

func main() {

	var c cat
    hit(c) //meow

	var ss speaker
	hit(ss) //panic

	ss = c // 将c赋值给ss
	hit(ss)
}
```


```go
package server

import (
	"context"
	"database/sql"
	"errors"
	"server/rtapi"
	"sync"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

const (
	KeyAchievementActionLog = "ALog"
	HistoryKeptDays = 30
)

var (
	ErrAlogHistoryNotFound = errors.New("alog history not found")
	ErrAlogHistoryAlreadyExist = errors.New("alog history already exists")
)

type AlogEntity interface {
	Get(t AchievementConditionType) (int64, bool)

	Set(t AchievementConditionType, value int64)
	SetIf(t AchievementConditionType, predict func(x int64) bool, value int64)
	SetIfGreater(t AchievementConditionType, value int64)
	SetIfLess(t AchievementConditionType, value int64)
	Increase(t AchievementConditionType, value int32)

	GetDirty() map[AchievementConditionType]int64
	GetAll() map[AchievementConditionType]int64

	SetAllClean()

	GetAcks() []*rtapi.Envelope

	Apply(ctx context.Context, logger *zap.Logger) StorageOpWrites
}

type AlogEntityImplement struct {
	alog    *rtapi.ActionLog
	db      *sql.DB
	userID  uuid.UUID
	dirties map[int32]int64
}

func (a *AlogEntityImplement) Get(t AchievementConditionType) (value int64, ok bool) {
	if a.alog.Logs == nil {
		return 0, false
	}
	value, ok = a.alog.Logs[int32(t)]
	return
}

func (a *AlogEntityImplement) Set(t AchievementConditionType, value int64) {
	if a.alog.Logs == nil {
		a.alog.Logs = make(map[int32]int64)
	}
	original, ok := a.Get(t)
	if original != value || !ok {
		a.alog.Logs[int32(t)] = value
		a.setDirty(t)
	}
}

func (a *AlogEntityImplement) SetIf(t AchievementConditionType, predict func(x int64) bool, value int64) {
	original, ok := a.Get(t)
	if !ok || predict(original) {
		a.Set(t, value)
	}
}

func (a *AlogEntityImplement) SetIfGreater(t AchievementConditionType, value int64) {
	a.SetIf(t, func(x int64) bool { return value > x }, value)
}

func (a *AlogEntityImplement) SetIfLess(t AchievementConditionType, value int64) {
	a.SetIf(t, func(x int64) bool { return value < x }, value)
}

func (a *AlogEntityImplement) Increase(t AchievementConditionType, delta int32) {
	original, _ := a.Get(t)
	a.Set(t, original+int64(delta))
}

func (a *AlogEntityImplement) GetDirty() map[AchievementConditionType]int64 {
	ret := make(map[AchievementConditionType]int64)
	for k, v := range a.dirties {
		ret[AchievementConditionType(k)] = v
	}
	return ret
}

func (a *AlogEntityImplement) GetDirtyPB() map[int32]int64 {
	ret := make(map[int32]int64)
	for k, v := range a.GetDirty() {
		ret[int32(k)] = v
	}
	return ret
}

func (a *AlogEntityImplement) GetAll() map[AchievementConditionType]int64 {
	ret := make(map[AchievementConditionType]int64)
	for k, v := range a.alog.GetLogs() {
		ret[AchievementConditionType(k)] = v
	}
	return ret
}

func (a *AlogEntityImplement) SetAllClean() {
	a.dirties = make(map[int32]int64)
}

func (a *AlogEntityImplement) Apply(ctx context.Context, logger *zap.Logger) StorageOpWrites {
	logger.Info("Applying Alog Entity to Ops", zap.Any("alog", a.GetDirty()))
	ops := StorageConvertPB2Ops(ctx, logger, a.userID, CollectionAchievement, KeyAchievementActionLog, a.alog)
	ops = append(ops, a.appendHistory(ctx, logger)...)
	return ops
}

func (a *AlogEntityImplement) appendHistory(ctx context.Context, logger *zap.Logger) (ops StorageOpWrites) {
	monitorDirty := ALog.CheckDirty(a.GetDirty())

	if monitorDirty {
		alogHistory := &rtapi.ActionLogBackup{}
		StorageReadProtoObject(ctx, logger, a.db, a.userID, CollectionAchievementBackup, KeyAchievementBackup, alogHistory)

		backup := &rtapi.ActionLog{
			Logs: make(map[int32]int64),
		}
		for k, v := range a.alog.Logs {
			if ALog.monitors[AchievementConditionType(k)] {
				backup.Logs[k] = v
			}
		}
		if alogHistory.Backups == nil {
			alogHistory.Backups = make(map[int64]*rtapi.ActionLog)
		}
		alogHistory.Backups[int64(GetCurrentServerDay())] = backup
		a.historyFilter(alogHistory, int64(GetCurrentServerDay()))
		ops = append(ops, StorageConvertPB2Ops(ctx, logger, a.userID, CollectionAchievementBackup, KeyAchievementBackup, alogHistory)...)
	}
	return ops
}

func (a *AlogEntityImplement) historyFilter(history *rtapi.ActionLogBackup, currentServerDay int64)  {
	recentOutOfRange := int64(0)
	for k := range history.Backups {
		if k > currentServerDay - HistoryKeptDays {
			continue
		}
		if k > recentOutOfRange {
			delete(history.Backups, recentOutOfRange)
			recentOutOfRange = k
		} else {
			delete(history.Backups, k)
		}
	}
}

func (a *AlogEntityImplement) GetAcks() []*rtapi.Envelope {
	acks := make([]*rtapi.Envelope, 0)
	acks = append(acks, &rtapi.Envelope{Message: &rtapi.Envelope_AchievementConditionUpdate{
		AchievementConditionUpdate: &rtapi.AchievementConditionUpdate{
			Conditions: a.GetDirtyPB(),
		}}})
	return acks
}

func (a *AlogEntityImplement) setDirty(t AchievementConditionType) {
	value := a.alog.Logs[int32(t)]
	a.dirties[int32(t)] = value
}

type ALogWithBackup struct {
	AlogEntityImplement
	readonly  bool
	serverDay int
}

func (a *ALogWithBackup) Apply(ctx context.Context, logger *zap.Logger) StorageOpWrites {
	panic("cannot apply readonly alog because this is a historical snapshot")
}

func (a *ALogWithBackup) Set(t AchievementConditionType, value int64) {
	if a.ReadOnly() {
		panic("cannot set a readonly alog")
	}
	a.AlogEntityImplement.Set(t, value)
}
func (a *ALogWithBackup) ReadOnly() bool {
	return a.readonly
}

type AlogGlobal struct {
	mu       sync.RWMutex
	monitors map[AchievementConditionType]bool
}

func (a *AlogGlobal) Get(ctx context.Context, userID uuid.UUID, logger *zap.Logger, db *sql.DB) AlogEntity {
	alog := &rtapi.ActionLog{}
	StorageReadProtoObject(ctx, logger, db, userID, CollectionAchievement, KeyAchievementActionLog, alog)

	entity := AlogEntityImplement{
		alog:    alog,
		db:      db,
		userID:  userID,
		dirties: make(map[int32]int64),
	}
	return &entity
}

func (a *AlogGlobal) InitHistory(ctx context.Context, userID uuid.UUID, logger *zap.Logger, db *sql.DB) error {
	alogHistory := &rtapi.ActionLogBackup{}
	err := StorageReadProtoObjectCheckEmpty(ctx, logger, db, userID, CollectionAchievementBackup, KeyAchievementBackup, alogHistory)
	
	if err == nil {
		return ErrAlogHistoryAlreadyExist
	}
	if !IsDBCacheEmptyError(err) {
		return err
	} 
	alog := &rtapi.ActionLog{}
	StorageReadProtoObject(ctx, logger, db, userID, CollectionAchievement, KeyAchievementActionLog, alog)

	alogHistory.Backups = make(map[int64]*rtapi.ActionLog)
	alogHistory.Backups[1] = &rtapi.ActionLog{
		Logs: make(map[int32]int64),
	}
	for k, v := range alog.Logs {
		if a.monitors[AchievementConditionType(k)] {
			alogHistory.Backups[1].Logs[k] = v
		}
	}
	StorageWriteProtoObject(ctx, logger, db, userID, CollectionAchievementBackup, KeyAchievementBackup, alogHistory)
	return nil
}

func (a *AlogGlobal) GetHistory(ctx context.Context, userID uuid.UUID, logger *zap.Logger, db *sql.DB, serverDay int) (AlogEntity, error) {
	alogHistory := &rtapi.ActionLogBackup{}

	err := StorageReadProtoObjectCheckEmpty(ctx, logger, db, userID, CollectionAchievementBackup, KeyAchievementBackup, alogHistory)
	if err != nil {
		return nil, ErrAlogHistoryNotFound
	}
    var alog *rtapi.ActionLog
	if alog_, ok := alogHistory.Backups[int64(serverDay)]; ok {
		alog = alog_
	} else {
		nearest := int64(0)
		for k, alog_ := range alogHistory.Backups {
			if k < int64(serverDay) && k > nearest {
				nearest = k
				alog = alog_
			}
		}
	}
	if alog != nil {
		entity := &ALogWithBackup{
			AlogEntityImplement: AlogEntityImplement{
				alog:    alog,
				db:      db,
				userID:  userID,
				dirties: make(map[int32]int64),
			},
			readonly:  true,
			serverDay: serverDay,
		}
		return entity, nil
	} 
	return nil, ErrAlogHistoryNotFound
}

func (a *AlogGlobal) RegisterMonitor(t ...AchievementConditionType) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, tt := range t {
		a.monitors[tt] = true
	}
}

func (a *AlogGlobal) CheckDirty(in map[AchievementConditionType]int64) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	for k := range in {
		if a.monitors[k] {
			return true
		}
	}
	return false
}

var ALog = AlogGlobal{
	monitors: make(map[AchievementConditionType]bool),
}

```



### 空接口

- 什么方法都没有包含
- 任何结构体都实现了它

```go
type xx interface {
    
}
```

没有必要起名字，通常定义成下面的格式：

```go
interface{} // 空接口
```

所有的类型都实现了空接口，也就是任意类型的变量都能保存到空接口中。

注意：`interface`是关键字，`interface{}`阿此时空接口类型。

##### 例

```go
func show(a interface{}) {
    fmt.Printf("type: %T value:%v\n", a, a)
}

func main() {
    var m1 map[string]interface{}
    m1 = make(map[string]interface{},16)
    m1["name"] = "eoawhefo"
    m1["age"] = 9000
    
    
    show(false) //type: bool value:false
    show(nil) //type: <nil> value:<nil>
    show(m1) //type: map[string]interface {} value:map[age:9000 name:eoawhefo]
}
```

### 类型断言

想要判断空接口中的值，这个时候可以使用类型断言，语法格式：

```go
x.(T)
```

##### 例

```go
func assign(a interface{}) {
    fmt.Printf("%T\n", a)
    str, ok := a.(string)
    if !ok {
        fmt.Println("No")
    } else {
        fmt.Println("Yes, it is", str)
    }
}
```

##### 例

```go
func assign(a interface{}) {
    fmt.Printf("%T\n", a)
    switch a.(type) {
    	case string:
        fmt.Println("It is a string.")
        case int,int64, int32, int16, int8:
        fmt.Println("It is an int.")
        case bool:
        fmt.Println("")
    }
}
```

