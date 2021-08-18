# async

Simple Golang Async

Run your functions in parallel, and waiting for all functions end.  

## Usage

### example 1

```golang

import (
	"fmt"
	"runtime"
	"time"

	"github.com/lily-lee/async"
)

func main() {
    a := async.New()
	a.SetPoolNum(10)
	a.SetTimeout(time.Second)

	for i := 0; i < 100; i++ {
		index := i
		a.AddFunc(func() error {
			fmt.Println("index: ", index, " goroutine num: ", runtime.NumGoroutine())
			return nil
		})
	}
	a.AddFunc(func() error {
		fmt.Println("handle a")
		return nil
	})

	a.AddFunc(func() error {
		fmt.Println("handle b")
		return nil
	})

	runtime.GOMAXPROCS(1)
	// a.Run() returns an error case any func has an error.
	if err := a.Run(); err != nil {
		fmt.Println("err: ", err)
	}
}

```


### example 2

```golang
import (
	"fmt"
	"runtime"
	"time"

	"github.com/lily-lee/async"
)

func main() {
	p := async.Parallel()
	p.SetTimeout(time.Second)
	p.SetPoolNum(10)

	for i := 0; i < 100; i++ {
		index := i
		p.AddFunc(async.Func{
			Tag: fmt.Sprintf("index-%d", index),
			F: func() error {
				fmt.Println("index: ", index, " goroutine num: ", runtime.NumGoroutine())
				return nil
			},
		})
	}

	runtime.GOMAXPROCS(1)
	// p.Run() returns all funcs exec result.
	errResultMap := p.Run()
	fmt.Println("errResultMap: ", errResultMap)
}
```

## TODO
- [x] Error
- [x] Timeout
- [x] SetPoolNum
- [ ] Result
- [ ] context.Context
