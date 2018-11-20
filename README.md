# async

Simple Golang Async

## Usage

```golang

import (
    "github.com/lily-lee/async"
)

func main() {
    a := async.New()
    a.AddFuncs(func() error {
        fmt.Println("here")
    })
    
    if err:= a.Run(); err != nil {
        // handle error
    }
}

```

## TODO
- [x] Error
- [ ] Timeout
- [ ] Result