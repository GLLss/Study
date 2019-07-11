# string不可以更改

## 不可更改内存

首先看一下string的真实结构：

![](../.gitbook/assets/image%20%281%29.png)

源码是这么定义header的：

```go
// StringHeader is the runtime representation of a string.
// It cannot be used safely or portably and its representation may
// change in a later release.
// Moreover, the Data field is not sufficient to guarantee the data
// it references will not be garbage collected, so programs must keep
// a separate, correctly typed pointer to the underlying data.
type StringHeader struct {
	Data uintptr
	Len  int
}
```

测试代码：

```go
package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func chanage(pstr *string) {
	*pstr = "world"
}

func main() {

	s := "hello"
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	fmt.Println(sh.Data, sh.Len)
	fmt.Println(&s)
	chanage(&s)
	fmt.Println(s)
	sh = *(*reflect.StringHeader)(unsafe.Pointer(&s))
	fmt.Println(sh.Data, sh.Len)
	fmt.Println(&s)

}

```

输出：

```text
4935491 5
0xc00000e1e0
world
4935546 5
0xc00000e1e0
```

可以看出：

* 修改string的值本质上是修改的Data和Len
* 修改string的值会重新分配内存，印证了官方不可以修改原内存的描述
* 直接修改和使用引用修改结果是一样的， 都会重新分配内存

