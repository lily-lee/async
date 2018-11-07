package async_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/lily-lee/async"
)

func TestAsync(t *testing.T) {
	i := 0
	a := async.New()
	a.AddFuncs(func() error {
		i++
		fmt.Println("a:::", i, time.Now().Unix())
		time.Sleep(3 * time.Second)
		fmt.Println("a done.")
		return nil
	})
	a.AddFuncs(func() error {
		i++
		fmt.Println("b:::", i, time.Now().Unix())
		return nil
	})
	a.AddFuncs(func() error {
		i++
		fmt.Println("c:::", i, time.Now().Unix())
		panic(errors.New("asf"))
		return nil
	})
	a.AddFuncs(func() error {
		panic("d panic errr 1")
		fmt.Println("d:::", i, time.Now().Unix())
		panic("d panic errr 2")
		return nil
	})
	a.AddFuncs(func() error {
		i++
		fmt.Println("e::", i)
		return nil
	}, func() error {
		i++
		fmt.Println("f::", i)
		return nil
	})

	a.Run()

	fmt.Println("i::::", i)
}
