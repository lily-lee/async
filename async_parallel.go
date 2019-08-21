package async

import (
	"sync"
	"time"
)

// Parallel return parallel handler
func Parallel() *parallel {
	p := new(parallel)
	p.init()
	return p
}

type parallel struct {
	wg       *sync.WaitGroup
	f        []Func
	e        errMap
	eC       chan err
	timeout  time.Duration
	poolNum  int
	poolChan chan Func
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

func (p *parallel) init() {
	p.wg = &sync.WaitGroup{}
	p.f = []Func{}
	p.e = errMap{mu: sync.Mutex{}, errors: map[string]error{}}
	p.timeout = defaultTimeout
	p.poolNum = defaultPoolNum
	p.poolChan = make(chan Func)
}

// AddFunc add one func to handler
func (p *parallel) AddFunc(f Func) {
	p.f = append(p.f, f)
}

// AddFunc add more functions to handler
func (p *parallel) AddFuncs(f ...Func) {
	p.f = append(p.f, f...)
}

// SetTimeout ...
func (p *parallel) SetTimeout(timeout time.Duration) *parallel {
	p.timeout = timeout

	return p
}

// SetPoolNum ...
func (p *parallel) SetPoolNum(n int) {
	p.poolNum = n
}

// Run return all errors if error exists
func (p *parallel) Run() map[string]error {
	p.eC = make(chan err, len(p.f))

	if p.poolNum > len(p.f) {
		p.poolNum = len(p.f)
	}

	go func() {
		for i := range p.f {
			p.poolChan <- p.f[i]
		}
		close(p.poolChan)
	}()
	p.wg.Add(p.poolNum)
	for i := 0; i < p.poolNum; i++ {
		go func() {
			defer p.wg.Done()
			for f := range p.poolChan {
				p.eC <- err{
					tag: f.Tag,
					err: handle(p.timeout, f.F),
				}
			}
		}()
	}
	p.wg.Wait()
	close(p.eC)
	for e := range p.eC {
		p.e.errors[e.tag] = e.err
	}

	result := p.e.errors

	p.init()

	return result
}
