package async

import (
	"sync"
	"time"
)

// Parallel return parallel handler
func Parallel() *parallel {
	return &parallel{
		wg:      &sync.WaitGroup{},
		f:       []Func{},
		e:       errMap{mu: sync.Mutex{}, errors: map[string]error{}},
		timeout: defaultTimeout,
	}
}

type parallel struct {
	wg      *sync.WaitGroup
	f       []Func
	e       errMap
	eC      chan err
	timeout time.Duration
}

type Func struct {
	Tag string
	F   func() error
}

type errMap struct {
	mu     sync.Mutex
	errors map[string]error
}

type err struct {
	tag string
	err error
}

// AddFunc add one func to handler
func (p *parallel) AddFunc(f Func) {
	p.f = append(p.f, f)
}

// AddFunc add more functions to handler
func (p *parallel) AddFuncs(f ...Func) {
	p.f = append(p.f, f...)
}

func (p *parallel) SetTimeout(timeout time.Duration) *parallel {
	p.timeout = timeout

	return p
}

// Run return all errors if error exists
func (p *parallel) Run() map[string]error {
	p.eC = make(chan err, len(p.f))

	for i := range p.f {
		p.wg.Add(1)
		go func(i int) {
			defer p.wg.Done()
			p.eC <- err{
				tag: p.f[i].Tag,
				err: handle(p.timeout, p.f[i].F),
			}
		}(i)
	}

	p.wg.Wait()
	close(p.eC)

	for e := range p.eC {
		p.e.errors[e.tag] = e.err
	}

	result := p.e.errors

	p.reset()

	return result
}

func (p *parallel) handle(i int) {
	defer handlePanic()
	var timer *time.Timer
	if p.timeout > 0 {
		timer = time.NewTimer(p.timeout)
		defer timer.Stop()
		go p.handleTimeout(i, timer)
	}

	func(t *time.Timer) {
		defer p.wg.Done()
		p.eC <- err{tag: p.f[i].Tag, err: p.f[i].F()}
		t.Stop()
	}(timer)
}

func (p *parallel) handleTimeout(i int, timer *time.Timer) {
	defer handlePanic()
	if timer != nil {
		<-timer.C
		p.eC <- err{tag: p.f[i].Tag, err: TimeoutErr}
		p.wg.Done()
	}
}

func (p *parallel) reset() {
	p.f = []Func{}
	p.e = errMap{mu: sync.Mutex{}, errors: map[string]error{}}
	p.timeout = defaultTimeout
}

func handle(timeout time.Duration, f func() error) error {
	if timeout > 0 {
		timer := time.NewTimer(timeout)
		defer timer.Stop()
		eC := make(chan error)
		go func() {
			eC <- f()
		}()
		select {
		case <-timer.C:
			return TimeoutErr
		case err := <-eC:
			return err
		}
	}

	return f()
}
