# async

Simple Golang Async

```golang

import (
    "github.com/lily-lee/async"
)

func main() {
    a := async.New()
    a.AddFuncs(func() error {
        fmt.Println("here")
    })
    
    a.Run()
}

```
