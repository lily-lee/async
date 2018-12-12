package async_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/lily-lee/async"
)

func TestAsync(t *testing.T) {
	ss := []string{}
	m := struct {
		A string
		B string
		C string
		D string
	}{}
	s := time.Now()
	ii := 0
	a := async.New()
	a.AddFunc(func() error {
		ii++
		ss = append(ss, "a")
		m.A = "a"
		fmt.Println("a:::", ii, time.Now().UnixNano())
		time.Sleep(3 * time.Second)
		return errors.New("aaaaa")
	})
	a.AddFunc(func() error {
		ii++
		ss = append(ss, "b")
		m.B = "b"
		fmt.Println("b:::", ii, time.Now().UnixNano())
		return nil
	})
	a.AddFuncs(func() error {
		ii++
		ss = append(ss, "c")
		m.C = "c"
		fmt.Println("c:::", ii, time.Now().UnixNano())
		return nil
	})
	a.AddFuncs(func() error {
		ii++
		ss = append(ss, "d")
		m.D = "d"
		fmt.Println("d:::", ii, time.Now().UnixNano())
		return errors.New("ddddddd")
	}, func() error {
		ii++
		ss = append(ss, "e")
		fmt.Println("e:::", ii, time.Now().UnixNano())
		return errors.New("eeeeeeeee")
	}, func() error {
		ii++
		ss = append(ss, "f")
		fmt.Println("f:::", ii, time.Now().UnixNano())
		return errors.New("ffffffff")
	})

	e := a.Run()
	fmt.Println("err: ", e)

	fmt.Println("ii::::", ii, time.Now().Sub(s).String(), ss, m)

	a.AddFuncs(func() error {
		fmt.Println("AAA")
		return nil
	}, func() error {
		fmt.Println("BBB")
		return nil
	})
	a.Run()
}
