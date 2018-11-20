# async

Simple Golang Async

Run your functions in parallel, and waiting for all functions end.  

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
    
    p := async.Parallel()
    p.AddFuncs(Func{Tag: "a", F: func() error {
        fmt.Println("a ", time.Now().UnixNano())
        return errors.New("a's error")
    }}, Func{Tag: "b", F: func() error {
        fmt.Println("b ", time.Now().UnixNano())
        return errors.New("b's error")
    }})
    
    errMap := p.Run()
    fmt.Printf("Parallel err: %+v\n", err
}

```

## TODO
- [x] Error
- [ ] Timeout
- [ ] Result
