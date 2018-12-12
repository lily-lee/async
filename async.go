package async

import (
	"log"
	"runtime/debug"
	"sync"
)

// New simple aSync handler
func New() *aSync {
	return &aSync{
		wg:    sync.WaitGroup{},
		funcs: []func() error{},
	}
}

type aSync struct {
	wg    sync.WaitGroup
	funcs []func() error
}

// AddFunc add one func to handler
func (a *aSync) AddFunc(f func() error) {
	a.funcs = append(a.funcs, f)
}

// AddFunc add more functions to handler
func (a *aSync) AddFuncs(f ...func() error) {
	a.funcs = append(a.funcs, f...)
}

// Run return only one error if error exists
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
			a.reset()
			return err
		}
	}
	close(errChan)
	a.reset()
	return nil
}

func (a *aSync) reset() {
	a.funcs = []func() error{}
}

func handlePanic() {
	if err := recover(); err != nil {
		log.Println("panic: ", err)
		debug.PrintStack()
	}
}
