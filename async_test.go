package async_test

import (
	"testing"
	"time"

	"fmt"

	"errors"

	"github.com/lily-lee/async"
)

func TestAsync(t *testing.T) {
	t.Run("test timeout", func(t *testing.T) {
		start := time.Now()
		a := async.New()
		a.AddFunc(func() error {
			time.Sleep(3 * time.Second)
			return errors.New("ErrA")
		})

		a.AddFunc(func() error {
			time.Sleep(3 * time.Second)
			return errors.New("ErrB")
		})

		err := a.SetTimeout(3 * time.Second).Run()
		fmt.Println(err)
		//if err != async.TimeoutErr {
		//	t.Fail()
		//}

		fmt.Println(time.Now().Sub(start).Seconds())
	})

	t.Run("test common error", func(t *testing.T) {
		start := time.Now()
		a := async.New()
		a.AddFunc(func() error {
			time.Sleep(3 * time.Second)
			return errors.New("ErrA")
		})

		a.AddFunc(func() error {
			time.Sleep(3 * time.Second)
			return errors.New("ErrB")
		})

		err := a.SetTimeout(5 * time.Second).Run()
		fmt.Println(err)
		//if err != async.TimeoutErr {
		//	t.Fail()
		//}

		fmt.Println(time.Now().Sub(start).Seconds())
	})
}
