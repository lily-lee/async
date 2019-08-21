package async

import (
	"sync"
	"time"
)

// New simple aSync handler
func New() *aSync {
	a := new(aSync)
	a.init()
	return a
}

type aSync struct {
	wg       *sync.WaitGroup
	funcs    []func() error
	timeout  time.Duration
	poolNum  int
	poolChan chan func() error
}

func (a *aSync) init() {
	a.wg = &sync.WaitGroup{}
	a.funcs = []func() error{}
	a.timeout = defaultTimeout
	a.poolNum = defaultPoolNum
	a.poolChan = make(chan func() error)
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

// SetPoolNum set num
func (a *aSync) SetPoolNum(n int) {
	a.poolNum = n
}

// Run return only one error if error exists
func (a *aSync) Run() error {
	errChan := make(chan error, len(a.funcs))
	if len(a.funcs) <= a.poolNum {
		a.poolNum = len(a.funcs)
	}

	go func() {
		for i := range a.funcs {
			a.poolChan <- a.funcs[i]
		}
		close(a.poolChan)
	}()

	a.wg.Add(a.poolNum)
	for i := 0; i < a.poolNum; i++ {
		go func() {
			defer a.wg.Done()
			for f := range a.poolChan {
				errChan <- handle(a.timeout, f)
			}
		}()
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

	a.init()

	return err
}
