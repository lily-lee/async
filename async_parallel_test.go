package async

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestParallel(t *testing.T) {
	p := Parallel()
	p.AddFuncs(Func{Tag: "a", F: func() error {
		fmt.Println("a ", time.Now().UnixNano())
		return errors.New("a's error")
	}}, Func{Tag: "b", F: func() error {
		fmt.Println("b ", time.Now().UnixNano())
		return errors.New("b's error")
	}}, Func{Tag: "c", F: func() error {
		fmt.Println("c ", time.Now().UnixNano())
		return errors.New("c's error")
	}}, Func{Tag: "d", F: func() error {
		fmt.Println("d ", time.Now().UnixNano())
		time.Sleep(3 * time.Second)
		fmt.Println("d done")
		return nil
	}})

	err := p.Run()
	fmt.Printf("Parallel err: %+v\n", err)
}
