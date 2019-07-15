# nil channel

在阅读nsq源码时发现，多个goroutine之间通讯channel在不使用时，并不是close而是设置为nil;一时好奇研究了下channel的close和设置为nil有什么不同，以下， 进入正题

### close channel

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

var r, w int

// 不断向channel c中发送[0,10)的随机数
func send(a chan int) {
	for {
		n := rand.Intn(10)
		a <- n
		w++
		fmt.Println(w, "PUT a <-", n)
	}
}

func channelSelect(a chan int) {
	// 1秒后，将向t.C通道发送时间点，使其可读
	t := time.NewTimer(1 * time.Second)
	for {
		select {
		case input, fClose := <-a:
			r++
			fmt.Println(r, "GET a ->", input, fClose)
		case <-t.C:
			fmt.Println("close a ......")
			//a = nil
			close(a)
		}
	}
}

func main() {
	a := make(chan int)
	go channelSelect(a)
	go send(a)
	// 给3秒时间让前两个goroutine有足够时间运行
	time.Sleep(300 * time.Second)
}

```

### 结果

```text
...
70086 PUT a <- 0
70087 PUT a <- 2
70087 GET a -> 2 true
close a ......
70088 GET a -> 0 false
...
70120 GET a -> 0 false
panic: send on closed channel

goroutine 6 [running]:
main.send(0xc00001e120)
        /root/workspace/go-path/src/learn/p.go:15 +0x58
created by main.main
        /root/workspace/go-path/src/learn/p.go:40 +0x7e
exit status 2
```

### 结论

* 已经close的channel不能写数据，否则panic
* 已经chose的channel人然可以读数据，读出的数据为对应类型的空值和关闭标识false
* 另，channel不能多次close否则panic

### 如果channel不是close，而是直接设置为nil呢？

修改代码：

```go
func channelSelect(a chan int) {
	// 1秒后，将向t.C通道发送时间点，使其可读
	t := time.NewTimer(1 * time.Second)
	for {
		select {
		case input, fClose := <-a:
			r++
			fmt.Println(r, "GET a ->", input, fClose)
		case <-t.C:
			fmt.Println("a = nil ......")
			a = nil
			//close(a)
		}
	}
}
```

### 结果：

```text
...
70235 GET a -> 8 true
70236 GET a -> 9 true
70236 PUT a <- 9
70237 PUT a <- 8
70237 GET a -> 8 true
a = nil ......
```

### 结论：

* channel为nil时，读写都会阻塞，因此select也不会进入此case
* 使用nil的好处就是不会迭代读取，禁用一个从channel读取数据的case，避免繁忙循环，提高性能

