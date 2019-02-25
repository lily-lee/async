package async

import (
	"log"
	"runtime/debug"
	"sync"
	"time"
)

// New simple aSync handler
func New() *aSync {
	return &aSync{
		wg:      &sync.WaitGroup{},
		funcs:   []func() error{},
		timeout: defaultTimeout,
	}
}

type aSync struct {
	wg      *sync.WaitGroup
	funcs   []func() error
	timeout time.Duration
}

// AddFunc add one func to handler
func (a *aSync) AddFunc(f func() error) {
	a.funcs = append(a.funcs, f)
}

// AddFunc add more functions to handler
func (a *aSync) AddFuncs(f ...func() error) {
	a.funcs = append(a.funcs, f...)
}

// SetTimeout set timeout, default timeout is 5 minutes
func (a *aSync) SetTimeout(timeout time.Duration) *aSync {
	a.timeout = timeout

	return a
}

// Run return only one error if error exists
func (a *aSync) Run() error {
	errChan := make(chan error, len(a.funcs))
	for i := range a.funcs {
		a.wg.Add(1)
		go func(i int) {
			defer a.wg.Done()
			errChan <- handle(a.timeout, a.funcs[i])
		}(i)
	}

	a.wg.Wait()
	close(errChan)

	var err error
	for e := range errChan {
		if e != nil {
			err = e
			break
		}
	}

	a.reset()

	return err
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
