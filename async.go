package async

import (
	"log"
	"runtime/debug"
	"sync"
)

type aSync struct {
	wg    sync.WaitGroup
	funcs []func() error
}

func New() *aSync {
	return &aSync{
		wg:    sync.WaitGroup{},
		funcs: []func() error{},
	}
}

func (a *aSync) AddFunc(f func() error) {
	a.funcs = append(a.funcs, f)
}

func (a *aSync) AddFuncs(f ...func() error) {
	if len(f) > 0 {
		for i := range f {
			a.funcs = append(a.funcs, f[i])
		}
	}
}

func (a *aSync) Run() {
	for i := range a.funcs {
		a.wg.Add(1)
		go func(i int) {
			defer handlePanic()
			defer a.wg.Done()
			a.funcs[i]()
		}(i)
	}
	a.wg.Wait()
}

func handlePanic() {
	if err := recover(); err != nil {
		log.Println("panic: ", err)
		debug.PrintStack()
	}
}
