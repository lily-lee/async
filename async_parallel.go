package async

import (
	"sync"
)

// Parallel return parallel handler
func Parallel() *parallel {
	return &parallel{
		wg: sync.WaitGroup{},
		f:  []Func{},
		e:  errMap{mu: sync.Mutex{}, errors: map[string]error{}},
	}
}

type parallel struct {
	wg sync.WaitGroup
	f  []Func
	e  errMap
	eC chan err
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

// Run return all errors if error exists
func (p *parallel) Run() map[string]error {
	p.eC = make(chan err, len(p.f))

	for i := range p.f {
		p.wg.Add(1)
		go func(i int) {
			defer handlePanic()
			defer p.wg.Done()
			p.eC <- err{tag: p.f[i].Tag, err: p.f[i].F()}
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

func (p *parallel) reset() {
	p.f = []Func{}
	p.e = errMap{mu: sync.Mutex{}, errors: map[string]error{}}
}
