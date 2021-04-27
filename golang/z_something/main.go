package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func F1(in chan interface{}, out chan interface{}, f func(data interface{}) interface{}, workerCount int) {
	poolLock := &sync.Mutex{}
	num := workerCount
	lock := func() bool {
		poolLock.Lock()
		defer poolLock.Unlock()
		if num > 0 {
			num--
			return true
		}
		return false
	}
	unlock := func() {
		poolLock.Lock()
		defer poolLock.Unlock()
		num++
	}
	lockMap := sync.Map{}
	index := 0
	for {
		select {
		case v := <-in:
			for {
				if lock() {
					break
				}
			}
			go func(vv interface{}, ff func(), n int) {
				defer func() {
					ff()
					lockMap.Store(n, 1)
				}()
				if n == 0 {
					out <- f(vv)
					return
				}
				for {
					if _, ok := lockMap.Load(n - 1); ok {
						out <- f(vv)
						lockMap.Delete(n - 1)
						break
					}
				}
			}(v, unlock, index)
			index++
		}
	}
}

func F2(v interface{}) interface{} {
	time.Sleep(time.Millisecond)
	return v
}

func main() {
	in, out := make(chan interface{}), make(chan interface{})
	var n = 10
	fmt.Printf("nihao:%s", "hiuohu")
	go F1(in, out, F2, n)

	index := 0
	go func() {
		for {
			in <- index
			index++
		}
	}()
	var result int
	for {
		select {
		case r := <-out:
			fmt.Printf("Result:%d, last:%d\n", r, result)
			if r.(int) < result {
				fmt.Errorf("ERR:%d, %d", r.(int), result)
				os.Exit(1)
			}
			result = r.(int)
		}
	}
}
