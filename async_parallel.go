package async

import (
	"sync"
)

type parallel struct {
	wg sync.WaitGroup
	f  []Func
	e  Error
	eC chan Err
}

type Error struct {
	mu     sync.Mutex
	errors map[string]error
}

type Func struct {
	Tag string
	F   func() error
}

func Parallel() *parallel {
	return &parallel{
		wg: sync.WaitGroup{},
		f:  []Func{},
		e:  Error{mu: sync.Mutex{}, errors: map[string]error{}},
	}
}

func (p *parallel) AddFunc(f Func) {
	p.f = append(p.f, f)
}

func (p *parallel) AddFuncs(f ...Func) {
	p.f = append(p.f, f...)
}

type Err struct {
	tag string
	err error
}

func (p *parallel) Run() map[string]error {
	p.eC = make(chan Err, len(p.f))

	for i := range p.f {
		p.wg.Add(1)
		go func(i int) {
			defer handlePanic()
			defer p.wg.Done()
			p.eC <- Err{tag: p.f[i].Tag, err: p.f[i].F()}
		}(i)
	}

	p.wg.Wait()

	close(p.eC)
	for e := range p.eC {
		p.e.errors[e.tag] = e.err
	}

	return p.e.errors
}
