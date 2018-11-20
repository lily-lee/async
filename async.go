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
	a.funcs = append(a.funcs, f...)
}

func (a *aSync) Run() error {
	errChan := make(chan error, len(a.funcs))
	for i := range a.funcs {
		a.wg.Add(1)
		go func(i int) {
			defer handlePanic()
			defer a.wg.Done()
			errChan <- a.funcs[i]()
		}(i)
	}
	a.wg.Wait()
	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
	}

	return nil
}

func handlePanic() {
	if err := recover(); err != nil {
		log.Println("panic: ", err)
		debug.PrintStack()
	}
}
