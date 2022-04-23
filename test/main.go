package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var x int64
var wg sync.WaitGroup
var lock sync.Mutex

func add() {
	// lock.Lock()
	// x++
	// lock.Unlock()
	atomic.AddInt64(&x, 1)
	wg.Done()
}

func main() {
	x := int64(200)
	ok := atomic.CompareAndSwapInt64(&x, 200, 100)
	fmt.Println(ok, x)
}